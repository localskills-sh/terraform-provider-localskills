package sso_connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var (
	_ resource.Resource                = &SsoConnectionResource{}
	_ resource.ResourceWithImportState = &SsoConnectionResource{}
)

type SsoConnectionResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &SsoConnectionResource{}
}

func (r *SsoConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_connection"
}

func (r *SsoConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the SAML SSO connection for a team on localskills.sh. This is a **singleton resource** — each team has at most one SSO connection. Deleting this resource disables SSO rather than removing the configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the SSO connection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the team (tenant) this SSO connection belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the SSO connection.",
				Required:    true,
			},
			"metadata_url": schema.StringAttribute{
				Description: "The URL to the IdP metadata XML. At least one of metadata_url or metadata_xml must be provided.",
				Optional:    true,
				Validators: []validator.String{
					frameworkvalidator.AtLeastOneOf(path.MatchRoot("metadata_url"), path.MatchRoot("metadata_xml")),
				},
			},
			"metadata_xml": schema.StringAttribute{
				Description: "The raw IdP metadata XML. At least one of metadata_url or metadata_xml must be provided.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					frameworkvalidator.AtLeastOneOf(path.MatchRoot("metadata_url"), path.MatchRoot("metadata_xml")),
				},
			},
			"default_role": schema.StringAttribute{
				Description: "The default role assigned to users who sign in via SSO. Must be one of: admin, member, viewonly.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					frameworkvalidator.OneOf("admin", "member", "viewonly"),
				},
			},
			"email_domains": schema.ListAttribute{
				Description: "List of email domains that are allowed to use SSO.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the SSO connection is enabled. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"require_sso": schema.BoolAttribute{
				Description: "Whether SSO is required for all users. Defaults to false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"idp_entity_id": schema.StringAttribute{
				Description: "The Identity Provider entity ID.",
				Computed:    true,
			},
			"idp_sso_url": schema.StringAttribute{
				Description: "The Identity Provider SSO URL.",
				Computed:    true,
			},
			"idp_slo_url": schema.StringAttribute{
				Description: "The Identity Provider SLO (Single Logout) URL.",
				Computed:    true,
			},
			"sp_entity_id": schema.StringAttribute{
				Description: "The Service Provider entity ID.",
				Computed:    true,
			},
			"sp_acs_url": schema.StringAttribute{
				Description: "The Service Provider ACS (Assertion Consumer Service) URL.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the SSO connection was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The timestamp when the SSO connection was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *SsoConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *SsoConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SsoConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildUpdateRequest(ctx, &plan)

	conn, err := r.client.UpdateSSOConnection(ctx, plan.TenantID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSO connection", err.Error())
		return
	}

	// Preserve metadata_xml from plan since API may not return raw XML
	preservedMetadataXML := plan.MetadataXML
	mapConnectionToState(conn, &plan)
	plan.MetadataXML = preservedMetadataXML

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SsoConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SsoConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve metadata_xml from state since API may not return raw XML
	preservedMetadataXML := state.MetadataXML

	conn, err := r.client.GetSSOConnection(ctx, state.TenantID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SSO connection", err.Error())
		return
	}

	mapConnectionToState(conn, &state)
	state.MetadataXML = preservedMetadataXML

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SsoConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SsoConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildUpdateRequest(ctx, &plan)

	conn, err := r.client.UpdateSSOConnection(ctx, plan.TenantID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating SSO connection", err.Error())
		return
	}

	// Preserve metadata_xml from plan
	preservedMetadataXML := plan.MetadataXML
	mapConnectionToState(conn, &plan)
	plan.MetadataXML = preservedMetadataXML

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SsoConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SsoConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// SSO connections are singletons — "delete" disables them
	enabled := false
	requireSso := false
	_, err := r.client.UpdateSSOConnection(ctx, state.TenantID.ValueString(), client.UpdateSsoRequest{
		Enabled:    &enabled,
		RequireSso: &requireSso,
	})
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error disabling SSO connection", err.Error())
	}
}

func (r *SsoConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by tenant_id only — SSO is a singleton per tenant
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), req.ID)...)
}

func buildUpdateRequest(ctx context.Context, plan *SsoConnectionModel) client.UpdateSsoRequest {
	displayName := plan.DisplayName.ValueString()
	enabled := plan.Enabled.ValueBool()
	requireSso := plan.RequireSso.ValueBool()

	updateReq := client.UpdateSsoRequest{
		DisplayName: &displayName,
		Enabled:     &enabled,
		RequireSso:  &requireSso,
	}

	if !plan.MetadataURL.IsNull() && !plan.MetadataURL.IsUnknown() {
		metadataURL := plan.MetadataURL.ValueString()
		updateReq.MetadataURL = &metadataURL
	}
	if !plan.MetadataXML.IsNull() && !plan.MetadataXML.IsUnknown() {
		metadataXML := plan.MetadataXML.ValueString()
		updateReq.MetadataXML = &metadataXML
	}
	if !plan.DefaultRole.IsNull() && !plan.DefaultRole.IsUnknown() {
		defaultRole := plan.DefaultRole.ValueString()
		updateReq.DefaultRole = &defaultRole
	}
	if !plan.EmailDomains.IsNull() && !plan.EmailDomains.IsUnknown() {
		var emailDomains []string
		plan.EmailDomains.ElementsAs(ctx, &emailDomains, false)
		updateReq.EmailDomains = emailDomains
	}

	return updateReq
}

func mapConnectionToState(conn *client.SsoConnection, state *SsoConnectionModel) {
	state.ID = types.StringValue(conn.ID)
	state.TenantID = types.StringValue(conn.TenantID)
	state.DisplayName = types.StringValue(conn.DisplayName)
	state.DefaultRole = types.StringValue(conn.DefaultRole)
	state.Enabled = types.BoolValue(conn.Enabled)
	state.RequireSso = types.BoolValue(conn.RequireSso)
	state.IdpEntityID = types.StringValue(conn.IdpEntityID)
	state.IdpSsoURL = types.StringValue(conn.IdpSsoURL)
	state.IdpSloURL = types.StringValue(conn.IdpSloURL)
	state.SpEntityID = types.StringValue(conn.SpEntityID)
	state.SpAcsURL = types.StringValue(conn.SpAcsURL)
	state.CreatedAt = types.StringValue(conn.CreatedAt)
	state.UpdatedAt = types.StringValue(conn.UpdatedAt)

	if conn.MetadataURL != "" {
		state.MetadataURL = types.StringValue(conn.MetadataURL)
	} else {
		state.MetadataURL = types.StringNull()
	}

	if len(conn.EmailDomains) > 0 {
		emailDomainValues, _ := types.ListValueFrom(context.Background(), types.StringType, conn.EmailDomains)
		state.EmailDomains = emailDomainValues
	} else {
		state.EmailDomains = types.ListNull(types.StringType)
	}
}

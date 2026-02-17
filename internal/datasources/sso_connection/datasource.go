package sso_connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &ssoConnectionDataSource{}

type ssoConnectionDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &ssoConnectionDataSource{}
}

func (d *ssoConnectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_connection"
}

func (d *ssoConnectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the SSO connection configuration for a team.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "The tenant (team) ID.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the SSO connection.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the SSO connection.",
				Computed:    true,
			},
			"metadata_url": schema.StringAttribute{
				Description: "The URL to the IdP metadata XML.",
				Computed:    true,
			},
			"default_role": schema.StringAttribute{
				Description: "The default role assigned to SSO users.",
				Computed:    true,
			},
			"email_domains": schema.ListAttribute{
				Description: "List of email domains allowed for SSO.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the SSO connection is enabled.",
				Computed:    true,
			},
			"require_sso": schema.BoolAttribute{
				Description: "Whether SSO is required for all users.",
				Computed:    true,
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
				Description: "The Identity Provider SLO URL.",
				Computed:    true,
			},
			"sp_entity_id": schema.StringAttribute{
				Description: "The Service Provider entity ID.",
				Computed:    true,
			},
			"sp_acs_url": schema.StringAttribute{
				Description: "The Service Provider ACS URL.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When the SSO connection was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When the SSO connection was last updated.",
				Computed:    true,
			},
		},
	}
}

func (d *ssoConnectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *ssoConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config SsoConnectionDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn, err := d.client.GetSSOConnection(ctx, config.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSO connection", err.Error())
		return
	}

	var state SsoConnectionDataModel
	state.TenantID = config.TenantID
	state.ID = types.StringValue(conn.ID)
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
		emailDomainValues, diags := types.ListValueFrom(ctx, types.StringType, conn.EmailDomains)
		resp.Diagnostics.Append(diags...)
		state.EmailDomains = emailDomainValues
	} else {
		state.EmailDomains = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

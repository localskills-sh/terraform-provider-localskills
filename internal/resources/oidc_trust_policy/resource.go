package oidc_trust_policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var (
	_ resource.Resource                = &OidcTrustPolicyResource{}
	_ resource.ResourceWithImportState = &OidcTrustPolicyResource{}
)

type OidcTrustPolicyResource struct {
	client *client.Client
}

func NewResource() resource.Resource {
	return &OidcTrustPolicyResource{}
}

func (r *OidcTrustPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oidc_trust_policy"
}

func (r *OidcTrustPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OIDC trust policy for a team on localskills.sh. OIDC trust policies allow CI/CD pipelines (GitHub Actions, GitLab CI) to authenticate using OpenID Connect tokens and exchange them for short-lived localskills.sh API tokens.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the OIDC trust policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the team (tenant) this policy belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the OIDC trust policy.",
				Required:    true,
			},
			"oidc_provider": schema.StringAttribute{
				Description: "The OIDC provider. Must be one of: github, gitlab.",
				Required:    true,
				Validators: []validator.String{
					frameworkvalidator.OneOf("github", "gitlab"),
				},
			},
			"repository": schema.StringAttribute{
				Description: "The repository identifier (e.g., 'org/repo').",
				Required:    true,
			},
			"ref_filter": schema.StringAttribute{
				Description: "Git ref filter pattern. Defaults to '*' (all refs).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("*"),
			},
			"environment_filter": schema.StringAttribute{
				Description: "Environment filter for the policy.",
				Optional:    true,
			},
			"skill_ids": schema.ListAttribute{
				Description: "List of skill IDs that this policy grants access to.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the policy is enabled. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"created_by": schema.StringAttribute{
				Description: "The user who created the policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the policy was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The timestamp when the policy was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *OidcTrustPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OidcTrustPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OidcTrustPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var skillIDs []string
	if !plan.SkillIDs.IsNull() && !plan.SkillIDs.IsUnknown() {
		resp.Diagnostics.Append(plan.SkillIDs.ElementsAs(ctx, &skillIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createReq := client.CreateOidcPolicyRequest{
		Name:       plan.Name.ValueString(),
		Provider:   plan.Provider.ValueString(),
		Repository: plan.Repository.ValueString(),
		RefFilter:  plan.RefFilter.ValueString(),
		SkillIDs:   skillIDs,
		Enabled:    plan.Enabled.ValueBool(),
	}
	if !plan.EnvironmentFilter.IsNull() && !plan.EnvironmentFilter.IsUnknown() {
		envFilter := plan.EnvironmentFilter.ValueString()
		createReq.EnvironmentFilter = &envFilter
	}

	policy, err := r.client.CreateOIDCPolicy(ctx, plan.TenantID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OIDC trust policy", err.Error())
		return
	}

	mapPolicyToState(ctx, policy, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OidcTrustPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OidcTrustPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, err := r.client.ListOIDCPolicies(ctx, state.TenantID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading OIDC trust policies", err.Error())
		return
	}

	var found *client.OidcTrustPolicy
	for i := range policies {
		if policies[i].ID == state.ID.ValueString() {
			found = &policies[i]
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	mapPolicyToState(ctx, found, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OidcTrustPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OidcTrustPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OidcTrustPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	provider := plan.Provider.ValueString()
	repository := plan.Repository.ValueString()
	refFilter := plan.RefFilter.ValueString()
	enabled := plan.Enabled.ValueBool()

	updateReq := client.UpdateOidcPolicyRequest{
		Name:       &name,
		Provider:   &provider,
		Repository: &repository,
		RefFilter:  &refFilter,
		Enabled:    &enabled,
	}

	if !plan.EnvironmentFilter.IsNull() && !plan.EnvironmentFilter.IsUnknown() {
		envFilter := plan.EnvironmentFilter.ValueString()
		updateReq.EnvironmentFilter = &envFilter
	}

	if !plan.SkillIDs.IsNull() && !plan.SkillIDs.IsUnknown() {
		var skillIDs []string
		resp.Diagnostics.Append(plan.SkillIDs.ElementsAs(ctx, &skillIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.SkillIDs = skillIDs
	}

	policy, err := r.client.UpdateOIDCPolicy(ctx, plan.TenantID.ValueString(), state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating OIDC trust policy", err.Error())
		return
	}

	mapPolicyToState(ctx, policy, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OidcTrustPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OidcTrustPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOIDCPolicy(ctx, state.TenantID.ValueString(), state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting OIDC trust policy", err.Error())
	}
}

func (r *OidcTrustPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected import ID in the format: tenant_id/policy_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func mapPolicyToState(ctx context.Context, policy *client.OidcTrustPolicy, state *OidcTrustPolicyModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(policy.ID)
	state.TenantID = types.StringValue(policy.TenantID)
	state.Name = types.StringValue(policy.Name)
	state.Provider = types.StringValue(policy.Provider)
	state.Repository = types.StringValue(policy.Repository)
	state.RefFilter = types.StringValue(policy.RefFilter)
	state.Enabled = types.BoolValue(policy.Enabled)
	state.CreatedBy = types.StringValue(policy.CreatedBy)
	state.CreatedAt = types.StringValue(policy.CreatedAt)
	state.UpdatedAt = types.StringValue(policy.UpdatedAt)

	if policy.EnvironmentFilter != nil {
		state.EnvironmentFilter = types.StringValue(*policy.EnvironmentFilter)
	} else {
		state.EnvironmentFilter = types.StringNull()
	}

	if len(policy.SkillIDs) > 0 {
		skillIDValues := make([]attr.Value, len(policy.SkillIDs))
		for i, id := range policy.SkillIDs {
			skillIDValues[i] = types.StringValue(id)
		}
		skillIDsList, d := types.ListValue(types.StringType, skillIDValues)
		diags.Append(d...)
		state.SkillIDs = skillIDsList
	} else {
		state.SkillIDs = types.ListNull(types.StringType)
	}
}

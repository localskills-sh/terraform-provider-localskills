package oidc_trust_policies

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/localskills/terraform-provider-localskills/internal/client"
)

var _ datasource.DataSource = &oidcTrustPoliciesDataSource{}

type oidcTrustPoliciesDataSource struct {
	client *client.Client
}

func NewDataSource() datasource.DataSource {
	return &oidcTrustPoliciesDataSource{}
}

func (d *oidcTrustPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oidc_trust_policies"
}

func (d *oidcTrustPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all OIDC trust policies for a team.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "The tenant (team) ID.",
				Required:    true,
			},
			"policies": schema.ListNestedAttribute{
				Description: "List of OIDC trust policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the policy.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the policy.",
							Computed:    true,
						},
						"oidc_provider": schema.StringAttribute{
							Description: "The OIDC provider (github or gitlab).",
							Computed:    true,
						},
						"repository": schema.StringAttribute{
							Description: "The repository identifier.",
							Computed:    true,
						},
						"ref_filter": schema.StringAttribute{
							Description: "Git ref filter pattern.",
							Computed:    true,
						},
						"environment_filter": schema.StringAttribute{
							Description: "Environment filter.",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether the policy is enabled.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "The user who created the policy.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "When the policy was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "When the policy was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *oidcTrustPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *oidcTrustPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OidcTrustPoliciesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, err := d.client.ListOIDCPolicies(ctx, config.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading OIDC trust policies", err.Error())
		return
	}

	var state OidcTrustPoliciesModel
	state.TenantID = config.TenantID
	state.Policies = make([]OidcTrustPolicyItemModel, len(policies))
	for i, p := range policies {
		state.Policies[i] = OidcTrustPolicyItemModel{
			ID:         types.StringValue(p.ID),
			Name:       types.StringValue(p.Name),
			Provider:   types.StringValue(p.Provider),
			Repository: types.StringValue(p.Repository),
			RefFilter:  types.StringValue(p.RefFilter),
			Enabled:    types.BoolValue(p.Enabled),
			CreatedBy:  types.StringValue(p.CreatedBy),
			CreatedAt:  types.StringValue(p.CreatedAt),
			UpdatedAt:  types.StringValue(p.UpdatedAt),
		}
		if p.EnvironmentFilter != nil {
			state.Policies[i].EnvironmentFilter = types.StringValue(*p.EnvironmentFilter)
		} else {
			state.Policies[i].EnvironmentFilter = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

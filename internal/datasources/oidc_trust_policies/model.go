package oidc_trust_policies

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OidcTrustPoliciesModel struct {
	TenantID types.String               `tfsdk:"tenant_id"`
	Policies []OidcTrustPolicyItemModel `tfsdk:"policies"`
}

type OidcTrustPolicyItemModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Provider          types.String `tfsdk:"oidc_provider"`
	Repository        types.String `tfsdk:"repository"`
	RefFilter         types.String `tfsdk:"ref_filter"`
	EnvironmentFilter types.String `tfsdk:"environment_filter"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	CreatedBy         types.String `tfsdk:"created_by"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

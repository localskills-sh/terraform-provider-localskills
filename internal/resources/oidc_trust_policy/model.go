package oidc_trust_policy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OidcTrustPolicyModel struct {
	ID                types.String `tfsdk:"id"`
	TenantID          types.String `tfsdk:"tenant_id"`
	Name              types.String `tfsdk:"name"`
	Provider          types.String `tfsdk:"oidc_provider"`
	Repository        types.String `tfsdk:"repository"`
	RefFilter         types.String `tfsdk:"ref_filter"`
	EnvironmentFilter types.String `tfsdk:"environment_filter"`
	SkillIDs          types.List   `tfsdk:"skill_ids"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	CreatedBy         types.String `tfsdk:"created_by"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

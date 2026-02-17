package scim_tokens

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ScimTokensModel struct {
	TenantID types.String         `tfsdk:"tenant_id"`
	Tokens   []ScimTokenItemModel `tfsdk:"tokens"`
}

type ScimTokenItemModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	LastUsedAt types.String `tfsdk:"last_used_at"`
	ExpiresAt  types.String `tfsdk:"expires_at"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

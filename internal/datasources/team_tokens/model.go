package team_tokens

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamTokensModel struct {
	TenantID types.String         `tfsdk:"tenant_id"`
	Tokens   []TeamTokenItemModel `tfsdk:"tokens"`
}

type TeamTokenItemModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	LastUsedAt     types.String `tfsdk:"last_used_at"`
	ExpiresAt      types.String `tfsdk:"expires_at"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedByName  types.String `tfsdk:"created_by_name"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
}

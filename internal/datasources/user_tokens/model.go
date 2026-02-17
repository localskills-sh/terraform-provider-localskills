package user_tokens

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserTokensModel struct {
	Tokens []UserTokenItemModel `tfsdk:"tokens"`
}

type UserTokenItemModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	LastUsedAt types.String `tfsdk:"last_used_at"`
	ExpiresAt  types.String `tfsdk:"expires_at"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

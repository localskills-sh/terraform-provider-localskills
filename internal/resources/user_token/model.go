package user_token

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserTokenModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	TokenValue types.String `tfsdk:"token_value"`
	LastUsedAt types.String `tfsdk:"last_used_at"`
	ExpiresAt  types.String `tfsdk:"expires_at"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

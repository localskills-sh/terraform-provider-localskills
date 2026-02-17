package scim_token

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ScimTokenModel struct {
	ID            types.String `tfsdk:"id"`
	TenantID      types.String `tfsdk:"tenant_id"`
	Name          types.String `tfsdk:"name"`
	ExpiresInDays types.Int64  `tfsdk:"expires_in_days"`
	TokenValue    types.String `tfsdk:"token_value"`
	LastUsedAt    types.String `tfsdk:"last_used_at"`
	ExpiresAt     types.String `tfsdk:"expires_at"`
	CreatedAt     types.String `tfsdk:"created_at"`
}

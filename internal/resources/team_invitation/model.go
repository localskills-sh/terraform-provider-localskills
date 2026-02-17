package team_invitation

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamInvitationModel struct {
	ID         types.String `tfsdk:"id"`
	TenantID   types.String `tfsdk:"tenant_id"`
	Email      types.String `tfsdk:"email"`
	Role       types.String `tfsdk:"role"`
	Token      types.String `tfsdk:"token"`
	InvitedBy  types.String `tfsdk:"invited_by"`
	ExpiresAt  types.String `tfsdk:"expires_at"`
	AcceptedAt types.String `tfsdk:"accepted_at"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

package team_audit_log

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamAuditLogModel struct {
	TenantID types.String        `tfsdk:"tenant_id"`
	Page     types.Int64         `tfsdk:"page"`
	Limit    types.Int64         `tfsdk:"limit"`
	Action   types.String        `tfsdk:"action"`
	Total    types.Int64         `tfsdk:"total"`
	Entries  []AuditLogEntryModel `tfsdk:"entries"`
}

type AuditLogEntryModel struct {
	ID           types.String `tfsdk:"id"`
	Action       types.String `tfsdk:"action"`
	ActorID      types.String `tfsdk:"actor_id"`
	ActorName    types.String `tfsdk:"actor_name"`
	ActorImage   types.String `tfsdk:"actor_image"`
	ResourceType types.String `tfsdk:"resource_type"`
	ResourceID   types.String `tfsdk:"resource_id"`
	Metadata     types.String `tfsdk:"metadata"`
	CreatedAt    types.String `tfsdk:"created_at"`
}

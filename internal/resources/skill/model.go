package skill

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SkillModel struct {
	ID             types.String `tfsdk:"id"`
	PublicID       types.String `tfsdk:"public_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
	Name           types.String `tfsdk:"name"`
	Slug           types.String `tfsdk:"slug"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Visibility     types.String `tfsdk:"visibility"`
	Content        types.String `tfsdk:"content"`
	Tags           types.List   `tfsdk:"tags"`
	CurrentVersion types.Int64  `tfsdk:"current_version"`
	CurrentSemver  types.String `tfsdk:"current_semver"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

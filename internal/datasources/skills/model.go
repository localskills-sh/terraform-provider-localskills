package skills

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SkillsDataSourceModel struct {
	TenantID   types.String     `tfsdk:"tenant_id"`
	Visibility types.String     `tfsdk:"visibility"`
	Type       types.String     `tfsdk:"type"`
	Query      types.String     `tfsdk:"query"`
	Tag        types.String     `tfsdk:"tag"`
	Skills     []SkillItemModel `tfsdk:"skills"`
}

type SkillItemModel struct {
	ID             types.String `tfsdk:"id"`
	PublicID       types.String `tfsdk:"public_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
	Name           types.String `tfsdk:"name"`
	Slug           types.String `tfsdk:"slug"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Visibility     types.String `tfsdk:"visibility"`
	Tags           types.List   `tfsdk:"tags"`
	CurrentVersion types.Int64  `tfsdk:"current_version"`
	CurrentSemver  types.String `tfsdk:"current_semver"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

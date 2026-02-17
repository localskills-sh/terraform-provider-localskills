package skill_versions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SkillVersionsDataSourceModel struct {
	SkillID  types.String            `tfsdk:"skill_id"`
	Versions []SkillVersionItemModel `tfsdk:"versions"`
}

type SkillVersionItemModel struct {
	ID          types.String `tfsdk:"id"`
	SkillID     types.String `tfsdk:"skill_id"`
	Version     types.Int64  `tfsdk:"version"`
	Semver      types.String `tfsdk:"semver"`
	ContentHash types.String `tfsdk:"content_hash"`
	Message     types.String `tfsdk:"message"`
	Format      types.String `tfsdk:"format"`
	FileCount   types.Int64  `tfsdk:"file_count"`
	CreatedBy   types.String `tfsdk:"created_by"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

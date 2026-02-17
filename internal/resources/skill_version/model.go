package skill_version

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SkillVersionModel struct {
	ID          types.String `tfsdk:"id"`
	SkillID     types.String `tfsdk:"skill_id"`
	Version     types.Int64  `tfsdk:"version"`
	Semver      types.String `tfsdk:"semver"`
	Bump        types.String `tfsdk:"bump"`
	Content     types.String `tfsdk:"content"`
	Message     types.String `tfsdk:"message"`
	ContentHash types.String `tfsdk:"content_hash"`
	Format      types.String `tfsdk:"format"`
	FileCount   types.Int64  `tfsdk:"file_count"`
	CreatedBy   types.String `tfsdk:"created_by"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

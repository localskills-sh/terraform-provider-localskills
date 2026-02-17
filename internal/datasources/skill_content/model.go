package skill_content

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SkillContentDataSourceModel struct {
	SkillID      types.String `tfsdk:"skill_id"`
	Version      types.String `tfsdk:"version"`
	Semver       types.String `tfsdk:"semver"`
	Range        types.String `tfsdk:"range"`
	Content      types.String `tfsdk:"content"`
	Format       types.String `tfsdk:"format"`
	SkillName    types.String `tfsdk:"skill_name"`
	SkillSlug    types.String `tfsdk:"skill_slug"`
	FetchVersion types.Int64  `tfsdk:"fetch_version"`
	FetchSemver  types.String `tfsdk:"fetch_semver"`
}

package skill_manifest

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SkillManifestDataSourceModel struct {
	SkillID     types.String `tfsdk:"skill_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	Files       types.List   `tfsdk:"files"`
}

package explore

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExploreDataSourceModel struct {
	Query  types.String         `tfsdk:"query"`
	Tag    types.String         `tfsdk:"tag"`
	Type   types.String         `tfsdk:"type"`
	Sort   types.String         `tfsdk:"sort"`
	Skills []ExploreSkillModel  `tfsdk:"skills"`
}

type ExploreSkillModel struct {
	ID             types.String `tfsdk:"id"`
	PublicID       types.String `tfsdk:"public_id"`
	Name           types.String `tfsdk:"name"`
	Slug           types.String `tfsdk:"slug"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Tags           types.List   `tfsdk:"tags"`
	CurrentVersion types.Int64  `tfsdk:"current_version"`
	CurrentSemver  types.String `tfsdk:"current_semver"`
	AuthorName     types.String `tfsdk:"author_name"`
	AuthorUsername types.String `tfsdk:"author_username"`
	Downloads      types.Int64  `tfsdk:"downloads"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

package team

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamDataSourceModel struct {
	TeamID      types.String `tfsdk:"team_id"`
	Slug        types.String `tfsdk:"slug"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Role        types.String `tfsdk:"role"`
}

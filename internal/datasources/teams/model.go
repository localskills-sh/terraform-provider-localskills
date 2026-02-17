package teams

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamsDataSourceModel struct {
	Teams []TeamModel `tfsdk:"teams"`
}

type TeamModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	Role        types.String `tfsdk:"role"`
}

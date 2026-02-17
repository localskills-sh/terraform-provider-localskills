package user_profile

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserProfileModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Name     types.String `tfsdk:"name"`
	Email    types.String `tfsdk:"email"`
	Image    types.String `tfsdk:"image"`
	Bio      types.String `tfsdk:"bio"`
}

package user_profile_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccUserProfileDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserProfileDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_user_profile.test", "id"),
					resource.TestCheckResourceAttrSet("data.localskills_user_profile.test", "email"),
				),
			},
		},
	})
}

func testAccUserProfileDataSourceConfig() string {
	return `
data "localskills_user_profile" "test" {}
`
}

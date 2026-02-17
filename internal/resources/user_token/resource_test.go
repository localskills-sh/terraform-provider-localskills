package user_token_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccUserTokenResource_basic(t *testing.T) {
	name := testutils.RandomName("tf-test-token")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserTokenConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_user_token.test", "id"),
					resource.TestCheckResourceAttr("localskills_user_token.test", "name", name),
					resource.TestCheckResourceAttrSet("localskills_user_token.test", "token_value"),
					resource.TestCheckResourceAttrSet("localskills_user_token.test", "created_at"),
				),
			},
			{
				ResourceName:            "localskills_user_token.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token_value"},
			},
		},
	})
}

func testAccUserTokenConfig(name string) string {
	return `
resource "localskills_user_token" "test" {
  name = "` + name + `"
}
`
}

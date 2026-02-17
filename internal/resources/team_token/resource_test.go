package team_token_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccTeamTokenResource_basic(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}
	name := testutils.RandomName("tf-test-team-token")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamTokenConfig(tenantID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_team_token.test", "id"),
					resource.TestCheckResourceAttr("localskills_team_token.test", "tenant_id", tenantID),
					resource.TestCheckResourceAttr("localskills_team_token.test", "name", name),
					resource.TestCheckResourceAttrSet("localskills_team_token.test", "token_value"),
					resource.TestCheckResourceAttrSet("localskills_team_token.test", "created_at"),
				),
			},
			{
				ResourceName:            "localskills_team_token.test",
				ImportState:             true,
				ImportStateId:           tenantID + "/",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token_value", "expires_in_days"},
			},
		},
	})
}

func testAccTeamTokenConfig(tenantID, name string) string {
	return `
resource "localskills_team_token" "test" {
  tenant_id      = "` + tenantID + `"
  name           = "` + name + `"
  expires_in_days = 90
}
`
}

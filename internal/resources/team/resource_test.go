package team_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccTeamResource_basic(t *testing.T) {
	testutils.TestAccPreCheck(t)
	name := testutils.RandomName("tf-test-team")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_team.test", "id"),
					resource.TestCheckResourceAttr("localskills_team.test", "name", name),
					resource.TestCheckResourceAttrSet("localskills_team.test", "slug"),
					resource.TestCheckResourceAttrSet("localskills_team.test", "created_at"),
					resource.TestCheckResourceAttrSet("localskills_team.test", "updated_at"),
				),
			},
			{
				ResourceName:      "localskills_team.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTeamResource_update(t *testing.T) {
	testutils.TestAccPreCheck(t)
	name := testutils.RandomName("tf-test-team")
	updatedName := testutils.RandomName("tf-test-team-updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_team.test", "name", name),
				),
			},
			{
				Config: testAccTeamConfigWithDescription(updatedName, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_team.test", "name", updatedName),
					resource.TestCheckResourceAttr("localskills_team.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccTeamConfig(name string) string {
	return `
resource "localskills_team" "test" {
  name = "` + name + `"
}
`
}

func testAccTeamConfigWithDescription(name, description string) string {
	return `
resource "localskills_team" "test" {
  name        = "` + name + `"
  description = "` + description + `"
}
`
}

package team_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccTeamDataSource_byID(t *testing.T) {
	testutils.TestAccPreCheck(t)
	teamName := testutils.RandomName("tf-test-team")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceByIDConfig(teamName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_team.test", "id"),
					resource.TestCheckResourceAttr("data.localskills_team.test", "name", teamName),
					resource.TestCheckResourceAttrSet("data.localskills_team.test", "slug"),
					resource.TestCheckResourceAttrSet("data.localskills_team.test", "role"),
				),
			},
		},
	})
}

func TestAccTeamDataSource_bySlug(t *testing.T) {
	testutils.TestAccPreCheck(t)
	teamName := testutils.RandomName("tf-test-team")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceBySlugConfig(teamName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_team.test", "id"),
					resource.TestCheckResourceAttr("data.localskills_team.test", "name", teamName),
					resource.TestCheckResourceAttrSet("data.localskills_team.test", "slug"),
				),
			},
		},
	})
}

func testAccTeamDataSourceByIDConfig(name string) string {
	return `
resource "localskills_team" "test" {
  name = "` + name + `"
}

data "localskills_team" "test" {
  team_id = localskills_team.test.id
}
`
}

func testAccTeamDataSourceBySlugConfig(name string) string {
	return `
resource "localskills_team" "test" {
  name = "` + name + `"
}

data "localskills_team" "test" {
  slug = localskills_team.test.slug
}
`
}

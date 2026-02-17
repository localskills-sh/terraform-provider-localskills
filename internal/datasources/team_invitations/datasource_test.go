package team_invitations_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccTeamInvitationsDataSource_basic(t *testing.T) {
	testutils.TestAccPreCheck(t)
	teamName := testutils.RandomName("tf-test-team")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamInvitationsDataSourceConfig(teamName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_team_invitations.test", "tenant_id"),
					resource.TestCheckResourceAttrSet("data.localskills_team_invitations.test", "invitations.#"),
				),
			},
		},
	})
}

func testAccTeamInvitationsDataSourceConfig(teamName string) string {
	return `
resource "localskills_team" "test" {
  name = "` + teamName + `"
}

data "localskills_team_invitations" "test" {
  tenant_id = localskills_team.test.id
}
`
}

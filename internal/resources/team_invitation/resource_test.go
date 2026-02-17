package team_invitation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccTeamInvitationResource_basic(t *testing.T) {
	testutils.TestAccPreCheck(t)
	teamName := testutils.RandomName("tf-test-team")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamInvitationConfig(teamName, "invite-test@example.com", "member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_team_invitation.test", "id"),
					resource.TestCheckResourceAttrSet("localskills_team_invitation.test", "tenant_id"),
					resource.TestCheckResourceAttr("localskills_team_invitation.test", "email", "invite-test@example.com"),
					resource.TestCheckResourceAttr("localskills_team_invitation.test", "role", "member"),
					resource.TestCheckResourceAttrSet("localskills_team_invitation.test", "token"),
					resource.TestCheckResourceAttrSet("localskills_team_invitation.test", "invited_by"),
					resource.TestCheckResourceAttrSet("localskills_team_invitation.test", "expires_at"),
					resource.TestCheckResourceAttrSet("localskills_team_invitation.test", "created_at"),
				),
			},
			{
				ResourceName:            "localskills_team_invitation.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccTeamInvitationImportStateIDFunc("localskills_team_invitation.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func testAccTeamInvitationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["tenant_id"] + "/" + rs.Primary.Attributes["id"], nil
	}
}

func testAccTeamInvitationConfig(teamName, email, role string) string {
	return `
resource "localskills_team" "test" {
  name = "` + teamName + `"
}

resource "localskills_team_invitation" "test" {
  tenant_id = localskills_team.test.id
  email     = "` + email + `"
  role      = "` + role + `"
}
`
}

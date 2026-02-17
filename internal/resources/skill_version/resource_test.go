package skill_version_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccSkillVersionResource_basic(t *testing.T) {
	skillName := testutils.RandomName("tf-test-skill-ver")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillVersionResourceConfig(skillName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_skill_version.test", "id"),
					resource.TestCheckResourceAttrSet("localskills_skill_version.test", "skill_id"),
					resource.TestCheckResourceAttr("localskills_skill_version.test", "content", "# Updated content\nVersion 2."),
					resource.TestCheckResourceAttr("localskills_skill_version.test", "message", "update content"),
					resource.TestCheckResourceAttrSet("localskills_skill_version.test", "version"),
					resource.TestCheckResourceAttrSet("localskills_skill_version.test", "semver"),
					resource.TestCheckResourceAttrSet("localskills_skill_version.test", "created_at"),
				),
			},
		},
	})
}

func testAccSkillVersionResourceConfig(skillName string) string {
	return fmt.Sprintf(`
resource "localskills_skill" "test" {
  tenant_id  = "default"
  name       = %q
  type       = "skill"
  visibility = "private"
  content    = "# Initial content"
}

resource "localskills_skill_version" "test" {
  skill_id = localskills_skill.test.id
  content  = "# Updated content\nVersion 2."
  message  = "update content"
  bump     = "minor"
}
`, skillName)
}

package skill_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccSkillResource_basic(t *testing.T) {
	name := testutils.RandomName("tf-test-skill")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_skill.test", "id"),
					resource.TestCheckResourceAttr("localskills_skill.test", "name", name),
					resource.TestCheckResourceAttr("localskills_skill.test", "type", "skill"),
					resource.TestCheckResourceAttr("localskills_skill.test", "visibility", "private"),
					resource.TestCheckResourceAttr("localskills_skill.test", "content", "# Test Skill\nThis is a test."),
					resource.TestCheckResourceAttrSet("localskills_skill.test", "slug"),
					resource.TestCheckResourceAttrSet("localskills_skill.test", "created_at"),
				),
			},
		},
	})
}

func TestAccSkillResource_update(t *testing.T) {
	name := testutils.RandomName("tf-test-skill")
	updatedName := name + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_skill.test", "name", name),
					resource.TestCheckResourceAttr("localskills_skill.test", "visibility", "private"),
				),
			},
			{
				Config: testAccSkillResourceConfigUpdated(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_skill.test", "name", updatedName),
					resource.TestCheckResourceAttr("localskills_skill.test", "visibility", "unlisted"),
					resource.TestCheckResourceAttr("localskills_skill.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccSkillResource_import(t *testing.T) {
	name := testutils.RandomName("tf-test-skill")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillResourceConfig(name),
			},
			{
				ResourceName:            "localskills_skill.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content"},
			},
		},
	})
}

func testAccSkillResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "localskills_skill" "test" {
  tenant_id  = "default"
  name       = %q
  type       = "skill"
  visibility = "private"
  content    = "# Test Skill\nThis is a test."
}
`, name)
}

func testAccSkillResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "localskills_skill" "test" {
  tenant_id   = "default"
  name        = %q
  type        = "skill"
  visibility  = "unlisted"
  content     = "# Test Skill\nThis is a test."
  description = "Updated description"
}
`, name)
}

package skill_versions_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccSkillVersionsDataSource_basic(t *testing.T) {
	name := testutils.RandomName("tf-test-ds-versions")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillVersionsDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_skill_versions.test", "versions.#"),
				),
			},
		},
	})
}

func testAccSkillVersionsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "localskills_skill" "test" {
  tenant_id  = "default"
  name       = %q
  type       = "skill"
  visibility = "private"
  content    = "# Test Skill"
}

data "localskills_skill_versions" "test" {
  skill_id = localskills_skill.test.id
}
`, name)
}

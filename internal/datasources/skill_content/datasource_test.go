package skill_content_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccSkillContentDataSource_basic(t *testing.T) {
	name := testutils.RandomName("tf-test-ds-content")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillContentDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_skill_content.test", "content"),
					resource.TestCheckResourceAttrSet("data.localskills_skill_content.test", "format"),
					resource.TestCheckResourceAttrSet("data.localskills_skill_content.test", "fetch_version"),
				),
			},
		},
	})
}

func testAccSkillContentDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "localskills_skill" "test" {
  tenant_id  = "default"
  name       = %q
  type       = "skill"
  visibility = "private"
  content    = "# Test Skill Content"
}

data "localskills_skill_content" "test" {
  skill_id = localskills_skill.test.id
}
`, name)
}

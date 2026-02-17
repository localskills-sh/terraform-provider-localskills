package skill_analytics_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccSkillAnalyticsDataSource_basic(t *testing.T) {
	name := testutils.RandomName("tf-test-ds-analytics")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillAnalyticsDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_skill_analytics.test", "total_downloads"),
					resource.TestCheckResourceAttrSet("data.localskills_skill_analytics.test", "unique_users"),
					resource.TestCheckResourceAttrSet("data.localskills_skill_analytics.test", "unique_ips"),
				),
			},
		},
	})
}

func testAccSkillAnalyticsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "localskills_skill" "test" {
  tenant_id  = "default"
  name       = %q
  type       = "skill"
  visibility = "private"
  content    = "# Test Skill"
}

data "localskills_skill_analytics" "test" {
  skill_id = localskills_skill.test.id
}
`, name)
}

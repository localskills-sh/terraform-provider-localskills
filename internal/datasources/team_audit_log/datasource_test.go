package team_audit_log_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccTeamAuditLogDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamAuditLogDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_team_audit_log.test", "total"),
					resource.TestCheckResourceAttrSet("data.localskills_team_audit_log.test", "entries.#"),
				),
			},
		},
	})
}

func testAccTeamAuditLogDataSourceConfig() string {
	return `
data "localskills_team_audit_log" "test" {
  tenant_id = "default"
}
`
}

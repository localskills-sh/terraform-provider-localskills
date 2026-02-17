package user_audit_log_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccUserAuditLogDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserAuditLogDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_user_audit_log.test", "total"),
					resource.TestCheckResourceAttrSet("data.localskills_user_audit_log.test", "entries.#"),
				),
			},
		},
	})
}

func testAccUserAuditLogDataSourceConfig() string {
	return `
data "localskills_user_audit_log" "test" {}
`
}

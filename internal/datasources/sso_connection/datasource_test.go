package sso_connection_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccSsoConnectionDataSource_basic(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSsoConnectionDataSourceConfig(tenantID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_sso_connection.test", "id"),
					resource.TestCheckResourceAttrSet("data.localskills_sso_connection.test", "display_name"),
					resource.TestCheckResourceAttrSet("data.localskills_sso_connection.test", "sp_entity_id"),
					resource.TestCheckResourceAttrSet("data.localskills_sso_connection.test", "sp_acs_url"),
				),
			},
		},
	})
}

func testAccSsoConnectionDataSourceConfig(tenantID string) string {
	return `
data "localskills_sso_connection" "test" {
  tenant_id = "` + tenantID + `"
}
`
}

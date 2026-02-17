package scim_tokens_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccScimTokensDataSource_basic(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScimTokensDataSourceConfig(tenantID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_scim_tokens.test", "tokens.#"),
				),
			},
		},
	})
}

func testAccScimTokensDataSourceConfig(tenantID string) string {
	return `
data "localskills_scim_tokens" "test" {
  tenant_id = "` + tenantID + `"
}
`
}

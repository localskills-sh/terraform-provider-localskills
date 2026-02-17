package oidc_trust_policies_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccOidcTrustPoliciesDataSource_basic(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcTrustPoliciesDataSourceConfig(tenantID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.localskills_oidc_trust_policies.test", "policies.#"),
				),
			},
		},
	})
}

func testAccOidcTrustPoliciesDataSourceConfig(tenantID string) string {
	return `
data "localskills_oidc_trust_policies" "test" {
  tenant_id = "` + tenantID + `"
}
`
}

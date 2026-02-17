package sso_connection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccSsoConnection_basic(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}
	rName := testutils.RandomName("tf-sso")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSsoConnectionConfig(tenantID, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_sso_connection.test", "id"),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "tenant_id", tenantID),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "display_name", rName),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "enabled", "true"),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "require_sso", "false"),
					resource.TestCheckResourceAttrSet("localskills_sso_connection.test", "sp_entity_id"),
					resource.TestCheckResourceAttrSet("localskills_sso_connection.test", "sp_acs_url"),
					resource.TestCheckResourceAttrSet("localskills_sso_connection.test", "created_at"),
					resource.TestCheckResourceAttrSet("localskills_sso_connection.test", "updated_at"),
				),
			},
			{
				ResourceName:            "localskills_sso_connection.test",
				ImportState:             true,
				ImportStateId:           tenantID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata_url", "metadata_xml"},
			},
		},
	})
}

func TestAccSsoConnection_update(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}
	rName := testutils.RandomName("tf-sso")
	rNameUpdated := rName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSsoConnectionConfig(tenantID, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "display_name", rName),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "enabled", "true"),
				),
			},
			{
				Config: testAccSsoConnectionConfigUpdated(tenantID, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "display_name", rNameUpdated),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "default_role", "admin"),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "email_domains.#", "1"),
					resource.TestCheckResourceAttr("localskills_sso_connection.test", "email_domains.0", "example.com"),
				),
			},
		},
	})
}

func testAccSsoConnectionConfig(tenantID, name string) string {
	return fmt.Sprintf(`
resource "localskills_sso_connection" "test" {
  tenant_id    = %q
  display_name = %q
  metadata_url = "https://idp.example.com/metadata"
}
`, tenantID, name)
}

func testAccSsoConnectionConfigUpdated(tenantID, name string) string {
	return fmt.Sprintf(`
resource "localskills_sso_connection" "test" {
  tenant_id     = %q
  display_name  = %q
  metadata_url  = "https://idp.example.com/metadata"
  default_role  = "admin"
  email_domains = ["example.com"]
}
`, tenantID, name)
}

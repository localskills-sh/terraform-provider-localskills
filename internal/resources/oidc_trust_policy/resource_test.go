package oidc_trust_policy_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/localskills/terraform-provider-localskills/internal/testutils"
)

func TestAccOidcTrustPolicy_basic(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}
	rName := testutils.RandomName("tf-oidc")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcTrustPolicyConfig(tenantID, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("localskills_oidc_trust_policy.test", "id"),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "name", rName),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "oidc_provider", "github"),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "repository", "myorg/myrepo"),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "ref_filter", "*"),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("localskills_oidc_trust_policy.test", "created_at"),
					resource.TestCheckResourceAttrSet("localskills_oidc_trust_policy.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccOidcTrustPolicy_update(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}
	rName := testutils.RandomName("tf-oidc")
	rNameUpdated := rName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcTrustPolicyConfig(tenantID, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "name", rName),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "ref_filter", "*"),
				),
			},
			{
				Config: testAccOidcTrustPolicyConfigUpdated(tenantID, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "name", rNameUpdated),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "ref_filter", "refs/heads/main"),
					resource.TestCheckResourceAttr("localskills_oidc_trust_policy.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccOidcTrustPolicy_import(t *testing.T) {
	tenantID := os.Getenv("LOCALSKILLS_TENANT_ID")
	if tenantID == "" {
		t.Skip("LOCALSKILLS_TENANT_ID must be set for acceptance tests")
	}
	rName := testutils.RandomName("tf-oidc")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcTrustPolicyConfig(tenantID, rName),
			},
			{
				ResourceName:      "localskills_oidc_trust_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOidcTrustPolicyImportStateIdFunc,
			},
		},
	})
}

func testAccOidcTrustPolicyImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["localskills_oidc_trust_policy.test"]
	if !ok {
		return "", fmt.Errorf("resource not found: localskills_oidc_trust_policy.test")
	}
	return rs.Primary.Attributes["tenant_id"] + "/" + rs.Primary.Attributes["id"], nil
}

func testAccOidcTrustPolicyConfig(tenantID, name string) string {
	return fmt.Sprintf(`
resource "localskills_oidc_trust_policy" "test" {
  tenant_id  = %q
  name       = %q
  oidc_provider = "github"
  repository = "myorg/myrepo"
}
`, tenantID, name)
}

func testAccOidcTrustPolicyConfigUpdated(tenantID, name string) string {
	return fmt.Sprintf(`
resource "localskills_oidc_trust_policy" "test" {
  tenant_id  = %q
  name       = %q
  oidc_provider = "github"
  repository = "myorg/myrepo"
  ref_filter = "refs/heads/main"
  enabled    = false
}
`, tenantID, name)
}

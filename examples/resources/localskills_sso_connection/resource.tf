resource "localskills_sso_connection" "okta" {
  tenant_id     = localskills_team.engineering.id
  display_name  = "Okta SSO"
  metadata_url  = "https://myorg.okta.com/app/abc123/sso/saml/metadata"
  default_role  = "member"
  email_domains = ["example.com", "myorg.com"]
  enabled       = true
  require_sso   = false
}

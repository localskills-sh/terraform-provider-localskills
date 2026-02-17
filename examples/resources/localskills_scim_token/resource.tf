resource "localskills_scim_token" "okta_provisioning" {
  tenant_id       = localskills_team.engineering.id
  name            = "Okta SCIM Provisioning"
  expires_in_days = 365
}

output "scim_token" {
  value     = localskills_scim_token.okta_provisioning.token_value
  sensitive = true
}

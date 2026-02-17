resource "localskills_team_token" "ci" {
  tenant_id       = localskills_team.engineering.id
  name            = "CI Pipeline Token"
  expires_in_days = 90
}

# Store the token securely -- it is only available at creation
output "ci_token" {
  value     = localskills_team_token.ci.token_value
  sensitive = true
}

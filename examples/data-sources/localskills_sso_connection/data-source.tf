data "localskills_sso_connection" "current" {
  tenant_id = localskills_team.engineering.id
}

output "sso_enabled" {
  value = data.localskills_sso_connection.current.enabled
}

output "sp_acs_url" {
  value = data.localskills_sso_connection.current.sp_acs_url
}

data "localskills_scim_tokens" "all" {
  tenant_id = localskills_team.engineering.id
}

output "scim_token_names" {
  value = [for t in data.localskills_scim_tokens.all.tokens : t.name]
}

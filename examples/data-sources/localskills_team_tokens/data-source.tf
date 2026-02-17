data "localskills_team_tokens" "all" {
  tenant_id = localskills_team.engineering.id
}

output "token_names" {
  value = [for t in data.localskills_team_tokens.all.tokens : t.name]
}

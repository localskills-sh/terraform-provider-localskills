resource "localskills_team" "engineering" {
  name        = "Engineering"
  description = "Platform engineering team"
}

output "team_slug" {
  value = localskills_team.engineering.slug
}

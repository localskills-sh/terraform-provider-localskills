# Look up a team by slug
data "localskills_team" "engineering" {
  slug = "engineering"
}

output "team_name" {
  value = data.localskills_team.engineering.name
}

output "my_role" {
  value = data.localskills_team.engineering.role
}

resource "localskills_team_invitation" "alice" {
  tenant_id = localskills_team.engineering.id
  email     = "alice@example.com"
  role      = "member"
}

resource "localskills_team_invitation" "bob_admin" {
  tenant_id = localskills_team.engineering.id
  email     = "bob@example.com"
  role      = "admin"
}

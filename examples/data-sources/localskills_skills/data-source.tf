data "localskills_skills" "public_rules" {
  tenant_id  = localskills_team.engineering.id
  visibility = "public"
  type       = "rule"
}

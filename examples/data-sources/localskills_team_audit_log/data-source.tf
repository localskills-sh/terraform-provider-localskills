# Fetch team audit log entries
data "localskills_team_audit_log" "recent" {
  tenant_id = localskills_team.engineering.id
}

# Filter by action type
data "localskills_team_audit_log" "member_changes" {
  tenant_id = localskills_team.engineering.id
  action    = "team.member.add"
  limit     = 20
}

output "team_audit_total" {
  value = data.localskills_team_audit_log.recent.total
}

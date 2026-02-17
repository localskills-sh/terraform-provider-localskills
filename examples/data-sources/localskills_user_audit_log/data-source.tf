# Fetch recent audit log entries
data "localskills_user_audit_log" "recent" {}

# Filter by action type
data "localskills_user_audit_log" "skill_creates" {
  action = "skill.create"
  limit  = 10
}

output "total_entries" {
  value = data.localskills_user_audit_log.recent.total
}

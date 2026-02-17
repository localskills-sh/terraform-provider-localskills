data "localskills_skill_analytics" "eslint" {
  skill_id = localskills_skill.eslint_rules.id
}

output "total_downloads" {
  value = data.localskills_skill_analytics.eslint.total_downloads
}

output "unique_users" {
  value = data.localskills_skill_analytics.eslint.unique_users
}

# Get content at a specific semver range
data "localskills_skill_content" "latest_v1" {
  skill_id = "sk_abc123"
  range    = "^1.0.0"
}

output "content" {
  value = data.localskills_skill_content.latest_v1.content
}

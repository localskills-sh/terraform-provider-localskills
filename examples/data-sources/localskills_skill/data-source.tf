data "localskills_skill" "example" {
  skill_id = "sk_abc123"
}

output "skill_name" {
  value = data.localskills_skill.example.name
}

# Publish a new patch version
resource "localskills_skill_version" "v1_1" {
  skill_id = localskills_skill.eslint_rules.id
  content  = file("${path.module}/rules/eslint-v1.1.md")
  message  = "Added rules for React hooks"
  bump     = "minor"
}

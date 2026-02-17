data "localskills_skill_manifest" "package" {
  skill_id = localskills_skill.eslint_rules.id
}

output "manifest_files" {
  value = data.localskills_skill_manifest.package.files
}

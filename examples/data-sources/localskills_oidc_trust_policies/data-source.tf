data "localskills_oidc_trust_policies" "all" {
  tenant_id = localskills_team.engineering.id
}

output "policy_names" {
  value = [for p in data.localskills_oidc_trust_policies.all.policies : p.name]
}

# Allow GitHub Actions from main branch to publish skills
resource "localskills_oidc_trust_policy" "github_deploy" {
  tenant_id    = localskills_team.engineering.id
  name         = "GitHub Deploy"
  oidc_provider     = "github"
  repository   = "myorg/skills-repo"
  ref_filter   = "refs/heads/main"
  skill_ids    = [localskills_skill.eslint_rules.id]
  enabled      = true
}

# Allow all branches for staging
resource "localskills_oidc_trust_policy" "github_staging" {
  tenant_id          = localskills_team.engineering.id
  name               = "GitHub Staging"
  oidc_provider           = "github"
  repository         = "myorg/skills-repo"
  ref_filter         = "*"
  environment_filter = "staging"
  enabled            = true
}

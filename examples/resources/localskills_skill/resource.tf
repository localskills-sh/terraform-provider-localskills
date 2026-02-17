# Create a public skill
resource "localskills_skill" "eslint_rules" {
  tenant_id   = localskills_team.engineering.id
  name        = "ESLint Rules"
  description = "Standard ESLint configuration for all TypeScript projects."
  type        = "rule"
  visibility  = "public"
  content     = file("${path.module}/rules/eslint.md")
  tags        = ["linting", "typescript", "standards"]
}

# Create a private team skill
resource "localskills_skill" "deploy_guide" {
  tenant_id   = localskills_team.engineering.id
  name        = "Deployment Guide"
  description = "Internal deployment procedures and checklists."
  type        = "skill"
  visibility  = "private"
  content     = "# Deployment Guide\n\nFollow these steps..."
}

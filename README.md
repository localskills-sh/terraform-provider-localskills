# Terraform Provider for Localskills

[![Terraform Registry](https://img.shields.io/badge/terraform-registry-blueviolet)](https://registry.terraform.io/providers/localskills/localskills/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/localskills/terraform-provider-localskills)](https://goreportcard.com/report/github.com/localskills/terraform-provider-localskills)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

The Localskills Terraform provider manages resources on [localskills.sh](https://localskills.sh), a multi-tenant skill sharing platform. It enables teams to define skills and rules, publish versioned content, manage team membership, configure authentication tokens, set up OIDC trust policies for CI/CD pipelines, and manage enterprise SAML SSO and SCIM integrations -- all as Terraform-managed infrastructure.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22 (for building from source)

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    localskills = {
      source  = "localskills/localskills"
      version = "~> 0.1"
    }
  }
}
```

## Authentication

The provider authenticates using API tokens issued by localskills.sh. Tokens use the `lsk_` prefix format.

**Environment variable** (recommended):

```sh
export LOCALSKILLS_API_TOKEN="lsk_your_token_here"
```

**Provider configuration** (not recommended for version-controlled files):

```hcl
provider "localskills" {
  api_token = var.localskills_token
}
```

The `base_url` defaults to `https://localskills.sh` and can be overridden via the `LOCALSKILLS_BASE_URL` environment variable or the provider `base_url` attribute.

## Quick Start

```hcl
provider "localskills" {}

# Create a team
resource "localskills_team" "engineering" {
  name        = "engineering"
  description = "Engineering team"
}

# Create a skill owned by the team
resource "localskills_skill" "code_review" {
  tenant_id   = localskills_team.engineering.id
  name        = "code-review-guidelines"
  type        = "skill"
  visibility  = "private"
  content     = file("${path.module}/skills/code-review.md")
  tags        = ["engineering", "process"]
}

# Publish a new version
resource "localskills_skill_version" "v1" {
  skill_id = localskills_skill.code_review.id
  content  = file("${path.module}/skills/code-review-v2.md")
  bump     = "minor"
  message  = "Add section on async reviews"
}
```

## Resources

| Resource | Description |
|---|---|
| [`localskills_skill`](docs/resources/skill.md) | Manages a skill or rule with content, visibility, and tags |
| [`localskills_skill_version`](docs/resources/skill_version.md) | Creates immutable versioned snapshots of skill content |
| [`localskills_team`](docs/resources/team.md) | Manages a team (tenant) on the platform |
| [`localskills_team_invitation`](docs/resources/team_invitation.md) | Sends an invitation to join a team |
| [`localskills_team_token`](docs/resources/team_token.md) | Manages team-scoped API tokens |
| [`localskills_user_token`](docs/resources/user_token.md) | Manages user-scoped API tokens |
| [`localskills_oidc_trust_policy`](docs/resources/oidc_trust_policy.md) | Configures OIDC trust policies for CI/CD token exchange |
| [`localskills_sso_connection`](docs/resources/sso_connection.md) | Manages the SAML SSO connection for a team |
| [`localskills_scim_token`](docs/resources/scim_token.md) | Manages SCIM provisioning tokens for identity providers |

## Data Sources

| Data Source | Description |
|---|---|
| [`localskills_skill`](docs/data-sources/skill.md) | Reads a single skill by ID |
| [`localskills_skills`](docs/data-sources/skills.md) | Lists skills for a team |
| [`localskills_skill_versions`](docs/data-sources/skill_versions.md) | Lists all versions of a skill |
| [`localskills_skill_content`](docs/data-sources/skill_content.md) | Reads the content of a specific skill version |
| [`localskills_skill_analytics`](docs/data-sources/skill_analytics.md) | Reads download and view analytics for a skill |
| [`localskills_skill_manifest`](docs/data-sources/skill_manifest.md) | Reads the package manifest of a skill |
| [`localskills_explore`](docs/data-sources/explore.md) | Queries the public skill explore feed |
| [`localskills_team`](docs/data-sources/team.md) | Reads a single team by ID or slug |
| [`localskills_teams`](docs/data-sources/teams.md) | Lists all teams the authenticated user belongs to |
| [`localskills_team_invitations`](docs/data-sources/team_invitations.md) | Lists pending invitations for a team |
| [`localskills_user_tokens`](docs/data-sources/user_tokens.md) | Lists API tokens for the authenticated user |
| [`localskills_team_tokens`](docs/data-sources/team_tokens.md) | Lists API tokens for a team |
| [`localskills_oidc_trust_policies`](docs/data-sources/oidc_trust_policies.md) | Lists OIDC trust policies for a team |
| [`localskills_sso_connection`](docs/data-sources/sso_connection.md) | Reads the SSO connection for a team |
| [`localskills_scim_tokens`](docs/data-sources/scim_tokens.md) | Lists SCIM provisioning tokens for a team |
| [`localskills_user_profile`](docs/data-sources/user_profile.md) | Reads the authenticated user's profile |
| [`localskills_user_audit_log`](docs/data-sources/user_audit_log.md) | Reads audit log entries for the authenticated user |
| [`localskills_team_audit_log`](docs/data-sources/team_audit_log.md) | Reads audit log entries for a team |

## Development

### Building

```sh
make build
```

### Running Tests

Unit tests:

```sh
make test
```

Acceptance tests (requires `LOCALSKILLS_API_TOKEN`):

```sh
make testacc
```

### Generating Documentation

```sh
make generate
```

This runs `tfplugindocs` to regenerate the `docs/` directory from schema definitions, templates, and examples.

### Linting

```sh
make lint
```

### Project Structure

```
.
├── main.go                    # Provider entry point
├── internal/
│   ├── provider/              # Provider configuration and registration
│   ├── client/                # HTTP client, models, and API methods
│   ├── resources/             # Terraform resource implementations
│   │   ├── skill/
│   │   ├── skill_version/
│   │   ├── team/
│   │   ├── team_invitation/
│   │   ├── team_token/
│   │   ├── user_token/
│   │   ├── oidc_trust_policy/
│   │   ├── sso_connection/
│   │   └── scim_token/
│   ├── datasources/           # Terraform data source implementations
│   └── testutils/             # Shared test helpers
├── templates/                 # tfplugindocs templates
├── examples/                  # Example Terraform configurations
└── docs/                      # Generated documentation (do not edit)
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

MPL-2.0 -- see [LICENSE](LICENSE) for details.

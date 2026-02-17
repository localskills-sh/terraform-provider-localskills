# Terraform Provider for Localskills

The Localskills Terraform provider manages resources on the [localskills.sh](https://localskills.sh) skill sharing platform. It enables teams to define skills, manage team membership, configure authentication tokens, set up OIDC trust policies for CI/CD pipelines, and manage enterprise SSO/SCIM integrations as code.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22 (for development only)

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    localskills = {
      source = "localskills/localskills"
    }
  }
}
```

## Authentication

The provider authenticates using API tokens issued by localskills.sh. Tokens use the `lsk_` prefix format.

You can provide your token in two ways:

1. **Environment variable** (recommended):

   ```sh
   export LOCALSKILLS_API_TOKEN="lsk_your_token_here"
   ```

2. **Provider configuration**:

   ```hcl
   provider "localskills" {
     api_token = "lsk_your_token_here"
   }
   ```

The `base_url` attribute defaults to `https://localskills.sh` and can be overridden via the `LOCALSKILLS_BASE_URL` environment variable or the provider configuration.

## Quick Start

```hcl
terraform {
  required_providers {
    localskills = {
      source = "localskills/localskills"
    }
  }
}

provider "localskills" {
  # api_token is read from LOCALSKILLS_API_TOKEN env var
}

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
  content     = "# Code Review Guidelines\n\nReview all PRs within 24 hours."
  tags        = ["engineering", "process"]
}
```

## Resources

| Resource | Description |
|---|---|
| `localskills_skill` | Manages a skill (or rule) with content, visibility, and tags |
| `localskills_skill_version` | Creates immutable versioned snapshots of skill content |
| `localskills_team` | Manages a team (tenant) on the platform |
| `localskills_team_invitation` | Sends an invitation to join a team |
| `localskills_team_token` | Manages team-scoped API tokens |
| `localskills_user_token` | Manages user-scoped API tokens |
| `localskills_oidc_trust_policy` | Configures OIDC trust policies for CI/CD token exchange |
| `localskills_sso_connection` | Manages the SAML SSO connection for a team |
| `localskills_scim_token` | Manages SCIM provisioning tokens for identity providers |

## Data Sources

| Data Source | Description |
|---|---|
| `localskills_skill` | Reads a single skill by ID |
| `localskills_skills` | Lists skills for a team |
| `localskills_skill_versions` | Lists all versions of a skill |
| `localskills_skill_content` | Reads the content of a specific skill version |
| `localskills_skill_analytics` | Reads analytics data for a skill |
| `localskills_skill_manifest` | Reads the manifest of a skill |
| `localskills_explore` | Queries the public skill explore feed |
| `localskills_team` | Reads a single team by ID |
| `localskills_teams` | Lists all teams the authenticated user belongs to |
| `localskills_team_invitations` | Lists pending invitations for a team |
| `localskills_user_tokens` | Lists all API tokens for the authenticated user |
| `localskills_team_tokens` | Lists all API tokens for a team |
| `localskills_oidc_trust_policies` | Lists OIDC trust policies for a team |
| `localskills_sso_connection` | Reads the SSO connection for a team |
| `localskills_scim_tokens` | Lists SCIM tokens for a team |
| `localskills_user_profile` | Reads the authenticated user's profile |
| `localskills_user_audit_log` | Reads audit log entries for the authenticated user |
| `localskills_team_audit_log` | Reads audit log entries for a team |

## Development

### Building

```sh
go build -o terraform-provider-localskills
```

### Testing

Run unit tests:

```sh
make test
```

Run acceptance tests (requires a valid API token):

```sh
make testacc
```

### Generating Documentation

```sh
make generate
```

## License

MPL-2.0 -- see [LICENSE](LICENSE) for details.

# Contributing to terraform-provider-localskills

Thank you for your interest in contributing to the Localskills Terraform provider. This document covers the process for submitting contributions.

## Getting Started

1. Fork the repository and clone your fork
2. Install dependencies:
   - [Go](https://golang.org/doc/install) >= 1.22
   - [Terraform](https://www.terraform.io/downloads.html) >= 1.0
3. Build the provider:
   ```sh
   make build
   ```
4. Run the tests:
   ```sh
   make test
   ```

## Development Workflow

### Branch Naming

Use descriptive branch names:

- `feat/add-webhook-resource` -- new features
- `fix/token-refresh-error` -- bug fixes
- `docs/update-skill-examples` -- documentation changes

### Making Changes

1. Create a branch from `main`
2. Make your changes following the existing code patterns
3. Add or update tests for any changed behavior
4. Run linting and tests:
   ```sh
   make lint
   make test
   ```
5. Regenerate documentation if you changed schemas or examples:
   ```sh
   make generate
   ```
6. Commit your changes and open a pull request

### Code Organization

Each resource and data source lives in its own package under `internal/resources/` or `internal/datasources/` with three files:

- `model.go` -- Terraform state model struct with `tfsdk` tags
- `resource.go` / `datasource.go` -- CRUD operations and schema definition
- `resource_test.go` / `datasource_test.go` -- acceptance tests

API client methods live in `internal/client/` and shared data models are in `internal/client/models.go`.

### Adding a New Resource

1. Create a new package under `internal/resources/<resource_name>/`
2. Define the model in `model.go`
3. Implement the resource in `resource.go` (Schema, Create, Read, Update, Delete, ImportState)
4. Add acceptance tests in `resource_test.go`
5. Add the corresponding client methods in `internal/client/`
6. Register the resource in `internal/provider/provider.go`
7. Create an example at `examples/resources/localskills_<name>/resource.tf`
8. Create a template at `templates/resources/<name>.md.tmpl`
9. Run `make generate` to produce the documentation

### Adding a New Data Source

Follow the same pattern as resources but under `internal/datasources/` and with `examples/data-sources/` and `templates/data-sources/`.

## Testing

### Unit Tests

Unit tests use `httptest.NewServer` to mock API responses. They run without any external dependencies:

```sh
make test
```

### Acceptance Tests

Acceptance tests run against the real localskills.sh API and require:

- `LOCALSKILLS_API_TOKEN` -- a valid API token
- `LOCALSKILLS_TENANT_ID` -- a team ID (for team-scoped resources)

```sh
make testacc
```

Acceptance tests are skipped automatically when the required environment variables are not set.

### Writing Tests

- Every resource and data source must have test coverage
- Use `testutils.RandomName()` to generate unique resource names
- Use `testutils.TestAccPreCheck(t)` in the `PreCheck` function
- Use `testutils.TestAccProtoV6ProviderFactories` for the provider factory

## Pull Requests

### Before Submitting

- [ ] Code compiles (`make build`)
- [ ] Linting passes (`make lint`)
- [ ] Unit tests pass (`make test`)
- [ ] Documentation is regenerated if schemas or examples changed (`make generate`)
- [ ] New resources/data sources have acceptance tests
- [ ] Commit messages are clear and descriptive

### PR Guidelines

- Keep PRs focused on a single change
- Include a description of what changed and why
- Link any related issues
- Add examples for new resources or data sources

## Reporting Issues

Open an issue on GitHub with:

- A clear description of the problem
- Steps to reproduce (Terraform config if applicable)
- Expected vs actual behavior
- Provider version and Terraform version

## License

By contributing, you agree that your contributions will be licensed under the [MPL-2.0 License](LICENSE).

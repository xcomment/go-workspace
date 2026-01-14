# Copilot Instructions for MongoDB Go Driver (v1.17.x)

## Project Overview
- This is the official MongoDB driver for Go, supporting MongoDB 3.6+ and Go 1.18+ (Go 1.22+ required for tests).
- The codebase is modular, with key packages: `mongo`, `bson`, `internal`, and `examples`.
- Most development occurs in `mongo/` (driver logic), `bson/` (serialization), and `internal/` (utilities, assertions).
- Example usage and integration patterns are in `examples/documentation_examples/`.

## Build, Test, and Lint Workflows
- **Build:**
  - `make build` (cross-compiles for multiple Linux architectures)
  - `make build-tests` (compiles all tests, does not run them)
  - `make build-compile-check` (runs `etc/compile_check.sh` for build validation)
- **Test:**
  - `make test` (runs all tests, requires local MongoDB on port 27017)
  - `make test-race`, `make test-cover`, `make test-short` for race detection, coverage, and short tests
  - Specialized targets for FaaS, AWS Lambda, and Atlas Data Lake (see Makefile)
- **Lint:**
  - `make lint` (requires `golangci-lint` and `lll`)
  - `make fmt` and `make check-fmt` for formatting
- **Modules:**
  - `make check-modules` ensures `go.mod` and `vendor/` are up-to-date

## Key Conventions and Patterns
- **Testing:**
  - Tests use Go's `testing` package, with custom assertions in `internal/assert`.
  - Integration tests require a running MongoDB instance.
  - Many tests use `mtest` for integration scenarios.
- **Assertions:**
  - Custom assertion helpers are in `internal/assert/` (adapted from `stretchr/testify`).
- **API Operations:**
  - Operation names are constants in `internal/driverutil/operation.go` (source: MongoDB command reference).
- **Examples:**
  - Aggregation and usage examples are in `examples/documentation_examples/examples.go`.
  - Example pipelines and data flows are provided for reference.

## External Dependencies
- Uses Go modules for dependency management.
- Requires `golangci-lint`, `lll`, and optionally `pre-commit` for linting and formatting.
- FaaS and AWS Lambda tests require Docker and AWS SAM CLI.

## Integration Points
- FaaS and cloud integrations are in `internal/test/faas/awslambda/`.
- Specialized test runners and entry points in `cmd/` (e.g., `testoidcauth`, `testentauth`).

## Contribution and CI
- See `docs/CONTRIBUTING.md` for detailed contribution and CI instructions.
- PRs should reference Jira tickets (GODRIVER-xxx) when applicable.
- Linting is enforced via pre-commit hooks and Makefile targets.

## Quick Start for AI Agents
- Use Makefile targets for all build, test, and lint operations.
- Reference `examples/documentation_examples/` for idiomatic driver usage.
- Follow custom assertion patterns from `internal/assert/` in tests.
- Ensure all new code is covered by tests and passes lint/format checks.
- For integration tests, ensure a local MongoDB instance is running.

---

For more details, see:
- [README.md](../pkg/mod/go.mongodb.org/mongo-driver@v1.17.4/README.md)
- [CONTRIBUTING.md](../pkg/mod/go.mongodb.org/mongo-driver@v1.17.4/docs/CONTRIBUTING.md)
- [Makefile](../pkg/mod/go.mongodb.org/mongo-driver@v1.17.4/Makefile)
- [Examples](../pkg/mod/go.mongodb.org/mongo-driver@v1.17.4/examples/documentation_examples/examples.go)
- [Assertions](../pkg/mod/go.mongodb.org/mongo-driver@v1.17.4/internal/assert/assertions.go)

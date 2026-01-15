# Copilot Instructions for MongoDB Go Driver (v1.17.x)

## Project Overview

## Build, Test, and Lint Workflows
  - `make build` (cross-compiles for multiple Linux architectures)
  - `make build-tests` (compiles all tests, does not run them)
  - `make build-compile-check` (runs `etc/compile_check.sh` for build validation)
  - `make test` (runs all tests, requires local MongoDB on port 27017)
  - `make test-race`, `make test-cover`, `make test-short` for race detection, coverage, and short tests
  - Specialized targets for FaaS, AWS Lambda, and Atlas Data Lake (see Makefile)
  - `make lint` (requires `golangci-lint` and `lll`)
  - `make fmt` and `make check-fmt` for formatting
  - `make check-modules` ensures `go.mod` and `vendor/` are up-to-date

## Key Conventions and Patterns
  - Tests use Go's `testing` package, with custom assertions in `internal/assert`.
  - Integration tests require a running MongoDB instance.
  - Many tests use `mtest` for integration scenarios.
  - Custom assertion helpers are in `internal/assert/` (adapted from `stretchr/testify`).
  - Operation names are constants in `internal/driverutil/operation.go` (source: MongoDB command reference).
  - Aggregation and usage examples are in `examples/documentation_examples/examples.go`.
  - Example pipelines and data flows are provided for reference.

## External Dependencies

## Integration Points

## Contribution and CI

## Quick Start for AI Agents


For more details, see:

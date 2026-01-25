# Agent Instructions

## Purpose

Short, machine-friendly instructions for GitHub Copilot agent mode. Follow these rules strictly when making code changes or tests.

## Top priorities

1. Working Code & Maintainability: produce correct, well-tested, maintainable implementations.
2. Clarity & Simplicity: write clear, simple, and idiomatic code.
3. Optimization: improve performance and resource usage where appropriate.

## Coding style

- Respect the existing project style and case syntax; match surrounding naming and formatting.
- Add concise comments for complex logic and algorithms.
- Break large/complex functions into smaller helper functions for readability and reuse.
- All code must pass the project's linter. See `.golangci.yml` for rules and expectations.

## Unit tests

- Write thorough unit tests covering regular and edge cases; do not write tests that fake passing.
- Where tests are required, achieve 80%–100% coverage for the targeted areas.
- Mock API responses in unit tests; all tests functions must rely on common test utilities where possible. Regroup common test setup/teardown logic into helper functions in a directory named `fake` under each corresponding package.
- Use table-driven tests for functions with multiple scenarios.

## Iteration completion

An iteration is complete only when:

- Linting passes.
- All tests pass.
- Either the user must intervene due to an error loop, or the requested task is functionally complete.

## Security, safety, performance

- Always apply security best practices and safe defaults.
- Prefer efficient algorithms and optimize hotspots where appropriate.

## Behavioral notes

- If behavior or requirements are unclear, ask the user before implementing major assumptions.
- Keep changes minimal and focused; avoid unrelated refactors unless requested.

## Makefile commands

- `make generate`: runs code generation from source code.
- `make lint`: runs `golangci-lint` and reports lint results.
- `make test`: runs unit tests.

These three commands should always pass before considering an iteration complete.

## Generated files

- NEVER edit generated files whose filenames begin with `zz_`. These files are overwritten by code generation. Instead, you must modify the generation source files.

## Reference

See the repository files and attachments for context and data examples when optimizing logic.

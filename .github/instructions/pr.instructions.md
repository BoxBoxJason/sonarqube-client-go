# Copilot PR Review Instructions

These instructions are for an automated Copilot agent performing PR reviews on this repository. They focus on code quality, maintainability, and correctness, and are intended to drive precise, repeatable review behavior by the agent.

## Scope

- **Repository:** entire codebase; prioritize packages under `pkg/`, `sonar/`, and top-level service packages.
- **Generated files:** NEVER modify generated files whose filenames begin with `zz_`. Flag PRs that attempt to change these files and suggest editing the generation source instead.

## Primary Review Goals

- **Maintainability:** Every `struct`, `function`, and `method` must have corresponding godoc comments. Document struct fields when their meaning is not trivial or when they are part of the public API.
- **Refactor Safety:** Any refactor must preserve 1:1 feature availability. Validate that behavior and public API surface remain unchanged; if behavior changes are intentional, require an explicit rationale in the PR description.
- **Duplication:** Avoid duplicate code. When similar logic appears in multiple places, prefer extracting a shared helper or utility and point to the new common location in the review.
- **Dead Code:** Flag unreachable, unused, or vestigial code and recommend removal. If retained intentionally (e.g., for future work), require an explanatory comment and a follow-up task reference.
- **Testing:** New features and public functions must include tests that cover regular cases and edge cases. Prefer table-driven tests for multi-scenario functions and mock external API interactions in tests where appropriate.
- **Performance:** Review performance-critical functions for algorithmic complexity and hotspot inefficiencies. Suggest algorithmic improvements or micro-optimizations when the function's runtime may impact user-visible latency or scale.

## Checklist for Each PR (agent actions)

- **Docs:** Ensure `godoc` comments exist for all exported and package-level symbols; request missing comments.
- **Behavioral Parity:** For refactors, verify call sites and public behavior match previous behavior; if diffs show changed semantics, require explicit justification.
- **Duplication Audit:** Search for repeated code patterns; when found, recommend consolidation and provide a candidate extraction suggestion.
- **Dead Code Detection:** Identify unused functions, types, or variables and flag them unless a clear justification exists in the PR.
- **Tests Presence & Quality:** Confirm tests accompany new functionality; check they include standard and edge-case inputs and that they are deterministic and clear.
- **Performance Notes:** If a function is in a hotspot (loops, frequently-called APIs, marshaling/unmarshaling), check for obvious inefficiencies (e.g., repeated allocations, unnecessary conversions) and recommend improvements.
- **Generated Files Safety:** Ensure `zz_*.go` files are not edited; if they are, request modification of the generation source instead and note which generator must be updated.

## Behavioral Rules for the Agent

- When suggesting code moves or extractions, include a precise code snippet illustrating the suggested helper and at least one call-site change.
- When flagging missing documentation or tests, provide an example comment or test-case scaffold to reduce churn for the author.
- Avoid making changes that alter public contracts without an explicit, documented justification in the PR.

## When to Block a PR

- Missing tests for new public behavior.
- Refactors that change behavior without documented rationale and migration path.
- Edits to generated files (`zz_*.go`) without updating the generator.
- Introduced duplicate implementations where a clear shared implementation already exists.

## Notes & References

- This guidance aligns with the repository's internal priorities: working code, maintainability, clarity, and performance. Use this as the basis for automated review comments and code-suggestion patches.

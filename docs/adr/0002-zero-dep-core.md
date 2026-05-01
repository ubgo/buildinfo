# ADR 0002 — Zero third-party dependencies in the core

**Status:** Accepted
**Date:** 2026-04-26

## Context

`buildinfo` is the kind of small utility users want to drop into any service without inheriting a transitive dependency tree. Even pulling in OpenTelemetry "just for the resource attributes" pollutes module graphs of users who don't care about OTEL.

## Decision

The core module (`github.com/ubgo/buildinfo`, repository root) **must not import any third-party package**. Standard library only.

A CI gate (`zero-dep-gate`) fails the build if `go.mod` of the core module gains any `require` line that is not under `github.com/ubgo`.

All integrations with third-party libraries live in adapter modules under `contrib/`. Each adapter has its own `go.mod` and freely depends on its target ecosystem.

## Consequences

- Core consumers download only the standard library.
- The core API surface stays small and testable without mocks.
- Adding a new integration requires creating a new adapter module — never adding a dep to the core.
- Adapter authors take on dep maintenance for their target ecosystem.

## Alternatives considered

- **Allow OpenTelemetry in core because almost everyone uses it.** Even if true, it punishes the small fraction that doesn't and makes upgrade decisions political. Rejected.
- **Build tags to gate optional imports.** Build tags interact badly with `go mod` — the dependency is still downloaded even when the build tag excludes the file. Rejected.

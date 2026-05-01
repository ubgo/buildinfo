# ADR 0005 — Population precedence: `-ldflags` over `runtime/debug`

**Status:** Accepted
**Date:** 2026-04-26

## Context

`Info` fields can be populated from two distinct sources:

1. **`-ldflags` overrides.** Variables (`Version`, `Commit`, `BuildTime`, `Branch`) declared in `ldflags.go` and stamped at build time via `go build -ldflags="-X github.com/ubgo/buildinfo.Version=1.2.3"`.
2. **`runtime/debug.ReadBuildInfo`.** Auto-populated VCS metadata Go embeds since 1.18 (`vcs.revision`, `vcs.time`, `vcs.modified`) plus dependency module list.

These sources can disagree (e.g. CI passes a release tag via `-ldflags` while `vcs.revision` reflects the underlying commit hash).

## Decision

`-ldflags` overrides **always win** over `runtime/debug` data when both are present.

`load()` populates fields in the following order:

1. Set `GoVersion`, `GOOS`, `GOARCH` from the `runtime` package (always available).
2. Read `runtime/debug.ReadBuildInfo`:
   - `vcs.revision` → `Commit`
   - `vcs.time` → `BuildTime`
   - `vcs.modified` → `Modified`
   - Iterate `Deps` to populate `Modules`.
3. Overlay `-ldflags` values on top of the above when the corresponding ldflags variable is non-empty.
4. Apply sentinel defaults (`"dev"` for `Version`, `"unknown"` for `Commit`, `BuildTime`, `Branch`) for any fields still empty.

`Branch` is **never populated by `runtime/debug`** — VCS metadata embedded by the toolchain does not include the branch. It must be set via `-ldflags` or it falls back to `"unknown"`.

## Consequences

- CI-driven release stamping always trumps inferred VCS metadata.
- Users who build without `-ldflags` still see meaningful `Commit`, `BuildTime`, `Modified` values.
- `Branch` is the only field where `runtime/debug` is silent.

## Alternatives considered

- **`runtime/debug` wins over ldflags.** Inverts the conventional CI pattern; release tags would be silently overwritten by the underlying commit hash. Rejected.
- **First non-empty wins (no precedence rule).** Order-dependent and surprising. Rejected.
- **Refuse to start if both sources disagree.** Hostile to developers; mismatched values are common (e.g. building from a tagged commit reports both the tag via ldflags and the commit hash via VCS). Rejected.

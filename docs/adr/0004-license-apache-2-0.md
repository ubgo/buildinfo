# ADR 0004 — License: Apache License 2.0

**Status:** Accepted
**Date:** 2026-05-02
**Supersedes:** an earlier draft choosing PolyForm Strict 1.0.0.

## Context

The `ubgo/*` ecosystem is built and maintained by a single owner who wants to encourage broad adoption — including by for-profit companies — while ensuring nobody can claim the work as their own without attribution.

An earlier draft of this ADR chose **PolyForm Strict 1.0.0**. On reading the canonical license text, that license restricts use to **noncommercial purposes only**, which would block the lib's intended audience (companies building services). That draft is therefore rejected.

## Decision

This repository is licensed under the **Apache License, Version 2.0**.

- License text: https://www.apache.org/licenses/LICENSE-2.0
- Copyright: `Copyright 2026 ubgo`
- A `NOTICE` file at the repository root carries the attribution string and is propagated to derivative works under section 4(d) of the license.

| Permitted | Required | Forbidden |
|-----------|----------|-----------|
| Commercial use | Keep upstream copyright + license + NOTICE | Trademark misuse (the name "ubgo" is reserved) |
| Modification | State significant changes | — |
| Distribution (source or binary) | — | — |
| Patent use (with grant) | — | — |
| Sublicensing under compatible terms | — | — |

## Consequences

- Companies' standard legal-review processes approve Apache-2.0 in seconds.
- Forks are technically permitted but must keep the original copyright + license + NOTICE — nobody can credibly claim sole authorship.
- Patent grant protects both the maintainer and downstream users from patent litigation.
- The lib joins the same license family as Kubernetes, OpenTelemetry, gRPC, Prometheus, Docker, etc. — frictionless adoption.

## Alternatives considered

- **PolyForm Strict 1.0.0.** Rejected once we read the canonical text — it is noncommercial-only, which contradicts the goal of company adoption.
- **MIT.** Permissive but lacks an explicit patent grant. Apache-2.0 is strictly stronger.
- **BUSL (Business Source License).** Designed for hosted-service competition (e.g. CockroachDB). Adds complexity (conversion-to-OSS clauses, additional-use grants) we do not need for libraries.
- **Functional Source License (FSL).** Newer "delayed open source" license used by Sentry. Same complexity concern as BUSL without the track record.
- **Custom proprietary EULA.** Custom legal drafting risk. Rejected.

## Trademark note

The name `ubgo` and any associated marks are not licensed under Apache-2.0. Forks may use the source code subject to the license but may not represent themselves as the official `ubgo` distribution.

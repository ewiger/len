# LIP-XXXX: Proposal Title

## Metadata

| Field | Value |
| --- | --- |
| LIP | XXXX |
| Title | Proposal Title |
| Status | Draft |
| Author | TBD |
| Created | YYYY-MM-DD |
| Updated | YYYY-MM-DD |
| Area | Surface Syntax, Semantics, Tooling |
| Target | len.l1 |

## Summary

State the proposal in 2 to 5 short paragraphs.

- what changes
- what stays core
- what the proposal explicitly does not introduce

## Motivation

Explain why the change is needed now.

Focus on design pressure, ambiguity, missing expressiveness, or unnecessary surface complexity.

## Problem Statement

Describe the current problem precisely.

- what existing syntax or semantics are unclear
- what user-facing confusion results
- what design risk appears if the issue is left unresolved

## Proposed Changes

Break the proposal into explicit decisions.

### 1. First Decision

State the rule.

### 2. Second Decision

State the rule.

### 3. Third Decision

State the rule.

## Grammar Sketch

Provide an illustrative grammar sketch, not necessarily the final parser grammar.

```ebnf
TopLevelDecl = ... ;
```

Call out additions, removals, and renamed forms.

## Desugaring / Lowering Model

Show how the surface proposal reduces to the core language.

Prefer explicit lowerings into:

- `type`
- `rel`
- `spec`
- `fn`

If a form cannot be lowered cleanly, state why.

## Examples

Include compact examples.

### Before / After

Show migration from old surface to proposed surface.

```len
# before
```

```len
# after
```

### Canonical Example

Show the preferred idiomatic form after the proposal.

## Rationale

Explain why this proposal fits `len` specifically.

Address these questions directly:

- why the proposal preserves `type` and `rel` as core primitives
- how `spec` remains the general mechanism for laws and definitions
- how `fn` remains the executable or constructive form
- why the new surface does not introduce accidental object-oriented or runtime semantics

## Backward-Compatibility Notes

List affected syntax, migration expectations, and parser or validator impact.

## Open Questions

Number unresolved design questions explicitly.

## Non-Goals

State what the proposal is not trying to solve.

## Future Proposal Notes

If more detailed companion documents are needed, list them here.

- `README.md`: proposal overview and decisions
- `plan.md`: implementation plan if accepted
*** End Patch
---
name: write-adr
description:
  "Write, draft, review, or edit an Architecture Decision Record (ADR)
  for the runkit project. Use when recording a new architectural decision, revising
  an existing ADR, or checking an ADR against project style rules. Triggers: write
  ADR, draft ADR, new ADR, create ADR, review ADR, architecture decision record."
argument-hint: "the decision to document, e.g., 'use gRPC for inter-service communication'"
---

# Write ADR

Guides drafting, reviewing, and editing Architecture Decision Records for the
runkit project. All ADRs follow the Nygard format with the project-specific
style rules defined in this skill.

The reference exemplar for this project is
`docs/adr/0002-rendezvous-hashing-for-workload-assignment.md`.

For voice and tone, refer to these ADRs from the dogma repository:

- [0009-immutable-keys.md](https://github.com/dogmatiq/dogma/blob/main/docs/adr/0009-immutable-keys.md)
- [0020-identifier-constraints.md](https://github.com/dogmatiq/dogma/blob/main/docs/adr/0020-identifier-constraints.md)
- [0022-remove-crud-application-support.md](https://github.com/dogmatiq/dogma/blob/main/docs/adr/0022-remove-crud-application-support.md)
- [0023-message-order-guarantees.md](https://github.com/dogmatiq/dogma/blob/main/docs/adr/0023-message-order-guarantees.md)
- [0026-event-stream-based-projection-occ.md](https://github.com/dogmatiq/dogma/blob/main/docs/adr/0026-event-stream-based-projection-occ.md)

## When to Use

- Drafting a new ADR for a technical or architectural decision
- Reviewing an existing ADR for style, correctness, or completeness
- Editing an ADR to fix issues flagged in a review

## Assets

- [ADR template](./assets/adr-template.md)

## ADR Status

The status value and relationship annotations are independent. The status
describes the state of this ADR; the annotations describe its relationships to
other ADRs. Both appear in the `## Status` section, with annotations as bullets
below the status value.

### Status values

| Value        | Meaning                                                                                |
| ------------ | -------------------------------------------------------------------------------------- |
| `Proposed`   | Drafted, not yet accepted; transient state during review                               |
| `Accepted`   | In force; includes explicit decisions _not_ to do something                            |
| `Deprecated` | No longer applicable; no replacement exists                                            |
| `Superseded` | No longer in force; replaced by one or more ADRs (named via `- Superseded by` bullets) |

### Admonitions for inactive ADRs

ADRs with status `Proposed`, `Superseded`, or `Deprecated` must include a
GitHub-style admonition immediately after the `## Status` section (before
`## Context`). Copy the admonition blocks exactly as shown below, including
whitespace and line breaks, because GitHub can render them incorrectly
otherwise:

**Proposed:**

```
> [!NOTE]
> This decision has not yet been accepted and is subject to change.
```

**Superseded:**

```
> [!WARNING]
> This decision has been superseded. Refer to the replacement ADR(s) listed
> above.
```

**Deprecated:**

```
> [!WARNING]
> This decision has been deprecated and is no longer applicable.
```

### Relationship annotations

Always written as bullets. Always add the counterpart annotation to each linked
ADR. Multiple annotations of the same kind are allowed.

| This ADR                           | Linked ADR                            | Meaning                                                                        |
| ---------------------------------- | ------------------------------------- | ------------------------------------------------------------------------------ |
| `- Supersedes [N. Title](file.md)` | `- Superseded by [N. Title](file.md)` | This ADR replaces the linked one; the linked ADR's status becomes `Superseded` |
| `- Amends [N. Title](file.md)`     | `- Amended by [N. Title](file.md)`    | This ADR modifies the linked one without replacing it; both remain `Accepted`  |
| `- References [N. Title](file.md)` | `- Referenced by [N. Title](file.md)` | This ADR cites the linked one for context only; no status change               |

## Drafting a New ADR

### 1. Determine the file number

List the files in `docs/adr/`. The next ADR gets the next sequential number,
zero-padded to four digits. File name pattern: `NNNN-title-with-dashes.md`.

Numbers are assigned at draft time and are provisional. If two concurrent drafts
claim the same number, the one that merges second must renumber. Abandoned
drafts are deleted; do not leave placeholder files or gaps.

### 2. Set the status

New ADRs start as `Proposed`. Do not change the status during drafting. See
[ADR Status](#adr-status) for allowed values and relationship annotations.

### 3. Write the title

The title states the decision, not the problem. It is a noun phrase, not a
question.

- Good: "Rendezvous hashing for workload assignment"
- Bad: "How should we assign workloads?"

**Negative decisions:** If the team explicitly decides _not_ to do something,
write and accept a normal ADR for it. The title is still a decision noun phrase:
"No Redis use in cache-layer", "Avoid gRPC for inter-service communication". Context
explains the reasoning; Consequences describe what not having X means going
forward. Do not use a special status for this -- it is a first-class `Accepted`
decision.

### 4. Draft the three sections

Copy [adr-template.md](./assets/adr-template.md) and fill it in.

Wrap regular paragraph text at 80 characters per line.

**Context** -- one to three paragraphs:

- State the problem or constraint that motivates the decision.
- Introduce every term that will appear in the Decision section as prose, woven
  into the narrative. If "candidate" appears in Decision, it must be defined
  here first. Do not introduce terms as a glossary list or bullet definitions.
- Do not propose or evaluate solutions in this section.

**Decision** -- explain what we will do and why:

- Use first-person plural throughout: "We will...", "We need...",
  "We considered..."
- Include a "Dismissed alternatives" subsection when alternatives were seriously
  considered. Acknowledge genuine advantages before explaining why each was
  rejected. Be specific about the reasons.
- Any pseudocode must use the exact terms from the prose.

**Consequences** -- inherent properties only:

- Describe what is true as a result of the decision, not what we hope will be
  true.
- Use "we are free to..." for future possibilities, not "future ADRs will..."
- Note any glossary terms this ADR introduces. Example: "This ADR introduces
  two terms to the glossary: rendezvous hashing and self-affinity."
- Do not cite rejected or superseded ADRs as the basis for a decision.

### 5. Add references

- Link external concepts to a well-regarded source on first use. Prefer
  Wikipedia for general concepts; use a more authoritative source (RFC, spec,
  official docs) when one exists.
- Link code identifiers to pkg.go.dev.
- Link RFCs to rfc-editor.org.
- Use markdown reference-style links. Collect them at the bottom of the file
  inside a `<!-- references -->` comment block. Keep the list alphabetized.

### 6. Run the pre-flight checklist

Work through every item in the [Style Checklist](#style-checklist).

## Reviewing an Existing ADR

Read the ADR, work through the [Style Checklist](#style-checklist), and report
all issues found. For each issue: quote the offending text, name the rule it
breaks, and suggest a corrected version.

## Style Checklist

### Structure

- [ ] Status value is one of the four allowed values; no other text on the same line
- [ ] Relationship annotations are bullets below the status value, using only the six allowed verbs
- [ ] Every relationship annotation has its counterpart added to the linked ADR
- [ ] `Proposed`, `Superseded`, and `Deprecated` ADRs include the exact required admonition after `## Status`, including whitespace and line breaks
- [ ] Title is a decision noun phrase, not a question or problem statement
- [ ] Exactly three top-level sections: Context, Decision, Consequences
- [ ] No extra top-level sections (no "Options", "Alternatives", "Background")
- [ ] File named `NNNN-title-with-dashes.md`
- [ ] One decision per ADR; split if scope has crept
- [ ] Date is today's date (new ADRs only)

### Voice and characters

- [ ] First-person plural throughout ("We will...", "We considered...")
- [ ] Conversational tone, not academic
- [ ] Regular paragraph text is wrapped at 80 characters per line
- [ ] No non-ASCII characters in prose: no em dashes, en dashes, or curly quotes
- [ ] ASCII punctuation only; hyphens, not dashes
- [ ] Non-ASCII characters in code blocks and formulas are acceptable

### Terminology

- [ ] Every term introduced before first use
- [ ] Consistent terminology throughout; no synonym alternation
- [ ] No unexplained jargon; written for average programmers
- [ ] No terms that collide with established meanings in other domains
- [ ] Dogma ecosystem terms (command, aggregate, process, etc.) used without
      redefinition

### Claims and evidence

- [ ] Performance claims quantified ("in the order of nanoseconds") or hedged
- [ ] No specific figures cited without a citable source
- [ ] Tradeoffs represented honestly; dismissed alternatives not strawmanned

### References

- [ ] External concepts linked to a well-regarded source on first use (Wikipedia, RFC, spec, or official docs)
- [ ] Code identifiers linked to pkg.go.dev
- [ ] RFCs linked to rfc-editor.org
- [ ] Reference-style links used throughout
- [ ] References collected at the bottom in a `<!-- references -->` comment block
- [ ] Reference list alphabetized

### Scope boundaries

- [ ] No undecided concepts or unfinished implementation details referenced
- [ ] Consequences are inherent properties, not aspirational claims
- [ ] No superseded or deprecated ADR is cited as the basis for a decision
- [ ] No "future ADRs will..." phrasing (use "we are free to..." instead)

### Dismissed alternatives (if present)

- [ ] Subsection appears under Decision, not as a top-level section
- [ ] Genuine advantages acknowledged before reasons for rejection
- [ ] Reasons are specific, not vague ("adds complexity without benefit here
      because...")

### Pseudocode (if present)

- [ ] Uses the same terms as surrounding prose (not synonyms or abbreviations)
- [ ] Minimal and readable

### Glossary

- [ ] Any new terms introduced by this ADR are called out in Consequences

---
name: glossary
description: "Add, update, or review glossary terms in a Dogmatiq project. Use when an ADR introduces new terms, a term definition needs revision, a new term needs to be cross-referenced, or the glossary index needs updating. Triggers: add glossary term, update glossary, glossary entry, define term, glossary."
argument-hint: "the term(s) to add or update, e.g., 'optimistic conflict resolution'"
---

# Maintain Glossary

Guides adding, updating, and reviewing glossary entries in Dogmatiq projects.
Each project keeps a `docs/glossary.md` that defines terms specific to that
project. Terms already defined in the [Dogma glossary] are not redefined —
they are cross-referenced only.

## When to Use

- A new term needs to be defined (whether or not it originates from an ADR)
- An ADR's Consequences section says "This ADR introduces X to the glossary"
- A term appears in project prose but has no glossary entry yet
- An existing definition needs correction or clarification
- Reviewing a glossary for completeness, style, or stale links

## File Structure

```
docs/glossary.md       # Project-level glossary
```

The file has three structural parts:

1. **Intro paragraph** — one sentence linking to the parent (Dogma) glossary
   and explaining what this file covers.
2. **Alphabetical index** — all 26 letters, one per line, separated by
   `•`. Letters with at least one entry are markdown links to their section
   anchor (`[R](#r)`). Letters with no entries are plain text.
3. **Letter sections** — `## A`, `## B`, etc. Each term is a `### Term name`
   heading followed by its definition. Only include letter sections that have
   at least one entry.
4. **Reference link blocks** at the bottom, inside HTML comment delimiters:
   - `<!-- anchors -->` — intra-glossary cross-reference anchors
   - `<!-- ADRs -->` — links to ADR files

## Adding a New Term

### 1. Check the upstream glossary

Before writing a definition, confirm the term is not already defined in
the [Dogma glossary]. If it is, do not redefine it. If project prose
needs to reference it, link directly to the upstream glossary entry.

### 2. Identify the source

If the term is introduced by an ADR, note the file path for a reference link.
Not all terms originate from an ADR — some are defined because they appear
frequently in project documentation without a formal decision record.

### 3. Place the entry in the correct letter section

- Terms are sorted alphabetically within their letter section, case-insensitive.
- If no section exists for that letter yet, create `## X` and insert it in
  alphabetical order among the existing sections.

### 4. Write the definition

Follow these rules:

- **Prose only.** No bullet lists or sub-headings inside a definition.
- **One paragraph as the primary definition.** Use additional paragraphs only
  when genuinely necessary.
- **Define the term itself, not how it is used.** The first sentence should
  complete "X is ..." or "An X is ...".
- **Follow `docs/style.md`** for line wrapping, character, and reference link
  conventions.
- **Cross-reference related terms** inline in the definition prose where it
  reads naturally. If a related term is not mentioned in the prose, add a
  trailing `See [term].` line. Do not use `See also`.
- **If an ADR introduced the term**, add `See [ADR-NNNN].` at the end.

### 5. Update the alphabetical index

If the letter was previously plain text in the index, change it to a link:
`[X](#x)`. Do not change letters that already link correctly.

### 6. Add reference links at the bottom

Add any new intra-glossary anchors to the `<!-- anchors -->` block and any new
ADR links to the `<!-- ADRs -->` block. Keep each block sorted alphabetically
by the link label. When a term is used in its plural form in prose, add a
separate plural reference link so that `[term-names]` reads naturally in source.
Use the form:

```markdown
<!-- anchors -->

[term-name]: #term-name
[term-names]: #term-name

<!-- ADRs -->

[ADR-NNNN]: adr/NNNN-title-with-dashes.md
```

## Updating an Existing Term

- Edit the definition prose in place; do not change the heading or anchor.
- If a term has been superseded or renamed, add a `See [replacement].` line
  and leave the original heading so existing links do not break.
- If the term was introduced by an ADR that has since been superseded,
  update the `See also [ADR-NNNN].` reference and the link in the
  `<!-- ADRs -->` block.

## Style Checklist

Work through every item when adding or reviewing entries. Also check the
`<!-- project-specific rules -->` comment at the end of the glossary file for
any additional conventions specific to this project.

### Index

- [ ] Index contains all 26 letters in alphabetical order, one per line
- [ ] Every letter with at least one entry is a link in the index
- [ ] Every letter without entries is plain text in the index
- [ ] Letters are separated by `•`

### Entries

- [ ] Follows `docs/style.md` (line wrapping, characters, reference links)
- [ ] Term heading uses the canonical capitalisation for the term (not all-caps
      or all-lowercase unless that is the proper form)
- [ ] Definition is prose, no bullet lists
- [ ] Term is not already defined in the upstream Dogma glossary
- [ ] Cross-references use reference-style links that resolve within the file
- [ ] Related terms are linked inline in prose or via a `See [term].` line, not
      `See also`
- [ ] If the term was introduced by an ADR, it is cited in a `See [ADR-NNNN].`
      line
- [ ] Entry is placed in alphabetical order within its letter section

### Reference link blocks

- [ ] `<!-- anchors -->` block contains an anchor for every cross-referenced
      term that lives in this file, including plural forms used in prose
- [ ] `<!-- ADRs -->` block contains a link for every cited ADR

<!-- references -->

[Dogma glossary]: https://github.com/dogmatiq/dogma/blob/main/docs/glossary.md

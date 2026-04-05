# Prose style

Common conventions for all markdown documentation in this repository.
Domain-specific rules live in their respective skill files; this document
covers the shared baseline.

## File naming

Use lowercase names for all documentation files (e.g. `style.md`,
`glossary.md`). Do not capitalize file names unless an external convention
requires it (e.g. `LICENSE`, `Makefile`).

## Line wrapping

Reflow paragraph text so that each line fits as many words as possible
within 80 characters.

## Characters

- No curly quotes or curly apostrophes in prose. Use straight ASCII quotes and
  apostrophes.
- Use a proper em dash character — not hyphen(s) — as a parenthetical separator.
  However, do not overuse em dashes. If a sentence or paragraph is heavy with
  them, consider restructuring with commas, colons, or separate sentences
  instead.
- Non-ASCII characters in code blocks and formulas are acceptable.

## Reference links

- Use markdown reference-style links throughout.
- Collect link definitions at the bottom of the file, inside HTML comment
  delimiters that label each group (e.g. `<!-- anchors -->`,
  `<!-- references -->`).
- Keep each group sorted alphabetically by label.
- When a term appears in plural form in prose, add a separate plural reference
  link (e.g. `[instructions]: #instruction`) so the source reads naturally.
- Every reference link definition must be used at least once in the document.
  Remove unused definitions.
- Every reference link used in prose must have a matching definition. Do not
  leave broken references.

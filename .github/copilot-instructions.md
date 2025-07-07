<!-- vale off -->

# Copilot instructions

This is a Go repository containing interfaces that define the contract between a
Dogma **application** and a Dogma **engine**.

API documentation is a central part of this repository. The primary audience is
application developers — the end-users of the Dogma API. Assume they're
proficient in Go and backend development, but may be unfamiliar with Dogma,
event-sourcing, or DDD.

All contributions should align with the style and expectations outlined below.

## Development

- Never push to `main`. Only push feature branches.
- Always rebase your branch on `main` before pushing:
  - Your branch's history must contain the head of `main`.
  - Your commits (and only your commits) must follow the head of `main`.
- Run `make` to build and test the project.
- Run `make lint precommit` before committing.
- Don't introduce uncovered code. Use `make artifacts/coverage/cover.out`.
- Don't create empty commits.
- Follow Go best practices and idiomatic patterns.
- Readability and understandability are paramount.
- Code without tests is incomplete.
- Use table-driven unit tests when possible.
- Maintain existing code structure and organization.
- Plan for "forward compatibility" by using placeholder options and similar
  strategies where appropriate.

## Documentation

- Code without API documentation is incomplete.
- Update the `CHANGELOG.md` file when making changes that affect the API.
- Suggest changes to the `docs/` folder when appropriate.
- `make lint` runs `vale` to check documentation quality.
- You don’t need to fix existing `vale` issues, but must not introduce new ones.

### Syntax

- Use [Go documentation conventions] for formatting and structure.
- Use US English spelling and grammar.
- Use Oxford commas.
- Hyphenate terms like "side-effect" and "real-time".
- Avoid non-ASCII punctuation.
- Use punctuation _outside_ quotation marks:
  - ✅ This is a "handler".
  - ❌ This is a "handler."
- Link to interface methods using syntax supported by `pkg.go.dev`:
  - ✅ `[Interface].Method`
  - ❌ `[Interface.Method]`

### Style

- Use precise, consistent, unambiguous language.
- Be concise, clear, and developer-friendly.
- The target audience is the application developer, not Dogma maintainers.
- Assume familiarity with Go, not Dogma.
- Match definitions from the [glossary].
- Use instructional phrasing to guide the reader to correct usage:
  - ✅ "Prefer conditionally disabling a handler ..."
  - ✅ "Use lowercase sentences with no trailing punctuation ..."
  - ✅ "Avoid constructing values of this type directly ..."
  - ❌ "The developer should prefer conditionally disabling a handler ..."
- Avoid constructions that misrepresent agency:
  - ✅ "The handler records an event."
  - ✅ "The engine persists the event."
  - ❌ "The command records an event."
  - ❌ "The application persists the event."
- Document intended usage, not just behaviour.
- Document design constraints and invariants where relevant.
- Ultrathink about whether changes to existing documentation alter the meaning.
- Think hard before deciding to ignore "suggestion" level Vale issues.
- Ultrathink before deciding to ignore "warning" level Vale issues.
- Don't ignore "error" level Vale issues.
- Reflow documentation to wrap at 80 characters, but don't split Markdown-style
  links across lines.
- Avoid [RFC 2119] style keywords in API documentation.
- Avoid [RFC 2119] style keywords in Markdown documentation, unless the document
  already includes the [RFC 2119] explanatory text.

## Repository structure

- All Go code is in the root of the repository.
- `docs/adr` — Architecture decision records that document highly visible or
  complex design decisions.
- `docs/concepts.md` — An overview of Dogma's concepts. Read this to gain a
  deeper understanding of the design.
- `docs/glossary.md` — Definitions of Dogma's terminology. Ensure that words in
  the glossary are only used with their documented meanings.
- `CHANGELOG.md` — A file that documents changes to the API.

[glossary]: ../docs/glossary.md
[Go documentation conventions]: https://go.dev/doc/comment
[RFC 2119]: https://datatracker.ietf.org/doc/html/rfc2119

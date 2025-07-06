# Copilot instructions

This is a Go repository containing interfaces that define a contract between a
Dogma "application" and a Dogma "engine".

The API documentation is a key part of the repository. The primary audience is
the "application developers" who are the end-users of the Dogma API. Assume that
application developers are familiar with Go and backend development but may not
be familiar with event-sourcing, DDD, or Dogma itself.

Please follow these guidelines when contributing:

## Development

- Check if your feature branch has diverged from `main` before pushing to GitHub.
  - If your branch diverges, rebase it and "force push" to GitHub.
  - The commit referenced by `main` must be in your branch's history.
  - Your commits (and only _your_ commits) must occur after the head of `main`.
  - Never push changes to any branch other than your feature branch.
- Do not introduce uncoverged changes. `make artifacts/coverage/cover.out`
  builds a coverage.
- Do not make empty commits.
- Run `make precommit` before pushing to GitHub.
- Run `make` to build and test the project.
- Run `vale` to check documentation for style and grammar issues. Don't worry
  about fixing existing issues, but do not introduce new issues.

## Repository structure

- All Go code is in the root of the repository.
- `docs/adr` — Architecture decision records that document highly visible or
  complex design decisions.
- `docs/concepts.md` — An overview of Dogma's concepts. Read this to gain a
  deeper understanding of the design.
- `docs/glossary.md` — Definitions of Dogma's terminology. Ensure that words in
  the glossary are only used with their documented meanings.
- `CHANGELOG.md` — A file that documents changes to the API.
- `artifacts/` — Uncommitted artifacts produced by `make`.

## Guidelines

- Follow Go best practices and idiomatic patterns.
- Readability and understandability are paramount.
- Code without API documentation is incomplete.
- Code without tests is incomplete. Use table-driven unit tests when possible.
- Maintain existing code structure and organization.
- Plan for "forward compability" by using placeholder options and similar
  strategies where appropriate.
- Suggest changes to the `docs/` folder when appropriate.
- Use US English grammar and orthography and Oxford commas.
- Update the `CHANGELOG.md` file when making changes that affect the API.

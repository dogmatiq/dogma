<!-- vale off -->

This is a Go repository containing interfaces that define a contract between a
Dogma "application" and a Dogma "engine". It contains very little executable
code. The API documentation is a key part of the repository.

Assume that application developers are familiar with Go and backend development
but may not be deeply versed in event-sourcing or DDD. They want clear,
concrete guidance on how to structure their code and reason about message
handling.

Some interfaces are implemented by "application developers", who are the end
users of the Dogma ecosystem. Others are implemented by "engine developers" and
used by the application.

Please follow these guidelines when contributing:

## Development Flow

- Check if there have been changes to the `main` branch since your last commit.
  You will need to fetch remotes to do this, as you wont have a local copy of
  `main` by default. Rebase your feature branch often to ensure that the commit
  referenced by `main` is in your branch's history and your changes occur after
  those changes on `main`.

- Run `make` to build & test the project after each change.

- Run `make coverage` to build coverage reports to verify that new code is
  covered by tests. CI will fail if a change introduces a decrease in
  coverage or new code is uncovered.

- Run `make precommit` to run any local pre-commit hooks. This runs the tests,
  builds for all target operating systems, and executes any other pre-commit
  hooks that are by the Makefiles, such as code generation and formatting.

- Run `vale` to check documentation for style and grammar issues. There are
  currently lots of existing issues, so don't worry about fixing them all at
  once, but do not introduce new issues.

## Repository Structure

- `docs/adr` — Architecture Decision Records (ADRs) for substantial or complex
  design decisions. You may be asked to write an ADR. Do not modify existing
  ADRs except to cross-link them to new ADRs, or update their status. The
  audience is primarily the Dogma maintainers, but they should be digestable by
  application developers too. Keep them simple.

- `docs/concepts.md` — Overview of the Dogma concepts and architecture. Read
  this to gain a deeper understanding of the design.

- `docs/glossary.md` — Definitions of Dogma's terminology. Some words have
  narrower meanings than their common usage. Ensure that any usage of words from
  this document are intended to carry the definition from the glossary.

- `*_nocoverage.go` - Files that are excluded from code coverage reports. Use
  these only for "tag" methods that exist solely to satisfy interface
  requirements and are never executed.

- `artifacts/` — All uncommitted artifacts, such as build outputs and coverage
  reports.

- `.makefiles/` - Vendored Makefiles, do not modify these directly.
- All code is in the root of the repository. Do not introduce subpackages.

## Key Guidelines

- Follow Go best practices and idiomatic patterns.

- Maintain existing code structure and organization.

- Write unit tests for new functionality. Use table-driven unit tests when possible.

- Changes to the interfaces can have significant implications for both
  application and engine developers. We are currently in a pre-release period,
  so backwards compatibility is not the primary concern, but sweeping changes
  can still be problematic — avoid them unless absolutely necessary or
  explicitly asked to do so.

- This repository makes heavy use of Go's "functional options" pattern. The
  option types are designed in such a way that a single option that can be used
  by multiple different functions. The intent is to keep the usage as readable
  as possible.

## Documentation

- Readability and understandability are paramount.
- Use `vale` to check for problems with new documentation.
- Document public APIs and complex logic.
- Suggest changes to the `docs/` folder when appropriate.
- Use US English grammar and orthography, with Oxford commas.

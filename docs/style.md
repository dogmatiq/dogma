# Dogma Documentation Strategy

## Audience

- Primary audience: mid-to-senior Go developers with at least a passing interest in
  message-driven and distributed applications.
- Readers may be unfamiliar with Dogma or Domain-Driven Design (DDD), so clarity and
  approachability are crucial.
- The documentation targets application developers who implement business logic using
  Dogma, rather than engine developers — though clear semantics are still important
  for engine maintainers.

## Writing Style and Tone

- Use clear, straightforward language that avoids unnecessary jargon.
- Prefer concise, economical phrasing where every word carries meaning.
- Avoid overly formal or spec-like language; instead, aim for a warm and helpful tone.
- Prioritize practical understanding over theoretical completeness.
- Emphasize concepts with concrete examples or references where helpful.

## Structure and Content

- Separate conceptual overviews (e.g., "Concepts" doc) from API/interface reference
  material.
- Document core ideas early, such as the distinction between application and engine,
  and key Dogma concepts (message handlers, aggregates, processes, integrations,
  projections).
- Explain terminology that may be confusing or misunderstood (e.g., “aggregate” in
  Dogma vs. common usage).
- Use modular documentation: a main README or overview document complemented by
  focused guides (e.g., testing, projections, engine development).
- Provide explicit guarantees or behaviors (e.g., message ordering) clearly and
  succinctly, referencing relevant ADRs or design documents.
- Include a “Getting Started” section linking to conceptual docs, examples, and API
  references.
- Maintain cross-links within docs for ease of navigation.

## Documentation Focus

- Document the contract between application and engine via the Dogma interfaces.
- Emphasize separation of concerns: business logic vs. infrastructure (messaging,
  persistence).
- Highlight the tooling available, including static analysis, testing utilities, and
  introspection.
- Clarify the scope and boundaries of Dogma, including what it does not cover (e.g.,
  UI, external API layers).
- Be transparent about ongoing development areas and future directions.

## Contribution and Maintenance

- Consider adding a CONTRIBUTING.md to explain documentation conventions, target
  audience, tone, and the distinction between application vs. engine developer
  perspectives.
- Encourage incremental improvements to keep docs aligned with ADRs and codebase
  evolution.
- Use examples liberally to ground abstract concepts.

---

This framework aims to ensure Dogma’s documentation remains clear, consistent, and
useful for its primary users — application developers building message-driven systems
in Go.

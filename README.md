<div align="center">

# Dogma

Build message-driven, event-sourced applications in Go.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/dogma)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dogma.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/dogma/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/dogma/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/dogma/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dogma/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/dogma)

</div>

## Overview

Dogma is a comprehensive suite of tools for building robust message-driven
applications in Go.

It provides an abstraction for describing your application's business logic with
strict separation from the "engine" responsible for message delivery and
persistence.

## Features

- **[Event sourcing]** — Every state change is persisted as an immutable
  domain event. This enables full auditability and allows read-optimized views
  to be built or rebuilt from the event history at any time.

- **Grounded in [Domain-Driven Design]** — Dogma adopts core DDD concepts to
  guide how applications are decomposed and how messages flow between
  components.

- **High-level testing** — The [testkit] module encourages verification of
  application behavior by making assertions about domain events rather than
  inspecting state. It integrates seamlessly with Go's standard [testing]
  package.

- **Native introspection** — Dogma's static analysis tools visualize message
  flow and application structure, enabling discovery of domain events across
  large codebases and multi-application projects.

- **Domain and infrastructure separation** — Domain logic is cleanly and
  strictly separated from infrastructure concerns such as message delivery,
  persistence and telemetry.

- **Type-agnostic** — Messages and application state can be any Go type that can
  be marshaled to a byte slice, with built-in support for JSON and Protocol
  Buffers.

- **Flexible persistence** — Support for a range of storage options such as
  PostgreSQL and Amazon DynamoDB, enabling use across diverse environments.

## Ecosystem

Dogma is a collection of Go modules that together provide the tools needed to
build, test, analyze, and run message-driven applications.

- [dogma] (this repository) — Defines the API for building applications.
- [testkit] — Utilities for testing Dogma applications.
- [projectionkit] — Utilities for building [projections] in popular database systems.

### Engines

An important Dogma concept is that of the [engine] — a Go module embedded within
your application binary that orchestrates message delivery, state persistence,
and the execution of application logic.

- [verity] — The original Dogma engine, designed for typical application loads
  in smaller deployments. While production-ready, it does not support scaling of
  a single application across multiple machines.

- [veracity] _(under development)_ — The next-generation Dogma engine built for
  horizontal scalability and distributed workloads. Long term, it will fully
  replace Verity, becoming _the_ production Dogma engine.

- [testkit] — A set of tools for testing Dogma applications. It includes an
  in-memory engine that can be used to execute and inspect application behavior
  without persisting state.

## Why "Dogma"?

The name _Dogma_ is a tongue-in-cheek nod to the project's strong opinions about
how message-driven applications should be structured. It's not about rigid
rule-following, but about embracing consistent patterns that enable rich tooling
and clarity in complex systems — without sacrificing flexibility where it
matters.

## Getting Started

If you're new to Dogma, we recommend starting with the [concepts] document to
gain a solid understanding of the core ideas and terminology used throughout the
ecosystem.

You can also explore the [example] application for a practical, working
implementation that demonstrates key concepts in action.

For a detailed reference, see the [API documentation].

<!-- references -->

[api documentation]: https://pkg.go.dev/github.com/dogmatiq/dogma
[concepts]: docs/concepts.md
[dogma]: https://github.com/dogmatiq/dogma
[domain-driven design]: https://en.wikipedia.org/wiki/Domain-driven_design
[engine]: docs/concepts.md#engines
[event sourcing]: https://martinfowler.com/eaaDev/EventSourcing.html
[example]: https://github.com/dogmatiq/example
[projectionkit]: https://github.com/dogmatiq/projectionkit
[projections]: docs/concepts.md#handlers
[testing]: https://pkg.go.dev/testing
[testkit]: https://github.com/dogmatiq/testkit
[veracity]: https://github.com/dogmatiq/veracity
[verity]: https://github.com/dogmatiq/verity

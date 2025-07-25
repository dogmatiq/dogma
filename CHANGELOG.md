# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html
[bc]: https://github.com/dogmatiq/.github/blob/main/VERSIONING.md#changelogs
[engine bc]: https://github.com/dogmatiq/.github/blob/main/VERSIONING.md#changelogs

## [Unreleased]

This release changes how handlers are added to an application. It introduces the
`ApplicationConfigurer.Routes()` method to replace the `RegisterAggregate()`,
`RegisterProcess()`, `RegisterIntegration()`, and `RegisterProjection()`
methods, which have been removed.

It also introduces a global message type registry. The application-defined types
that implement `Command`, `Event`, and `Timeout` must now be added to the
registry before they can be used in handler routes.

### Added

- Added `HandlerRoute` and `MessageRoute` interfaces.
- Added `ViaAggregate()`, `ViaProcess()`, `ViaIntegration()`, and
  `ViaProjection()` for use with `ApplicationConfigurer.Routes()`.
- Added `AggregateHandlerRoute`, `ProcessHandlerRoute`,
  `IntegrationHandlerRoute` and `ProjectionHandlerRoute` types.
- Added `RegisterCommand()`, `RegisterEvent()`, and `RegisterTimeout()` for
  adding types to the message type registry. These functions panic if called
  with a pointer to a type that uses non-pointer receivers to implement the
  `Message` interface. If a type implements `Command`, `Event`, or `Timeout`
  using pointer receivers then a pointer type must be used; otherwise, a
  non-pointer type must be used.
- Added `RegisteredMessageTypes()` and `RegisteredMessageTypeFor()` for querying
  the message type registry.
- Added `Now()` method to `AggregateCommandScope`, `ProcessEventScope`,
  `ProcessTimeoutScope`, `IntegrationCommandScope` and `ProjectionEventScope`.
- Added `WithIdempotencyKey()` option for `CommandExecutor`.
- **[ENGINE BC]** Added `Routes()` method to `ApplicationConfigurer`.
- **[ENGINE BC]** Added `MessageValidationScope` interface, which is embedded in
  `CommandValidationScope`, `EventValidationScope` and `TimeoutValidationScope`.
- **[ENGINE BC]** Added `MessageValidationScope.IsNew()` method.
- **[ENGINE BC]** Added `RecordedAt()` method to `EventValidationScope`.
- **[ENGINE BC]** Added `ScheduledAt()` and `ScheduledFor()` methods to
  `TimeoutValidationScope`.

### Changed

- **[BC]** `HandlesCommand()`,`ExecutesCommand()`, `HandlesEvent()`,
  `RecordsEvent()` and `SchedulesTimeout()` now panic if the message type isn't
  in the message registry.
- **[BC]** Changed `ProjectionMessageHandler.HandleEvent()` to explicitly use
  event stream IDs and event offsets for OCC, instead of abstract versioned
  resources.
- **[ENGINE BC]** Implementations of `MessageRoute` now use
  `RegisteredMessageType` instead of `reflect.Type`.

### Removed

- **[BC]** Removed `Route`, use `MessageRoute` instead.
- **[BC]** Removed `ApplicationConfigurer.RegisterAggregate()`.
- **[BC]** Removed `ApplicationConfigurer.RegisterProcess()`.
- **[BC]** Removed `ApplicationConfigurer.RegisterIntegration()`.
- **[BC]** Removed `ApplicationConfigurer.RegisterProjection()`.
- **[BC]** Removed `RegisterAggregateOption`.
- **[BC]** Removed `RegisterProcessOption`.
- **[BC]** Removed `RegisterIntegrationOption`.
- **[BC]** Removed `RegisterProjectionOption`.
- **[BC]** Removed `ProjectionConfigurer.DeliveryPolicy()`
- **[BC]** Removed `ProjectionScope.IsPrimaryDelivery()`
- **[BC]** Removed `ProjectionDeliveryPolicy`
- **[BC]** Removed `UnicastProjectionDeliveryPolicy`
- **[BC]** Removed `BroadcastProjectionDeliveryPolicy`

## [0.15.0] - 2024-10-03

### Removed

- **[BC]** Removed `Message.Validate()`.

### Changed

- Bumped minimum Go version to v1.23.
- **[BC]** The `Validate()` methods on the `Command`, `Event` and `Timeout`
  interfaces now require a `CommandValidationScope`, `EventValidationScope` or
  `TimeoutValidationScope` argument, respectively.

### Added

- Added `CommandValidationScope`, `EventValidationScope` and
  `TimeoutValidationScope` interfaces. These interfaces are currently empty, but
  methods will be added in the future without causing application-facing
  compatibility issues.

### Deprecated

- The `Message` interface is **no longer deprecated** as it sees widespread use
  within engine implementations. Applications should continue to use the
  more-specific `Command`, `Event` and `Timeout` interfaces wherever possible,
  which no longer share compatible method sets.

## [0.14.3] - 2024-09-27

### Removed

- **[ENGINE BC]** Removed deprecated `fixtures` package.

## [0.14.2] - 2024-08-21

### Removed

- **[ENGINE BC]** Removed generic `fixtures.TestCommand`, `TestEvent`,
  `TestTimeout` and related `TypeA` through `TypeZ` types.

## [0.14.1] - 2024-08-18

### Deprecated

- The `fixtures` sub-package, which is used internally to test Dogma engine and
  toolkit implementations, is now deprecated. It will be removed in a future
  release.

## [0.14.0] - 2024-08-17

### Changed

- **[BC]** The `Command`, `Event` and `Timeout` interfaces are no longer aliases
  for `Command`, they're distinct types. At this stage these interfaces are
  method-compatible with the `Message` interface, but they will diverge in a
  future release.

### Removed

This release removes the "timeout hint" feature (see [ADR-21]).
Application implementers are free to apply their own context deadlines when
handling messages.

- **[BC]** Removed `ProcessMessageHandler.TimeoutHint()`
- **[BC]** Removed `IntegrationMessageHandler.TimeoutHint()`
- **[BC]** Removed `ProjectionMessageHandler.TimeoutHint()`
- **[BC]** Removed `NoTimeoutHintBehavior`

### Deprecated

- Marked the `Message` interface as deprecated. It may be removed in a future
  release.

## [0.13.1]

### Added

- **[ENGINE BC]** Added `Disable()` method to handler configurer interfaces.
- Added placeholder `DisableOption` type for forward-compatibility.

## [0.13.0] - 2024-03-26

### Added

- **[ENGINE BC]** Added placeholder option parameters to the following methods
  and functions:
  - `ApplicationConfigurer.RegisterAggregate()`
  - `ApplicationConfigurer.RegisterProcess()`
  - `ApplicationConfigurer.RegisterIntegration()`
  - `ApplicationConfigurer.RegisterProjection()`
  - `CommandExecutor.ExecuteCommand()`
  - `HandlesCommand()`
  - `ExecutesCommand()`
  - `HandlesEvent()`
  - `RecordsEvent()`
  - `SchedulesTimeout()`

### Removed

This release marks 9 months since the release of [0.12.0], which deprecated
multiple elements of the API. Those elements have been removed in this release.

- **[BC]** Remove deprecated message routing methods (use `.Route()` instead)
  - `AggregateConfigurer.ConsumesCommandType()`
  - `AggregateConfigurer.ProducesEventType()`
  - `ProcessConfigurer.ConsumesEventType()`
  - `ProcessConfigurer.ProducesCommandType()`
  - `ProcessConfigurer.SchedulesTimeoutType()`
  - `IntegrationConfigurer.ConsumesCommandType()`
  - `IntegrationConfigurer.ProducesEventType()`
  - `ProjectionConfigurer.ConsumesEventType()`
- **[BC]** Removed `DescribableMessage` interface and `DescribeMessage()`
- **[BC]** Removed `ValidateableMessage` interface and `ValidateMessage()`
- **[BC]** Removed `EventRecorder` interface

## [0.12.1] - 2023-06-14

### Changed

- **[BC]** Application and handler identity names are now limited to 255
  bytes in length. This is a change to the specification/documentation only.

## [0.12.0] - 2023-04-09

This release aligns the Dogma API with best practices that have emerged since
the last release.

**Although this release includes a large number of changes there should be no
breaking changes to applications that are already following these best
practices.**

- Use [RFC 9562] UUIDs for identity keys
- Implement `MessageDescription()` on message types
- Implement `Validate()` methods on message types

Otherwise, most significant change is the introduction of `Routes()` methods to
handler configurer interfaces. Implementers should use `Routes()` in preference
to the existing `Consumes*Type()` and `Produces*Type()` methods, which are now
deprecated.

The `Routes()` API accepts arguments that use [type parameters] to communicate
message types. It also offers more extensible interface that allows future
support for per-message routing configuration without further breaking changes.

[type parameters]: https://go.dev/tour/generics/1
[rfc 9562]: https://www.rfc-editor.org/rfc/rfc9562.html

### Added

- **[BC]** Added `MessageDescription()` method to `Message` interface
- **[BC]** Added `Validate()` method to `Message` interface
- Added `Command`, `Event` and `Timeout` as aliases for `Message` in preparation for stricter static message typing

#### Routing API

- **[ENGINE BC]** Added `Routes()` methods to handler configurer interfaces
- Added `HandlesCommand()`
- Added `RecordsEvent()`
- Added `HandlesEvent()`
- Added `ExecutesCommand()`
- Added `SchedulesTimeout()`
- Added `Route` interface
- Added `AggregateRoute` interface
- Added `ProcessRoute` interface
- Added `IntegrationRoute` interface
- Added `ProjectionRoute` interface

#### Projection delivery policies

- Added `ProjectionConfigurer.DeliveryPolicy()`
- Added `ProjectionScope.IsPrimaryDelivery()`
- Added `ProjectionDeliveryPolicy`
- Added `UnicastProjectionDeliveryPolicy`
- Added `BroadcastProjectionDeliveryPolicy`

### Changed

- **[BC]** Handler and application identity keys must now be an RFC 9562 UUID string

### Deprecated

The new `Routes()` API supersedes the following methods:

- Deprecated `ConsumesCommandType()` methods, use `Routes()` with `HandlesCommand[T]()` instead
- Deprecated `ProducesCommandType()` methods, use `Routes()` with `ExecutesCommand[T]()` instead
- Deprecated `ConsumesEventType()` methods, use `Routes()` with `HandlesEvent[T]()` instead
- Deprecated `ProducesEventType()` methods, use `Routes()` with `RecordsEvent[T]()` instead
- Deprecated `SchedulesTimeoutType()` methods, use `Routes()` with `SchedulesTimeout[T]()` instead

Because `Message` now has `MessageDescription()` and `Validate()` methods, the
following elements are no longer necessary:

- Deprecated `DescribableMessage`
- Deprecated `DescribeMessage()`
- Deprecated `ValidatableMessage`
- Deprecated `ValidateMessage()`

No engines except [testkit] are able to provide a meaningful implementation of
`EventRecorder` without fundamental changes to the definition of an application.

- Deprecated `EventRecorder`, use an `IntegrationMessageHandler` instead

[testkit]: https://github.com/dogmatiq/testkit

## [0.11.1] - 2021-03-01

### Fixed

- Fix signature of `NoTimeoutMessagesBehavior.HandleTimeout()` to match `ProcessMessageHandler` interface

## [0.11.0] - 2021-02-23

### Added

- **[BC]** Add `ProjectionCompactScope.Now()`

### Changed

- **[BC]** `ProcessMessageHandler.HandleEvent()` now takes an `ProcessRoot` parameter
- **[BC]** `ProcessMessageHandler.HandleTimeout()` now takes an `ProcessRoot` parameter
- Handlers can now call `Process[Event|Timeout]Scope.ExecuteCommand()` and `ScheduleTimeout()` after `End()`

### Removed

- **[BC]** Remove `ProcessEventScope.Begin()`
- **[BC]** Remove `ProcessEventScope.HasBegun()` and `ProcessTimeoutScope.HasBegun()`
- **[BC]** Remove `ProcessEventScope.Root()` and `ProcessTimeoutScope.Root()`

## [0.10.0] - 2020-11-11

### Added

- **[BC]** Add `ProjectionMessageHandler.Compact()` and `NoCompactBehavior`
- Add `ValidatableMessage` interface and `ValidateMessage()`

## [0.9.0] - 2020-11-06

### Changed

- **[BC]** `AggregateMessageHandler.HandleCommand()` now takes an `AggregateRoot` parameter
- **[BC]** `fixtures.AggregateRoot` now stores all its historical events internally
- `AggregateCommandScope.Destroy()` no longer requires a prior call to `RecordEvent()`
- Handlers can now call `AggregateCommandScope.RecordEvent()` after `Destroy()`

### Removed

- **[BC]** Remove `AggregateCommandScope.Root()`

## [0.8.0] - 2020-11-03

### Changed

- Handlers can now call `AggregateCommandScope.Root()` for non-existent aggregate instances
- `AggregateCommandScope.Destroy()` is now a no-op for non-existent aggregate instances
- `AggregateRoot.ApplyEvent()` no longer has a requires an `UnexpectedMessage` panic

### Removed

- **[BC]** Remove `AggregateCommandScope.Create()`
- **[BC]** Remove `AggregateCommandScope.Exists()`
- **[BC]** Remove `StatelessAggregateRoot` and `StatelessAggregateBehavior`

## [0.7.0] - 2020-11-03

### Added

- **[BC]** Add `AggregateCommandScope.Exists()`
- **[BC]** Add `ProcessEventScope.HasBegun()` and `ProcessTimeoutScope.HasBegun()`

### Changed

- Allow engines to call `AggregateRoot.ApplyEvent()` with historical events
- Clarify semantics of `AggregateMessageHandler.New()` and `ProcessMessageHandler.New()`
- Clarify semantics surrounding creating an aggregate instance within the same scope that destroyed it
- Clarify semantics surrounding re-beginning a process instance within the same scope that ended it

## [0.6.3] - 2020-01-14

### Changed

- Clarify comparison semantics for identity names and keys

## [0.6.2] - 2019-12-09

### Fixed

- Exclude `fixtures.AggregateRoot.ApplyEventFunc` from JSON serialization

## [0.6.1] - 2019-11-19

### Added

- Add `DescribeMessage()` and the `DescribableMessage` interface
- Add the `fixtures` package, which contains message fixtures and mocks of Dogma interfaces

## [0.6.0] - 2019-08-01

### Changed

- **[BC]** `ProjectionMessageHandler` now uses an OCC strategy for event deduplication

### Removed

- **[BC]** Remove `ProjectionEventScope.Key()`

## [0.5.0] - 2019-07-24

### Added

- Applications and handlers are now assigned an immutable "key"
- **[BC]** Add `k` and `v` parameters to `ProjectionMessageHandler.HandleEvent()`
- **[BC]** Add `ProjectionMessageHandler.Recover()` and `Discard()`
- **[BC]** Add `ProcessMessageHandler.TimeoutHint()`
- **[BC]** Add `IntegrationMessageHandler.TimeoutHint()`
- **[BC]** Add `ProjectionMessageHandler.TimeoutHint()`
- **[BC]** Add `ProcessTimeoutScope.ScheduledFor()`
- **[BC]** Add `ProcessEventScope.RecordedAt()`
- Add `NoTimeoutHintBehavior`

### Changed

- **[BC]** Replace configurer `Name()` methods with `Identity()`
- **[BC]** Rename `NoTimeoutBehavior` to `NoTimeoutMessagesBehavior`
- **[BC]** Rename `ProjectionEventScope.Time()` to `RecordedAt()`

## [0.4.0] - 2019-04-17

### Added

- Document what strings constitute valid application and handler names
- **[BC]** Add `ProcessConfigurer.SchedulesTimeoutType()`

## [0.3.0] - 2019-02-26

### Added

- **[BC]** Require handlers to declare the message types they produce

### Changed

- **[BC]** Rename `RouteXXXType()` configurer methods to `ConsumesXXXType()`

## [0.2.0] - 2019-02-14

### Added

- **[BC]** Add `ProjectionEventScope.Key()` and `Time()`

## [0.1.0] - 2019-02-06

- Initial release

<!-- references -->

[unreleased]: https://github.com/dogmatiq/dogma
[0.1.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.1.0
[0.2.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.2.0
[0.3.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.3.0
[0.4.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.4.0
[0.5.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.5.0
[0.6.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.0
[0.6.1]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.1
[0.6.2]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.2
[0.6.3]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.3
[0.7.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.7.0
[0.8.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.8.0
[0.9.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.9.0
[0.10.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.10.0
[0.11.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.11.0
[0.11.1]: https://github.com/dogmatiq/dogma/releases/tag/v0.11.1
[0.12.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.12.0
[0.12.1]: https://github.com/dogmatiq/dogma/releases/tag/v0.12.1
[0.13.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.13.0
[0.13.1]: https://github.com/dogmatiq/dogma/releases/tag/v0.13.1
[0.14.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.14.0
[0.14.1]: https://github.com/dogmatiq/dogma/releases/tag/v0.14.1
[0.14.2]: https://github.com/dogmatiq/dogma/releases/tag/v0.14.2
[0.14.3]: https://github.com/dogmatiq/dogma/releases/tag/v0.14.3
[0.15.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.15.0

<!-- adr references -->

[ADR-21]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0021-remove-handler-timeout-hints.md

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Changed

- `AggregateCommandScope.Root()` can now be called for non-existent aggregate instances
- `AggregateCommandScope.Destroy()` is now a no-op for non-existent aggregate instances
- `AggregateRoot.ApplyEvent()` no longer has a hard requirement to panic with `UnexpectedMessage`

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
- Clarify semantics of surrounding creating an aggregate instance within the same scope as it was destroyed
- Clarify semantics of surrounding re-beginning a process instance within the same scope as it was ended

## [0.6.3] - 2020-01-14

### Changed

- Clarify comparison semantics for identity names and keys

## [0.6.2] - 2019-12-09

### Fixed

- Exclude `fixtures.AggregateRoot.ApplyEventFunc` from JSON serialization

## [0.6.1] - 2019-11-19

### Added

- Add `DescribeMessage()` and the `DescribableMessage` interface
- Add the `fixtures` package, which contains message fixtures and mocks of various Dogma interfaces

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
[Unreleased]: https://github.com/dogmatiq/dogma
[0.1.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.1.0
[0.2.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.2.0
[0.3.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.3.0
[0.4.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.4.0
[0.5.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.5.0
[0.6.0]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.0
[0.6.1]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.1
[0.6.2]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.2
[0.6.3]: https://github.com/dogmatiq/dogma/releases/tag/v0.6.3
[1.0.0-rc.0]: https://github.com/dogmatiq/dogma/releases/tag/v1.0.0-rc.0

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->

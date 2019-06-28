# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

- Applications and handlers are now assigned an immutable "key"

### Changed

- **[BC]** Replace configure `Name()` methods with `Identity()`

## [0.4.0]

### Added

- Document what strings consititute valid application and handler names
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

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->

# 2. Document API Changes

Date: 2018-12-07

## Status

Accepted

## Context

We need to advertise a meaningful history of changes to the Dogma API
specification for both application developers and engine developers.

The types of changes that have been made should be clearly identified, with
special attention drawn to changes that are not backwards compatible.

### Proposals

- Maintain a `CHANGELOG.md` as per the recommendations of [Keep a Changelog]
- Additionally, begin changelog entries that describe a BC break with `**[BC]**`
- Periodically tag releases, using [semantic versioning]

## Decision

A changelog will be maintained as per [Keep a Changelog]. Unreleased changes
should be added to the changelog as they are made.

Git tags will be named according to the rules of [semantic versioning].
Additionally, tag names are to be prefixed with a `v` as required by [Go modules].

## Consequences

Developers will have a single source of information about changes to the API.

Compatibility between versions can be clearly determined by examining the
version number.

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html
[Go modules]: https://github.com/golang/go/wiki/Modules#modules

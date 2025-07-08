# Architecture Decision Records

This project uses [Architecture Design Records], or ADRs, to keep track of the
decisions made about the design of the API. [adr-tools] is used to manipulate
the ADR documents.

<!-- references -->

[Architecture Design Records]: http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
[adr-tools]: https://github.com/npryce/adr-tools

## Index

* [1. Record architecture decisions](0001-record-architecture-decisions.md)
* [2. Document API changes](0002-document-api-changes.md)
* [3. Aggregate lifetime control](0003-aggregate-lifetime-control.md)
* [4. ADR process](0004-adr-process.md)
* [5. Routing of commands to processes](0005-routing-of-commands-to-processes.md)
* [6. Stateless aggregates and processes](0006-stateless-aggregates-and-processes.md)
* [7. Location of examples](0007-location-of-examples.md)
* [8. Location of testing features](0008-location-of-testing-features.md)
* [9. Immutable application and handler keys](0009-immutable-keys.md)
* [10. Handler timeout hints](0010-handler-timeout-hints.md)
* [11. Message timing information](0011-message-timing-information.md)
* [12. Comparison of Identifiers](0012-identifier-comparison.md)
* [13. Aggregate and process instance existence checks](0013-instance-exists-check.md)
* [14. Applying historical events to aggregate instances](0014-apply-historical-events-to-aggregates.md)
* [15. Routing unrecognized messages](0015-routing-unrecognized-messages.md)
* [16. Automatic aggregate creation](0016-automatic-aggregate-creation.md)
* [17. Recreation of aggregate instances after destruction](0017-recreate-aggregate-after-destruction.md)
* [18. Compacting projection data](0018-projection-compaction.md)
* [19. Automatic process creation](0019-automatic-process-creation.md)
* [20. Constraints on Identifier Values](0020-identifier-constraints.md)
* [21. Remove handler timeout hints](0021-remove-handler-timeout-hints.md)
* [22. Remove CRUD application support](0022-remove-crud-application-support.md)
* [23. Message order guarantees](0023-message-order-guarantees.md)
* [24. Permanently end processes](0024-permanently-end-processes.md)
* [25. Prevent reverting ended processes](0025-prevent-reverting-ended-processes.md)

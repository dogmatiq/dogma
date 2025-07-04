# Architecture Decision Records

This project uses [Architecture Design Records], or ADRs, to keep track of the
decisions made about the design of the API. [adr-tools] is used to manipulate
the ADR documents.

<!-- references -->

[Architecture Design Records]: http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
[adr-tools]: https://github.com/npryce/adr-tools

## Index

* [1. Record architecture decisions](0001-record-architecture-decisions.md)
* [2. Document API Changes](0002-document-api-changes.md)
* [3. Aggregate Lifetime Control](0003-aggregate-lifetime-control.md)
* [4. ADR Process](0004-adr-process.md)
* [5. Routing of commands to processes](0005-routing-of-commands-to-processes.md)
* [6. Stateless aggregates and processes](0006-stateless-aggregates-and-processes.md)
* [7. Location of examples](0007-location-of-examples.md)
* [8. Location of Testing Features](0008-location-of-testing-features.md)
* [9. Immutable Application and Handler Keys](0009-immutable-keys.md)
* [10. Handler Timeout Hints](0010-handler-timeout-hints.md)
* [11. Message Timing Information](0011-message-timing-information.md)
* [12. Comparison of Identifiers](0012-identifier-comparison.md)
* [13. Aggregate and Process Instance Existance Checks](0013-instance-exists-check.md)
* [14. Applying Historical Events to Aggregate Instances](0014-apply-historical-events-to-aggregates.md)
* [15. Routing Unrecognized Messages](0015-routing-unrecognized-messages.md)
* [16. Automatic Aggregate Creation](0016-automatic-aggregate-creation.md)
* [17. Recreation of Aggregate Instances After Destruction](0017-recreate-aggregate-after-destruction.md)
* [18. Compacting Projection Data](0018-projection-compaction.md)
* [19. Automatic Process Creation](0019-automatic-process-creation.md)
* [20. Constraints on Identifier Values](0020-identifier-constraints.md)
* [21. Remove handler timeout hints](0021-remove-handler-timeout-hints.md)
* [22. Remove CRUD application support](0022-remove-crud-application-support.md)
* [23. Message Order Guarantees](0023-message-order-guarantees.md)
* [24. Permanently End Processes](0024-permanently-end-processes.md)
* [25. Prevent Reverting Ended Processes](0025-prevent-reverting-ended-processes.md)

# 31. Require retries for idempotency-keyed commands

Date: 2026-04-06

## Status

Accepted

## Context

The `WithIdempotencyKey()` option allows callers to attach a stable key to a
command passed to `CommandExecutor.ExecuteCommand()`. Engines can use that key
to detect duplicate attempts to execute the same operation.

Retries and idempotency keys are coupled. Retries are only safe when the caller
reuses the same idempotency key, and an idempotency key is only useful when the
caller can retry a failed attempt. This pairing allows an engine to detect
duplicate retries without adding extra synchronous recovery work to
`CommandExecutor.ExecuteCommand()`.

Our API documentation currently frames retry behavior as guidance. That wording
is not strong enough for designs that rely on caller retry as the only recovery
path for idempotency-keyed submissions.

## Decision

We will define a stronger contract for commands submitted with
`WithIdempotencyKey()`.

By providing an idempotency key, the caller accepts responsibility for retrying
failed submissions with the same key.

Engine implementations may rely on that caller behavior as the sole recovery
mechanism for idempotency-keyed commands. If the caller does not retry after a
failure during acceptance, the command may be silently lost.

## Consequences

The role of `WithIdempotencyKey()` becomes explicit: it identifies repeated
submissions and establishes a retry responsibility for the caller.

This ADR does not define what makes two attempts the same business operation,
or what business outcome counts as success. Applications define those
semantics.

Engine implementations are free to optimize acceptance for idempotency-keyed
commands around caller retry, instead of also maintaining an engine-managed
recovery path for the same submissions.

Callers that provide an idempotency key but do not retry failed submissions
accept a risk of silent command loss.

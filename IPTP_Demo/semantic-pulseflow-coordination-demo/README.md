# PulseFlow Coordination Demo

A minimal runnable demo showing how independently defined services can coordinate through explicit state rather than direct calls or a central workflow.

This demo is intentionally single-process. Its purpose is to show the semantic model clearly: services react to explicit state rather than direct calls or central orchestration. It does not yet introduce inter-process exchange; that is the next step, where IPTP can carry the same intention signals between independently running components.

## The problem

In many systems where components are developed independently, the logic that determines **when something should happen** ends up:

- spread across code
- hidden in assumptions
- hard to inspect or debug over time

## The idea

Instead of services calling each other directly, each service reacts only to explicit state in a shared field.

In this demo:

- `auth_service` reacts to `login_submitted:Y`
- `cart_service` reacts to `cart_submitted:Y`
- `order_service` reacts only when both `user_authenticated:Y` and `cart_valid:Y` are present

That means the triggering conditions are visible in data rather than buried in control flow.

## Why the pass loop matters

The pass loop is not an orchestrator telling services what to do.

It simply visits each service and lets that service decide whether it is ready based on the current field state.

This means:
- services do not depend on direct calls from other services
- readiness depends on explicit state, not service position
- even when services are listed in an awkward order, execution still follows semantic dependency

The loop provides opportunity, not orchestration; the field determines readiness.

## Run

```bash
npm start
```

or:

```bash
node run.js
```

## What this demonstrates

- no service calls another service
- no central orchestration logic
- coordination happens through explicit state transitions

## Key line

> Each service reacts to state — not to other services.

## Relation to IPTP

This demo is the execution-side companion to IPTP.

- **IPTP** explores how a message can carry explicit intent + state markers
- **PulseFlow demo** shows how independently defined services can react to that explicit state

The protocol alone is abstract; this demo makes the coordination pattern visible.

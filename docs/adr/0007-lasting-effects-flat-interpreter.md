# Lasting "remainder of the turn" effects as a flat interpreter

## Context

Some effects last "for the remainder of the turn" and fire on a later event: Full
Moon gains Æmber whenever you play a creature, Charge! deals damage whenever you
play a creature, Crystal Hive gains Æmber whenever a creature reaps, Dimension Door
makes reaping steal instead of gain. The natural implementation — store the
deferred `Effect` (a closure) in state and run it when the event happens — is
impossible, because state is flat and pointerless (ADR 0005) and cannot hold a
`func` or `Effect`. The lazy alternative — a bespoke `if g.State.Foo …` block in
the play/reap path for each such card — hardcodes every effect into the hot path
and makes effects that share a timing window impossible to order.

## Decision

Route every lasting effect through a **flat registry** that is a second,
deliberately smaller interpreter. State holds enum-tagged records
`LastingEffect{On Event, Do lastingAction, Controller, Amount}`; `lastingActionOf`
translates a composed effect into an action tag, and `game_lasting.go` fires and
queries the records. Two flavors:

- A **reaction** runs after an event — `AddLasting` registers it, `emitLasting` at
  the event site gathers every reaction the actor owns and, when several fire at
  once, lets the active player order them.
- A **replacement** changes an event's own outcome — the event site queries
  `lastingReplacement` and applies `Instead{Of, With}` in place.

## Consequences

- Simultaneous triggers on one event order for free (the active player orders
  them), because they all route through one dispatch.
- Adding a reaction on an existing event = support its `Do` in `lastingActionOf`
  and `resolveReaction`; a new event = one `Event` value, one
  `emitLasting`/`lastingReplacement` call at the site, and its text. The play/reap
  path's structure never changes.
- This is a real sub-language, not an ad-hoc pile of `if`s — keep the enum dispatch
  centralized; if the action set grows, formalize it as a tiny instruction set with
  one `apply(action, amount, ctx)` switch.

This registry exists because of ADR 0005; it is the flat-state answer to "do X
later".

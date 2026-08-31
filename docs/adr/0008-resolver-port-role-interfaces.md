# The Resolver port, segregated into role interfaces

## Context

Effects have to change the game — deal damage, move cards, gain Æmber, forge keys.
If they reached into `*Game` or `GameState` directly, the card domain would be
coupled to the runtime and to the flat-state internals (ADR 0005), every card
could touch anything, and there would be no single place that answers "what can a
card actually do?"

## Decision

Effects touch the game **only** through a `Resolver` port, held via
`EffectContext` — never `*Game` or `GameState`. This is Ports & Adapters: the
effect AST is the domain, `Resolver` is the port, `*Game` is the adapter. Every
change a card can make is one `Resolver` method, so the interface is the complete,
auditable catalogue of card-facing capability. `Resolver` is intentionally wide,
so it is **composed from focused role interfaces** rather than being one flat list:
`StateReader`, `EconomyResolver`, `CreatureResolver`, `CombatResolver`,
`ZoneResolver`, `TurnResolver`, `ChoiceResolver`, `Logger`. A new capability is a
method on the role it belongs to; a genuinely new cluster becomes a new role
embedded in `Resolver`.

## Consequences

- `Resolver` is the single, reviewable list of everything a card can do.
- An effect — or a test double — depends on just the role it needs, not the whole
  runtime.
- The wide interface is a deliberate deviation from Go's small-interface taste,
  paid for by the role segregation. Add a method to the right role; do not append
  to a flat list, and let a truly new cluster of capability become a new role.

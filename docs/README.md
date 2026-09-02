# Vactrol documentation

Index of every doc in the repo. Start with **architecture** for how the code fits
together and **CONTEXT** for what the words mean.

## Start here

- [architecture.md](architecture.md) — how the codebase fits together: the
  pointerless state, the card-effect AST, and how a turn and an ability flow.
- [../CONTEXT.md](../CONTEXT.md) — the glossary of vocabulary shared across the
  engine, card database, and deck generation.
- [../README.md](../README.md) — project overview, quick start, and layout.

## Design & decisions

- [deck-generation.md](deck-generation.md) — the procedural deck-generation design
  (philosophy, pipeline, distribution, scoring).
- [style-guide.md](style-guide.md) — coding style, composition, and naming.
- [roadmap.md](roadmap.md) — the long-term vision and phased plan.
- [adr/](adr/) — Architecture Decision Records (the "why" behind hard-to-reverse
  choices):
  - [0001](adr/0001-intrusive-linked-lists-for-extensible-flat-state.md) —
    intrusive linked lists for unbounded per-card data in flat state
  - [0002](adr/0002-per-zone-capacities.md) — per-zone capacities, no shared arena
  - [0003](adr/0003-deckgen-inverted-dependency.md) — deck generation is a pure
    function below the card facade
  - [0004](adr/0004-template-materialize-seam.md) — card templates materialize to
    concrete defs at generation time
  - [0005](adr/0005-flat-pointerless-comparable-state.md) — flat, pointerless,
    comparable game state (the keystone)
  - [0006](adr/0006-effects-render-and-resolve-themselves.md) — effects render
    their own text and carry themselves out
  - [0007](adr/0007-lasting-effects-flat-interpreter.md) — lasting effects as a
    flat interpreter
  - [0008](adr/0008-resolver-port-role-interfaces.md) — the Resolver port,
    segregated into role interfaces
  - [0009](adr/0009-tactic-type-omni-as-versatile.md) — card type "Tactic"; Omni
    as Versatile + Action
  - [0010](adr/0010-invalid-zero-validated-at-init.md) — invalid-zero sentinels
    validated at card init
  - [0011](adr/0011-logs-narrate-resolved-outcomes.md) — the game log narrates
    resolved outcomes, not card text
  - [0012](adr/0012-first-class-turn-phases.md) — turn phases are first-class
    engine state
  - [0013](adr/0013-end-of-turn-last-and-ordered-triggers.md) — end-of-turn
    abilities resolve last; simultaneous triggers are ordered

## Rules

- [rulebook.md](rulebook.md) — Vactrol's own rulebook, **generated** from the
  engine (`mage gen`); the authoritative source for how mechanics work.
- [keyforge-master-rulebook.md](keyforge-master-rulebook.md) — the official
  KeyForge rulebook, kept as a faithful reference only.
- [card-wording-rules.md](card-wording-rules.md) — the curated conventions every
  card's printed text must follow.

## Contributing

- [testing.md](testing.md) — the testing layers and what to test where.
- Agent/contributor rules live in the `AGENTS.md` files:
  [root](../AGENTS.md), [internal/engine](../internal/engine/AGENTS.md),
  [internal/cards](../internal/cards/AGENTS.md).

## Planning

- [todo.md](todo.md) — the running list of things to build and open design
  questions.

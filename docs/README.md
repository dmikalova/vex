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
  - [0014](adr/0014-style-gallery-on-real-components.md) — the Style gallery
    renders real components, gated at runtime
  - [0015](adr/0015-selected-card-lifts-an-enlarged-copy.md) — clicking a card
    lifts an enlarged copy with its actions on it
  - [0016](adr/0016-under-a-per-host-facedown-capable-attachment.md) — under, a
    per-host, out-of-play, facedown-capable attachment
  - [0017](adr/0017-one-deckgen-set-per-source-set-with-legacy-pool.md) — one
    deck-generation Set per source set, other sets as its legacy pool
  - [0018](adr/0018-rulebook-term-registry.md) — the rulebook is a typed term
    registry, complete by construction
  - [0019](adr/0019-controlled-rules-voice.md) — a controlled Rules voice for
    every player-facing surface
  - [0020](adr/0020-controlled-code-comments.md) — code comments follow the same
    plain, controlled style
  - [0021](adr/0021-reprints-are-full-set-members.md) — a reprint is a full member
    of the reprinting set's deck-generation pool
  - [0022](adr/0022-iconography-is-a-visitor-not-an-effect-method.md) — card
    iconography is a Visitor over the effect AST, not a node method
  - [0023](adr/0023-card-gallery-served-page-with-in-client-filtering.md) — the
    card gallery is a served page with in-client filtering
  - [0024](adr/0024-generic-counters-global-side-table.md) — generic card counters
    live in a global side-table, not a field per kind
  - [0025](adr/0025-deck-list-from-retained-generated-roster.md) — the deck list is
    rendered from the retained generated roster
  - [0026](adr/0026-creatures-playable-as-upgrades.md) — creatures can be played as
    upgrades
  - [0027](adr/0027-houses-alphabetical.md) — houses are ordered alphabetically
  - [0028](adr/0028-take-control-lifo-stack.md) — take control is a LIFO stack of
    control effects, not one source per card
  - [0029](adr/0029-settle-at-resolution-boundaries.md) — destruction settles at
    resolution boundaries, not after every state change
  - [0030](adr/0030-a-card-out-of-play-takes-no-further-part.md) — a card out of
    play takes no further part
  - [0031](adr/0031-zone-movement-is-one-mechanism.md) — zone movement is one
    mechanism; the KeyForge verbs are sugar over it
  - [0033](adr/0033-runtime-type-conversion-via-lastingtype.md) — a card's runtime
    type is `TypeOf`, overridable in play via `LastingType`

## Rules

- The **rulebook** lives as the web client's `/rulebook` and `/glossary` pages,
  rendered live from the engine's typed rulebook term registry (ADR 0018); the
  authoritative source for how mechanics work.
- [keyforge-master-rulebook.md](keyforge-master-rulebook.md) — the official
  KeyForge rulebook, kept as a faithful reference only.
- [card-wording-rules.md](card-wording-rules.md) — the curated conventions every
  card's printed text must follow.
- [keyforge-divergences.md](keyforge-divergences.md) — the Vactrol⇄KeyForge
  divergence register: where and why Vactrol departs from KeyForge.

## Contributing

- [testing.md](testing.md) — the testing layers and what to test where.
- Agent/contributor rules live in the `AGENTS.md` files:
  [root](../AGENTS.md), [internal/engine](../internal/engine/AGENTS.md),
  [internal/cards](../internal/cards/AGENTS.md).

## Planning

- [todo.md](todo.md) — the running list of things to build and open design
  questions.
- [todo-agent.md](todo-agent.md) — the agent's scratchpad of concrete work items.
  Agents write here, never in `todo.md`; a done item is deleted, not marked done,
  so the file always shows only outstanding work.

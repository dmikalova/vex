# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Cite the ADR or doc that decided an item where one exists.

## Engine refactor sweep (survey of `internal/engine`)

File-naming decision for the splits: keep each family's existing top-level prefix
so `ls` groupings stay intact — `target.go` is a bare concept file, so its splits
are `target_*.go`; `effect_condition.go` / `effect_count.go` are part of the
`effect_*` vocabulary, so their splits keep the prefix (`effect_condition_*.go`,
`effect_count_*.go`). Ratchet the split convention into `internal/engine/AGENTS.md`.

Structural (decompose / atomize / recompose):

- **A3** — Selection mode (chosen / random / each) is a `Selector` Strategy, not a
  node per (verb × zone × mode): fold `PurgeFromHand` / `PurgeRandomFromHand` /
  `PurgeEachFromHand` / `PurgeCreatureFromHand` (and the archive/discard mirrors)
  onto one node per verb carrying a selector. **Do all zones**, not just hand.

Cleanup (splits / comments / naming — not method refactors):

- **B2** — Split `target.go` (1479 lines) into `target_*.go` (type + kinds vs
  selection/filter helpers).
- **B4** — Split `resolver.go` (1246) and `text.go` (1106): separate the role-
  interface declarations from the `*Game` method bodies; group the text helpers.

## Card wording / authoring

- **Orator Hissaro** could read: "Play: Exalt and ready each neighboring creature.
  For the remainder of the turn, those creatures belong to house Saurian."
  Blocked on two engine additions: making `Exalt` a `combinable` (so `Ready`+`Exalt`
  fold to "ready and exalt each neighboring creature" — but the fold must keep
  "exalt N times" for `Amount > 1`), and a pronoun form of `BelongToHouse` so the
  second sentence reads "those creatures belong to …" instead of repeating the
  target. Today it renders correctly but verbosely (target repeated three times).
- **Borr-Nit** and similar could be atomized and recomposed further (decompose
  fused effects into shared nodes).

## Card catalog / provenance

_No outstanding items._

## Web — mobile, previews, layout

_No outstanding items._

## Tooling / tests

_No outstanding items._

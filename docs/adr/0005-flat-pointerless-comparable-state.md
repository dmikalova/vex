# Flat, pointerless, comparable game state

## Context

The engine is built to be searched: an MCTS/self-play bot needs to clone a
position, explore, and discard it, millions of times. A conventional mutable
object graph — `GameState` as a web of pointers, slices, and maps — makes each
clone a deep copy with allocation and GC pressure, and makes two positions
impossible to compare with `==`. That cost sits on the hottest path in the whole
system.

## Decision

`GameState` is a **value struct of fixed-size arrays** (`[maxCards]CardCore`,
`[2]Zone`, …) with **no pointers, slices, or maps**, so `GameState.FastCopy()` is
a plain value copy with no allocation. Cards are referenced by `LocalID uint8`
into a read-only `catalog` of definitions held **separately** from `GameState` and
shared across every clone. Value types that are embedded in or compared against
state stay **comparable**: `Target` (compared `== Target{}` by `ConstantAbility`)
is a wide flag struct rather than a `[]filter`, and an "optional int" is a paired
`x int` + `hasX bool`, never a `*int`.

## Consequences

- `FastCopy` is allocation-free, so cloning a position for search is essentially
  free — the property the rest of the engine bends to preserve.
- State **cannot hold a closure or an `Effect`/`func` value**, so "do X later" has
  to be flat, enum-tagged data. That constraint is what forces the lasting-effects
  registry (ADR 0007).
- `Target` is a comparable flag-soup, not a slice of filters; per-card filters go
  through builder methods and set-relative rules through a `Selector`, not new
  slice fields.
- Fixed capacities mean zones are sized to their own bounds (ADR 0002) and an
  unbounded per-card collection (a creature's upgrades) needs an intrusive linked
  list rather than a slice (ADR 0001).
- Keeping definitions in a separate read-only catalog keeps them out of every
  copy.

This is the keystone constraint of the engine: ADRs 0001, 0002, and 0007 are all
consequences of it.

# 24. Generic counters live in one global side-table, not a field per kind

## Context

A growing family of cards place a card-specific marker on a creature or artifact
in play: a doom counter (Wretched Doll), and — across the sets still to
implement — mineralize, corrosion, ignorance, mutation, paint, glory, time,
warrant, fuse, and Election's Yea/Nay tallies, among others. Each is the same
shape: a named token that sits on an in-play card, stacks (a second doom counter
on the same creature is legal), does nothing on its own, and is read only by the
cards that care about it.

Doom counters were first modelled the obvious way: a bespoke `DoomCounters int16`
field on `CardCore`. Repeating that per kind does not scale. `CardCore` is copied
whole into every undo/MCTS snapshot, `maxCards = 128` of them per state, so each
`int16` field costs `2 * 128 = 256` bytes per snapshot. Two dozen kinds would add
several kilobytes to a `GameState` that is otherwise ~6 KB, roughly doubling the
cost of the value copy that makes undo and MCTS cheap (ADR 0005). A dense
`[NumKinds]uint8` per card has the same problem — it pays for every kind on every
card whether or not the counter is present.

The saving grace is sparsity. A card rarely carries more than one or two kinds,
and a whole game very rarely has more than two kinds of generic counter in play
at once. Storing the counters that actually exist, rather than a slot for every
kind on every card, is far cheaper.

Nothing in the engine dedups logically-equal states, so a sparse, order-by-
construction store is safe: only the web tests compare a state to its own
snapshot, and a deterministic array survives a value copy unchanged. This is the
same reasoning that lets zone lists and the upgrade chain be order-by-
construction (ADR 0001).

## Decision

Generic counters live in one **global side-table** on `GameState`, not in
`CardCore`:

```go
type CounterEntry struct {
    Card LocalID    // the in-play card the counter sits on
    Kind CounterKind
    N    uint8      // how many, saturating at 255
}

// on GameState:
Counters     [64]CounterEntry
CounterCount uint8
```

- **One entry per `(card, kind)` pair**, with `N` folding the count in. Two doom
  counters on one creature is `{card, Doom, 2}`, not two entries; thirty time
  counters is one entry. This is what keeps the table small even for the cards
  that pile counters high (time counters have been seen near thirty; Election
  holds up to six Yea and six Nay at once).
- **Compacted on removal** with the existing shift-left `listRemoveAt` pattern —
  no tombstones, order preserved by construction.
- **Shed on leave-play** in the single `removeFromPlay` funnel every exit passes
  through, dropping every entry for that card. Counters only ever sit on in-play
  cards, so nothing else has to clear them.
- **Reads are O(M ≤ 64) scans.** We spend a little compute to save the space; the
  table is tiny and the scan is a linear walk over a flat array.
- **Overflow panics.** A 65th distinct `(card, kind)` pair is a caught invariant
  violation the sim and fuzz harness surface, never a silent drop. `M = 64`
  (192 bytes + 1) is comfortably above any real board.

Counter **identity** is one global `CounterKind uint8` enum. Yea and Nay are
separate kinds — Election holds both at once. Each kind carries a display noun
for rendering card text (`"doom counter"`); that is its only per-kind data.

Doom counters are the first client: the bespoke `DoomCounters` field is deleted
and re-expressed as `CounterKind.Doom`, and the doom-specific engine surface is
generalized — `PlaceDoomCounter` → `PlaceCounter{Kind, Target, Amount}`,
`DoomCounterInPlay` → `CounterInPlay{Kind}`, `Target.WithDoomCounter()` →
`Target.WithCounter(kind)`, and the `AddDoomCounter`/`DoomCounterOn` resolver
methods → `PlaceCounter`/`CountersOn`.

`PowerCounter`, `Damage`, and `Amber` stay bespoke `CardCore` fields. They are not
generic markers: power counters change a creature's power, damage drives
destruction, and Æmber-on-card is captured/exalted and released to a player when
the card leaves play. They are read on the hot path by identity, not by a card
that names them.

### What the counter layer does not own

Some cards say a marker also grants a static property while it is present —
mineralize makes a creature an artifact, ignorance blanks its text, mutation
gives it the Mutant trait, paint makes it a Mars creature. That static grant is a
separate `ConstantAbility` composed on top of a `CounterInPlay{Kind, Target:
This}` read (ADR-style composition), not something the counter table knows about.
The counter layer owns only place, remove, and read.

### Rulebook

One rulebook section, **"Generic Counters"**: a counter is a card-placed marker
whose meaning is defined entirely by the card that places it, followed by the
flat list of counter names. There is **no per-kind rulebook term** — a new
`CounterKind` builds without a rulebook entry, unlike a keyword or card type
(ADR 0018). The display noun each kind carries is for card-text rendering only.

## Consequences

- Adding a counter kind costs one enum value and its display noun. Snapshot size
  is unchanged; the table only grows with counters that actually exist on the
  board, and only when they do.
- Reads walk a ≤64-entry array instead of dereferencing a field. This is a
  deliberate compute-for-space trade; the array is small and the walk is linear.
- The store is order-by-construction and safe to value-copy, consistent with the
  zone lists and upgrade chain (ADR 0001, ADR 0005). It must never be made to
  dedup or reorder for "canonical" equality — nothing depends on that, and it
  would cost more than it saves.
- Leave-play cleanup lives in exactly one place (`removeFromPlay`). Because the
  store is off `CardCore`, `resetCore`'s "zero the struct, every field resets for
  free" guarantee no longer covers counters; the funnel clears them explicitly.
- A pathological 65th distinct `(card, kind)` pair panics rather than silently
  dropping a counter. If a real card ever approaches the bound, raise `M`.

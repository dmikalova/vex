# Per-zone capacities (36 / 72), no shared arena

## Context

Each per-player zone (`Hand`, `Deck`, `Battleline`, `Discard`, `Artifacts`,
`Archives`, `Purge`) is a fixed-capacity `LocalID` array plus a count, held by
value in `GameState` so a position copies with no allocation — the invariant
behind cheap cloning for AI search. The zones originally shared one `zoneCap = 40`
constant. That was wrong in two directions at once: it over-allocated the zones
that can only ever hold a single 36-card deck, and it _under_-allocated the zones
that can hold cards from both decks (a player can control an opponent's creatures
and artifacts, or archive cards from either deck), which can reach 72 — a latent
overflow, since `add` writes past the end with no bound check. A shared arena (one
backing store all zones index into, or a per-card "which zone am I in" tag with
zones as filtered views) would shrink the total further, since a `LocalID` lives
in exactly one place at a time.

## Decision

Size each zone to its own bound rather than a single global cap, and do **not**
introduce a shared arena. Deck-bounded zones (deck, hand, discard, purge — they
only ever receive their owner's cards) use a 36-slot `deckList`; zones that can
hold cards from both decks (battle line, artifact row, archives) use a 72-slot
`wideList`. Both are plain fixed-array value types; because Go generics cannot
slice a type parameter whose array sizes differ, the list logic lives in shared
free functions and each type is a thin wrapper.

## Consequences

- Correct by construction: the wide zones can no longer overflow at 40, and the
  deck-bounded zones stop reserving slots they can never use.
- The zones keep O(1) indexed access, in-place ordering, and shuffling — which an
  arena or per-card zone tag would give up (an arena needs compaction on every
  move; a tag turns ordered operations into O(n) scans). Unlike upgrades (ADR
  0001), the ordered zones are neither sparse nor order-loose, so the linked-list
  trick that fit upgrades does not fit here.
- Two list types cost a little duplication (a thin method wrapper per type) but
  keep the flat, pointerless, comparable-value layout intact.
- A full arena remains the reserved option if per-zone sizing is ever
  insufficient; nothing here forecloses it.

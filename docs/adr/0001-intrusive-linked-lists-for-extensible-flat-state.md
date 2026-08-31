# Intrusive linked lists for unbounded per-card data in flat state

## Context

`GameState` is a flat, pointerless, comparable value (a struct of fixed-size
arrays) so a position copies with no allocation — the load-bearing invariant
behind cheap cloning for the search bot (ADR 0005). Because state holds no slices, maps, or
pointers, any "a card owns a variable number of other cards" relationship has to
be represented some other way. The first cut modelled attached upgrades as a
fixed `[6]LocalID` array plus a count on each card. A soak/fuzz run then played a
creature a seventh upgrade — legal in KeyForge, which sets no upgrade limit — and
the write ran off the end of the array and panicked.

## Decision

Store unbounded, order-loose, parent-bounded per-card collections as an
**intrusive singly-linked list threaded through `CardCore`**, using `+1`
encoding (0 = "none", value `n` = `LocalID n-1`) for the link fields — the same
idiom as `ControlPlus`. Upgrades use three bytes: `FirstUpgradePlus` (the head of
a host's chain), `NextUpgradePlus` (the next sibling on the same host), and
`HostPlus` (an upgrade's back-link to its host). Attaching threads onto the tail
(preserving attach order); detaching stitches the neighbours back together so the
chain stays whole when an element in the middle leaves play. The linkage and its
no-allocation iterator live in `internal/engine/game_upgrades.go`.

## Consequences

- No arbitrary cap, and the state stays flat, pointerless, and comparable — a
  slice or pointer would have broken the copy-is-a-value-copy invariant, and a
  fixed array either caps the game illegally or wastes space on every card.
- The zero value means "unattached" for free (no separate count field to keep in
  sync), and the back-link makes "which creature am I on?" O(1).
- Reads are O(k) in the chain length and must go through the
  `firstUpgrade`/`nextUpgrade` iterator rather than indexing; there is no random
  access. This is the right shape only for collections that are **sparse,
  order-loose, and bounded to one parent** — upgrades fit; the ordered zones
  (deck/hand/discard), which need shuffling and index access, deliberately do
  not and stay fixed-capacity arrays.

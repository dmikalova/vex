# One deck-generation Set per source set, cross-set cards as its legacy pool

## Context

deckgen was built around a single `Set` — the pool a deck draws from (ADR 0003,
ADR 0004). Its design always anticipated **legacy cards**: KeyForge occasionally
fills a slot with a card from a _different_ set that shares the slot's House, and
the pipeline reserved a per-slot `Legacy` flag and a `Tuning.LegacyRate` for it.
That seam sat inert while only one set (Call of the Archons) existed — there was
no "other set" to draw a legacy card from.

With a second set (Age of Ascension) now registering cards, the seam has to be
wired: a slot in a CotA deck may roll an AoA card of the same House, and vice
versa. Two questions had to be answered together — where the set boundary comes
from, and how a legacy card differs from a maverick (a card _rehoused_ into a
slot's House, ADR 0004). A legacy card is **not** rehoused: it keeps its own
printed House and only fills a slot of that same House in another set's deck.

## Decision

A **source set is one `internal/cards/sets/<slug>` package**, identified by a
card's first provenance `Ref` (`rc.Provenance[0].Set.Name`). The `cards`
aggregator groups every registered card by that name (`bySet`), walks the sets in
release order (`provenance.Sets()`), and builds **one `deckgen.Set` per source
set** — `DeckSets()`. Each set's own cards form its pool; the cards of _every
other_ set become its **legacy pool**, attached with a new builder step:

```go
deckgen.NewSet(name, own, deckgen.DefaultTuning()).WithLegacy(legacy)
```

`WithLegacy` buckets the legacy cards by House (`legacyByHouse`), **skipping**
houseless Specials, `Connected` cards, and `HouseNone` — the same cards the main
pool excludes — and keeps each legacy card's own House (no rehousing). During
generation, `fillSlot` rolls `Tuning.LegacyRate` per slot; on a hit it draws from
`legacyByHouse[house]` and commits the slot with `Legacy: true`, leaving the
card's House untouched.

`DeckSet()` — the single-set entry point every current caller (`match`, `sim`)
uses — returns `DeckSets()[0]`, the **first released set** (CotA). That set is the
fully-implemented base, so generated decks are always full even while later sets
are only partially implemented; the later sets ride along only as legacy cards.

## Consequences

- The set boundary is provenance-derived, so adding a set is just adding its
  package and cards — `DeckSets()` picks it up, gives it a pool, and folds it into
  every other set's legacy pool with no wiring.
- Legacy and maverick stay cleanly distinct: maverick rehouses (`def.House =
slotHouse`) and routes through `Materialize` (ADR 0004); legacy keeps its House
  and is a pool-membership decision made in `fillSlot`. They combine freely, as the
  design always intended (a legacy-maverick is emergent).
- `fillSlot` now draws an extra RNG value per slot for the legacy roll, which
  shifts the deterministic `(set, seed) → deck` sequence. There are no golden-deck
  snapshots, so this is invisible today; if one is ever added it must be
  regenerated against the post-legacy stream.
- `DeckSet()` deliberately does **not** track the newest set. Until a later set is
  fully implemented, defaulting to it would deal near-empty decks; keeping the base
  set as the default trades "newest content" for "always a full, legal deck." When
  a later set is complete, this is the one line to revisit.
- The choice is scoped to _legacy cards_ (same-House, other-set). The rarer
  house-level overlays the design notes anticipate — a whole deck-House whose pool
  or House comes from another set (legacy/maverick _Houses_) — remain a separate,
  still-deferred seam.

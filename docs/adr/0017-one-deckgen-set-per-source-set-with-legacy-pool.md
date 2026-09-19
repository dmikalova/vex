# One deck-generation Set per source set, cross-set cards as its legacy pool

> **Status: superseded in part by [ADR 0021](0021-reprints-are-full-set-members.md).**
> This ADR treats a set's own cards as its pool and _every other set's_ cards as a
> same-House **legacy** pool drawn on a per-slot `Tuning.LegacyRate` roll. ADR 0021
> then decided that a card **reprinted** into a set is a full member of that set's
> own pool, not a legacy draw — so a reprint is native to every set that prints it,
> and legacy stays reserved for cards that a set never reprinted. Read the legacy
> mechanism below as it still stands; read "each set's own cards form its pool"
> through ADR 0021's lens (own pool = printed **or reprinted** here).

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
card's declared **home set** (`rc.Set`). The `cards`
aggregator groups every registered card by that set (`bySet`), walks the sets in
release order (`provenance.Sets()`), and builds **one `deckgen.Set` per source
set** — `DeckSets()`. Each set's own cards form its pool; the cards of _every
other_ set are its **legacy pool**, supplied as one shared `Legacy` value built
once for the whole catalog and attached to every set with a builder step:

```go
deckgen.NewSet(name, own, deckgen.DefaultTuning()).WithLegacy(shared)
```

> **Amended:** membership was originally keyed on a card's _first provenance_
> `Ref` (`rc.Provenance[0]`). It is now keyed on an explicit home set — `card.InSet`,
> defaulting to the first provenance tag when a card declares none — so deck
> generation **never reads `Provenance`** (which is pure coverage bookkeeping) and a
> Vactrol-invented card with no provenance (Master of 4/5, an Anomaly Expansion
> Shard) still joins a pool. `setName` reads `rc.Set`; `nativeSet` resolves it.

`NewLegacy` buckets every registered card by House **and rarity** (with a flat
per-House fallback), tagging each with the set it came from and **skipping**
houseless Specials, `Connected` cards, and `HouseNone` — the same cards the main
pool excludes — and keeps each legacy card's own House (no rehousing). The one
`*Legacy` is shared by every set rather than copied per set; a set draws only the
entries whose set differs from its own Name, so a legacy slot pulls a card printed
in another set. During generation, `fillSlot` rolls `Tuning.LegacyRate` per slot;
on a hit it draws a legacy card of the slot's **rolled rarity** for its House,
falling back to any rarity of that House when it has no legacy card of that
rarity, and commits the slot with `Legacy: true` at the drawn card's own rarity,
leaving its House untouched.

`DeckSet()` — the single-set entry point every current caller (`match`, `sim`)
uses — returns `DeckSets()[0]`, the **first released set** (CotA). That set is the
fully-implemented base, so generated decks are always full even while later sets
are only partially implemented; the later sets ride along only as legacy cards.

## Consequences

- The set boundary is provenance-derived, so adding a set is just adding its
  package and cards — `DeckSets()` picks it up, gives it a pool, and folds it into
  every other set's legacy pool with no wiring.
- Legacy and maverick stay cleanly distinct: maverick rehouses (`def.House =
slotHouse`) and routes through`Materialize` (ADR 0004); legacy keeps its House
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

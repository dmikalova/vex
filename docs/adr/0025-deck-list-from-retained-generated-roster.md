# The deck list reads a retained generated roster, not the live piles

## Context

A player's **Deck** is the 36-card result of generating from one Set with one
seed — 3 House pods of 12 Slots, each Slot carrying an intrinsic Rarity and its
Maverick/Legacy provenance (ADR 0004, ADR 0017). We want a **Deck list**: a
public, always-open view of that full roster, reached from a deck icon on the
Player bar, laid out as one column per House with each card's type and Rarity.

The roster is generated once at match setup and then thrown away.
`match.SetupDecksFor` calls `deckgen.Generate(set, seed)`, adds each Slot's card
to the engine deck, and returns only three derived slices — the deck's Houses,
its Maverick ids, and its Legacy ids — discarding the ordered `deckgen.Deck`
itself. So the ordered roster (which House holds which cards, in which Slots, at
what Rarity) exists nowhere after setup.

Two tempting shortcuts both leak or lie:

- **Regenerate on demand.** `Generate(set, seed)` is reproducible only within a
  single version of a Set's pool (ADR 0017's note, and the `Deck` doc comment).
  As cards are implemented the pool shifts, so regenerating later can yield a
  different 36 cards than the ones actually in play.
- **Read the live piles.** The engine's `Deck(player)` zone is the current draw
  pile, not the roster; reading it would both omit cards already drawn/played and
  expose the draw order, which the existing zone-roster popover deliberately hides
  by sorting (view_board.go `readableZoneIDs`).

The Deck list must therefore show the _static roster as generated_, and it must
be able to show it for **both** players — today every game is the open format,
but sealed/hidden-list formats are anticipated.

## Decision

**Retain the generated roster at match setup and thread a read-only projection to
the client; never regenerate and never read the piles.**

`SetupDecksFor` keeps the `deckgen.Deck` it already generates and, alongside the
Houses/Maverick/Legacy slices it returns today, exposes a per-player **roster
projection**: for each of the 3 Houses, its 12 rows of `{card definition, Rarity,
Maverick, Legacy}` in Slot order. The projection is a plain value copied into the
web game state next to the existing `deckHouses`/`mavericks`/`legacy` fields — no
pointers into deckgen, no live engine handles.

The web client renders the Deck list from that projection alone. Visibility is
gated by a single predicate, `deckListVisible(viewer, owner) bool`, that returns
`true` today; a future sealed/hidden-list format flips it without touching the
render path. The popover is public by default because the current format is open.

Presentation choices that fall out of this decision:

- **Sort within a House column by type then name** (Creature, Artifact, Upgrade,
  Tactic), not by draw order — the projection is static, so there is no order to
  leak and no collector number to sort on (the engine card carries none; numbers
  live only in provenance bookkeeping).
- **Rarity renders as one compact shape per tier**, distinct by silhouette rather
  than colour (triangle/square/pentagon/hexagon for Common/Uncommon/Rare/Special,
  the existing link glyph for Connected). Colour alone is rejected for
  accessibility; the card face keeps its existing 1–4 diamond pips.

## Consequences

- The Deck list is always truthful: it shows the exact 36 cards dealt, immune to
  later pool changes, because the roster is captured at the one moment it is known
  to match play.
- Nothing about the live piles reaches the list, so it cannot leak draw order or
  desync with what a player has already drawn — the hidden-info boundary the
  zone-roster popover guards is preserved by construction.
- `SetupDecksFor` gains one more return (or one field on a returned struct); its
  callers that ignore the roster (`sim`) are unaffected, and the roster rides
  next to data the client already holds.
- Visibility is one predicate, so opening the door to sealed formats later is a
  one-line change, not a render-path rewrite — and until then both players' lists
  are open, matching the only format that exists.
- The rarity shapes are a Deck-list-scoped asset set; if the single-glyph scheme
  is liked, unifying the card face and gallery onto it later is additive, not a
  migration.

# 31. Zone movement is one mechanism; the KeyForge verbs are sugar over it

## Context

Five effect families move cards between zones, and each has grown a swarm of
bespoke types: Archive (eight), Discard (nine), Purge (nine), Shuffle (four), Put
(eight). They differ in parameters, not mechanic. Every one of them:

1. names a **source** zone — a battleline, a hand, a discard pile, the top of a
   deck, archives;
2. **selects** cards from it — a chosen one, any number chosen, every matching
   one, a random one, the top N — narrowed by house, type, name, or trait;
3. moves them to a **destination** — hand, top/bottom of deck, deck-shuffled,
   archives, discard, or out of the game;
4. sometimes binds the moved card in context (`ctx.It`) for a follow-up, or
   tallies how many moved (`ctx.Produced`) for a "this way" count.

Naming each combination a new type hides that it is one mechanic. `ArchiveTopOfDeck`
and `ArchiveTopOfDiscard` differ only in the source zone; `DiscardRandomFromHand`
and `DiscardRandomFromArchives` only in the source zone; `PurgeFromHand`,
`PurgeCreatureFromHand`, and `PurgeEachFromHand` only in the selection strategy.
The families sprawl because the axes above are spelled into names instead of
fields.

This is the same shape the engine already factored for Æmber quantities — the
`{Amount, By, Per}` cluster shared by `StealAember`/`LoseAember`/`CaptureAember`
through helpers, not a copy per effect (see "Reused effect shapes get one shared
helper" in `internal/engine/AGENTS.md`). Zone movement is the next such shape, and
the largest.

## Decision

**All movement of a card between zones is one engine mechanism**, parameterized by
the four axes above: source zone, selection strategy, destination, and the optional
`ctx.It`/`ctx.Produced` side effects. Where the source is play the selection is the
existing `Target`/`Selection` vocabulary; from a pile it is a small zone-selection
strategy (chosen / any-number / all / random / top-N, plus the house/type/name/trait
filters). The destination is the existing `Destination` (`To.Hand`, `To.TopOfDeck`,
`To.Archives`, …) extended with the two terminal destinations, discard and
out-of-the-game (purge).

**The mechanism is not public.** A card author never reaches for a generic
`Move`/`Put`; they reach for the KeyForge verb — `Archive`, `Discard`, `Purge`,
`Shuffle…IntoDeck`, `Put…` — each a thin authoring struct with flat, ergonomically
named fields that builds the shared mechanism with its destination fixed.
`PurgeCard` sets the destination to out-of-the-game; `ArchiveTop{From: Deck}` sets
it to archives; `DiscardCard{Player, Zones, Selection, Amount, AnyNumber}` sets it
to discard. The printed card says the verb, so the verb is what the author writes.

**The sugar delegates; it does not embed.** This follows the quantity decision:
Go composite literals do not promote embedded fields, and the `card` facade aliases
the engine structs (`card.PurgeCard = engine.PurgeCard`), so embedding a shared
`movement` struct would turn every flat authoring site into a nested one across the
whole card database, buying nothing. Each verb keeps its own flat fields; only the
behavior is shared.

### Why not one public `Put`

KeyForge has no single "move a card" verb. It has archive, discard, purge, shuffle,
and put into play / into hand / on top of deck. The printed card text names the
verb, and the printed text is the authoring vocabulary (ADR 0006: effects render
themselves, and the card is authored in the words it prints). A public generic
`Move`/`Put` would read nothing like the card and would push the author to think in
engine terms. Keep the verbs; share the engine underneath them.

## Consequences

- A new zone-movement card reaches for an existing verb and adds a selection filter
  or a destination — not a new effect type. The families stop sprawling.
- The movement is exercised once, at the mechanism; each verb's sugar needs only
  its own text-and-wiring test.
- The DiscardTop-versus-DiscardRandom question dissolves. "The top card" and "a
  random card" are two **selection strategies** of one Discard, so
  `DiscardCard{..., Selection}` carries either (`Chosen` or `Random`) — no
  `DiscardTop` versus `DiscardRandom` split, and calling a random discard "top of
  deck" never comes up. `AnyNumber` covers the "discard as many as you like"
  variant on the same struct.
- Realized incrementally, not in one sweep. Today `PurgeCard` (zone + count +
  up-to), `PutChosen` / `PutFromPlay` (target + `Destination`), the `To`
  destinations, and the `ctx.Produced` tally are partial realizations. Each family
  is folded into the mechanism as its cards are touched; this ADR is the target the
  folds converge on, so a new movement effect is written to fit it rather than
  adding one more bespoke type to unwind later.

# Gigantic tutors are a cross-set reservoir cluster pulled per gigantic

## Context

Mass Mutation ships four **gigantic tutors** — It's Coming…, Build Your Champion,
Digging Up the Monster, and Tomes Gigantica — Tactics that search a player's deck
and discard pile for gigantic creatures and move them to a zone. Unlike ordinary
cards, a tutor is never part of a house's draftable pool. In KeyForge a gigantic
tutor is added to a deck only _because_ that deck contains a gigantic creature,
and it takes the gigantic's house, not a house of its own — a houseless card
stamped to the pod it joins.

This is the same shape as the shard cluster (ADR 0017/0021's reservoir cards, ADR
0036's `ClusterPool`): cards undraftable on their own that enter a deck through a
catalog-wide cross-set cluster rather than a set's normal pool. The existing
cluster strategies did not cover the tutors' trigger, though. `OnePerHouse` places
one member per **deck house**; nothing placed a member per **gigantic**, into that
gigantic's house. The tutors need exactly that: one tutor per gigantic base in the
deck, stamped to the gigantic's house.

The reprint machinery (ADR 0021) already dissolves the duplication problem: each
tutor is one Go definition carrying every provenance ref it reprints (It's Coming…
is Mass Mutation #117 and Menagerie #305), so the cluster has exactly one member
per tutor name and dedupe is automatic.

## Decision

The gigantic tutors are members of one cross-set reservoir cluster, `Tutors`, with
a new cluster strategy `PerGigantic`.

### `PerGigantic` is a cluster strategy, not a bespoke path

`PerGigantic` joins `OnePerHouse` in the `ClusterStrategy` enum. Like the shard
cluster, its members are houseless reservoir cards (`Profile.Houseless`) that are
undraftable on their own; they enter a deck only through the pull. A set gains the
tutors by attaching the catalog `ClusterPool` with `WithClusters`, exactly as it
gains any other cross-set cluster — no new pool type.

`PerGigantic` fires off **deck structure** (a gigantic being present), not off a
member being drawn, so unlike a drawn cluster it needs no trigger, lead, or
rollable member — `validateClusters` skips those checks for it. `buildClusters`
still guarantees at least one member, so the pull always has something to choose.

### The pull happens per gigantic base, into the gigantic's pod and house

Deck generation places the tutor after the gigantic art half, in the pod chain:
`placeGiganticTutor(placeGiganticArt(…))`. For each pod it counts the gigantic
base slots and gathers the free ordinary slots in one pass up front — a placed
tutor is itself an ordinary (`GiganticNone`) card, so gathering the free slots
before placing anything keeps a second tutor from overwriting the first. It then
places one tutor per base, each chosen at random from the cluster's members, into
a free slot with `Special: true` so `materialize` rehouses the houseless tutor to
the pod's house (the same `Rehouse` path a maverick uses). A pod with no free slot
pulls nothing; a set with no `PerGigantic` cluster pulls nothing.

`expandPodClusters` skips `PerGigantic` clusters the same way it skips
`OnePerHouse`: their placement is driven by structure, not by the per-pod cluster
expansion.

## Consequences

- A new reservoir mechanic that keys off deck structure rather than a drawn member
  is a new `ClusterStrategy` value plus its placement pass, not a bespoke branch in
  the pod builder — the same extension shape `OnePerHouse` established.
- The tutors search with the `Gigantic` card filter, which admits either half of a
  gigantic creature. It's Coming… (search for "either half", one card, into hand,
  then shuffle) reuses the existing choose-one `Search` node directly. The three
  MoMu tutors search for "two halves of a gigantic creature" — up to two gigantic
  halves, which the rulebook rules need not belong to the same creature — and every
  deck search shuffles afterward (rulebook: an entire-deck search is always
  shuffled). That count-capped half-search is the `Search` node's `Max int` field:
  the zero value takes exactly one card, `Any` takes every match, and `Max` caps
  the take at that many optional choices (the MoMu tutors set `Max: 2`). Digging Up
  the Monster's mid-effect shuffle (shuffle after finding the halves, before placing
  them on top of the deck) is expressed by leading the search with `Shuffle{}` and
  targeting `To.TopOfDeck`, so the found halves land atop an already-shuffled deck.
  All four tutors are now implemented.
- Because a tutor is one Go definition carrying all its provenance refs, the
  cluster has one member per tutor name regardless of how many sets reprint it.

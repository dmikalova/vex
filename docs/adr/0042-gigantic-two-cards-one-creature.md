# Gigantic: two cards, one creature, through enter/leave funnels

## Context

A gigantic creature is printed across two physical cards — a **base half**
carrying power, armor, traits, keywords, abilities, and the rarity icon, and an
**art half** carrying only the bonus icons. Both share the same name, House, and
card type (Creature). The Master Rulebook's rules are: you play a gigantic only
by having **both halves in hand** and playing them together as one creature; it
counts as **two cards while out of play but one creature while in play**; playing
it is **playing one card** (legal on turn one); and after it leaves play the two
halves are separate cards again. Only the art half has bonus icons; only the base
half has the text box.

This collides with several invariants the engine is built on:

- `GameState` is flat, pointerless, and comparable (ADR 0005): a two-card
  relationship cannot be a pointer, slice, or map — it must be an optional
  `LocalID+1` byte, like the Upgrade (ADR 0001) and Under (ADR 0016) chains.
- The battleline derives flank, neighbor, and center from **adjacency of IDs**,
  and every board count assumes one ID per creature — so a gigantic must occupy
  exactly **one** battleline slot.
- Both halves must survive as distinct cards with distinct `LocalID`s, because
  out of play they sit in separate zones and are individually targetable (Infurnace
  purges an art half from a discard pile for its Æmber pips; the base half has
  none).
- Every play, put-into-play, and leave-play path is a potential leak point where
  one half could enter or exit alone.

Two further problems are specific to gigantics. **Playing** must unify normal
hand play with effect-plays (Exhume from discard, Wild Wormhole from deck-top,
Æmberlution from hand, Saurian Egg from a just-discarded pile) and with
put-into-play (Overlord Greking from discard). **Deck generation** must guarantee
both halves plus exactly one of four tutor cards, all in the gigantic's House.

## Decision

### One creature, two IDs — the base half is the representative

Both halves register as ordinary cards with their own `LocalID`s. While in play,
the **base half's ID sits in the battleline and _is_ the creature**: every read
(power, armor, damage, counters, traits, keywords, abilities, exhaust, upgrades,
House, type) comes from it, exactly as for any battleline creature. The **art
half is in play but slot-less**, linked to the base by a new optional
`LocalID+1` field pair on `CardCore` (the third intrusive relationship after
Upgrades and Under), and contributes exactly one thing to the combined creature:
its **bonus icons**, read through the link. So `Target.Creature.WithoutBonusIcons()`
(Wail of the Damned) sees the art half's pips and spares a gigantic that has any.
The art half holds no battleline slot, is not directly targetable by "choose a
creature", and leaves play only with its base. `InvariantError` gains a
conservation/backlink check for the link, alongside the Upgrade and Under checks.

Rejected alternatives: **overloading the Under chain** (its whole contract is
"out of play", so every Under reader, `Peekable`, and `discardUnder` would need a
guard) and **putting both IDs in the battleline** (breaks the one-slot invariant
— flank/neighbor/center are adjacency-derived and every board count would
double).

### Pairing is by name and role; forming is a one-to-one matching

A `Gigantic` keyword plus a `GiganticRole` (Base or Art) on the immutable
`CardDefinition` establish the pairing; partners match on **same name, opposite
role**, so any base of "Deusillus" pairs with any art of "Deusillus" — which is
why two copies in a deck give "either base + either art", never "two bases".
Forming a gigantic is a **one-to-one matching that consumes each physical half at
most once**, so no single resolution ever forms more gigantics than the matching
yields.

### One file authors the pair; registration splits it

A gigantic is authored as **one** `card.New(...)` marked `card.Gigantic()` with a
**single provenance** — the base half's collector number. The source catalog
(`internal/cards/provenance`) lists only the base, typed `gigantic creature base`
(e.g. Deusillus is MM 244 with no MM entry for its art half), so the gigantic is
**one source card** for coverage; the art half is a synthetic engine artifact and
carries **no provenance of its own**. Registration splits the authored card into
two `CardDefinition`s by a fixed, mechanical rule — **bonus icons → art; power,
armor, traits, keywords, abilities, rarity → base; name, House, type → both** —
emitting the base (carries text and rarity, leads its clusters) and the art
(`Rarity.Connected`, bonus icons, no provenance). This reuses the
Template/materialize seam (ADR 0004); deck generation then sees the ordinary
base-leads-art cluster it already handles. "One file" is an authoring convenience
only: it still registers as two cards with two `LocalID`s. (The MM gigantics print
no bonus icons, so their art halves carry none — the split rule still routes any
future printed icon to the art half.)

### Enter and leave play are two central funnels

Only two seams know that two cards make one; everything else is oblivious,
operating on a single half out of play or on the one creature in play.

- **Enter funnel** — shared by `playCardFromZone` (every true play) and
  `putIntoPlay` (every put-into-play). When the entering card is a gigantic half,
  it recruits its partner by the one-to-one matching from a **partner-source set**
  the initiating effect supplies — default `{hand}`; **Overlord Greking** passes
  the specific destroyed pair (from discard, no hand reach, so killing one gigantic
  forms exactly one); **Saurian Egg** passes its just-discarded set then hand. Both
  halves enter together **before anything else**, then the standard resolution runs
  once: constant abilities → bonus icons (from the art half) → the single after-play
  window (`Play:` and "after you play a creature" together). A **lone half cannot be
  played** — "cannot" overrides "must" — so a failed play (no partner available)
  returns the half to its source zone, including a deck-top play that puts the half
  back on top of the deck.

- **Leave funnel** — each mover that relocates a card out of play loops over the
  gigantic's halves (`giganticHalves`) and reuses the existing single-card
  teardown (`leavePlayDestroyed`). After the ward gate (a ward saves the whole
  creature → neither half moves), both halves go to the **same** destination as
  **two separate cards**. `discardDestroyed`,
  `purgeFromPlay`, `putIntoHand`, `putOnTopOfDeck`, `putIntoArchives`,
  `putIntoDeckShuffled`, `swapAcrossZones`, and `removeFromAnyZone` all route
  through it, so purge, discard, bounce, archive, deck, and shuffle each split the
  pair. (`GraftUnder` is not wired: a gigantic has no Graft keyword, so a grafted
  card is never a gigantic half.)

### Deck generation: base leads two clusters; tutors are cluster-only houseless

The base half leads two `ByLead` clusters: a `PullExact` cluster guaranteeing its
art half, and a `RandomCount{Min:1, Max:1}` cluster over the four tutor cards
(It's Coming, Build Your Champion, Digging Up the Monster, Tomes Gigantica). The
tutors are **`Houseless` + `Connected`**: houseless so a pulled tutor is stamped
with the gigantic's House, and Connected so a tutor can **only** enter a deck by
being pulled — it never rolls in the normal pool or the special slot. This needs
deckgen's houseless routing to respect `Connected` (skip the special-slot roll).
Mass Mutation and More Mutation are one combined set, so all four tutors are
present and the random-of-four is live.

## Consequences

- **Every play/put/leave chokepoint must route through a funnel or a half
  leaks.** The exhaustive list to wire: enter — `playCardFromZone`, `putIntoPlay`;
  leave — `removeFromPlay`/`leavePlayDestroyed`, `discardDestroyed`,
  `purgeFromPlay`, `putIntoHand`, `putOnTopOfDeck`, `putIntoArchives`,
  `putIntoDeckShuffled`, `swapAcrossZones`, `GraftUnder`, `removeFromAnyZone`. A
  new mover added later must go through `leavePlayTo`.
- **`Houseless` + `Connected` is a new deckgen combination.** Today a houseless
  card always rolls in the special slot; the tutors require "houseless but
  cluster-only", so the routing in `NewSet` must check `Connected` before bucketing
  a houseless card into the special pool.
- **The split is registration-time, not runtime.** `mage gen` renders one
  combined doc comment per gigantic; the two `CardDefinition`s it produces are what
  the catalog and coverage see.
- **Anti-duplication rides on scoped partner sources.** The correctness of "Greking
  plays one, Saurian Egg with doubles plays two" depends entirely on each effect
  passing the right partner-source set to the enter funnel; a blanket same-name hand
  search would over-form.
- **The art half is a third conserved relationship.** After Upgrades and Under,
  this confirms ADR 0001's intrusive-list idiom generalizes again; the invariant
  walk must account the link so a slot-less in-play art half is never seen as lost.

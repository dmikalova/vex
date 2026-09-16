# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Work that has **no consuming card in an implemented set yet** does not belong
  here — park it in [todo-future-set.md](todo-future-set.md), keyed to the set that
  first needs it.
- Cite the ADR or doc that decided an item where one exists.
- **When you need a decision from the human, ask in the reply itself, grill-me
  style** — a numbered list of `❓ **Q1** - **title**: <question>` with a `➡️`
  recommended answer under each — at the end of the turn, then stop. Do **not**
  reach for an interactive question tool: under autopilot it is auto-answered with
  "work autonomously" and the human never sees it. Plain-text questions at the end
  of the turn are the channel the human actually reads.

## Mass Mutation — deferred mechanics (grill before building)

The `massmutation` set is stubbed and its easy cards are being implemented. Two
mechanics are **deferred for a grilling session** because each needs a design
decision (a new node, a deckgen-time hook, or a cluster shape) before any of its
cards can be authored. Each group's cards are `//go:build todo` stubs in
`internal/cards/sets/massmutation/`. **Grill each group, decide the shape, then
implement the whole cluster together** (implement-cards: shape for the cluster,
not the first card).

### 1. Gigantics (two cards, one big creature) + their tutors

A **gigantic** creature is a single large creature split across **two cards** (a
**base half** carrying the text/stats and an **art half** carrying only bonus
icons); both must be in play together to form the creature (ADR 0042).

**Engine + facade structure — DONE (all green, engine 100% covered):**

- Engine mechanic: `GiganticRole {None,Base,Art}` on `CardDefinition`;
  `CardCore.GiganticPartnerPlus` symmetric link; enter funnel (`playGigantic`
  recruits the opposite-role same-name half from hand, records one play of the
  base) and leave funnel (both halves travel together); bonus icons resolve once
  through the link; invariants account the slot-less art half.
- Facade: `card.Gigantic(name, house, rarity, opts...)` builds the pair and
  registers **only the base** (the two halves share a name, which the database's
  unique-name rule forbids for two registered cards). The synthetic art half —
  no provenance — rides on the base's `GenerationProfile.GiganticArt`.
  `card.GiganticArt(base)` returns it for tests/tools.
- Rulebook: a standalone "Gigantic creature" term under `SectionCardType`
  (no `Gigantic` keyword — the engine pairs by `GiganticRole`, not a keyword).

**Still to do:**

- **Deckgen placement**: the deckgen cluster machinery is entirely name-keyed
  (`countMember`, `freeClusterSlot`, PullExact's `Name != lead`), so two
  same-named halves **cannot** use clusters. Add a post-fill pod pass that, for
  each base-half slot, places its carried `Profile.GiganticArt` into a second
  slot of the same pod (so a deck that draws a gigantic draws both halves).
- **Power/armor data gap**: the provenance JSON lists only the base half, typed
  "gigantic creature base", and carries **no power** for it (the stub generator
  therefore emits no `Power:` line). Confirm authoritative power/armor for each
  gigantic before authoring — do not guess stats.
- **Cards** (un-stub base + author abilities; art half is synthetic via
  `card.Gigantic`): **MM (3):** Deusillus (MM244 Saurian, Play: capture all
  opponent Æmber + deal 5 to an enemy creature; Fight/Reap: move 1Æ from itself
  to the common supply + deal 2 to each enemy creature), Ultra Gravitron
  (MM125 Logos), Niffle Kong (MM422 Untamed). **MoMu (11):** Tormax, Wretched
  Anathema, Horizon Saber, Ascendant Hester, Sirs Colossus, Bawretchadontius,
  Boosted B4-RRY, Dodger's 10, Cadet Allison, J43G3R V, Titanic Bumblebird.
- **Tutors (4)** that fetch gigantic halves — build _after_ gigantics exist:
  It's Coming… (MM #117), Build Your Champion, Digging Up the Monster, Tomes
  Gigantica (MoMu #002-004). These are Houseless + Connected (cluster-only,
  pod-stamped to the gigantic's house); deckgen must route them.
- **Partner-source corner cases**: Saurian Egg (MM210) and Overlord Greking
  (Call of the Archons) put a gigantic pair into play from a scoped source
  (`putIntoPlay` is not yet gigantic-aware — wire at the card layer).

### 2. Mass Mutation cards still to implement (normal cards, grill one by one)

The small-primitive bottleneck is cleared: every remaining non-gigantic Mass
Mutation card is a **normal card** that needs no new effect node, so they are
tackled one at a time or in small themed groups, grilling each for its wording and
the smallest engine seam it needs. Run `mage tool:missing -set=massmutation` for
the live list; `mage tool:nextCard -set=massmutation` names the next stub.

### 3. Remaining cards

006 Drecker
011 Mark of Dis
012 Mindfire
021 Essence Scale
028 Picaroon
029 Relentless Creeper
035 Etan's Jar
042 Painmail
049 The Pale Star
051 Turnkey
100 Animator
102 Cyber-Clone
113 The Archivist
126 Ardent Hero
127 Bull-wark
131 Commandeer
143 Angry Mob
153 Purify
159 Book of Malefaction
165 Lady Loreena
174 Purifier of Souls
200 Blast from the Past
206 Gladiodontus
210 Saurian Egg
222 High Priest Torvus
224 Legion's March
257 Shoulder Id
263 "Borrow"
267 Lucky Dice
287 Mole
292 Shadowsaurus
303 Blast Shielding
335 Ambassador Liu
337 Commander Dhrxgar
347 J.O.N. Cargo
383 Growth Surge
394 Aemberlution
395 Blossom Drake
396 Chonkers
403 Mercy, Malkin Queen

## Snapshot: an ability reads its card as it was immediately before leaving play

Decided in the bonus-icon grill but **general to every ability**: when a card
leaves play (destroyed, purged, returned, sacrificed) _while_ an ability of its
(or an ability referencing it) is resolving, the ability must consider the card as
it was immediately before it left — a snapshot of its stats / traits / house /
position / bonus icons / Æmber — not read stale or absent state. This refines
[ADR 0030](adr/0030-a-card-out-of-play-takes-no-further-part.md): the card takes no
_new_ action, but an in-flight ability that references it uses its pre-departure
snapshot.

**This is a full verification pass, not one card.** Audit every effect/ability
that reads a card's live properties during resolution (power / armor / traits /
house / position / bonus icons / Æmber-on-card) and confirm that a card which has
just left play mid-resolution is read from a snapshot rather than from cleared
state. Bonus-icon resolution already benefits — icons live on the immutable
`CardDefinition`, so a purged/destroyed card still exposes them. The pass must find
the cases that read _mutable_ state (counters, damage, position, granted stats or
keywords) and either capture a snapshot at the leave-play boundary or resolve the
reference against a captured value. **Decisions needed:** where the snapshot is
captured (a per-resolution frame vs. a leave-play hook), how it stays flat and
comparable (ADR 0005), and which reads actually need it (many effects finish before
their source can leave). Grill the shape, then sweep.

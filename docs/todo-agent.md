# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Cite the ADR or doc that decided an item where one exists.
- **When you need a decision from the human, ask in the reply itself, grill-me
  style** — a numbered list of `❓ **Q1** - **title**: <question>` with a `➡️`
  recommended answer under each — at the end of the turn, then stop. Do **not**
  reach for an interactive question tool: under autopilot it is auto-answered with
  "work autonomously" and the human never sees it. Plain-text questions at the end
  of the turn are the channel the human actually reads.

## Maverick houses complete every OnePerHouse cycle to all 9 Houses (ADR 0036)

Deck generation will grow a **maverick House pod** (deck-generation.md §1): very
rarely a whole pod is a House not in the deck's Set at all, drawn wholesale from
another Set. That breaks the current per-set `OnePerHouse` gate: a maverick Saurian
pod dropped into an Age of Ascension deck fires the Shard cycle for a House AoA has
no Shard for. So **when maverick houses land**:

- Expand the `OnePerHouse` completeness gate (`Set.validateClusters`) from the
  set's own Houses to **every House a maverick pod can introduce** — effectively all
  nine. AoA (and every Shard-bearing set) then needs a Shard for **all** Houses,
  including Saurian and Star Alliance, which AoA never printed.
- Add those missing Shards as `//go:build todo` stubs so the gate stays honest
  until they are implemented (Saurian and Star Alliance were introduced after AoA,
  so their Shards are Vactrol-generative, like Master of 4/5).
- The Ambassador and Plant templates already cover this for free — they materialize
  per partner House from `DeckHouses`, so a maverick House gets its Ambassador/Plant
  with no new card. Only the Shard cycle, which is concrete per-House cards, needs
  the stubs.

Until then, a legacy-drawn Shard does not fire the cycle either: `expandClusters`
reads only the drawing set's own `clusters`. Revisit both together when maverick/
legacy whole-pod overlays are built.

## Engine distillation sweep — 20 findings (one-off nodes → cohesive systems)

The theme: three sets now cover the game's surface, so card-named one-off nodes
that fuse several effects for a single card should collapse into small
parameterized systems, and the Resolver/log families they spawned should collapse
with them. Grouped by mechanic so the shared primitive lands once.

### House-permission system pass — implement ADR 0037 (`MayPlayOrUse`)

The out-of-house permission grants are spread across five effect nodes —
`GrantFight` (house-filtered), `GrantFightAnyHouse`, `MayActFriendlyHouse` + its
`HouseGrant` bitset, `MayPlayOffHouse` (exclusion/count/type), and
`MayUseFriendlyArtifacts` — each with its own Resolver method (`GrantFightForHouse`,
`GrantFightAnyHouse`, `GrantUseForHouse`, `GrantPlayForHouse`,
`GrantUseArtifactsAnyHouse`, `GrantOffHousePermit`) and log record
(`FightGrantedForHouse` …). H1/H5 turned out only **partially** generalized, so
folding the Resolver methods and logs (H2/H3/H6) on their own would be churn.
ADR 0037 decides the target shape: one `MayPlayOrUse{Houses, Grant, Types, Count}`
node — `Houses` = named / chosen / any / all-but-a-named / houses-you-control;
`Grant` = a GrantPlay|GrantFight|GrantUse bitset; `Types` = a card-type filter whose
**zero value means all types** (opt in to creatures/artifacts); `Count` = 0 for
unlimited. One `GrantMayPlayOrUse` resolver seam, one `MayPlayOrUseGranted` log.
Implement it: fold **all five** nodes (the exclusion/count of `MayPlayOffHouse` and
the artifact subject of `MayUseFriendlyArtifacts` are now the `Houses`/`Count`/`Types`
axes, not siblings), add `GrantFight` to the bitset, collapse the resolver/log/state
families. Keep engine coverage at 100%.

### Spend-N-cards scaling — producer→consumer count (folds old H8, H16, H17)

De-fuse "spend N cards, then a benefit scaled by N" into one shared mechanic: a
**producer** sub-effect (discard / purge / archive / destroy / sacrifice / return /
lose Æmber …) records how many things it acted on, and the **consumer** that
immediately follows reads that count to scale its amount. Model it as a
producer→consumer link (a `Then`/`Sequence`-threaded `ctx.Produced` tally), **scoped
per-effect, not per-ability**: every KeyForge "… this way" clause binds to exactly
one preceding same-verb spend, so per-effect scoping keeps a card with two spends in
one ability unambiguous (Codex of True Names destroys itself _and_ returns Fiends;
"for each card returned this way" must count only the Fiends). It is **generic** —
the count scales any consumer effect, not just damage.

Card examples confirming the shape (all "this way" = only what this effect moved):

- Destructive Analysis / Martian Civil War — purge N from archives, deal 2 damage
  per (`PurgeArchivesForDamage`).
- Harvest Time — choose a trait, purge each, gain Æmber per
  (`PurgeEachOfChosenTrait`).
- "Discard your hand, deal 1 damage to each creature for each card discarded this
  way"; "discard your hand, draw a card per"; "discard any number, steal 1 Æmber
  per"; Obsidian Forge — sacrifice any number, reduce a key's cost per.

Both former H16/H17 stragglers are this pattern, not `Selection`-over-a-pile nodes,
so they resolve once the tally exists. Needs careful 100%-coverage work; do with a
review, not unattended.

### Generalize `PlayFromOpponent{Source}` (was H9)

One node for "play a card from your opponent," parameterized over the source zone
(random-from-archives vs top-of-deck) and touching the web icon map once. Worth
building ahead of the cards so they drop in effortlessly. Cards: Murkens and Fidget
(today); Gracchan Reform, Talent Scout, Blank Check, and likely more (later). Design
the `Source` axis to cover archives-random and deck-top up front.

### `EachPlayer{Do}` combinator (was H14)

A controller-rescoping wrapper so "each player does X" authors as one effect. A bare
`EachPlayer{Sequence{...}}` is the wrong shape for cards that run a phase for **all**
players before the next (discard all hands, then refill all), so the wrapper must
rescope "you" / frames / trigger ordering per player, not just loop a `Sequence`.
~10 cards need it; build the wrapper now so they drop in later. Design deliberately
— the subtlety is pronouns and trigger windows.

## Engine — effects & mechanics

- **Generalize top-of-deck look (`LookAtTop{Amount, Then}`).** Replace
  Philophosaurus's hardcoded `LookAtTopSort` with a general `LookAtTop{Amount,
Then}` that composes existing zone-routing sub-effects (hand / archive / discard
  / reorder), folding Eyegor, Lay of the Land, and future Vandalize. This is a
  **zone slice** (top N of deck). Separate follow-up primitive: a **result-set**
  reference — "the exact cards a preceding sub-effect just moved" (may be empty even
  when the discard is non-empty) — for Junk Restoration / Haunting Measures and
  scrap-adjacent cards. Different primitives; do not block one on the other. First
  implement the result set against cards that act on a **result set of one**. The
  "purge/archive the top card of the discard" cards are the `Amount: 1` case of a
  discard-top variant under this same family.
- **Cooperative Hunting → an iterator over `DealDamage`.** "Deal 1 damage X
  times, choosing any creature each time" is distinct from Sack of Coins ("deal X
  to one creature"): Cooperative Hunting reuses the `DealDamage` primitive wrapped
  in a repeat/iterator that re-prompts for a target each iteration. All selections
  are made first, then all the damage is dealt **simultaneously**. Model as a
  `ForEach`-style iterator wrapping `DealDamage{Amount: 1}` with a per-iteration
  target choice and deferred simultaneous application.
- **Unify the deck "dig until" family — Old Boomy folds in with Lethologica
  (deferred).** There are two dig-until-a-match nodes that share one shape (walk the
  top of a deck one card at a time, do something to each, stop on a match or when the
  deck runs out, then optionally act on the found card): `DiscardDeckUntil` +
  `PutDiscardedIntoHand` (Sound the Horns, Invasion Portal, and unimplemented
  **Lethologica** — "Discard from the top of your deck until you discard a Logos card
  … put it into your hand") and `RevealDeckUntilHouse` (**Old Boomy** — reveal and
  archive each until a Brobnar card or you choose to stop). Fold them into one
  parameterized dig — axes: **per-card fate** (discard vs archive), **terminator**
  (type/house filter), **stop choice** (Old Boomy lets the player stop; the discard
  digs do not), **visibility** (reveal-to-both vs private discard), and **found-card
  fate** (leave in `ctx.It` for a follow-up vs nothing). Deliberately **kept out of
  the `chooseFromTop` read-and-route consolidation** (that core is fixed-N; this is a
  variable-count gated loop returning a bool — cramming it in bloats both). Do it
  when Lethologica is implemented, alongside Old Boomy, with a review — not
  unattended.

## Card wording / authoring

- **Typed card-name references — convert remaining literals + add a guard.**
  Decision taken (option a): cards reference another card by `OtherCard.Name`, or by
  a `const XName` the lead card declares when a plain `.Name` would form a package
  var-init cycle (the Hyde/Velum and Igon Green/Terrible pattern — the follower
  reads the lead's const). **Done:** Help from Future Self → `Timetraveller.Name`;
  Igon the Terrible → `IgonTheGreenName` const; the two self-references (Disruption
  Field, Uncharted Lands) now use `card.Target.GrantingCard`, which resolves to the
  granting card and renders the `{card}` placeholder (its own name) — no literal.
  **Remaining:** add a source-scanning test that fails on a string literal in any
  card-name field (`SearchForName.Name`, `ControlsNamed.Name`, `NamedCardPurged.Name`,
  `ItIsNamed.Name`, `ReturnNamedToHand.Name`, `PutFromDiscard.Name`,
  `AttachSelfTo.Host`, `.Named(...)`), allow-listing `card.New`'s own-name literal,
  the `const XName` declarations, and `Cluster.Name`.
- **Generic zone-move `Match` predicate — do now.** Replace the `Type` / `Trait` /
  `OrTrait` triple on `PutFromDiscard` (and its zone-movement siblings) with one
  composable `Match` predicate expressing a union of type/trait clauses (e.g.
  `Upgrade` OR trait `Robot`), so there is one place to test "any match in the
  discard". The Chief Engineer Walls empty-`May` bug is already fixed via
  `vacuous()`; this is the consolidation refactor on top, keeping the shared
  `matches` / `vacuous` helper.

- **Orator Hissaro** could read: "Play: Exalt and ready each neighboring creature.
  For the remainder of the turn, those creatures belong to house Saurian."
  Blocked on two engine additions: making `Exalt` a `combinable` (so `Ready`+`Exalt`
  fold to "ready and exalt each neighboring creature" — but the fold must keep
  "exalt N times" for `Amount > 1`), and a pronoun form of `BelongToHouse` so the
  second sentence reads "those creatures belong to …" instead of repeating the
  target. Today it renders correctly but verbosely (target repeated three times).
- **Borr-Nit** and similar could be atomized and recomposed further (decompose
  fused effects into shared nodes).

## Card catalog / provenance

_No outstanding items._

## Web — mobile, previews, layout

- **Card iconography glyph sweep (blank/missing asset gaps).** Sweep every card
  for missing/blank glyph assets (use the totality test that names uncovered
  effects, and eyeball the strips) and fill the gaps. Goal: no card renders a blank
  or broken glyph in its icon strip. (The Collar of Subordination zero-Target
  `TakeControl` glyph is fixed — it now renders the "this creature" host noun.)
- **A card is occasionally not dimmed on a new turn (stale selection carryover).**
  Sometimes a card that should be dimmed/unplayable at the start of a turn stays
  lit. Suspected cause: a per-slot selection/highlight state from the previous turn
  is not cleared when the turn advances, so the slot keeps its old dimmed/lit
  computation. Find where the client resets per-slot selection on turn change and
  ensure the dim state is recomputed (not carried over) when a new turn starts.
  Needs a concrete repro first.

## Logging

- **Effect-use logs should narrate the outcome, not just "uses X's action
  ability" (ADR 0011).** `ActionAbilityUsed` prints the generic "P0 uses X's action
  ability" with no detail of what happened. `ShuffleIntoDeck` now records a
  `ShuffledIntoDeck` outcome (naming the public discard cards, counting the hidden
  hand/archives cards). Remaining: sweep the other effects that still lean on the
  bare `ActionAbilityUsed` line for their narration.

## Tooling / tests / docs

_No outstanding items._

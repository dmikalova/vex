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

## Clusters + deck-house templates (ADR 0036, ADR 0004)

Design settled with the human (grill). Two orthogonal axes: **clusters** decide
which other cards are placed; **templates/materialize** decide what one drawn card
becomes. Build one piece at a time, in this order.

**Done:** `DeckHouses [3]House` on `SlotContext`; the cluster framework
(`deckgen/cluster.go` types + `card.Cluster`/`card.InCluster`/`card.LeadsCluster`
/`card.ClusterStrategy`/`card.ClusterTrigger` facade); the deck-wide `OnePerHouse`
pass with maverick substitution (`generator.expandClusters`); the
complete-by-construction gate (`Set.validateClusters` panics if a set can deck a
House with no member); the seven Age of Ascension Shards wired onto the cluster
(`shard_cluster.go`, each Shard `card.InCluster(shardCluster)` + `OneCopyPerDeck`);
the pod-local cluster pass (`generator.expandPodClusters`) placing `WholePool`,
`RandomCount`, `SelfPull`, `PullExact`, and `Pull`; the `SelfPull(Min, Mean)`
strategy (a Min-plus-Poisson count capped at `PodSize`); the `PullExact` strategy
(one partner per lead instance, always `ByLead`); the `Pull` strategy (a per-partner
Min-plus-Poisson count, each pulled partner setting its own rate with `card.Pulled`,
always `ByLead`); the four Horsemen migrated to a `WholePool`/`ByLead` cluster
(`callofthearchons/horsemen_cluster.go`, Pestilence leads); Plague Rat wired to an
inline `SelfPull(3, 5)` cluster (`ageofascension/plague_rat.go`); the three exact
pullers migrated to `PullExact` clusters — Timetraveller→Help from Future Self
(`callofthearchons`, `timetravellerCluster`), Hyde→Velum and Igon the Green→Igon the
Terrible (`worldscollide`); and **every** remaining puller migrated to `Pull`
clusters — Grumpus Tamer→War Grumpus, Ortannu→Ortannu's Binding, Bear Flute→Ancient
Bear, Faygin→Urchin, Troop Call→Niffle Ape/Queen, Chain Gang→Subtle Chain, and the
nine ship blasters→their officers. No card uses `card.Connects` any more.

### 1. Retire the dead Connection type

Every puller is now a cluster, so the whole `Connection` mechanism is dead code kept
alive only by its own tests. Retire it for one placement mechanism, not two:

- Remove the `card.Connects`/`card.Pull`/`card.PullSometimes` facade
  (`internal/card/generation.go`), the `Connection`/`ConnectedCard` types and
  `GenerationProfile.Connection` (`deckgen/materialize.go`), and the
  `expandConnections`/`pullCopies`/`placePartners`/`freeSlot`/`pullCopies` legacy
  path plus its `g.expandConnections(...)` call site in `generate.go`.
- Drop `connection_test.go` and the `Connection`-branch of the reachability tests in
  `cards_test.go` (`TestReferencedCardIsConnected`/`TestConnectedCardIsPulled` now
  only need the cluster branch).
- Keep 100% coverage as the dead code and its tests come out together.

### 2. Plants (template, partner house)

The Ambassador cycle is **done**: the six concrete AoA Ambassadors are collapsed
into one Sanctum template (`ageofascension/house_ambassador.go`, "House Ambassador")
that materializes the partner House's Ambassador from `DeckHouses`, keyed on the
engine House enum so a new House gets its Ambassador for free. **Plants (Shadows)
remain unimplemented** — build them as the same kind of template. Plant's ability
is a `TriggerAfterChooseHouse` reaction (already exists, ADR 0035); the novel part
is only the partner-house binding (materialize per non-Shadows House in
`DeckHouses`, mirroring `ambassadorFor`). First find the plant source cards — grep
the set
provenance JSONs (esp. `ageofascension.json`) for the plant cycle to confirm the
concrete names, home House, and collector numbers before authoring the template.

### 3. Maverick houses complete every OnePerHouse cycle to all 9 Houses

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

### 4. Plain per-House bane (template, trait-derived)

- Compute most-common-Trait per House **once at init** from the registry, cached
  (no committed `mage gen` table — the registry is already frozen per binary, same
  as the pool). O(1) at generation; no MCTS cost.
- **Deterministic tie-break** among tied top Traits: seed a shuffle from
  `hash(sorted(names of the cards in the tied top Traits))` — order-independent, and
  not alphabetical (which would bias low letters). Pick from the shuffle.
- Migrate master-of-N onto a single template entry too, for Axis-B consistency.
- Plain bane name = `"<Trait>'s Bane"` ("Demon's Bane"); handle trait word-form
  corner cases (already-possessive, plural, multi-word) as they arise.

### Deferred (own ADR)

- **Boosted combinatorial bane**: choose 3 Houses, "Play: Destroy an x, y, and z
  creature", C(9,3)=84 variants, rarity `1/84·Rare`, generative name by splicing
  consonant-vowel-consonant sections of the three Traits. Persistence rides the
  same computed-at-init table + `snapshotVersion` discipline. Write its ADR before
  building.

### Cleanup found in current deckgen/materialize/connections (do during the migration)

- `expandConnections` (internal/deckgen/generate.go) fuses gather-pullers,
  roll-counts, count-present, and place into one long method with intertwined
  `wanted`/`present`/`counted`/`protected`/`rehouse` bookkeeping. Decompose while
  generalizing to deck scope, or it gets worse. Candidate for the refactor-sweep
  skill.
- `Connection`/`ConnectedCard` are stringly-typed (`Name string`, resolved via
  `set.lookup`); the cluster API should keep the compile-time card-symbol linkage
  `Connects`/`Pull` already give and not regress to raw strings.
- `docs/deck-generation.md` §5 still describes the "connections" model — rewrite to
  "clusters + strategies" as part of the migration so the doc does not diverge from
  ADR 0036 / CONTEXT.md.
- `deckgen.SlotContext`'s `Maverick/Legacy/Special` bools and the "templates — a
  future addition" comments in materialize.go / generation.go become stale once the
  first real templates land; refresh the comments then (ADR 0006/0020).

## Engine distillation sweep — 20 findings (one-off nodes → cohesive systems)

The theme: three sets now cover the game's surface, so card-named one-off nodes
that fuse several effects for a single card should collapse into small
parameterized systems, and the Resolver/log families they spawned should collapse
with them. Grouped by mechanic so the shared primitive lands once.

### Deferred, with the reason (so they are not re-proposed blindly)

- **H2 / H3 / H6** — collapsing the house-permission Resolver methods and the
  mirror log records (`FightGrantedForHouse` … `HouseForcedNextTurn` …) only pays
  off done _together_ with a fuller house-system pass (H1/H5 fully generalized,
  including the any-house and wager variants). Doing the logs now, while the
  effects are only partially merged, is churn. Bundle into one deliberate
  house-system ADR + pass.
- **H8** — the highest-value de-fusion (spend N cards → benefit scaled by N), but
  it needs a new `ctx.Produced`-style "cards spent this way" tally and careful
  100%-coverage work; worth doing with a review, not unattended.
- **H9** — `PlayFromOpponent{Source}`: modest (one node saved), touches the web
  client icon map, and the two behaviors genuinely diverge (random-from-archives
  vs top-of-deck). Low priority.
- **H10** — `NeighborsOfThis` + `NeighborsSharingHouse` already share the
  `neighbors()` helper; a merged `Neighbors{Subject, SameHouseOnly}` adds two
  unused subject×filter combos for near-zero dedup. **Keep separate** (rejected).
- **H11** — `SourceNeighborsAllOfHouse` is a lone one-off; fold only if a second
  neighbor-condition card appears. No rule to write yet.
- **H16 / H17** — rejected: the "discard-random / owner-gain-purge stragglers"
  have no clean pile-`Selection` target. The candidates (`PurgeEachOfChosenTrait`,
  Harvest Time; `PurgeArchivesForDamage`, Destructive Analysis = H8) are
  in-play / spend-cards composites, not `Selection`-over-a-pile nodes. (H13's
  shuffle-from-discard fold is done; `ShuffleNamedFromDiscardIntoDeck` stays its
  own node because its printed article "a Subtle Chain" is card-specific — the
  reason is a code comment there.)
- **H14** — `EachPlayer{Do}` combinator: the card discards _all_ hands then refills
  _all_, so `EachPlayer{Sequence{...}}` is wrong ordering; a correct `EachPlayer`
  is a controller-rescoping wrapper (subtle around "you", frames, and trigger
  ordering). Design deliberately, not unattended.
- **H20** — `Fraction` (creature count) vs the Æmber-pool "half loss": different
  subjects, different rounding contexts. **Keep separate** unless a card wants a
  pool `Fraction`; record as rejected for now.

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

## Card wording / authoring

- **`card.Name` typed card references.** Introduce a `card.Name` named string and
  convert the authoring-facing name fields (`AttachSelfTo.Host`, `Target.Named`,
  `ArchiveGrantingUpgrade`, `ControlsNamed`, `NamedCardPurged`, `ItIsNamed`,
  `InPlay.Name`, `PutFromDiscard.Name`, `ReturnNamedToHand`, `SearchForName`,
  `ShuffleNamedFromDiscardIntoDeck`). The engine-internal resolver `Name(id)
string` port stays `string`. Rationale (durable): a card that references another
  card by name is always within its own set and connected to it, so the value set is
  closed — the type prevents arbitrary strings.
  NEEDS A DECISION (attempted, reverted): the change does **not** compose with the
  codebase's actual authoring pattern. Cards reference another card not by a string
  literal but by `OtherCard.Name` — the referenced card's `CardDefinition.Name`,
  which is `string`. A named `CardName` field does **not** accept a typed `string`
  variable, so every reference (~30 call sites across the set packages) would need
  `card.Name(OtherCard.Name)` wrapping — uglier authoring, defeating the ergonomic
  goal. The only clean alternative is to make `CardDefinition.Name` itself a
  `CardName`, but `def.Name` is consumed as a plain `string` in 50+ engine sites
  (logs, text rendering, invariants, the `Name(id) string` resolver), so that
  ripples `string(...)` conversions through the whole engine for marginal type
  safety. Pick one before implementing: (a) accept `card.Name("Literal")`-style
  references and drop the `OtherCard.Name` pattern; (b) make `CardDefinition.Name` a
  `CardName` and chase the engine-wide `string()` ripple; or (c) drop the item.
- **Uncharted Lands self-reference.** Replace `Named("Uncharted Lands")` with a
  self-referencing target (`card.Target.Source`) that resolves to the granting card
  and renders its own name. Depends on the `card.Name` work.
- **Generic zone-move `Match` predicate.** Replace the `Type` / `Trait` / `OrTrait`
  triple on `PutFromDiscard` (and its zone-movement siblings) with one composable
  `Match` predicate expressing a union of type/trait clauses (e.g. `Upgrade` OR
  trait `Robot`). Gives one place to test "any match in the discard".
  DEFERRED (nicety, not a bug): the concrete **Chief Engineer Walls** bug — the
  `May` prompt appearing when the discard has no upgrade or Robot — is already
  fixed; `PutFromDiscard` now implements `vacuous()` (shared `matches` helper), so
  `May` skips the empty choice. The `Match`-predicate consolidation across the
  zone-movement siblings remains as an optional refactor.

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

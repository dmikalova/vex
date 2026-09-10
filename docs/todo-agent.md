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

## Engine refactor sweep (survey of `internal/engine`)

File-naming decision for the splits: keep each family's existing top-level prefix
so `ls` groupings stay intact — `target.go` is a bare concept file, so its splits
are `target_*.go`; `effect_condition.go` / `effect_count.go` are part of the
`effect_*` vocabulary, so their splits keep the prefix (`effect_condition_*.go`,
`effect_count_*.go`).

Structural (decompose / atomize / recompose):

- **A3** — Selection mode (chosen / random / each) is a new `Selection` Strategy
  (distinct from the set-relative `Refinement`), not a node per (verb × zone ×
  mode). Concrete strategies `Chosen{House, Mandatory}` / `Random{}` /
  `Each{Type, ExceptHouse}` each render their own fragment ("a card" / "a random
  card" / "each creature") and carry their own mode-specific fields on the
  strategy (no union fields on the node). Express declinable / `Then`-gate / tally
  as optional-capability interfaces mirroring `leadingRefinement` /
  `OptionChooser` / `Orderer`.
  - **Done:** the `Selection` strategy (`effect_selection.go`) + the hand-purge
    fold — `PurgeFromHand` / `PurgeRandomFromHand` / `PurgeEachFromHand` are now
    one `PurgeFromHand{Player, Selection}`. Added the `ChooseRandom` resolver-port
    method so `Random` returns an id the node moves like `Chosen`/`Each` (removed
    the orphaned `PurgeRandomFromHand` port method). Also folded
    `PurgeCreatureFromHand` in: `Chosen` gained a `Type` filter, and
    `PurgeFromHand` now puts a single purged card in context (`ctx.It`) so Custom
    Virus's trait-share destroy still finds it. Then the discard fold —
    the discard fold — `DiscardRandomFromHand`, the chosen `DiscardFromHand`, and
    `DiscardHand` ("each") folded into one node, then the cross-zone step:
    `DiscardRandomFromArchives` merged in and the node renamed
    `DiscardFromHand`→`DiscardCard{Player, Zone, Selection, Amount, AnyNumber}`
    across all 17 hand cards plus Tantadlin. `Zone` (invalid-zero, Hand|Archives)
    picks the source; the archives path added an `Archives` reader and a
    `DiscardCardFromArchives(owner, id)` mover (RNG-equivalent — one `Intn(len)`
    draw, and it fires no hand-discard reactions). Tantadlin was reworded to owner
    voice ("Your opponent discards a random card from their archives"), so
    `selectionOwnerActs` drives voice uniformly with no zone condition; that voice
    change is recorded in `docs/keyforge-divergences.md`. `Each` gained an
    `OfChosenHouse` filter (Deep Probe). Removed the `DiscardRandomFromHand`,
    `DiscardHand`, and `DiscardRandomFromArchives` nodes, their resolver-port
    methods, the now-dead `typeNoun`/`matchesTypes`/`cardTypeNoun` helpers, and the
    unused `card.Types` facade.
  - **Remaining — deferred; the last two mirror-folds are low value.** A usage
    census settled it: each remaining verb-mirror removes **one node used by ~one
    card** while carrying a real wrinkle, so folding them one at a time is a poor
    trade against the churn:
    - `ArchiveRandomFromHand` — **1 card**. Its target `ArchiveFromHand` is the
      richest zone node (`Per`/`Or`/`UpTo`/`Revealed` + filters, **18 cards**).
      Adding a `Selection` there risks destabilising 18 cards to absorb 1.
    - `PurgeEachFromDiscard` (**1 card**, Soldiers to Flowers) + `PurgeCard`
      (**6 cards**) do **not** fold under a generic `Chosen`: `PurgeCard` commits
      to **one** discard pile before purging `Amount`/`UpTo` (Creeping Oblivion:
      "up to 2 cards from **a** discard pile"), while flattening both piles into
      one candidate list would let a multi-purge cross piles — a behavior change.
      The pile-commit is pile-scoping logic `Chosen` does not model, and
      `PurgeEachFromDiscard` also carries a `GainOwnerAember` rider.
    - **Recommendation:** treat the cross-zone unification as one deliberate
      architectural pass (the real ADR 0031 payoff — one node per verb across all
      zones, absorbing the voice fork, the two-pile commit, and the rich archive
      node together), **or** leave it here: the high-value same-zone folds
      (purge-from-hand 3→1, discard-from-hand 3→1) are done and the residue is
      low-value. Do not grind the mirrors one at a time.

## Engine distillation sweep — 20 findings (one-off nodes → cohesive systems)

The theme: three sets now cover the game's surface, so card-named one-off nodes
that fuse several effects for a single card should collapse into small
parameterized systems, and the Resolver/log families they spawned should collapse
with them. Grouped by mechanic so the shared primitive lands once.

### Done this pass

- **H7** — `pickCards(ctx, prompt, limit, optional, avail)` in `target_select.go`
  now backs the five copy-pasted pick-one-at-a-time loops; ratcheted into
  `internal/engine/AGENTS.md`.
- **H18** — `SkipForgePhase` moved to `effect_forge.go` beside the forge family.
- **H5 (partial)** — `ForceOpponentActiveHouse` / `ForbidOpponentActiveHouse` /
  `ForceOpponentActiveHouseOfFought` folded into one
  `ConstrainOpponentActiveHouse{Source, Forbid}` (Source ∈ chosen/fought). The
  wager (`WagerOpponentChoosesChosenHouse`) and mirror
  (`ForbidSameActiveHouseNextTurn`) stay separate — their resolution and POV
  genuinely diverge.
- **H1 (partial)** — `GrantFightForChosenHouse` + `GrantFightForFriendlyHouse`
  folded into one `GrantFight{House}` (HouseNone = chosen house). `GrantFightAnyHouse`
  stays separate — it calls a different Resolver method.
- **H19** — `UnforgeKey`-vs-`ForgeKey` rejection already lives in the code comment;
  no doc change needed.

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
- **H12** — `SearchForName` + `SearchDeck`: **DONE (Q3, de-fused not merged).**
  Kept the two searches separate but pulled the trailing shuffle out of
  `SearchDeck` into a standalone `ShuffleDeck{}` effect, so a search never bundles
  its own shuffle (fixed Grumpus Tamer silently skipping its shuffle). The
  search-implies-shuffle rule is now enforced by `TestSearchIsFollowedByShuffle`
  (a reflection lint: any ability whose effect tree contains a search type must
  also contain a `Shuffle*` type).
- **H13 / H16 / H17** — these are all the A3 zone-movement/`Selection` fold, not
  standalone merges: `ShuffleNamed`-vs-`ShuffleMatching` is single-vs-each
  resolution divergence, and the discard-random / owner-gain-purge stragglers are
  `Selection.Random`/`Selection.Each` cases. Do under the A3 pass (already flagged
  above as a deliberate, non-blind pass), not piecemeal.
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
- **Generic `Repeat{}` node (decompose `ExaltToRepeat` + `RepeatWhile`).** One node
  `Repeat{Do, Max, While, Cost}`: `Max` is a hard iteration cap, `While` an optional
  condition re-checked each iteration (they may combine), `Cost` a small **cost**
  strategy (call it `Cost`, not `Pay` — in KeyForge vernacular "pay" means give the
  opponent Æmber; exalt is the only concrete cost today). Text must key off the
  semantics: `Max: 1` renders **"repeat the preceding effect"** (once); an unbounded
  repeat renders **"repeat this effect"** (Rule of Six caps it at 5 more). The
  current bug is `ExaltToRepeat` always emitting "the preceding effect" while
  repeating unboundedly. Phalanx Strike and **Tribute** become
  `Repeat{Max: 1, Cost: exalt}`; Bait and Switch becomes `Repeat{While: …}`.
- **`Reap:` and "after a creature reaps" share one ordered timing window.** A
  reaping creature's own `Reap:` ability and every lasting "after target reaps"
  reaction (Bloodshard Imp, Candle Unit) resolve in the **same** window and the
  active player orders them (KeyForge lets the active player order simultaneous
  triggers; ADR 0013). Today the creature's `Reap:` fires ahead of the lasting
  reactions as a separate step. Fold the reaper's own triggered ability into the
  same `emitLasting(EventReap, …)` gather-and-order pass so all of them queue
  together and are orderable.
- **Constant abilities gain from the card that just entered play (Harmonia,
  Hunting Witch).** When a card is played its constant abilities take effect and
  its bonus icons resolve **before** the after-play window opens, so a
  self-referential "gain when a creature enters play" constant ability sees the
  card itself as an eligible entrant. Today Harmonia and Hunting Witch do not gain
  from their own entrance. Sequence play as: enter play → constant abilities live
  → resolve bonus icons → open the Play/after-play window (with the new card
  already counted).
- **Cooperative Hunting → an iterator over `DealDamage`.** "Deal 1 damage X
  times, choosing any creature each time" is distinct from Sack of Coins ("deal X
  to one creature"): Cooperative Hunting reuses the `DealDamage` primitive wrapped
  in a repeat/iterator that re-prompts for a target each iteration. All selections
  are made first, then all the damage is dealt **simultaneously**. Model as the
  generic `Repeat{}` node (above) wrapping `DealDamage{Amount: 1}` with a
  per-iteration target choice and deferred simultaneous application.

## Card wording / authoring

- **`card.Name` typed card references.** Introduce a `card.Name` named string and
  convert the authoring-facing name fields (`AttachSelfTo.Host`, `Target.Named`,
  `ArchiveGrantingUpgrade`, `ControlsNamed`, `NamedCardPurged`, `ItIsNamed`,
  `InPlay.Name`, `PutFromDiscard.Name`, `ReturnNamedToHand`, `SearchForName`,
  `ShuffleNamedFromDiscardIntoDeck`). The engine-internal resolver `Name(id)
string` port stays `string`. Rationale (durable): a card that references another
  card by name is always within its own set and connected to it, so the value set is
  closed — the type prevents arbitrary strings.
- **Uncharted Lands self-reference.** Replace `Named("Uncharted Lands")` with a
  self-referencing target (`card.Target.Source`) that resolves to the granting card
  and renders its own name. Depends on the `card.Name` work.
- **Generic zone-move `Match` predicate.** Replace the `Type` / `Trait` / `OrTrait`
  triple on `PutFromDiscard` (and its zone-movement siblings) with one composable
  `Match` predicate expressing a union of type/trait clauses (e.g. `Upgrade` OR
  trait `Robot`). Gives one place to test "any match in the discard" — which also
  fixes **Chief Engineer Walls**: skip the prompt as a vacuous choice when the
  discard has no upgrades or robots.

- **Orator Hissaro** could read: "Play: Exalt and ready each neighboring creature.
  For the remainder of the turn, those creatures belong to house Saurian."
  Blocked on two engine additions: making `Exalt` a `combinable` (so `Ready`+`Exalt`
  fold to "ready and exalt each neighboring creature" — but the fold must keep
  "exalt N times" for `Amount > 1`), and a pronoun form of `BelongToHouse` so the
  second sentence reads "those creatures belong to …" instead of repeating the
  target. Today it renders correctly but verbosely (target repeated three times).
- **Borr-Nit** and similar could be atomized and recomposed further (decompose
  fused effects into shared nodes).
- **Tentacus — reverse the clause direction.** Its ability is a cost the opponent
  pays to act, so it should read as an "in order to" restriction:
  "In order to use an artifact, your opponent must give you 1 Æmber" — the
  restriction clause names the price and who pays it, rather than framing it from
  the controller's side. Check the wording conventions in
  [card-wording-rules.md](card-wording-rules.md) and whether a reversed-clause
  restriction node already exists before adding one.

## Card catalog / provenance

_No outstanding items._

## Web — mobile, previews, layout

- **Player-bar swipe through icons.** A touch that begins inside `.score-pill`
  disables swipe tracking, so ending on an icon/`.tip` kills the horizontal scroll.
  Let the player bar scroll horizontally on swipe even when the touch starts on an
  icon.
- **Player-bar tooltips render inside the bar.** Tooltips (CSS `.tip::after`,
  z-index 30) should overlay above the bar but now appear clipped inside it — fix
  the overflow/stacking regression so they float over the bar again.
- **Creature-as-upgrade play buttons.** Replace the single `Play` button above a
  lifted creature-that-can-be-an-upgrade with two explicit buttons, `Play creature`
  and `Play upgrade`. `Play upgrade` skips the flank step (upgrades take no flank);
  `Play creature` keeps the flank prompt. Removes the later Creature/Upgrade sidebar
  prompt.
- **Zone-select modal jumps to top.** After selecting a card from the zone viewer,
  a deferred `OnUpdate` remeasure repositions the modal/lifted card so it jumps to
  the top a moment later. Stop the jump.
- **Sidebar warnings one per line.** `.restrictions` is a wrapping flex row, so
  Stealth Mode and Sensor Chief Garcia share a line when wide. Make each warning its
  own line (column layout).
- **Narp + Universal Translator illegal reap.** Universal Translator's "use" path
  offers Reap on a creature neighboring Narp even though Narp forbids it; gate the
  offered actions through the restriction check (`CanUseTo`) so the illegal Reap is
  omitted.
- **Universal Translator use-buttons above the creature.** When Universal
  Translator selects a creature, show the normal use buttons (Fight / Reap / Action)
  **above the creature**, like a normal selection, instead of as sidebar option
  buttons.
- **Text-less cards keep an empty textbox.** Cards with no trait/rules text (e.g.
  Virtuous Works) currently skip the textbox; render an empty black textbox that
  fills the card space instead.
- **Manual mode: graft / place under / to hand.** With a card selected, add zone
  options **"Graft"** (face up) and **"Place under"** (face down) that then prompt
  the user to click an in-play host (creature or artifact). Add a **"To hand"**
  action available when an upgrade or a card-under-a-host is selected, sending it to
  hand.
- **Deploy: choose the side first, then the creature.** Today the prompt is
  {left-flank button, deploy-left toggle, deploy-right toggle, right-flank button}
  and you deploy by clicking a creature with the toggle set. Replace it with a
  {Deploy left, Deploy right} button pair; once the side is chosen, present the
  creature selection. (Engine already models Deploy via
  `deployPosition`/`chooseFlank`/`choosePosition`; this is the client prompt flow.)

## Logging

_No outstanding items._

## Tooling / tests / docs

- **Design-patterns audit + doc.** Inventory the patterns actually in use
  (Interpreter/effect AST, Visitor, Strategy — Chooser/Refinement/Count/Condition,
  interface-segregated Resolver port, functional-options builder, value-type undo
  snapshot), check which are already documented (ADRs, style-guide, engine AGENTS),
  and fill the gaps in one place.

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
- **Sequence the work yourself — never ask the human what order to build in.**
  Once work is approved, deriving the order is your job, not a decision to hand
  back. Order by dependency first (a shared primitive or a renderer helper lands
  before the cards that consume it), then easiest-win first inside each
  dependency tier (the mechanical, no-rules-risk change before the one that needs
  a judgement call). Say what the order is and why in one line; do not offer a
  menu of plans.
- **When you need a decision from the human, ask in the reply itself, grill-me
  style** — a numbered list of `❓ **Q1** - **title**: <question>` with a `➡️`
  recommended answer under each — at the end of the turn, then stop. Do **not**
  reach for an interactive question tool: under autopilot it is auto-answered with
  "work autonomously" and the human never sees it. Plain-text questions at the end
  of the turn are the channel the human actually reads.

---

## gocognit gate: web renderers deferred to their own passes

The engine pass is done: the `gocognit` gate's global `min-complexity` is now 40,
with every genuinely nested engine/deckgen seam split (fight, NewCard, hasKeyword,
the text.go renderers, Target.Text, Target.filter) and 100% coverage held. What
remains are the web renderers, kept exempted by name in the `.golangci.yaml`
`exclusions` `text:` list until each gets its own pass:

- `effectGlyphs` (internal/web/icon.go, 90) — glyph dispatch; wants its own
  glyph-family grilling session (see docs/todo.md "Split out glyphs more in
  icon.go"). Split along glyph families, then delete its exclusion.
- `installTips` (internal/web/game_lifecycle.go, 49) — tip installation; split by
  tip group, then delete its exclusion.

When both are split, lower the global toward the tool default (30) and re-check
the 31–39 engine stragglers (`triggeredBy`, `RenderAbility`, `Harness.location`,
`allowedHouses`) — left untouched this pass — one at a time.

## Post-Mass-Mutation cleanup sweep

Everything below was decided in a grilling session after Mass Mutation landed.
The decision is recorded with each item — **the decision wins over a
contradicting code comment**, which is by definition describing the behaviour the
item exists to change. Items marked **INVESTIGATE** were approved in principle
but need a design answer before code; report the answer back before building.

Some items say "and find the others like it". Those are real work: the named card
is the example that surfaced, not the whole set.

### Card text rendering

The printed text a card's nodes render. Fixing a node's `Text()` changes every
card that uses it; hand-editing a card's doc comment is useless because
`mage gen` overwrites it.

- **Threshold ladders must not repeat the clause.** `galactic_census` and
  `a_fair_game` render the same 9-word clause three times ("If there are 3 or more
  houses represented among creatures in play, gain 1 Æmber. If there are 5 or
  more … If there are 6 or more …"). Group it — but the grouped wording must stay
  obvious to a player; the first attempt ("for each of 3, 5, and 6 houses") was
  rejected as unclear. Try variants and bring back two or three before committing.
- **Symmetric cards use the `EachPlayer` template.** Two 310-char cards render the
  same paragraph twice, once per player; `binate_rupture` and
  `champions_challenge` are the same shape. Fold them onto the existing each-player
  mechanism. Order is the **rules default** — the active player chooses who
  resolves first for each player — not a hardcoded "starting with you". Sweep for
  the other candidates while you are in there.
- **Ambiguous `it` across a sequence.** `destructive_analysis` renders `Deal 2
damage to a creature and purge any number of cards from your archives, and for
each card purged this way, deal 2 damage to it` — by the second clause `it`
  could be the purged card. **INVESTIGATE**: first establish whether `ctx.It`
  still holds the creature (text bug) or not (behaviour bug), then find the other
  cards where a later clause re-binds `ctx.It` before a back-reference. A general
  fix is preferred over per-card rewording.
- **Conjoined `It*` conditions collapse in the renderer, not the card.**
  `mercy_malkin_queen` renders `if it is a friendly creature and it is a Cat
creature`; it should read `if it is a friendly Cat creature`. **Decision: the
  renderer does this automatically** — do not fix it card-by-card.
- **`city_gates` keeps its two-capture shape, with clearer text.** The card
  deliberately diverges from printed KeyForge to drop the "otherwise": the subtle
  resolution difference is accepted. Render it `A friendly creature captures 1
Æmber from your opponent. If it is a Dinosaur creature, it captures 1 Æmber from
your opponent.` Look for the other `Otherwise`/`Choose one:` cards that can take
  the same simplification — **not all of them can**, so judge each.
- **Sweep for the other cards that want `Target.MatchingAny()`.** The union axis
  now exists (`emp_blast` reads `Stun each Mars or Robot creature`, and a Mars
  Robot is one member of that set rather than two). `MatchingAny` disjoins the
  house and trait axes only, and degrades to a conjunction when there is no second
  adjective to join. Hunt the remaining cards rendering `each X … and each Y …`
  over the same base noun.
- **Towards deleting `card.Sentences{}` (foundation only).** The goal is for
  `card.Sequence` to work out its own sentence breaks so a card implementation has
  exactly one text representation. `help_from_future_self` (`…put it into your
hand, and shuffle your discard pile into your deck` — a dangling conjunction
  joining two complete operations) is the motivating case. Full delivery is out of
  scope for this sweep; **lay the foundation**: give `Sequence.Text()` a notion of
  whether a rendered clause is a full sentence, and use it for the obvious cases.
  Do not remove `Sentences{}` yet.
- **Card-text fan-in: find the remaining ungrouped repeats.** The
  creature/artifact/upgrade trigger fan-in already exists; the sweep did not find
  a case where it fails. Hunt for the other renderers that join whole clauses
  instead of grouping the varying noun.

### Game log

Every change to a rendered log string requires bumping `snapshotVersion`
([internal/web/game.go](../internal/web/game.go), currently 19), because persisted
logs store rendered prose. **Do these as one batch so the version bumps once.**

The log's voice rules — lowercase mechanics, source card first, no left-to-right
backtracking — are now written into `internal/engine/log.go`'s file doc comment,
so a new entry is written to them. Keep that comment as the authority; the items
below are what is still unapplied.

Items:

- **Finish the "action" audit.** The log entry, its emitter, and the play path
  are renamed to Tactic (`TacticPlayed`, `playTacticCard`,
  `emitTacticPlayedBeforeResolve`, `consumeNextTacticToHand`). What remains is the
  rest of the tree: every identifier and comment still saying "action" must be
  either a genuine `Action:` ability or renamed to Tactic. `ActionAbilityUsed` and
  `Trigger.Action` are correct as they stand — they name the `Action:` ability.
- **Bonus-icon entries keep an explicit "bonus", with one shared helper.** The
  four entries (`log_bonus.go:22` Æmber, `:57` capture, `:71` damage, `:84` draw)
  should share a framing helper, but each verb's wording is peculiar — `bonus
draw` was deliberate. **INVESTIGATE**: propose a template plus its rendering for
  **each** of the four bonus icons, keeping "bonus" explicit, and bring them back
  together before implementing.
- **Audit every remaining parenthetical in the log.** The exalt one is gone. The
  only place a parenthetical has seemed reasonable is fights, and that is itself
  slated for rework — so treat a parenthetical as a defect unless argued.
- **Make "every state change is logged" a tested invariant.** `Exhaust` and
  `ReadyIfFirstUse` now emit their entries, which closes the two known gaps — but
  nothing stops the next one. Walk the state transitions and assert an entry.
  Expect a few legitimate exceptions; the point is to discover what they are and
  push the test as far as it goes before reaching for an escape hatch.
- **Improve `countNoun`'s discoverability.** The `1 keys` bug existed because the
  helper was there and not found. See the engine-implementation doc item below.

### Zone movement: one mechanism per KeyForge verb

The unifying decision: **each KeyForge verb — archive, discard, purge, put — is
one mechanism with a source/destination axis, not a family of near-duplicate
nodes.** The natural-language node (`Archive{…}`) should cover its verb's whole
surface.

- **Investigated; the shape is decided. Build it incrementally, but keep going
  until the whole sweep is done.** The four verbs stay as the authoring surface;
  what folds is everything under them. Findings and decisions:
  - **The private mover has landed.** `game_move.go` holds `moveCard(id, from, to
zoneRef, entry LogEntry) bool` over a `zoneRef{Player, Zone}` pair and a
    `cardPile` interface (archives are a `wideList`, the other four piles are
    `deckList`s, so the mover reaches all five through `add`/`remove`). The nine
    repeats in `game_archive.go`, `game_purge.go`, and `resolver_game.go` now go
    through it. The "resting zone" concept is gone \u2014 the existing `Zone` enum is
    the axis \u2014 and play is deliberately excluded, since entering and leaving it
    runs upgrade/damage/ward lifecycles the mover must not know about.
  - **Build on `crossZoneMover`** (`effect_cross_zone.go`, `{Player, Dest,
Sources}` with `zoneCards`/`gather`/`originOf`/`inZone`), which already backs
    `ReturnNamedToHand`, `Search`, and
    `ShuffleIntoDeck`. It is the fold's destination, not a rival.
  - **Start with purge** — the smallest complete verb (`PurgeCard` 9 cards,
    `PurgeFromHand` 7, `PurgeCreature` 21, `PurgeSource` 3, `PurgeArchives` 1,
    `PurgeArchivedCardThen` 1). Done; see the two entries below.
  - **Purge's movement is now fully on the seam.** `Destination.moveFrom`'s
    `destPurged` matrix was missing its `Archives` and `Deck` rows, so
    `PurgeArchives`, `PurgeArchivedCardThen`, and `effect_deck.go`'s `IntoPurge`
    closure each called `ctx.Resolver.PurgeFromX` directly and bypassed ADR 0031.
    The rows are filled in and all four call sites (plus `PurgeSource`'s in-play
    branch) now go through `purgeFrom`. **`Resolver.PurgeFromPlay` is the only
    purge the effect layer still names directly, and only from `MarkPlayedActionPurged`'s
    sibling branch — revisit it with the resolution-zone question below.**
  - **Correction: the purge nodes do NOT fold into one, and here is the evidence.**
    The earlier "pure subsets dissolve" list was drawn from names and usage counts,
    not from reading the code; reading it contradicts the list. `PurgeCard` (from a
    discard pile) and `PurgeFromHand` differ in more than a source zone: `PurgeCard`
    _chooses among piles_ (`piles()` prompts when both hold a match) while
    `PurgeFromHand` names a player outright; `PurgeCard` renders "up to N" through
    `object()` while `PurgeFromHand` renders "you may purge"; `PurgeCard` tallies
    `PurgedAemberBonus` and adds to `Produced.Purged` while `PurgeFromHand` assigns
    it and puts a lone card in `ctx.It`; only `PurgeFromHand` implements the
    `declinable`/`resolveOptional` May protocol. `PurgeArchives` is "any number",
    which no `Selection` expresses — `Chosen` picks exactly one card, so
    `PurgeCard` fakes "up to N" by looping `count()` times. Merging these would
    produce one node whose fields are only valid in particular combinations, which
    is the fused-effect anti-pattern wearing a fold's clothes.
    **So the movement is the shared thing, and it is now shared. Selection
    semantics and printed-text templates are what make these separate nodes, and
    they should stay separate.** Re-check the other three verbs the same way — by
    reading each node, not by comparing names — before folding anything.
  - **Re-opened: the pile-purges DO look foldable, and the axis list is below.**
    The "they are genuinely different" verdict above conflated two kinds of
    difference. Sorting them honestly: most are **accidents of implementation
    history**, and only two are **rules**.
    _Accidents (a merged node should just do the right thing everywhere):_
    `PurgeFromHand` assigns `Produced.Purged` where `PurgeCard` accumulates it
    (accumulating is correct across two piles); `PurgeCard` tallies
    `PurgedAemberBonus` and `PurgeFromHand` does not, though cards in hand carry
    bonus icons too; `PurgeFromHand` binds a lone card to `ctx.It` and `PurgeCard`
    does not, though a single-card purge is nameable either way; only
    `PurgeFromHand` implements the `declinable`/`resolveOptional` May protocol.
    _Rules (these must survive any fold):_ a purge **from play** selects with a
    `Target` rather than a pile `Filter`, so `PurgeCreature` stays its own node \u2014
    note this is now a **targeting** distinction only, not a lifecycle one, since
    the leave-play teardown moves under the mover (below); and a **chooser can
    only choose from a zone it can see**, which the two entries below turn from a
    validation rule into computed behavior.
    _Printed text is derivable, not a difference:_ "a discard pile" vs "your
    opponent's hand" is just ChosenPlayer-vs-named crossed with the zone noun.
    **Proposed shape**, following the human's steer (always name the zones, take
    an array, split may/up-to into strategies):
    `Purge{Zones []Zone, Player Player, Filter CardFilter, House HouseMatcher,
Count Count, GainOwnerAember bool}` where the pick strategy is `Count`
    (`Exactly{N}` / `UpTo{N}` / `AnyNumber{}` / `All{}` / `Random{N}`) rather than
    today's `Selection`+`Amount` pair \u2014 which is what removes the invalid
    combination, since `Amount` is meaningless beside an `Each` Selection and
    "any number" (`PurgeArchives`) is not expressible as a `Selection` at all.
    `Sources []Zone` has precedent on `Search` and `crossZoneMover`. Build it on
    `crossZoneMover`, and do the same pass for Archive/Discard/Put/Shuffle after.
  - **Better than validating visibility: compute it, and let the chooser degrade
    to random by itself.** Instead of rejecting `Chosen` against a hidden zone,
    give the engine `visibleTo(zone, owner, chooser)` and have the pick strategy
    consult it \u2014 a chooser aimed at a zone it cannot see picks uniformly at random.
    That makes the invalid combination unrepresentable rather than merely rejected,
    and it pays off twice, because `Text()` derives "at random" from the same fact,
    so printed text cannot desync from which zone a card reaches into (ADR 0006).
    **It is not "hidden zone \u21d2 random", though**, and modelling it that way would
    be wrong: visibility is a function of all three of zone, owner, and chooser \u2014
    a discard pile is public, a deck is hidden to both players, and archives are
    hidden to the opponent but visible to their owner. A card can also **grant** a
    look (Philophosaurus looks at the top 3 of its own deck), so the grant has to
    be an input to the function, not an exception around it. **Subsume, do not sit
    beside, the half-model that already exists:** `ownerActsSelection` /
    `Random.ownerActs` (`effect_selection.go:79-90,214-216`) already encodes "a
    random pick from a hidden hand is attributed to its owner". That is the same
    idea reached from the other end and must fold into the new function.
  - **The leave-play lifecycle should move under the mover, and most of it already
    has.** `removeFromPlay` (`game_leaves_play.go:99-108`) is already documented as
    the funnel every real exit runs through, and `leavePlayDestroyed` is already
    the shared teardown (unlist, discard upgrades, discard under-cards, release
    \u00c6mber, reset core). `discardDestroyed` and `purgeFromPlay` are then the _same
    five lines_, differing only in destination pile and log entry \u2014 which is
    exactly `moveCard(id, from, to, entry)`. So lift the in-play exclusion in
    `game_move.go` and have the mover run the teardown when the source is play;
    purge, archive, bounce, and shuffle then stop knowing about the lifecycle.
    **Hazard, stated precisely: a ward check belongs at a removal _attempt_, never
    at the _filing_ of a removal already settled.** Two ward checks are not by
    themselves wrong \u2014 destruction and purge are separate attempts, so a creature
    whose ward absorbs a destruction, is granted a fresh ward by another creature's
    Destroyed ability, and is then purged absorbs the purge too. The engine already
    does this correctly: `filterUndestroyed` consults ward at the destruction
    attempt, `purgeFromPlay` at the purge attempt, and `discardDestroyed` \u2014 the
    filing step that runs after the Destroyed window resolves \u2014 deliberately does
    not, so a ward granted mid-window cannot resurrect a creature the destruction
    already claimed. **Both rules are now pinned by
    `TestWardAbsorbsDestructionAndPurgeSeparately` and
    `TestDiscardDestroyedIgnoresWard` in `game_leaves_play_test.go`.** The
    consequence for the fold: `moveCard` must never consult ward itself \u2014 it tears
    down and files, and ward stays at the attempt sites above. Gigantic halves
    (ADR 0042) are uniform and fold cleanly.
  - **Still believed to be thin wrappers, but UNVERIFIED by reading:**
    `ArchiveDiscardedThisWay`, `DiscardTop`,
    `PutRevealedCard`, `PutDiscardedIntoHand`,
    `PutDiscardedIntoPlay`, `PutFromHand`, `PutNamedIntoHand`, `PutItIntoHand`,
    `PutChosen` (vs `PutFromPlay`). Verify each before touching it.
  - **`ArchivePurgedCard` was verified and is GONE.** It was
    `ArchiveCard{Zone: Purged, Selection: Chosen{}}` spelled out by hand. The purge
    pile is now an exported `Zone`, so a card may name it as a **source** (the only
    one that does is Universal Recycle Bin) while still never naming it as a
    destination — setting a card aside stays the `Purge` verb over the unexported
    `toPurged` (ADR 0031). `Zone.noun()` gained the "purge pile" arm it was
    silently missing (it fell through to "discard pile"). Text moved from "archive
    a purged card you own" to "archive a card from your purge pile", which is the
    template every other `ArchiveCard` already uses (rule 17 names the source zone).
    Pinned by `TestArchiveCardFromPurge` and `TestArchiveCardFromPurgeEmpty`.
  - **These carry a genuinely distinct rule and survive:**
    `ArchiveGrantingUpgrade` (`ctx.Upgrade`), `ArchiveFromPlay` (batch +
    `ctx.Produced.Archived`), `DiscardArchives` (active-player ordering +
    randomized opponent archives), `DiscardHand` (`playersInOrder`),
    `DiscardOpponentArchivesOrDeckTop` (two unlike hidden sources),
    `PurgeArchivedCardThen` (cost gate), `PurgeCreature` (play-or-discard
    fallback), `PurgeCard` (\u00c6mber bonus tally), `PutNextTacticIntoHand` (lasting
    replacement), the under-chain nodes (`game_under.go:117-123` \u2014 explicitly not
    zone movement), bare `Shuffle`, `ShuffleFriendlyCardsIntoDeck`.
  - **`Until{}` was investigated and deliberately NOT built; the real atom was the
    filter, and it has landed.** `DiscardUntil` is the only `*Until` node in the
    engine, so an `Until{}` combinator would have had no second consumer to prove
    its shape \u2014 speculative generality. What the dig _did_ duplicate was a card
    matcher: it hand-rolled `Type`/`Name`/`ExceptTrait` inline while `Search`
    already expressed the same axes as `Filter CardFilter` + `House HouseMatcher`.
    `CardFilter` gained the one axis it was missing (`ExceptTrait`, the negation
    `Target` already had), and `DiscardUntil` now carries `Filter`+`House` in
    exactly `Search`'s shape. **Revisit an `Until{}` combinator only when a second
    `*Until` card actually appears** \u2014 the loop body (reveal, record, test, offer a
    stop) is then the thing to extract.
  - **Fold `ArchiveSource`/`PurgeSource` into a `Self{}` selection.** The redirect
    half of this is **DONE**: `PurgePlayedAction*`/`ArchivePlayedAction*` are gone
    from `GameState`, replaced by per-card `CardCore.ResolvingDest` and one port
    method `RedirectResolvingCard(id, dest)`. Per-card was load-bearing — the old
    single-slot fields were a latent nesting bug, since a resolving Tactic can play
    another card that also redirects itself (Wild Wormhole into Causal Loop,
    `TestResolvingCardRedirectIsPerCard`).
    - **The investigated premise was wrong, and the truth is simpler.** The plan
      was to leave a resolving tactic in `State.Hand[owner]` with a `Resolving`
      flag so it would have a `from` zone. It does not need one: `playTacticCard`'s
      caller has already removed it from **every** zone, so a resolving card is in
      no zone at all. That is why Labwork cannot archive itself (the hand no longer
      holds it) and why the redirect is automatically generic over the origin zone
      — Wild Wormhole plays off the deck and the redirect still lands. No
      `Resolving` flag was added, because nothing reads one.
    - **What remains:** nothing structural. The duplicated "in play → move, else →
      redirect" body is now one seam, `Destination.moveSource(ctx)` in
      `destination.go`, and `ArchiveSource`/`PurgeSource` are each a one-line
      `Resolve` over it (`ToArchives.moveSource` / `toPurged.moveSource`). A future
      "return {self} to hand" is one more line.
    - **They were deliberately NOT folded into one `MoveSource{To: …}` node.** That
      would put the destination on the card, and ADR 0031 hides `toPurged` from
      cards precisely so a card writes the verb (`Purge`) and never the terminal
      destination. Two verb-named nodes over one shared seam respects that; one
      destination-parameterized node would not.
- **Purge is one mechanism** — partly done, remainder scoped.
  - **Done:** `PurgeFromHand` is gone. `PurgeCard` now carries a `Zone` (`Hand`
    or `Discard`) and its `Player` means **whose copy of that zone**, with
    `ChosenPlayer` keeping the "controller picks among the eligible sides"
    reading and `EachPlayer` the "both at once" reading. `whoseHand` generalized
    into `whoseZone(Player, Zone)`, so the phrase table is one function. All 16
    call sites migrated with **zero card-text changes** (`mage gen` rewrote 0
    comments), and `ctx.It` is now bound uniformly whenever exactly one card is
    purged rather than only by the hand node.
  - **`PurgeArchives` was deliberately left out of this pass.** It is
    `PurgeCard{Zone: Archives}` in every respect except its count: it purges
    "any number", which `Amount int` cannot say. Fold it when the `Count`
    strategy (`Exactly`/`UpTo`/`AnyNumber`/`All`) replaces `Selection`+`Amount`
    — adding an `AnyNumber bool` next to `Amount` now would be the boolean-flag
    smell the `Count` item exists to remove. One consumer
    (`destructive_analysis`), so there is no pressure.
  - **`PurgeCreature` survives** — it is Target-based and sources from play, and
    it carries the leaves-play-mid-resolution fallback that is the next item.
  - **`PurgeArchivedCardThen` survives** — its purge is a _cost gate_, not a
    selection.
- **A creature that leaves play mid-resolution is followed into its visible
  zone.** The rule is now the shared `currentZone(ctx, id)` seam in
  `effect_cross_zone.go`: it reports where a selected card can still be reached
  (in play under its controller, or its owner's discard pile) and `PurgeCreature`
  reads it instead of carrying its own fallback. Pinned by
  `TestPurgeCreatureFollowsIntoDiscard`. **Still to do:** the other mechanisms
  that move a creature which might die — archive, discard, return to hand — still
  fizzle instead of following. Route each through `currentZone`, and write the
  rule into the rulebook since it is player-facing.
- **The `Shuffle*` family is folded.** `ShuffleFromDiscard`,
  `ShuffleNamedFromDiscardIntoDeck`, and `ShuffleChosenCreaturesFromZones` are one
  `ShuffleIntoDeck{From []Zone, Selection, AnyNumber, Count}` over
  `crossZoneMover`. In-play is now two exported **source-only** zones like
  `Purged`: `InPlay` (battleline + artifacts, for a future card that shuffles
  upgrades out of play) and `Battleline` (the creature row alone, what Song of
  Spring names). They are separate because they are separate printed nouns —
  `Zone.noun()` must never guess which one a card meant. `ShuffleIntoDeck` also
  takes an explicit `Player` and tallies `Produced.Moved` per **owner**, so a
  controlled enemy card shuffled home credits its owner's deck. Two rewords fell
  out and are recorded in `docs/card-wording-rules.md`: a card name takes no
  article ("shuffle Subtle Chain", not "a Subtle Chain"), and Song of Spring
  drops "friendly" because "your hand, discard pile, or battleline" already
  scopes it.
- **Finish making upgrades "in play" everywhere.** Decided: _in play_ means every
  card someone controls in play — creatures, artifacts, **and the upgrades on
  either**. A card placed _under_ another card is **not** in play; it is in the
  out-of-play zone that is "under" its host. So an effect that says "cards in
  play" without naming a type reaches upgrades, while one that says Creature or
  Artifact narrows afterwards. **Done:** the one authoritative enumeration is
  `resolverCardsInPlay(ctx, p)` in `effect_cross_zone.go` (each host's upgrades
  listed ahead of the host), and `CardsInPlay`, `HousesInPlay`, untyped
  `HousesAmong`, and `ActiveHouseMatchesNoCardsInPlay` all read it — the last of
  which was hand-rolling the traversal and missing the upgrades on an artifact.
  Pinned by `TestCardsInPlayCountsUpgrades`. **Still row-only, and each needs the
  same treatment:**
  - `Game.inPlay` (`game_read.go:580`) and `resolverInPlay`
    (`effect_deck.go:681`) — the two "is this card in play?" predicates. They
    must change **together with** the selection paths below: if a Target starts
    returning an upgrade and the predicate still says it is not in play, every
    removal that re-checks reachability will silently fizzle.
  - `allCardsInPlay` and `cardsInPlayOf` (`target_select.go:496`, `:503`), which
    back `TargetEachCardInPlay` and `TargetEachFriendlyCardInPlay`. The
    creature-or-artifact target kinds name their types and must **not** change.
  - `crossZoneMover.zoneCards` for `InPlay` (`effect_cross_zone.go`). Moving an
    attached upgrade needs the mover to detach it from its host first. This is now
    the **gating item for the Timequake fold below** — it used to be blocked behind
    a leave-play sweep that has since been cancelled.
  - `Game.allInPlay` (`game_read.go:1054`) has ~14 callers in ability scanning,
    phase processing, and invariants. Decide per caller whether an upgrade
    belongs in that scan; do **not** change it wholesale.
- **Fold `ShuffleFriendlyCardsIntoDeck` (Timequake) into `ShuffleIntoDeck`.**
  `ShuffleIntoDeck` now carries the per-owner `Produced.Moved` tally that
  Timequake's `Draw{Per: CardsShuffledIntoDeck}` reads, so the shape would be
  `ShuffleIntoDeck{Player: Controller, From: []Zone{InPlay}, Selection: Each{}}`.
  The two nodes disagree about upgrades: the shared mover's in-play move runs
  `putIntoDeckShuffled` → `leavePlayTeardown` → `discardUpgrades`, so a shuffled
  creature's upgrades are **discarded**, while `shuffleFriendlyInPlayIntoDeck`
  detaches each upgrade and shuffles it in on its own.
  **The blocker has changed.** It used to read "blocked until the leave-play sweep
  replaces the proactive teardown"; that sweep is cancelled, so the eager
  `discardUpgrades` is permanent and the disagreement will not resolve itself. The
  remaining route is to make the mover reach upgrades as cards in their own right:
  `crossZoneMover.zoneCards(InPlay)` must list them and the mover must detach an
  upgrade from its host before moving it, so `ShuffleIntoDeck` moves creature and
  upgrade as two members of one `simultaneously` batch and never reaches
  `discardUpgrades` for them at all. `SwapDeckAndDiscard` stays separate — a
  distinct mechanism.
- **`discardDestroyed` / `purgeFromPlay`** (`game_leaves_play.go:24`, `:35`)
  repeat the same gigantic-halves loop. **INVESTIGATE**: destruction participates
  in the Destroyed window and purge can be invoked during it — establish whether
  the timing entanglement is real before folding, and report.

### Conditions and counts: atoms, not wrappers

The shared decision: a threshold is a **`Count` plus a comparison**, never a
bespoke condition type wrapping `>=`. Basic arithmetic at a card call site is
acceptable; flat pointerless state (ADR 0005) is not negotiable.

- **`Overwhelmed` and `ControlsMoreCreatures` stay two conditions. Decided: do
  not merge them.** Overwhelmed is a pure count of every creature on each side;
  `ControlsMoreCreatures` compares a trait on each side (Pismire compares Mutant
  counts), so the two ask different questions and a merged node would carry a
  `Player` × `Trait` combination no card uses. `ControlsMoreCreatures` has instead
  been decomposed onto the shared atoms: its threshold is now
  `CountIs{Count: ExcessCreatures{…}, Is: AtLeast, Amount: 1}`, its counted noun
  comes from `CardFilter.noun()`, and its board-wide third-person wording is its
  own `symmetricCondText` method rather than a type switch in `text.go`.
- **Subject is a field, not a name prefix. Split landed; finish the collapse.**
  The two concepts are now separate and honest: `ItNoun` (`itnoun.go`) is the
  wording choice — which noun the text prints in place of "it" — and carries the
  field name `Noun`; `Subject` (`subject.go`) is the real referent, with `It` (the
  card in context) and `This` (the card the ability is printed on), resolved by
  `Subject.card(ctx)` and rendered by `Subject.name()`. The investigation found
  only these **two** referents: `It*` conditions read `ctx.It`, while `This*` and
  `Source*` both read `ctx.Source` — there is no `ctx.This`, so `This` and the
  source are one referent.
  `ItHasAember` and `ThisHasAember` are merged into `HasAember{Subject}`.
  **The other nine are deliberately left alone, and here is why.** Checking each
  for a counterpart on the other referent found **none**: there is no
  `ItIsReady`, no `ThisIsStunned`, no `SourceIsNamed`. `HasAember` was the only
  question the card pool actually asks of both referents, so the remaining nine
  have exactly one referent each and their prefix is not hiding an axis — it is
  accurate. Adding a `Subject` field to each now would be speculative
  generality: the unused branch could not be reached by any card, only by a
  unit test written to keep the coverage gate at 100%. Their `CondText` also
  varies more than the referent does (`SourceIsFighting` says "if fighting" so
  the constant-ability renderer can reframe it; `ItIsStunned` says "that
  creature", not "it"), so the merged text would not fall out of `Subject.name()`
  anyway. **When a second card does ask one of these of the other referent, add
  `Subject` to that one condition then** \u2014 the enum and its `card`/`name` helpers
  are already in place, so it is a one-node change.

### Effect composition

- **Collapse a repeated subject to a pronoun in the renderer.** Folding
  `ReadyIfFirstUse` into `Conditional{SourceFirstUseThisTurn, Ready{This}}` cost
  Rocket Boots its pronoun: the composed text reads "If this is the first time
  this creature has been used this turn, ready **this creature**", where the
  bespoke node hard-coded "ready **it**". The composition is right and stays; what
  is missing is a renderer pass that replaces a subject already named earlier in
  the same sentence with "it". Build it alongside the conjoined-`It*`-condition
  collapse below — both are the same job: the renderer, not the node, decides how
  a repeated subject reads.
- **`PlayOrUse` and friends are decomposable.** The instinct that a single combined
  prompt is the mechanic does not mean the _node_ must be bespoke: the axis is a
  grant of `Play || Use`. Reshape `PlayOrUse` (`effect_play_or_use.go`) and
  `DiscardOpponentArchivesOrDeckTop` (`effect_discard_opponent_source.go`) so the
  permitted-verb set is a value, not a node name. `PutItIntoHand`
  (now in `effect_put_into_zone.go`) is a silent runtime zone probe — it folds
  into the follow-the-creature rule above rather than into a visible choice.
- **Tide is next-set work** — `TideIsHigh` having no consumer is expected; leave it.

### Engine folds

- **One lasting-registry scan. INVESTIGATE**: three scans exist — reactions
  (`game_lasting.go:320`), replacements (`:520`), and a bespoke Before-Fight scan
  (`:373`). Approved in principle; bring back **worked examples of each scan and
  what the tension is** (Before-Fight matches by subject rather than owning actor;
  replacements return one record) so the unified matcher does not become a
  branchy mini-language.
- **Start-of-turn and end-of-turn order abilities the same way.** End-of-turn
  gathers one ordered window (`game_phase.go:194`); start-of-turn fires per-card
  then does a separate cross-player loop (`:71`, `:91`). **Decision: the ordering
  logic is the same for both; the only difference is the window (start or end).**
  Governed by ADR 0013 — update it if the fold changes what it describes.
- **One in-play traversal, and prefer the mutation-safe one.**
  `game_read.go:583`, `:601`, `:1054` each know the two physical zones, and 15+
  scans repeat it. **INVESTIGATE the tradeoff** the human raised: can we just use
  the single mutation-safe (snapshot) implementation everywhere and drop the
  callback variant? Report the allocation cost on the hot paths (MCTS load — use
  `mage profile`) before deciding.
- **Upgrade teardown shares its prefix.** Four paths
  (`game_leaves_play.go:150`, `:163`, `:180`, `:278`) repeat detach / release
  control / clear counters / reset. Fold the prefix; wrappers keep only
  destination and logging. **INVESTIGATE**: the blocker raised was event
  attribution — produce a concrete example of what attribution is used for, then
  show the fold preserving it (expected to need only a minor extension).
- **`TargetKind` has three interpretations. INVESTIGATE**: `isChosen`
  (`target_select.go:165`), a large text switch (`target.go:534`), and the
  base-set switch. Target is deliberately broad, so find the clean lines to cut
  across rather than forcing one opaque descriptor table, and propose them.
- **Stat aggregation folds.** Power (`game_read.go:103`), armor (`:154`), and the
  repeated printed/upgrade/temporary/constant/blanking logic (`:312`, `:327`,
  `:345`). Prefer narrow helpers (`upgradeStatBonus`, `constantStatBonus`) over
  one generic stat switch; **double-check** that Power, armor, Assault, Hazardous
  and Splash really share the pipeline before merging.
- **Constant-ability traversal folds — except blanking. INVESTIGATE**: four scans
  (`game_read.go:189`, `:261`, `:282`, `:300`) share scan-allInPlay / check-active
  / check-affects, but `constantBlanksText` deliberately avoids recursive source
  blanking. Find a path that keeps that distinction explicit rather than hiding it
  behind a mode flag, and report.
- **Boolean clusters become small comparable option structs.** `SetStatOverride`
  is done (it takes two `StatMask` values now). Still open:
  `PutIntoBattlelineAsCreature(…, right, temporary bool)` (`resolver.go:329`),
  `GrantTextBox(…, remainderOfTurn bool)` (`:394` — a `Duration` would say more
  than a bool), and `mitigateDamage(…, ignoreArmor bool, …)`
  (`game_combat.go:505`). Safe as long as the structs do **not** enter
  `GameState`. While doing these, **sweep for the other call sites** that would
  read better as comparable records.
- **`ZoneResolver` and `CreatureResolver` role placement.** `ZoneResolver`
  (`resolver.go:455`) has grown into a grab-bag (ordinary moves, play-from-zone,
  bonus icons, under-cards, shuffles, archive/purge); `CreatureResolver` holds
  board-wide duration-scoped mutations (`:291`, `:295`, `:359`) despite being
  documented as per-card state. Re-home against ADR 0008 — but this is a broad
  port change, so land it after the zone-movement fold above, which will change
  the method list anyway.

### Rulebook, tests, and coverage

- **Close the effect catalog (ADR 0018).** `ruleterm_test.go:8` exempts effects
  from the completeness test, so a player-facing mechanic can ship undescribed.
  Build the closed effect catalog and enforce it, then write terms for the gaps it
  surfaces.
- **Raise the minimum bar for a card test. INVESTIGATE**: today
  `internal/cards/testfiles_test.go:34` only checks that a matching `_test.go`
  exists, so a test can pass while asserting nothing about the card's ability.
  Engine tests carry most of the load (cards hold no code), so the goal is modest:
  work out how to express "a new card's test must at minimum cover these
  aspects" and enforce what can be enforced. Propose the rule before building it.

### Web client

- **Prompt source by identity, not name. BLOCKED on an engine change.**
  `internal/web/view_controls.go` matches `Def(id).Name == def.Name`, so two
  copies in play can preview the wrong card's live house/Maverick state. The web
  client **cannot** fix this on its own: the engine's `Chooser` prompt API hands
  the client only the source card's _name_, with no `LocalID`, so there is no
  identity to match on. The fix is to thread a source `LocalID` through every
  `Chooser` method (and the `replayChooser`, `sessionChooser`, and MCTS chooser
  that implement them), then have `promptSourceHouse` take the id directly. That
  is an engine port change spanning `resolver.go`, every prompt call site, and
  three implementations — do it as its own piece of work, not as a web tweak.
- **Engine regressions must survive a bad snapshot.** `game_persist.go:92`
  recovers a replay panic and deletes the snapshot, so a genuine regression and an
  old incompatible save look identical and the reproduction is destroyed.
  Incompatible saves may still be tossed. **INVESTIGATE** the mechanism the human
  asked for: have the local dev server capture a failing client snapshot to disk
  the way soak and fuzz findings are captured, so `mage check` (or a sibling
  target) surfaces it for an agent to investigate, cover, fix, and clean up.
  Mirror the existing corpus workflow.

### Tooling and hygiene

- **HUMAN: untrack `cardcheck` and `err.txt`.** Both files are **deleted from the
  working tree** and `/cardcheck` is now in `.gitignore`, so the next commit drops
  them from the index. No agent action remains — agents do not run git write
  commands (see AGENTS.md), so this is recorded here only so it is not forgotten.
  Nothing in the tree was found writing `err.txt`; it looks like a one-off stray
  redirect.
- **Replace `qmark`.** The quickmark markdown linter is a brew-installed Rust
  binary that `mage check` silently skips when absent (`magefiles/lint.go:25`),
  and it has a known emphasis-scanner bug (see repo memory). **INVESTIGATE** a
  better markdown linter — ideally Go, ideally something `golangci-lint` or a
  `go run`-pinned tool can provide — and migrate, so the gate stops being able to
  pass without checking Markdown.
- **`docs/deck-generation.md` status block is fixed** — it now says the pipeline
  is built and marks § 6 (scoring and band-targeting) as the only design-only
  part. What remains: `internal/scoring` does not exist at all, so if § 6 is not
  going to happen soon, consider moving it to `docs/roadmap.md` rather than
  leaving a whole section of a "how it works" doc describing nothing.
- **`docs/roadmap.md` is labelled** as forward-looking notes rather than
  current-state documentation, so a sweep should leave it alone. Nothing to do.
- **`internal/session` vs web replay. Investigated twice; BLOCKED on an engine
  gap, and the gap is now measured.** `internal/session` is the newer,
  ADR-0040-blessed driver: it is built on the
  `engine.Stepper`/`Command`/`Request`/`View` seam in `suspend.go`, owns
  `{version, seed, sets, []Command}` plus the undo cursor, and is fully tested —
  but **nothing imports it**. `internal/web/replay.go` is a second, live
  implementation of the same job with its own ad-hoc `input` enum. So the
  migration was started and never completed, not the other way round.
  **The decision stands: migrate the web client onto `internal/session` and delete
  `internal/web/replay.go`.** What the second investigation established is _why it
  cannot start yet_:
  - **The engine's command vocabulary only implements half of ADR 0039.** The ADR
    says a command is "a root action _or_ one answer to a choice", but
    `engine.CommandKind` has only the answer half — `CommandPickCard`,
    `CommandDecline`, `CommandOption`, `CommandPosition`, `CommandReaction`. Of
    `replay.go`'s 27 `inputKind` values, **4** map onto a `Command`. The 11 real
    root actions (`ChooseHouse`, `PlayCreature`/`Artifact`/`Action`/`Upgrade`,
    `DiscardFromHand`, `Reap`, `Unstun`, `UseAction`, `Fight`, end-turn) and the 13
    manual/debug roots have no representation at all.
  - **`session.Action` is a fixed closure, not a live driver.** `Stepper` takes one
    `func(*Game)` that must already encode the whole turn loop, and `Session`
    exposes no way to say "run this root action next". A web client needs each
    click to be the next command. So root actions must become a suspension point:
    the turn loop yields a `RequestAction`, the client answers with a root
    `Command`.
  - **Legality enumeration exists but is not collected.** `CanPlay`, `CanDiscard`,
    and `CanUse` already gate each root action (`internal/sim/sim.go:159-169` walks
    them), so `Request.LegalCommands()` for a root request is assemblable from
    parts that exist — no new rules work, just a gatherer.
  - **`Orderer` is NOT a blocker.** `suspendChooser` does not implement it, so an
    order prompt degrades to repeated `ChooseCreature` — which is the engine's
    documented fallback and is fully expressible as a run of `CommandPickCard`s.
    The web's single drag-to-order widget is a UI affordance over that run, not a
    missing command kind.
  - **Manual/debug mode is the open question.** 13 `Manual*` calls edit state with
    no rule checks and sit entirely outside `Request`/`Command`. Either they become
    command kinds (and a saved match can contain force-edits) or manual mode is
    declared a web-only escape hatch that voids the Record.
  - **Undo granularity differs.** Web undo works at root-action boundaries
    (`rootMarks`); `Session.Undo(n)` works at raw command index. Once root actions
    are commands the web still needs its own boundary bookkeeping on top.
  - **Also update ADR 0039's and ADR 0040's stale opening lines** ("the engine code
    that realizes it does not exist yet") when this lands.
  - **DECIDED (human, this session): build `RequestAction` in the engine.** One
    command vocabulary, not two — a second root-command type beside
    `Request`/`Command` would recreate the two-sources-of-truth problem ADR 0039
    exists to kill, and the ADR already says a command is "a root action _or_ an
    answer to a choice". The engine grows a canonical turn loop that suspends with
    a `RequestAction` carrying the legal set, and resumes on a root `Command`;
    `internal/web`'s hand-rolled turn driver then disappears into the engine.
  - **DECIDED (human, this session): manual/debug roots become command kinds.** A
    playtester who force-edits the board and then hits a bug has produced exactly
    the reproduction worth keeping, so the Record must carry the force-edits rather
    than be voided by them. It also opens manual mode as a source for the style
    page's log gallery (ADR 0046). **Constraint: the simulator must never run with
    manual mode** — `internal/sim` drives legal play only, so the manual command
    kinds must be unreachable from the sim's driver and a test should assert it.
    Decide whether that is a build-tag split, a flag on the driver, or a
    `LegalCommands()` that simply never offers them (preferred: the last, since it
    keeps one vocabulary and makes the restriction a property of legality rather
    than of who is asking).
- **An engine-implementation doc, mirroring card-implementation.md.
  INVESTIGATE**: the `1 keys` bug happened because `countNoun` existed and was not
  found. Decide whether the fix is a new engine-side capability catalog (the
  helpers, log entry types, text helpers, and invariants an engine change should
  reach for) or a section in an existing doc, and propose it.

## Text & log style gallery (ADR 0046)

Decided in a grilling session: give `/style` a section showing an example of
every kind of rendered card text and game-log line, sampled from real games and
anchored to an enumerable engine catalog. Build in dependency order — the engine
catalog is the keystone the rest consumes, so it lands first. Stages are each
independently green.

- **Stages 1 and 2 are done.** Stage 1: `engine.LogEntrySamples()` is the exported
  catalog (one value per `LogEntry` variant), `TestLogEntrySamplesTotality` fails
  the build if a `Text(Namer)` type has no sample, and `TestLogEntryText` consumes
  the catalog. Stage 2: `sim.Play(script)` hands back a played game; `sampleLog`
  ([internal/web/style_log_sample.go](../internal/web/style_log_sample.go)) plays
  seeded scripts, retains games that add a new kind, and greedily set-covers their
  bubbles, reporting the kinds no game produced (24 of 130 in a 300-game run — the
  manual-edit entries, concede, restore, and rare unimplemented-card effects) for
  the synthetic fallback. Tested in `style_log_sample_test.go`.
- **Stage 3 — the two `/style` sections (web, ungated). Game log is done; card
  text is partly done.**
  - _Game log_ (**done**): `logSection`
    ([internal/web/style_log_section.go](../internal/web/style_log_section.go))
    lazily samples games on a button click (wasm-safe), draws the set-cover hero
    gallery with production `logBlockView`, drills into a clicked bubble's full
    game log with the bubble highlighted, and renders every unobserved kind from
    its `LogEntrySamples()` instance badged "synthetic — not observed in N games".
    Tested in `style_log_section_test.go`.
  - _Card text_ (**partly done**): `cardTextSection`
    ([internal/web/style_cardtext_section.go](../internal/web/style_cardtext_section.go))
    shows one specimen per **trigger kind** (ranging `engine.Triggers()`, printed
    only), one per **continuous ability** field (constant/granted, key cost,
    Æmber bonus, Æmber-cannot-be-stolen, spendable Æmber, forge Æmber, upgrade),
    one per **target shape** (distinct `Target.Text()` phrase gathered from every
    ability effect by a reflection walk, `walkEffectTargets`, so it is data-driven
    and gapless), one per **duration** (ranging `engine.Durations()`, captioned by
    `Duration.String()`, a real card carrying that window found via the same
    reflection walk `walkEffectDurations`; gaps drawn where no loaded card uses a
    window), and a **Combiners** subsection of constructed minimal pairs
    (`combinerRows`) isolating the article (`a`/`an`) and single-vs-collective
    quantifier helpers. Real cards matched over `cards.All()`, gaps drawn not
    skipped, terms linked to `/rulebook`. Tested in `style_cardtext_section_test.go`.
    Duration needed a small engine surface — `engine.Durations()` + a canonical
    `Duration.String()` short label — because there is **no** single printed
    duration phrase: the same window renders differently per effect and flips
    prefix/suffix, so the caption names the window and the card shows the phrasing.
    **Still to add**: `Condition`/`Count` composition — the last axis with no
    card-side text hook yet. Expose an enumerable/renderable hook on the engine
    side, then add a subsection in the same "walk `def.Abilities[].Effect`,
    enumerate, randomMatch" shape.
  - Tests assert coverage (every catalogued kind sampled-or-flagged) and no-panic
    render — never markup or wording (ADR 0014).

## Handoff: state of the zone-movement sweep

Written for an agent who was not in the session that produced it. Everything
below is current as of the last `mage check`.

### What just landed

The zone-movement consolidation is at a clean stopping point. `simultaneously`
(`game_settle.go`) is now the one moment primitive: it holds the settling flag
**and** the deferred leave-play queue across a batch, so a multi-card effect moves
every card before any of them reacts. `Destination.move` is variadic and
`PutFromPlay.put` runs its whole selection inside one batch.

The part most likely to surprise you: a deferred `Leaves Play:` window resolves
**after** its card has left play, which the source-in-play guard in
`resolveTriggered` would normally skip. `triggeredAbility.fromLeave` is the
exemption, and it is the fourth one in that guard — the others are tactics,
`TriggersFromDiscard`, and duration reactions. This is written up as a Refinement
section in ADR 0030. If a `Leaves Play:` ability ever appears to do nothing, that
guard is the first place to look.

Two deliberate asymmetries, both documented at their seam, neither a bug:

- `fireScheduledOnLeave` (Turnkey's armed forced forge) still fires immediately
  rather than joining the queue. It is an armed effect, not a trigger, and
  deferring it would mean its flat-state entry has to outlive the card.
- `flushDeferredLeaves` orders its window by `g.State.ActivePlayer` (ADR 0013)
  while `settleDestroyed` beside it takes the resolving controller. The two
  genuinely differ when a batch takes cards from both players out at once.

### What was cancelled, and why it must not be re-proposed

`settleBoard` and the orphaned-upgrade sweep. The backlog said to strip
`discardUpgrades`/`discardUnder` out of `leavePlayTeardown` and let a sweep find
upgrades whose host had gone. Investigating it killed it: `settleDestroyed`
already loops to a fixpoint, a dangling upgrade is **already** an invariant
violation that `-tags assert` runs never trip, and the eager discard settles at
the same boundary a sweep would. It would have been cleanup for a state that
cannot occur, in the hottest path in the engine. The reasoning now lives on
`leavePlayTeardown`'s doc comment, which is the seam it governs.

That cancellation also took the `returnUpgradesToHand` fold with it — that item
depended on the sweep to make upgrades reachable as cards in play in their own
right, so the helper and its `ZoneResolver.ReturnUpgradesToHand` port method stay.

**It also moved the Timequake blocker.** The `ShuffleFriendlyCardsIntoDeck` fold
used to be "blocked until the sweep lands". The sweep is never landing, so the
route is now `crossZoneMover.zoneCards(InPlay)` listing upgrades plus a detach in
the mover. Both items above say so; do not re-derive it.

### Next up

1. `crossZoneMover.zoneCards` for `InPlay` — now the gating item, since the
   Timequake fold sits behind it.
2. The Timequake fold itself. **Check `git status` first**: a sibling agent was
   actively rewriting `effect_shuffle_*`, `effect_purge`, `effect_forge`, and
   `zone.go` throughout this session. Stay out until their work lands.
3. The remaining uniform-in-play groups, in the "Follow the creature" item above.

### Gate status

`mage check` is green on everything except engine coverage, which sits at 99.8%
from the sibling agent's in-flight files — `effect_shuffle_discard.go` (0.0% ×5),
`keycolor.go:35`, `log_render.go:92`, `assert_off.go:13`,
`effect_damage.go:112`. None of them belong to this work; verify with
`mage cover` before assuming a regression is yours. Worth knowing:
`go test -tags assert ./...` enables the invariant checks at turn boundaries and
is not part of `mage check` — it is the fastest way to prove a state-corruption
claim, and it is what retired the orphan sweep.

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
- **`resurgence` second pick says `another`.** Render `If that creature is a
Mutant, put **another** creature from your discard pile into your hand` — and
  make the implementation actually forbid re-picking, since the text now promises
  it.
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
- **Grouped targets are singular, so a card that matches twice acts once.**
  `emp_blast` renders `Stun each Mars creature and each Robot creature`. Adopt
  `each Mars or Robot creature`: stun is idempotent so it does not matter there,
  but the same renderer produces `deal 1 damage to each Mars creature and each
Robot creature`, where a Mars Robot reading as 2 damage is wrong. **Decision:
  the grouped-union form is correct and a card matching both halves is affected
  once.**
- **Event-scoped conditions say `this`, never `it`.** `aember_conduction_unit`
  renders `if it is the first time a creature has reaped this turn` — `it` is
  reserved for a card. Use `this` for an event referent and apply it consistently
  everywhere an event, not a card, is the subject.
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
    `ReturnNamedToHand`, `Search`, `ShuffleNamedFromDiscardIntoDeck`, and
    `ShuffleChosenCreaturesFromZones`. It is the fold's destination, not a rival.
  - **Start with purge** — the smallest complete verb (`PurgeCard` 9 cards,
    `PurgeFromHand` 7, `PurgeCreature` 21, `PurgeSource` 3, `PurgeArchives` 1,
    `PurgeArchivedCardThen` 1) — and land the mover collapse with it.
  - **These are pure subsets or thin wrappers and dissolve:**
    `ArchivePurgedCard`, `ArchiveDiscardedThisWay`, `DiscardTop`, `PurgeFromHand`,
    `PurgeArchives`, `PutRevealedCard`, `PutDiscardedIntoHand`,
    `PutDiscardedIntoPlay`, `PutFromHand`, `PutNamedIntoHand`, `PutItIntoHand`,
    `PutChosen` (vs `PutFromPlay`), `ShuffleNamedFromDiscardIntoDeck`,
    `ShuffleChosenCreaturesFromZones`.
  - **These carry a genuinely distinct rule and survive:**
    `ArchiveGrantingUpgrade` (`ctx.Upgrade`), `ArchiveFromPlay` (batch +
    `ctx.Produced.Archived`), `DiscardArchives` (active-player ordering +
    randomized opponent archives), `DiscardHand` (`playersInOrder`),
    `DiscardOpponentArchivesOrDeckTop` (two unlike hidden sources),
    `PurgeArchivedCardThen` (cost gate), `PurgeCreature` (play-or-discard
    fallback), `PurgeCard` (\u00c6mber bonus tally), `PutNextTacticIntoHand` (lasting
    replacement), the under-chain nodes (`game_under.go:117-123` \u2014 explicitly not
    zone movement), bare `Shuffle`, `ShuffleFriendlyCardsIntoDeck`.
  - **`DiscardUntil` becomes a general `Until{}` dig-loop primitive** that other
    `*Until` effects can compose over, rather than staying a bespoke node.
  - **A resolving tactic is still in hand** \u2014 it is simply not in hand to be
    discarded. That reclassifies `ArchiveSource`/`PurgeSource`, which were held
    back as "genuinely distinct" only because the resolving tactic appeared to be
    in no zone at all. **Open design question: does this want an explicit
    resolution zone?** Settle that before folding those two.
- **Purge is one mechanism** — `PurgeCreature`, `PurgeFromHand`, `PurgeArchives`,
  `PurgeCard`, `ArchivePurgedCard`, `NamedCardPurged` collapse onto it. Approved
  outright, but **held until the verb-fold investigation above reports**: purge is
  one of the four verbs that fold, so doing it standalone would set a
  source/destination shape the other three then have to match. Build it as part of
  that consolidation, not before it.
- **A creature that leaves play mid-resolution is followed into its visible
  zone.** `PurgeCreature` (`effect_purge.go:221`) has a fallback: if the target
  already died, find it in discard. **That is a general rule, not a purge
  feature** — it belongs to _every_ mechanism that moves a creature which might
  die (archive, discard, return to hand, …). If an ability is mid-resolution and
  its subject leaves play into a **visible** zone, the ability continues on it
  there. Lift the fallback into the shared movement mechanism, and write the rule
  into the rulebook/engine guidance since it is player-facing.
- **`ShuffleFromDiscard` takes a `Selection` strategy**, absorbing
  `ShuffleNamedFromDiscardIntoDeck` (one consumer, `chain_gang`).
- **Fold the `Shuffle*` family** (`ShuffleChosenCreaturesFromZones`,
  `ShuffleFriendlyCardsIntoDeck`, `ShuffleFromDiscard`) into one shape with zone
  and selection axes. **`SwapDeckAndDiscard` is a genuinely distinct mechanism and
  stays separate.**
- **Ownership is looked up from the card**, not passed around. The objection that
  archives can hold abducted cards does not block the fold: the mover takes an
  explicit destination (this card, to this player's archives) and derives owner
  from the card.
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
- **`internal/session` vs web replay. Investigated; finish the migration.**
  `internal/session` is the newer, ADR-0040-blessed driver: it is built on the
  `engine.Stepper`/`Command`/`Request`/`View` seam in `suspend.go`, owns
  `{version, seed, sets, []Command}` plus the undo cursor, and is fully tested —
  but **nothing imports it**. `internal/web/replay.go` is a second, live
  implementation of the same job with its own ad-hoc `input` enum. So the
  migration was started and never completed, not the other way round.
  **The decision: migrate the web client onto `internal/session` and delete
  `internal/web/replay.go`.** The wrinkle to plan for is that ADR 0040 says the
  engine code realizing the pure step function "does not exist yet" — the engine
  still pulls choices through `Chooser` — so establish how much of `replay.go`'s
  chooser/prompt replay `Stepper` already absorbs before starting.
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
    and gapless), and a **Combiners** subsection of constructed minimal pairs
    (`combinerRows`) isolating the article (`a`/`an`) and single-vs-collective
    quantifier helpers. Real cards matched over `cards.All()`, gaps drawn not
    skipped, terms linked to `/rulebook`. Tested in `style_cardtext_section_test.go`.
    **Still to add**: the axes with no card-side text hook yet — `Condition`/`Count`
    composition and **durations** (`engine.Duration` renders no exported text; the
    effect that uses it does). Expose an enumerable/renderable hook for these on
    the engine side, then add a subsection in the same "walk `def.Abilities[].Effect`,
    enumerate, randomMatch" shape.
  - Tests assert coverage (every catalogued kind sampled-or-flagged) and no-panic
    render — never markup or wording (ADR 0014).

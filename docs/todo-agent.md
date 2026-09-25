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

## Closed effect catalog + a `RulesBearing` marker (ADR 0018)

**Decided** (the human said yes to the marker). Two halves, in this order: the
catalog first because the marker has nothing to iterate over without it.

1. **Build the closed catalog of `Effect` implementations.** `ruleterm_test.go`
   currently exempts effects from `TestClosedCatalogsAreComplete` because there is
   no enumeration of them — there are ~131 implementations and nothing lists them.
   ADR 0018 makes the rulebook complete _by construction_, so this hole means an
   effect node can ship a new rule with no rulebook term and the build stays green.
   The catalog is mechanical to build and also unblocks totality tests for the
   Visitor passes (the AI/MCTS scoring walk and any future static analysis), which
   today silently fall through their `default:` for a node nobody added.

2. **Gate rulebook-term enforcement behind a per-node `RulesBearing` marker.** Do
   _not_ force every node in the catalog to carry a term: `Sequence`,
   `ChooseHouseThen`, `Conditional` and friends are plumbing — they express no
   rule a player needs described, and demanding a term for them would produce
   filler rulebook entries, which is worse than the hole. The marker names the
   nodes that _do_ carry a rule, and only those are checked for a term.

## `ForgeKey` bakes "purge self" into a shared node

**Decided: strip the purge out of `ForgeKey` and push self-removal up to the two
cards that actually print it.** Found while diagnosing a sim invariant; not fixed
in-session because `internal/engine` was contended and the blast radius is 17
cards' generated text.

What is wrong today, in `internal/engine/effect_forge.go`:

- `Text()` unconditionally appends `" -> purge " + SelfName`.
- `Resolve()` ends with `if forged { PurgeSource{}.Resolve(ctx) }`.

That is a card's behavior welded into a shared mechanic. The evidence
(`mage tool:lookup` against the provenance files):

- **Key Charge** (CotA #325) prints "Play: Lose 1 Aember. If you do, you may forge
  a key at current cost." — **no self-removal at all**. Vex currently renders
  `… -> purge Key Charge.`
- **Imperial Forge** (WC #222) — no self-removal either.
- **Epic Quest** (CotA #231) and **[REDACTED]** (AoA #139) say **sacrifice**,
  which [card-wording-rules.md](card-wording-rules.md) renders as **`Destroy
<self>`**, not purge.

Nothing in [keyforge-divergences.md](keyforge-divergences.md) licenses the purge,
so it is a bug, not a recorded divergence.

The work: remove the purge from `ForgeKey.Text()` and `ForgeKey.Resolve()`, add an
explicit `Destroy <self>` to Epic Quest and [REDACTED], then `mage gen` and fix the
card tests it moves. The 17 users of `ForgeKey{` are `imperial_forge`,
`data_forge`, `forging_an_alliance`, `obsidian_forge`, `the_colosseum`, `triumph`,
`might_makes_right`, `nightforge`, `redacted`, `epic_quest`, `chota_hazri`,
`key_charge`, `key_abduction`, `key_of_darkness`, `turnkey`, `desire`, `keyfrog` —
check each against its printed text while you are there.

Trap for whoever picks this up: `redacted.go`'s generated doc comment reads
"… -> purge [REDACTED]" although the card definition contains no purge. The purge
comes from inside `ForgeKey.Text()`. Also, `[REDACTED]` is a real card name, not a
log redaction.

## gocognit gate: one exclusion left, and the next ratchet step

The engine and web passes are done: the `gocognit` gate's global `min-complexity`
is 30, with every genuinely nested engine/deckgen/cardtest/web seam split (fight,
NewCard, hasKeyword, the text.go renderers, Target.Text, Target.filter,
triggeredBy, RenderAbility, Harness.location, allowedHouses, `renderMarkdown`,
`computeFlashes`, `installKeyShortcuts`, `installTips`) and 100% coverage held.
One exclusion remains in the `tools.golangci` `exclusions` `text:` list in
`mklv.config.json`:

- `effectGlyphs` (internal/web/icon.go, 90) — glyph dispatch; wants its own
  glyph-family grilling session (see docs/todo.md "Split out glyphs more in
  icon.go"). Split along glyph families, then delete its exclusion.

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

**Do not trust an item's "these cards are the same shape" grouping without
checking each card's printed text with `mage tool:lookup` first.** A retired item
here claimed three cards shared one fix and none of them did: one was already a
single each-player effect and needed no fold, one spells both halves out on the
printed card deliberately and must not be folded, and only the third was the
shape the item described. The same trap retired the card-text fan-in sweep: an
affix search reported a dozen candidates and all but three were either
already-folded output or an asymmetric rule that must repeat (Savage Clash spares
the most powerful enemy and the **least** powerful friendly).

### Game log

Every change to a rendered log string requires bumping `snapshotVersion`
([internal/web/game.go](../internal/web/game.go), currently 21), because persisted
logs store rendered prose. **Bumping it costs nothing** — a stale snapshot is
flushed and the next deal is fresh — so change a log string when it is wrong and
bump. Do not batch unrelated log work to keep the number low.

The log's voice rules — lowercase mechanics, source card first, no left-to-right
backtracking, no parentheticals — are written into `internal/engine/log.go`'s file
doc comment, so a new entry is written to them. Keep that comment as the
authority; the items below are what is still unapplied.

Items:

- **"Every state change is logged": half landed, half needs a decision.**
  `TestEveryStateFieldDeclaresItsNarration` ([internal/engine/state_test.go](../internal/engine/state_test.go))
  now classifies all 72 `GameState` fields through `fieldNarration` — narrated
  directly, narrated by the card that caused it, or never (only `PRNG`) — and a
  new field fails the build until it is classified. That is the ratchet; it does
  not prove the log is complete.
  **INVESTIGATE the played-game half.** A version that plays whole games and
  fails when a step changes the state without appending an entry was written and
  **deleted**, because it could not fail: at the granularity the sim driver
  exposes (`StartTurn` / `doAction` / `EndPlayPhase`) each step already emits many
  entries, so suppressing `CardsDrawn` — or every entry in a whole phase — still
  passed. The check only has teeth per mutation, which needs a hook the engine
  does not have. Decide where that hook goes and what it may cost (it must not
  burden MCTS rollouts; `mage profile` measures) before writing it again. Do not
  re-add a step-level version: a test that cannot fail is worse than none.

### Zone movement: one mechanism per KeyForge verb

The sweep is **done** and its entries are deleted. The decision it settled, for
context when reading the code: each KeyForge verb — archive, discard, purge, put
— is one mechanism with a source/destination axis, not a family of near-duplicate
nodes. `PutCard` / `PurgeCard` / `DiscardCard` share one shape (`Zones` +
`Selection` + `Quantity` over `crossZoneMover`), every relocation out of play goes
through `leavePlayTeardown` / `fileFromPlay` / `leavePlayInto`, and
`Destination.moveFrom` is the one source/destination matrix (ADR 0031). The nodes
that did **not** fold each say why in their own doc comment; do not re-litigate
them from their names.

One piece is deliberately left open:

- **Widen `Game.allInPlay` per caller.** _In play_ means every card someone
  controls in play — creatures, artifacts, **and the upgrades on either**. A card
  placed _under_ another card is **not** in play. The reads, predicates, and
  selection paths are already widened onto `resolverCardsInPlay` (pinned by
  `TestCardsInPlayCountsUpgrades`, `TestInPlayCountsUpgrades`, and
  `TestEachCardInPlayReachesUpgrades`), and the creature-or-artifact kinds were
  split onto row-only `creaturesAndArtifacts` / `creaturesAndArtifactsOf` because
  they name their types. What is left is `Game.allInPlay` (`game_read.go`), still
  row-only with **69 callers across 18 files** in ability scanning, phase
  processing, combat restrictions, key-cost calculation, and invariants. Decide
  per caller whether an upgrade belongs in that scan; do **not** change it
  wholesale. At that size, expect to split the callers into two named helpers
  rather than to review all 69 in one sitting.

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
  `Subject` to that one condition then** — the enum and its `card`/`name` helpers
  are already in place, so it is a one-node change.

### Engine folds

- **No implemented card puts an upgrade into play.** `putIntoPlay` has only
  `Creature` and `Artifact` arms, so the two `AfterUpgradeEnters` cards (Armory
  Officer Nel, Commander Dhrxgar) can never fire from a put-into-play. Fine now;
  a trap for the first card that needs it.
- **One in-play traversal. Mostly landed; only `allInPlay` is left.** The "15+
  scans repeat it" claim in the old wording was wrong — a sweep for
  `State.Artifacts[` found the two-zone knowledge in only four places, and three
  are now folded: `choosableHouses` goes through `allInPlay`, and
  `placeUnderController` goes through `removeFromPlayRows` (the shared remove
  under both a real exit and a change of control — its doc comment already
  claimed control used it, which was untrue until the fold). The tradeoff the
  human raised (drop the callback variant for the mutation-safe snapshot one) is
  **moot**: there is no callback variant; every traversal already returns a fresh
  slice (`allInPlay`, `battlelineCopy`, `resolverCardsInPlay`). What remains is
  `inPlay`, which keeps its own two-zone loop on purpose — it is a `contains`
  predicate over both players and allocates nothing, so routing it through
  `allInPlay` would add a slice to a hot path. Measure with `mage profile` before
  changing that.
- **`fireLastingBeforeFight` cannot honour `Once`.** It resolves reactions
  directly rather than through a trigger window, and the window
  (`game_abilities.go`) is the only place that consumes `Once` and removes the
  record. So a one-shot Before-Fight record would fire on every fight. No card
  installs one today. **Decided:** when one does, route this scan through the
  window path so it inherits both `Once` and the ordering prompt — do **not**
  re-derive the filters inside the Before-Fight scan, which would give the
  registry a second, divergent matcher. The scan's doc comment records this.
- **Boolean clusters become small comparable option structs.** `SetStatOverride`
  is done (it takes two `StatMask` values now), and so are the two bools that
  encoded a timing window: `PutIntoBattlelineAsCreature` and `GrantTextBox` both
  take a `Duration` now (`RemainderOfPlayerTurn` or `UntilCardLeavesPlay`), which
  reused the existing vocabulary instead of inventing an option struct —
  `TurnIntoCreature` had a `Duration` field all along and was flattening it to a
  bool at the port boundary. Still open: `mitigateDamage(…, ignoreArmor bool, …)`
  (`game_combat.go`). Note it is **one** bool, not a cluster, and both call sites
  pass a named `ignoreArmor` variable rather than a literal, so it reads fine
  today; do it only as part of a wider damage-options change, not on its own.
  While doing these, **sweep for the other call sites** that would read better as
  comparable records.
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
  the way soak and fuzz findings are captured, so `mage ci:check` (or a sibling
  target) surfaces it for an agent to investigate, cover, fix, and clean up.
  Mirror the existing corpus workflow.

### Tooling and hygiene

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
  - **The root-action Command vocabulary now exists — half landed.** ADR 0039 says
    a command is "a root action _or_ one answer to a choice"; `engine.CommandKind`
    used to have only the answer half (`CommandPickCard`, `CommandDecline`,
    `CommandOption`, `CommandPosition`, `CommandReaction`). The **11 legal root
    actions** now have a Command representation and a dispatcher:
    `CommandChooseHouse`, `CommandPlayCreature`/`Artifact`/`Tactic`/`Upgrade`,
    `CommandDiscardFromHand`, `CommandReap`, `CommandUnstun`, `CommandUseAction`,
    `CommandFight`, `CommandEndTurn` (suspend.go), performed by
    `Game.ApplyAction(Command)` (command_action.go), which mirrors `replay.go`'s
    `dispatchRoot` so that hand-rolled driver can fold into it. The Command struct
    grew the root fields it needs (`House`, `Hand`, `Left`, `Card2`). The **13
    manual/debug roots** now also have Command kinds (`CommandSetManual`,
    `CommandManual{Move,Ready,Exhaust,Attach,Place,Detach,Amber,Unforge,ForgeColor,Chains,House,AddCard}`)
    plus the fields they need (`Player`, `Delta`, `Name`), performed by
    `Game.ApplyManual(cmd, resolveCard)` (command_manual.go), mirroring
    `replay.go`'s `dispatchManual`. `LegalActions` never offers a manual kind, so
    the sim never drives manual mode (`TestLegalActionsNeverOffersManualKinds`).
  - **`session.Action` is a fixed closure, not a live driver.** `Stepper` takes one
    `func(*Game)` that must already encode the whole turn loop, and `Session`
    exposes no way to say "run this root action next". A web client needs each
    click to be the next command. So root actions must become a suspension point:
    the turn loop yields a `RequestAction`, the client answers with a root
    `Command`.
  - **Legality enumeration is now collected.** `Game.LegalActions(player) []Command`
    (command_action.go) gathers the legal root-action set from the gates that
    already exist (`CanPlay`, `CanDiscard`, `canUseTo`, `FightTargets`,
    `CanUseArtifact`, `canUnstun`): the house-choice set when
    `Phase == PhaseChooseHouse` (No House when the allowed set is empty), and in
    `PhasePlay` a creature play per meaningful flank, a discard, reap, a
    `CommandFight` per legal target, an Action, a stunned creature's `CommandUnstun`,
    each artifact Action, and `CommandEndTurn`. The stun path offers only Unstun
    (any use just sheds the stun) and an empty battleline offers one flank (both
    place identically). The contract is pinned by
    `TestLegalActionsAreAllApplicable`: every command `LegalActions` offers,
    `ApplyAction` on the same state accepts.
  - **The canonical turn loop now exists.** `Game.RunMatch(firstPlayer)`
    (command_action.go) deals the game and drives it to a winner: it asks the
    active player's `ActionChooser` for their next root action from
    `LegalActions` and performs it via `ApplyAction`, over and over, handing off a
    turn an effect ended without a paired `StartTurn` (Omega, Book of leQ) between
    actions. **DECIDED (human): the turn loop owns `StartGame`** — so the mulligan
    prompts surface through the same suspending choosers as every later decision
    and become recorded answers; the loop takes `firstPlayer` as a parameter. It is
    the standing `Action` a real match hands the session, so a client no longer
    hand-rolls a turn driver. Covered by `TestRunMatchDrivesToWinnerThroughEndedTurn`
    (walks setup, a house choice, a play, the Omega hand-off, and the win).
  - **`Orderer` is NOT a blocker.** `suspendChooser` does not implement it, so an
    order prompt degrades to repeated `ChooseCreature` — which is the engine's
    documented fallback and is fully expressible as a run of `CommandPickCard`s.
    The web's single drag-to-order widget is a UI affordance over that run, not a
    missing command kind.
  - **Manual/debug mode: command kinds landed; how the SESSION applies them is the
    open design knot.** The 13 `Manual*` edits now have Command kinds and
    `Game.ApplyManual` (above). But `internal/session` cannot yet apply one: its
    `Apply(cmd)` requires `request.IsLegal(cmd)`, and a manual force-edit answers no
    `Request` — it is an out-of-band poke, not a turn-loop answer. Worked-out design
    for the web rewrite: (1) `Session.ApplyManual(cmd)` applies the edit directly to
    the live `s.game` via `game.ApplyManual` (the stepper goroutine is parked in
    `yield`, so `s.game` is not being touched concurrently) and records it with
    `barrier=false`; (2) replay (`start`/`Undo`/`Load`) routes each recorded command
    by kind — manual kinds via `ApplyManual`, everything else via the stepper; (3)
    the `resolveCard` func for `CommandManualAddCard` is injected into the `Session`
    (it needs the card pool the engine cannot hold); (4) because an out-of-band edit
    staled the cached `RequestAction`, `Session.Apply` validates a `RequestAction`
    command against **live** `s.game.LegalActions(player)` rather than the cached
    `request.Actions` slice — the stepper's `ChooseAction` does not validate and
    `ApplyAction` re-checks live, so a manual edit between the yield and the answer
    is safe.
  - **Undo granularity differs.** Web undo works at root-action boundaries
    (`rootMarks`); `Session.Undo(n)` works at raw command index. Once root actions
    are commands the web still needs its own boundary bookkeeping on top.
  - **Also update ADR 0039's and ADR 0040's stale opening lines** — DONE: both now
    say the engine seam exists and only the web migration remains.
  - **PROGRESS (earlier session): the whole engine-side seam is complete.** The 11
    legal root actions have Command kinds + fields (suspend.go), `Game.ApplyAction`
    performs one (command_action.go, mirrors `dispatchRoot`), `Game.LegalActions`
    enumerates the legal set (bound to ApplyAction by
    `TestLegalActionsAreAllApplicable`), `RequestAction` + `Request.Actions` +
    `suspendChooser.ChooseAction` surface that set through the Stepper, and
    `Game.RunMatch` is the canonical turn loop that owns `StartGame` and drives the
    match through `ChooseAction`/`ApplyAction` — all at 100% coverage.
  - **PROGRESS (this session): the manual/debug command vocabulary landed too.**
    The 13 manual Command kinds + fields (suspend.go), `Game.ApplyManual`
    (command_manual.go), and the sim-safety assertion
    (`TestLegalActionsNeverOffersManualKinds`) are in at 100% coverage, and the
    ADR opening lines are corrected. So the **entire engine + ADR side is done**.
    **Remaining — only the web client rewrite is left, and it is a large,
    browser-validated epic** (the human gave the green light to use Playwright):
    - Extend `internal/session` with the manual-command support designed above
      (`ApplyManual` + replay routing + injected `resolveCard` + live
      `RequestAction` validation). `internal/session` is ungated — test it well but
      100% is not required.
    - Invert the web client's concurrency model: today it drives each root action on
      a background goroutine and renders prompts through a **live** `webChooser`
      that posts-and-blocks and records its own `input`s. The session model has no
      live chooser — prompts surface as `session.Pending()` `Request`s answered by
      `session.Apply(cmd)`, and `RunMatch` (not the client) owns the turn loop. So
      every click becomes a `Command` fed to `session.Apply`, and the prompt UI
      reads `Pending()` instead of `g.choosing`/`chooserCandidates`. This touches
      game_action.go, game_chooser.go, game_play.go, game_manual.go,
      game_persist.go, game_setup.go, undo/redo, and the flash/log-group
      bookkeeping.
    - Persistence: `save`/`resume` write/read `session.Record` (seed, sets,
      `[]Command`) instead of the `[]input` snapshot; bump `snapshotVersion`.
    - Undo/redo: `Session.Undo(n)` works at raw command index; the web still needs
      its own root-action-boundary bookkeeping (`rootMarks`) on top, plus a redo
      stack (session has no redo).
    - Delete `internal/web/replay.go` and rewrite the ~310 web tests that reference
      `g.inputs`, `replayGame`, `driveReplay`, `rootMarks`, and the `input` type.
    - Validate with the go-app unit tests AND Playwright (browser). Must build for
      host AND js/wasm (`mage webWasm`).
      This was NOT started this session: it is architecturally loaded (the
      concurrency inversion), can only be fully validated in a browser, and half a
      rewrite would leave the web package non-compiling. It is a coherent next chunk,
      ideally driven with the human present for the Playwright loop.
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
    manual mode** — RESOLVED as chosen: `LegalActions` simply never offers a manual
    command kind, so the restriction is a property of legality rather than of who is
    asking, and one command vocabulary is kept. Asserted by
    `TestLegalActionsNeverOffersManualKinds`.

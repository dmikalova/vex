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

The engine pass is done: the `gocognit` gate's global `min-complexity` is now 30,
with every genuinely nested engine/deckgen/cardtest seam split (fight, NewCard,
hasKeyword, the text.go renderers, Target.Text, Target.filter, triggeredBy,
RenderAbility, Harness.location, allowedHouses) and 100% coverage held. What
remains are the web renderers, kept exempted by name in the `.golangci.yaml`
`exclusions` `text:` list until each gets its own pass:

- `effectGlyphs` (internal/web/icon.go, 90) — glyph dispatch; wants its own
  glyph-family grilling session (see docs/todo.md "Split out glyphs more in
  icon.go"). Split along glyph families, then delete its exclusion.
- `installTips` (internal/web/game_lifecycle.go, 49) — tip installation; split by
  tip group, then delete its exclusion.
- `installKeyShortcuts` (internal/web/game_lifecycle.go, 39) — keyboard-shortcut
  registration; split by shortcut group, then delete its exclusion.
- `computeFlashes` (internal/web/game_action.go, 31) — per-event flash selection;
  split by flash source, then delete its exclusion.
- `renderMarkdown` (internal/web/markdown.go, 31) — inline markdown rendering;
  split by span kind, then delete its exclusion.

When these are split, lower the global toward the tool default and re-check for
any new stragglers one at a time.

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
shape the item described.

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
  row-only with ~14 callers in ability scanning, phase processing, and invariants.
  Decide per caller whether an upgrade belongs in that scan; do **not** change it
  wholesale.

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

- **No implemented card puts an upgrade into play.** `putIntoPlay` has only
  `Creature` and `Artifact` arms, so the two `AfterUpgradeEnters` cards (Armory
  Officer Nel, Commander Dhrxgar) can never fire from a put-into-play. Fine now;
  a trap for the first card that needs it.
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

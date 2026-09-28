# Engine design guide

`internal/engine` is the pure game engine: it imports nothing upward and is held
to 100% coverage, so a new code path needs an engine test (a card test does not
count toward the engine gate). This file is the design ideal: the patterns the
engine is built from, where each belongs, and the deliberate tradeoffs, some of
which look like smells but are forced by a constraint. The _why_ behind the
composition rules is the style guide's "Composition and design".

## Two hard constraints shape everything

When something looks un-idiomatic, check it against these and their ADRs before
"fixing" it.

1. **Flat, pointerless, comparable state** (ADR 0005). `GameState` is a value
   struct of fixed arrays (`[maxCards]CardCore`, `[2]Zone`), so
   `GameState.FastCopy()` is a plain value copy for MCTS cloning and undo is a
   snapshot. Never introduce a pointer, slice, or map into `GameState` or into
   a value compared against it.
   - State holds no closure or `Effect`; "do X later" is flat, enum-tagged data
     (the lasting registry, ADR 0007).
   - Values compared against state (notably `Target`, `== Target{}` in
     `ConstantAbility`) stay comparable: a wide flag struct, not a `[]filter`;
     paired `x int` + `hasX bool`, not `*int`.
   - Cards are `LocalID uint8` indices into a read-only `catalog` shared across
     clones.
2. **Effects touch the game only through the `Resolver` port** (ADR 0008). An
   `Effect` holds a `Resolver` via `EffectContext`, never `*Game` or
   `GameState`; the port is the complete, auditable catalogue of card-facing
   capability (the effect AST is the domain, `Resolver` the port, `*Game` the
   adapter).

## File organization

The filename says which responsibility a file holds:

- `game.go` + `game_*.go` — methods on `*Game`, grouped by area
  (`game_turn.go`, `game_play.go`, `game_read.go`, `game_abilities.go`,
  `game_combat.go`, `game_leaves_play.go`, …). A new `Game` method goes in the
  matching `game_*.go`, unless it is one mechanic's own gate or read.
- `effect.go` + `effect_*.go` — the effect AST, one file per mechanic. A new
  effect goes in `effect_<mechanic>.go`. **A mechanic's own gates and reads
  stay in its file**, even when they are `*Game` methods: Alpha's
  `barredByAlpha`, Omega's `endStepIfOmega`, Deploy's
  `deployPosition`/`chooseFlank`/`choosePosition`, copy-stats'
  `copiedStatsSource` and `CopyStats`, the text box's `grantedTextBoxSources`
  and `GrantTextBox`. A method that is not part of one mechanic's seam goes in
  the matching `game_*.go` (`SetActiveHouse` is in `game_turn.go`).
- `mechanic_*.go` — a self-contained ambient mechanic that is neither an effect
  node nor a `Game`-method area: a small flat state type plus its reads and
  gates (`mechanic_tide.go`, `mechanic_toll.go`, `mechanic_counter.go`, …).
  A mechanic already clustered by a family prefix (`house_*`) stays there, and
  execution-model plumbing (`suspend.go`, ADR 0040) keeps its own concept file.
- Everything else is a concept file named for what it defines (`card.go`,
  `state.go`, `target.go`, `duration.go`, `destination.go`, `resolver.go`,
  `text.go`): a pure value, enum, or display type with no rules. A new enum or
  value type gets its **own** `<concept>.go` (destinations are in
  `destination.go`, not `game_destination.go`; `game_` is only for `Game`
  methods).
- When a vocabulary file grows unwieldy, split it by category into files that
  keep the family prefix, so `ls` groups the family: `effect_condition.go`
  keeps the framework (the interface, `Comparison`, `Conditional`, `Or`,
  `CountIs`) and its conditions move to `effect_condition_board.go` /
  `_source.go` / `_it.go` / `_turn.go`; `effect_count.go` splits into
  `_board.go` / `_turn.go` / `_produced.go`; `target.go` into
  `target_refinement.go` + `target_select.go`; `resolver.go` into
  `resolver_game.go`; `text.go` into `text_helpers.go`.
- Tests sit beside their source in `<name>_test.go`, named for the **source
  file that defines the symbols they exercise**, never for the card or
  scenario that motivated them: a `TakeControl` test goes in
  `effect_control_test.go`, not `effect_saurian_egg_test.go`. Split a test that
  spans two mechanics across both mechanics' test files. The one exception is
  `helpers_test.go`, the shared test infrastructure.
- The `internal/card` facade mirrors these namespace and type files
  (`duration.go` → `card/duration.go`, `destination.go` →
  `card/destination.go`), and the rules in this file govern it too.

## Patterns in use

- **Interpreter — the effect AST (`Effect`)** (ADR 0006). Every node has
  `Text()` (renders English) and `Resolve(ctx)` (carries it out), so printed
  text cannot desync from behavior. A new mechanic is almost always a new node
  in `effect_<mechanic>.go`, not a new branch in the `Game` runtime.
- **A brand-new `Effect` node or `card` facade type needs a grill-me session
  (the `grilling` skill) the human signs off on**, whether or not you came
  through `implement-cards`. A new field, `Strategy`, `Target` filter, or
  `Count` on an existing node needs no ceremony. A new node widens the shared
  vocabulary permanently, and cramming a mechanic into an existing node in a
  way that is not clean and composable is the same smell: stop and present the
  grill (why it is necessary, the cluster of real cards that need it, the
  alternatives rejected, a before/after authoring comparison). If the grill
  does not clearly favour the new node, extend an existing one.
- **Composite — `Sequence`, `Conditional`, `ChooseHouseThen`, `Repeat`, …**
  compose child `Effect`s and recurse `validateEffect` into them; compose small
  nodes over one fused node. `Sequence` is the only ordered composite: its
  children are separate sentences by default ("A. B."), a run of folding
  children conjoins ("destroy a creature and an artifact"), and a
  `Conditional` gate keeps its consequence joined so it visibly covers every
  clause. There is no per-child sentence wrapper and no second list node.
- **Strategy** — see below; reach for it first when behavior varies along an
  axis.
- **Ports & Adapters — `Resolver`** (ADR 0008); see the roles below.
- **Facade — `internal/card`** wraps the engine so card files never import
  `engine`; `card.New(name, …, WithFoo)` is a functional-options builder.
- **Null Object / invalid-zero** (ADR 0010). `FirstChooser` is the default
  strategy; the `playerUnset` / `targetUnset` / `durationUnset` / `eventUnset`
  sentinels make an omitted required field a card-init `validate()` error, not
  a silent default.
- **Visitor** — a whole-tree operation that is not part of a card's identity;
  see below.

## Strategy: keep the text with the behavior

A strategy carries **both** its behavior and its printed-text fragment, so it
plugs into the AST without desync. When behavior varies along an axis, model
the axis as a small strategy that renders its own text, not a new
`Effect`/`Target` field or a `bool`.

- **`Chooser` (`game.go`)** is the decision strategy, swapped per frontend:
  `FirstChooser` (bot/tests, deterministic), `suspendChooser` (interactive: it
  yields a `Request` the web client answers), `bridgeChooser` (test harness).
  Extend it with an **optional capability interface** found by type assertion
  (`OptionChooser`, `Orderer`) with a graceful fallback, not by widening the
  base `Chooser`.
- Every asking capability takes a leading `PromptSource` (`prompt.go`) naming
  the card that raised the prompt, then the already-rendered prompt text
  (`renderPrompt` substitutes `SelfName`). The source is **identity, not
  text**: a client points at the card with it and never rebuilds the sentence.
  Add new prompt context as a `PromptSource` field, not a parameter on all 68
  chooser implementations. `ReactionChooser` and `BadgeChooser` take no source:
  a window is not raised by one card (each `OrderableReaction` names its own),
  and a badge decorates a pick that already carries one.
- **`Refinement` (`target.go`)** is a set-relative rule (`refine` + `clause`)
  such as `MostPowerful` or `Except(MostPowerful)`: it narrows the ids and
  contributes a phrase. Add a `Refinement`, not another `Target` bool, for a
  whole-set rule.
- **`Count` and `Condition`** are value/predicate strategies with paired text
  (`CountText` / `CondText`). A number that scales with the board is a `Count`;
  a branch is a `Condition` fed to `Conditional`.
- **`CreatureVerb`** is a per-creature verb strategy for `OnChooseCreature`.
- **`RepeatGate` (`effect_repeat.go`)** is the axis `Repeat` varies along:
  **why the loop stops**. `While` stops only when its `Condition` fails (so the
  Rule of Six, ADR 0043, may be the only bound); `WhileYouDo` also stops when
  the effect does nothing; `MayWhileYouDo` also when the controller declines;
  `ByExalting` never loops, offering one more resolution paid by exalting a
  creature. Gates that ask whether the effect happened print rule 5's `->`;
  `While` prints two plain sentences. A new repeat shape is a new gate, not a
  new node.
- Changing `Stepper` (`suspend.go`) or how an interactive action ends: read
  [docs/engine/stepper.md](../../docs/engine/stepper.md) first.

## Reused effect shapes get one shared helper

When a shape (a field cluster and the logic turning it into a value, phrase, or
validation error) shows up in a third effect, factor its logic into one helper;
each effect keeps its own authoring fields. Prefer a shared helper over a
shared embeddable value type. Before writing or extending an effect's amount,
`Per`/`Times` scaling, `Produced` tally, zone movement, multi-pick, or
redistribution, read
[docs/engine/effect-shapes.md](../../docs/engine/effect-shapes.md): it names the
helpers already factored (`scaled`, `errAmountOr`, `poolAmount`,
`crossZoneMover`, `pickCards`, `placeAmong`, …).

## `Resolver` is segregated into role interfaces

`StateReader` (reads) · `EconomyResolver` (Æmber/keys/chains) ·
`CreatureResolver` (per-card in-play state) · `BoardResolver` (duration-scoped
rules over creatures collectively) · `CombatResolver` (damage, destruction,
ability-driven fight/reap/action) · `PlayResolver` (playing a card from a zone,
putting one into play, the play sequence) · `ZoneResolver` (movement between
zones without playing, plus draw) · `TurnResolver` (turn-scoped grants and the
lasting registry) · `ChoiceResolver` (ordering and choosing) · `Logger`.

A new engine capability is a method on the role it belongs to, implemented on
`*Game`. If it fits no role, a new area is emerging: add a small role interface
and embed it in `Resolver`.

## Every interface and enum this package declares is classified

The totality tests (`catalog_test.go`, `catalog_enum_test.go`,
`catalog_interface_test.go`; ADR 0018) fail the build on anything
unclassified. `mage tool:census` only reports what the census still owes, so it
reads while half-filled: per node family (`engine.Families()`) the node types
no census row covers and the rows naming no type, how the interfaces are
classified, and the rulebook terms still to write or unclaimed, counting the
text-bearing value enums (`engine.Enums()`).

- A new `type X interface` fails `TestInterfaceTotality` until `Interfaces()`
  in `catalog_interface.go` has a row for it. The row says either **a
  catalogued family** (the interface a card's meaning varies along, whose
  implementations print their own text; the row names its catalog, a
  `Families()` entry or `LogEntrySamples` for `LogEntry`, so a new strategy axis
  gets a `catalog_<family>.go`, as `Quantity`, `Gather` and `TopAct` did), or
  **not a family, with the reason** (a port such as `Resolver`'s roles,
  `Chooser` and its capabilities, `Namer`; an optional capability that only
  reshapes a catalogued node's text, such as `negatable`, `combinable`,
  `framedClause`; an internal shape two value types share, such as
  `cardPile`). Write the reason so the next reader can judge from the row alone
  whether it is wrong or has since grown text-printing members. When in doubt,
  classify it as a family and write the catalog.
- A new `Keyword`, `Duration`, or `CounterKind` constant fails the build until
  its enumerating function (`Keywords()`, `Durations()`, `CounterKinds()`)
  returns it, or the enum's `Excluded` list names it with the reason it is not
  a real member.
- **Keep the rules current.** Adding or changing a keyword, trigger, card type,
  or rule-bearing effect registers or updates its rulebook term in the same
  change.

## New whole-tree operations: a type-switch Visitor, not a new AST method

`Text()` and `Resolve()` are intrinsic to a card's identity. An operation that
is not (an MCTS value estimate, a static analysis, serialization) is a
standalone function that type-switches over `Effect` in one file, never a third
method on every node. Follow `lastingActionOf` (`effect_lasting.go`). Keep each
such function in one file, next to a comment listing the effects it covers; the
coverage gate forces every branch.

## The lasting registry is the flat-state interpreter

"For the remainder of the turn" behavior (Full Moon, Charge!, Crystal Hive,
Dimension Door) is a smaller interpreter, because flat state cannot hold an
`Effect` closure (ADR 0007). State holds flat
`LastingEffect{On Event, Do lastingAction, Controller, Amount}` records;
`lastingActionOf` maps a composed effect to an enum tag and `game_lasting.go`
fires and queries them.

- Never add a bespoke `if g.State.Foo…` block to the play or reap path
  (`PlayCreature`, `reapWith`); that hardcodes each effect into the hot path and
  makes effects sharing a timing window unorderable. Everything routes through
  the registry as a **reaction** or a **replacement**.
- A card in play granting every creature an ability is a `ConstantAbility`
  with a `Target` and `Granted` abilities, never a global override in
  `discardDestroyed` or another leave-play path.
- A **reaction** runs after an event: every site firing the event folds its
  reactions into that event's trigger window with `lastingReactions`, ordered
  with card abilities through the flat `ReactionChooser` port (ADR 0013). An
  event with no card ability of its own still gets a window
  (`afterDestroyedReactions` folds `EventEnemyCreatureDestroyed` in). A
  **replacement** changes an event's outcome (`lastingReplacement` +
  `Instead{Of, With}`).
- To add a reaction on an existing event, support its `Do` in `lastingActionOf`
  and `resolveReaction`. A new event is an `Event` value, one
  `lastingReactions` or `lastingReplacement` call at the site, and its
  `clause`/`gerund` text. Never restructure the play/reap hot path; keep the
  enum dispatch centralized.
- **Modifying pending damage is a replacement, not a reaction.** "Whenever a
  creature takes damage, it takes an additional N" (Lethal Distraction) is to be
  modelled as a replacement on a "damage about to be dealt" event that
  increases the pending amount, retiring the bespoke `TakesExtraDamage`
  (ADR 0031; not yet built).
- **Æmber flow is replaced at either endpoint.** The destination is
  `EventAemberAddedToPool` (Ether Spider captures Æmber before it lands); the
  source is `EventAemberTakenFromPool` (Po's Pixies draws a steal or capture
  from the supply). Both ride the one `Instead{Of, With, Player}` in
  `CardDefinition.Replaces`, read by a scanner scoped to the watched pool
  (`aemberCaptorFor`, `AemberTakenFromSupply`). A new redirect is a `Replaces`
  on the matching event plus a `Replacement` enum member, never a bespoke
  `CardDefinition` bool. A scanner returns the **card it found**, and every
  replaced outcome narrates through `replacementLine` (`log_replacement.go`):
  `<cause> has <actor> <outcome>, instead of <displaced>`.
- Card-side authoring of lasting effects is in
  [docs/card-implementation.md](../../docs/card-implementation.md).

## Side-tables: house constraints, counters, take control

Forward house choice (`HouseConstraints`), card-placed generic counters
(`Counters`), and the take-control stack (`Controls`) are flat side-tables on
`GameState`, never a field per card or per kind. Before changing one, or adding
a mechanic that reaches a future house choice, places a marker, or seizes
control, read [docs/engine/side-tables.md](../../docs/engine/side-tables.md).

## Power is dynamic: settle destruction after anything that can lower it

Power is computed (`Power` in `game_read.go`: base, upgrades, power counters,
temporary bonuses, constant abilities in play), so a creature can become
destroyable with no effect touching it. `shouldDestroy` names that state and
`settleDestroyed` notices it, but only where it is called.

- **Every code path that can lower a creature's power ends by settling
  destruction**: a new counter, constant, blank, control swap, `Per`-scaled
  bonus, or card leaving play ends in `g.settleDestroyed(controller)`, directly
  or through `removeFromPlay`. A missed settle fails silently and surfaces turns
  later as a sim invariant ("in 2 places", "power above damage").
  `AddPowerCounter`, `BlankEnemyText`, `SetAember`, and the key-forge paths all
  settle; new siblings must too.
- **Cards leaving play at once leave together.** A batch removal (archive,
  destroy, return, shuffle) must not let the first removal's settle destroy a
  later member: Epic Quest archives "Lion" Bautrem and the neighbor it buffed,
  so the neighbor is archived, not destroyed and then archived into two zones.
  Hold the `settling` flag across the batch and settle once at the end
  (`destroyBatch`, `putIntoArchivesEach`); never loop a per-card leave-play
  call over a pre-selected list.
- **A card is not "destroyed" until it reaches the discard pile.** During
  `destroyTogether` every dying creature stays in play while the batch's
  `Destroyed:` abilities resolve; only then does each go to discard. The
  **"after ... destroyed"** reactions (`emitAfterCreatureDestroyed` for Neffru,
  `emitAfterEnemyDestroyed` for Pile of Skulls) fire after that move, so a
  creature killed in the same batch can neither be chosen by them nor react
  (Pile of Skulls cannot capture onto a friendly creature that died in the same
  combat). Never fire an "after destroyed" reaction, or offer a prompt that
  could pick a dying creature, inside the destruction window.

## Event, ability, and effect verbs: emit → trigger → resolve

- **`emit<Event>`** announces an event and fans out to every listener
  (`emitReapWindow`, `emitEnters`, `emitLeavesPlay`); its site folds duration
  reactions into the window with `lastingReactions`.
- **`trigger…`** resolves a _single card's_ abilities matching a trigger
  (`triggerAbilities(id, TriggerAfterReap, …)`); emitters call it per card.
- **`resolve…`** carries out one effect or ability (`Effect.Resolve`,
  `resolveReaction`, `resolveUpgradePlay`).
- "Fire" is for prose ("a trigger fires"), never a method name.

## The game log is typed entries, not sentences

The log narrates what **resolved**, not what a card promised (ADR 0011). There
is no `Logf`.

- **A new log line is a new `LogEntry` variant**: a small comparable struct in
  `log_<family>.go` carrying the outcome's ids and amounts, plus a
  `Text(Namer) string`. Record it with `g.record(entry)` (or `Resolver.Record`
  inside an effect), at the site the outcome lands. Carry the `LocalID` and let
  `Text` ask the `Namer`; never format a name into a string at the call site.
- **Attribution is a `Frame`, not a line.** Open one around the resolution
  (`closeFrame := g.openFrame(Frame{Actor, Source, Trigger, Grantor})`) and
  every entry recorded inside inherits it. A card's printed text is never a
  log line.
- Whole-tree passes over an entry are extrinsic (`RenderEntry` in
  `log_render.go`); never add a second method to every entry.
- **Reuse the shared phrasing.** `log.go` has `namedCards`, `because(text, on)`,
  and `nameMoved`; `text_helpers.go` has `countNoun`/`plural` and `indefinite`.
  Never hand-roll `card(s)`, `noun + "s"`, or `"a " + noun`. An entry that is
  another entry plus context renders the base:
  `LastingAemberGained.Text` is `because(AemberGained{…}.Text(n), e.On)`.
- **Zones decide whether a card may be named.** `Zone.public()` says which
  zones both players see (discard, the board, purged; not hand, archives, or
  deck), and `nameMoved(n, id, from, to)` names a card once either end is
  public, else "a card". A zone-change entry states both zones and calls
  `nameMoved`, never `Namer.Name`. A hidden card becomes nameable only through a
  `Reveal`, which records its own entry.
- **Discards resolve one at a time, one log line each**; never group them into
  "discards X and Y" (the planned scrap mechanic acts on each discard).
- Recording is switchable (`SetRecording`), so search pays nothing for it.
- **Two ratchets fail the build on unnarrated state.** `fieldNarration`
  (`state_test.go`) classifies every `GameState` field by dotted path
  (`Cards.Damage`, `Controls.Controller`) as narrated directly, narrated by its cause, or never; a new field fails
  `TestEveryStateFieldDeclaresItsNarration` until classified.
  `methodNarration` (`narrationaudit/narration.go`) classifies every method of
  the mutating `Resolver` roles, and its `Auditor` (installed by `cardtest` and
  the seeded sim batch) fails a directly-narrating method that changes state
  without appending an entry. A new port method therefore decides where its
  line comes from; "it narrates" means recording inside the method body.
- Rewording a log entry, or changing how a recorded command resolves, makes old
  saved matches replay wrong: bump `snapshotVersion` in `internal/web` (see
  `internal/web/AGENTS.md`, "Persistence").

## Prompts: name the card, and let the player click a target

- **The active player makes every choice**, an unbreakable KeyForge rule: the
  chooser for a placement, target, order, or option is `ActivePlayer`, even
  when it lands in the opponent's play area (the flank a Treachery creature or
  a give-to-opponent seize enters on, where a creature put into play under the
  opponent lands; the rulebook: whenever a creature enters play or changes
  control, the active player chooses its flank). Never pass a non-active
  controller as the chooser. A creature entering a battleline funnels through a placement seam —
  `deployPosition`/`chooseFlank` (play, put into play), `placeGainedOnFlank`
  (control gains, gifts), `placeSeizedOnFlank` (seize) — each asking the active
  player, never assuming a flank (ADR 0010). A new placement site uses a seam,
  never `Battleline[...].add` directly.
- Ask through `pickCreature`/`pickCard` whenever the answer is a card, so a
  frontend highlights candidates and takes a click; `ChooseOption` only for
  non-card options (a house, a key colour, yes/no).
- Refer to the source card with `SelfName` (`"fully heal "+SelfName`), never a
  literal name.
- Phrase the prompt as the card's instruction, lowercase and imperative
  ("choose a creature to attach {self} to").
- An optional choice ("may", "up to N") must let the player stop: the prompt is
  declinable (a frontend shows _Done_), and it is **not** short-circuited when
  one candidate remains. `pickCreature` auto-takes a sole candidate because it
  models a mandatory choice; an optional one needs its own path.

## Deliberate tradeoffs: handle them, don't "fix" them

- **Wide `Resolver`** (ADR 0008): add to the right role.
- **`Filter` as a wide comparable struct + paired kind/amount values** (ADR 0005).
  `Target` is `Target{Kind, Filter, Refinement, Neighbors}` and must stay
  comparable, which rules out a `[]filter` slice and `*int` optionals. A new
  **per-card** axis is a **field on `Filter`** (`filter.go`) plus its census row
  in `catalog_filter.go` — the census discovers the family as `Filter`'s
  exported fields, so an uncatalogued axis is a red build. Route **set-relative**
  rules — anything needing the whole candidate set — through `Refinement`
  instead. An optional number is a named comparable kind-and-amount value
  (`PowerBound`), never a `*int`. Do **not** add a new `Target` field for a
  narrowing: `Target` has four. The line between the two is in ADR 0005: a
  filter is decidable one candidate at a time however much board it reads; a
  refinement needs the whole set.
- **Enum-tagged lasting records instead of stored closures** (ADR 0007).
- **`panic` in `EffectContext.PlayerFor` on `playerUnset`** (ADR 0010) is
  acceptable **only** because `validate()` at init keeps real cards from
  reaching it. Never let a _computed_ `Player` reach `PlayerFor`; reject unset
  at the boundary.
- **`Game` implements all of `Resolver`**: kept in check by the `game_*.go`
  split and the role split. If an area outgrows its file, split the file.

## Before you add: check the surface that already exists

The `1 keys` bug (a hand-rolled plural where `countNoun` existed) is the failure
this prevents. Each entry names the **authority file**; enumerate it rather
than trusting this list.

- **Word a number or a name — `text_helpers.go`**: `countNoun(n, noun)` (never
  `noun+"s"` or `if n == 1`), `plural`, `indefinite`.
  `grep -n '^func ' internal/engine/text_helpers.go`.
- **Word a log line — `log.go`**: `namedCards`, `because`, `nameMoved`. Find
  the entry family first:
  `grep -rln 'Namer) string' internal/engine/log_*.go`.
- **Settle destruction — `game_settle.go` / `game_leaves_play.go`**:
  `settleDestroyed`, `removeFromPlay`; never re-derive destruction inline.
- **Read legality or a derived stat — `game_read.go`, `game_play.go`,
  `game_abilities.go`**: `Power(id)`, `CanPlay`/`CanDiscard`/`CanUse`. Add a
  read to the matching file and its `Resolver` role, after enumerating the
  role.
- **A new state-changing outcome needs its log line**: `internal/sim`'s
  step-narration invariant fails otherwise.
- **The facade by usage — `mage tool:nodeUsage`**: every exported name in
  `internal/card` by category, with card uses (`CARDS`) and other references
  (`OTHER`). Use it to find the neighbours a new node should match and to audit
  that a thin node is built from atoms. Low usage is not a defect, and `CARDS`
  0 with non-zero `OTHER` is not dead code (type surface, registry plumbing, an
  enum member no card names); only both at zero is unused. Narrow with
  `-max=<n>` and `-category=<substring>`
  (`mage tool:nodeUsage -max=1 -category=damage`).

## Adding a mechanic: the decision order

1. Can an existing `Effect` express it by changing a `Target`, `Count`,
   `Condition`, or `Refinement`? Prefer that; no new type.
2. Is it a new _node_ (after the grill)? Add an `Effect` (or
   `Condition`/`Count`/`Refinement`) in the matching `effect_*.go` /
   `target.go`, with `Text()` + `Resolve()` (+ `validate()` for an illegal
   field combination), and a facade alias in `internal/card`.
3. A new engine capability? A method on the right `Resolver` role, implemented
   on `*Game` in the matching `game_*.go`.
4. "For the remainder of the turn"? The lasting registry, not the play/reap
   path.
5. An extrinsic whole-tree operation? A type-switch function.
6. Does it say what happened? A `LogEntry` variant in `log_<family>.go`,
   recorded where the outcome lands.

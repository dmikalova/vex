# AGENTS.md

Repo-wide guidance for agents. See `internal/cards/AGENTS.md` for
card-authoring specifics.

## "Reminder" means write it down

When a request says **"Reminder"** (or "remember this", "document this"), it is
not asking for a one-off fix — it is asking you to capture the rule in the
appropriate durable place so it holds for future work: the relevant `AGENTS.md`,
a `docs/` page, an ADR, or a code comment on the seam it governs. Make the change
_and_ record the rule; a fix without the write-down is only half the task.

## An explained mechanic becomes a test

When the human **explains a nuanced game mechanic** — how ward interacts with the
destroyed tag, when a "Destroyed:" ability still reaches a card, which of two
similar timings wins — that explanation is a specification, and the turn is not
finished until it is **pinned by a test**. Write the test even when the code
already behaves correctly: the point is not to fix a bug but to stop the next
agent from "simplifying" the nuance away, because a rule that lives only in a
conversation is a rule that gets refactored out. A passing new test is a success,
not a wasted change.

Alongside the test, **add the documentation if it does not already exist** — a
short comment on the seam the rule governs (a doc comment on the funnel, not a
note buried in a caller), and a rulebook term or `docs/` line when the rule is
player-facing. Cite the test by name in that comment so the claim and its proof
are findable from each other. Do not write a markdown file just to record the
explanation; the test and the seam comment are the durable record.

Where the mechanic contradicts what the code does, the human's explanation is the
authority — but say so plainly, and fix the code rather than quietly reshaping
the explanation to match what is already there.

## Agent todos go in `docs/todo-agent.md`, not `docs/todo.md`

`docs/todo.md` is the **human's** personal list — do not write into it. When a
request implies work you want to queue, plan, or show progress on, write it in
[docs/todo-agent.md](docs/todo-agent.md): a scratchpad where you translate a
request into concrete, grouped work items. When an item is **done, delete it**
(do not mark it done) — the file only ever shows what is still outstanding, so it
reads as a live surface for coordinating with the human on what you mean to do
next. Group items by area or mechanic so related work is built together.

A `todo-agent.md` item is a **handoff to a future agent who was not in the
conversation that wrote it**, so it must carry the decision, not just the task.
When you record an item, write down _what was decided and why_ — the chosen
behavior, the cards affected, the expected text — so the next agent does not have
to reconstruct it from code that may already be stale. When you pick an item up,
the recorded decision **wins over a contradicting code comment**: a comment that
disagrees with the item is out of date (it describes the behavior the item exists
to change), so fix the comment to match the decision — do not treat the comment as
evidence the item is wrong and re-litigate it.

Work that is decided but has **no consuming card in an implemented set yet** goes
in [docs/todo-future-set.md](docs/todo-future-set.md) instead, keyed to the set
that first needs it. When you implement or stub a set, scan that file for items
naming it and build the primitive alongside its first real consumer (the
`implement-cards` and `stub-cards` skills both point there).

**You sequence approved work; the human does not.** Once work is agreed, do not
end a turn asking which item to do first or offering a menu of plans. Derive the
order: dependency first (a shared primitive, renderer helper, or enum fold lands
before the cards and call sites that consume it), then easiest-win first within a
tier (the mechanical change with no rules risk before the one needing a judgement
call). State the order in one line and start.

## Build, test, and lint through `mage`

Run all build/test/format/coverage tasks through `mage`, not raw `go`
commands. The generic gate comes from project-standards' shared `ci` targets
(imported in `magefiles/build.go`); the rest are vex's own targets, which wrap
the project's conventions (comment/rulebook generation, the simulator, the card
tools):

- `mage ci:fix && mage ci:check` — **the local validator**; run it before
  considering work done. `ci:fix` applies every autofix (go fix, golangci-lint
  `--fix`, goimports, golines, gci, goldmark-lint `--fix`, misspell `-w`) and
  regenerates the generated configs. `ci:check` only verifies — it writes no
  file and fails on anything `ci:fix` would change — then runs format, tidy,
  build, vet, lint, markdown, spell, secrets, commits and drift, then the tests
  and the coverage gates. It must print `ALL GREEN`.
- `mage ci:build` — build all packages for the host, then the web client for
  js/wasm (into a temporary directory; `mage webWasm` writes `web/app.wasm`).
- `mage ci:test` — run all tests.
- `mage testRun <pattern>` — run only the tests matching a name pattern.
- `mage ci:vet` — run `go vet`.
- `mage ci:lint` — run golangci-lint, fixing nothing.
- `mage ci:format` — check formatting without writing (`mage ci:fix` formats).
- `mage ci:markdown` / `mage ci:spell` — lint Markdown, check spelling.
- `mage ci:cover` — report coverage for each gated area, which must stay at 100%.
  The areas are `internal/engine`, the card definitions under
  `internal/cards/sets`, `internal/cards/cardtest`, and `internal/deckgen`; they
  are listed as `ci.CoverGates` in `magefiles/build.go`. `internal/web` is
  deliberately ungated.
- `mage generateComments` — rewrite each card's doc comment from its definition.
- `mage gen` — regenerate card comments and the rulebook. It is deterministic and
  regenerates from the tree, so **just run it** — you never need to worry about
  what else is in flight, guard against sibling edits, or hand-write a comment to
  match the generator. Run `mage gen` freely whenever a definition changes.
- `mage docs` — serve this module's Go documentation at `http://localhost:6060`
  (pkgsite, the pkg.go.dev renderer); read-only, no gate depends on it.
- `mage debug` — replay a simulated game with the game log on and print the log
  tail next to the invariant violation that ended it. With no `-script` it finds
  the first failing game in the fixed-seed property batch `mage ci:test` plays; pass
  `-script` the hex a failure printed to replay that one, and `-tail` to widen
  the log (`mage debug -script=<hex> -tail=200`).

  When an invariant names a card, the violation is a symptom — the cause is
  usually an earlier line and a card no longer in the frame. Read the named card's
  whole lifecycle, not just the tail: widen `-tail` (or `mage trace` the game to a
  file) and grep the log for the card by name to see when it entered, what damage
  and power it showed, and which other card was buffing, blanking, capturing, or
  neighboring it. A creature that dies "for no reason" almost always lost a buff a
  card that has since left play was granting — so identify the cards that were in
  play around it, not only the card the invariant printed.

  Read the whole deck, not only the cards in the log. The card that causes an
  invariant may never appear in the tail — the log names the victim, not the
  culprit. `mage debug` prints both players' full deck lists above the log for
  exactly this reason: scan every card for the mechanic that could produce the bad
  state — a power reducer for a 0-power creature, a blanker for a creature that
  lost its ability, an attachment for a stat that drifted. Suspect the mechanic
  first, then find which card in the deck carries it; do not assume the only
  suspects are the cards printed near the violation.

- `mage trace` — play the fixed-seed property games once with the game log on and
  write every line to `tmp/sim/trace.log` (gitignored), so a whole game reads end
  to end. Where `mage debug` shows the tail of the game that broke, a trace is the
  full log of games that pass. `-count` sets how many games (default 1), `-out` the
  destination (`mage trace -count=25 -out=tmp/sim/mine.log`).
- `mage profile` — profile the engine under the whole-game simulator (the closest
  proxy for MCTS load), write CPU and allocation profiles to `tmp/sim`, and print
  the per-game counts next to a committed baseline so a regression shows as a delta.
  Default runs the 1000 seeded games; `-random` runs a fresh random outlier batch;
  `-save` re-blesses the baseline. It never blocks and never gates — the baseline is
  advisory and lives outside `mage ci:check` and CI. See `docs/testing.md`.
- `mage profileServer` — open a profile from the last `mage profile` run in the
  interactive pprof web UI (flame graph, call graph, source view); blocks until
  stopped. `-mem` serves the allocation profile, `-random` the outlier run's.
- `mage corpusPrune` — replay every entry in `FuzzPlay`'s seed corpus and rewrite
  it as one minimized entry per bug that still reproduces, dropping the entries
  whose bug is fixed. The corpus is the list of open findings, not an archive of
  every script a soak ever saw; run this after fixing a soak or fuzz find.

When you add or update a mage target that writes a binary, keep the output in an
ignored location (for example `./bin`) or add it under `ignore.git` in
`mklv.config.json`. Generated binaries should not be committed.

Tool configs are generated, never hand-edited: `.golangci.yaml`,
`.markdownlint-cli2.yaml`, `.commitlint.yaml`, `.gitleaks.toml` and `.ruleguard.go`
are project-standards' base configs merged with the overrides in
`mklv.config.json`. Change an exclusion or a lint setting there, then run
`mage ci:fix` to regenerate; `ci:check` fails on a hand edit (drift). The same
goes for `.gitignore` and `.dockerignore`: add entries under `ignore.git` /
`ignore.docker` in `mklv.config.json`, because the weekly conformance run
rewrites both files from the universal templates plus those additions.

`mage` lists each target with the **first sentence** of its doc comment, so keep
that sentence short enough to fit one 80-column line (roughly 55 characters after
mage's indent) and put every detail in the sentences after it.

Scratch files — throwaway scripts, captured logs, profiles, diff dumps — go in
`./tmp` (gitignored), never in the system `/tmp`. Keeping them inside the
workspace avoids permission prompts for paths outside it. Create it on demand
with `mkdir -p tmp`.

`./tmp` is gitignored but it is still inside the module, so `mage ci:check` (which
walks `./...`) will build, vet, lint, and coverage-check any Go package left
there. Delete a scratch Go package before running the gate.

When editing Markdown (including `AGENTS.md` and skill docs), keep files
markdownlint-clean.

Card research (so you never need a throwaway grep/JSON script). These live under
the `tools` mage namespace, invoked with a colon (`mage tool:stub`):

- `mage tool:lookup "<name>"` — print every source card whose name contains the
  query, with set code, collector number, house/type/rarity, printed text, and a
  ready-made `card.Provenance(...)` call.
- `mage tool:missing` — list the source cards in a set not yet tagged by an
  implemented card (the cards still to implement). With no set chosen it opens an
  interactive ↑/↓ picker; pass `-set=<slug>` to name one directly (slugs match the
  files in `internal/cards/provenance` minus `.json`, e.g. `callofthearchons`).
- `mage tool:nextCard` — print the next unimplemented card whose stub still
  carries the `//go:build todo` constraint, in collector-number order — the card
  to build next. With no set chosen it opens the interactive ↑/↓ picker; pass
  `-set=<slug>` to name one directly. This is the driver of the `implement-cards`
  workflow: build the card it names, drop the build tag, and run it again.
- `mage tool:coverage` — per-source-set count of cards covered by an implemented
  card's provenance Ref. Pass `-new` (`mage tool:coverage -new`) to count only the
  cards a set introduces, excluding the ones it reprints from an earlier set.

- `mage tool:nodeUsage` — inventory the card-authoring facade: every exported name
  in `internal/card`, grouped by the category its declaration block documents,
  with how many card definitions use it (`CARDS`, rarest first) and how many other
  files reference it (`OTHER` — the facade's internals, engine, web, tools, tests),
  plus a summary of the whole facade. Use it to find the neighbours a new
  node should be shaped alongside, and to audit that a thinly-used node is built
  from reusable atoms. **Low usage is not a defect** — half the card pool is
  unimplemented — and **`CARDS` 0 with a non-zero `OTHER` is not dead code**: it is
  type surface, registry plumbing, or an enum family member a card never writes by
  name. Only both columns at zero means nothing reaches for it. Narrow with
  `-max=<n>` and `-category=<substring>`, e.g.
  `mage tool:nodeUsage -max=1 -category=damage`.

- `mage tool:gameSize` — report the in-memory `GameState` size (the cost of one
  undo snapshot), the number of implemented cards, and — after building a fresh
  wasm — the shipped web bundle's size raw and compressed (brotli and gzip, the
  levels `WebAssets` ships), split into the WASM bundle, the other assets, and
  the total.

- `mage tool:stub "<setSlug>"` — scaffold a build-excluded (`//go:build todo`) stub
  file for every unimplemented card in a set, each carrying the printed text and a
  TODO marker. Excluded stubs do not compile or register, so the card database and
  coverage stay honest until a card is actually implemented; to implement one,
  remove the build tag and write the real ability. It also (re)generates the set
  package's `0set.go`, cataloging the cards the set reprints from earlier sets so
  they join its deck-generation pool as full members (ADR 0021). See the
  `stub-cards` skill (`.agents/skills/stub-cards`) for the stubbing workflow and
  the `implement-cards` skill (`.agents/skills/implement-cards`) for building the
  stubs into real cards.

Run `mage -l` to see every target.

## Long shell commands may span multiple lines

A long `&&`/`||` command chain does not have to be one physical line. A line that
ends with `&&`, `||`, or `|` continues on the next line, so a multi-step command
can be written as, e.g.:

```sh
mage gen &&
  mage ci:fix &&
  mage ci:check
```

Prefer this layout over one unreadable line when chaining several build/test
steps.

## Navigate and rename through the Go language server, not text substitution

Go has first-class language tooling, and the editor's symbol tools are backed by
`gopls`. It understands the type graph, so it distinguishes a method from a
same-named field, ignores a word in a comment or a string literal, and updates
every reference across the module. Reach for it first:

- **Renaming any identifier** — use the rename-symbol tool. **Never** rename with
  `sed`/`perl` across the tree. Many of this repo's identifiers are ordinary
  English (`Creature`, `Card`, `Power`, `House`, `Source`, `Target`), so a
  bare-word substitution silently rewrites comments, log strings, and unrelated
  types. A past sweep corrupted `g.State.Battleline[0]` into `g.State.InPlay[0]`
  exactly this way.
- **Finding a symbol's uses** — use the list-code-usages tool rather than `grep`
  when you want real references. `grep` is still the right tool for a literal
  string, a text pattern, or a survey across non-Go files.

Two failure modes, both about **aiming** the tool rather than the tool itself:

- **An ambiguous locator renames the wrong symbol.** The rename tool resolves a
  position from the _first_ line matching the `lineContent` you pass, then renames
  whatever symbol sits there. A bare field declaration is often not unique — this
  file has had `\tCreature LocalID` in two different structs — so the rename lands
  on the wrong one, compiles, and passes every test. Pass a locator that is unique
  and names its owner (a call site like `AemberCaptured{Creature: id,` beats the
  declaration line), and **`git diff` after every rename** to confirm it hit the
  symbol you meant. A wrong-but-consistent rename is invisible to the gate.
- **A stale language-server view rejects a valid rename**, citing errors in a file
  that is no longer on disk (a test file whose symbols moved). Confirm with
  `go build ./...`; if the tree really builds, the rename is safe and the server
  is behind. Do not fall back to text substitution because of this — retry, or
  make the edit by hand at the few call sites the usages tool reports.

Related: a deleted **comment** line is invisible to build, vet, lint, and tests.
When you insert a function directly above an existing one, re-read the seam (or
check `git diff`) to confirm you did not swallow the first line of the next
declaration's doc comment.

## Multiple agents may be running

More than one agent can be working in this repo at the same time. If the build,
vet, or tests fail because of a change you did **not** make — an unfamiliar file,
a symbol you never touched, an in-progress edit that doesn't yet compile — assume
another agent is mid-change. Wait a little and try again rather than "fixing" or
reverting their work. Only act on failures that stem from your own changes.

Prefer **targeted `go test`** for your own work and save the full `mage ci:check`
for when you actually need the whole-tree gate. `mage ci:check` is the single most
contended command in a multi-agent run: it builds and lints everything, so it
catches every sibling's mid-edit as a failure that is not yours. Before you run
it, glance at `git status --short internal/engine internal/web` — if a shared
package is mid-edit, the gate will fail on their work, not yours. When it fails
only on files outside your change set, record it and move on; do not chase a red
gate you did not cause.

## Leave git alone

Do not perform or propose git operations. Do not stage, commit, amend, branch,
tag, push, or pull, and do not end a reply by offering to commit or asking
whether the work should be committed — the human handles all of that. The only
exception is when you are explicitly asked to run a specific git command.

Read-only inspection (`git status`, `git diff`, `git log`) is fine when you
genuinely need it to answer a question.

## Never delete a failing test to get to green

A test that fails after your change is evidence, not an obstacle. Do **not**
delete, skip, or weaken it to make `mage ci:check` pass. Assume the test is right
and your change is wrong until you can state, in the commit or your summary,
exactly which rule or ADR makes the old expectation incorrect. Only then rewrite
the test — and rewrite it to assert the new correct behavior, never to assert
nothing. The same goes for an assertion inside a test: dropping the assertion
that caught you is deleting the test.

If a test genuinely blocks a change you were asked to make and you are not sure
the old behavior is wrong, leave the test failing and say so.

## Style and design live in `docs/style-guide.md`

The repo's coding style — how to interpret a request, when to refactor, how to
design for composition, naming and KeyForge vernacular, safety, comments, and
formatting — lives in one place: [docs/style-guide.md](docs/style-guide.md).
Read it before writing or reshaping code. The load-bearing summary:

- Read every request as **idiomatic, composable Go that will keep being
  extended**; when a request is ambiguous, pick what a senior Go engineer would
  find easiest to build on, not the shortest path to green.
- **Implement the mechanic, not the card** — decompose fused effects,
  parameterize over enums, reuse the shared vocabularies (`Target`, events,
  strategies), and treat a one-off name as a smell.
- **One consumer is fine; un-extendable is not.** A node only one card uses is
  not a defect — about half the card pool is unimplemented and some cards are
  genuinely unique. The test is whether a _second_ card doing a very similar thing
  could extend it (a field, a `Strategy`, another enum value) or would force a
  rewrite. So build every node out of atoms even when it has one consumer: a
  threshold is a `Count` plus a comparison, never a hard-coded `>=`; a subject is
  a field, never a name prefix. `mage tool:nodeUsage` shows the facade by
  category with consumer counts.
- **Refactoring is welcome and preferred over working around code**; keep it
  focused, keep everything green (including 100% `internal/engine` coverage), and
  lean on the tests to catch regressions. To refactor or clean up a whole area
  rather than one call site, use the `refactor-sweep` skill
  (`.agents/skills/refactor-sweep`), which also writes each finding back into
  these docs so it does not grow back.

The sections below stay here because they are structural facts about _where
things go_, not style.

## Engine design patterns and constraints

The full design ideal for `internal/engine` — the patterns it is built from, the
two hard constraints that shape it, and the deliberate tradeoffs (with how to
handle each) — lives in `internal/engine/AGENTS.md`. Read it before reshaping an
engine seam. The load-bearing rules that affect how you add anything:

- **Effects are an Interpreter AST** (ADR 0006). Every `Effect` renders its own
  text (`Text()`) and carries itself out (`Resolve()`), so printed card text can
  never desync from behavior. A new mechanic is almost always a new node in
  `effect_<mechanic>.go`, not a new branch in the `Game` runtime.
- **Vary behavior with a Strategy that also renders its own text.** When behavior
  changes along an axis, model the axis as a small strategy — a `Chooser` (or its
  optional-capability interfaces `OptionChooser`/`Orderer`), a `Refinement`, a
  `Count`, or a `Condition` — each of which carries both its behavior and its text
  fragment. Reach for this before adding another `Target` field or a `bool`.
- **The `Resolver` port is segregated into role interfaces** (ADR 0008)
  (`StateReader`, `EconomyResolver`, `CreatureResolver`, `CombatResolver`,
  `ZoneResolver`, `TurnResolver`, `ChoiceResolver`, `Logger`). A new engine
  capability is a method added to the role it belongs to, not to a flat list.
- **New whole-tree operations that are not part of a card's identity** (AI/MCTS
  scoring, static analysis, serialization) are a standalone type-switch function
  over `Effect` in one file — the Go-idiomatic Visitor — **never** a third method
  on every node. `Text()`/`Resolve()` are intrinsic; heuristics are not.
- **Flat, pointerless, comparable state is non-negotiable** (ADR 0005). State
  holds no closures, so "do X later" is flat enum-tagged data (the lasting registry), and
  values compared against state (e.g. `Target`) stay comparable — which is why
  `Target` is a flag struct with paired `x`/`hasX` fields rather than slices or
  pointers. Do not introduce a pointer/slice/map into `GameState` or a
  state-compared value type.

## File organization

`internal/engine` splits by responsibility, and the filename says which:

- `game.go` + `game_*.go` — the `Game` runtime: methods on `*Game`, grouped by
  area (`game_turn.go`, `game_play.go`, `game_read.go`, `game_abilities.go`,
  `game_combat.go`, `game_leaves_play.go`, …). A new `Game` method goes in the
  matching `game_*.go` — except a keyword mechanic's own gate methods, which live
  with the keyword (below).
- `effect.go` + `effect_*.go` — the effect AST, one file per mechanic
  (`effect_aember.go`, `effect_damage.go`, …). A new effect goes in
  `effect_<mechanic>.go`. **A mechanic's own gates and reads stay in its file**, so
  the mechanic reads as one unit even when some of it is a `*Game` method rather
  than an effect node: Alpha's `barredByAlpha`, Omega's `endStepIfOmega`, Deploy's
  `deployPosition`/`chooseFlank`/`choosePosition`, copy-stats' `copiedStatsSource`
  and `CopyStats`, the text box's `grantedTextBoxSources` and `GrantTextBox`. A
  `*Game` method that is _not_ part of one mechanic's seam goes in the matching
  `game_*.go` instead (`SetActiveHouse` lives in `game_turn.go`).
- `mechanic_*.go` — a self-contained ambient game mechanic that is neither an
  effect node nor a `Game`-method area: a small flat state type plus the reads and
  gates that govern it (`mechanic_tide.go`, `mechanic_toll.go`,
  `mechanic_bonus.go`, `mechanic_turnstat.go`, `mechanic_counter.go`,
  `mechanic_forge_guard.go`, `mechanic_spend_as_pool.go`). Reach for this so a
  mechanic clusters here instead of sitting alone as an orphaned concept file. A
  mechanic already clustered by a family prefix (`house_*`) stays there, and
  execution-model plumbing (`suspend.go`, ADR 0040) is not a game mechanic and
  keeps its own concept file.
- Everything else is a small type/data file named after the concept it defines:
  `card.go`, `state.go`, `types.go`, `target.go`, `duration.go`, `destination.go`,
  `resolver.go`, `text.go`. A concept file is a pure value/enum/display type with
  no rules of its own (a mechanic that carries rules is a `mechanic_*.go`, above).
  A new enum or value type gets its **own** `<concept>.go` — e.g. destinations live
  in `destination.go`, not `game_destination.go`; the `game_` prefix is only for
  `Game` methods.

Tests live beside their source in `<name>_test.go`, named after the **source file
that defines the symbols they exercise** — not after the card or scenario that
motivated them. A test for `TakeControl` (defined in `effect_control.go`) goes in
`effect_control_test.go`; do not leave it in a card-named `effect_saurian_egg_test.go`
or a scenario-named `deals_no_damage_when_attacked_test.go`. When a would-be test
file spans two mechanics, split its functions across the two mechanics' test files
rather than inventing a third name. The only test file not named after a source is
`helpers_test.go`, the shared test infrastructure. The `card` facade mirrors the
engine's namespace/type files: `duration.go` → `card/duration.go`,
`destination.go` → `card/destination.go`.

## Lasting "remainder of the turn" effects (event-driven, never hardcoded)

Some effects last "for the remainder of the turn" and attach to a later game event
— Full Moon gains Æmber whenever you play a creature, Dimension Door makes reaping
steal instead of gain. Do **not** add a bespoke `if g.State.Foo…` block to the play
or reap path (`PlayCreature`, `reapWith`) for each of these — that hardcodes every
effect into the hot path and makes effects sharing a timing window impossible to
order. Everything routes through the flat registry in `game_lasting.go`, as either
a **reaction** to an event or a **replacement** of its outcome (ADR 0007).

The same rule covers a card in play that says every creature gains an ability: it
is a `ConstantAbility` with a `Target` and `Granted` abilities, never a global
override in `discardDestroyed` or another leave-play path.

How to author either one is in
[docs/card-implementation.md](docs/card-implementation.md); the engine-side
registry and what adding an event costs are in `internal/engine/AGENTS.md`.

## KeyForge vernacular

Names must stay within KeyForge's own vocabulary, not generic gaming terms —
`MostPowerful`, not `Strongest`; `AemberCannotBeStolen`, not
`AemberTheftImmune` (theft and immunity are not KeyForge words). Use `cannot` for
a standing restriction or immunity, `absorbs` for a resource that is **spent**
stopping something (armor absorbing damage, a ward absorbing damage or a
destruction), and reserve `prevent` for a standing effect that refuses an outcome
without being used up. The full sourcing order (provenance
files → existing implementations → closest KeyForge phrasing) is in the naming
section of [docs/style-guide.md](docs/style-guide.md).

## Writing abilities (card authoring)

The engine's card-implementation vocabulary — every effect node, target and
filter, condition, count, trigger, duration, and card-level option, with how to
look one up — is cataloged in
[docs/card-implementation.md](docs/card-implementation.md). Read it before
concluding a card needs a new mechanic.

Author cards through the `card` facade in the multiline ability style —
`card.WithAbility(` breaks onto its own line, the trigger and effect share the
next line, and each effect struct field goes on its own line when there is more
than one (single-field effects stay inline):

```go
card.WithAbility(
  card.Trigger.Play, card.CannotFight{
    Player:   card.Opponent,
    Duration: card.Duration.NextTurn,
  }),
```

See `internal/cards/AGENTS.md` for the full card-authoring guide (file layout,
generated doc comments, wording rules, and tests). Every card's printed text must
obey the curated wording conventions in
[docs/card-wording-rules.md](docs/card-wording-rules.md) — most load-bearing:
`Sacrifice <self>` is written `Destroy <self>`, and an **`Omni:` ability is
authored as `Versatile` (a keyword on the card) plus a `Trigger.Action` ability**,
never a bespoke Omni trigger (the engine has none; ADR 0009).

## Rules voice and rulebook maintenance

Every player-facing surface — printed card text, the game log, the rulebook page
(`/rulebook`, rendered from the engine's typed rulebook term registry), and the
prose that frames it — is written in one
controlled **Rules voice** (ADR 0019): short declarative sentences, one
instruction per sentence, controlled vocabulary (one word per meaning), no
flourish. Card text follows the wording conventions in
[docs/card-wording-rules.md](docs/card-wording-rules.md); rulebook rules are plain
declarative sentences; bound examples are written Given/When/Then. The resource is
spelled **Æmber** everywhere. New and edited code comments follow the same plain
style (ADR 0020); the bulk comment migration is separate.

Three standing invariants:

- **Keep the rules current.** When you add or change a keyword, trigger, card
  type, or rule-bearing effect, register or update its rulebook term in the same
  change. The rulebook is a typed registry complete by construction (ADR 0018), so
  an undescribed keyword fails the build and a stale committed rulebook fails
  `mage ci:check`.
- **Comment and implementation must agree.** A doc comment says what its code
  does, never what it no longer does (ADR 0006). Where the behavior is non-obvious,
  bind the claim to an example (a cited engine test/scenario) rather than trusting
  prose.
- **Reference the KeyForge Master Rulebook, follow the Vex voice.** The
  KeyForge Master Rulebook ([docs/keyforge-master-rulebook.md](docs/keyforge-master-rulebook.md))
  is the wording authority **only for what Vex has not already decided**. Where
  Vex deliberately diverges from KeyForge, **Vex wins**, and the divergence
  is recorded in the Vex⇄KeyForge divergence register
  ([docs/keyforge-divergences.md](docs/keyforge-divergences.md)) — never silently
  overwritten by a later "match KeyForge".
- **Template consistency and simplicity across Vex outrank exact KeyForge
  wording.** Matching the printed KeyForge text is not a goal in itself and is
  never a blocker. A **meaning-preserving** reword that makes a card read in the
  same voice and template as the rest of Vex — or that lets a mechanic
  decompose into shared nodes instead of a bespoke one — is **welcome**, not a
  reluctant exception: prefer it whenever it simplifies the templating or makes
  card text more consistent. Such a reword changes only phrasing, so it needs no
  divergence-register entry (the register catalogs changes to a **rule or a name**,
  not house-voice alignment). Only when a reword actually changes a rule or a name
  does it belong in the register.

## Speak the lingo

[CONTEXT.md](CONTEXT.md) is the glossary: the game, cards and sets, deck
generation, scoring, and the client's UI areas. Prefer its term over the generic
one, and when a term is missing add it there rather than coining a second name.
When a reply uses a term the reader may not have met, name it plainly — the point
is that both sides end up speaking the same dialect. A sampler of what that
sounds like:

- **The game** — a one-shot card is a **Tactic**, not an "action card"; Æmber
  sits in a **pool**, is **captured**, **stolen**, or **exalted**, never "spent
  as mana"; a creature **reaps**, **fights**, or is **used**, and doing so
  **exhausts** it (not "taps"); cards sit in a **battleline** whose ends are its
  **flanks**; removal is **destroy**, **purge**, or **put into hand/discard**,
  never "exile" or "bounce"; you **forge a key**, you do not "score".
- **The client** — **player bar**, **play zone**, **board row**, **row label**,
  **midline**, **hand row**, **zone counts**, **sidebar**, **game log**,
  **turn HUD**, **prompt**, **action bar**, **zone viewer**, **card preview**.
- **The engine** — a mechanic is an **effect node** in the **effect AST**; an
  axis of behavior is a **Strategy** (a `Chooser`, `Refinement`, `Count`, or
  `Condition`); the `Resolver` is a **port** split into **role interfaces**; a
  "rest of the turn" effect is a **lasting effect**, either a **reaction** to an
  event or a **replacement** of its outcome.
- **The wider craft** — the effect tree is the **Interpreter** pattern, and a
  whole-tree pass over it is a **Visitor**; splitting `Resolver` by role is the
  **interface segregation principle**; `card.New(name, …, WithFoo)` is a
  **functional-options** builder; the bot is **MCTS** (Monte Carlo Tree Search);
  flat comparable state makes `GameState` a **value type**, which is what lets
  undo be a **snapshot** rather than an inverse operation.

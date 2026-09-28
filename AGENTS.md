# AGENTS.md

Rules for every task in this repo. A directory's own `AGENTS.md` adds the rules
for working there.

## "Reminder" means write it down

A request saying **"Reminder"** (or "remember this", "document this") asks for
the rule to be recorded where future work finds it, as well as applied: the
relevant `AGENTS.md`, a `docs/` page, an ADR, or a comment on the seam it
governs. A fix without the write-down is half the task. Record a rule in the
narrowest file every task that needs it reads; this file takes only what nearly
every task needs.

## An explained mechanic becomes a test

When the human explains a nuanced game mechanic (how ward meets the destroyed
tag, when a "Destroyed:" ability still reaches a card, which of two timings
wins), the turn is not done until a test pins it, even when the code already
behaves correctly, so the next agent cannot "simplify" the nuance away. A
passing new test is a success.

- Add the missing documentation too: a short comment on the seam the rule
  governs (the funnel, not a caller) that cites the test by name, plus a
  rulebook term or `docs/` line when the rule is player-facing. Write no
  markdown file just to record the explanation.
- Where the explanation contradicts the code, the explanation wins: say so
  plainly and fix the code.

## No agent backlog in the tree

- Outstanding agent work is tracked outside the tree by whatever drives the
  session. Create no markdown backlog. `docs/todo.md` is the human's list: read
  it if useful, never write to it.
- A recorded work item is a handoff to an agent who was not in the
  conversation: record what was decided and why (behavior, cards affected,
  expected text). When you pick one up, the recorded decision wins over a
  contradicting code comment; fix the comment rather than re-litigate the item.
- Decided work with no consuming card in an implemented set yet goes in
  [docs/todo-future-set.md](docs/todo-future-set.md), keyed to the set that
  first needs it. When you implement or stub a set, scan it for that set and
  build each primitive with its first real consumer.
- **You sequence approved work; the human does not.** End no turn asking which
  item comes first or offering a menu of plans. Order by dependency (a shared
  primitive, renderer helper, or enum fold before its consumers), then easiest
  win within a tier. State the order in one line and start.

## Checks, configs, and scratch files

- The full gate (`mage ci:fix`, then `mage ci:check`, which must print
  `ALL GREEN`) runs when the session finishes. Run only the narrowest check that
  answers your question:
  - `go test ./internal/<pkg>/ -run '<Pattern>'` for one package;
    `mage testRun '<Pattern>'` for matching tests in every package.
  - `go build ./...` or `go vet ./internal/<pkg>/` to compile-check.
  - Coverage: `internal/engine`, `internal/cards/sets` (counted over the
    `./internal/cards/...` tests), `internal/cards/cardtest`, and
    `internal/deckgen` are held at 100%; `internal/web` is ungated. Check one
    with `go test -coverprofile=tmp/cover.out ./internal/engine/` and
    `go tool cover -func=tmp/cover.out`.
- Tool configs are generated; never hand-edit `.golangci.yaml`,
  `.markdownlint-cli2.yaml`, `.commitlint.yaml`, `.gitleaks.toml`,
  `.ruleguard.go`, `.gitignore`, or `.dockerignore`. Change a lint setting or
  exclusion under `tools` in `mklv.config.json`, and an ignore entry under
  `ignore.git` / `ignore.docker`; `ci:fix` regenerates the files and `ci:check`
  fails on a hand edit.
- Scratch files (throwaway scripts, logs, profiles, diff dumps) go in `./tmp`
  (gitignored; `mkdir -p tmp`), never the system `/tmp`. The gate walks
  `./...`, so delete any scratch Go package in `./tmp` before you finish.
- Break a long `&&` / `||` / `|` chain after the operator, one step per line.
- Before editing any Markdown, read [docs/AGENTS.md](docs/AGENTS.md).
- When the gate reports a simulator invariant violation, read
  [docs/sim-debugging.md](docs/sim-debugging.md) before debugging it.
- Every mage target is described in [docs/building.md](docs/building.md).

## Navigate and rename through the Go language server

The editor's symbol tools are backed by `gopls`, which knows the type graph.

- Rename every identifier with the rename-symbol tool, never `sed`/`perl`
  across the tree. Many identifiers are English words (`Creature`, `Card`,
  `Power`, `House`, `Source`, `Target`); a bare-word sweep once turned
  `g.State.Battleline[0]` into `g.State.InPlay[0]`.
- Find a symbol's real uses with the list-code-usages tool. `grep` is for a
  literal string, a text pattern, or non-Go files.
- Aim the rename. The tool renames whatever sits on the first line matching the
  `lineContent` you pass, and a bare field declaration (`\tCreature LocalID`)
  can occur in two structs. Pass a unique locator that names its owner (a call
  site like `AemberCaptured{Creature: id,`), and `git diff` after every rename:
  a wrong-but-consistent rename passes the gate.
- When the server rejects a valid rename citing a file no longer on disk,
  confirm with `go build ./...`; if the tree builds, retry, or edit by hand the
  call sites the usages tool reports. Never fall back to text substitution.
- A deleted comment line passes build, vet, lint, and tests. After inserting a
  function above another, check the seam (or `git diff`) for the next
  declaration's doc comment.

## Leave git alone

Do not stage, commit, amend, branch, tag, push, or pull, and do not propose git
operations or offer to commit; the human handles all of it. The one exception
is an explicit request to run a specific git command. Read-only `git status`,
`git diff`, and `git log` are fine when you need them.

## Never delete a failing test to get to green

A test that fails after your change is evidence, not an obstacle. Do **not**
delete, skip, or weaken it to pass the gate. Assume the test is right and your
change is wrong until you can state, in the commit or your summary, exactly
which rule or ADR makes the old expectation incorrect; only then rewrite it, to
assert the new correct behavior, never to assert nothing. Dropping the
assertion that caught you is deleting the test. If a test blocks a change you
were asked to make and you are not sure the old behavior is wrong, leave it
failing and say so.

## Style and design: `docs/style-guide.md`

Read [docs/style-guide.md](docs/style-guide.md) before writing or reshaping
code. The load-bearing summary:

- Read every request as **idiomatic, composable Go that will keep being
  extended**; when it is ambiguous, pick what a senior Go engineer would find
  easiest to build on, not the shortest path to green.
- **Implement the mechanic, not the card**: decompose fused effects,
  parameterize over enums, reuse the shared vocabularies (`Target`, events,
  strategies), and treat a one-off name as a smell.
- **One consumer is fine; un-extendable is not.** About half the card pool is
  unimplemented, so a node one card uses is not a defect. The test is whether a
  _second_ similar card could extend it (a field, a `Strategy`, another enum
  value) or would force a rewrite. Build every node from atoms: a threshold is
  a `Count` plus a comparison, never a hard-coded `>=`; a subject is a field,
  never a name prefix.
- **Refactoring is welcome and preferred over working around code.** Keep it
  focused and green, and lean on the tests. To refactor a whole area rather
  than one call site, use the `refactor-sweep` skill
  (`.agents/skills/refactor-sweep`).

## KeyForge vernacular

Names stay in KeyForge's own vocabulary: `MostPowerful`, not `Strongest`;
`AemberCannotBeStolen`, not `AemberTheftImmune`. Use `cannot` for a standing
restriction or immunity, `absorbs` for a resource **spent** stopping something
(armor absorbing damage, a ward absorbing damage or a destruction), and
`prevent` only for a standing effect that refuses an outcome without being used
up. The sourcing order for a new word is in the naming section of the style
guide.

## Rules voice

Every player-facing surface — card text, the game log, the `/rulebook` page and
the prose around it — is written in one controlled **Rules voice** (ADR 0019):
short declarative sentences, one instruction per sentence, one word per
meaning, no flourish. Card text follows
[docs/card-wording-rules.md](docs/card-wording-rules.md); bound examples are
Given/When/Then. Spell the resource **Æmber** everywhere.

- New and edited code comments follow the same plain style (ADR 0020), and a
  doc comment says what its code does, never what it no longer does
  (ADR 0006). Bind a non-obvious claim to a cited engine test or scenario.
- The KeyForge Master Rulebook
  ([docs/keyforge-master-rulebook.md](docs/keyforge-master-rulebook.md), a
  converted PDF) is the authority **only for what Vex has not decided** — in
  the implementation, the `/rulebook` registry, or these docs. Where Vex
  deliberately diverges, Vex wins, and the divergence is recorded in
  [docs/keyforge-divergences.md](docs/keyforge-divergences.md), never silently
  overwritten by a later "match KeyForge".
- Consistency and simplicity across Vex outrank exact KeyForge wording. A
  meaning-preserving reword that fits the Vex template, or lets a mechanic
  decompose into shared nodes, is welcome and needs no divergence entry; the
  register records only a changed rule or name.

## Speak the lingo

[CONTEXT.md](CONTEXT.md) is the glossary: the game, cards and sets, deck
generation, scoring, and the client's UI areas. Use its term over the generic
one; add a missing term there rather than coin a second name. Name a term
plainly when the reader may not have met it. In the game, a one-shot card is a
**Tactic** (not an "action card"); Æmber sits in a **pool** and is
**captured**, **stolen**, or **exalted** (never "spent as mana"); a creature
**reaps**, **fights**, or is **used**, which **exhausts** it (not "taps");
cards sit in a **battleline** whose ends are **flanks**; removal is
**destroy**, **purge**, or **put into hand/discard** (never "exile" or
"bounce"); you **forge a key**, you do not "score".

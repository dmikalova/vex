---
name: refactor-sweep
description: Sweep a bounded area of this repo for refactors and cleanup — decompose fused effects, consolidate near-duplicates, comment what the names do not carry, split oversized files, regroup misplaced code, hunt bugs and anti-patterns — then ratchet each finding into the docs so it cannot come back. Use when the user wants code refactored, cleaned up, tidied, simplified, audited, or reorganized.
---

A **sweep** takes one bounded area, collects every **finding** in it, fixes them
cheapest-first, and then **ratchets** each one: writes the rule that made it a
finding into the doc that governs that area, so the next agent does not
reintroduce it. A sweep that only fixes code is half a sweep — the same shapes
grow back within a few rounds of card work.

Read the root `AGENTS.md` and `docs/style-guide.md` before starting; read
`internal/engine/AGENTS.md`, `internal/cards/AGENTS.md`, or
`internal/web/AGENTS.md` for whichever area the sweep covers. Those four files
plus `docs/style-guide.md` are both the standard you judge against and the place
most ratchets land.

## 1. Bound the sweep

Restate the area before touching code: a package (`internal/engine`), a family
(`effect_*.go`), a single file, or a class of finding across the repo ("every
struct field that needs a comment"). An unbounded "clean things up" means the
package the user was last working in.

Other agents work this tree at the same time. Sweep only files you can see are
settled — if `mage check` fails on a symbol you never touched, that is someone
else mid-change, so leave it, say so, and sweep elsewhere.

## 2. Survey before deciding

Gather the candidate list first; picking findings by memory finds only the file
you just read.

```sh
mage check                                   # the baseline: know what was already red
mage lint                                    # golangci-lint: unused code, shadowing, staticcheck
wc -l $(git ls-files '<area>/*.go' | grep -v _test) | sort -rn | head -20
grep -rn 'fmt\.Print\|println(\|TODO\|FIXME\|XXX' <area> --include='*.go' | grep -v _test
git status --porcelain                       # stray files another agent has not committed
```

Undocumented struct fields (a candidate list, not a to-do list — comment only the
ones whose name does not carry them):

```sh
awk '/^type .* struct \{/{s=1;next} /^\}/{s=0} s && /^\t[A-Za-z]/ && prev !~ /^\t*\/\// {print FILENAME":"FNR": "$0} {prev=$0}' $(git ls-files '<area>/*.go' | grep -v _test)
```

Then read the area's files end to end. The greps find debris; the findings that
matter — a fused effect, a rule reimplemented in the client, a file that has
become two files — only show up in a read.

## 3. What counts as a finding

Judge against `docs/style-guide.md` and the area's `AGENTS.md`. Each class below
names what to hunt; all of them are worth a pass in any sweep.

**Fusion.** One node doing what composition should. "A, then B if C" belongs in
the card as `Sequence{A, Conditional{C, B}}`, with values threaded through
`EffectContext` (`ctx.It`, `ctx.ChosenHouse`, `ctx.Produced.*`) — not inside a
bespoke effect. A node named after a card rather than a mechanic is the loudest
tell.

**Duplication.** Two effects differing by a constant or a side are one effect
parameterized over an enum; two conditions asking the same question of different
subjects are one condition with a `Player`; two card files repeating a five-line
ability want a shared composite. Retire the shape you replaced and re-express its
callers in the same commit.

**Ladder violations.** A change belongs at the cheapest rung that can carry it: a
field or Strategy on an existing effect (a `Count`, `Selector`, `Condition`,
`Chooser`) beats a new node; a new node beats a new `Resolver` capability; that
beats new state. A new node that only varies an existing one along one axis is a
Strategy wearing a node's clothes.

**Misplacement.** `internal/engine`'s filenames are its index: `game_*.go` is
`*Game` methods only, `effect_<mechanic>.go` is one mechanic, and any other
concept gets its own `<concept>.go`. A `Game` method in an `effect_` file, an
enum living inside `game_`, a `Resolver` method on the wrong role interface, and
a rule reimplemented in `internal/web` instead of read from the engine are all
findings.

**Oversized files.** A file past ~250 lines is a prompt to look, not a defect. It
earns a split only when it holds two concepts that a reader would look for
separately; split along that seam and name each half for its concept. Test files
split the same way as their source.

**Missing comments.** The repo's rule is one short line stating what the code
cannot show on its own. A finding is a type, field, or function whose _why_,
unit, invariant, or zero-value meaning is unstated — not one whose name already
says it. Rewrite a comment that restates its next line rather than adding another.

**Bugs and loopholes.** Read for the ones tests miss: a zero value that is a legal
value and so cannot signal "unset" (ADR 0010 — validate at init, pair a value with
`hasX`), an unchecked bound or index, a `Target` compared against state that has
stopped being comparable, an effect whose `Text()` no longer matches what
`Resolve()` does, a lasting effect that clears on the wrong turn, an error
swallowed instead of returned.

**Dead weight.** Unused exports, a helper with one caller that should be inlined,
a stub file left after the real implementation landed, a debug print, an orphan
file with no obvious owner.

## 4. Fix cheapest-first, and keep the suite honest

Order the findings by blast radius and take the cheap ones first: comments and
deletions, then in-file moves, then file splits, then reshaping an effect, then
touching a seam (`Resolver`, `GameState`, the `card` facade) last. Land each one
green rather than batching a dozen unverified edits.

The card and engine tests pin both behavior and rendered text, so refactor
freely — but a test that goes red is evidence. Never delete, skip, or weaken one
to reach green; either fix your change, or state the rule that makes the old
expectation wrong before rewriting the test to assert the new correct behavior.

A wording change means editing the effect's `Text()` in
`internal/engine/effect_*.go` and re-running `mage generateComments`; card doc
comments are generated and hand-edits are overwritten.

## 5. Ratchet every finding

For each finding, ask what would stop it coming back, and act on the answer.
Findings cluster: three instances of one shape mean the rule was never written
down, and writing it is worth more than the three fixes.

- A rule about **style, naming, or composition** → `docs/style-guide.md`.
- A rule about **where code goes**, or a structural fact about the repo → the
  root `AGENTS.md`.
- A rule about an **engine seam** (effect AST, Strategy, `Resolver` roles, flat
  state) → `internal/engine/AGENTS.md`.
- A rule about **authoring a card or its test** → `internal/cards/AGENTS.md`.
- A rule about **the client** (prompts, handlers, snapshots, CSS) →
  `internal/web/AGENTS.md`.
- A rule about **printed card text** → `docs/card-wording-rules.md`.
- A **decision with real tradeoffs** you had to weigh → a new `docs/adr/NNNN-*.md`,
  with the alternatives you rejected and why. Point at it from the AGENTS file
  that governs the area.
- A **new term** you found yourself explaining → `CONTEXT.md`.

Write the positive rule ("author an Omni ability as `Versatile` plus a
`Trigger.Action`"), not the ban. Prune while you are in there: a sweep that
corrects a doc should also delete the line the correction made stale, and a rule
that now contradicts the code is itself a finding.

A one-off with no general rule behind it ratchets to nothing. Say so and move on
rather than inventing a rule to have written one.

## 6. Close green and report

```sh
mage gen && mage check      # gen = card comments + rulebook; check must print ALL GREEN
```

`mage check` includes the 100% `internal/engine` coverage gate: new engine code
needs its test in the matching `internal/engine/effect_*_test.go` or
`game_*_test.go`.

Report the sweep as a table of findings — what it was, where, what you did, and
what you ratcheted (or why nothing) — so the ratchets are reviewable as a set.
Then hand back, or bound the next area.

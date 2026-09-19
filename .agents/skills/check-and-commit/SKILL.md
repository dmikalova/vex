---
name: check-and-commit
description: Get `mage check` fully green, then stage, commit, and push the work in one go. Use when the user wants to finish up by verifying the gate and committing ("check and commit", "/check-and-commit", "get it green and push"). Assumes this is the only agent running in the repo.
---

This skill takes the repo from "work in progress" to "pushed", in order:

1. Get `mage check` to print `ALL GREEN`.
2. Stage everything, commit, and push.

**This skill is the one exception to the repo's "leave git alone" rule.** The
user has explicitly asked for the commit and push, so run them — but only as the
final step, only after the gate is green, and only within this skill.

## Assume you are the only agent running

Every safeguard the repo carries for concurrent agents is off here. A build,
vet, lint, or test failure is **yours** — there is no sibling mid-edit to blame,
so do not wave a red gate off as someone else's work. Fix it. An unfamiliar file
or an unexpected diff is part of the change set you are about to commit, so
account for it rather than working around it.

## 1. Get `mage check` green

Run the full gate:

```sh
mage check
```

It formats in place, then builds, vets, lints, markdown-lints, tests, and
coverage-checks the whole tree (all four gated areas must stay at 100%). It must
print `ALL GREEN`.

If it fails, fix the cause and run it again. Repeat until it is green. Do **not**
advance to the commit while anything is red, and do **not** reach for the
escape hatches the repo forbids: never delete or weaken a failing test to get to
green, and never bypass a check (`--no-verify` and the like). A red gate is a
real problem to solve, not an obstacle to route around.

Common fixes:

- Formatting drift — `mage check` formats in place, so a second run often clears
  a pure-format failure. Re-run before hand-editing.
- Stale generated output — run `mage gen` if a card comment or the rulebook is
  out of date, then re-run the gate.
- Coverage below 100% in a gated area — add the missing test for the new code.
  See "Closing a coverage gap" below for the exact commands to find the red line.

### Closing a coverage gap

`mage cover` reports each gated area's total and, when one is below 100%, names
every function still short with its percentage. To turn that into "the exact
uncovered line", generate a profile and read the zero-count blocks — do **not**
guess which branch is missing:

1. **Get the per-function shortfall.** `mage cover` is authoritative and prints
   the offenders. To iterate faster on one area, profile it directly. The engine
   coverage gate runs under the `assert` build tag, so match it or the numbers
   drift:

   ```sh
   cd internal/engine &&
     go test -tags assert -coverprofile=tmp/eng.cov . &&
     go tool cover -func=tmp/eng.cov | grep -v '100.0%'
   ```

   (Use `tmp/` inside the repo, not the system `/tmp` — but delete the scratch
   profile before the gate, since `mage check` walks the whole tree.)

2. **Find the uncovered line, not the function.** The `-func` view gives a
   percentage; the profile itself gives the line. Each block line is
   `file:startLine.col,endLine.col numStmts count` — a trailing `0` is an
   unexecuted block. Filter to the function's line range:

   ```sh
   awk -F: '$0 ~ /text.go/ && $2>=837 && $2<=913' tmp/eng.cov | grep ' 0$'
   ```

   Then open that line and see which branch it is — an untaken `if`, a `default`
   arm, a `break` when candidates run out before a max, a flag field never set
   true. Write one focused test that drives exactly that branch.

3. **Confirm and clean up.** Re-run the `-func` filter to see the function hit
   100%, then delete the scratch profile so it does not linger in the tree:

   ```sh
   go tool cover -func=tmp/eng.cov | grep '<funcName>' && rm -f tmp/eng.cov
   cd - >/dev/null && mage cover   # authoritative: all four areas at 100.0%
   ```

### Diagnosing a sim / `FuzzPlay` failure

When the gate fails inside `internal/sim` (a `FuzzPlay/<seed>` case reporting an
`invariant violated ...` line), the failure is a whole simulated game, not a
single unit. Work it like this:

1. **Find the failing seed and its script.** Run the sim tests alone and read
   the failure — it prints the seed and the long `script:` hex:

   ```sh
   go test ./internal/sim/ 2>&1 | tail -30
   ```

2. **Replay the script with the game log.** Stash the hex in a shell var (it is
   long) and replay it. `mage debug` prints **both players' full deck lists**
   above the log tail, which is the whole point — the invariant names the victim,
   not the culprit:

   ```sh
   S=<script-hex-from-the-failure>
   mage debug -script=$S -tail=60
   ```

   Widen `-tail`, or dump the whole game to a file to read it end to end:

   ```sh
   mage debug -script=$S -tail=100000 > tmp/sim/bug.txt 2>&1
   grep -nE 'deck \(|VIOLATION' tmp/sim/bug.txt   # deck lists + the violation
   ```

3. **Suspect the mechanic, then find the card that carries it.** Read _both_
   decks, not only the cards in the log tail. Scan every card for the mechanic
   that could produce the bad state — a power reducer for a 0-power creature, a
   blanker for a creature that lost its ability, an attachment for a stat that
   drifted. A creature that dies (or fails to die) "for no reason" almost always
   lost or gained a constant a card that has since left the tail was granting.
   Grep the log for a suspect by name to see its whole lifecycle:

   ```sh
   grep -niE 'King of the Crag|Bingle|Shadow of Dis' tmp/sim/bug.txt
   ```

4. **Fix the mechanic, not the seed.** Never delete or weaken the invariant or
   the seed to get green. Once you find the root cause, add a focused engine or
   card regression test that fails before the fix and passes after (prove it:
   run the new test, revert the fix, watch it fail, restore the fix), then
   confirm the seed passes:

   ```sh
   go test ./internal/sim/ 2>&1 | tail -5
   ```

## 2. Stage, commit, and push

Only once `mage check` prints `ALL GREEN`, run the commit and push exactly as
the user asked:

```sh
git add -A && git commit -a -m "feat: implement cards" && git push
```

Then report the pushed commit and stop. Do not amend, force-push, rebase, or
touch history — this skill's entire git footprint is the single add-commit-push
above.

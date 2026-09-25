---
name: implement-cards
description: Work through a KeyForge set's unimplemented cards in this repo until the amount asked is done. Use when the user wants to implement a set's backlog ("implement more Call of the Archons cards", "keep going", "iterate until the set is done"), or triage what is easy.
---

The backlog is worked **one card at a time**, driven by `mage tool:nextCard`.
Ask the tool for the next unimplemented card, build it, drop its build tag, and
ask again. Each card leaves the engine slightly richer and the set measurably
more covered. The loop repeats until the stop condition is met.

When several stubbed cards share a mechanic, you may jot them grouped in
[docs/todo-agent.md](../../../docs/todo-agent.md) and build the shared primitive
once, then knock out the group — a planning overlay on the collector-number
`nextCard` order, not a replacement for it. Delete each entry as its card lands (a
done item is removed, never marked done).

The goal is **cards implemented**. Keep moving through the backlog: prefer
implementing the next card to polishing the last one. This is **one continuous
run through the entire stop condition** — never frame it to the user as a
"multi-session grind", a "first batch", or work that will be "continued later".When the user says "do all the cards", they mean do all the cards now, in this
run, without handing back partway. Do not propose stopping, do not ask whether to
keep going, and do not offer a status report as a substitute for finishing. Tests
matter, and you should get your own work passing; but if a card is a reasonable
best-effort implementation and the only failures come from **another agent's**
in-progress changes (an unfamiliar file, a symbol you never touched, a gate that
was green before you started), note it and keep going rather than stalling on
someone else's edit.

**Do not call `task_complete`, and do not end your turn, until the stop condition
below is met.** A card that compiles, a green `mage ci:check`, a reconnaissance
conclusion that "the rest all need new mechanics", and a written status summary
are ALL checkpoints, never the finish — reaching one means ask for the next card
and build it, not stop and report. "Every remaining card needs a new mechanic" is
the normal state of a backlog run, not a blocker: building the mechanic IS the
job, so keep building them one after another. The ONLY things that authorize
ending the run are: the stop condition is met, every remaining card is genuinely
blocked by something you cannot build (not merely "needs new work"), or the user
interrupts. After each card lands, the default and automatic next action is to run
`mage tool:nextCard` and start the next card — no pause, no check-in.

Read `internal/cards/AGENTS.md` (authoring + tests), `docs/card-implementation.md`
(the catalog of every engine capability a card can compose from),
`docs/card-wording-rules.md` (rendered-text rules), and the root `AGENTS.md`
(composability) before starting.

## 1. Fix the stop condition, then set up

Restate the stop condition before touching code, because it is the ONLY thing
that ends the run — nothing else authorizes `task_complete`.
It is one of: **the whole set**, **a count** ("ten more cards"), or a
**qualifier** ("all the Mars cards", "everything that isn't a lasting effect").
An unqualified "keep going" or "iterate" means the whole set. Hold the count
explicitly (in a session memory file) and check it after each card; until it is
reached, keep going.

The backlog must already be stubbed — every unimplemented card is a
`//go:build todo` file that `mage tool:nextCard` can hand you. If a set has not
been stubbed yet, run the **stub-cards** skill first.

**Scan [docs/todo-future-set.md](../../../docs/todo-future-set.md) before you
start.** It holds decided work parked against a future set — primitives with no
consumer in an implemented set yet. If an item names the set you are implementing
(or a card it introduces), fold it into the run and build the primitive alongside
its first real consumer, then delete the item when it lands.

Then, once per run:

```sh
mage tool:coverage             # per-set covered/total — the number the run moves
mage tool:nextCard -set=<slug> # the next //go:build todo card: its file path,
                               # stats, printed text, and card.Provenance(...) call
```

Set slugs match the files in `internal/cards/provenance/` minus `.json`. `nextCard`
walks the set's missing cards in collector-number order and stops at the first one
still carrying a `//go:build todo` stub — so implementing a card (dropping its
build tag) is what advances the tool to the next one. With no `-set` it opens an
interactive ↑/↓ picker.

**Keep a triage memo in session memory.** Record the running state so a resumed
run does not re-triage from scratch: which cards are done, and the mechanics you
have already built (with the cards each unblocked). When you notice, while reading
one card, that a later card wants the same mechanic, jot it there. Update it as
cards land.

## 2. The loop: next card, build, repeat

Run `mage tool:nextCard`, build the card it names, then run it again. For each
card:

1. **Decide whether it is easy or gated.** A card is **easy** when its whole text
   composes from capabilities that already exist. Check
   [docs/card-implementation.md](../../../docs/card-implementation.md) first — it
   catalogs every effect node, target and filter, refinement, condition, count,
   trigger, and card-level option by name, so a mechanic you have not met is found
   by reading rather than rebuilt. Then confirm a primitive's exact fields in the
   facade (`internal/card/effects.go`, `target.go`, `options.go`) or in a similar
   existing card. An easy card is built directly (_Implementing one card_ below).
2. **A gated card needs a mechanic that does not exist yet** — usually a new
   field, `Strategy`, target filter, count, refinement, condition, or cross-turn
   hook, which you build directly (_Building a mechanic_ below). A gate that would
   need a **brand-new effect node** is different: it is subject to the new-node
   grill gate (_Building a mechanic_), so do **not** add the node until it clears.
   Then implement the card on top of the mechanic.
3. **After a mechanic lands, cash it in.** Before returning to strict `nextCard`
   order, implement any other unimplemented card that the same mechanic now
   unblocks — that is what makes it a mechanic instead of a one-off. Use the
   triage memo (and a quick scan of the remaining stubs' printed text) to find
   them. Only when the mechanic is fully cashed in do you go back to `nextCard`.

**Verify with a targeted `go test`, not `mage ci:check`, inside the loop.** A full
gate run is slow and its `ALL GREEN` is a false finish line that invites stopping,
so `mage ci:check` is an **end-of-run** step, not a per-card one. Verify each card as
it lands by running `go test` for **exactly the card or mechanic you changed**:

```sh
go test ./internal/cards/sets/<slug>/ -run Test<Name>   # the card you just wrote
go test ./internal/engine/ -run Test<Mechanic>          # the mechanic it uses
```

Name the specific `-run` pattern — the card's `Test<Name>` and, for a gated card,
the engine `Test<Mechanic>` you added — so the run is a few seconds, not the whole
suite. Add `mage ci:build` when a change spans packages. Only the two targeted tests
matter per card; save `mage ci:check` for step 3.

### Building a mechanic

**A new effect node is a last resort, gated behind a grill.** Extending the
engine in reasonable, composable ways needs no ceremony — a new field or
`Strategy` (a `Chooser`, `Refinement`, `Count`, or `Condition`), a new `Target`
filter, or a new count on an existing effect is the normal way a set grows the
engine, and you just build it. But **introducing a brand-new `Effect` node in the
AST — or cramming a mechanic into an existing node in a way that is not clean and
composable — is forbidden until you have run a full grill-me session** (the
`grilling` skill) and the human has signed off. The engine's whole design is that
behaviour composes from a small vocabulary of self-rendering nodes (ADR 0006); a
new node widens that vocabulary permanently, so it must be argued for, not slipped
in. Reaching a genuine new-node need is one of the **authorized reasons to pause
the run** — unlike a status check-in, which is never allowed — so present the
grill and stop. Put the questions in the reply as end-of-turn plain text (the
`❓`/`➡️` convention in [docs/todo-agent.md](../../../docs/todo-agent.md)), not
through an interactive tool. The grill must put on the table:

- **Why it is necessary** — the mechanic the existing nodes genuinely cannot
  express, not merely a shape that would be more convenient as its own node.
- **The example cards** — the cluster of real cards (this set and future) that
  need it, with their printed text, so the node is shaped for the whole group.
- **Alternatives rejected** — what a new field, `Strategy`, filter, or count on an
  existing node would look like, and precisely why it does not work (composition
  lost, an illegal state made representable, wording it cannot render).
- **A before/after comparison** — how the affected cards author today (or the ugly
  cram) versus how they author on the proposed node, so the win is concrete.

If the grill does not clearly land in favour of the new node, extend an existing
one instead. When in doubt — when the extension feels like a cram rather than a
clean fit — stop and grill rather than pushing it in.

**A node with one consumer is fine; a node that cannot take a second one is not.**
About half the card pool is unimplemented, and some cards genuinely do a unique
thing, so "only one card uses it" is never on its own a defect. The test is:
_if another card did a very similar thing, could it extend this node — a field, a
`Strategy`, another enum value — or would the node have to be rewritten to be
composable?_ So write every node, even a single-consumer one, out of atoms: the
threshold is a `Count` plus a comparison, not a hard-coded `>=`; the subject is a
field, not a name prefix; the amount is a `Count`, not an `int` literal welded to
one card's number. A node named after a card, or one whose `Resolve` inlines a
comparison an existing `Count`/`Condition` already expresses, is the failure this
rule catches. `mage tool:nodeUsage` lists the facade by category with each node's
consumer count — use it to find the neighbours a new node should be shaped
alongside, not as a hit list of "unused" nodes to delete.

Once the node is justified (or you are extending an existing one cleanly):

1. **Shape it for the whole cluster, not the first card.** Name and shape the
   mechanic so every card that wants it can use it. Prefer the cheapest engine
   surface: a new field or Strategy on an existing effect is cheaper than a new
   effect; a new effect is cheaper than a new `Resolver` capability, which is
   cheaper than new state.
2. **Add its engine test** in the matching `internal/engine/effect_*_test.go` as
   you go; `mage ci:cover` gates `internal/engine` at 100%, and a card test does not
   count toward it.
3. **File the mechanic in the rulebook.** A player-facing mechanic is only
   finished when a player can look it up, so register a `RuleTerm` for it in the
   matching `internal/engine/ruleterms_<section>.go` — `effect`, `keyword`,
   `ability`, `cardtype`, `combat`, or `turn` — with a `Title` and a `Body` that
   is the entry's text. The `/rulebook` and `/glossary` pages render the registry
   live (ADR 0018), and the completeness test fails the build if a closed catalog
   (keyword, trigger, card type) has a member with no term, so an entry can never
   drift from the code it describes. `docs/keyforge-master-rulebook.md` is the
   guide to what belongs: if the official rulebook explains the term to a player,
   ours must too. Two code sites that are one rule share a `Title` (a `Subtitle`
   groups them beneath it). While you are in the file, add a term for any
   **existing** mechanic beside it that is missing one — an undocumented neighbour
   is a finding, not the status quo.
4. **Teach the client to play the mechanic**, if it needs anything new. A mechanic
   that asks the player a question the browser client cannot ask is only half
   built. Read `internal/web/AGENTS.md` and check the new mechanic against it: a
   new `Chooser` prompt shape needs a case in `game_chooser.go` and a prompt in
   `view_controls.go`; a new zone or card state needs to be drawn in
   `view_board.go`/`view_card.go`; a new player action needs a keyboard route in
   `game_lifecycle.go` and a Tab stop in `game_nav.go`. Never reimplement the rule
   in the client — ask the engine, and add the reader to `internal/engine` if it
   does not exist.

### Implementing one card

1. Delete the stub (`command rm <snake>.go` — `rm -f` is blocked by an alias) and
   write the real `internal/cards/sets/<slug>/<snake>.go`.
2. Seed it with a bare `// <Card Name>` comment above a
   `var Name = card.New(…)` declaration — the house, `card.Type`, rarity, a
   `card.Provenance(card.<Set>, n)`, then the `With*` options. The card TYPE
   "action" is `card.Type.Tactic` (wording rule 19). Follow the one-field-per-line
   struct style in `internal/cards/AGENTS.md`. When an ability names the card's
   **own** house, write `card.House.Self` rather than repeating the house — but a
   card naming a _different_ house (Take That, Smarty Pants is about Logos
   creatures) spells that house out.
3. Write `<snake>_test.go` with the `ct.Play` harness — a `func Test<Name>` with
   `t.Run` subtests. A sole target auto-resolves; with 2+ candidates answer via
   `h.P1.ClickCard(handle)` / `h.P1.ClickOption(name)`. Set up a damaged creature
   with `handle.Damaged(n)`; read chains via `h.Game().State.Chains[0]`. A
   recurring `Trigger.StartOfTurn` / `EndOfTurn` ability fires from in play, not
   on play: loop back to its owner's next turn with `h.P1.EndTurn()`, then
   `h.P2.ChooseHouse(other)` and `h.P2.EndTurn()` — the start-of-turn ability
   resolves during that turn hand-off, before choose-house, so a sole-target
   effect needs no click.
4. Run `mage generateComments` — it rewrites the card and test doc comments from
   the definition. Read the generated text against `docs/card-wording-rules.md`.
   A wording fix means changing the effect's `Text()` in
   `internal/engine/effect_*.go`, never hand-editing the comment.
5. If a card's printed text is deliberately reworded — to dodge a mechanic that
   is not worth the state it would cost, or to simplify — add the rule to
   `docs/card-wording-rules.md` as a numbered section, stated as a **general**
   rule with the affected cards listed, not as a one-card exception. That file is
   the only record of why the rendered text differs from the printed card, and
   also serves as a guide for implementing similar cards correctly.

Watch the `create_file` dup-first-line bug: after creating `.go` files, check
`line1 == line2` and drop the dup
(`for f in ...; do [ "$(sed -n 1p "$f")" = "$(sed -n 2p "$f")" ] && sed -i '' '1d' "$f"; done`).

## 3. Verify, then keep going

The targeted `go test` runs in step 2 are what verify each card as it lands.
`mage ci:check` is the **final** validation, run at the end to confirm everything in
your changes is working together — not after every card:

```sh
mage gen && mage ci:fix && mage ci:check  # gen = comments + rulebook; check prints ALL GREEN
mage tool:coverage                        # confirm the set's count moved
```

Run this once a mechanic and all the cards it unblocked have landed (and again
before you hand back), so `mage ci:check` validates the whole batch of your changes
rather than a single card. Aim to leave the tree with `mage ci:check` printing
`ALL GREEN` and the set's count higher than it started. Every mechanic you added should appear on the `/rulebook`
page, and nothing you retired should still be listed.

But a green gate is a checkpoint, not a finish line: the run's purpose is to keep
converting stubs into implemented cards. Do not stall chasing a green gate you did
not break. If `mage ci:check` fails only on **another agent's** in-progress change —
a file you never touched, a symbol you did not add, a check that was green before
your edits — record it briefly and move on to the next card rather than reverting
or "fixing" their work. Get _your_ changes passing; leave theirs alone. Report
progress as it lands rather than saving one summary for the end — but a report is
not a handoff: after reporting, immediately run `mage tool:nextCard` and begin the
next card. **Hand back (and only then call `task_complete`) solely when the stop
condition is met** or you have run out of implementable cards.

## 4. Running the backlog in parallel with subagents

A large backlog is faster to clear with several subagents working at once. The
one hazard is that they all edit the **same set package**, so a sibling's
half-written file transiently breaks everyone's `go build`/`go test`. Remove the
contention rather than serializing:

- **Triage first, on the orchestrator.** Before spawning anyone, read the whole
  backlog and classify every card into **easy** (composes from the facade),
  **medium** (one small new primitive), and **gated** (a genuinely new mechanic).
  Write the classification to the session triage memo. The subagents implement
  from this list; they do not each re-triage the set.
- **Know which regime you are in.** Early in a set the backlog is mostly easy, so
  parallel card agents win big. In a mature set the easy cards are gone and
  **every remaining card needs a primitive** — there parallel card agents have
  little to do until the engine grows, so the throughput lever is the engine, not
  the fan-out. Do not spin up empty-handed card agents against a mechanic-gated
  backlog; grow the engine first (next bullet).
- **Phase the engine ahead of the cards; do not overlap them.** The engine is a
  single-writer resource — only the mechanics agent touches `internal/engine` +
  `internal/card`, and a card agent cannot build a gated card until its primitive
  exists. So run the mechanics work as its **own phase**: land a primitive (or a
  cluster of them), reconcile centrally, and only **then** fan out card agents
  over the cards those primitives just unblocked. Overlapping a card-agent wave
  with an in-flight engine refactor is the main source of wasted time — the
  refactor breaks the shared build mid-wave and every card agent stalls on a
  failure that is not theirs. A **brand-new effect node is never added inside a
  subagent** — it is gated behind the human grill (_Building a mechanic_), so a
  mechanics phase that discovers it needs one surfaces it to the orchestrator to
  grill the human, rather than adding the node on its own.
- **Batch the primitives, not one-per-wave.** A mechanics phase that lands one
  primitive unblocks ~3 cards and pays the full reconcile cost for them. Group the
  triage's medium/gated cards by the primitive they share and have the mechanics
  agent build the whole **cluster** of related primitives in one phase (each with
  its engine test and rulebook term), so one reconcile releases a dozen cards
  instead of three. Order clusters by how many cards each unblocks.
- **Give each card subagent a disjoint set of cards** (hence disjoint files, since
  one card is one `<snake>.go` + `<snake>_test.go`). No two agents touch the same
  file. Hand them the card **names** and let them locate the stub; include the
  printed text or a one-line mechanic hint so they need no extra lookup.
- **Forbid `mage generateComments` / `mage gen` / `mage ci:fix` / `mage ci:check`
  inside subagents.** `generateComments` rewrites **every** card's doc comment
  repo-wide, and `ci:fix` reformats and autofixes every file — those are the real
  cross-agent races. Card agents write code + a targeted `go test` only; the
  **orchestrator runs `mage generateComments` once, centrally**, after a wave
  finishes, then the full `mage gen && mage ci:fix && mage ci:check`.
- **Quiesce before the central gate.** `mage ci:check` is the most contended command
  in the run — it catches every sibling's mid-edit as a failure. Run it only
  between waves, when your own agents have returned, and first glance at
  `git status --short internal/engine internal/web`: if a package you did not
  touch is mid-edit, wait for it to settle rather than diagnosing its red. When
  the gate fails only on files outside your wave's change set, record it and move
  on — do not chase a red gate you did not cause.
- **Tell every subagent that a transient build break from a file it did not write
  is expected** — wait a moment and retry the `go test`; never edit or revert
  another agent's file.
- **Make skips explicit.** A card agent that discovers its "easy" card actually
  needs a new primitive must **leave the stub, skip the card, and report it** with
  the missing primitive — never invent engine surface (that is the mechanics
  agent's job) and never force a broken implementation.
- **Waves, not one shot.** A wave is: a mechanics phase that lands a primitive
  cluster, then a fan-out of card agents over disjoint batches of the cards it
  unblocked (plus any still-easy cards). When they return, the orchestrator runs
  the central `generateComments` + `mage ci:fix && mage ci:check`, folds in the results, and plans
  the next wave's cluster. Keep the waves going until the stop condition is met.
- **When the mechanics agent adds new `Effect` nodes, the central `mage ci:check`
  will fail two ungated spots the subagents cannot see — fix them in the
  reconcile:** `internal/web` `TestIconTotality` fails until each new effect gets a
  glyph `case` in `internal/web/icon.go` `effectGlyphs` (the `iconFallbackAllowed`
  list is empty, so every effect needs a mapping — reuse the nearest sibling's
  asset), and `internal/cards` `TestOptionsAreInCanonicalOrder` fails if a card
  uses a `WithX` option no card used before until that option gets an
  `optionRank` entry in `cardorder_test.go`.

The stop condition, the "keep going" discipline, and the "leave another agent's
failures alone" rule from sections 1–3 all still hold — parallelism changes
_how many cards are in flight_, not _when the run ends_.

---
name: implement-cards
description: Work through a KeyForge set's unimplemented cards in this repo until the amount asked is done. Use when the user wants to implement a set's backlog ("implement more Call of the Archons cards", "keep going", "iterate until the set is done"), or triage what is easy.
---

The backlog is worked **one card at a time**, driven by `mage tool:nextCard`.
Ask the tool for the next unimplemented card, build it, drop its build tag, and
ask again. Each card leaves the engine slightly richer and the set measurably
more covered. The loop repeats until the stop condition is met.

The goal is **cards implemented**. Keep moving through the backlog: prefer
implementing the next card to polishing the last one. This is **one continuous
run through the entire stop condition** — never frame it to the user as a
"multi-session grind", a "first batch", or work that will be "continued later".
When the user says "do all the cards", they mean do all the cards now, in this
run, without handing back partway. Do not propose stopping, do not ask whether to
keep going, and do not offer a status report as a substitute for finishing. Tests
matter, and you should get your own work passing; but if a card is a reasonable
best-effort implementation and the only failures come from **another agent's**
in-progress changes (an unfamiliar file, a symbol you never touched, a gate that
was green before you started), note it and keep going rather than stalling on
someone else's edit.

**Do not call `task_complete`, and do not end your turn, until the stop condition
below is met.** A card that compiles, a green `mage check`, a reconnaissance
conclusion that "the rest all need new mechanics", and a written status summary
are ALL checkpoints, never the finish — reaching one means ask for the next card
and build it, not stop and report. "Every remaining card needs a new mechanic" is
the normal state of a backlog run, not a blocker: building the mechanic IS the
job, so keep building them one after another. The ONLY things that authorize
ending the run are: the stop condition is met, every remaining card is genuinely
blocked by something you cannot build (not merely "needs new work"), or the user
interrupts. After each card lands, the default and automatic next action is to run
`mage tool:nextCard` and start the next card — no pause, no check-in.

Read `internal/cards/AGENTS.md` (authoring + tests), `docs/card-wording-rules.md`
(rendered-text rules), and the root `AGENTS.md` (composability) before starting.

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

Then, once per run:

```sh
mage tool:coverage           # per-set covered/total — the number the run moves
SET=<slug> mage tool:nextCard # the next //go:build todo card: its file path,
                              # stats, printed text, and card.Provenance(...) call
```

Set slugs match the files in `internal/cards/provenance/` minus `.json`. `nextCard`
walks the set's missing cards in collector-number order and stops at the first one
still carrying a `//go:build todo` stub — so implementing a card (dropping its
build tag) is what advances the tool to the next one. With no `SET` it opens an
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
   composes from the facade (`internal/card/effects.go`, `target.go`,
   `options.go`): the effect nodes (`DealDamage`, `GainAember`, `Stun`, `Destroy`,
   `PutFromPlay`, `PurgeCreature`, `CaptureAember`, `Draw`, `GainChains`,
   `OnChooseCreature`, …), the targets (`card.Target.*` with chainable filters
   `.PowerAtMost()`, `.OfHouse()`, `.WithTrait()`, `.Damaged()`, `.Neighboring()`,
   `.UpTo()`, …), and the composites (`Sequence` of `Sentence`-wrapped effects,
   `ChooseOne`, `Conditional`, `Then`). Grep the effect files or a similar
   existing card to confirm a primitive's exact fields before using it. An easy
   card is built directly (_Implementing one card_ below).
2. **A gated card needs a mechanic that does not exist yet** — a new effect,
   target filter, count, selector, condition, or cross-turn hook. Build the
   mechanic (_Building a mechanic_ below), then implement the card on top of it.
3. **After a mechanic lands, cash it in.** Before returning to strict `nextCard`
   order, implement any other unimplemented card that the same mechanic now
   unblocks — that is what makes it a mechanic instead of a one-off. Use the
   triage memo (and a quick scan of the remaining stubs' printed text) to find
   them. Only when the mechanic is fully cashed in do you go back to `nextCard`.

**Verify with a targeted `go test`, not `mage check`, inside the loop.** A full
gate run is slow and its `ALL GREEN` is a false finish line that invites stopping,
so `mage check` is an **end-of-run** step, not a per-card one. Verify each card as
it lands by running `go test` for **exactly the card or mechanic you changed**:

```sh
go test ./internal/cards/sets/<slug>/ -run Test<Name>   # the card you just wrote
go test ./internal/engine/ -run Test<Mechanic>          # the mechanic it uses
```

Name the specific `-run` pattern — the card's `Test<Name>` and, for a gated card,
the engine `Test<Mechanic>` you added — so the run is a few seconds, not the whole
suite. Add `mage build` when a change spans packages. Only the two targeted tests
matter per card; save `mage check` for step 3.

### Building a mechanic

1. **Shape it for the whole cluster, not the first card.** Name and shape the
   mechanic so every card that wants it can use it. Prefer the cheapest engine
   surface: a new field or Strategy on an existing effect is cheaper than a new
   effect; a new effect is cheaper than a new `Resolver` capability, which is
   cheaper than new state.
2. **Add its engine test** in the matching `internal/engine/effect_*_test.go` as
   you go; `mage cover` gates `internal/engine` at 100%, and a card test does not
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
2. Seed it with a bare `// <Card Name>` comment above
   `var Name = card.New("Name", card.House.X, card.Type.Y, card.Rarity.Z,
   card.Provenance(card.<Set>, n), With*...)`. The card TYPE "action" is
   `card.Type.Tactic` (wording rule 19). Follow the one-field-per-line struct style
   in `internal/cards/AGENTS.md`. When an ability names the card's **own** house,
   write `card.House.Self` rather than repeating the house — but a card naming a
   _different_ house (Take That, Smarty Pants is about Logos creatures) spells that
   house out.
3. Write `<snake>_test.go` with the `ct.Play` harness — a `func Test<Name>` with
   `t.Run` subtests. A sole target auto-resolves; with 2+ candidates answer via
   `h.P1.ClickCard(handle)` / `h.P1.ClickOption(name)`. Set up a damaged creature
   with `handle.Damaged(n)`; read chains via `h.Game().State.Chains[0]`.
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
`mage check` is the **final** validation, run at the end to confirm everything in
your changes is working together — not after every card:

```sh
mage gen && mage check    # gen = comments + rulebook; check prints ALL GREEN
mage tool:coverage        # confirm the set's count moved
```

Run this once a mechanic and all the cards it unblocked have landed (and again
before you hand back), so `mage check` validates the whole batch of your changes
rather than a single card. Aim to leave the tree with `mage check` printing
`ALL GREEN` and the set's count higher than it started. Every mechanic you added should appear on the `/rulebook`
page, and nothing you retired should still be listed.

But a green gate is a checkpoint, not a finish line: the run's purpose is to keep
converting stubs into implemented cards. Do not stall chasing a green gate you did
not break. If `mage check` fails only on **another agent's** in-progress change —
a file you never touched, a symbol you did not add, a check that was green before
your edits — record it briefly and move on to the next card rather than reverting
or "fixing" their work. Get _your_ changes passing; leave theirs alone. Report
progress as it lands rather than saving one summary for the end — but a report is
not a handoff: after reporting, immediately run `mage tool:nextCard` and begin the
next card. **Hand back (and only then call `task_complete`) solely when the stop
condition is met** or you have run out of implementable cards.

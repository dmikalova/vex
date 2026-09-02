---
name: implement-cards
description: Work through a KeyForge set's unimplemented cards in this repo, one round at a time. Use when the user wants to implement or stub a set's backlog ("implement more Call of the Archons cards", "stub the rest", "keep going", "iterate until the set is done"), triage what is easy, or pick the next mechanic to build.
---

A set's backlog is worked in **rounds**. A round builds one missing mechanic and
then cashes in every card that mechanic unblocks, so each round leaves the engine
slightly richer and the set measurably more covered. Rounds repeat until the stop
condition is met.

Read `internal/cards/AGENTS.md` (authoring + tests), `docs/card-wording-rules.md`
(rendered-text rules), and the root `AGENTS.md` (composability) before starting.

## 1. Fix the stop condition, then set up

Restate the stop condition before touching code, because it decides when to hand
back: **the whole set**, **a count** ("ten more cards", "three rounds"), or a
**qualifier** ("all the Mars cards", "everything that isn't a lasting effect").
An unqualified "keep going" or "iterate" means the whole set.

Then, once per run:

```sh
mage coverage             # per-set covered/total — the number the run moves
mage missing <setSlug>    # each remaining card with stats, printed text, and a
                          # ready-made card.Provenance(...) call
mage stub <setSlug>       # scaffold a build-excluded stub for every one of them
```

Set slugs match the files in `internal/cards/provenance/` minus `.json`.

`mage stub` writes `internal/cards/sets/<slug>/<snake>.go` starting with
`//go:build todo`, so it is left out of the build, vet, test, lint, gencomments,
and the card registry — coverage stays honest. Each stub carries the card's
printed text, a `// TODO(stub)` marker, and a vanilla `card.New(...)` skeleton.
Existing files are never overwritten, so re-run it freely as rounds land.

## 2. Round one: harvest the free cards

The first round builds no mechanic. Walk the whole backlog and implement every
card that is already **easy** — its text composes entirely from primitives that
exist today. This clears the noise so later rounds choose between genuine
mechanics rather than tripping over cards that never needed one.

A card is easy when its whole text composes from the facade
(`internal/card/effects.go`, `target.go`, `options.go`): the effect nodes
(`DealDamage`, `GainAember`, `Stun`, `Destroy`, `PutFromPlay`, `PurgeCreature`,
`CaptureAember`, `Draw`, `GainChains`, `OnChooseCreature`, …), the targets
(`card.Target.*` with chainable filters `.PowerAtMost()`, `.OfHouse()`,
`.WithTrait()`, `.Damaged()`, `.Neighboring()`, `.UpTo()`, …), and the composites
(`Sequence` of `Sentence`-wrapped effects, `ChooseOne`, `Conditional`, `Then`).
Grep the effect files or a similar existing card to confirm a primitive's exact
fields before using it.

Everything else is **gated**: it needs a mechanic that does not exist yet — a new
effect, target filter, count, selector, condition, or a cross-turn hook. Leave a
gated card as its `//go:build todo` stub and, when the design is already clear,
add a one-line note after the TODO marker naming the mechanic it waits on.

## 3. Each later round: build the cheapest gate, then cash it in

1. **Group the gated cards by the mechanic they wait on.** Read the remaining
   stubs' printed text and cluster them: "needs a count of X", "needs a lasting
   reaction on event Y", "needs a target filter for Z".
2. **Pick the cheapest gate** — the smallest engine surface that unblocks the most
   cards. A new field or Strategy on an existing effect is cheaper than a new
   effect; a new effect is cheaper than a new `Resolver` capability, which is
   cheaper than new state. A gate that frees a single card is worth building only
   when nothing cheaper is left.
3. **Build the mechanic against the whole cluster, not the first card.** Name and
   shape it so every card in the group can use it — that is what makes it a
   mechanic instead of a card. Add its engine test in the matching
   `internal/engine/effect_*_test.go` as you go; `mage cover` gates
   `internal/engine` at 100%, and a card test does not count toward it.
4. **Teach the client to play the mechanic**, if it needs anything new. A mechanic
   that asks the player a question the browser client cannot ask is only half
   built. Read `internal/web/AGENTS.md` and check the new mechanic against it: a
   new `Chooser` prompt shape needs a case in `game_chooser.go` and a prompt in
   `view_controls.go`; a new zone or card state needs to be drawn in
   `view_board.go`/`view_card.go`; a new player action needs a keyboard route in
   `game_lifecycle.go` and a Tab stop in `game_nav.go`. Never reimplement the rule
   in the client — ask the engine, and add the reader to `internal/engine` if it
   does not exist.
5. **Implement every card the gate unblocks**, then close the round (step 4).

State at the top of each round which gate you picked and which cards it frees, so
the round has a visible bound.

### Implementing one card

1. Delete the stub (`command rm <snake>.go` — `rm -f` is blocked by an alias) and
   write the real `internal/cards/sets/<slug>/<snake>.go`.
2. Seed it with a bare `// <Card Name>` comment above
   `var Name = card.New("Name", card.House.X, card.Type.Y, card.Rarity.Z,
   card.Provenance(card.<Set>, n), With*...)`. The card TYPE "action" is
   `card.Type.Tactic` (wording rule 19). Follow the one-field-per-line struct style
   in `internal/cards/AGENTS.md`.
3. Write `<snake>_test.go` with the `ct.Play` harness — a `func Test<Name>` with
   `t.Run` subtests. A sole target auto-resolves; with 2+ candidates answer via
   `h.P1.ClickCard(handle)` / `h.P1.ClickOption(name)`. Set up a damaged creature
   with `handle.Damaged(n)`; read chains via `h.Game().State.Chains[0]`.
4. Run `mage generateComments` — it rewrites the card and test doc comments from
   the definition. Read the generated text against `docs/card-wording-rules.md`.
   A wording fix means changing the effect's `Text()` in
   `internal/engine/effect_*.go`, never hand-editing the comment.

Watch the `create_file` dup-first-line bug: after creating `.go` files, check
`line1 == line2` and drop the dup
(`for f in ...; do [ "$(sed -n 1p "$f")" = "$(sed -n 2p "$f")" ] && sed -i '' '1d' "$f"; done`).

## 4. Close every round green

```sh
mage gen && mage check    # gen = comments + rulebook; check must print ALL GREEN
mage coverage             # confirm the set's count moved
```

A round is finished when `mage check` prints `ALL GREEN` and the set's count has
gone up. Then start the next round at step 3, or hand back if the stop condition
is met. Report each round's gate, cards, and new count as it lands rather than
saving one summary for the end.

## 5. Consolidate every third round

Mechanics added one round at a time drift toward a pile of near-duplicates. Every
third round, before starting the next gate, sweep the mechanics the run has added
against the "Composition and design" section of `docs/style-guide.md` and
`internal/engine/AGENTS.md`, and act on what you find:

- **Consolidate** near-duplicates. Two effects differing by a constant or a side
  are one effect parameterized over an enum; two conditions asking the same
  question of different subjects are one condition with a `Player`.
- **Decompose** anything fused. "A, then B if C" belongs in the card as
  `Sequence{A, Conditional{C, B}}`, not in one bespoke node; thread values through
  `EffectContext` (`ctx.It`, `ctx.ChosenHouse`, `ctx.Produced.*`).
- **Simplify** down the ladder: a field or Strategy on an existing effect (a
  `Count`, `Selector`, `Condition`, `Chooser`) beats a new node, and a new node
  beats new state.
- **Re-express** the cards already using the old shape, and delete what the
  consolidation retired. The card and engine tests pin both behavior and rendered
  text, so refactor freely and let the suite catch regressions.

Close the sweep the same way a round closes: `mage gen && mage check`.

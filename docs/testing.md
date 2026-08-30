# Writing tests

This guide explains the testing options in Vactrol, when to reach for each, and
what to actually test at each layer. It is the human-facing companion to the
per-directory rules: [`internal/cards/AGENTS.md`](../internal/cards/AGENTS.md) has
the full card-harness API reference, and
[`internal/engine/AGENTS.md`](../internal/engine/AGENTS.md) explains the engine
design the tests pin down. See also [architecture.md](architecture.md) for how the
pieces fit together.

## The short version

- **Engine rule or effect?** Write an **engine test** (`package engine`) against
  small blueprints. The engine is held at **100% coverage**, so every branch needs
  one.
- **A specific card's behavior?** Write a **card test** with the `ct` harness
  (`ct.Play`), one per card, reading like a game.
- **Shared setup / cross-cutting logic (e.g. `match`)?** Write an ordinary
  package test.
- **Frontend glue (`tui`, `web`)?** Largely untested by design (TTY/DOM-bound);
  push logic worth testing down into the engine or `match`.

## Coverage philosophy: what is gated and why

`mage cover` measures **only `internal/engine`, and requires 100%**. That is
deliberate:

- The engine is where the value and the risk concentrate — the rules — and it has
  no UI or I/O to dilute the measurement, so 100% is both meaningful and
  achievable.
- A new engine code path (a new effect branch, a new rule edge case) **must** come
  with an engine test, or the gate fails. This is a feature: it forces you to
  cover the branch where it lives rather than hoping a card test happens to reach
  it.
- Card set packages are essentially data (`var X = card.New(...)`), so their own
  statement coverage is not the useful signal; their tests exist to pin
  _behavior_ and _rendered text_, not to hit lines.

Consequence to internalize: **if you delete or change an engine test, re-check
coverage** — a card test exercising the same path does not count toward the engine
gate (`mage cover` doesn't run card tests).

## Option 1 — Engine tests (`package engine`)

Internal tests in `package engine` (not `engine_test`) so they can use unexported
internals directly. They exercise the rules against a few **blueprints** defined
in the test files, not cards from the database — `engine` must not import the card
packages that import it.

Two shapes, depending on what you're testing:

**A single effect node** — construct it, assert its `Text()`, then `Resolve` it
against a hand-built context and assert the state change. This is the canonical
way to cover an `Effect`/`Count`/`Condition`/`Target`:

```go
func TestHealEffect(t *testing.T) {
    g := NewGame("A", "B", 1)
    src := g.AddToBattleline(testCreature("src", 5), 0)
    g.State.Cards[src].Damage = 3
    ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

    h := Heal{Amount: 2, Target: Target{Kind: TargetThisCreature}}
    if h.Text() != "heal 2 damage from "+SelfName {   // text and…
        t.Errorf("text = %q", h.Text())
    }
    h.Resolve(ctx)                                     // …behavior, from one node
    if g.State.Cards[src].Damage != 1 {
        t.Errorf("damage = %d, want 1", g.State.Cards[src].Damage)
    }
}
```

**A game rule** — build a small board with `started(t)` (a game with player 0
active and a house chosen) or `NewGame`, drive it through the public `Game`
methods (`Play*`, `Reap`, `Fight`, `EndTurn`, …), and assert state and returned
errors.

Helpers available (in `helpers_test.go`):

- `started(t)` — a ready game; `testCreature(name, power, opts...)` — a quick
  blueprint; `exGiant()/exBruteStrength()/exBattleFury()/exAutocannon()` — richer
  example cards.
- `handIdx`/`handIdxByID` — locate a card in hand.
- Custom **choosers** (`orderLastChooser`, `orderRejectChooser`, `orderAllChooser`)
  to drive or reject choices deterministically; the default `FirstChooser` picks
  the first candidate.

**What to test here:**

- Every effect's `Text()` _and_ `Resolve()` — they're one node, so cover both.
- Every rule branch and edge case: clamping (Æmber/damage never below zero),
  simultaneous combat and destruction ordering, empty-deck reshuffle, illegal
  actions returning the right sentinel error, `validate()` rejecting bad field
  combos, and the odd-but-legal boundaries (`maxCards`, over-heal flooring).
- Rendering rules that fold or reword text (e.g. Fight/Reap merging).

Because the gate is 100%, the practical rule is: **the test that makes a new
branch reachable lives next to that branch.**

## Option 2 — Card tests (`ct` harness)

Every card has its own `snake_case_test.go` in its set package, built on the
`cardtest` harness (imported as `ct`). A card test declares the whole scenario up
front and then plays it:

```go
func TestAmmoniaClouds(t *testing.T) {
    t.Run("deals 3 damage to each creature", func(t *testing.T) {
        var toughFoe ct.Card
        h := ct.Play(t, ct.Setup{
            P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(AmmoniaClouds)},
            P2: ct.Side{InPlay: ct.Cards(
                ct.Bind(&toughFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
                ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2)),
            )},
        })

        h.P1.Play(AmmoniaClouds)

        h.Expect(toughFoe).At(ct.PlayArea).Damage(3)
    })
}
```

The pieces (full reference in [`internal/cards/AGENTS.md`](../internal/cards/AGENTS.md)):

- `ct.Play(t, ct.Setup{P1, P2, Seed})` builds the match with P1's house chosen.
  Each `ct.Side` sets `House`, the zones via `ct.Cards(...)`, and `Amber`/`Keys`.
- `ct.Creature/Artifact/Tactic/Upgrade(opts...)` build **vanilla** cards — a body
  with no baggage — for isolating the card under test; options `ct.OfHouse`,
  `ct.Power`, `ct.Keywords`, ….
- `ct.Bind(&handle, def)` names a placed card to reference later;
  `ct.Upgraded(host, ups...)` attaches upgrades at setup.
- Players act by def **or** handle: `h.P1.Play/Reap/Fight/UseAction/EndTurn/
ChooseHouse`. A choice among several candidates pauses — answer with
  `h.P1.ClickCard(x)` / `h.P1.ClickOption("...")` and assert it with
  `h.P1.ExpectPrompt("...").Source("Card")`. A sole candidate auto-resolves.
- Assert with `h.Expect(defOrHandle).Damage/Power/Armor/AmberOn/Exhausted/Ready/
Stunned/At(zone)` and `h.P1.ExpectAmber/ExpectKeys`. Drop to `h.Game()` for
  anything the fluent API doesn't cover.

**What to test here:**

- The card's _observable behavior_ end-to-end: play it (or reap/fight/use it) and
  assert what changed on the board and in the pools.
- The interesting branches of _that card_ — its choice being declined, its
  optional repeat stopping, a trigger firing or not firing.
- Vanilla creatures/artifacts still get a one-line test (play + stat check): it
  documents the card and guards its stats, and the card `var` loads at package
  init so untested-but-registered cards still count as covered.

**What NOT to test here:**

- Don't re-prove engine mechanics a card merely _uses_ — combat math, clamping,
  ordering — that's the engine's job and its 100% gate already covers it. Test
  what _this card_ adds.
- Don't reach into engine internals or assert redundant zone counts that `At(zone)`
  already implies. Card tests use the public API and exported card `var`s only.
- Don't hand-assert the full rendered card text in a card test; the card-text doc
  comment is generated and the `Text()` assertions live in the engine tests.

## Option 3 — Package / integration tests

Ordinary tests for logic that isn't the rules engine or a single card — e.g.
`internal/match` (deck construction) or the `cards` aggregator (database validity:
unique names, real houses/types, deterministic ordering, every creature/artifact
has a trait). Use internal (`package match`) tests when you want to cover
unexported helpers, and an external (`package <pkg>_test`) test when you only need
the public surface. Seed any RNG so results are reproducible (`match.New(..., seed)`).

## Conventions that apply everywhere

- **describe / it via `t.Run`.** Model behavior as subtests: `func TestFoo` is the
  "describe", each `t.Run("does X", ...)` is an "it", and a `setup := func(t){…}`
  closure is the "beforeEach". This frees the `func TestFoo` doc comment for the
  generated card-text block (`mage generateComments` splices it above card tests).
- **Determinism.** Games are seeded (`NewGame(_, _, seed)`, `ct.Setup{Seed}`,
  `match.New(_, _, seed)`); the default `FirstChooser` and the harness's
  auto-resolve keep choices predictable. Never rely on wall-clock or map-iteration
  order.
- **Test at the boundary.** Prefer the public API and observable state over poking
  internals — except in `package engine` tests, whose whole point is to reach the
  unexported rule you're covering.

## Running tests

```sh
mage test         # go test ./...  (whole suite)
mage cover        # engine coverage, must print 100%
mage check        # the full gate: fmt, build, vet, lint, test, cover

# focused runs while iterating:
go test ./internal/engine/ -run TestHeal
go test ./internal/cards/sets/callofthearchons/ -run TestAmmoniaClouds
```

# Card authoring guide

This directory holds the card database. Package `cards` (`cards.go`) is the
aggregator; each set lives in its own subpackage under `sets/`
(e.g. `sets/callofthearchons/`) and every card is one self-registering file.
Shared test helpers live in the sibling `cardtest/` package — it is a test
helper (stdlib `httptest`/`iotest` convention), not a set, so it sits outside
`sets/`.

## Rules reference

When a rules question isn't settled by the vactrol implementation itself or by
[docs/rulebook.md](../../docs/rulebook.md), the converted KeyForge Master Rulebook
at [docs/keyforge-master-rulebook.md](../../docs/keyforge-master-rulebook.md) is
available for reference. Take it with a grain of salt: it is a converted PDF, and
while much of it matches how vactrol behaves, some nuances have been intentionally
simplified — so the implementation and this repo's own docs win where they
disagree.

## File layout

- One card per file, named `snake_case.go` after the card (e.g. `dust_imp.go`,
  `ammonia_clouds.go`). The matching test is `snake_case_test.go`.
- Author cards through the `card` facade only
  (`github.com/dmikalova/vactrol/internal/card`). Use the grouped
  namespaces — `card.House.X`, `card.Type.X`, `card.Rarity.X`,
  `card.Keyword.X`, `card.Trigger.X`, `card.Target.X` (and the relative-player
  values `card.Controller` / `card.Opponent`) —
  never the raw engine package.
- `card.New(...)` self-registers the card (via `init`); there is no central
  list to update. Adding a card is just adding a file.
- Tag each card's origin with `card.Provenance(card.<Set>, <number>)` (e.g.
  `card.Provenance(card.CotA, 1)`) as the first option — this links it to the
  original KeyForge card it derives from and drives the TUI Provenance coverage
  view. Source-set aliases (`card.CotA`, `card.DT`, …) and the catalogs live in
  `internal/cards/provenance`. Optional and repeatable; omit for wholly original
  cards.

## Style

The doc comment above each card is **generated, not hand-written** — run
`mage generateComments` (it rewrites every card's comment from its definition,
the same details the TUI card box shows via `engine.RenderCardText`). Write the
`card.New(...)` call and let the generator fill in the comment:

```go
// Ammonia Clouds
//
//	House:  Mars
//	Type:   Action
//	Rarity: Common
//
//	Play: Deal 3 damage to each creature.
var AmmoniaClouds = card.New(
	"Ammonia Clouds",
	card.House.Mars,
	card.Type.Action,
	card.Rarity.Common,
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Target: card.Target.EachCreature,
		}),
)
```

The generated comment (tab-indented so godoc renders the block preformatted;
`gofmt` requires the blank `//` line after the title) is:

1. `// <Card Name>` — the card's display name.
2. A labeled, colon-aligned block: `House`, `Type`, `Rarity`, then `Power` and
   `Armor` for creatures, `Æmber` when the card has an Æmber bonus, and `Traits`
   when it has traits.
3. After a blank `//`, the printed rules text (keywords, upgrade modifier, and
   ability lines) — omitted entirely for vanilla cards with no rules text.

So authoring a card is just adding the `card.New(...)` file and running
`mage generateComments`; a placeholder `// <Name>` line above the var is enough
to seed it. `card.New(...)`:

- The four positional arguments — name, `card.House.X`, `card.Type.X`,
  `card.Rarity.X` — each on their own line.
- Each `card.With*` option on its own line, in this order: `WithPower`,
  `WithArmor`, `WithAemberBonus`, `WithTraits`, `WithKeywords`, `WithStatic`,
  `WithAbility`.
- `card.WithAbility(` breaks onto the next line; the trigger and effect share a
  line (`card.Trigger.Play, card.DealDamage{`).
- **Every struct literal with two or more fields is written one field per line**,
  for readability and consistent diffs. This covers effects
  (`card.DealDamage{Amount: …, Target: …}`), the modifier/value structs passed to
  the `With*` options (`WithStatic(card.StaticModifier{PowerBonus: …, HazardousBonus: …})`,
  `WithAttackDamage(card.AttackDamage{Amount: …, FlankOnly: …})`,
  `WithConstantAbility(card.ConstantAbility{…})`), and nested effects the same
  way. A struct with a single field stays inline (e.g. `card.GainAember{Amount: 1}`,
  `card.Stun{Target: card.Target.This}`).
- A granted / `WithAbilities` entry keeps its trigger and effect on one line
  (`{Trigger: card.Trigger.Reap, Effect: card.DealDamage{`), mirroring the
  `WithAbility(trigger, effect)` form; if that effect has two or more fields,
  break its fields onto their own lines as above.
- Slice elements that are themselves single-field or empty structs stay inline
  within the slice (e.g. `Verbs: []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}}`).

Run `gofmt -w .` after editing (it aligns the fields; it does not add the line
breaks, so the one-field-per-line layout above is the author's responsibility).

## Wording rules

The generated comment is the card's printed text, produced by the effect AST's
`Text()` methods. That text must obey [docs/card-wording-rules.md](../../docs/card-wording-rules.md)
— the curated wording conventions (front-load `for each`; `gains` not `gets`;
`Put` not `Return`; result gates `A -> B`; `Aember` not `Æmber` in curated
source; capital `Damage` only where dealt; self-reference by name; etc.).

When adding or reviewing a card, read its generated comment against those rules.
**Call out any line that violates a rule**, and **auto-apply the fix when it is
obvious** — because the text is generated from the AST, a wording fix means
changing the effect rendering, not hand-editing the comment (see
[../engine/AGENTS.md](../engine/AGENTS.md) for how effects render their own text,
and prefer reshaping a `Target`/`Count`/`Selector` over adding a new effect):

- If a whole family of cards renders wrong (e.g. a `for each` clause appearing at
  the end instead of the front), fix it once in the effect's `Text()` in
  `internal/engine/effect_*.go`, then re-run `mage generateComments`.
- If only one card reads wrong, it usually means the wrong effect/target was
  chosen — fix the `card.New(...)` definition.
- If a rule is genuinely ambiguous or the "fix" would change game behavior, don't
  silently apply it — flag it for the user instead.

After any wording change, run `mage generateComments` (regenerates every card's
comment) and the engine tests (the `Text()` assertions live in
`internal/engine/effect_*_test.go`).

## Tests

- Every card has its own `snake_case_test.go` in the same set package, built on
  the `cardtest` harness (`internal/cards/cardtest`, imported as `ct`). A test
  reads like the game: declare the board with `ct.Play`, act through the players,
  and assert with `h.Expect`.
- The card-text doc comment above `func Test<Card>` is generated by
  `mage generateComments` (it mirrors the card box, the same block the card's own
  source file carries) — don't hand-write it. Put the behavior description in a
  `t.Run("...")` subtest name instead, in the Keyteki-style describe/it shape
  (the `func Test<Card>` is the "describe", each `t.Run` is an "it").
- Sketch of the shape:

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

- The pieces (`h` is the `*ct.Harness` returned by `ct.Play`):
  - `ct.Play(t, ct.Setup{P1, P2})` builds the game with player 1's house chosen.
    Each `ct.Side` sets `House`, the zones `InPlay`/`Hand`/`Deck`/`Discard`/
    `Archives` (via `ct.Cards(...)`), and `Amber`/`Keys`.
  - `ct.Creature/Artifact/Tactic/Upgrade(...)` build vanilla cards — a "body with
    no baggage" for isolating a mechanic — with options `ct.OfHouse`, `ct.Power`,
    `ct.Armor`, `ct.Keywords`, `ct.PowerBonus`, … (default house Brobnar).
  - `ct.Bind(&handle, def)` names a placed card so you can reference it later;
    `ct.Upgraded(host, upgrades...)` attaches upgrades at setup.
  - Players act by card definition **or** handle: `h.P1.Play/Reap/Fight/UseAction`,
    `h.P1.EndTurn/ChooseHouse`. A choice among several candidates pauses; answer it
    with `h.P1.ClickCard(x)` or `h.P1.ClickOption("...")`, and assert the prompt
    with `h.P1.ExpectPrompt("...").Source("Card")`. A sole candidate auto-resolves.
  - Assert with `h.Expect(defOrHandle).Damage/Power/Armor/AmberOn/Exhausted/
Ready/Stunned/At(zone)` (chainable) and `h.P1.ExpectAmber/ExpectKeys`. Reach
    the raw engine via `h.Game()` for anything the fluent API doesn't cover.
- Set test packages import `card` (for `card.House.X`) and `ct`
  (`internal/cards/cardtest`); they use the public engine API and exported card
  `var`s (no reaching into engine internals).

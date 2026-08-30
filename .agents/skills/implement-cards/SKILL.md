---
name: implement-cards
description: Implement or stub the remaining unimplemented cards of a KeyForge set in this repo. Use when the user wants to work through a set's backlog (e.g. "implement more Call of the Archons cards", "stub the rest of the cards", "what's left to implement"), triage which cards are easy, or scaffold stub files for the hard ones.
---

Work through a set's unimplemented cards: stub every one, implement the easy ones
for real, and leave the hard ones as build-excluded stubs marked for later.

Read `internal/cards/AGENTS.md` (authoring + tests), `docs/card-wording-rules.md`
(rendered-text rules), and the root `AGENTS.md` (composability) before starting.

## 1. See what's left

```sh
mage missing <setSlug>    # e.g. callofthearchons — lists each card with stats,
                          # printed text, and a ready-made card.Provenance(...) call
mage coverage             # per-set covered/total
mage lookup "<name>"      # one card's full details
```

Set slugs match the files in `internal/cards/provenance/` minus `.json`.

## 2. Scaffold stubs for the whole backlog

```sh
mage stub <setSlug>
```

This writes a build-excluded stub `internal/cards/sets/<slug>/<snake>.go` for
every unimplemented card. Each starts with `//go:build todo`, so it is left out of
the build, vet, test, lint, gencomments, and the card registry — coverage stays
honest. The stub carries the card's printed text and a `// TODO(stub)` marker plus
a vanilla `card.New(...)` skeleton (house/type/rarity/power/traits already filled).
Existing files are never overwritten, so it is safe to re-run as you make progress.

## 3. Triage: easy vs. hard

A card is **easy** when its whole text composes from primitives that already exist
in the facade (`internal/card/effects.go`, `target.go`, `options.go`): the effect
nodes (`DealDamage`, `GainAember`, `Stun`, `Destroy`, `PutFromPlay`, `PutUpTo`,
`PurgeCreature`, `CaptureAember`, `Draw`, `GainChains`, `OnChooseCreature`, …), the
targets (`card.Target.*` with chainable filters `.PowerAtMost()`, `.OfHouse()`,
`.WithTrait()`, `.Damaged()`, `.Neighboring()`, …), and the composites (`Sequence`
of `Sentence`-wrapped effects, `ChooseOne`, `Conditional`). Grep the effect files
or an existing similar card to confirm a primitive's exact fields before using it.

A card is **hard / unsure** when it needs a mechanic that does not exist yet (a new
effect, target filter, count, selector, or a cross-turn/lasting hook). **Leave it
as its `//go:build todo` stub** — do not force a wrong or non-compiling
implementation. If the design is clear, add a one-line note after the TODO marker.

## 4. Implement an easy card

1. Delete the stub file (`command rm <snake>.go` — note `rm -f` is blocked by an
   alias) and create a real `internal/cards/sets/<slug>/<snake>.go`.
2. Seed it with just a `// <Card Name>` comment above
   `var Name = card.New("Name", card.House.X, card.Type.Y, card.Rarity.Z,
card.Provenance(card.<Set>, n), With*...)`. The card TYPE "action" is
   `card.Type.Tactic` (wording rule 19). Follow the one-field-per-line struct
   style in `internal/cards/AGENTS.md`.
1. Write `<snake>_test.go` with the `ct.Play` harness (a `func Test<Name>` with
   `t.Run` subtests). A sole target auto-resolves; when 2+ candidates exist (e.g. a
   self-targeting `Stun`, or "an artifact" that includes the source), answer with
   `h.P1.ClickCard(handle)` / `h.P1.ClickOption(name)`. Set up a damaged creature
   with `handle.Damaged(n)`; read chains via `h.Game().State.Chains[0]`.
2. Run `mage generateComments` — it rewrites the card and test doc comments from
   the definition. Read the generated rules text against `docs/card-wording-rules.md`.
   A wording fix means changing the effect's `Text()` in
   `internal/engine/effect_*.go`, never hand-editing the comment.

Watch the `create_file` dup-first-line bug: after creating `.go` files, check
`line1 == line2` and drop the dup
(`for f in ...; do [ "$(sed -n 1p "$f")" = "$(sed -n 2p "$f")" ] && sed -i '' '1d' "$f"; done`).

## 5. Verify

```sh
mage gen && mage check    # gen = comments + rulebook; check must print ALL GREEN
mage coverage             # confirm the set's count went up
```

`mage cover` gates only `internal/engine` at 100% — a new engine primitive needs an
engine test, not just a card test. If you added or reshaped an effect to fit a card,
add its `Text()`/`Resolve()` coverage in `internal/engine/effect_*_test.go`.

## When a card almost fits an existing primitive

Prefer reshaping the primitive (a new `Target` filter, a `Count`, a `Selector`, a
`Duration`) over a bespoke one-off effect, and re-express the cards already using
it — the tests pin each card's behavior and rendered text, so refactor freely and
let the suite catch regressions. See the "Composition and design" section of
`docs/style-guide.md` and `internal/engine/AGENTS.md`.

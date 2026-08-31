# Card templates materialize to concrete defs at generation time

## Context

Some cards are not fully determined until a deck is generated. MM mutants compose
two houses into one card (7 home × 6 other = 42 outcomes) chosen at generation;
Master of 1/2/3 is one family with an enumerable parameter; a maverick that names
its own house (Pitlord's "choose Dis as your active house") must rebind that house
to the slot it lands in, or it reads nonsensically. Enumerating every outcome as a
concrete card is impossible for the combinatorial cases and wasteful for the rest.
But the engine's `CardDefinition` is immutable, flat, pointerless, comparable data
held in a read-only catalog, and `GameState` holds no closures — so a "card that is
still a function of parameters" cannot exist inside the engine.

## Decision

Give a card two representations. A **concrete** `engine.CardDefinition` is what the
engine consumes, unchanged. A **template** is a parameterized blueprint that lives
in the card facade layer, occupies a single pool entry, and carries a build closure
— but implements a port that `deckgen` owns:

```go
type Materializer interface {
    Materialize(ctx SlotContext, r *rand.Rand) engine.CardDefinition
}
```

Generation runs a per-slot **Materialize** stage (before validation) that compiles
each slot to a concrete def: binding template parameters, rehousing a maverick
(`def.House = slotHouse`), and binding home-house self-references to the slot's
house. A template exposes a fixed _face_ — `{House, Type, Rarity,
GenerationProfile}` — so pool bucketing and rarity rolls treat it exactly like a
concrete card; only its effects, name, and parameters defer to `Materialize`.
Concrete cards get a trivial identity materializer. Enhancements and distortions
are a **separate** deck-wide finishing pass, not part of per-slot materialization,
because they redistribute across cards.

## Consequences

- The build closure lives in the card layer and runs at generation time; it never
  enters `GameState`, and its output is a flat pointerless `CardDefinition`. The
  engine's copy-is-a-value-copy invariant (ADR 0001) is untouched, and there is
  zero per-variant engine code.
- 42 mutants become one base template per home house plus a shared generator for
  the second half; families collapse to one pool entry each; a self-house maverick
  reads correctly with no bespoke handling.
- `Materialize` is the single seam for "produce the concrete card that fills this
  slot," so maverick rehousing, home-house binding, template families, and
  connection expansion all route through one mechanism rather than special cases in
  the fill path.
- The engine still needs new bonus/distortion representation before the finishing
  pass can land non-Æmber enhancements; that is deferred, and the seam does not
  depend on it.

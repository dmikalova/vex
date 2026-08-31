# Deck generation is a pure function below the card facade

## Context

Procedural deck generation (choose 3 houses, fill 36 slots by rarity, resolve
mavericks/legacies/connections, land enhancements) is a substantial subsystem that
today does not exist — `internal/match` just repeats a house pool and shuffles. It
needs to live somewhere, and the natural first instinct is to put it beside or
inside the card facade (`internal/card`), since that is where the card registry and
the provenance/sidecar metadata already live. But generation only needs a _built
pool of cards plus tunables_, not the registry, not provenance (which is
author-only coverage metadata), and not the engine runtime. It also must stay
trivially testable with synthetic pools, and eventually be shared with the
self-play bot's MCTS scorer (a procedural, heuristics-based player — not a
generative model).

## Decision

Put generation in a new package `internal/deckgen` that depends **only on
`engine`'s data types** (`House`, `CardType`, `Rarity`, `CardDefinition`) — never
the `Game` runtime, never the `card` facade, never `provenance`. Invert the
dependency so `deckgen` defines the pure types (`Set`, `Deck`, `Slot`, the
`Materializer` port, `GenerationProfile`) and the `card`/`cards` layer depends on
`deckgen` to assemble a `Set` and attach profiles. The public surface is a pure
function:

```go
func Generate(set Set, seed int64) Deck
```

Final layering: `cards → card → deckgen → engine`, with `bot/scoring → deckgen`
later. `match` remains the only code that pushes a generated `Deck` into a live
`*engine.Game`.

## Consequences

- `Generate` is a pure, deterministic function of `(Set, seed)`, unit-testable
  with a hand-built synthetic `Set` and no engine game, no registry, no I/O.
- Depending on `engine`'s _data_ types (which are flat, comparable, runtime-free)
  costs nothing and avoids the alternative — a `deckgen`-local `Card` interface or
  duplicated enums — which would add boxing and an interface with a single
  implementer, un-idiomatic in Go. "Orthogonal to the engine" means orthogonal to
  the _runtime_, not to its value types.
- Provenance never enters generation; set membership is derived from the set
  packages, not from provenance refs.
- The inversion is the reason templates and scoring can define their types in
  `deckgen`/`scoring` and have the card facade adapt to them, rather than the other
  way around — see ADR 0004.

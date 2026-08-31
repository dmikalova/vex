# Effects render their own text and carry themselves out

## Context

In the usual card-game implementation a card's rules text is one artifact (data,
or a hand-written string) and its behavior is another (code). The two are free to
drift: a wording tweak that isn't mirrored in the logic, or a logic fix that isn't
mirrored in the text, and the card now lies about what it does. For a game whose
whole surface is hundreds of small rules, that drift is a constant, silent bug
source.

## Decision

Every card ability is a tree of `Effect` nodes — an **Interpreter AST**. Each node
implements **both** `Text()` (renders its own English) and `Resolve(ctx)` (carries
itself out). One value drives both, so printed text can never desync from
behavior; card doc comments are generated from the tree (`mage generateComments`).
A new mechanic is almost always a **new node** in `effect_<mechanic>.go`, not a new
branch in the `Game` runtime. When behavior varies along an axis, the axis is a
small **Strategy that also renders its own text** (`Chooser`, `Selector`, `Count`,
`Condition`) rather than a new field or `bool`.

## Consequences

- Rules text and behavior are the same value read two ways — they cannot disagree.
- The cost is ~40 node types and no compiler-checked exhaustiveness for
  cross-cutting operations. `Text()`/`Resolve()` are the only methods intrinsic to
  a card's identity; any **extrinsic** whole-tree operation (bot value, static
  analysis, serialization) is a standalone type-switch function in one file, not a
  third method on every node.
- Effects reach the game only through the `Resolver` port (ADR 0008), never
  `*Game` or `GameState` directly.

See `internal/engine/AGENTS.md` for how to add a node and where each pattern lives.

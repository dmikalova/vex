# Engine design guide

`internal/engine` is the pure game engine. It imports nothing upward (see the
dependency layering in the root `AGENTS.md`) and is the one package held to 100%
test coverage (`mage cover`). This document is the **design ideal** for the
engine: the patterns it is built from, where each belongs, and the deliberate
tradeoffs — including the ones that look like smells but are forced by a
constraint, and how to handle them so the next mechanic still lands cleanly.

Read this alongside [docs/style-guide.md](../../docs/style-guide.md) ("Composition
and design") and the file-organization rules in the root `AGENTS.md`. This guide is
the _why_; those are the _what goes where_.

## Two hard constraints shape everything

Almost every non-obvious choice in the engine traces back to one of these. When
something looks un-idiomatic, check it against these before "fixing" it. Both are
recorded as ADRs — read them for the full rationale and the rejected alternatives.

1. **Flat, pointerless, comparable state** (ADR 0005). `GameState` is a value
   struct of fixed arrays (`[maxCards]CardCore`, `[2]Zone`) — no pointers, slices,
   or maps — so `GameState.FastCopy()` is a plain value copy for MCTS cloning. What
   that forces in day-to-day work:
   - State cannot hold a closure or `Effect`; "do X later" is flat, enum-tagged
     data (the lasting registry, ADR 0007), not a stored `func`.
   - Values compared against state (notably `Target`, `== Target{}` in
     `ConstantAbility`) stay comparable: a wide flag struct, not a `[]filter`;
     paired `x int` + `hasX bool`, not `*int`.
   - Cards are `LocalID uint8` indices into a read-only `catalog` shared across
     clones.

2. **Effects touch the game only through the `Resolver` port** (ADR 0008). An
   `Effect` holds a `Resolver` via `EffectContext`, never `*Game` or `GameState`;
   the port is the complete, auditable catalogue of card-facing capability (Ports &
   Adapters: the effect AST is the domain, `Resolver` the port, `*Game` the
   adapter).

## Patterns in use (and where each belongs)

- **Interpreter — the effect AST (`Effect`)** (ADR 0006). Every node has `Text()`
  (renders English) and `Resolve(ctx)` (carries it out); one value drives both, so
  printed card text can never desync from behavior. A new mechanic is almost always
  a new `Effect` node in `effect_<mechanic>.go`, not a new branch in the `Game`
  runtime.
- **Composite — `Sequence`, `Conditional`, `ChooseHouseThen`, `MayRepeat`, …**
  compose child `Effect`s and recurse `validateEffect` into them. Prefer composing
  small nodes over one fused node (root `AGENTS.md`: "decompose fused effects").
- **Strategy — the `Chooser` family, and the `Selector` / `Count` / `Condition`
  trio.** See the next section; this is used heavily and should keep being the
  first tool reached for when behavior varies along an axis.
- **Ports & Adapters — `Resolver`** (ADR 0008). See constraint 2 and the
  segregation section.
- **Facade — the `internal/card` package** wraps the engine for authoring so card
  files never import `engine`.
- **Null Object / invalid-zero** (ADR 0010). `FirstChooser` is the default
  strategy; `playerUnset` / `targetUnset` / `durationUnset` / `eventUnset`
  sentinels make an omitted required field an error caught at card-init
  `validate()` rather than a silent default.

## Strategy: use it, and keep the text with the behavior

Strategy is the engine's second backbone after Interpreter. The mature form here
is a strategy that carries **both** its behavior and its printed-text fragment, so
it plugs into the AST without desync:

- **`Chooser` (`game.go`) is the decision strategy**, swapped per frontend:
  `FirstChooser` (bot/tests, deterministic), `teaChooser` (TUI), `webChooser`
  (web), `bridgeChooser` (test harness). The engine never knows _how_ a decision
  is made. Extend a chooser's capability with an **optional capability
  interface** discovered by type assertion — `OptionChooser`, `Orderer` — with a
  graceful fallback when unimplemented. That is the idiomatic-Go form of Strategy
  (cf. `io.WriterTo`, `http.Flusher`); prefer it over widening the base `Chooser`.
- **`Selector` (`target.go`)** is a set-relative refinement (`refine` + `clause`)
  such as `ExceptMostPowerful`. It narrows the ids _and_ contributes a phrase, so
  niche "compare candidates to each other" rules compose onto any `Target`
  **without a field per rule**. When you are tempted to add another `Target` bool
  for a whole-set rule, add a `Selector` instead.
- **`Count` and `Condition`** are pluggable value/predicate strategies, each with
  paired text (`CountText` / `CondText`). A number that scales with the board is a
  `Count`, not a bespoke effect; a branch is a `Condition` fed to `Conditional`.
- **`CreatureVerb`** is a per-creature verb strategy for `OnChooseCreature`.

Rule of thumb: when behavior varies along an axis, model the axis as a small
strategy interface that also renders its own text — not a new `Effect`/`Target`
field or a `bool`.

## `Resolver` is segregated into role interfaces — add to the right role

The port and its rationale are ADR 0008. In practice it is composed from focused
role interfaces, not one flat list:

`StateReader` (reads) · `EconomyResolver` (Æmber/keys/chains) ·
`CreatureResolver` (per-card in-play state) · `CombatResolver` (damage,
destruction, ability-driven fight/reap/action) · `ZoneResolver` (card movement
between zones + draw) · `TurnResolver` (turn-scoped grants + the lasting
registry) · `ChoiceResolver` (ordering + choosing) · `Logger`.

**When adding a mechanic that needs a new engine capability**, add the method to
the role interface it belongs to (and implement it on `*Game`). Do not append to a
flat list. If a new method fits no existing role, that is a signal a new role —
and probably a new area of the game — is emerging: add a small role interface and
embed it in `Resolver`.

## New whole-tree operations: a type-switch "Visitor", not a new AST method

`Text()` and `Resolve()` are **intrinsic** to a card's identity, so they live on
each `Effect` node. An operation that is **not** part of a card's identity — an
MCTS/AI value estimate, a static analysis like "does this effect ever target the
opponent", serialization — should **not** become a third method on every node.
Adding it that way forces touching ~40 types for a concern cards do not care
about.

Instead, write a standalone function that type-switches over `Effect` in one
file (the Go-idiomatic Visitor):

```go
// aiValue estimates the board value of an effect for the search AI. It is an
// extrinsic operation over the effect AST, kept out of the Effect interface so
// the card vocabulary stays about identity, not heuristics.
func aiValue(e Effect) int {
    switch e := e.(type) {
    case GainAember:
        return e.Amount
    case DealDamage:
        return e.Amount // ...
    default:
        return 0
    }
}
```

`lastingActionOf` in `effect_lasting.go` is already this shape — a centralized
type switch that maps effects to a flat action. Follow it. The tradeoff (a type
switch is not exhaustively checked by the compiler) is handled by keeping each
such operation in **one** file next to a comment listing which effects it covers,
and by the 100% coverage gate forcing every branch to be exercised.

## The lasting registry is the flat-state interpreter

"For the remainder of the turn" behavior (Full Moon, Charge!, Crystal Hive,
Dimension Door) is a second, deliberately smaller interpreter, because flat state
cannot store an `Effect` closure — the decision and its rationale are ADR 0007.
State holds flat `LastingEffect{On Event, Do lastingAction, Controller, Amount}`
records; `lastingActionOf` maps a composed effect to an enum tag and
`game_lasting.go` fires/queries them. A **reaction** runs after an event
(`AddLasting` + `emitLasting`, ordered when several fire); a **replacement**
changes an event's outcome (`lastingReplacement` + `Instead{Of, With}`).

To add one: a reaction on an existing event = support its `Do` in `lastingActionOf`
and `resolveReaction`; a new event = an `Event` value, one
`emitLasting`/`lastingReplacement` call at the site, and its `clause`/`gerund`
text. You never restructure the play/reap hot path. Keep the enum dispatch
centralized.

## Event, ability, and effect verbs: emit → trigger → resolve

Three tiers of verb, kept distinct so a method name says which level it works at:

- **`emit<Event>`** announces a game event and fans out to everything listening —
  `emitCardPlayed`, `emitCreatureEnters`, `emitEnemyDestroyed`, `emitLasting`. Use
  it at an event site that dispatches to responders (triggered abilities and the
  lasting registry).
- **`trigger…`** resolves a *single card's* abilities matching a trigger —
  `triggerAbilities(id, TriggerAfterReap, …)`. The emitters call into it per card.
- **`resolve…`** carries out one specific effect or ability — `Effect.Resolve`,
  `resolveReaction`, `resolveUpgradePlay`.

Reserve the domain word *"fire"* for prose ("a trigger fires", "an ability
fires"): it is what an ability does in response to an emitted event, never the name
of a method that does the emitting.

## Prompts: name the card, and let the player click a target

Every prompt the engine raises reaches a human, so write it as the card's own
sentence and route it through the channel a UI can render as a board interaction.

- **Ask through `pickCreature`/`pickCard`, not `ChooseOption`, whenever the answer
  is a card.** A card choice is made by clicking the card, so a frontend
  highlights the candidates and takes a click. A list of card *names* as buttons is
  a fallback, not the design; only use `ChooseOption` when the options are not
  cards (a house, a key colour, "yes"/"no").
- **Refer to the source card with `SelfName`, never a literal name.** The engine
  substitutes the asking card's name into the prompt (`renderPrompt`), so
  `"fully heal "+SelfName` reads as "fully heal Chuff Ape" at runtime and matches
  the printed text. Hard-coding a name in a prompt is a one-off smell.
- **Phrase the prompt as the card's instruction**, lowercase and imperative
  ("choose a creature to attach {self} to"), so the prompt and the printed text
  are the same sentence.
- **When the choice is optional ("may", "up to N"), the player must be able to
  stop**: the prompt is declinable, and a frontend shows a *Done* button beside
  the highlighted candidates. An optional choice must therefore **not** be
  short-circuited when only one candidate remains — that would take the choice
  away. `pickCreature` auto-takes a sole candidate precisely because it models a
  *mandatory* choice; an optional one needs its own path.

These are conscious tradeoffs; the "why" is in the linked ADRs. Handle them as
described; do not "fix" them into a regression of a constraint.

- **Wide `Resolver`** (ADR 0008). Handled by the role-interface segregation above —
  add to the right role, and let a genuinely new cluster become a new role.
- **`Target` flag-soup + paired `x` / `hasX` fields** (ADR 0005). `Target` must
  stay comparable, which rules out a `[]filter` slice and `*int` optionals.
  Handled by: (a) route per-card filters through the existing builder methods
  (`WithTrait`, `OfHouse`, `PowerAtMost`, …); (b) route **set-relative** rules
  through `Selector` instead of adding another field; (c) only add a new `Target`
  field when a per-card filter genuinely has no `Selector` form. If the field count
  ever truly hurts, the in-constraint move is a single fixed-size comparable filter
  descriptor, **not** a slice.
- **Enum-tagged lasting records instead of stored closures** (ADR 0007). Handled by
  the flat-state-interpreter discipline above.
- **`panic` in `EffectContext.PlayerFor` on `playerUnset`** (ADR 0010). Acceptable
  **only** because `NewCard` runs `validate()` at init, so a real card can never
  reach it — the panic fires at authoring time, not mid-game. Never let a
  _computed_ `Player` reach `PlayerFor`; reject unset at the boundary.
- **`Game` as a large type implementing all of `Resolver`.** A single live-match
  façade, kept in check by the `game_*.go` file split (by area) and the `Resolver`
  role split (by capability). Add a `Game` method to the matching `game_*.go`; if
  an area outgrows its file, split the file, don't grow the type's responsibilities
  silently.

## Adding a mechanic: the decision order

1. Can an existing `Effect` express it by changing a `Target`, `Count`,
   `Condition`, or `Selector`? Prefer that — no new type.
2. Is it a new _node_? Add an `Effect` (or `Condition`/`Count`/`Selector`) in the
   matching `effect_*.go` / `target.go`, with `Text()` + `Resolve()` (+ `validate()`
   if it has an illegal field combo), and a facade alias in `internal/card`.
3. Does it need a new engine capability? Add the method to the right `Resolver`
   role interface and implement it on `*Game` in the matching `game_*.go`.
4. Is it "for the remainder of the turn"? Route it through the lasting registry,
   not the play/reap path.
5. Is it an extrinsic whole-tree operation (AI, analysis)? A type-switch function,
   not a new interface method.

Keep everything green including 100% `internal/engine` coverage (`mage check`).

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
- **Composite — `Sequence`, `Sentences`, `Conditional`, `ChooseHouseThen`,
  `MayRepeat`, …** compose child `Effect`s and recurse `validateEffect` into them.
  Prefer composing small nodes over one fused node (root `AGENTS.md`: "decompose
  fused effects"). The two ordered composites differ only in how they _read_:
  `Sequence` conjoins its children into one compound instruction ("a, and b"),
  `Sentences` renders each child as its own sentence ("A. B."). Pick by how the
  printed card reads. Do **not** wrap individual children to change their
  punctuation — there is no per-child sentence wrapper, and a genuinely mixed card
  nests instead: `Sentences{A, Sequence{B, C}}` reads "A. B, and C."
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
  `FirstChooser` (bot/tests, deterministic), `webChooser`
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

## Reused effect shapes get one shared helper, not a copy per effect

Effects recur in **shapes** — the same little field cluster and the same logic to
turn it into a value, a phrase, or a validation error — across unrelated
mechanics. When a shape shows up in a third effect, factor its logic into one
helper the effects share; each effect keeps its own authoring fields (so the card
API stays ergonomic), but not its own copy of the behavior. The shapes factored so
far:

- **Scale by a `Per` count.** A magnitude that reads "N for each ..." is a base
  times a `Per Count`. `scaled(base, per, ctx)` (`effect_count.go`) is the value
  companion to `forEach`'s text — `Draw`, `GainAember`, `DealDamage`,
  `AddPowerCounter`, `CaptureAember`, and the economy amounts all scale through it.
  `Per` means the same thing on every effect: multiply one target's amount. The
  distinct "re-run choosing a fresh target" axis is `Times` (`Repeat.Times`,
  `CaptureAember.Times`), never an overloaded `Per`.
- **Amount, or an alternative magnitude.** Many effects offer a fixed `Amount`
  or one other way to say the same size — a `By` pool share, an `All`/`Fully`
  whole-quantity flag. `errAmountOr(name, alt, amount, altSet)` (`effect.go`) is
  the shared "set one, not both" validation, beside the `errUnset*` helpers.
- **How much Æmber against a pool.** `StealAember`, `LoseAember`, `CaptureAember`
  (and, minus the share, `GainAember`) express the amount as a fixed base scaled by
  `Per`, or a `By Loss` share (`Half`, `AllBut(n)`, `AllAember`). The economy
  helpers in `effect_aember.go` — `poolAmount(base, by, per, ctx, pool)` for the
  value (uncapped; the caller `min`s against the pool and records the capped
  figure) and `aemberObject(amount, by, possessive)` for the printed object — carry
  that. A `bool` like Capture's `All` is `By: AllAember` internally; keep the
  authoring `bool` only where the printed wording genuinely differs ("captures
  **all** your opponent's Æmber" has no "from" clause).
- **Count what the previous effect produced.** A "for each ... this way" magnitude
  is a `Count` reading a tally the prior effect left on `ctx.Produced` —
  `CreaturesHealed`, `DamageHealed`, `CardsDestroyed`. A new such tally is a field
  on `Produced` plus a small `Count`, never a bespoke fused effect.

Prefer a shared **helper** over a shared embeddable value type. Now that `Per`
means one thing everywhere, `{Amount, By, Per}` genuinely is uniform across
`StealAember`/`LoseAember`/`CaptureAember` (and `{Amount, Per}` in `GainAember`,
which has no pool to take a share of), so a `Quantity{Amount, By, Per}` would fit.
It still does not pay off. Go composite literals do not promote embedded fields, and
the `card` facade aliases these structs (`card.StealAember = engine.StealAember`),
so a `Quantity` would turn every flat authoring site — `{Amount: 2, Per: X}` — into a
nested `{Quantity: Quantity{Amount: 2, Per: X}}` across ~120 card call sites plus
their tests, buying nothing the helpers above do not already carry. Factor the
logic, keep the fields.

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
- **`trigger…`** resolves a _single card's_ abilities matching a trigger —
  `triggerAbilities(id, TriggerAfterReap, …)`. The emitters call into it per card.
- **`resolve…`** carries out one specific effect or ability — `Effect.Resolve`,
  `resolveReaction`, `resolveUpgradePlay`.

Reserve the domain word _"fire"_ for prose ("a trigger fires", "an ability
fires"): it is what an ability does in response to an emitted event, never the name
of a method that does the emitting.

## The game log is typed entries, not sentences

The log narrates what **resolved**, not what a card promised (ADR 0011). There is
no `Logf`: the engine cannot write a sentence into the log at all.

- **A new log line is a new `LogEntry` variant**, a small comparable struct in
  `log_<family>.go` (`log_aember.go`, `log_zone.go`, `log_play.go`, …) carrying the
  ids and amounts of the outcome plus a `Text(Namer) string` that words it. Record
  it with `g.record(entry)` (or `Resolver.Record` from inside an effect). Never
  format a name into a string at the call site — carry the `LocalID` and let
  `Text` ask the `Namer`, so the client can link the card and the log stays honest
  about hidden zones.
- **Attribution is a `Frame`, not a line.** "Because Bumpsy's Play ability said
  so" is not narration — it is context. Open a frame around the resolution
  (`closeFrame := g.openFrame(Frame{Actor, Source, Trigger, Grantor})`, deferred or
  called at the end) and every entry recorded inside inherits it. A card's printed
  text is never a log line.
- **Whole-tree passes over an entry are extrinsic**, the same Visitor rule as
  effects: `RenderEntry` splits a rendered entry into card-linkable segments from
  the outside, in `log_render.go`. Do not add a second method to every entry.
- **A `Text` method reuses the shared phrasing, it does not re-roll it.**
  `log.go` holds the helpers every entry draws on — `namedCards` for a list of card
  names, `because(text, on)` to suffix an event clause, `nameMoved` for a card
  crossing zones — and `text.go` holds `countNoun`/`plural` for "1 card" vs
  "3 cards" and `indefinite` for "a"/"an". Never hand-roll a `card(s)` placeholder,
  a `noun + "s"` plural, or a bare `"a " + noun`. When an entry is another entry
  plus context, build the base entry and render it: `LastingAemberGained.Text` is
  `because(AemberGained{…}.Text(n), e.On)`.
- **Whether a card may be named is decided by its zones, not by the entry.**
  `Zone.public()` says which zones both players can see (discard pile, the board,
  the purged pile — not a hand, archives, or a deck), and `nameMoved(n, id, from,
to)` names a card once either end of the move is public and calls it "a card"
  otherwise. An entry that narrates a zone change states the two zones and calls
  `nameMoved`; it never calls `Namer.Name` directly, so no one entry can leak a
  hand or a deck on its own initiative. A card in a hidden zone becomes nameable
  only through a `Reveal`, which records its own entry.
- Recording is switchable (`SetRecording`), so a search that plays thousands of
  games pays nothing for narration.

## Prompts: name the card, and let the player click a target

Every prompt the engine raises reaches a human, so write it as the card's own
sentence and route it through the channel a UI can render as a board interaction.

- **Ask through `pickCreature`/`pickCard`, not `ChooseOption`, whenever the answer
  is a card.** A card choice is made by clicking the card, so a frontend
  highlights the candidates and takes a click. A list of card _names_ as buttons is
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
  stop**: the prompt is declinable, and a frontend shows a _Done_ button beside
  the highlighted candidates. An optional choice must therefore **not** be
  short-circuited when only one candidate remains — that would take the choice
  away. `pickCreature` auto-takes a sole candidate precisely because it models a
  _mandatory_ choice; an optional one needs its own path.

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
6. Does it need to say what happened? A `LogEntry` variant in `log_<family>.go`,
   recorded at the site the outcome actually lands — never a formatted sentence.

Keep everything green including 100% `internal/engine` coverage (`mage check`).

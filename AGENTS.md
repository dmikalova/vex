# AGENTS.md

Repo-wide guidance for agents. See `internal/cards/AGENTS.md` for
card-authoring specifics.

## Build, test, and lint through `mage`

Run all build/test/format/coverage tasks through `mage`, not raw `go`
commands. The targets wrap the project's conventions (engine coverage gate,
comment/rulebook generation, gofmt), so use them:

- `mage build` — build all packages.
- `mage test` — run all tests.
- `mage vet` — run `go vet`.
- `mage fmt` / `mage fmtCheck` — format, or check formatting without writing.
- `mage cover` — run engine tests and report coverage (kept at 100%).
- `mage generateComments` — rewrite each card's doc comment from its definition.
- `mage gen` — regenerate card comments and the rulebook.
- `mage check` — the full green gate (fmt-check, build, vet, test, coverage);
  run this before considering work done. It must print `ALL GREEN`.

Card research (so you never need a throwaway grep/JSON script):

- `mage lookup "<name>"` — print every source card whose name contains the
  query, with set code, collector number, house/type/rarity, printed text, and a
  ready-made `card.Provenance(...)` call.
- `mage missing <setSlug>` — list the source cards in a set not yet tagged by an
  implemented card (the cards still to implement); slugs match the files in
  `internal/cards/provenance` minus `.json` (e.g. `callofthearchons`).
- `mage coverage` — per-source-set count of cards covered by an implemented
  card's provenance Ref.

Run `mage -l` to see every target.

## Long shell commands may span multiple lines

A long `&&`/`||` command chain does not have to be one physical line. A line that
ends with `&&`, `||`, or `|` continues on the next line, so a multi-step command
can be written as, e.g.:

```sh
mage gen &&
  mage cover &&
  mage check
```

Prefer this layout over one unreadable line when chaining several build/test
steps.

## Multiple agents may be running

More than one agent can be working in this repo at the same time. If the build,
vet, or tests fail because of a change you did **not** make — an unfamiliar file,
a symbol you never touched, an in-progress edit that doesn't yet compile — assume
another agent is mid-change. Wait a little and try again rather than "fixing" or
reverting their work. Only act on failures that stem from your own changes.

## Interpreting requests

Read every request through the lens of **idiomatic, maintainable, composable,
well-written Go that will keep being extended** as more of the game is
implemented. When a request is ambiguous or underspecified, choose the option a
senior Go engineer would find clearest and easiest to build on later — not merely
the shortest path to a passing build.

## Refactoring is welcome

Do **not** contort a new feature to fit the current implementation. If a mechanic
lands more cleanly after reshaping existing code — renaming, splitting, or
generalizing a type, effect, or seam — prefer the refactor. A good fit for the
new feature (and the features that will follow it) matters more than preserving
today's shape. Keep such refactors focused, keep everything green (including 100%
`internal/engine` coverage), and leave the design better than you found it.

## Design for composition, not for the card in front of you

Every card is an instance of a more general mechanic. Implement the **mechanic**,
not the card. Before adding a type, ask: what is the smallest orthogonal piece
this card needs, and how would the *next* card reuse it? Prefer many small pieces
that snap together over one bespoke effect that does everything a single card
happens to want.

Concretely:

- **Decompose fused effects.** A card that "does A, then B if C" is a
  `Sequence` of an effect, a `Conditional`, and a condition — not one
  `DoAThenBIfC` effect. Bonkers Killing Machine is `DiscardTopOfEachDeck` →
  `DestroyOfEachDiscardedHouse` → `Conditional{CardsDestroyedFewerThan, Destroy}`,
  not a single `DiscardAndDestroyByHouse`.
- **Thread state through the context, don't fuse producers and consumers.** When
  one step produces a value a later step consumes (cards revealed, houses
  discarded, creatures healed, a card put in focus as `ctx.It`), record it on the
  `EffectContext` and let any following effect/condition read it. This is how
  `Heal` + `CreaturesHealed`, `RevealTopOfDeck` + `ItIsOfHouse` + `PlayRevealedCard`
  compose without a combined type.
- **Parameterize over enums, not new types.** A new "at least N of a house"
  wants a reusable `CardsPlayed{House, Amount}` count, not a
  `CardsPlayedOfHouseAtLeast`. A new lifetime wants a `Duration` field
  (`EndOfTurn`, `UntilThisLeavesPlay`), not a `…ForRemainderOfTurn` variant. A new
  "which house" wants a `HouseChoice`/reference, not a `…OfChosenHouse` /
  `…OfActiveHouse` pair.
- **Reuse the shared vocabularies.** `Target` already filters by house, trait,
  type, chosen/active house, and set-relative selectors — reach for it (or extend
  it) before inventing a parallel filter. Events (`EventCreaturePlayed`,
  `EventReap`, `EventCreatureDestroyed`, …) with a subject already drive triggers,
  lasting reactions (`ForRemainderOfTurn`), and replacements (`Instead`, `Replace`)
  — extend that spine rather than adding a one-off flag on `CardDefinition`.
- **Self-reference through existing seams.** "Destroy this Upgrade" is
  `Destroy{Target: This}` (the destroy path detaches an attached upgrade), not a
  `DestroyThisUpgrade`. "This creature captures the Æmber" is a replacement of the
  add-to-pool event, not a `CapturesOpponentAember` bool.

A one-off is a smell: if a name encodes a specific card's whole sentence
(`ReadyAndBelongToHouseAfterYouPlayCreature`, `PlayTopOfDeckOfChosenHouse`), it is
hiding reusable pieces. Split it. The cost of a slightly larger card definition is
worth an engine other cards can build on. When a card genuinely needs something new,
add the smallest general primitive and one card that uses it — never speculative,
but always shaped so the next card can reach for it.

**Prefer refactoring cards to improve composability over working around them.**
When a new card almost fits an existing primitive but not quite, reshape the
primitive (and re-express the cards already using it) rather than bolting on a
special case or a parallel one-off. A card whose text shifts slightly to share a
cleaner mechanic is a better outcome than two near-duplicate effects. The tests
are what make this safe: they pin each card's behavior (and rendered text), so a
refactor that breaks something surfaces immediately. Lean on them — refactor
freely, keep everything green, and let the suite catch regressions instead of
avoiding the change.


## Engine design patterns and constraints

The full design ideal for `internal/engine` — the patterns it is built from, the
two hard constraints that shape it, and the deliberate tradeoffs (with how to
handle each) — lives in `internal/engine/AGENTS.md`. Read it before reshaping an
engine seam. The load-bearing rules that affect how you add anything:

- **Effects are an Interpreter AST.** Every `Effect` renders its own text
  (`Text()`) and carries itself out (`Resolve()`), so printed card text can never
  desync from behavior. A new mechanic is almost always a new node in
  `effect_<mechanic>.go`, not a new branch in the `Game` runtime.
- **Vary behavior with a Strategy that also renders its own text.** When behavior
  changes along an axis, model the axis as a small strategy — a `Chooser` (or its
  optional-capability interfaces `OptionChooser`/`Orderer`), a `Selector`, a
  `Count`, or a `Condition` — each of which carries both its behavior and its text
  fragment. Reach for this before adding another `Target` field or a `bool`.
- **The `Resolver` port is segregated into role interfaces** (`StateReader`,
  `EconomyResolver`, `CreatureResolver`, `CombatResolver`, `ZoneResolver`,
  `TurnResolver`, `ChoiceResolver`, `Logger`). A new engine capability is a method
  added to the role it belongs to, not to a flat list.
- **New whole-tree operations that are not part of a card's identity** (AI/MCTS
  scoring, static analysis, serialization) are a standalone type-switch function
  over `Effect` in one file — the Go-idiomatic Visitor — **never** a third method
  on every node. `Text()`/`Resolve()` are intrinsic; heuristics are not.
- **Flat, pointerless, comparable state is non-negotiable.** State holds no
  closures, so "do X later" is flat enum-tagged data (the lasting registry), and
  values compared against state (e.g. `Target`) stay comparable — which is why
  `Target` is a flag struct with paired `x`/`hasX` fields rather than slices or
  pointers. Do not introduce a pointer/slice/map into `GameState` or a
  state-compared value type.


## File organization

`internal/engine` splits by responsibility, and the filename says which:

- `game.go` + `game_*.go` — the `Game` runtime: methods on `*Game`, grouped by
  area (`game_turn.go`, `game_play.go`, `game_read.go`, `game_abilities.go`,
  `game_combat.go`, `game_leaves_play.go`, …). A new `Game` method goes in the
  matching `game_*.go`.
- `effect.go` + `effect_*.go` — the effect AST, one file per mechanic
  (`effect_aember.go`, `effect_damage.go`, …). A new effect goes in
  `effect_<mechanic>.go`.
- Everything else is a small type/data file named after the concept it defines:
  `card.go`, `state.go`, `types.go`, `target.go`, `duration.go`, `destination.go`,
  `resolver.go`, `text.go`. A new enum or value type gets its **own**
  `<concept>.go` — e.g. destinations live in `destination.go`, not
  `game_destination.go`; the `game_` prefix is only for `Game` methods.

Tests live beside their source in `<name>_test.go`. The `card` facade mirrors the
engine's namespace/type files: `duration.go` → `card/duration.go`,
`destination.go` → `card/destination.go`.

## Lasting "remainder of the turn" effects (event-driven, never hardcoded)

Some effects last "for the remainder of the turn" and attach to a later game event
— Full Moon gains Æmber whenever you play a creature, Charge! deals damage whenever
you play a creature, Crystal Hive gains Æmber whenever a creature reaps, Dimension
Door makes reaping steal instead of gain. Do **not** add a bespoke `if g.State.Foo…`
block to the play or reap path (`PlayCreature`, `reapWith`) for each of these — that
hardcodes every effect into the hot path and makes effects sharing a timing window
impossible to order. Instead, everything routes through the flat registry in
`game_lasting.go` (`AddLasting(on Event, do lastingAction, controller, amount)`).
Because the state is a flat, pointerless value it cannot hold effect closures, so a
lasting effect is a small enum-tagged record `LastingEffect{On Event, Do
lastingAction, Controller, Amount}`.

There are two flavors:

- A **reaction** runs *after* an event. The event site emits **one** dispatch —
  `g.fireLasting(EventCreaturePlayed, actor, subject)` / `g.fireLasting(EventReap,
  …)` — which gathers every reaction the actor owns for that event and, when several
  fire at once, lets the controller **order** them (KeyForge lets the active player
  order simultaneous triggers), resolving each via `resolveReaction`. So "gain Æmber
  after you play a creature" and "deal damage after you play a creature" order for
  free. Authored as `card.ForRemainderOfTurn{On: card.Event.CreaturePlayed, Do:
  card.GainAember{...}}` (Do is a small composed effect — `GainAember` or
  `DealDamage` to an enemy creature).
- A **replacement** changes an event's *own outcome* before it happens. The event
  site queries the registry (`g.lastingReplacement(player, EventReapAember)`) and
  applies the replacement in place — `gainReapAember` steals instead of gaining when
  Dimension Door is active. Authored as `card.Instead{Of: card.Event.ReapAember,
  With: card.Steal}`.

Adding a reaction on an existing event = supporting its `Do` in `lastingActionOf` +
`resolveReaction`; a new event = an `Event` value, one `fireLasting`/
`lastingReplacement` call at that site, and the `clause`/`gerund` text. You never
touch the play/reap path's structure. `EndTurn` drops a player's entries via
`clearLasting`.

## Constant-granted abilities

When a card in play says that every creature gains an ability, author it as a
`ConstantAbility` with a `Target` and `Granted` abilities. Annihilation Ritual is
the model: it grants each creature the ability that purges that creature when it
is destroyed.

```go
card.WithConstantAbility(card.ConstantAbility{
  Target: card.Target.EachCreature,
  Granted: []card.Ability{{
    Trigger: card.Trigger.Destroyed,
    Effect:  card.PurgeCreature{Target: card.Target.This},
  }},
}),
```

The engine gathers every Destroyed ability for all creatures being destroyed,
then lets the active player order them. A `PurgeCreature{Target: card.Target.This}`
ability takes its creature out of play immediately, so that creature's remaining
Destroyed abilities do not resolve and final destruction cleanup does not move it
to its discard pile. Never implement this kind of card as a global override in
`discardDestroyed` or another leave-play path.

## KeyForge vernacular

Names — types, methods, fields, effects, targets, everything — must stay within
KeyForge's own vocabulary rather than generic gaming terms. Say
`ExceptMostPowerfulCreature`, not `ExceptStrongest`; `CannotFight`, not
`PreventFight`. When you need the right word, source it in this order:

1. The **provenance files** (`internal/cards/provenance`) and the original card
   text they point at — the canonical wording.
2. **Existing implementations** in this repo — reuse an established term instead
   of coining a synonym.

If neither has a term for the concept, pick the phrasing closest to how KeyForge
cards are actually written, and flag it for the user.

## Writing abilities (card authoring)

Author cards through the `card` facade in the multiline ability style —
`card.WithAbility(` breaks onto its own line, the trigger and effect share the
next line, and each effect struct field goes on its own line when there is more
than one (single-field effects stay inline):

```go
card.WithAbility(
	card.Trigger.Play, card.CannotFight{
		Player:   card.Opponent,
		Duration: card.Duration.NextTurn,
	}),
```

See `internal/cards/AGENTS.md` for the full card-authoring guide (file layout,
generated doc comments, wording rules, and tests).

# Flat, pointerless, comparable game state

## Context

The engine is built to be searched: an MCTS/self-play bot needs to clone a
position, explore, and discard it, millions of times. A conventional mutable
object graph — `GameState` as a web of pointers, slices, and maps — makes each
clone a deep copy with allocation and GC pressure, and makes two positions
impossible to compare with `==`. That cost sits on the hottest path in the whole
system.

## Decision

`GameState` is a **value struct of fixed-size arrays** (`[maxCards]CardCore`,
`[2]Zone`, …) with **no pointers, slices, or maps**, so `GameState.FastCopy()` is
a plain value copy with no allocation. Cards are referenced by `LocalID uint8`
into a read-only `catalog` of definitions held **separately** from `GameState` and
shared across every clone. Value types that are embedded in or compared against
state stay **comparable**: `Target` (compared `== Target{}` by `ConstantAbility`)
is a wide flag struct rather than a `[]filter`, and an "optional int" is a paired
`x int` + `hasX bool`, never a `*int`.

### One comparable `Filter` narrows every set of cards

Every narrowing a card writes is one comparable `engine.Filter`, exposed as
`card.Filter`: the axes a card must satisfy — its type, house, traits, name,
power, state, and place in a battleline — in one struct literal. `Target` is
`Target{Kind, Filter, Refinement, Neighbors}`, and the nodes that point at a card
in a pile (`Search`, `DiscardUntil`) or at the creatures a card refuses damage
from take the same `Filter`.

A `Filter` holds **no slice and no pointer**, which is the constraint above
applied to it rather than a separate rule. That is what removed `CardFilter.Or`:
its one consumer, Chief Engineer Walls, writes `MatchAny: true` instead. An
optional number is still a paired kind-and-number value, now a named comparable
type — `PowerBound{Kind, Amount}` rather than an `x`/`hasX` pair, which is the
same decision with a name on it.

**Two unions, told apart by the sentence they print.** `MatchAny` is the union
_inside one noun phrase_, joined with "or": "each Mars or Robot creature" (EMP
Blast), "Mutant creatures or creatures with power 5 or higher" (Ardent Hero). It
spans every axis, not only the identity ones. `AnyOf`, a `Refinement`, is the
union of _separately quantified phrases_, joined with "and": "each Dinosaur
creature and each creature with power 6 or higher" (Regrettable Meteor). Do not
collapse them into one mechanism: the difference is printed text, not logic.

**Filter is a subtype of Refinement, not its sibling.** A `Filter` answers
"does this one card match"; a `Refinement` answers "which of these candidates do
I keep". Every per-card test is trivially a set rule, so `Filter` implements
`Refinement`; the reverse is false, because `MostPowerful`, `KeepPerSide`,
`PortionPerSide` and `SamePowerAsChosen` need the whole candidate set — to
compare candidates, or to prompt across them — and because `AnyOf` holds a slice
and `PowerLessThan` holds a `Count` interface, so a `Filter` holding one would
stop being comparable.

The line: **a filter is any test decidable on one candidate at a time, however
much board it reads; a refinement is a rule that needs the whole candidate set.**
Reading the board is not the test — `SharesTrait`, `SharesHouseWithNeighbors`,
`HouseWithMostCreatures`, `HouseWithAtLeast` and `WithoutSharedTrait` all consult
the board and are filter axes. A bound with a literal number, or no number at
all, is a filter kind (`Power.AtLeast(6)`, `Power.LessThanSource()`); a bound
whose threshold is a live `Count` stays a refinement (`PowerLessThan`). A bare
`Filter` in `Target.Refinement` is rejected by `Target.validate`, naming the
`Filter` field: a filter is a refinement only inside a union combinator, which is
the only place the interface earns anything.

**Point at the card.** A node takes a `Filter` when it selects from, counts, or
tests _actual cards that exist somewhere_ — can you point at the card the field
is deciding about. `CardsInPlay` can (a battleline), `Search` can (a deck),
`ItIs` can (the card in context). A `Type` field that only names a noun is not a
filter: `CardsPurged` reads an integer tally and never looks at a card.

**A node that arms a rule into game state keeps flat scalar parameters.** The
reason is size and staleness, not typing. `GameState` is ~4024 bytes and the
project has deliberately shrunk it (4232 → 4112 → 4024) because one snapshot is
one undo and one MCTS node; it holds eight lasting slots and four play bars, so
arming a ~25-field `Filter` in place of the 24 bytes `House`+`Type`+`Trait` take
today would add over a kilobyte — a third again on the whole state. And `Filter`
carries context-relative axes (`Houses.Contextual`, `Except.It`, `Except.Focus`,
`SharesTrait`, `Neighboring`, the position axes) that mean nothing once stored
and evaluated on someone else's turn. This covers `CannotPlay`,
`PlayersCannotPlay`, `NextPlayed`, `DamageOthersAfterUsingTrait`,
`CreatureBar.Houses` and the lasting-reaction record. **Revisit when state size
allows**: `Bar[T comparable]` would already accept a comparable `Filter`, so the
objection is cost, not the type system.

**Filter renders one noun phrase; the consumer supplies the quantifier.**
`Filter` owns the adjectives and the trailing clauses, and splits its rendering
in two so the consumer can write between the halves: `Target` prints
"each …"/"a …"/"another …", and a consumer that prints its subject in one piece
asks for the whole noun instead.

House matching inside a `Filter` is the one role-typed `HouseMatcher` of ADR
0038, which notes that `Target` needed a manual `houseReplaced` only because its
house fields were unexported. They are exported now, so the self-house reflection
pass descends into a `Target` like any other struct and that method is gone.

## Consequences

- `FastCopy` is allocation-free, so cloning a position for search is essentially
  free — the property the rest of the engine bends to preserve.
- State **cannot hold a closure or an `Effect`/`func` value**, so "do X later" has
  to be flat, enum-tagged data. That constraint is what forces the lasting-effects
  registry (ADR 0007).
- `Target` is `Target{Kind, Filter, Refinement, Neighbors}`: a comparable `Filter`
  holds the per-card axes, and set-relative rules go through a `Refinement`, never
  a new slice field. A new axis is a field on `Filter`, not a new `Target` field.
- `Refinement` is the one field that keeps a `Target` from being unconditionally
  comparable, because `AnyOf` holds a slice. Isolating the hazard to a single
  field is the point: everything else in a `Target` is a value.
- `engine.CardFilter` and `engine.DamageSourceMatcher` are gone, folded into
  `Filter`, so there is one narrowing vocabulary rather than three that drift.
- Fixed capacities mean zones are sized to their own bounds (ADR 0002) and an
  unbounded per-card collection (a creature's upgrades) needs an intrusive linked
  list rather than a slice (ADR 0001).
- Keeping definitions in a separate read-only catalog keeps them out of every
  copy.

This is the keystone constraint of the engine: ADRs 0001, 0002, and 0007 are all
consequences of it.

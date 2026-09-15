# 38. House matching is one role-typed `HouseMatcher`, rendered as a noun qualifier

## Context

"Which house(s) does this clause admit" is asked all over the engine — a card
played from hand, a creature chosen on the board, a card revealed, a house counted,
a permission granted. Today the same question is answered in at least three
unrelated shapes:

- **`House` + `Except bool`** — `PlayFrom` and `PlayOrUse` carry a house plus a
  boolean that _inverts_ it, so `House` means "Star Alliance" when `Except` is false
  and "every house but Star Alliance" when it is true. The field's meaning flips
  under a sibling flag. `validate()` has to reject the "`Except` set, no house"
  combination that the type still lets an author write.
- **`House` + a sibling `ExceptHouse`** — `Target`, the `Chosen`/`Each`/`Random`
  selections, and the discard/reveal effects carry two house fields, at most one of
  which is meaningful. Nothing at the type level stops both being set.
- **`HouseSelector`** — `MayPlayOrUse` (ADR 0037) already folded its house axis into
  one flat, comparable value (`Kind` + `House`) covering named / chosen / any /
  every-house-but-one / houses-you-control. This is the shape the other two want to
  be, but it also carries a grant-only kind (`SelectControlled`) and renders whole
  grant _clauses_ ("play or use a non-Star Alliance card"), not the woven noun
  qualifier a `Target` needs ("a friendly **Mars** creature").

So the same concept has a boolean-inversion form, a two-field form, and a good
composable form that is welded to grants. Each new house-filtered wording picks one
of the three, and the "non-`<house>`" text is re-derived in every effect's `Text()`.
This is the sister fragmentation to ADR 0037: that ADR unified the house-**grant**
axis; this one unifies the house-**match** axis that filters, targets, and predicates
all share.

## Decision

**A single flat, comparable `HouseMatcher` is the one way to say "cards of / not of
/ any particular house", and it renders as a noun qualifier that each effect wraps in
its own phrase.** The grant selector becomes a distinct _superset_ type, so a filter
cannot express a grant-only kind, and the one set-relative house notion stays a
`Refinement`.

```go
type HouseMatcher struct {
    Kind  HouseMatchKind // Any (zero) / Named / Except / Chosen / Active / Contextual
    House House          // the named or excluded house; SelfHouse resolves at build time
}
```

- **The zero value means "any house" — no restriction.** A house filter is usually
  absent (most `PlayFrom` name no house), so the permissive case is the one an author
  leaves unwritten. This is the opposite of the grant selector, whose zero is
  _invalid_ because a grant must always name whose cards it frees — a difference that
  is exactly why the two are different types.
- **`Kind` covers the per-card notions:** a named house (`Named`), every house but
  one (`Except`), the enclosing `ChooseHouseThen`'s house (`Chosen`), the active
  house (`Active`), and the house of the card in context / `ctx.It` (`Contextual`).
  The `SelfHouse` sentinel is carried in the exported `House` field, so
  `resolveSelfHouse` rewrites it by reflection with no per-type seam (ADR/self-house):
  `Target` needed its manual `houseReplaced` only because its house fields were
  unexported.
- **`HouseMatcher` renders a noun qualifier, not a clause.** One method turns the
  matcher into "Mars", "non-Sanctum", "of the chosen house", "of that card's house",
  or "" (Any), and each consuming effect welds that into its own subject —
  `PlayFrom` into "play a **non-Star Alliance** card", `Target` into "a friendly
  **Mars** creature". The house wording lives in one place instead of being
  re-derived per effect (ADR 0006 still holds: the matcher renders itself).
- **The grant selector is a distinct superset.** `MayPlayOrUse` keeps its own
  selector type, which carries a `HouseMatcher` plus the grant-only `Controlled`
  ("every house you have a card in play for"). Because a filter or target field is
  typed `HouseMatcher`, it _cannot_ hold `Controlled` — the illegal combination is
  unrepresentable at the type level rather than caught by a lint (the type-split
  chosen over one shared struct guarded by per-effect `validate()`).
- **`houseWithMostCreatures` stays a `Refinement`, not a matcher kind.** It compares
  houses to each other by creature count, so it is a set-relative rule (like
  `MostPowerful`), not a per-card test. Folding it into `HouseMatcher` would make the
  matcher mean "any narrowing" and lose the line ADR 0006's `Refinement` draws.

`HouseMatcher` is adopted **everywhere a card's house is matched against a
criterion**: the pile-selection effects (`PlayFrom`, `PlayOrUse`, `Chosen`/`Each`/
`Random`, `Search*`, `Reveal`, `Discard`, `PutFromHand`, `PutUnderFromHand`), the
in-play `Target`, and the read-only house predicates (`CardsPlayed`, the discard-count
conditions). Grants keep the superset. The rollout is staged — introduce the matcher
and its renderer, then convert one family at a time, green between each — but the end
state is one house vocabulary.

The card-**type** axis has no matching problem: `CardTypes` is already a set-of-types
bitset and any subset is valid for any effect, so this ADR touches only houses. The
one related type cleanup rides along: `PlayFrom`'s split `Type CardType` /
`Types CardTypes` pair collapses to a single `Types CardTypes` set, and the named
`card.Types.NonCreature` / `card.Types.Artifacts` aliases are deleted in favor of an
explicit `card.Types.Of(card.Type.Artifact, …)` so a card spells out the types it
admits.

## Consequences

- A new house-filtered wording is a `HouseMatcher` value, not a new field pair or a
  new inverting boolean. `Com. Officer Kirby` reads
  `House: card.Houses.Except(card.House.Self)` /
  `Types: card.Types.Of(card.Type.Artifact, card.Type.Upgrade, card.Type.Tactic)`
  instead of `House: card.House.Self, Except: true, Types: card.Types.NonCreature` —
  the definition now reads as the printed card does.
- Illegal states are unrepresentable: no "`Except` set with no house", no "both
  `House` and `ExceptHouse` set", and no grant-only `Controlled` on a filter — each
  was previously either a `validate()` rejection or nothing at all.
- The "non-`<house>`" / "of the chosen house" wording is rendered once, by the
  matcher, so every filtered subject across the engine phrases houses identically and
  cannot desync (ADR 0006).
- Rejected — **one shared struct guarded by `validate()`** (Q9 option A): keeps a
  single type but lets a filter hold `Controlled` and be rejected only at init
  (ADR 0010). The role-typed split makes the grant/filter boundary structural, which
  is the whole reason the user asked to "split the types so a method can only have the
  correct ones."
- Rejected — **reuse the `Refinement` combinators (`card.Not(house)`)** for house
  negation: `Not`/`AnyOf` are set-relative rules over creatures; overloading them for
  a per-card house test would erase the concept `Refinement` names.
- Rejected — **leave `Target` on its own two-field house filter**: doing the
  pile-selection effects now and `Target` later means building the shared
  noun-qualifier renderer and the matcher type, then a second migration to reuse them.
  One vocabulary in one ticket avoids the duplicated effort.

# 37. Out-of-house permission is one `MayPlayOrUse` grant, not a node per wording

## Context

A family of cards lets a player act with cards outside their active house for the
turn — fight, reap, use, or play from hand. Each printed wording grew its own
effect node, its own resolver method, and (for the ones that log) its own record:

- **`GrantFight{House}`** — "each friendly creature of that house may fight"
  (Brothers in Battle reads the chosen house, Signal Fire names Brobnar). Backed by
  `GrantFightForHouse`, recorded as `FightGrantedForHouse`, held in
  `State.MayFightHouse[player]`.
- **`GrantFightAnyHouse{}`** — "each friendly creature may fight" (Follow the
  Leader, Horseman of War). Backed by `GrantFightAnyHouse`.
- **`MayActFriendlyHouse{House, Grant}`** — "you may play or use a Mars card" /
  "you may use friendly Sanctum creatures" (the Ambassador cycle plays and uses;
  Sigil of Brotherhood and Ritual of the Hunt only use). Backed by
  `GrantPlayForHouse` and `GrantUseForHouse`, held in `State.MayPlayHouse[player]`
  and `State.MayUseHouse[player]`.
- **`MayPlayOffHouse{Except, Controlled, NotType, Grant, Count}`** — the Star
  Alliance "play a non-Star-Alliance card" cycle (United Action, Commander Kirby,
  CXO Taber). It frees cards by _exclusion_ (every house but one, or every house you
  control a card in), bounds a _count_, and restricts a _type_.
- **`MayUseFriendlyArtifacts{}`** — Scientifical Hack frees _every_ friendly
  artifact whatever its house. Backed by `GrantUseArtifactsAnyHouse`.

These are the same mechanic — a this-turn permission to treat an out-of-house card
as if it were in the active house — expressed five ways. Fight is the narrow case of
use; use and play are two halves of the same "act with this house" grant that
`MayActFriendlyHouse` already carries as a `HouseGrant` bitset. `MayPlayOffHouse`
selects its houses by exclusion instead of by name and bounds a count, and
`MayUseFriendlyArtifacts` restricts the grant to artifacts — but both are still "for
the rest of the turn you may act with cards that are not your active house." A sixth
wording ("that house may reap", "you may play a non-Brobnar creature") would be a
sixth node, and the resolver/log/state families keep growing one method per wording.
This is the same fragmentation ADR 0035 found in the forward house-**choice**
constraints (three paired slots → one table); this ADR is its sister for the
house-**permission** grants.

## Decision

**Out-of-house permission grants are one effect node, `MayPlayOrUse`, over four
axes**, backed by one resolver-method family and one log record.

```go
type MayPlayOrUse struct {
    Houses HouseSelector // named / chosen / any / all-but-a-named / houses-you-control
    Grant  HouseGrant    // GrantPlay | GrantFight | GrantUse (a bitset)
    Types  CardTypes    // zero value = all card types; narrow to creatures or artifacts
    Count  int           // 0 = unlimited; N bounds how many cards the grant frees
}
```

- **`Houses`** selects whose cards the grant frees: a named house (Signal Fire's
  Brobnar), the chosen house (Brothers in Battle, `HouseNone` reading
  `ctx.ChosenHouse`), any house (Follow the Leader), every house but a named one
  (the Star Alliance exclusion), or every house you control a card in. This is the
  axis that `GrantFight`/`GrantFightAnyHouse` split into two nodes and that
  `MayPlayOffHouse` carried as its `Except`/`Controlled` pair.
- **`Grant`** is the existing `HouseGrant` bitset, extended with `GrantFight` beside
  `GrantPlay` and `GrantUse`. `GrantFight` is the narrow "fight only" case; `GrantUse`
  remains the full "fight, reap, or Action:" use. The bitset renders the KeyForge
  clause — "may fight", "may use", "may play or use".
- **`Types`** restricts the card types the grant reaches. Its **zero value means all
  types** — the common case ("a Mars card", "a non-Star-Alliance card") frees
  everything, so a card must _opt in_ to a narrower subject (creatures only, or
  artifacts only). This axis absorbs `MayUseFriendlyArtifacts` (any house +
  `GrantUse` + artifacts) and `MayPlayOffHouse`'s `NotType`. It is a new bitset type
  **`CardTypes`** (a set of allowed types); the helper that lists every real type,
  formerly `CardTypes()`, was renamed to the unexported `allCardTypes()` to free the
  name. `OffHousePermit.NotType` (a single excluded `CardType`) became
  `OffHousePermit.Types CardTypes` to store it, so the permit frees a card only when
  its type is admitted rather than not-excluded.
- **`Count`** bounds how many cards the grant frees; **zero means unlimited**, the
  common case. Only the Star Alliance cycle sets it.

The node folds `GrantFight`, `GrantFightAnyHouse`, `MayActFriendlyHouse`,
`MayPlayOffHouse`, and `MayUseFriendlyArtifacts` into one. Its `Resolve` records the
grant through a single resolver seam (`GrantMayPlayOrUse(player, Houses, Grant,
Types, Count)`), and a single `MayPlayOrUseGranted{Player, Houses, Grant, Types,
Count}` log record narrates it — replacing `FightGrantedForHouse` and its siblings.
The per-player state slots (`MayFightHouse`, `MayUseHouse`, `MayPlayHouse`, the
off-house permit, the any-house artifact flag) remain flat, comparable state
(ADR 0005); the ready phase clears them exactly as today.

The name is **`MayPlayOrUse`** — explicit about the two verbs the node can free, so
a reader never has to know that "use" silently also means "play from hand." The
rendered card text still narrows to the KeyForge clause the axes select ("may fight",
"may use", "may play or use"); only the Go node name carries the full span.

## Consequences

- A new out-of-house wording ("that house may reap", "play up to two non-Brobnar
  cards") is a new `Grant` bit, `HouseSelector` case, or `Count`/`Types` value — not
  a new node, resolver method, and log record. The wager/must/cannot table of
  ADR 0035 and this grant node together cover the whole house-permission surface with
  two small systems instead of a dozen paired slots.
- `MayPlayOrUse`'s `Text()` renders every combination from the four axes, so the
  printed clause can never desync from the grant (ADR 0006).
- Because `Types` and `Count` default to "all" and "unlimited", the plain grants
  (Signal Fire, Follow the Leader, the Ambassador cycle) name only `Houses` and
  `Grant`; the Star Alliance cycle and Scientifical Hack are the same node with two
  more fields set, not separate types.
- Rejected — **keep the nodes split**: leaves fight outside the bitset and exclusion
  and artifacts in their own types, so the next wording is the next node and the
  resolver/log/state families keep growing one method per card. The whole reason this
  reached a decision.
- Rejected — **name it `MayUse`**: KeyForge's own verb for acting with an in-play
  card, but it hides that the node also frees _playing_ from hand. `MayPlayOrUse`
  spells out the coverage; the card text still uses the shorter KeyForge clause.

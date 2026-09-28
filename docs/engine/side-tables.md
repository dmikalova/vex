# Engine side-tables: house constraints, counters, control

Three mechanics keep flat, fixed-cap tables on `GameState` instead of a field
per card or per kind (ADR 0005). A new mechanic of the same kind extends its
table.

## Forward house choice is one delayed constraint table

Cards that reach forward to a player's **next house choice** — Control the Weak
and Snag (must a house), Tezmal and Snag's Mirror (cannot a house), Snaglet (a
wager on the house) — do **not** each grow a state slot and resolver method.
They arm entries in `State.HouseConstraints[player]` / `HouseConstraintCount`
for this turn and `…Next` for the pending one (`maxHouseConstraints`).
`StartTurn` promotes the starting player's `…Next` entries and clears them
(ADR 0035).

- A `HouseConstraint` is a comparable enum-tagged record (`Kind`: must-house,
  must-creature, cannot-house, or wager).
- `ChooseHouse` resolves the **whole table at choice time** through
  `allowedHouses`: start from `choosableHouses` (identity houses plus the
  houses of cards the player controls in play in their own right, read now),
  remove every cannot (cannot overrides must), and if any must survives,
  restrict to the surviving musts. An empty set means **no active house**, a
  valid outcome chosen explicitly as No House.
- A must-creature entry stores a `LocalID` and reads that creature's house
  **now**, so Snag respects house changes between the fight and the choice.
  Never freeze a house at arm time.
- Wagers are scanned by `resolveHouseWagers` at the same choice and pay their
  predictor when the chosen house matches.
- Two reads stay separate: `AllowedHouses` (choosable minus cannot) drives the
  choice and the client's picker, so a controlled off-house creature widens
  the choice; `PlayerHasHouse` stays **identity-only**, because
  `ItIsOffIdentity` (Sneklifter) means "on the identity card" and a controlled
  off-house card must not count as in-identity.
- A new forward-choice mechanic is a new `Kind` plus its handling in
  `allowedHouses`/`resolveHouseWagers`, not a new state slot.

## Generic counters are a global side-table, not a field per kind

Card-placed markers that matter only to the cards that read them (a doom
counter on Wretched Doll, and its kin) live in
`Counters [maxCounterEntries]CounterEntry` + `CounterCount` on `GameState`.
`CardCore` is copied whole into every snapshot, so a field per kind would bloat
`FastCopy` for markers that are almost always absent (ADR 0024).

- To add a counter: one `CounterKind` value in `counter.go` with its display
  noun. No rulebook term (the one "Generic Counters" entry covers them all), no
  state field, no per-kind resolver method.
- Place and read through `PlaceCounter{Kind, Target, Amount}`, the
  `CounterInPlay{Kind}` condition, the `Target.WithCounter(kind)` filter, and
  the `PlaceCounter`/`CountersOn` resolver methods. A per-card count folds into
  its entry's `N` (saturating); a card sheds every entry through the one
  `removeFromPlay` funnel.
- Do **not** add a `FooCounters int16` to `CardCore`. Power counters, damage,
  and Æmber-on-card stay bespoke: they are read on the hot path and change a
  creature's power or fate.

## Take control is a LIFO stack, not one source per card

Two abilities can seize the same card; when the newer ends, the older (if still
in effect) takes over, and only when the last ends does the card revert to its
owner. The stack is `Controls [maxControlEntries]ControlEntry` + `ControlCount`
(`game_control.go`); a card's controller is the top entry naming it, and
`ControlPlus` on `CardCore` is only the O(1) cache of that top, refreshed on
every push and pop (ADR 0028).

- Take control by pushing one entry through
  `takeControl(id, controller, source)`, for creatures and artifacts alike. `source` is the card whose
  leaving play ends the grant; an **UntilCardLeavesPlay** grant names the
  seized card itself.
- A take first supersedes its own source's earlier entry on the card
  (`supersedeControl` drops the matching `(card, source)` pair), so re-taking
  replaces rather than stacks, bounding the stack by the board.
- Every exit funnels through `removeFromPlay`, which calls
  `releaseControlHeldBy(id)` and `clearControls(id)`.
- Do **not** add a second control field to `CardCore` or special-case
  artifacts; a permanent grant is a self-sourced stack entry.

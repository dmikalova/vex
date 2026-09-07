# 27. Houses are ordered alphabetically

## Context

`House` is a small `uint8` enum (`internal/engine/types.go`) that lives inside the
flat, value-copyable `GameState`. Its members were originally listed in the order
the houses were released — Brobnar, Dis, Logos, Mars, Sanctum, Shadows, Untamed,
then Saurian and Star Alliance appended as later sets added them. That release
order carried no meaning in the code: nothing keys off a house's numeric value
except state that is indexed by house (for example
`PlayPermissionsUsedThisTurn [2][NumHouses]uint8`), and the printed name is looked
up through `houseNames`, so the enum's order is invisible to players.

Appending each new house kept the enum from ever reading in a predictable order,
which made the list harder to scan and left `houseNames` easy to knock out of sync
with the enum (a real bug: the array was left in release order after Saurian and
Star Alliance were slotted into the middle, so names no longer lined up with their
enum values).

## Decision

List the houses **alphabetically**: Brobnar, Dis, Logos, Mars, Sanctum, Saurian,
Shadows, Star Alliance, Untamed. `HouseNone` stays first (value 0, the unset
zero), and `SelfHouse` stays last as the resolve-away sentinel that never indexes
state. `houseNames` is kept in the same order as the enum, and `NumHouses` is
`int(Untamed) + 1` — Untamed is now the last real house.

Alphabetical order is arbitrary but stable and predictable: a new house slots into
its alphabetical place, the list always reads the same way, and `houseNames` has
one obvious correct order to match.

## Consequences

- **The enum's numeric values are the on-disk order of persisted state.** State
  arrays indexed by house (`PlayPermissionsUsedThisTurn`, and any future
  `[NumHouses]` array) serialize in enum order, so reordering the enum invalidates
  every stored snapshot. The web client's `snapshotVersion` is bumped (11 → 12) so
  a stale snapshot is flushed and re-dealt rather than restored against the new
  ordering.
- Inserting a future house at its alphabetical position shifts the values of every
  house after it and so is itself a state-version bump — the same cost as any
  enum reorder. This is acceptable: houses are added rarely, and the client
  already flushes on a version bump.
- Nothing in the rules or rendering depends on house order, so no card behavior or
  printed text changes.

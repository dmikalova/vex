# Invalid-zero sentinels validated at card init

## Context

Effects and targets carry required enum fields — `Player`, `Trigger`, target kind,
`Duration`, `Destination`, `Zone`, `Event`. Every Go enum has a zero value, and if
that zero means a real option, an author who forgets a field gets silent,
wrong-but-plausible behavior that surfaces only mid-game, far from the mistake.

## Decision

Give each such enum an **invalid-zero sentinel** (`playerUnset`, `triggerUnset`,
`targetUnset`, `durationUnset`, `zoneUnset`, `eventUnset`, …) rather than letting
zero be a valid option. `NewCard` runs `validate()` at initialization and **panics
if a required field is unset**, so the failure fires at authoring/startup time, not
during a game. A residual guard (`EffectContext.PlayerFor` panics on
`playerUnset`) covers the otherwise-impossible case; a _computed_ `Player` must
never be allowed to reach it — unset is rejected at the boundary, in `validate()`.

## Consequences

- A forgotten required field is caught deterministically at init, as a loud panic
  with the card in hand, instead of a silent zero-value default deep in a match.
- The mid-game panic guard is acceptable **only** because `validate()` makes it
  unreachable for any real card; it is a belt-and-suspenders assert, not error
  handling.
- `validate()` runs once per card at construction (`NewCard`) — at startup, or at
  deck-generation time for a materialized card — never in a turn or an MCTS
  rollout, so it is off the hot path and stays on in production (it also fail-fasts
  on malformed generated cards). The runtime state checks that *would* be hot-path
  expensive are the separate, `-tags assert`-gated `assertInvariants`
  (`assert_on.go` / `assert_off.go`), which compile to a no-op in production.
- This is a uniform convention: a new required enum field should add its own
  sentinel and validate it at init.

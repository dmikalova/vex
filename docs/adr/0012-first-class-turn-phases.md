# Turn phases are first-class engine state

## Context

The turn lifecycle was three coarse methods that each fused several parts of a
turn: `BeginTurn` reset state and forged a key, `ChooseHouse` set the active house
and offered the archives, and `EndTurn` resolved end-of-turn abilities, readied,
and drew. Between `ChooseHouse` and `EndTurn` there was no engine concept at all —
the frontend simply called play and use in a loop.

Three separate requirements converge on the same missing concept:

- **The game log** (ADR 0011) must demarcate each part of a turn. Half of the
  demarcations had no call site to hang on, and a turn where the player plays
  nothing would have shown no play boundary at all.
- **"At the start of your turn" abilities** had nowhere to resolve. There was no
  such trigger, and forging happened before any such ability could have run.
- **An effect that ends the current phase early** (Omega, from a later set) cannot
  be expressed without a phase to end.

A fourth pressure is latent: forging was buried inside `BeginTurn`, leaving no
seam for an effect that cancels a forge.

## Decision

A turn is eight ordered **phases**, and the current one is a `Phase` value in
`GameState` — a `uint8` with an invalid zero (ADR 0010), so it stays flat and
comparable (ADR 0005):

1. Start of turn — "at the start of your turn" abilities resolve
2. Forge — the mandatory forge, after those abilities
3. Choose a house
4. Archives
5. Play
6. Ready
7. Draw
8. End of turn — "at the end of your turn" abilities resolve

The engine advances phases and **blocks only on input**. Start of turn, forge,
ready, draw, and end of turn have no player decision and run to completion when
entered; choose a house, archives, and play wait for the frontend. Ending a phase
early is a flag the phase loop reads, so an effect skips the remainder of the
current phase without special-casing any one phase.

KeyForge's rulebook calls these divisions "steps". Vactrol calls them **phases**
everywhere — engine identifiers, generated rulebook, rendered card text, and the
game log — rather than carrying two words for one concept.

## Consequences

- `BeginTurn`, `ChooseHouse`, and `EndTurn` split apart, changing the public turn
  API and every test that drove a turn through them.
- Forging gains a seam it did not have, so a forge-cancelling effect becomes
  expressible rather than needing a special case.
- `TriggerStartOfTurn` exists and has a defined place to resolve.
- The resolution order changes; see ADR 0013.
- The frontend gains a real turn skeleton to render, which is what lets the client
  group the log by phase.

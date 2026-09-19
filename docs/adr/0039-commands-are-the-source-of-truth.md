# 39. Commands are the source of truth; state, log, and undo are replay projections

This decision records the target architecture for match persistence. It guides a
staged refactor. The engine seam that realizes it now exists — the
`Command`/`Request` step function and `Stepper` (suspend.go), the
`ApplyAction`/`ApplyManual`/`LegalActions`/`RunMatch` root-action and force-edit
vocabulary (command_action.go, command_manual.go), and the `internal/session`
driver that owns `{version, seed, sets, []Command}` and replays it. The remaining
step is migrating the web client off its own hand-rolled `input`/replay log
(`internal/web/replay.go`) onto that driver.

## Context

Today a match persists its `engine.GameState` directly. The web client serializes
a `snapshot{Version, Seed, State, Log, ...}` to local storage, and undo is a
client-side stack of `GameState` copies. Three problems follow from persisting the
_outcome_ rather than the _cause_:

- **Undo is lost on reload.** The undo stack is in-memory only and is not part of
  the saved snapshot, so a reload restores the current state but not the history
  behind it.
- **Logs do not hydrate.** A typed `LogEntry` is an interface value that does not
  survive JSON (ADR 0011 keeps it outside the flat state). On save each entry is
  flattened to the prose it was narrated with and restored as a `RestoredEntry`,
  which carries no card IDs and cannot render card mentions or be inspected as its
  original type.
- **Two sources of truth.** The saved `GameState` and the inputs that produced it
  can disagree, and there is no way to replay a match — which a networked and
  shareable-replay future needs as its foundation.

The engine is deterministic given its initial seed, its sets, and the ordered
player inputs. Nothing exploits that.

## Decision

**The authoritative record of a match is its ordered command log plus the initial
seed and sets.** A command is one player input crossing the engine boundary — a
root action or one answer to a choice the engine asked. State, the typed log, and
the undo history are all _derived_ by replaying the commands from a fresh game.

- **The persisted artifact is `{version, seed, sets, []Command}`** (plus transient
  UI state). `GameState` is no longer serialized for saves at all, which removes
  the second source of truth.
- **Reload replays the command log from the start**, regenerating the exact state
  and the fully typed log with recording on. This deletes the lossy
  `savedLine`/`RestoredEntry` path. Replay is cheap — the simulator plays over a
  thousand whole games per second — so no periodic state snapshot is kept; the
  command log is the save.
- **Undo replays the log up to command N.** Undo is unrestricted up to the nearest
  _information barrier_ — a command whose resolution stepped the PRNG or revealed a
  hidden zone. Crossing a barrier requires the opponent's consent in networked
  play and is unrestricted in hotseat.
- **Randomness lives inside `GameState`** as one counter-based PRNG (flat and
  comparable, ADR 0005), replacing the `*rand.Rand` that hung off `Game`, so replay
  is bit-exact. Per-player streams are deferred; they buy only isolation and are
  additive later.
- **A client renders through a projection seam**, `Project(state, viewer) → View`,
  identity today. A future server swaps identity for redaction so a client sees
  only its own hidden zones, without the client reading raw `GameState`.
- **No forward compatibility.** A command log is tagged with an engine and catalog
  version and is refused on mismatch. Shared replays are same-version only.

## Consequences

- Hotseat, fully client-side, one player driving both sides, is the first target
  and stays permanently supported; server-authoritative redaction is designed for
  but not built here.
- MCTS is unaffected: it drives the engine over in-memory `FastCopy` copies with
  recording off and never touches the command log.
- Reload cost is replay cost, bounded by game length rather than by a saved-state
  size. If reload latency ever bites, a serializable _typed_ log checkpoint can be
  added then, not now.

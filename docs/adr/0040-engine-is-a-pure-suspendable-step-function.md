# 40. The engine is a pure suspendable step function

This decision records the target shape of the engine's resolution loop. It guides
a staged refactor. The engine seam that realizes it now exists — the suspendable
`Stepper` and `Command`/`Request` step function (suspend.go), driven by
`internal/session`; the web client's migration onto it is the remaining step. It
is the counterpart to ADR 0039, which makes commands the source of truth.

## Context

The engine resolves an action by _pulling_ choices synchronously: mid-resolution
it calls a `Chooser` on `Game` and blocks until the answer comes back. That model
has two limits the command-sourced, networked future cannot live with:

- **It blocks.** A remote or asynchronous player may answer seconds or hours after
  the engine asks. A blocking pull cannot represent a decision that has not
  arrived yet.
- **Choices are Go interface values.** A pulled choice is not a serializable input,
  so it cannot be recorded as a command, replayed, or sent over a wire.

## Decision

**The engine is a pure suspendable step function**,
`Advance(state, command) → (state, request, stepInfo)`. It never blocks. When
resolution needs input it _yields_ rather than pulls:

- **`Request` is a thin decision-point marker** — which player owes a decision, and
  in what context — not an enumerated candidate list. Because every client carries
  the whole engine, each holder re-derives the legal answers itself through a
  separate on-demand query (`LegalCommands(state)` / `IsLegal(state, command)`). The
  authority never trusts a transmitted option list; it always re-derives.
- **`Chooser` inverts.** It stops being an engine callback and becomes a
  `Request → Command` answerer on the driving side: the UI answers interactively,
  `FirstChooser` answers deterministically for tests, the MCTS bot answers for
  rollouts. The engine no longer knows who answers.
- **`stepInfo` reports whether the step crossed an information barrier** — stepped
  the PRNG or revealed a hidden zone — which the undo rule in ADR 0039 reads.
- **A new `internal/session` drives the loop.** It owns `{version, seed, sets,
[]Command}` and the undo cursor, wraps the pure engine, and exposes apply, undo,
  and view. `internal/match` still sets a match up (decks, houses); `internal/web`
  becomes a thin projection over the session; a future server wraps the same core.

## Consequences

- The refactor touches the resolution loop and every site that pulls a choice.
- MCTS stays cheap: `Advance` allocates nothing (the request is thin), enumeration
  is pay-for-what-you-use, and rollouts drive the engine directly over `FastCopy`
  copies, bypassing the session entirely.
- Asynchronous, networked, and replayable play all reduce to feeding the same pure
  step function a stream of commands from different sources.

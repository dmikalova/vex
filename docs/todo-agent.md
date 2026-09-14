# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Work that has **no consuming card in an implemented set yet** does not belong
  here — park it in [todo-future-set.md](todo-future-set.md), keyed to the set that
  first needs it.
- Cite the ADR or doc that decided an item where one exists.
- **When you need a decision from the human, ask in the reply itself, grill-me
  style** — a numbered list of `❓ **Q1** - **title**: <question>` with a `➡️`
  recommended answer under each — at the end of the turn, then stop. Do **not**
  reach for an interactive question tool: under autopilot it is auto-answered with
  "work autonomously" and the human never sees it. Plain-text questions at the end
  of the turn are the channel the human actually reads.

## Event-sourced replay (ADR 0039, 0040)

Decided in ADR 0039 (commands are the source of truth; state, log, and undo are
replay projections) and ADR 0040 (the engine is a pure suspendable step
function). Glossary terms in [CONTEXT.md](../CONTEXT.md): Command, Replay,
Session, Request, Information barrier, Projection/View. Build staged, each stage
independently green; hotseat-first, server redaction designed-for but not built.

- **Stage 2 — suspendable step function.** DONE for the engine core: `Advance`
  (via `Stepper`), thin `Request`, `Command`, `LegalCommands`/`IsLegal`, and
  `StepInfo` live in `internal/engine/suspend.go`, reusing the existing effects
  through a `suspendChooser` that yields a `Request` instead of pulling. The
  synchronous `Chooser` path is deliberately kept for MCTS rollouts and tests
  (allocation-free over `FastCopy`, per ADR 0040's consequences). The remaining
  inversion — pointing the live UI at the suspendable path — folds into Stage 5's
  web rewire. `StepInfo` currently flags PRNG steps; hidden-reveal detection and
  re-deriving the legal set purely from `GameState` (Request carries the candidate
  context today) are later refinements.
- **Stage 3 — `internal/session` driver.** DONE as an additive package:
  `internal/session` owns `Record{Version, Seed, Sets, []Command}` and the undo
  cursor, wraps the engine via the Stage-2 `Stepper`, and exposes
  `Apply`/`Undo`/`Pending`/`View`/`Record`/`Load`. Undo replays 0→N from a fresh
  deal (state is a projection of the log, not an inverse op); `CrossesBarrier(n)`
  reports whether an undo rewinds past an information barrier (free in hotseat,
  consent to cross in networked play). `Setup`/`Action` are injected so the session
  does not depend on deck generation — `internal/match` stays setup-only. The
  driving `Action` is supplied by the caller; expressing the web's whole-game turn
  loop as a suspendable `Action` (so root actions become `Command`s too) is part of
  Stage 5's rewire.
- **Stage 5 — rewire `internal/web` onto the session.** Persist only
  `{version, seed, sets, []Command, ui}` (stop serializing `GameState`); reload
  replays from 0, regenerating exact state + fully typed log; delete the lossy
  `savedLine`/`RestoredEntry` path; drop the client undo/redo stacks in favor of
  session undo. Version-tag the log and refuse on mismatch (no forward compat).
  NOTE: this is the riskiest stage — it rewrites the web action/chooser/persist
  path and cannot be visually verified from here. It also carries the residual
  Stage-2 UI inversion (point the live UI at the suspendable `Stepper`) and needs
  the whole-game turn loop expressed as a session `Action` with root-action
  `Command`s. Approach behind the green web test suite; land incrementally.

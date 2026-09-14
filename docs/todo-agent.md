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

- **Stage 1 — PRNG into `GameState`.** Replace the `*rand.Rand` on `Game` with one
  counter-based PRNG (PCG/philox; state is one or two `uint64`s) stored in
  `GameState`, so it is flat/comparable (ADR 0005) and captured by `FastCopy`.
  Route every shuffle/random-select site through it. Per-player streams deferred
  (isolation-only, additive later). This unblocks bit-exact replay.
- **Stage 2 — suspendable step function.** Turn resolution into
  `Advance(state, command) → (state, request, stepInfo)` that yields a thin
  `Request` instead of pulling a `Chooser`; add on-demand `LegalCommands`/`IsLegal`
  derived from local state; invert `Chooser` implementations (UI, `FirstChooser`,
  MCTS bot) into `Request → Command` answerers. `stepInfo` reports information-
  barrier crossings (PRNG step / hidden reveal). Keep `Advance` allocation-free for
  MCTS.
- **Stage 3 — `internal/session` driver.** New layer owning
  `{version, seed, sets, []Command}` + undo cursor; wraps the engine; exposes
  apply/undo/view. Undo replays 0→N, free up to the nearest information barrier,
  consent to cross in networked mode. `internal/match` keeps match setup only.
- **Stage 4 — projection seam.** `Project(state, viewer) → View`, identity now;
  clients render from the View, never raw state.
- **Stage 5 — rewire `internal/web` onto the session.** Persist only
  `{version, seed, sets, []Command, ui}` (stop serializing `GameState`); reload
  replays from 0, regenerating exact state + fully typed log; delete the lossy
  `savedLine`/`RestoredEntry` path; drop the client undo/redo stacks in favor of
  session undo. Version-tag the log and refuse on mismatch (no forward compat).

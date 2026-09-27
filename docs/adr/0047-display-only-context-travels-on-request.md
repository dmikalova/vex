# 47. The client holds no chooser; display-only context travels on `Request`

This decision records what happens when the UI needs something from the engine
that is not an answer. It is the standing rule left behind by the ADR 0040
switchover, and the counterpart to ADR 0039 (commands are the source of truth)
and ADR 0045 (the UI is proven to render every prompt kind).

## Context

Before the switchover the web client implemented `engine.Chooser` and its
optional capabilities. The engine pulled a decision by calling into the client,
so anything else the UI wanted could be pushed the same way: the engine called a
method, the client drew something, and the call returned. `BadgeChooser`
(`internal/engine/game.go`) is exactly that shape — `PreviewBadge(SelectionBadge)`
pushes the "3 damage" a Festering Touch pick will deal so the board can badge each
candidate.

The switchover removed that channel. The client holds no `Chooser`: it drives a
`*session.Session`, reads the pending `Request`, and answers with a `Command`
(ADR 0040). The only `Chooser` left on the interactive path is `suspendChooser`
(`internal/engine/suspend.go`), which lives inside the engine and answers nothing
— it yields. So there is no longer anyone for the engine to push to, and the
question the badge raised will keep recurring: every future "the UI needs to know
X about this decision" has to arrive somehow.

Two answers are available and only one is safe. A callback field on `Session`, or
a small interface the client implements and the session calls, would work and
would look local and cheap. It is the live chooser rebuilt under another name: it
puts client code back on the engine's resolution path, makes the thing it carries
unrecordable and unreplayable, and reintroduces exactly what the ADR 0040 line
exists to remove.

## Decision

**The client holds no chooser. Anything the UI needs that is not an answer travels
as display-only context on `Request`.**

- A new need is a **new field on `engine.Request`**, populated at the suspension
  point where the engine already knows the value, and read by the client off
  `session.Pending()`.
- Display-only context is **context for a decision, not part of it**. It has no
  mirror in `Command`, so it can never enter the command log or change a replay,
  and `Request.LegalCommands` and `Request.IsLegal` ignore it, so it can never
  change what answers are legal.
- **Forbidden**: a callback field on `Session`, or a live interface the client
  implements for the engine or the session to call. Either one quietly
  reintroduces the live chooser.

Two instances, one built and one queued:

- **The badge preview (built).** `Request.Badge` carries the `SelectionBadge` an
  effect previewed before its choose loop. `suspendChooser` stashes it and stamps
  it on every `Request` yielded until the effect clears it
  (`TestStepperStampsThePreviewedBadgeOnRequests`); the client mirrors it in
  `syncBadge` (`internal/web/game_badge.go`). The pushed `BadgeChooser` capability
  survives for the pulled path (tests, and any chooser that draws immediately);
  the suspendable path routes the same value through `Request`.
- **The prompt's source `LocalID` (queued).** `Request.Source` is the source
  card's _name_, so `internal/web/view_controls.go` matches on it and two copies
  of a card in play can preview the wrong one's live house or Maverick state. The
  fix is a source `LocalID` on `Request`, not a lookup callback — which is now a
  small change, because the switchover deleted `replayChooser` and `webChooser`
  and left one interactive implementation to thread it through.

## Consequences

- The engine keeps one direction of flow on the interactive path: commands in,
  requests out. Nothing the client owns runs inside resolution.
- `Request` grows a field per display need. That is the intended cost: a field is
  flat, comparable, inspectable in a test, and visibly outside the command log,
  where a callback is none of those.
- A display-only field is invisible to replay by construction, so adding one never
  requires a `session.Version` bump — unlike a change to what a command means.
- The rule applies to the interactive path only. `FirstChooser`, `bridgeChooser`
  and the MCTS bot still pull, and a pushed optional capability stays the right
  shape for them.

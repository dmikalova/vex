# 28. Take control is a LIFO stack of control effects, not one source per card

## Context

Several cards take control of another player's card. Harland Mindlock and Collar
of Subordination seize an enemy creature "until this leaves play"; Sneklifter and
Spangler Box hand an artifact over for good; Anahita the Trader gives one away.
The rule is uniform across all of them: control is granted by an ability, that
grant lasts for some duration (end of turn, "until the granting card leaves play",
or forever), and while it is in effect the granting ability's controller controls
the card. Ownership never changes and still decides which pile the card returns to
when it leaves play.

The first model stored control as a single override on `CardCore`: a `ControlPlus`
byte (0 = owner controls, otherwise controller+1) plus a `ControlSource LocalID`
naming the one card whose leaving play reverts it. That single source cannot
express two grants at once. When a second ability takes control of a card already
seized, it overwrote `ControlSource`; when its source then left play, the card
reverted straight to its **owner**, skipping the first grant that was still in
effect. Artifacts were worse — `takeControlOfArtifact` set no source at all, so a
seized artifact could never be handed back, even by a duration that should end.

KeyForge resolves overlapping control the way every "most recent effect wins"
situation resolves: the newest grant applies, and when it ends the next-newest
grant still in effect takes over. That is a **LIFO stack** of control effects per
card, not a single slot.

## Decision

Control effects live in one **global stack** on `GameState`, mirroring the generic
counter side-table (ADR 0024):

```go
type ControlEntry struct {
    Card       LocalID // the seized card
    Controller uint8   // player index (0 or 1) this grant hands the card to
    Source     LocalID // the card whose leaving play ends this grant
}

Controls     [maxControlEntries]ControlEntry
ControlCount uint8
```

The entries are ordered by when they were applied. A card's current controller is
the `Controller` of its **most recently pushed** entry; entries beneath it are the
fallback a later one reverts to when it is removed. `takeControl` pushes an entry
and relists the card under the new controller. `releaseControlHeldBy(source)`
drops every entry that `source` holds and re-derives each affected card's
controller from the entry beneath — the LIFO fallback — reverting to the owner
only when no entry remains.

`ControlPlus` stays on `CardCore` as a fast, O(1) cache of the current controller
(read on every "friendly"/"enemy" check), refreshed whenever a grant is pushed or
popped. The stack is consulted only when control changes.

A **Forever** grant (an artifact taken for good) names the seized card _itself_ as
its `Source`, so no leaving source reverts it; it is shed only when the card leaves
play. That unifies the permanent and reverting cases: every exit funnels through
`removeFromPlay`, which calls `releaseControlHeldBy(id)` to revert what the leaving
card controlled and `clearControls(id)` to drop the leaving card's own entries.
`putIntoPlay` under a non-owner pushes a self-sourced entry the same way.

This subsumes the two old methods (`takeControl` for creatures,
`takeControlOfArtifact` for artifacts) into one `takeControl(id, controller,
source)` that handles both card types, and it fixes the latent bug where a
creature source (Harland Mindlock) leaving play did not revert what it held —
`removeFromPlay` now reverts control for every exit, not only for a leaving
upgrade.

## The stack tracks the board, not the game length

A naive append-only stack would grow once per control change over a whole game,
not per board state: a card taken permanently again and again would push a
never-popped entry each time and eventually overflow any fixed bound. A
control-heavy game is longer than a small cap, so the bound has to track the board,
not the number of takes.

A fresh take therefore **supersedes** its own source's earlier entry on the card
(`supersedeControl`) before pushing a new one: it drops any entry with the same
`(card, source)` pair. This keeps each source to one entry per card it holds, so
the live count is at most the distinct sources controlling a card. Every permanent
take names the seized card itself as its source, so repeated permanent takes share
a `(card, source)` key and collapse to one entry the same way — no special case for
"permanent" is needed. Sources that leave play shed their entries
(`releaseControlHeldBy`), so the live count is bounded by the board — the same
sparsity argument as ADR 0024 — and the fixed cap is a defensive ceiling far above
it, not a game-length budget.

An entry buried under a later one from a different source is harmless: the only way
to pop the entry above it is that source (or the card) leaving play, and a
permanent take's source is the card itself, so anything beneath a permanent entry
is cleared with the card when it leaves. Superseding by `(card, source)` alone is
therefore enough; there is no need to also collapse unrelated entries.

## Consequences

- Overlapping control resolves correctly: the newest grant wins, and removing it
  falls back to the next grant still in effect (`TestControlStackIsLIFO`), reverting
  to the owner only when the stack empties.
- Repeated control of the same card does not accumulate: taking it permanently over
  and over holds at one entry (`TestPermanentControlDoesNotAccumulate`), and a
  source re-taking a card it already holds dedupes
  (`TestRevertibleControlDedupesBySource`).
- Artifacts are revertible on the same footing as creatures; a permanent grant is
  just a self-sourced entry, not a special case.
- The stack is flat, pointerless, and comparable, so it survives the value-copy
  snapshot that makes undo and MCTS cheap (ADR 0005), and it is sparse, so a card
  never seized costs nothing per snapshot (the same reasoning as ADR 0024).
- `CardCore.ControlSource` is gone; control history lives in the stack. The
  persisted-state layout changed, so `snapshotVersion` is bumped.
- A live control entry past capacity panics rather than silently dropping a grant,
  matching the counter table's full-table invariant.

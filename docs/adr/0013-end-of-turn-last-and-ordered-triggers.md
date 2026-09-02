# End-of-turn abilities resolve last, and simultaneous triggers are ordered

## Context

Two resolution-order bugs were found while designing the phase machine (ADR 0012).
Both had test expectations pinning the wrong behavior, so recording why the old
expectations are wrong is a precondition for changing them.

**End-of-turn abilities fired first, not last.** `EndTurn` resolved every
`TriggerEndOfTurn` ability, _then_ readied cards, _then_ drew. An ability that
reads "at the end of your turn" therefore saw a board where nothing had readied
and the hand had not refilled — the state at the end of the play phase, not the
end of the turn.

**Simultaneous triggers fired in board order.** When several abilities triggered
on one event, `triggerAbilities` resolved them in scan order: printed abilities,
then upgrade-granted, then constant-granted, player 0 before player 1. KeyForge
gives the active player the choice of order, which matters whenever one trigger
changes whether another can resolve. Destruction was the sole exception —
`destroyTogether` hand-rolled the player's choice, so the correct behavior existed
but only in one place.

## Decision

**End of turn is the last phase**, resolving after ready and draw (ADR 0012).

**Every multi-source trigger window is ordered by the active player.**
`triggerAbilities` routes through the existing `orderByChoice` helper whenever more
than one ability fires on an event, and `destroyTogether`'s hand-rolled ordering is
replaced by the same call. A single trigger is forced and never prompts; the
default chooser keeps scan order, so ordering only becomes interactive under a real
frontend.

## Consequences

- Existing tests that asserted the old order are wrong and are rewritten to assert
  the new one — not deleted, and not weakened.
- Cards whose end-of-turn ability depended on seeing unreadied creatures change
  behavior. This is the correction, not a regression.
- The ordering prompt is a new interaction point: a frontend that ignores it falls
  back to scan order, which is the previous behavior.
- The game log can state the order a player chose, because there now is one
  (ADR 0011).

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
- **Identical abilities are auto-ordered, never prompted.** Ordering compares
  ability _identity_ — trigger plus rendered text — not card, because the same
  card can resolve differently as the board changes. Abilities that would resolve
  identically are collapsed to one representative before the prompt, so a window
  of only identical abilities is silent and a window with a distinct one asks only
  about the distinct abilities. See `orderTriggered` in `game_abilities.go`.
- **A frontend may offer Auto-resolve on an ordering prompt.** Whether an order is
  worth arranging is the player's call, not the engine's, so the client's
  `Orderer` implementation adds an Auto-resolve button that answers with a random
  order in one click. The engine is unchanged: the default chooser still keeps
  scan order.
- **A printed `Reap:`/`Fight:`/`Action:`/`Play:` ability orders in the same window
  as the bystander reactions to that action.** Each printed timing ability is
  semantically a reaction to its own event ("after this creature reaps, do …"), so
  a use/play verb gathers the acting card's own printed ability together with every
  bystander "after a creature reaps/is used/fights/is played" reaction into one
  `orderTriggered` pass (`reapReactions`, `fightReactions`, `actionReactions`,
  `playCreatureReactions`, `afterPlayReactions` in `game_abilities.go`), rather than
  firing the printed ability first as a separate step. Every entry is ordered by the
  active player but resolves for its own controller, so a bystander's "after an enemy
  creature reaps: gain Æmber" still credits its owner. Triggering a named ability
  directly (Replicator makes a creature reap) is not performing the action, so it
  resolves the ability alone and opens no window (`TriggerAbilityOf`).
- **Duration reactions order in the same window as the card-sourced set, through
  the flat `ReactionChooser` port.** A "for the remainder of the turn" reaction (Full
  Moon, Charge!, Crystal Hive) lives in the flat lasting registry (ADR 0007) and has
  no in-play source card. Each unified window folds its duration reactions into the
  same `orderTriggered` pass as window entries (`lastingReactions` builds them;
  `resolveWindow` resolves them through `resolveReaction`) rather than resolving them
  in a trailing `emitLasting(…)` window of its own. The whole window — card abilities
  and duration reactions alike — is one flat labeled list: `orderTriggered` renders
  each entry as an `OrderableReaction` (a card ability shows its source card and
  rendered text, a duration reaction its rendered effect) and asks the active player,
  through the optional `ReactionChooser` capability, to pick the next reaction to
  resolve, repeating until one remains. A window whose entries are _all_ identical is
  auto-ordered and never prompts, since their order cannot matter; any distinct entry
  makes the whole window ordered in full (identical entries included), because
  resolving one entry can change what another would do. A `Chooser` without
  `ReactionChooser` (the AI, the simulator) keeps the gathered order — the card
  abilities in collection order, then the duration reactions in registry order — so
  folding never reorders the card abilities that already resolved there. The cardtest
  harness (`bridgeChooser.ChooseReaction`, scripted by `Player.Order(cards…)`) and the
  web client (`webChooser.ChooseReaction`, rendering the reactions as a labeled list)
  both implement the port. `emitLasting` remains for the events that have no card
  window of their own (an enemy creature destroyed).

# Turn cleanup runs after end-of-turn abilities, in the end-of-turn phase

## Context

`readyPhase` (`internal/engine/game_phase.go`) did two unrelated jobs. It readied
the active player's cards, and it also lifted everything that expires with the
turn: `revertTemporaryCreatures`, `GrantedKeywords`, `LostKeywords`,
`ConsideredFlank`, `TempPowerBonus`, `TempArmorBonus`, `TempAssaultBonus`,
`TextBoxTurnSourcePlus`, `TempHouse`, the armor refresh, the six `CannotX` bars,
the `TurnHistory` roll into last-turn, the `MayFight*` / `MayUse*` / `MayPlay*`
permits, `KeyCostBump`, `KeyCostPerHouse`, and `clearLasting`.

Ready is phase 6; end of turn is phase 8. So a "remainder of the turn" effect
expired two phases before the end-of-turn abilities that are supposed to see it.

The KeyForge Master Rulebook says the opposite. Answering an Animator plus
Fangtooth Cavern question
([keyforge-master-rulebook.md](../keyforge-master-rulebook.md)): "Animator's
ability creates a lasting effect that expires at the end of the turn, which is
after Fangtooth Cavern's 'end of turn' effect resolves." Vex reverted Animator
inside `readyPhase`, before that window. The same ruling covers a temporary
house: "At the end of the turn, the upgraded neighbor will belong only to House
Logos" — Vex cleared `TempHouse` in `readyPhase`.

This was the unfinished half of
[ADR 0013](0013-end-of-turn-last-and-ordered-triggers.md). That ADR moved the
end-of-turn ability window to run last, after ready and draw, because an
end-of-turn ability was otherwise seeing "the state at the end of the play phase,
not the end of the turn". It moved the window but left the turn's cleanup welded
into the ready phase, so the abilities resolved last while the effects they must
observe had already expired.

## Decision

**The ready phase only readies.** It unexhausts the active player's cards and
records `CardsReadied`; `TriggerEndOfReadyStep` still fires there.

**Every turn-scoped expiry moves into the end-of-turn phase, after the ability
window.** `endOfTurnPhase` already had a cleanup tail (`clearScheduled`,
`clearExpiredContinuous`) running after `resolveWindow`; the rest of the cleanup
joins it as `expireTurnScoped`. The phase's order is:

1. Resolve end-of-turn abilities and the effects scheduled into this window. The
   active player orders them (ADR 0013).
2. Settle destruction.
3. Clean up, in whatever order the engine finds convenient.

**Step 3 is not player-orderable.** Only step 1 is a window of card abilities,
and only a window is ordered by the active player. Cleanup is bookkeeping with no
card-visible interleaving, so it runs in a fixed engine order.

**No ninth phase.** The turn keeps its eight phases
([ADR 0012](0012-first-class-turn-phases.md)). Cleanup is the tail of the
end-of-turn phase, not a phase of its own, so the rulebook's turn structure, the
log's `PhaseBegan` entries, and the turn HUD are unchanged.

**Armor refresh moves with the rest.** The `ArmorRemaining` / `ArmorStripped`
reset is an expiry ("armor is reduced for the remainder of the turn"), not part
of readying, so it belongs in the cleanup tail.

## Consequences

- Tests that pinned the old ordering are rewritten, not deleted: the ruling cited
  above is the rule that makes the old expectation wrong.
  `TestAnimatorRevertsAfterEndOfTurnAbilities` pins the worked example.
- Cards with a "remainder of the turn" effect and cards with an end-of-turn
  ability now interact as the rulebook says. Animator plus an end-of-turn ability
  is the worked example.
- A card animated for the turn now lives as a creature through ready, draw, and
  the end-of-turn window, so it faces those phases' destruction settles. An
  animated card with no power dies rather than quietly reverting first; Animator
  gives its target three +1 power counters, so the printed card is unaffected.
- No divergence-register entry: the change moves toward the printed rulebook
  rather than away from it.
- A ready phase that only unexhausts can run through
  `CreatureResolver.SetExhausted`, which puts it inside the narration audit at the
  Resolver port for free. The cleanup tail stays off the port: its writes are not
  capabilities any card has, and `resolver.go` describes the port as the full
  explicit catalogue of what an effect may do.

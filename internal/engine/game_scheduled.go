package engine

// maxScheduled bounds how many effects can be scheduled into one turn's end-of-turn
// window at once — generous for the handful of cards that schedule a delayed effect,
// and a fixed size keeps the state flat.
const maxScheduled = 4

// scheduledAction is the flat, enum-tagged identity of an effect scheduled to
// resolve in the active player's end-of-turn window. State holds no effect closures
// (ADR 0005), so a scheduled effect is stored as this tag and rebuilt into its
// Effect by scheduledEffectOf when the window fires.
type scheduledAction uint8

const (
	schedUnset scheduledAction = iota
	// schedDestroyEachCreature is Ragnarok's board wipe: destroy each creature.
	schedDestroyEachCreature
)

// ScheduledEffect is one effect armed during the turn to resolve in the active
// player's end-of-turn window, alongside the in-play "at the end of your turn"
// abilities (ADR 0013). Source is the card that armed it, for attribution.
type ScheduledEffect struct {
	Source LocalID
	Do     scheduledAction
}

// scheduledEffectOf rebuilds a scheduled action into the Effect the end-of-turn
// window resolves.
func scheduledEffectOf(a scheduledAction) Effect {
	if a == schedDestroyEachCreature {
		return Destroy{Target: Target{Kind: TargetEachCreature}}
	}
	return nil
}

// ScheduleAtEndOfTurn arms an effect to resolve in the active player's end-of-turn
// window, dropping it silently when the schedule is full. The schedule survives the
// ready phase (which runs before end of turn) and is cleared as the window fires.
func (g *Game) ScheduleAtEndOfTurn(source LocalID, do scheduledAction) {
	if int(g.State.ScheduledCount) >= maxScheduled {
		return
	}
	g.State.Scheduled[g.State.ScheduledCount] = ScheduledEffect{Source: source, Do: do}
	g.State.ScheduledCount++
}

// scheduledEndOfTurn returns the effects scheduled into this turn's end-of-turn
// window as window entries, so they order alongside the in-play end-of-turn
// abilities. The active player is their actor.
func (g *Game) scheduledEndOfTurn(player int) []triggeredAbility {
	var pending []triggeredAbility
	for i := 0; i < int(g.State.ScheduledCount); i++ {
		s := g.State.Scheduled[i]
		pending = append(pending, triggeredAbility{
			source:  s.Source,
			grantor: s.Source,
			ability: Ability{Trigger: TriggerEndOfTurn, Effect: scheduledEffectOf(s.Do)},
			actor:   int8(player),
		})
	}
	return pending
}

// clearScheduled empties the end-of-turn schedule once its window has fired.
func (g *Game) clearScheduled() {
	for i := 0; i < int(g.State.ScheduledCount); i++ {
		g.State.Scheduled[i] = ScheduledEffect{}
	}
	g.State.ScheduledCount = 0
}

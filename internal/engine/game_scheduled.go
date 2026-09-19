package engine

// maxScheduled bounds how many effects can be scheduled at once — the handful
// resolving in a turn's end-of-turn window plus the few armed to resolve when a
// card leaves play. A fixed size keeps the state flat.
const maxScheduled = 4

// scheduledAction is the flat, enum-tagged identity of an effect scheduled to
// resolve later — in the active player's end-of-turn window, or when its source
// leaves play. State holds no effect closures (ADR 0005), so a scheduled effect is
// stored as this tag and rebuilt into its Effect by scheduledEffectOf when its
// window fires.
type scheduledAction uint8

const (
	schedUnset scheduledAction = iota
	// schedDestroyEachCreature is Ragnarok's board wipe: destroy each creature.
	schedDestroyEachCreature
	// schedOpponentForgesKeyFree is Turnkey's delayed consequence: the source
	// card's controller's opponent forges a key at no cost.
	schedOpponentForgesKeyFree
)

// ScheduledEffect is one effect armed to resolve later, either in the active
// player's end-of-turn window alongside the in-play "at the end of your turn"
// abilities (ADR 0013, Duration RemainderOfPlayerTurn) or when its source card
// leaves play, however many turns on (Duration UntilThisLeavesPlay). Source is the
// card that armed it, both for attribution and — for a leave-play schedule — to
// match the card whose exit fires it.
type ScheduledEffect struct {
	Source   LocalID
	Do       scheduledAction
	Duration Duration
}

// scheduledEffectOf rebuilds a scheduled action into the Effect its window
// resolves.
func scheduledEffectOf(a scheduledAction) Effect {
	switch a {
	case schedDestroyEachCreature:
		return Destroy{Target: Target{Kind: TargetEachCreature}}
	case schedOpponentForgesKeyFree:
		return ForgeKey{
			Player:     Opponent,
			FreeOfCost: true,
		}
	}
	return nil
}

// scheduledActionOf maps an Effect an author armed to the flat action the schedule
// stores for it, reporting whether the schedule can carry it. It is the inverse of
// scheduledEffectOf, used by the arming effects (ScheduleOnLeave) to reduce a
// carried Effect to its enum tag.
func scheduledActionOf(e Effect) (scheduledAction, bool) {
	if f, ok := e.(ForgeKey); ok && f.Player == Opponent && f.FreeOfCost {
		return schedOpponentForgesKeyFree, true
	}
	return schedUnset, false
}

// ScheduleAtEndOfTurn arms an effect to resolve in the active player's end-of-turn
// window, dropping it silently when the schedule is full. The schedule survives the
// ready phase (which runs before end of turn) and is cleared as the window fires.
func (g *Game) ScheduleAtEndOfTurn(source LocalID, do scheduledAction) {
	g.schedule(source, do, RemainderOfPlayerTurn)
}

// ScheduleOnLeave arms an effect to resolve when the source card leaves play, at
// whatever later point that is (Turnkey's forced forge). Unlike an end-of-turn
// schedule it survives across turns — it is swept only when its source leaves play
// (fireScheduledOnLeave), never by the end-of-turn window.
func (g *Game) ScheduleOnLeave(source LocalID, do scheduledAction) {
	g.schedule(source, do, UntilThisLeavesPlay)
}

// schedule appends one scheduled effect, dropping it silently when the schedule is
// full.
func (g *Game) schedule(source LocalID, do scheduledAction, dur Duration) {
	if int(g.State.ScheduledCount) >= maxScheduled {
		return
	}
	g.State.Scheduled[g.State.ScheduledCount] = ScheduledEffect{
		Source:   source,
		Do:       do,
		Duration: dur,
	}
	g.State.ScheduledCount++
}

// scheduledEndOfTurn returns the effects scheduled into this turn's end-of-turn
// window as window entries, so they order alongside the in-play end-of-turn
// abilities. Leave-play schedules are not end-of-turn entries, so they are skipped.
// The active player is the entries' actor.
func (g *Game) scheduledEndOfTurn(player int) []triggeredAbility {
	var pending []triggeredAbility
	for i := 0; i < int(g.State.ScheduledCount); i++ {
		s := g.State.Scheduled[i]
		if s.Duration != RemainderOfPlayerTurn {
			continue
		}
		pending = append(pending, triggeredAbility{
			source:  s.Source,
			grantor: s.Source,
			ability: Ability{
				Trigger: TriggerEndOfTurn,
				Effect:  scheduledEffectOf(s.Do),
			},
			actor: int8(player),
		})
	}
	return pending
}

// clearScheduled drops the end-of-turn schedule once its window has fired, keeping
// the leave-play schedules, which live until their source leaves play.
func (g *Game) clearScheduled() {
	n := 0
	for i := 0; i < int(g.State.ScheduledCount); i++ {
		if g.State.Scheduled[i].Duration == RemainderOfPlayerTurn {
			continue
		}
		g.State.Scheduled[n] = g.State.Scheduled[i]
		n++
	}
	for i := n; i < int(g.State.ScheduledCount); i++ {
		g.State.Scheduled[i] = ScheduledEffect{}
	}
	g.State.ScheduledCount = uint8(n)
}

// fireScheduledOnLeave resolves and removes every leave-play schedule armed by the
// card now leaving play. It runs inside emitLeavesPlay while the card is still on
// the board, so its controller is still readable — Turnkey's opponent is the player
// whose key it unforged. An entry is removed before it resolves, so a consequence
// that itself moves cards cannot fire it twice.
func (g *Game) fireScheduledOnLeave(id LocalID) {
	for i := 0; i < int(g.State.ScheduledCount); {
		s := g.State.Scheduled[i]
		if s.Duration != UntilThisLeavesPlay || s.Source != id {
			i++
			continue
		}
		g.removeScheduledAt(i)
		if eff := scheduledEffectOf(s.Do); eff != nil {
			eff.Resolve(&EffectContext{
				Resolver:   g,
				Source:     id,
				Controller: g.controller(id),
			})
		}
	}
}

// removeScheduledAt drops the schedule entry at index i, sliding the rest down and
// zeroing the freed tail slot so the flat state stays canonical.
func (g *Game) removeScheduledAt(i int) {
	for j := i; j < int(g.State.ScheduledCount)-1; j++ {
		g.State.Scheduled[j] = g.State.Scheduled[j+1]
	}
	g.State.ScheduledCount--
	g.State.Scheduled[g.State.ScheduledCount] = ScheduledEffect{}
}

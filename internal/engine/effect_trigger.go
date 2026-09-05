package engine

import "fmt"

// Triggering another card's ability is not using that card. The card does not
// exhaust, its use is not recorded, and nothing that watches for a card being
// used fires — only the abilities printed under the named trigger resolve, and
// they resolve for the player whose effect reached for them, as if that player
// controlled the card.
// TriggerAbility resolves the abilities a chosen card carries for one trigger as
// if the effect's controller controlled that card. Replicator reaps by
// triggering the reap effect of another creature in play.
type TriggerAbility struct {
	// Trigger names which of the target's abilities resolve — its Play, Fight, or
	// Reap effect.
	Trigger Trigger
	// Target names the card whose abilities are triggered.
	Target Target
}

// validate requires a target and one of the three action triggers a card can
// name as "the <verb> effect of" another card.
func (e TriggerAbility) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("TriggerAbility")
	}
	if triggerEffectNoun(e.Trigger) == "" {
		return fmt.Errorf("TriggerAbility: Trigger must be Play, Fight, or Reap")
	}
	return nil
}

// triggerEffectNoun names an action trigger the way a card refers to the ability
// hanging off it — "the reap effect of another creature". A trigger with no such
// name renders empty, which is what validate rejects.
func triggerEffectNoun(t Trigger) string {
	switch t {
	case TriggerAfterPlay:
		return "play"
	case TriggerAfterFight:
		return "fight"
	case TriggerAfterReap:
		return "reap"
	default:
		return ""
	}
}

// Text renders the effect, e.g. "trigger the reap effect of another creature".
// "As if you controlled that creature" is reminder text for what the node already
// does, so rule 1 strips it.
func (e TriggerAbility) Text() string {
	return "trigger the " + triggerEffectNoun(e.Trigger) + " effect of " + e.Target.Text()
}

// Resolve triggers the named abilities of each selected card for the effect's
// controller. Only cards that actually carry the trigger are offered, so the
// choice is never a wasted one. A chain of these — two Replicators reaching for
// each other — stops at the Rule of Six.
func (e TriggerAbility) Resolve(ctx *EffectContext) {
	if ctx.Resolver.TriggerDepth() >= RuleOfSix {
		return
	}
	carries := func(id LocalID) bool { return ctx.Resolver.HasTrigger(id, e.Trigger) }
	for _, id := range e.Target.selectWith(ctx, false, carries) {
		ctx.It, ctx.HasIt = id, true
		ctx.Resolver.TriggerAbilityOf(ctx.Controller, id, e.Trigger)
	}
}

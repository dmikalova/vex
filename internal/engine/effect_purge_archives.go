package engine

import "fmt"

// PurgeArchivesForDamage lets the controller purge any number of cards from their
// own archives, then deals Amount damage to Target for each card purged this way —
// Destructive Analysis's "you may purge any number of cards from your archives to
// deal an additional 2 damage to the same creature for each card purged this way".
// The controller purges one card at a time and may stop at any point, so purging
// none is allowed and deals no additional damage.
type PurgeArchivesForDamage struct {
	// Amount is the damage dealt to Target for each card purged.
	Amount int
	// Target is the creature the additional damage lands on.
	Target Target
}

// validate requires a target and a positive per-card amount.
func (e PurgeArchivesForDamage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PurgeArchivesForDamage")
	}
	if e.Amount < 1 {
		return fmt.Errorf("PurgeArchivesForDamage: amount must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "purge any number of cards from your archives to
// deal an additional 2 damage to it for each card purged this way".
func (e PurgeArchivesForDamage) Text() string {
	return fmt.Sprintf(
		"purge any number of cards from your archives to deal an additional %d "+
			"damage to %s for each card purged this way",
		e.Amount, e.Target.Text())
}

// Resolve purges cards from the controller's archives one at a time until they
// decline, then deals Amount damage to each targeted creature for every card
// purged.
func (e PurgeArchivesForDamage) Resolve(ctx *EffectContext) {
	purged := 0
	for {
		cands := ctx.Resolver.Archives(ctx.Controller)
		if len(cands) == 0 {
			break
		}
		chosen, ok := ctx.ChooseCardOptional(
			"Choose a card to purge from your archives", cands)
		if !ok {
			break
		}
		ctx.Resolver.PurgeFromArchives(ctx.Controller, chosen)
		purged++
	}
	if purged == 0 {
		return
	}
	var hits []DamageTarget
	for _, id := range e.Target.Select(ctx) {
		hits = append(hits, DamageTarget{ID: id, Amount: e.Amount * purged})
	}
	if len(hits) > 0 {
		ctx.Resolver.DealDamage(ctx.Controller, hits)
	}
}

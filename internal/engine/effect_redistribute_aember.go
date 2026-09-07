package engine

import "fmt"

// RedistributeCapturedAember gathers all the Æmber sitting on the creatures of one
// side and lets the controller place it back among that side's creatures however
// they choose — Equalize redistributes the Æmber on friendly creatures among
// friendly creatures, then the Æmber on enemy creatures among enemy creatures. The
// total is conserved: every unit taken off is placed back on some creature of the
// same side.
type RedistributeCapturedAember struct {
	// Side names whose creatures the redistribution ranges over: Controller for the
	// controller's own creatures, Opponent for the enemy's.
	Side Player
}

// validate accepts the effect; both Side values are meaningful.
func (RedistributeCapturedAember) validate() error { return nil }

// Text renders the effect, e.g. "redistribute the Æmber on friendly creatures
// among friendly creatures".
func (e RedistributeCapturedAember) Text() string {
	who := "friendly"
	if e.Side == Opponent {
		who = "enemy"
	}
	return fmt.Sprintf("redistribute the Æmber on %s creatures among %s creatures", who, who)
}

// Resolve removes all the Æmber from the side's creatures into a pool, then places
// it back one unit at a time onto a creature of that side the controller chooses.
// With no creatures or no Æmber to move, it does nothing.
func (e RedistributeCapturedAember) Resolve(ctx *EffectContext) {
	side := ctx.PlayerFor(e.Side)
	creatures := ctx.Resolver.Battleline(side)
	pool := 0
	for _, id := range creatures {
		if held := ctx.Resolver.AmberOn(id); held > 0 {
			ctx.Resolver.AddAmberOn(id, -held)
			pool += held
		}
	}
	for ; pool > 0; pool-- {
		id, ok := ctx.ChooseCreature("Place 1 Æmber", creatures)
		if !ok {
			id = creatures[0]
		}
		ctx.Resolver.AddAmberOn(id, 1)
	}
}

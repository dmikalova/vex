package engine

// MoveAemberToPool moves 1 Æmber from one of the controller's cards (Æmber sitting
// on a creature or artifact from capture or exalt) into their pool — Selwyn the
// Fence. It does nothing when none of their cards carry Æmber.
type MoveAemberToPool struct{}

// Text renders the effect.
func (MoveAemberToPool) Text() string {
	return "move 1 \u00c6mber from one of your cards to your pool"
}

// Resolve has the controller pick one of their cards holding Æmber and move a
// single Æmber from it into their pool.
func (MoveAemberToPool) Resolve(ctx *EffectContext) {
	var candidates []LocalID
	for _, id := range append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Artifacts(ctx.Controller)...) {
		if ctx.Resolver.AmberOn(id) > 0 {
			candidates = append(candidates, id)
		}
	}
	if len(candidates) == 0 {
		return
	}
	id, ok := ctx.ChooseCreature("Choose one of your cards", candidates)
	if !ok {
		return
	}
	ctx.Resolver.AddAmberOn(id, -1)
	ctx.Resolver.SetAember(ctx.Controller, ctx.Resolver.Aember(ctx.Controller)+1)
	ctx.Resolver.Logf("%s moves 1 Æmber from %s to their pool", ctx.Resolver.PlayerName(ctx.Controller), ctx.Resolver.Name(id))
}

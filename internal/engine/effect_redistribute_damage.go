package engine

// RedistributeDamage lets the controller choose a player and then optionally move
// all the damage sitting on that player's creatures back among that player's
// creatures however they choose (Entropic Manipulator). The total is conserved:
// every unit taken off is placed back on some creature of the same player. A unit
// may be piled onto a creature past its power, which destroys it.
type RedistributeDamage struct{}

// validate accepts the effect; it has no configuration.
func (RedistributeDamage) validate() error { return nil }

// Text renders the effect's full printed sentence.
func (RedistributeDamage) Text() string {
	return "Redistribute the damage among a player's creatures"
}

// Resolve asks the controller which player to redistribute, then — if they
// accept — gathers all the damage off that player's creatures and places it back
// one unit at a time onto a creature of that player the controller chooses.
// Creatures whose damage reaches their power are destroyed afterward.
func (RedistributeDamage) Resolve(ctx *EffectContext) {
	r := ctx.Resolver
	choice := ctx.ChooseOption(
		"Choose a player",
		[]string{r.PlayerName(ctx.Controller), r.PlayerName(ctx.Opponent())},
	)
	player := ctx.Controller
	if choice == 1 {
		player = ctx.Opponent()
	}
	creatures := r.Battleline(player)
	pool := 0
	for _, id := range creatures {
		pool += r.Damage(id)
	}
	if pool == 0 || len(creatures) == 0 {
		return
	}
	if ctx.ChooseOption("Redistribute damage?", []string{"Yes", "No"}) != 0 {
		return
	}
	assigned := make(map[LocalID]int, len(creatures))
	for _, id := range creatures {
		r.SetDamage(id, 0)
	}
	for ; pool > 0; pool-- {
		id, ok := ctx.ChooseCreature("Place 1 damage", creatures)
		if !ok {
			id = creatures[0]
		}
		assigned[id]++
		r.SetDamage(id, assigned[id])
	}
	var dying []LocalID
	for _, id := range creatures {
		if p := r.Power(id); p > 0 && r.Damage(id) >= p {
			dying = append(dying, id)
		}
	}
	r.DestroyEach(ctx.Controller, dying)
}

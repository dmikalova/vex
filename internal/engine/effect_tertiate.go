package engine

// Tertiate destroys one third of all enemy creatures and one third of all
// friendly creatures, rounding each count up. It reads both counts from the
// pre-destruction board, lets the controller choose which creatures on each side
// are destroyed, then destroys the whole chosen set together so their Destroyed
// abilities see one another still in play.
type Tertiate struct{}

// Text renders the effect.
func (e Tertiate) Text() string {
	return "destroy one third of all enemy creatures and one third of all " +
		"friendly creatures (rounding up each time)"
}

// Resolve destroys ceil(E/3) enemy and ceil(F/3) friendly creatures, chosen by
// the controller, all at once.
func (e Tertiate) Resolve(ctx *EffectContext) {
	doomed := e.chooseThird(ctx, ctx.Opponent())
	doomed = append(doomed, e.chooseThird(ctx, ctx.Controller)...)
	Destroy{}.destroy(ctx, doomed)
}

// chooseThird asks the controller to pick ceil(len/3) of player p's creatures,
// one at a time, returning the chosen ids.
func (e Tertiate) chooseThird(ctx *EffectContext, p int) []LocalID {
	line := ctx.Resolver.Battleline(p)
	want := (len(line) + 2) / 3
	picked := map[LocalID]bool{}
	chosen := make([]LocalID, 0, want)
	for len(chosen) < want {
		cands := make([]LocalID, 0, len(line))
		for _, id := range line {
			if !picked[id] {
				cands = append(cands, id)
			}
		}
		pick, ok := ctx.ChooseCard("Choose a creature to destroy", cands)
		if !ok {
			break
		}
		picked[pick] = true
		chosen = append(chosen, pick)
	}
	return chosen
}

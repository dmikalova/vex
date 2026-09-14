package engine

// placeAmong hands pool units back one at a time, each onto a creature the
// controller chooses (falling back to the first when they decline), calling place
// for each unit placed. It is the shared inner loop of the redistribution effects
// — Equalize's Æmber and Entropic Manipulator's damage: the drain and any
// follow-up differ per resource, but handing the pool back one unit at a time does
// not. The caller guarantees creatures is non-empty when pool > 0.
func placeAmong(
	ctx *EffectContext, creatures []LocalID, prompt string, pool int, place func(LocalID),
) {
	for ; pool > 0; pool-- {
		id, ok := ctx.ChooseCreature(prompt, creatures)
		if !ok {
			id = creatures[0]
		}
		place(id)
	}
}

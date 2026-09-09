package engine

import "fmt"

// DestroyAllExceptChosen has the controller keep up to FriendlyKept of their own
// creatures and up to EnemyKept enemy creatures, then destroys every other
// creature on both sides — Unnatural Selection keeps three per side and destroys
// the rest. When a side holds fewer than its keep count, all of that side's
// creatures are kept and none are destroyed there. The survivors are read from the
// pre-destruction board, so the whole doomed set is destroyed at once and their
// Destroyed abilities see one another still in play.
type DestroyAllExceptChosen struct {
	// FriendlyKept and EnemyKept are how many creatures the controller keeps on
	// their own side and the enemy side.
	FriendlyKept int
	EnemyKept    int
}

// validate requires a positive keep count on each side.
func (e DestroyAllExceptChosen) validate() error {
	if e.FriendlyKept <= 0 {
		return fmt.Errorf("DestroyAllExceptChosen: FriendlyKept must be positive")
	}
	if e.EnemyKept <= 0 {
		return fmt.Errorf("DestroyAllExceptChosen: EnemyKept must be positive")
	}
	return nil
}

// Text renders the printed two-sentence effect.
func (e DestroyAllExceptChosen) Text() string {
	return fmt.Sprintf(
		"choose %d friendly creatures and %d enemy creatures. Destroy each "+
			"other creature",
		e.FriendlyKept, e.EnemyKept,
	)
}

// Resolve keeps the controller's chosen creatures on each side, then destroys
// every other creature in play all at once.
func (e DestroyAllExceptChosen) Resolve(ctx *EffectContext) {
	kept := map[LocalID]bool{}
	e.chooseKept(ctx, ctx.Controller, e.FriendlyKept, kept)
	e.chooseKept(ctx, ctx.Opponent(), e.EnemyKept, kept)

	var doomed []LocalID
	for p := 0; p < 2; p++ {
		for _, id := range ctx.Resolver.Battleline(p) {
			if !kept[id] {
				doomed = append(doomed, id)
			}
		}
	}
	Destroy{}.destroy(ctx, doomed)
}

// chooseKept asks the controller to pick min(keep, len) of player p's creatures to
// keep, one at a time, marking each in kept.
func (e DestroyAllExceptChosen) chooseKept(ctx *EffectContext, p, keep int, kept map[LocalID]bool) {
	line := ctx.Resolver.Battleline(p)
	want := keep
	if want > len(line) {
		want = len(line)
	}
	for chosen := 0; chosen < want; chosen++ {
		cands := make([]LocalID, 0, len(line))
		for _, id := range line {
			if !kept[id] {
				cands = append(cands, id)
			}
		}
		pick, ok := ctx.ChooseCard("Choose a creature to keep", cands)
		if !ok {
			break
		}
		kept[pick] = true
	}
}

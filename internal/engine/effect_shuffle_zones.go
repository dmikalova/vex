package engine

// ShuffleChosenCreaturesFromZones shuffles any number of the controller's
// creatures — chosen from their hand, discard pile, and battleline — back into
// their deck as one grouped shuffle. An optional House filter narrows the
// eligible creatures. It is Song of Spring: "Shuffle any number of friendly
// Untamed creatures from your hand, discard pile, or battleline back into your
// deck."
type ShuffleChosenCreaturesFromZones struct {
	// House restricts the choice to creatures it admits; an unset matcher allows any.
	House HouseMatcher
}

// validate has no required fields.
func (e ShuffleChosenCreaturesFromZones) validate() error { return nil }

// Text renders the effect, e.g. "shuffle any number of friendly Untamed creatures
// from your hand, discard pile, or battleline into your deck".
func (e ShuffleChosenCreaturesFromZones) Text() string {
	noun := "friendly " + e.House.qualifyNoun("creatures")
	return "shuffle any number of " + noun +
		" from your hand, discard pile, or battleline into your deck"
}

// Resolve lets the controller pick eligible creatures one at a time until they
// decline or none remain, shuffling each into their deck from whichever zone it
// sits in as one grouped batch.
func (e ShuffleChosenCreaturesFromZones) Resolve(ctx *EffectContext) {
	mover := crossZoneMover{
		Player:  ctx.Controller,
		Dest:    ToDeckShuffled,
		Sources: []Zone{Hand, Discard, inPlay},
	}
	ctx.Resolver.BeginShuffleBatch()
	picked := map[LocalID]bool{}
	for {
		cands := mover.gather(ctx, func(id LocalID) bool {
			return !picked[id] && e.eligible(ctx, id)
		})
		if len(cands) == 0 {
			break
		}
		id, ok := ctx.ChooseCardOptional(
			"Choose a creature to shuffle into your deck", cands)
		if !ok {
			break
		}
		picked[id] = true
		mover.move(ctx, id)
	}
	ctx.Resolver.EndShuffleBatch()
}

// eligible reports whether one card is a creature the optional house filter admits.
func (e ShuffleChosenCreaturesFromZones) eligible(ctx *EffectContext, id LocalID) bool {
	if !ctx.Resolver.IsCreature(id) {
		return false
	}
	return e.House.matches(ctx, id)
}

package engine

// ShuffleChosenCreaturesFromZones shuffles any number of the controller's
// creatures — chosen from their hand, discard pile, and battleline — back into
// their deck as one grouped shuffle. An optional House filter narrows the
// eligible creatures. It is Song of Spring: "Shuffle any number of friendly
// Untamed creatures from your hand, discard pile, or battleline back into your
// deck."
type ShuffleChosenCreaturesFromZones struct {
	// House restricts the choice to creatures of this house; HouseNone allows any.
	House House
}

// validate has no required fields.
func (e ShuffleChosenCreaturesFromZones) validate() error { return nil }

// Text renders the effect, e.g. "shuffle any number of friendly Untamed creatures
// from your hand, discard pile, or battleline into your deck".
func (e ShuffleChosenCreaturesFromZones) Text() string {
	noun := "friendly creatures"
	if e.House != HouseNone {
		noun = "friendly " + e.House.String() + " creatures"
	}
	return "shuffle any number of " + noun +
		" from your hand, discard pile, or battleline into your deck"
}

// Resolve lets the controller pick eligible creatures one at a time until they
// decline or none remain, shuffling each into their deck from whichever zone it
// sits in as one grouped batch.
func (e ShuffleChosenCreaturesFromZones) Resolve(ctx *EffectContext) {
	ctx.Resolver.BeginShuffleBatch()
	picked := map[LocalID]bool{}
	for {
		cands := e.candidates(ctx, picked)
		if len(cands) == 0 {
			break
		}
		id, ok := ctx.ChooseCardOptional(
			"Choose a creature to shuffle into your deck", cands)
		if !ok {
			break
		}
		picked[id] = true
		e.shuffle(ctx, id)
	}
	ctx.Resolver.EndShuffleBatch()
}

// candidates lists the controller's not-yet-chosen creatures in their hand,
// discard pile, and battleline that match the optional house filter.
func (e ShuffleChosenCreaturesFromZones) candidates(
	ctx *EffectContext,
	picked map[LocalID]bool,
) []LocalID {
	var out []LocalID
	consider := func(ids []LocalID) {
		for _, id := range ids {
			if picked[id] || !ctx.Resolver.IsCreature(id) {
				continue
			}
			if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
				continue
			}
			out = append(out, id)
		}
	}
	consider(ctx.Resolver.Hand(ctx.Controller))
	consider(ctx.Resolver.Discard(ctx.Controller))
	consider(ctx.Resolver.Battleline(ctx.Controller))
	return out
}

// shuffle moves one chosen creature into the deck from whichever zone it is in.
func (e ShuffleChosenCreaturesFromZones) shuffle(ctx *EffectContext, id LocalID) {
	if resolverInPlay(ctx, id) {
		ctx.Resolver.PutIntoDeckShuffled(id)
		return
	}
	for _, h := range ctx.Resolver.Hand(ctx.Controller) {
		if h == id {
			ctx.Resolver.ShuffleFromHandIntoDeck(id)
			return
		}
	}
	ctx.Resolver.ShuffleFromDiscardIntoDeck(id)
}

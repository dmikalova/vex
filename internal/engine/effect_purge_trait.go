package engine

import "sort"

// PurgeEachOfChosenTrait lets the controller name a trait, then purges every card
// in play carrying it and pays each player 1 Æmber for every card they controlled
// that the purge took (Harvest Time). Choosing a trait no card in play has is
// pointless, so the choice is limited to the traits actually present.
type PurgeEachOfChosenTrait struct{}

// validate has nothing to reject: the effect takes no fields.
func (e PurgeEachOfChosenTrait) validate() error { return nil }

// Text renders the effect.
func (e PurgeEachOfChosenTrait) Text() string {
	return "choose a trait, then purge each card with that trait. " +
		"Each player gains 1 Æmber for each card they controlled that was purged this way"
}

// Resolve gathers the traits worn by cards in play, asks the controller to pick
// one, purges every card with it, and pays each player for their losses.
func (e PurgeEachOfChosenTrait) Resolve(ctx *EffectContext) {
	// The cards each player controls in play, snapshotted before any purge shifts
	// the battlelines.
	controlled := map[int][]LocalID{}
	present := map[Trait]bool{}
	for _, p := range [2]int{ctx.Controller, ctx.Opponent()} {
		ids := append(ctx.Resolver.Battleline(p), ctx.Resolver.Artifacts(p)...)
		controlled[p] = ids
		for _, id := range ids {
			for tr := traitUnset + 1; tr < traitCount; tr++ {
				if ctx.Resolver.HasTrait(id, tr) {
					present[tr] = true
				}
			}
		}
	}
	if len(present) == 0 {
		return
	}

	traits := make([]Trait, 0, len(present))
	for tr := range present {
		traits = append(traits, tr)
	}
	sort.Slice(traits, func(i, j int) bool { return traits[i].String() < traits[j].String() })
	labels := make([]string, len(traits))
	for i, tr := range traits {
		labels[i] = tr.String()
	}
	chosen := traits[ctx.ChooseOption("Choose a trait", labels)]

	purged := map[int]int{}
	for _, p := range [2]int{ctx.Controller, ctx.Opponent()} {
		for _, id := range controlled[p] {
			if !ctx.Resolver.InPlay(id) || !ctx.Resolver.HasTrait(id, chosen) {
				continue
			}
			ctx.Resolver.PurgeFromPlay(id)
			purged[p]++
		}
	}

	for _, p := range [2]int{ctx.Controller, ctx.Opponent()} {
		if purged[p] == 0 {
			continue
		}
		if capturer, ok := ctx.Resolver.GainAember(p, purged[p]); ok {
			ctx.Resolver.Record(AemberCapturedInsteadOfGain{
				Creature: capturer,
				Player:   p,
				Amount:   purged[p],
			})
			continue
		}
		ctx.Resolver.Record(AemberGained{Player: p, Amount: purged[p]})
	}
}

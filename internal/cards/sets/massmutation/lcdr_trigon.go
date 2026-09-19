package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LCdr. Trigon
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant
//
//	Reap: Discard the top card of your deck. Resolve that card's bonus icons.
var LCdrTrigon = set.New(
	"LCdr. Trigon",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "324"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Reap, card.Sentences{Effects: []card.Effect{
			card.DiscardTop{Amount: 1, Player: card.Controller},
			card.ResolveBonusIcons{Target: card.Target.Triggering},
		}}),
)

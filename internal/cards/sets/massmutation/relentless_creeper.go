package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Relentless Creeper
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	After you choose Dis as your active house, you may put Relentless Creeper from your discard pile into your hand.
var RelentlessCreeper = set.New(
	"Relentless Creeper",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "029"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithTriggersFromDiscard(),
	card.WithAbility(
		card.Trigger.AfterChooseHouse, card.Conditional{
			Cond: card.ChoseHouse{House: card.House.Self},
			Then: card.May{
				Do: card.PutCard{Zones: []card.Zone{card.Discard},
					Selection:   card.Self{},
					Destination: card.To.Hand,
				},
			},
		}),
)

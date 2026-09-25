package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Nature's Call
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Put up to 3 creatures into their owners' hands.
var NaturesCall = set.New(
	"Nature's Call",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "329"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PutChosen{
			Quantity:    card.UpTo{N: card.Fixed(3)},
			Target:      card.Target.EachCreature,
			Destination: card.To.Hand,
		}),
)

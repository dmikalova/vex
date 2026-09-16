package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Dark Queen Gloriana
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Mutant
//
//	Play: Put a friendly non-Untamed creature into its owner's hand.
//	Enhance Æmber Æmber.
var DarkQueenGloriana = set.New(
	"Dark Queen Gloriana",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "397"),
	card.WithEnhance(card.Bonus.Aember, card.Bonus.Aember),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Play, card.PutFromPlay{
			Target:      card.Target.FriendlyCreature.House(card.Houses.Except(card.House.Self)),
			Destination: card.To.Hand,
		}),
)

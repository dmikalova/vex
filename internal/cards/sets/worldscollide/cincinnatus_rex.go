package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Cincinnatus Rex
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  4
//	Traits: Dinosaur • Soldier
//
//	If there are no enemy creatures in play, destroy Cincinnatus Rex.
//	Fight: You may exalt Cincinnatus Rex. Ready each other friendly card.
var CincinnatusRex = set.New(
	"Cincinnatus Rex",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "215"),
	card.WithPower(6),
	card.WithArmor(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithDestroyedWhen(card.InPlay{
		Player: card.Opponent,
		Type:   card.Type.Creature,
		None:   true,
	}),
	card.WithAbility(
		card.Trigger.Fight, card.May{Do: card.Sentences{Effects: []card.Effect{
			card.Exalt{
				Target: card.Target.This,
				Amount: 1,
			},
			card.Ready{Target: card.Target.EachFriendlyCardInPlay.Other()},
		}}}),
)

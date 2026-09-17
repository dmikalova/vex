package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Dexus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon
//
//	After your opponent plays a creature on their right flank, your opponent loses 1 Æmber.
var Dexus = set.New(
	"Dexus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "124"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(card.Trigger.AfterEnemyCardPlayed, card.Conditional{
		Cond: card.OnFlank{OfIt: true, Where: card.RightFlank},
		Then: card.LoseAember{Player: card.Opponent, Amount: 1},
	}),
)

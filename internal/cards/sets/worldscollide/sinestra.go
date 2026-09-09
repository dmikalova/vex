package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Sinestra
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon
//
//	After your opponent plays a creature on their left flank, they lose 1A.
var Sinestra = card.New(
	"Sinestra",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "116"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(card.Trigger.AfterEnemyCardPlayed, card.Conditional{
		Cond: card.ItIsOnNamedFlank{},
		Then: card.LoseAember{Player: card.Opponent, Amount: 1},
	}),
)

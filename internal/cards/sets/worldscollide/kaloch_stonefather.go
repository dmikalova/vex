package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Kaloch Stonefather
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant • Leader
//
//	While Kaloch Stonefather is in the center of your battleline, each friendly creature gains skirmish.
var KalochStonefather = set.New(
	"Kaloch Stonefather",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "41"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant, card.Traits.Leader),
	card.WithConstant(card.ConstantAbility{
		Target:        card.Target.EachFriendlyCreature,
		Keywords:      card.Keywords(card.Keyword.Skirmish),
		WhileInCenter: true,
	}),
)

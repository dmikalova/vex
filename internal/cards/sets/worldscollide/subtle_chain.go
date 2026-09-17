package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Subtle Chain
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Your opponent discards a random card from their hand.
var SubtleChain = set.New(
	"Subtle Chain",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "262"),
	card.InCluster(card.Pulled(chainGangCluster, 1, 2)),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.DiscardCard{
			Player:    card.Opponent,
			Zones:     []card.Zone{card.Hand},
			Selection: card.Random{Count: 1},
		},
	),
)

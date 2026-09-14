package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Subtle Chain
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent discards a random card from their hand.
var SubtleChain = set.New(
	"Subtle Chain",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "262"),
	card.InCluster(card.Pulled(chainGangCluster, 1, 2)),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.DiscardCard{Player: card.Opponent, Zone: card.Hand, Selection: card.Random{}},
	),
)

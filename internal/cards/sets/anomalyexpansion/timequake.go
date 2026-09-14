package anomalyexpansion

import "github.com/dmikalova/vactrol/internal/card"

// Timequake
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Shuffle each friendly card in play into your deck. For each card shuffled into your deck this way, draw a card.
var Timequake = set.New(
	"Timequake",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.WC, "A09"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.ShuffleFriendlyCardsIntoDeck{},
			card.Draw{Amount: 1, Per: card.CardsShuffledIntoDeck{}},
		}}),
)

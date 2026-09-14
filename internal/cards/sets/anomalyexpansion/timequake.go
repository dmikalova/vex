package anomalyexpansion

import "github.com/dmikalova/vactrol/internal/card"

// Timequake
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Shuffle each friendly card in play into your deck. Draw a card for each card shuffled into your deck this way.
var Timequake = set.New(
	"Timequake",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.WC, "A09"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Shuffle{FromPlay: true}),
)

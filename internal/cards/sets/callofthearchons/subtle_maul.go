package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Subtle Maul
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Your opponent discards a random card from their hand.
var SubtleMaul = card.New(
	"Subtle Maul",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 294),
	card.WithTraits("Weapon"),
	card.WithAbility(
		card.Trigger.Action, card.DiscardRandom{Player: card.Opponent}),
)

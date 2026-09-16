package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Scrivener Favian
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	When you resolve a Capture bonus icon, steal 1 Æmber instead.
//	Enhance Capture Capture.
var ScrivenerFavian = set.New(
	"Scrivener Favian",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "155"),
	card.WithEnhance(card.Bonus.Capture, card.Bonus.Capture),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	card.WithBonusInstead(card.BonusInstead{
		From:    card.Bonus.Capture,
		Instead: card.StealAember{Amount: 1},
	}),
)

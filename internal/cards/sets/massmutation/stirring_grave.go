package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Stirring Grave
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Archive a creature from your discard pile.
var StirringGrave = set.New(
	"Stirring Grave",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "015"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveCard{
			Zone:      card.Discard,
			Selection: card.Chosen{Type: card.Type.Creature},
		}),
)

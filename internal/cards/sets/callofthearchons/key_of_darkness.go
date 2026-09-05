package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Key of Darkness
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Forge a key at +6 Æmber current cost, or +2 if your opponent has no Æmber.
var KeyOfDarkness = card.New(
	"Key of Darkness",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 273),
	card.WithAbility(
		card.Trigger.Play, card.ForgeKey{
			Extra: 6,
			Or: card.OrAmount{
				Amount: 2,
				When:   card.OpponentAember{Is: card.Exactly, Amount: 0},
			},
		}),
)

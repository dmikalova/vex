package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Forging an Alliance
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Forge a key at +7 Æmber current cost, reduced by 1 Æmber for each house represented among cards in play.
var ForgingAnAlliance = card.New(
	"Forging an Alliance",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "331"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ForgeKey{
			Extra: 7,
			ReducedBy: card.HousesAmong{
				Player: card.EachPlayer,
			},
		}),
)

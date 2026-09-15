package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Forging an Alliance
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Forge a key at +7 Æmber current cost, reduced by 1 Æmber for each house represented among cards in play -> purge Forging an Alliance.
var ForgingAnAlliance = set.New(
	"Forging an Alliance",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "331"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ForgeKey{
			Extra: 7,
			ReducedBy: card.HousesAmong{
				Player: card.EachPlayer,
			},
		}),
)

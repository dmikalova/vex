package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Data Forge
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Forge a key at +10 Æmber current cost, reduced by 1 Æmber for each card in your hand -> purge Data Forge.
var DataForge = card.New(
	"Data Forge",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "148"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ForgeKey{
			Extra: 10,
			ReducedBy: card.CardsInHand{
				Player: card.Controller,
				House:  card.AnyHouse,
			},
		}),
)

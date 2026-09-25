package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Data Forge
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Forge a key at +10 Æmber current cost, reduced by 1 Æmber for each card in your hand -> purge Data Forge.
var DataForge = set.New(
	"Data Forge",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "148"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ForgeKey{
			Extra: 10,
			ReducedBy: card.CardsInHand{
				Player: card.Controller,
				House:  card.AnyHouse,
			},
		}),
)

package massmutation

import "github.com/dmikalova/vex/internal/card"

// Vault's Blessing
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: For each friendly Mutant creature, each player gains 1 Æmber.
var VaultsBlessing = set.New(
	"Vault's Blessing",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "391"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.EachPlayer,
			Amount: 1,
			Per: card.CardsInPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				Trait:  card.Traits.Mutant,
			},
		}),
)

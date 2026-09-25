package massmutation

import "github.com/dmikalova/vex/internal/card"

// Hedonistic Intent
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Exalt each flank creature.
var HedonisticIntent = set.New(
	"Hedonistic Intent",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "207"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.EachCreature.OnFlank(),
			Amount: 1,
		}),
)

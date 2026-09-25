package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Diplo-Macy
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Alpha.
//	Play: Until the start of your next turn, each creature gains, "Before Fight: Exalt this creature."
var DiploMacy = set.New(
	"Diplo-Macy",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "218"),
	card.WithBonus(card.Bonus.Aember),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.GainAbility{
			Target:   card.Target.EachCreature,
			Duration: card.Duration.StartOfPlayerNextTurn,
			Ability: card.Ability{
				Trigger: card.Trigger.BeforeFight,
				Effect: card.Exalt{
					Target: card.Target.This,
					Amount: 1,
				},
			},
		}),
)

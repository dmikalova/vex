package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Diplo-Macy
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Alpha.
//	Play: Until the start of your next turn, each Creature gains, "Before Fight: Exalt this Creature."
var DiploMacy = set.New(
	"Diplo-Macy",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "218"),
	card.WithAemberBonus(1),
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

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Adaptoid
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	After you play a card with a bonus icon, choose one:
//	- For the remainder of the turn, Adaptoid gains +2 armor
//	- For the remainder of the turn, Adaptoid gains assault 2
//	- For the remainder of the turn, Adaptoid gains, "Fight: Steal 1 Æmber."
//	Enhance Capture Damage Draw.
var Adaptoid = set.New(
	"Adaptoid",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "099"),
	card.WithEnhance(card.Bonus.Capture, card.Bonus.Damage, card.Bonus.Draw),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.AfterCardPlayed, card.Conditional{
			Cond: card.ItHasBonusIcon{},
			Then: card.ChooseOne{
				Options: []card.Effect{
					card.GainStats{
						Target: card.Target.This,
						Armor:  2,
					},
					card.GainAssault{
						Target: card.Target.This,
						Amount: card.Fixed(2),
					},
					card.GainAbility{
						Target:   card.Target.This,
						Duration: card.Duration.RemainderOfPlayerTurn,
						Ability: card.Ability{
							Trigger: card.Trigger.Fight,
							Effect:  card.StealAember{Amount: 1},
						},
					},
				},
			},
		}),
)

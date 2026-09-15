package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mad Prophet Gizelhart
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  3
//	Traits: Leader • Priest
//
//	Action: If Mad Prophet Gizelhart is in the center of your battleline, fully heal each non-Mutant creature. For each creature healed this way, gain 1 Æmber.
var MadProphetGizelhart = set.New(
	"Mad Prophet Gizelhart",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "167"),
	card.WithPower(4),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Leader, card.Traits.Priest),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.SourceInCenterOfBattleline{},
			Then: card.Sentences{
				Effects: []card.Effect{
					card.Heal{
						Fully:  true,
						Target: card.Target.EachCreature.ExceptTrait(card.Traits.Mutant),
					},
					card.GainAember{
						Player: card.Controller,
						Amount: 1,
						Per:    card.CreaturesHealed{},
					},
				},
			},
		}),
)

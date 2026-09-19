package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mutation of Fury
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Choose a creature. The chosen creature gains assault 3 and the Mutant trait until the start of your next turn.
var MutationOfFury = set.New(
	"Mutation of Fury",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Connected,
	card.Provenance(card.MM, "414"),
	card.InCluster(darkHarbingerCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.GainUntilNextTurn{
				Effects: []card.Effect{
					card.GainAssaultUntilNextTurn{
						Target: card.Target.TheChosenCreature,
						Amount: card.Fixed(3),
					},
					card.GainTrait{
						Target: card.Target.TheChosenCreature,
						Trait:  card.Traits.Mutant,
					},
				},
			},
		}),
)

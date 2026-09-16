package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mutation of Instinct
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Choose a creature - the chosen creature gains skirmish and the Mutant trait until the start of your next turn.
var MutationOfInstinct = set.New(
	"Mutation of Instinct",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Connected,
	card.Provenance(card.MM, "415"),
	card.InCluster(darkHarbingerCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.GainUntilNextTurn{
				Effects: []card.Effect{
					card.GainKeyword{
						Target:  card.Target.TheChosenCreature,
						Keyword: card.Keyword.Skirmish,
					},
					card.GainTrait{
						Target: card.Target.TheChosenCreature,
						Trait:  card.Traits.Mutant,
					},
				},
			},
		}),
)

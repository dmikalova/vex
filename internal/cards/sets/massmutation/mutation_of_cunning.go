package massmutation

import "github.com/dmikalova/vex/internal/card"

// Mutation of Cunning
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Choose a creature. The chosen creature gains elusive and the Mutant trait until the start of your next turn.
var MutationOfCunning = set.New(
	"Mutation of Cunning",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Connected,
	card.Provenance(card.MM, "413"),
	card.InCluster(darkHarbingerCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.GainUntilNextTurn{
				Effects: []card.Effect{
					card.GainKeywords{
						Target:   card.Target.TheChosenCreature,
						Keywords: []card.KeywordValue{card.Keyword.Elusive},
						Duration: card.Duration.StartOfPlayerNextTurn,
					},
					card.GainTrait{
						Target: card.Target.TheChosenCreature,
						Trait:  card.Traits.Mutant,
					},
				},
			},
		}),
)

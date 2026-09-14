package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Creed of Nature
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Creed of Nature. Choose a Creature - for the remainder of the turn, it gains skirmish and assault equal to its power.
var CreedOfNature = set.New(
	"Creed of Nature",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "385"),
	card.WithTraits(card.Traits.Power),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{
			Effects: []card.Effect{
				card.Destroy{Target: card.Target.This},
				card.ChooseCreatureThen{
					Target: card.Target.Creature,
					Then: card.ForDuration{
						Duration: card.Duration.RemainderOfPlayerTurn,
						Effects: []card.Effect{
							card.GainKeywordForTurn{
								Target:  card.Target.Triggering,
								Keyword: card.Keyword.Skirmish,
							},
							card.GainAssault{
								Target: card.Target.Triggering,
								Amount: card.PowerOfChosen{},
							},
						},
					},
				},
			},
		}),
)

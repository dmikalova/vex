package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Spectral Tunneler
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Choose a creature - for the remainder of the turn, it is considered a flank creature, and it gains, "Reap: Draw a card."
var SpectralTunneler = card.New(
	"Spectral Tunneler",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 133),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Sequence{Effects: []card.Effect{
				card.ConsiderFlank{Target: card.Target.Triggering},
				card.GainAbility{
					Target: card.Target.Triggering,
					Ability: card.Ability{
						Trigger: card.Trigger.Reap,
						Effect:  card.Draw{Amount: 1},
					},
				},
			}},
		}),
)

package massmutation

import "github.com/dmikalova/vex/internal/card"

// Resurgence
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Put a creature from your discard pile into your hand. If it is a Mutant creature, put another creature from your discard pile into your hand.
//	Enhance Draw.
var Resurgence = set.New(
	"Resurgence",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "375"),
	card.WithEnhance(card.Bonus.Draw),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PutCard{Zones: []card.Zone{card.Discard},
				Selection:   card.Chosen{Filter: card.Filter{Type: card.Type.Creature}},
				Destination: card.To.Hand,
				Bind:        true,
			},
			card.Conditional{
				Cond: card.ItIs{
					Filter: card.Filter{Type: card.Type.Creature, Trait: card.Traits.Mutant},
				},
				Then: card.PutCard{
					Zones: []card.Zone{card.Discard},
					Selection: card.Chosen{
						Filter: card.Filter{Type: card.Type.Creature, Except: card.Except.It},
					},
					Destination: card.To.Hand,
				},
			},
		}}),
)

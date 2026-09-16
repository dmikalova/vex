package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Resurgence
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Put a creature from your discard pile into your hand. If it is a Mutant creature, put a creature from your discard pile into your hand.
//	Enhance Draw.
var Resurgence = set.New(
	"Resurgence",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "375"),
	card.WithEnhance(card.Bonus.Draw),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.PutFromDiscard{
				Selection:   card.Chosen{Type: card.Type.Creature},
				Destination: card.To.Hand,
				Bind:        true,
			},
			card.Conditional{
				Cond: card.ItIsOfTrait{Trait: card.Traits.Mutant},
				Then: card.PutFromDiscard{
					Selection:   card.Chosen{Type: card.Type.Creature},
					Destination: card.To.Hand,
				},
			},
		}}),
)

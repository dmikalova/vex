package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mighty Javelin
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Mighty Javelin. Deal 4 damage to a creature.
var MightyJavelin = card.New(
	"Mighty Javelin",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 24),
	card.WithAemberBonus(1),
	card.WithTraits("Weapon"),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{
			Effects: []card.Effect{
				card.Destroy{Target: card.Target.This},
				card.DealDamage{
					Amount: 4,
					Target: card.Target.Creature,
				},
			},
		}),
)

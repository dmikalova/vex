package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mothergun
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Reveal any number of Mars cards from your hand, and for each card revealed this way, deal 1 damage to a creature.
var Mothergun = card.New(
	"Mothergun",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, 187),
	card.WithTraits("Weapon"),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.Reveal{
					Player: card.Controller,
					House:  card.House.Mars,
				},
				card.DealDamage{
					Amount: 1,
					Target: card.Target.Creature,
					Per:    card.CardsRevealed{},
				},
			},
		}),
)

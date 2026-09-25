package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Mothergun
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Reveal any number of Mars cards from your hand. For each card revealed this way, deal 1 damage to a creature.
var Mothergun = set.New(
	"Mothergun",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, "187"),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.RevealHand{
					Player: card.Controller,
					House:  card.Houses.Named(card.House.Self),
				},
				card.DealDamage{
					Amount: 1,
					Target: card.Target.Creature,
					Per:    card.CardsRevealed{},
				},
			},
		}),
)

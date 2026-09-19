package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Orbital Bombardment
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Reveal any number of Mars cards from your hand. For each card revealed this way, deal 2 damage to a creature.
var OrbitalBombardment = set.New(
	"Orbital Bombardment",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "172"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.RevealHand{
					Player: card.Controller,
					House:  card.Houses.Named(card.House.Self),
				},
				card.ForEach{
					Times: card.CardsRevealed{},
					Do: card.DealDamage{
						Target: card.Target.Creature,
						Amount: 2,
					},
				},
			},
		}),
)

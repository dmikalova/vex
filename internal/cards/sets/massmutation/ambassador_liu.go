package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Ambassador Liu
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant • Politician
//
//	Action: Discard a card from your hand. If it is a Dis or Shadows card, steal 1 Æmber. If it is a Logos or Untamed card, gain 2 Æmber. If it is a Sanctum or Saurian card, Ambassador Liu captures 3 Æmber from your opponent.
var AmbassadorLiu = set.New(
	"Ambassador Liu",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "335"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Politician),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.DiscardCard{
				Player:    card.Controller,
				Zones:     []card.Zone{card.Hand},
				Selection: card.Chosen{},
				Bind:      true,
			},
			card.Conditional{
				Cond: card.Or{Conditions: []card.Condition{
					card.ItIs{House: card.Houses.Named(card.House.Dis)},
					card.ItIs{House: card.Houses.Named(card.House.Shadows)},
				}},
				Then: card.StealAember{Amount: 1},
			},
			card.Conditional{
				Cond: card.Or{Conditions: []card.Condition{
					card.ItIs{House: card.Houses.Named(card.House.Logos)},
					card.ItIs{House: card.Houses.Named(card.House.Untamed)},
				}},
				Then: card.GainAember{
					Player: card.Controller,
					Amount: 2,
				},
			},
			card.Conditional{
				Cond: card.Or{Conditions: []card.Condition{
					card.ItIs{House: card.Houses.Named(card.House.Sanctum)},
					card.ItIs{House: card.Houses.Named(card.House.Saurian)},
				}},
				Then: card.CaptureAember{
					Amount: 3,
					Target: card.Target.This,
					Source: card.Opponent,
				},
			},
		}}),
)

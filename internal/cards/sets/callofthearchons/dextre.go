package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dextre
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Play: Dextre captures 1 Æmber from your opponent.
//	Destroyed: Put Dextre on top of its owner's deck.
var Dextre = card.New(
	"Dextre",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 138),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
	card.WithAbility(
		card.Trigger.Destroyed, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.TopOfDeck,
		}),
)

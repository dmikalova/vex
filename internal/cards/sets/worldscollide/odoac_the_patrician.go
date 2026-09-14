package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Odoac the Patrician
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Dinosaur • Politician
//
//	While Odoac the Patrician has Æmber on it, your Æmber cannot be stolen.
//	Play: Odoac the Patrician captures 1 Æmber from your opponent.
var OdoacThePatrician = set.New(
	"Odoac the Patrician",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "188"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithAemberCannotBeStolenWhileItHasAember(),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)

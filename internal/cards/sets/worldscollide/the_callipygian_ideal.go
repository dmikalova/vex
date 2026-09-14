package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Callipygian Ideal
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This Creature gains, "You may spend Æmber on this Creature as if it were in your pool."
//	Play: Exalt this Creature.
var TheCallipygianIdeal = set.New(
	"The Callipygian Ideal",
	card.House.Saurian,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "212"),
	card.WithStatic(card.StaticModifier{SpendAemberOnCard: true}),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.This,
			Amount: 1,
		}),
)

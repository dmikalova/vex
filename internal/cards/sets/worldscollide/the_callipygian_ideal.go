package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Callipygian Ideal
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "You may spend Æmber on this creature as if it were in your pool."
//	Play: Exalt this creature.
var TheCallipygianIdeal = set.New(
	"The Callipygian Ideal",
	card.House.Saurian,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "212"),
	card.WithStatic(card.StaticModifier{SpendAemberOnCard: card.SpendScope.Controller}),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.This,
			Amount: 1,
		}),
)

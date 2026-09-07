//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// C.Ae.N.D.L.E. Unit
var CAeNDLEUnit = card.New(
	"C.Ae.N.D.L.E. Unit",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 146),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureReaps, card.Draw{Amount: 1}),
	card.WithAbility(
		card.Trigger.Action, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)

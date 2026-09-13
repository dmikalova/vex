package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quant
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Reap: Play a non-Logos Tactic.
var Quant = card.New(
	"Quant",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "137"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.PlayFrom{
			From:   card.Hand,
			House:  card.House.Self,
			Except: true,
			Type:   card.Type.Tactic,
		}),
)

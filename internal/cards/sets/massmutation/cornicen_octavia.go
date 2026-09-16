package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Cornicen Octavia
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Action: Cornicen Octavia captures 2 Æmber from your opponent.
var CornicenOctavia = set.New(
	"Cornicen Octavia",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "188"),
	card.InCluster(card.Pulled(monumentToOctaviaCluster, 1, 1)),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Action, card.CaptureAember{
			Amount: 2,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)

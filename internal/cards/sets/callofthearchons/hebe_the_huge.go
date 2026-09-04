package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hebe the Huge
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Giant • Knight
//
//	Play: Deal 2 damage to each other undamaged creature.
var HebeTheHuge = card.New(
	"Hebe the Huge",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 36),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachCreature.Other().Undamaged(),
		}),
)

//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// YzphyzKnowdrone
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Armor:  1
//	Traits: Martian • Scientist
//
//	Play: Archive a card. You may purge an archived card to stun a creature.
var YzphyzKnowdrone = card.New(
	"Yzphyz Knowdrone",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 210),
	card.WithPower(3),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)

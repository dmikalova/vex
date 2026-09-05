//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// CybergiantRig
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "At the end of your turn, this creature loses a +1 power counter."
//	Play: Fully heal this creature and give it a +1 power counter for each damage healed.
var CybergiantRig = card.New(
	"Cybergiant Rig",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 37),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)

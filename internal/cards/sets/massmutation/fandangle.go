package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Fandangle
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Witch
//
//	While you have 4 or more Æmber, your non-Untamed creatures enter play ready.
var Fandangle = set.New(
	"Fandangle",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "365"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Witch),
	card.WithFriendlyEntersPlayReady(card.EntersReadyGrant{
		Type:        card.Type.Creature,
		MinAember:   4,
		ExceptHouse: card.House.Self,
	}),
)

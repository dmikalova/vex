package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Babbling Bibliophile
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Cyborg • Scientist
//
//	Reap: Draw 2 cards.
var BabblingBibliophile = set.New(
	"Babbling Bibliophile",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "129"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.Draw{Amount: 2}),
)

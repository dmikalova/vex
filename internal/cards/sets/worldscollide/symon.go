package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Symon
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Alien • Thief
//
//	Skirmish.
//	Fight: Put the Creature Symon fought on top of its owner's deck.
var Symon = card.New(
	"Symon",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "247"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Alien, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.PutFromPlay{
			Target:      card.Target.CreatureFought,
			Destination: card.To.TopOfDeck,
		}),
)

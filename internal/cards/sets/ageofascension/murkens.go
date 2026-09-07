package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Murkens
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Choose one:
//	- Play a random card from your opponent's archives
//	- Play the top card of your opponent's deck.
var Murkens = card.New(
	"Murkens",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "290"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{
			Options: []card.Effect{
				card.PlayRandomFromOpponentArchives{},
				card.PlayTopOfOpponentDeck{},
			},
		}),
)

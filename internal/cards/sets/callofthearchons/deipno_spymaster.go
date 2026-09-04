package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Deipno Spymaster
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive, Versatile.
//	Action: Use a friendly creature.
var DeipnoSpymaster = card.New(
	"Deipno Spymaster",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 299),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive, card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature,
			Verbs:  []card.CreatureVerb{card.UseVerb{}},
		}),
)

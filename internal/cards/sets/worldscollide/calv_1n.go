package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CALV-1N
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Robot
//
//	Fight/Reap: Draw a card.
//	CALV-1N may be played as an upgrade instead of a creature, with the text: "This creature gains, 'Fight/Reap: Draw a card.'"
var CALV1N = set.New(
	"CALV-1N",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "308"),
	card.WithPower(2),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Robot),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Fight, Effect: card.Draw{Amount: 1}},
			{Trigger: card.Trigger.Reap, Effect: card.Draw{Amount: 1}},
		},
	}),
	card.WithPlayableAsUpgrade(),
	card.WithAbility(card.Trigger.FightReap, card.Draw{Amount: 1}),
)

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Armory Officer Nel
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Alien
//
//	After an upgrade enters play, draw a card.
//	Enhance Draw.
var ArmoryOfficerNel = set.New(
	"Armory Officer Nel",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "319"),
	card.WithEnhance(card.Bonus.Draw),
	card.WithPower(4),
	card.WithTraits(card.Traits.Alien),
	card.WithAbility(
		card.Trigger.AfterUpgradeEnters, card.Draw{Amount: 1}),
)

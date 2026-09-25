package massmutation

import "github.com/dmikalova/vex/internal/card"

// Bot Bookton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Reap: Play the top card of your deck.
var BotBookton = set.New(
	"Bot Bookton",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "067"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	card.WithAbility(card.Trigger.Reap, card.PlayTopOfDeck{}),
)

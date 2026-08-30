package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Witch of the Eye
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Witch
//
//	Reap: Put a card from your discard pile into your hand.
var WitchOfTheEye = card.New(
	"Witch of the Eye",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 368),
	card.WithPower(3),
	card.WithTraits("Human", "Witch"),
	card.WithAbility(
		card.Trigger.Reap, card.PutFromDiscard{Destination: card.To.Hand}),
)

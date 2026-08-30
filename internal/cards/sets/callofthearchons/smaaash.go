package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Smaaash
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Play: Stun a creature.
var Smaaash = card.New(
	"Smaaash",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 46),
	card.WithPower(5),
	card.WithTraits("Giant"),
	card.WithAbility(
		card.Trigger.Play, card.Stun{Target: card.Target.Creature}),
)

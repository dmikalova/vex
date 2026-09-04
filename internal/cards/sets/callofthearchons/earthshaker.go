package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Earthshaker
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Giant
//
//	Play: Destroy each creature with power 3 or lower.
var Earthshaker = card.New(
	"Earthshaker",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 31),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{Target: card.Target.EachCreature.PowerAtMost(3)}),
)

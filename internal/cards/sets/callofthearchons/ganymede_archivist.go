package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ganymede Archivist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Reap: Archive a card from your hand.
var GanymedeArchivist = card.New(
	"Ganymede Archivist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 142),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(card.Trigger.Reap, card.ArchiveFromHand{Count: 1}),
)

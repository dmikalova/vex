package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Incubation Chamber
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Versatile.
//	Action: Reveal a Mars creature from your hand and archive it.
var IncubationChamber = card.New(
	"Incubation Chamber",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 186),
	card.WithTraits(card.Traits.Location),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.ArchiveFromHand{
			Count:    1,
			Type:     card.Type.Creature,
			House:    card.House.Self,
			Revealed: true,
		}),
)

package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Novu Archaeologist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	Action: Archive a card from your discard pile.
var NovuArchaeologist = set.New(
	"Novu Archaeologist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "147"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Action,
		card.ArchiveCard{
			Zone:      card.Discard,
			Selection: card.Chosen{},
		},
	),
)

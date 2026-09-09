package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Information Officer Gray
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Play/Fight/Reap: You may reveal a non-Star Alliance card from your hand and archive it.
var InformationOfficerGray = card.New(
	"Information Officer Gray",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "312"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithPlayFightReap(card.May{
		Do: card.ArchiveFromHand{
			Amount:      1,
			Revealed:    true,
			ExceptHouse: card.House.Self,
		},
	}),
)

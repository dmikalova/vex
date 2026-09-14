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
var InformationOfficerGray = set.New(
	"Information Officer Gray",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "312"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(card.Trigger.PlayFightReap, card.May{
		Do: card.ArchiveCard{
			Zone:      card.Hand,
			Selection: card.Chosen{House: card.Houses.Except(card.House.Self), Optional: true},
			Revealed:  true,
		},
	}),
)

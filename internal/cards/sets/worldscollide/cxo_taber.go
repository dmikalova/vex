package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CXO Taber
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Alien • Krxix
//
//	Fight/Reap: Play or use a non-Star Alliance card.
var CXOTaber = set.New(
	"CXO Taber",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "309"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Alien, card.Traits.Krxix),
	card.WithAbility(card.Trigger.FightReap, card.PlayOrUse{
		House: card.Houses.Except(card.House.Self),
	}),
)

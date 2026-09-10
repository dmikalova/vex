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
//	Fight/Reap: You may play or use one non-Star Alliance card this turn.
var CXOTaber = card.New(
	"CXO Taber",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "309"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Alien, card.Traits.Krxix),
	card.WithFightOrReap(card.MayPlayOffHouse{
		Except: card.House.Self,
		Grant:  card.GrantPlay | card.GrantUse,
		Count:  1,
	}),
)

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Lieutenant Khrkhar
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Alien • Handuhan
//
//	Taunt, Hazardous 3.
var LieutenantKhrkhar = card.New(
	"Lieutenant Khrkhar",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 299),
	card.WithPower(5),
	card.WithTraits(card.Traits.Alien, card.Traits.Handuhan),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithHazardous(3),
)

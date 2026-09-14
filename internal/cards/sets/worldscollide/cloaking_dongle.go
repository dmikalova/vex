package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Cloaking Dongle
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This Creature and each of its neighbors gains elusive.
var CloakingDongle = set.New(
	"Cloaking Dongle",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, "294"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		KeywordGrants: []card.KeywordGrant{{
			Keywords:  card.Keywords(card.Keyword.Elusive),
			Host:      true,
			Neighbors: true,
		}},
	}),
)

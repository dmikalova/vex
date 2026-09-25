package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Custom Virus
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Custom Virus. Purge a creature from your hand. Destroy each creature that shares a trait with it.
var CustomVirus = set.New(
	"Custom Virus",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "183"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Weapon),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.PurgeCard{
				Zones:     []card.Zone{card.Hand},
				Player:    card.Controller,
				Selection: card.Chosen{Type: card.Type.Creature},
			},
			card.Destroy{Target: card.Target.EachCreature.SharingTrait()},
		}}),
)

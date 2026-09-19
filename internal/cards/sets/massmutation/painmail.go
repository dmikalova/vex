package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Painmail
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "After a player chooses Dis as their active house, archive Painmail. Destroy this creature."
var Painmail = set.New(
	"Painmail",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.MM, "042"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger:    card.Trigger.AfterChooseHouse,
			EachPlayer: true,
			Effect: card.Conditional{
				Cond: card.ChoseHouse{House: card.House.Dis},
				Then: card.Sequence{Effects: []card.Effect{
					card.ArchiveGrantingUpgrade{},
					card.Destroy{Target: card.Target.This},
				}},
			},
		}},
	}),
)

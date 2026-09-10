package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tertiate
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy one third of all enemy creatures and one third of all friendly creatures (rounding up each time).
var Tertiate = card.New(
	"Tertiate",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "232"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Refine(card.PortionPerSide(card.ThirdRoundedUp)),
		}),
)

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tendrils of Pain
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to each creature. If your opponent forged a key during their previous turn, deal 3 damage to each creature.
var TendrilsOfPain = set.New(
	"Tendrils of Pain",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "64"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.DealDamage{
				Target: card.Target.EachCreature,
				Amount: 1,
			},
			card.Conditional{
				Cond: card.ForgedKey{
					Player:   card.Opponent,
					Previous: true,
				},
				Then: card.DealDamage{
					Target: card.Target.EachCreature,
					Amount: 3,
				},
			},
		}}),
)

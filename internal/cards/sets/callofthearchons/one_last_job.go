package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// One Last Job
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Purge each friendly Shadows creature. For each creature purged this way, steal 1 Æmber.
var OneLastJob = set.New(
	"One Last Job",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "277"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PurgeCreature{
				Target: card.Target.EachFriendlyCreature.House(card.Houses.Named(card.House.Self)),
			},
			card.StealAember{
				Amount: 1,
				Per:    card.CardsPurged{Type: card.Type.Creature},
			},
		}}),
)

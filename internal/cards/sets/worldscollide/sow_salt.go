package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Sow Salt
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Alpha.
//	Play: Until the start of your next turn, creatures cannot be used to reap.
var SowSalt = set.New(
	"Sow Salt",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "230"),
	card.WithBonus(card.Bonus.Aember),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.CreaturesCannot{
			Action:   card.UseKind.Reap,
			Duration: card.Duration.StartOfPlayerNextTurn,
		}),
)

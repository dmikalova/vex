package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Hold the Line
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: For each creature your opponent controls in excess of you, draw a card.
var HoldTheLine = set.New(
	"Hold the Line",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "346"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Draw{
			Amount: 1,
			Per:    card.ExcessCreatures{Player: card.Opponent},
		}),
)

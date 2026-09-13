package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Saury About That
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy a Creature -> its controller gains 1 Æmber.
var SauryAboutThat = card.New(
	"Saury About That",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "228"),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First:  card.Destroy{Target: card.Target.Creature},
			Result: card.GainAember{Player: card.ItsController, Amount: 1},
		}),
)

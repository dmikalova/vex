package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Disruption Field
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Your opponent's keys cost +1 Æmber for each disruption counter on Disruption Field.
//	This creature gains, "Fight/Reap: Put a disruption counter on Disruption Field."
var DisruptionField = set.New(
	"Disruption Field",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "328"),
	card.WithBonus(card.Bonus.Aember),
	card.WithKeyCost(
		card.KeyCostChange(card.Opponent, 1).
			Per(card.CountersOnThis{Kind: card.Counter.Disruption}),
	),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.PlaceCounter{Amount: 1,
			Kind:   card.Counter.Disruption,
			Target: card.Target.GrantingCard,
		}),
	}),
)

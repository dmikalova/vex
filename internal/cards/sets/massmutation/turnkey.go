package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Turnkey
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Demon
//
//	Play: Unforge one of your opponent's keys -> when Turnkey leaves play, your opponent forges a key at no cost.
var Turnkey = set.New(
	"Turnkey",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "051"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.UnforgeKey{Player: card.Opponent},
			Result: card.ScheduleOnLeave{
				Do: card.ForgeKey{Player: card.Opponent, FreeOfCost: true},
			},
		}),
)

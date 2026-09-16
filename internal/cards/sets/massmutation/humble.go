package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Humble
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Exhaust a creature -> move 3 Æmber from the chosen creature to the common supply.
var Humble = set.New(
	"Humble",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "208"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.Exhaust{
				Target: card.Target.Creature,
				Bind:   true,
			},
			Result: card.MoveAemberToSupply{
				Amount: 3,
				Target: card.Target.TheChosenCreature,
			},
		}),
)

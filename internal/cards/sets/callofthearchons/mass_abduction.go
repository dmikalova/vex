package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mass Abduction
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Put up to 3 enemy damaged creatures into your archives.
var MassAbduction = card.New(
	"Mass Abduction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 169),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PutChosen{
			Amount:      3,
			UpTo:        true,
			Target:      card.Target.EachEnemyCreature.Damaged(),
			Destination: card.To.Archives.Yours(),
		}),
)

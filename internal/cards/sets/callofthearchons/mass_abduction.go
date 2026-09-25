package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Mass Abduction
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Put up to 3 enemy damaged creatures into your archives.
var MassAbduction = set.New(
	"Mass Abduction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "169"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PutChosen{
			Quantity:    card.UpTo{N: card.Fixed(3)},
			Target:      card.Target.EachEnemyCreature.Damaged(),
			Destination: card.To.Archives.Yours(),
		}),
)

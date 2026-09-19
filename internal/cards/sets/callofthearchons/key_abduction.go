package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Key Abduction
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Put each Mars creature into its owner's hand. Forge a key at +9 Æmber current cost, reduced by 1 Æmber for each card in your hand -> purge Key Abduction.
var KeyAbduction = set.New(
	"Key Abduction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "166"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PutFromPlay{
				Target:      card.Target.EachCreature.House(card.Houses.Named(card.House.Self)),
				Destination: card.To.Hand,
			},
			card.ForgeKey{
				Extra: 9,
				ReducedBy: card.CardsInHand{
					Player: card.Controller,
					House:  card.AnyHouse,
				},
			},
		}}),
)

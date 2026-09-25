package massmutation

import "github.com/dmikalova/vex/internal/card"

// Survey
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Look at the top 2 cards of your deck and discard 1.
//	Enhance Draw.
var Survey = set.New(
	"Survey",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "316"),
	card.WithBonus(card.Bonus.Aember),
	card.WithEnhance(card.Bonus.Draw),
	card.WithAbility(
		card.Trigger.Play, card.LookAtTopOfDeck{
			Amount: 2,
			Then: []card.TopAct{
				card.ChooseAndMove{
					Cards: 1,
					Dest:  card.Into.Discard,
				},
			},
		}),
)

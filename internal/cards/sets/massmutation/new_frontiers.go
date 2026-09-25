package massmutation

import "github.com/dmikalova/vex/internal/card"

// New Frontiers
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a house. Reveal the top 3 cards of your deck. Archive each card of the chosen house and discard the others.
var NewFrontiers = set.New(
	"New Frontiers",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "326"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.RevealTopOfDeck{
				Amount: 3,
				Then: []card.TopAct{card.PartitionByChosenHouse{
					Matching: card.Into.Archives,
					Rest:     card.Into.Discard,
				}},
			},
		}),
)

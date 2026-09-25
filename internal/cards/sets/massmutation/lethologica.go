package massmutation

import "github.com/dmikalova/vex/internal/card"

// Lethologica
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Discard cards from the top of your deck until you discard a Logos card or run out of cards -> put the discarded card into your hand.
var Lethologica = set.New(
	"Lethologica",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "075"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.DiscardUntil{
				Player: card.Controller,
				House:  card.Houses.Named(card.House.Self),
			},
			Result: card.PutDiscardedIntoHand{},
		}),
)

package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Sound the Horns
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Discard cards from the top of your deck until you discard a Brobnar creature or run out of cards -> put the discarded creature into your hand.
var SoundTheHorns = set.New("Sound the Horns",
	card.House.Brobnar, card.Type.Tactic, card.Rarity.Uncommon,
	card.Provenance(card.CotA, "15"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.DiscardUntil{
				Player: card.Controller,
				Filter: card.Filter{Type: card.Type.Creature},
				House:  card.Houses.Named(card.House.Self),
			},
			Result: card.PutDiscardedIntoHand{Type: card.Type.Creature},
		}),
)

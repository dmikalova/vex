package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sound the Horns
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Discard cards from the top of your deck until you discard a Brobnar Creature or run out of cards -> put the discarded Creature into your hand.
var SoundTheHorns = set.New("Sound the Horns",
	card.House.Brobnar, card.Type.Tactic, card.Rarity.Uncommon,
	card.Provenance(card.CotA, "15"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.DiscardUntil{
				Type:  card.Type.Creature,
				House: card.House.Self,
			},
			Result: card.PutDiscardedIntoHand{Type: card.Type.Creature},
		}),
)

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Horizon Saber
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  11
//	Armor:  2
//	Traits: Robot
//
//	Play/Fight/Reap: Search your deck and discard pile for a card, reveal it, and put it into your archives, and shuffle your discard pile into your deck.
var HorizonSaber = set.Gigantic(
	"Horizon Saber",
	card.House.Logos,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "078"),
	card.WithPower(11),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.PlayFightReap, card.Sequence{
			Effects: []card.Effect{
				card.Search{
					Sources: []card.Zone{card.Deck, card.Discard},
					Reveal:  true,
					Dest:    card.To.Archives,
				},
				card.Shuffle{Zones: []card.Zone{card.Discard}},
			},
		}),
)

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Rockatiel
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast • Mutant
//
//	Elusive, Hazardous 1.
//	Play/Reap: Shuffle up to 2 creatures into their owners' decks.
var Rockatiel = set.New(
	"Rockatiel",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "429"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast, card.Traits.Mutant),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithHazardous(1),
	card.WithAbility(
		card.Trigger.PlayReap, card.PutChosen{
			Quantity:    card.UpTo{N: card.Fixed(2)},
			Target:      card.Target.EachCreature,
			Destination: card.To.DeckShuffled,
		}),
)

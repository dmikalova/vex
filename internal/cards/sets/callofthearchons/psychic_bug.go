package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Psychic Bug
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Æmber:  1
//	Traits: Cyborg • Insect
//
//	Play/Reap: Reveal your opponent's hand.
var PsychicBug = card.New(
	"Psychic Bug",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 149),
	card.WithPower(2),
	card.WithAemberBonus(1),
	card.WithTraits("Cyborg", "Insect"),
	card.WithPlayReap(card.Reveal{Player: card.Opponent}),
)

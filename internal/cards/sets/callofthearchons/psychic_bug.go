package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Psychic Bug
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Bonus:  Æmber
//	Traits: Cyborg • Insect
//
//	Play/Reap: Reveal your opponent's hand.
var PsychicBug = set.New(
	"Psychic Bug",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "149"),
	card.WithBonus(card.Bonus.Aember),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Insect),
	card.WithAbility(card.Trigger.PlayReap, card.RevealHand{Player: card.Opponent}),
)

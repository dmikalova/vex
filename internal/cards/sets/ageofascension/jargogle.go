package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Jargogle
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Beast • Mutant
//
//	Elusive.
//	Play: Put a card from your hand facedown under Jargogle.
//	Destroyed: If it is your turn, play the card under Jargogle. Otherwise, archive the card under Jargogle.
var Jargogle = set.New(
	"Jargogle",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "131"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast, card.Traits.Mutant),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.PutUnderFromHand{FaceDown: true}),
	card.WithAbility(
		card.Trigger.Destroyed, card.Conditional{
			Cond: card.ItIsYourTurn{},
			Then: card.PlayCardUnder{},
			Else: card.ArchiveCardUnder{},
		}),
)

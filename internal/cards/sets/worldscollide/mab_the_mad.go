package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mab the Mad
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Æmber:  1
//	Traits: Faerie
//
//	Reap: Shuffle Mab the Mad into its owner's deck.
var MabTheMad = card.New(
	"Mab the Mad",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 378),
	card.WithPower(2),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Faerie),
	card.WithAbility(
		card.Trigger.Reap, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.DeckShuffled,
		}),
)

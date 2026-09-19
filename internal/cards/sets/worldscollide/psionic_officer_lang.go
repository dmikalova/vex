package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Psionic Officer Lang
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	After an enemy creature reaps, archive the top card of your deck.
var PsionicOfficerLang = set.New(
	"Psionic Officer Lang",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "337"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Conditional{
			Cond: card.ItIsEnemy{},
			Then: card.ArchiveCard{
				Zone:      card.Deck,
				Selection: card.Top{},
			},
		}),
)

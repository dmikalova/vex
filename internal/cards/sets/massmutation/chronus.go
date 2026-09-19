package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Chronus
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	After you resolve a Draw bonus icon, you may archive a card from your hand.
//	Enhance Draw Draw.
var Chronus = set.New(
	"Chronus",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "084"),
	card.WithEnhance(card.Bonus.Draw, card.Bonus.Draw),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(card.Trigger.AfterBonusDraw, card.May{
		Do: card.ArchiveCard{
			Zone:      card.Hand,
			Selection: card.Chosen{},
		},
	}),
)

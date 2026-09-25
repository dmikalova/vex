package massmutation

import "github.com/dmikalova/vex/internal/card"

// Mutagenesis Researcher
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Scientist
//
//	Enhance Æmber Capture Damage Draw.
var MutagenesisResearcher = set.New(
	"Mutagenesis Researcher",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "325"),
	card.WithEnhance(card.Bonus.Aember, card.Bonus.Capture, card.Bonus.Damage, card.Bonus.Draw),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
)

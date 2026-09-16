package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Ensign El-Samra
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Mutant
//
//	Action: Reveal a card from your hand. Resolve that card's bonus icons.
//	Enhance Draw Draw Draw.
var EnsignElSamra = set.New(
	"Ensign El-Samra",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "340"),
	card.WithEnhance(card.Bonus.Draw, card.Bonus.Draw, card.Bonus.Draw),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.RevealChosenFromHand{},
			card.ResolveBonusIcons{Target: card.Target.Triggering},
		}}),
)

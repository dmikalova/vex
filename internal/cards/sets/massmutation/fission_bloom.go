package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Fission Bloom
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Power
//
//	Action: The next time you play a card this turn, resolve each of its bonus icons an additional time.
//	Enhance Draw.
var FissionBloom = set.New(
	"Fission Bloom",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "087"),
	card.WithEnhance(card.Bonus.Draw),
	card.WithTraits(card.Traits.Power),
	card.WithAbility(
		card.Trigger.Action, card.ExtraBonusIconResolution{}),
)

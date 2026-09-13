package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Creed of Nurture
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Creed of Nurture. Reveal a Creature from your hand and choose a Creature in play - for the remainder of the turn, the chosen Creature gains the text box of the revealed Creature.
var CreedOfNurture = card.New(
	"Creed of Nurture",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "386"),
	card.WithTraits(card.Traits.Power),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{
			Effects: []card.Effect{
				card.Destroy{Target: card.Target.This},
				card.LendTextBoxFromHand{},
			},
		}),
)

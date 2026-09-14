package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Sci. Officer Morpheus
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Shapeshifter • Scientist
//
//	After a Creature is played, if it is a friendly Creature, trigger the play effect of it.
var SciOfficerMorpheus = set.New(
	"Sci. Officer Morpheus",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "318"),
	card.WithPower(2),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Shapeshifter, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.AfterCreaturePlayed, card.Conditional{
			Cond: card.ItIsFriendly{},
			Then: card.TriggerAbility{
				Trigger: card.Trigger.Play,
				Target:  card.Target.Triggering,
			},
		}),
)

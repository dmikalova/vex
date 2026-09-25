package anomalyexpansion

import "github.com/dmikalova/vex/internal/card"

// Nizak, The Forgotten
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Dragon • Psion
//
//	While fighting, Nizak, The Forgotten gains invulnerable.
//	After a creature is destroyed in a fight with Nizak, The Forgotten, put it into its owner's hand.
var NizakTheForgotten = set.New(
	"Nizak, The Forgotten",
	card.House.Brobnar,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A05"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Dragon, card.Traits.Psion),
	card.WithConstant(card.ConstantAbility{
		Target:         card.Target.This,
		Keywords:       card.Keywords(card.Keyword.Invulnerable),
		WhileCondition: card.SourceIsFighting{},
	}),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.PutItIntoHand{}),
)

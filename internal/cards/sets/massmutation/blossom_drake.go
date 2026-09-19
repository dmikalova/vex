package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Blossom Drake
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dragon
//
//	Blossom Drake gains +1 power for each artifact in play.
//	Each artifact's text box is considered blank, except for traits.
var BlossomDrake = set.New(
	"Blossom Drake",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "395"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dragon),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This,
		PowerBonus: 1,
		Per:        card.ArtifactsInPlay{},
	}),
	card.WithConstant(card.ConstantAbility{
		Target:    card.Target.EachArtifact,
		BlankText: true,
	}),
)

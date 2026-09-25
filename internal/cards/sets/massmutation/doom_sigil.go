package massmutation

import "github.com/dmikalova/vex/internal/card"

// Doom Sigil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	If there are no creatures in play, destroy Doom Sigil.
//	Each creature gains poison.
var DoomSigil = set.New(
	"Doom Sigil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "277"),
	card.WithTraits(card.Traits.Power),
	card.WithConstant(card.ConstantAbility{
		Target:   card.Target.EachCreature,
		Keywords: card.Keywords(card.Keyword.Poison),
	}),
	card.WithDestroyedWhen(card.CardsInPlay{
		Player: card.EachPlayer,
		Type:   card.Type.Creature,
		None:   true,
	}),
)

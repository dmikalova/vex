package massmutation

import "github.com/dmikalova/vex/internal/card"

// The Archivist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Cyborg
//
//	Instead of picking up all of your archives, you may pick up any number of cards in your archives.
var TheArchivist = set.New(
	"The Archivist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "113"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Cyborg),
	card.WithConstant(card.ConstantAbility{
		SelectiveArchivePickup: true,
	}),
)

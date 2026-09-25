package massmutation

import "github.com/dmikalova/vex/internal/card"

// Safe House
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Action: Archive a friendly creature from play.
var SafeHouse = set.New(
	"Safe House",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "274"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.ArchiveFromPlay{Target: card.Target.FriendlyCreature}),
)

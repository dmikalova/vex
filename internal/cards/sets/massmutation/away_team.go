package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Away Team
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Alien • Human • Robot
//
//	Destroyed: Archive each upgrade on Away Team from play.
var AwayTeam = set.New(
	"Away Team",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "320"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Alien, card.Traits.Human, card.Traits.Robot),
	card.WithAbility(
		card.Trigger.Destroyed, card.ArchiveFromPlay{
			Target: card.Target.EachUpgradeOnThis,
		}),
)

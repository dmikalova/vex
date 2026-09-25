package massmutation

import "github.com/dmikalova/vex/internal/card"

// Wild Bounty
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: The next time you play a card this turn, resolve each of its bonus icons an additional time.
//	Enhance Æmber Æmber.
var WildBounty = set.New(
	"Wild Bounty",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "392"),
	card.WithEnhance(card.Bonus.Aember, card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ExtraBonusIconResolution{}),
)

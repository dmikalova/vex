package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ancient Power
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Ward each friendly creature with Æmber on it.
var AncientPower = set.New(
	"Ancient Power",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "198"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Ward{
			Target: card.Target.EachFriendlyCreature.WithAember(),
		}),
)

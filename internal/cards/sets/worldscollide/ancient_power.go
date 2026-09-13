package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ancient Power
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Ward each friendly Creature with Æmber on it.
var AncientPower = card.New(
	"Ancient Power",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "198"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Ward{
			Target: card.Target.EachFriendlyCreature.WithAember(),
		}),
)

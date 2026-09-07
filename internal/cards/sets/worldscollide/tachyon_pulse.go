package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tachyon Pulse
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each artifact, and exhaust each creature with an upgrade.
var TachyonPulse = card.New(
	"Tachyon Pulse",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "340"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.EachArtifact},
			card.Exhaust{Target: card.Target.EachCreature.WithUpgrade()},
		}}),
)

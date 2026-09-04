package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Strange Gizmo
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	After you forge a key, destroy each creature and each artifact.
var StrangeGizmo = card.New(
	"Strange Gizmo",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 134),
	card.WithTraits(card.Traits.Item),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.AfterForgeKey, card.Sequence{Effects: []card.Effect{
		card.Destroy{Target: card.Target.EachCreature},
		card.Destroy{Target: card.Target.EachArtifact},
	}}),
)

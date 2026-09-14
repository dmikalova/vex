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
//	After you forge a key, destroy each Creature and each Artifact.
var StrangeGizmo = set.New(
	"Strange Gizmo",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "134"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(card.Trigger.AfterForgeKey, card.Sequence{Effects: []card.Effect{
		card.Destroy{Target: card.Target.EachCreature},
		card.Destroy{Target: card.Target.EachArtifact},
	}}),
)

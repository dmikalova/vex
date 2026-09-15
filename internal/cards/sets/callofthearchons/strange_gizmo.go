package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Strange Gizmo
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	After you forge a key, destroy each creature and each artifact.
var StrangeGizmo = set.New(
	"Strange Gizmo",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "134"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(card.Trigger.AfterForgeKey, card.Sequence{Effects: []card.Effect{
		card.Destroy{Target: card.Target.EachCreature},
		card.Destroy{Target: card.Target.EachArtifact},
	}}),
)

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Customs Office
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Your opponent must give you 1 Æmber in order to play an artifact.
var CustomsOffice = card.New(
	"Customs Office",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 285),
	card.WithTraits(card.Traits.Location),
	card.WithRestrictions(card.Restrictions{Toll: card.Toll{
		Action: card.TollOn.PlayArtifact,
		Amount: 1,
	}}),
)

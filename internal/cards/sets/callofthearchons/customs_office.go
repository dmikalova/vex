package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Customs Office
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	In order to play an Artifact, your opponent must give you 1 Æmber.
var CustomsOffice = card.New(
	"Customs Office",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "285"),
	card.WithTraits(card.Traits.Location),
	card.WithRestrictions(card.Restrictions{Toll: card.Toll{
		Action: card.TollOn.PlayArtifact,
		Amount: 1,
	}}),
)

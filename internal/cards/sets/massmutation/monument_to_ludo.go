package massmutation

import (
	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/clusters"
)

// praefectusLudoName names Praefectus Ludo for the discard check; the card lives in
// the WorldsCollide package this set cannot import, so it is referenced by const.
const praefectusLudoName = "Praefectus Ludo"

// Monument to Ludo
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: If Praefectus Ludo is in your discard pile, move 2 Æmber from a creature to the common supply. Otherwise, move 1 Æmber from a creature to the common supply.
var MonumentToLudo = set.New(
	"Monument to Ludo",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.MM, "194"),
	card.LeadsCluster(clusters.Ludo),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.NamedCardInDiscard{Name: praefectusLudoName},
			Then: card.MoveAemberToSupply{
				Amount: 2,
				Target: card.Target.Creature,
			},
			Else: card.MoveAemberToSupply{
				Amount: 1,
				Target: card.Target.Creature,
			},
		}),
)

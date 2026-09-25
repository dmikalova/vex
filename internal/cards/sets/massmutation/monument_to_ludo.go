package massmutation

import (
	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/cards/clusters"
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
//	Action: Move 1 Æmber from a creature to the common supply. If Praefectus Ludo is in your discard pile, move 1 Æmber from the chosen creature to the common supply.
var MonumentToLudo = set.New(
	"Monument to Ludo",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.MM, "194"),
	card.LeadsCluster(clusters.Ludo),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.MoveAemberToSupply{
				Amount: 1,
				Target: card.Target.Creature,
				Bind:   true,
			},
			card.Conditional{
				Cond: card.NamedCardInDiscard{Name: praefectusLudoName},
				Then: card.MoveAemberToSupply{
					Amount: 1,
					Target: card.Target.TheChosenCreature,
				},
			},
		}}),
)

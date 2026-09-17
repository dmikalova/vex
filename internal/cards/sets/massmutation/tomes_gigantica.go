package massmutation

import (
	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/clusters"
)

// Tomes Gigantica
//
//	House:  None
//	Type:   Tactic
//	Rarity: Special
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature, reveal them, and put them into your hand. Shuffle your deck. Purge Tomes Gigantica.
var TomesGigantica = set.New(
	"Tomes Gigantica",
	card.House.None,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "004"),
	card.InCluster(clusters.Tutors),
	card.Houseless(),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.Search{
					Sources: []card.Zone{card.Deck, card.Discard},
					Filter:  card.Filter{Gigantic: true},
					Max:     2,
					Reveal:  true,
					Dest:    card.To.Hand,
				},
				card.Shuffle{},
				card.PurgeSource{},
			},
		}),
)

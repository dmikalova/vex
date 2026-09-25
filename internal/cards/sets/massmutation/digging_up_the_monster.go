package massmutation

import (
	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/cards/clusters"
)

// Digging Up the Monster
//
//	House:  None
//	Type:   Tactic
//	Rarity: Special
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature, reveal them, shuffle your deck, and put them into the top of your deck.
var DiggingUpTheMonster = set.New(
	"Digging Up the Monster",
	card.House.None,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "003"),
	card.InCluster(clusters.Tutors),
	card.Houseless(),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Search{
			Sources:              []card.Zone{card.Deck, card.Discard},
			Filter:               card.Filter{Gigantic: true},
			Max:                  2,
			Reveal:               true,
			ShuffleBeforePlacing: true,
			Dest:                 card.To.TopOfDeck,
		}),
)

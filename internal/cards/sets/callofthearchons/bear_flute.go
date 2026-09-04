package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bear Flute
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Fully heal an Ancient Bear. If there are no Ancient Bears in play, search your deck and discard pile and put each Ancient Bear from them into your hand -> shuffle your discard pile into your deck.
var BearFlute = card.New(
	"Bear Flute",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 340),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Heal{
				Fully:  true,
				Target: card.Target.Creature.Named("Ancient Bear"),
			},
			card.Conditional{
				Cond: card.InPlay{
					Player: card.EachPlayer,
					Type:   card.Type.Creature,
					Name:   "Ancient Bear",
					None:   true,
				},
				Then: card.Then{
					First:  card.SearchForName{Name: "Ancient Bear", All: true},
					Result: card.ShuffleIntoDeck{Zones: []card.Zone{card.Discard}},
				},
			},
		}}),
)

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Wretched Doll
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Play: Put a doom counter on a creature.
//	Action: Destroy each creature with a doom counter. Put a doom counter on a creature.
var WretchedDoll = set.New(
	"Wretched Doll",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "107"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Play, card.PlaceCounter{
			Kind:   card.Counter.Doom,
			Target: card.Target.Creature,
		}),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Destroy{
				Target: card.Target.EachCreature.WithCounter(card.Counter.Doom),
			},
			card.PlaceCounter{
				Kind:   card.Counter.Doom,
				Target: card.Target.Creature,
			},
		}}),
)

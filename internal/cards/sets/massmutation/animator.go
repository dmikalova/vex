package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Animator
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Give an artifact three +1 power counters. Move it to a flank of its controller's battleline as a creature with versatile for the remainder of the turn.
var Animator = set.New(
	"Animator",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "100"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.AddPowerCounter{
				Target: card.Target.Artifact,
				Amount: 3,
			},
			card.TurnIntoCreature{
				Target:    card.Target.TheChosenCreature,
				Duration:  card.Duration.RemainderOfPlayerTurn,
				Versatile: true,
			},
		}}),
)

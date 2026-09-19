package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Blast from the Past
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Exalt a friendly creature. Archive a Saurian creature from your discard pile. Deal damage equal to its power to an enemy creature.
var BlastFromThePast = set.New(
	"Blast from the Past",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "200"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Exalt{Target: card.Target.FriendlyCreature, Amount: 1},
			card.ArchiveCard{
				Zone: card.Discard,
				Selection: card.Chosen{
					Type:  card.Type.Creature,
					House: card.Houses.Named(card.House.Saurian),
				},
				Bind: true,
			},
			card.DealDamage{
				AmountFrom: card.PowerOfChosen{},
				Target:     card.Target.EnemyCreature,
			},
		}}),
)

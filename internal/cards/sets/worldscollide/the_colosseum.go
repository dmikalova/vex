package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Colosseum
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Versatile.
//	After an enemy Creature is destroyed while fighting, put a glory counter on The Colosseum.
//	Action: If there are 6 or more glory counters on The Colosseum, remove 6 glory counters from The Colosseum, and forge a key at current cost -> purge The Colosseum.
var TheColosseum = card.New(
	"The Colosseum",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "233"),
	card.WithTraits(card.Traits.Location),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.AfterEnemyDestroyedFighting, card.PlaceCounter{
			Kind:   card.Counter.Glory,
			Target: card.Target.This,
		}),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.CountersOnThisAtLeast{Kind: card.Counter.Glory, N: 6},
			Then: card.Sequence{Effects: []card.Effect{
				card.RemoveCounters{
					Kind:   card.Counter.Glory,
					Target: card.Target.This,
					Amount: 6,
				},
				card.ForgeKey{},
			}},
		}),
)

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Crassosaurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Dinosaur • Politician
//
//	Elusive.
//	Play: Crassosaurus captures 10 Æmber from any combination of players. If there are fewer than 10 Æmber on it, purge Crassosaurus.
var Crassosaurus = set.New(
	"Crassosaurus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "217"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.CaptureFromAnyPlayer{Amount: 10},
			card.Conditional{
				Cond: card.Not{Cond: card.CountIs{
					Count:  card.AemberOnThis{},
					Is:     card.AtLeast,
					Amount: 10,
				}},
				Then: card.PurgeCreature{Target: card.Target.This},
			},
		}},
	),
)

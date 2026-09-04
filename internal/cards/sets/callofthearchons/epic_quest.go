package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Epic Quest
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Quest
//
//	Versatile.
//	Play: Archive each friendly Knight trait creature from play.
//	Action: If you have played 7 or more Sanctum cards this turn, destroy Epic Quest, and forge a key at no cost.
var EpicQuest = card.New(
	"Epic Quest",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 231),
	card.WithTraits(card.Traits.Quest),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveFromPlay{
			Target: card.Target.EachFriendlyCreature.WithTrait(card.Traits.Knight),
		}),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.CardsPlayed{
				Player: card.Controller,
				House:  card.House.Self,
				Amount: 7,
			},
			Then: card.Sequence{Effects: []card.Effect{
				card.Destroy{Target: card.Target.This},
				card.ForgeKey{FreeOfCost: true},
			}},
		}),
)

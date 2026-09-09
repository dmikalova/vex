package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Big One
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	After a creature is played, put a fuse counter on The Big One. If there are 10 or more fuse counters on The Big One, destroy each creature and each artifact.
var TheBigOne = card.New(
	"The Big One",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "50"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.AfterCreaturePlayed, card.Sentences{Effects: []card.Effect{
			card.PlaceCounter{
				Kind:   card.Counter.Fuse,
				Target: card.Target.This,
			},
			card.Conditional{
				Cond: card.CountersOnThisAtLeast{
					Kind: card.Counter.Fuse,
					N:    10,
				},
				Then: card.Sequence{Effects: []card.Effect{
					card.Destroy{Target: card.Target.EachCreature},
					card.Destroy{Target: card.Target.EachArtifact},
				}},
			},
		}}),
)

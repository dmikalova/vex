package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Destroy Them All!
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy an artifact and a creature and an upgrade.
var DestroyThemAll = card.New(
	"Destroy Them All!",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 179),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.Artifact},
			card.Destroy{Target: card.Target.Creature},
			card.Destroy{Target: card.Target.Upgrade},
		}}),
)

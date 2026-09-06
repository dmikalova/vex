package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Cybergiant Rig
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "At the end of your turn, give this creature a -1 power counter."
//	Play: Fully heal this creature, and for each damage healed this way, give this creature a +1 power counter.
var CybergiantRig = card.New(
	"Cybergiant Rig",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 37),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.EndOfTurn,
			Effect: card.AddPowerCounter{
				Target: card.Target.This,
				Amount: -1,
			},
		}},
	}),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Heal{Fully: true, Target: card.Target.This},
			card.AddPowerCounter{
				Target: card.Target.This,
				Amount: 1,
				Per:    card.DamageHealed{},
			},
		}}),
)

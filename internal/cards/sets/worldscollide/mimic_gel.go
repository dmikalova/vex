package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mimic Gel
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Shapeshifter • Mutant
//
//	Play: Choose another creature - give Mimic Gel +1 power counters equal to its power, and Mimic Gel gains the text box of the chosen creature.
var MimicGel = set.New(
	"Mimic Gel",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "170"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Shapeshifter, card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature.Other(),
			Then: card.Sequence{
				Effects: []card.Effect{
					card.AddPowerCounter{
						Target: card.Target.This,
						Equal:  card.PowerOfChosen{},
					},
					card.GainTextBox{
						Target: card.Target.This,
						Source: card.Target.TheChosenCreature,
					},
				},
			},
		}),
)

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Chain Gang
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	After you play Subtle Chain, ready Chain Gang.
//	Action: Steal 1 Æmber. Shuffle a Subtle Chain from your discard pile into your deck.
var ChainGang = card.New(
	"Chain Gang",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "252"),
	card.Connects(card.Pull(SubtleChain, 1)),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.AfterCardPlayed, card.Conditional{
			Cond: card.ItIsNamed{Name: SubtleChain.Name},
			Then: card.Ready{Target: card.Target.This},
		}),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{
			Effects: []card.Effect{
				card.StealAember{Amount: 1},
				card.ShuffleNamedFromDiscardIntoDeck{Name: SubtleChain.Name},
			},
		}),
)

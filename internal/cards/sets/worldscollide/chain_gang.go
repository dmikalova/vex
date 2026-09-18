package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// chainGangCluster pulls a Subtle Chain or two into Chain Gang's pod — a Pull
// cluster, at least one averaging two (ADR 0036).
var chainGangCluster = card.Cluster{
	Name:     "Chain Gang",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Chain Gang
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	After you play Subtle Chain, ready Chain Gang.
//	Action: Steal 1 Æmber. Shuffle Subtle Chain from your discard pile into your deck.
var ChainGang = set.New(
	"Chain Gang",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "252"),
	card.LeadsCluster(chainGangCluster),
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
				card.ShuffleIntoDeck{
					Player: card.Controller, From: []card.Zone{card.Discard},
					Selection: card.Named{Name: SubtleChain.Name},
				},
			},
		}),
)

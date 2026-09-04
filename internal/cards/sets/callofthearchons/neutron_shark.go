package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Neutron Shark
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Mutant
//
//	Play/Fight/Reap: Destroy an enemy creature or artifact and a friendly creature or artifact, and discard the top card of your deck -> if the discarded card is not a Logos card, repeat this effect.
var NeutronShark = card.New(
	"Neutron Shark",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 146),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Mutant),
	card.WithPlayFightReap(card.RepeatOnCondition{
		Do: card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.EnemyCreatureOrArtifact},
			card.Destroy{Target: card.Target.FriendlyCreatureOrArtifact},
			card.DiscardTopOfDeck{Player: card.Controller},
		}},
		Cond: card.ItIs{
			House:   card.House.Self,
			Not:     true,
			Subject: card.Subject.DiscardedCard,
		},
	}),
)

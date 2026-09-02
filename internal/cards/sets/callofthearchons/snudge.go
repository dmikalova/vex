package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Snudge
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Demon
//
//	Fight/Reap: Put an artifact or flank creature into its owner's hand.
var Snudge = card.New(
	"Snudge",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 97),
	card.WithPower(4),
	card.WithTraits("Demon"),
	card.WithFightOrReap(card.PutFromPlay{
		Target:      card.Target.CreatureOrArtifact.OnFlank(),
		Destination: card.To.Hand,
	}),
)

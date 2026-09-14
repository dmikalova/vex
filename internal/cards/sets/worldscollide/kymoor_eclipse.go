package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Kymoor Eclipse
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Shuffle each flank Creature into its owner's deck.
var KymoorEclipse = set.New(
	"Kymoor Eclipse",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "243"),
	card.WithAbility(
		card.Trigger.Play, card.PutFromPlay{
			Target:      card.Target.EachCreature.OnFlank(),
			Destination: card.To.DeckShuffled,
		}),
)

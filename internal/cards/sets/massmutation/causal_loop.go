package massmutation

import "github.com/dmikalova/vex/internal/card"

// Causal Loop
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Archive 2 cards from your hand. Archive Causal Loop.
var CausalLoop = set.New(
	"Causal Loop",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "083"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.ArchiveCard{
					Zone:      card.Hand,
					Selection: card.Chosen{},
					Quantity:  card.Takes{N: card.Fixed(2)},
				},
				card.ArchiveSource{},
			},
		}),
)

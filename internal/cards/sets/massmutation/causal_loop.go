package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Causal Loop
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Archive a card from your hand. Archive Causal Loop.
var CausalLoop = set.New(
	"Causal Loop",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "083"),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.ArchiveCard{
					Zone:      card.Hand,
					Selection: card.Chosen{},
					Amount:    1,
				},
				card.ArchiveSource{},
			},
		}),
)

package massmutation

import "github.com/dmikalova/vex/internal/card"

// Look What I Found!
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Omega.
//	Play: Put a tactic, artifact, creature, and upgrade from your discard pile into your hand.
var LookWhatIFound = set.New(
	"Look What I Found!",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "402"),
	card.WithKeywords(card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PutCard{Zones: []card.Zone{card.Discard},
				Selection:   card.Chosen{Type: card.Type.Tactic},
				Destination: card.To.Hand,
			},
			card.PutCard{Zones: []card.Zone{card.Discard},
				Selection:   card.Chosen{Type: card.Type.Artifact},
				Destination: card.To.Hand,
			},
			card.PutCard{Zones: []card.Zone{card.Discard},
				Selection:   card.Chosen{Type: card.Type.Creature},
				Destination: card.To.Hand,
			},
			card.PutCard{Zones: []card.Zone{card.Discard},
				Selection:   card.Chosen{Type: card.Type.Upgrade},
				Destination: card.To.Hand,
			},
		}}),
)

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Z-Particle Tracker
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Connected
//	Bonus:  Æmber
//
//	This creature gains, "Fight: Search your deck for an upgrade, reveal it, and put it into your hand. Shuffle your deck."
var ZParticleTracker = set.New(
	"Z-Particle Tracker",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Connected,
	card.Provenance(card.MM, "354"),
	card.InCluster(zForceCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Fight,
			Effect: card.Sentences{
				Effects: []card.Effect{
					card.Search{
						Sources: []card.Zone{card.Deck},
						Filter:  card.Filter{Type: card.Type.Upgrade},
						Reveal:  true,
						Dest:    card.To.Hand,
					},
					card.Shuffle{},
				},
			},
		}},
	}),
)

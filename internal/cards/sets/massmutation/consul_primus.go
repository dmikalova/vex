package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Consul Primus
var ConsulPrimus = set.New(
	"Consul Primus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "187"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithEnhance(card.Bonus.Capture),
	card.WithAbility(
		card.Trigger.Reap, card.MoveAember{
			Amount: 1,
			From:   card.Target.Creature,
			Onto:   card.Target.OtherCreature,
		}),
)

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Titanic Bumblebird
//
//	House:  Untamed
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  8
//	Traits: Beast • Insect
//
//	Play/Reap: Destroy an enemy creature -> give a friendly creature +1 power counters equal to power of creatures destroyed this way.
var TitanicBumblebird = set.Gigantic(
	"Titanic Bumblebird",
	card.House.Untamed,
	card.Rarity.Rare,
	card.Provenance(card.MoMu, "395"),
	card.WithPower(8),
	card.WithTraits(card.Traits.Beast, card.Traits.Insect),
	card.WithAbility(
		card.Trigger.PlayReap, card.Then{
			First: card.Destroy{Target: card.Target.EnemyCreature},
			Result: card.AddPowerCounter{
				Target: card.Target.FriendlyCreature,
				Equal:  card.PowerDestroyedThisWay{},
			},
		}),
)

package massmutation

import "github.com/dmikalova/vex/internal/card"

// Ultra Gravitron
//
//	House:  Logos
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  10
//	Armor:  3
//	Traits: Robot
//
//	Play: Archive the top 5 cards of your deck.
//	Fight/Reap: Discard a card from your archives -> purge a creature. Resolve that card's bonus icons.
var UltraGravitron = set.Gigantic(
	"Ultra Gravitron",
	card.House.Logos,
	card.Rarity.Rare,
	card.Provenance(card.MM, "125"),
	card.WithPower(10),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveCard{
			Zone:      card.Deck,
			Selection: card.Top{},
			Quantity:  card.Takes{N: card.Fixed(5)},
		}),
	card.WithAbility(
		card.Trigger.FightReap, card.Then{
			First: card.DiscardCard{
				Player:    card.Controller,
				Zones:     []card.Zone{card.Archives},
				Selection: card.Chosen{},
			},
			Result: card.Sequence{Effects: []card.Effect{
				card.PurgeCreature{Target: card.Target.Creature},
				card.ResolveBonusIcons{Target: card.Target.Triggering},
			}},
		}),
)

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

const angryMobName = "Angry Mob"

// Angry Mob
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Before Fight: You may discard cards from the top of your deck until you discard an Angry Mob or run out of cards -> put the discarded creature into your hand.
var AngryMob = set.New(
	angryMobName,
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "143"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(
		card.Trigger.BeforeFight, card.May{Do: card.Then{
			First:  card.DiscardUntil{Player: card.Controller, Name: angryMobName},
			Result: card.PutDiscardedIntoHand{Type: card.Type.Creature},
		}}),
)

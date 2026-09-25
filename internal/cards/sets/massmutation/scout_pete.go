package massmutation

import "github.com/dmikalova/vex/internal/card"

// Scout Pete
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Alien
//
//	Play/Fight/Reap: Look at the top card of your deck and you may discard that card.
var ScoutPete = set.New(
	"Scout Pete",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "311"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Alien),
	card.WithAbility(card.Trigger.PlayFightReap, card.LookAtTopOfDeck{
		Amount: 1,
		Then:   []card.TopAct{card.MayDiscardLookedAt{}},
	}),
)

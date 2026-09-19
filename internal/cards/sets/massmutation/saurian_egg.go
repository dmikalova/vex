package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Saurian Egg
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Armor:  5
//	Bonus:  Æmber
//	Traits: Dinosaur • Egg
//
//	Versatile.
//	Saurian Egg cannot fight.
//	Saurian Egg cannot reap.
//	Action: Discard the top 2 cards of your deck. For each Saurian creature discarded this way, put it into play ready. Give it three +1 power counters. If you discard a Saurian creature this way, destroy Saurian Egg.
var SaurianEgg = set.New(
	"Saurian Egg",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "210"),
	card.WithBonus(card.Bonus.Aember),
	card.WithPower(1),
	card.WithArmor(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Egg),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithCannotBeUsedTo(card.UseKind.Fight, card.UseKind.Reap),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.DiscardTop{Player: card.Controller, Amount: 2},
			card.ForEachDiscarded{
				House: card.Houses.Named(card.House.Self),
				Type:  card.Type.Creature,
				Do: card.Sequence{Effects: []card.Effect{
					card.PutIntoPlay{Target: card.Target.Triggering, Ready: true},
					card.AddPowerCounter{Target: card.Target.Triggering, Amount: 3},
				}},
			},
			card.Conditional{
				Cond: card.DiscardedThisWay{
					House: card.Houses.Named(card.House.Self),
					Type:  card.Type.Creature,
				},
				Then: card.Destroy{Target: card.Target.This},
			},
		}}),
)

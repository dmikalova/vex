package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tireless Crocag
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Giant
//
//	Versatile.
//	Tireless Crocag cannot reap.
//	If there are no enemy creatures in play, destroy Tireless Crocag.
var TirelessCrocag = card.New(
	"Tireless Crocag",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 47),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithCannotBeUsedTo(card.UseKind.Reap),
	card.WithDestroyedWhen(card.InPlay{
		Player: card.Opponent,
		Type:   card.Type.Creature,
		None:   true,
	}),
)

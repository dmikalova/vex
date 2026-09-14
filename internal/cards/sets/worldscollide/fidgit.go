package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Fidgit
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Faerie • Thief
//
//	Elusive.
//	Reap: Discard a random card from your opponent's archives or the top card of their deck. If that card is a Tactic, play it as if it were yours.
var Fidgit = set.New(
	"Fidgit",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "254"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Faerie, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.Sentences{Effects: []card.Effect{
			card.DiscardOpponentArchivesOrDeckTop{},
			card.Conditional{
				Cond: card.ItIs{Type: card.Type.Tactic, Subject: card.Subject.ThatCard},
				Then: card.PlayItFromOpponentDiscard{},
			},
		}},
	),
)

package callofthearchons_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	cota "github.com/dmikalova/vactrol/internal/cards/sets/callofthearchons"
)

// Sniffer
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Ally
//
//	Action: For the remainder of the turn, each creature loses elusive.
func TestSniffer(t *testing.T) {
	var sniffer, hunter, hider ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Mars,
			InPlay: ct.Cards(
				ct.Bind(&sniffer, cota.Sniffer),
				ct.Bind(&hunter, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&hider, ct.Creature(ct.Power(1), ct.Keywords(card.Keyword.Elusive))),
			),
		},
	})

	h.P1.UseAction(sniffer)
	// Elusive would otherwise soak this first attack, leaving the hider alive.
	h.P1.Fight(hunter, hider)
	h.Expect(hider).At(ct.Discard)
}

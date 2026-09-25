package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Cybergiant Rig
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "At the end of your turn, give this creature a -1 power counter."
//	Play: Fully heal this creature. For each damage healed this way, give this creature a +1 power counter.
func TestCybergiantRig(t *testing.T) {
	var host ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			Hand:   ct.Cards(CybergiantRig),
			InPlay: ct.Cards(ct.Bind(&host, ct.Creature(ct.Power(4)))),
		},
	})
	host.Damaged(3)

	h.P1.Play(CybergiantRig)

	// Fully healed, and a +1 power counter per damage healed (3).
	h.Expect(host).Damage(0)
	h.Expect(host).Power(7)

	// At the end of the turn the host sheds one of those counters.
	h.P1.EndTurn()
	h.Expect(host).Power(6)
}

package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Might Makes Right
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: You may destroy any number of friendly creatures with total power of 25 or more - forge a key at no cost.
func TestMightMakesRight(t *testing.T) {
	t.Run("sacrificing 25 total power forges a key for free", func(t *testing.T) {
		var big, med ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(MightMakesRight),
				InPlay: ct.Cards(
					ct.Bind(&big, ct.Creature(ct.Power(13))),
					ct.Bind(&med, ct.Creature(ct.Power(12))),
				),
			},
		})

		h.P1.Play(MightMakesRight)
		h.P1.ClickCard(big)
		h.P1.ClickCard(med)

		h.Expect(big).At(ct.Discard)
		h.Expect(med).At(ct.Discard)
		h.P1.ExpectKeys(1)
	})
}

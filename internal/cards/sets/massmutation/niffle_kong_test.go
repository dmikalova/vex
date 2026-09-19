package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Niffle Kong
//
//	House:  Untamed
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  12
//	Armor:  2
//	Traits: Mutant • Niffle
//
//	Play: Search your deck and discard pile for any number of Niffle creatures, reveal them, and put them into your hand. Shuffle your deck.
//	Fight/Reap: You may destroy a friendly Niffle creature -> deal 3 damage to a creature. Steal 1 Æmber. Destroy an enemy artifact.
func TestNiffleKong(t *testing.T) {
	t.Run("play tutors every Niffle creature to hand", func(t *testing.T) {
		var kong, niffleInDeck, niffleInDiscard, plainInDeck ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&kong, NiffleKong), card.GiganticArt(NiffleKong)),
				Deck: ct.Cards(
					ct.Bind(&niffleInDeck, ct.Creature(ct.Traits(card.Traits.Niffle))),
					ct.Bind(&plainInDeck, ct.Creature(ct.Traits(card.Traits.Beast))),
				),
				Discard: ct.Cards(
					ct.Bind(&niffleInDiscard, ct.Creature(ct.Traits(card.Traits.Niffle))),
				),
			},
		})

		h.P1.Play(kong)

		h.Expect(niffleInDeck).At(ct.Hand)    // every Niffle creature is tutored
		h.Expect(niffleInDiscard).At(ct.Hand) // from both the deck and discard pile
		h.Expect(plainInDeck).At(ct.Deck)     // the non-Niffle card stays behind
	})

	t.Run("reap sacrifices a Niffle to strike, steal, and break an artifact",
		func(t *testing.T) {
			var kong, niffleAlly, enemy, artifact ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Untamed,
					InPlay: ct.Cards(
						ct.Bind(&kong, NiffleKong),
						ct.Bind(&niffleAlly, ct.Creature(
							ct.Power(2), ct.Traits(card.Traits.Niffle))),
					),
				},
				P2: ct.Side{
					Amber: 2,
					InPlay: ct.Cards(
						ct.Bind(&enemy, ct.Creature(ct.Power(5))),
						ct.Bind(&artifact, ct.Artifact()),
					),
				},
			})

			h.P1.Reap(kong)
			h.P1.ClickCard(niffleAlly) // accept the may and choose the sacrifice
			h.P1.ClickCard(enemy)      // deal 3 damage to the enemy creature

			h.Expect(niffleAlly).At(ct.Discard)
			h.Expect(enemy).Damage(3)
			h.Expect(artifact).At(ct.Discard)
			h.P1.ExpectAmber(2) // 1 from reaping, 1 stolen
			h.P2.ExpectAmber(1)
		})

	t.Run("declining the sacrifice does nothing", func(t *testing.T) {
		var kong, niffleAlly, enemy, artifact ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&kong, NiffleKong),
					ct.Bind(&niffleAlly, ct.Creature(
						ct.Power(2), ct.Traits(card.Traits.Niffle))),
				),
			},
			P2: ct.Side{
				Amber: 2,
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.Power(5))),
					ct.Bind(&artifact, ct.Artifact()),
				),
			},
		})

		h.P1.Reap(kong)
		h.P1.ClickDone() // decline the optional sacrifice

		h.Expect(niffleAlly).At(ct.PlayArea)
		h.Expect(enemy).Damage(0)
		h.Expect(artifact).At(ct.PlayArea)
		h.P1.ExpectAmber(1) // only the reap bonus; nothing was stolen
		h.P2.ExpectAmber(2)
	})
}

// TestIsNiffleCreature checks the deck-generation predicate Niffle Kong uses to
// guarantee its deck a pool of Niffle creatures to fetch (card.PullsMatching),
// the cross-set family Troop Call also pulls (Niffle Ape, Niffle Queen).
func TestIsNiffleCreature(t *testing.T) {
	cases := []struct {
		name string
		def  card.Definition
		want bool
	}{
		{
			"niffle creature",
			card.Definition{
				Type:   card.Type.Creature,
				Traits: []card.Trait{card.Traits.Niffle},
			},
			true,
		},
		{
			"non-niffle creature",
			card.Definition{
				Type:   card.Type.Creature,
				Traits: []card.Trait{card.Traits.Beast},
			},
			false,
		},
		{
			"niffle artifact",
			card.Definition{
				Type:   card.Type.Artifact,
				Traits: []card.Trait{card.Traits.Niffle},
			},
			false,
		},
	}
	for i := range cases {
		tc := &cases[i]
		if got := isNiffleCreature(tc.def); got != tc.want {
			t.Errorf("isNiffleCreature(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

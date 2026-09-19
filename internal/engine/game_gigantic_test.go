package engine

import (
	"errors"
	"testing"
)

// TestGiganticPartnerLinkIsSymmetric links two halves, checks each finds the
// other, that giganticHalves returns both, and that unlinking clears both ends.
func TestGiganticPartnerLinkIsSymmetric(t *testing.T) {
	g := started(t)
	base := g.Register(NewCard("Deusillus", Sanctum, Creature, Rare, WithPower(15)), 0)
	art := g.Register(NewCard("Deusillus", Sanctum, Creature, Rare), 0)

	g.linkGiganticPartners(base, art)

	if p, ok := g.giganticPartner(base); !ok || p != art {
		t.Fatalf("base partner = %d,%v, want %d,true", p, ok, art)
	}
	if p, ok := g.giganticPartner(art); !ok || p != base {
		t.Fatalf("art partner = %d,%v, want %d,true", p, ok, base)
	}

	if got := g.giganticHalves(base); len(got) != 2 || got[0] != base || got[1] != art {
		t.Fatalf("giganticHalves(base) = %v, want [%d %d]", got, base, art)
	}

	g.unlinkGigantic(base)

	if _, ok := g.giganticPartner(base); ok {
		t.Error("base should have no partner after unlink")
	}
	if _, ok := g.giganticPartner(art); ok {
		t.Error("art should have no partner after unlink")
	}
}

// TestGiganticHalvesLoneCard confirms an unlinked card reports only itself and
// that unlinking a lone card is a no-op.
func TestGiganticHalvesLoneCard(t *testing.T) {
	g := started(t)
	lone := g.Register(NewCard("lone", Brobnar, Creature, Common, WithPower(1)), 0)

	if got := g.giganticHalves(lone); len(got) != 1 || got[0] != lone {
		t.Fatalf("giganticHalves(lone) = %v, want [%d]", got, lone)
	}
	g.unlinkGigantic(lone)
	if _, ok := g.giganticPartner(lone); ok {
		t.Error("lone card should have no partner")
	}
}

// TestHasBonusIconsSeesLinkedArtHalf confirms a gigantic base half prints no
// icons alone but exposes its linked art half's icons while in play (ADR 0042).
func TestHasBonusIconsSeesLinkedArtHalf(t *testing.T) {
	g := started(t)
	base := g.Register(NewCard("Deusillus", Sanctum, Creature, Rare, WithPower(15)), 0)
	art := g.Register(NewCard("Deusillus", Sanctum, Creature, Rare, WithBonus(BonusAember)), 0)

	if g.HasBonusIcons(base) {
		t.Fatal("base half alone should print no bonus icons")
	}
	if !g.HasBonusIcons(art) {
		t.Fatal("art half prints its own bonus icons")
	}

	g.linkGiganticPartners(base, art)
	if !g.HasBonusIcons(base) {
		t.Fatal("linked base should expose the art half's bonus icons")
	}
}

// playedGigantic builds the in-play shape of a gigantic for player p: the base
// half holds a battleline slot and the art half is in play slot-less, linked. It
// stands in for the enter funnel so leave-play paths can be tested on their own.
func playedGigantic(g *Game, p int) (base, art LocalID) {
	base = g.AddToBattleline(giganticBase("Deusillus", 15), p)
	art = g.Register(giganticArt("Deusillus"), p)
	g.linkGiganticPartners(base, art)
	return base, art
}

// TestGiganticLeavesPlayTogether checks each relocation out of play sends both
// halves to the same zone as two separate cards, leaving the state sound.
func TestGiganticLeavesPlayTogether(t *testing.T) {
	cases := []struct {
		name string
		move func(g *Game, base LocalID)
		has  func(g *Game, id LocalID) bool
	}{
		{
			"discard",
			(*Game).discardDestroyed,
			func(g *Game, id LocalID) bool { return g.State.Discard[0].contains(id) },
		},
		{
			"purge",
			(*Game).purgeFromPlay,
			func(g *Game, id LocalID) bool { return g.State.Purge[0].contains(id) },
		},
		{
			"hand",
			(*Game).putIntoHand,
			func(g *Game, id LocalID) bool { return g.State.Hand[0].contains(id) },
		},
		{
			"topOfDeck",
			(*Game).putOnTopOfDeck,
			func(g *Game, id LocalID) bool { return g.State.Deck[0].contains(id) },
		},
		{
			"archives",
			(*Game).putIntoArchives,
			func(g *Game, id LocalID) bool { return g.State.Archives[0].contains(id) },
		},
		{
			"shuffled",
			(*Game).putIntoDeckShuffled,
			func(g *Game, id LocalID) bool { return g.State.Deck[0].contains(id) },
		},
		{
			// Abduction files a card in the abductor's archives, not its owner's, so
			// both halves must land on the same side or the gigantic is torn in two.
			"abducted",
			func(g *Game, base LocalID) { g.PutIntoYourArchives(base, 1) },
			func(g *Game, id LocalID) bool { return g.State.Archives[1].contains(id) },
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := started(t)
			base, art := playedGigantic(g, 0)
			tc.move(g, base)
			if !tc.has(g, base) || !tc.has(g, art) {
				t.Fatalf("%s: want both halves in zone, got base=%v art=%v",
					tc.name, tc.has(g, base), tc.has(g, art))
			}
			if _, ok := g.giganticPartner(base); ok {
				t.Error("halves should be unlinked out of play")
			}
			if err := g.InvariantError(); err != nil {
				t.Fatalf("state unsound after %s: %v", tc.name, err)
			}
		})
	}
}

// TestGiganticSwapTakesBothHalves swaps an in-play gigantic for a card in the
// discard: both halves go to the discard and the resting card enters play.
func TestGiganticSwapTakesBothHalves(t *testing.T) {
	g := started(t)
	base, art := playedGigantic(g, 0)
	resting := g.AddToDiscard(testCreature("Sub", 3), 0)

	g.swapAcrossZones(base, resting)

	if !g.State.Discard[0].contains(base) || !g.State.Discard[0].contains(art) {
		t.Fatal("both gigantic halves should be in the discard after the swap")
	}
	if !g.inPlay(resting) {
		t.Error("the resting card should have entered play")
	}
	if err := g.InvariantError(); err != nil {
		t.Fatalf("state unsound after swap: %v", err)
	}
}

// TestManualMoveTearsDownBothGiganticHalves confirms the manual removal path tears
// down a gigantic's art half so it never dangles.
func TestManualMoveTearsDownBothGiganticHalves(t *testing.T) {
	g := started(t)
	base, art := playedGigantic(g, 0)

	g.ManualMove(base, ManualDiscard)

	if g.inPlay(base) || g.inPlay(art) {
		t.Error("neither half should be in play after removal")
	}
	if _, ok := g.giganticPartner(art); ok {
		t.Error("the art half should be unlinked after removal")
	}
}

// brobnarGiganticHalves builds a playable gigantic in the active house (Brobnar):
// a base half carrying the power, an art half carrying one Æmber bonus icon.
func brobnarGiganticHalves(power int) (base, art CardDefinition) {
	base = NewCard("Colossus", Brobnar, Creature, Rare, WithPower(power))
	base.GiganticRole = GiganticBase
	art = NewCard("Colossus", Brobnar, Creature, Rare, WithBonus(BonusAember))
	art.GiganticRole = GiganticArt
	return base, art
}

// TestPlayGiganticFromHand plays a gigantic with both halves in hand: the base
// half takes the battleline slot, the art half is linked slot-less, both leave
// the hand, and the art half's bonus icon resolves once.
func TestPlayGiganticFromHand(t *testing.T) {
	g := started(t)
	baseDef, artDef := brobnarGiganticHalves(15)
	base := g.AddToHand(baseDef, 0)
	art := g.AddToHand(artDef, 0)
	before := g.Aember(0)

	got, err := g.PlayCreature(0, handIdxByID(g, 0, base), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got != base {
		t.Fatalf("played id = %d, want base %d", got, base)
	}
	if !g.inPlay(base) {
		t.Error("the base half should be on the battleline")
	}
	if p, ok := g.giganticPartner(base); !ok || p != art {
		t.Errorf("base should be linked to art %d, got %d,%v", art, p, ok)
	}
	if g.State.Hand[0].contains(base) || g.State.Hand[0].contains(art) {
		t.Error("both halves should have left the hand")
	}
	if g.Aember(0) != before+1 {
		t.Errorf("art half's Æmber bonus should resolve once: Æmber %d, want %d",
			g.Aember(0), before+1)
	}
	if err := g.InvariantError(); err != nil {
		t.Fatalf("state unsound: %v", err)
	}
}

// TestPlayGiganticFromArtHalf plays the art half's hand card: the base half still
// becomes the battleline representative.
func TestPlayGiganticFromArtHalf(t *testing.T) {
	g := started(t)
	baseDef, artDef := brobnarGiganticHalves(15)
	base := g.AddToHand(baseDef, 0)
	art := g.AddToHand(artDef, 0)

	got, err := g.PlayCreature(0, handIdxByID(g, 0, art), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got != base {
		t.Fatalf("played id = %d, want base %d even when the art half was clicked",
			got, base)
	}
	if !g.inPlay(base) {
		t.Error("the base half holds the battleline slot")
	}
}

// TestPlayGiganticLoneHalfCannotBePlayed confirms a gigantic half with no partner
// in hand cannot be played: the play fails and the half stays in hand.
func TestPlayGiganticLoneHalfCannotBePlayed(t *testing.T) {
	g := started(t)
	baseDef, _ := brobnarGiganticHalves(15)
	base := g.AddToHand(baseDef, 0)

	_, err := g.PlayCreature(0, handIdxByID(g, 0, base), false)
	if !errors.Is(err, ErrGiganticNoPartner) {
		t.Fatalf("PlayCreature lone half = %v, want ErrGiganticNoPartner", err)
	}
	if !g.State.Hand[0].contains(base) {
		t.Error("the lone half should stay in hand after a failed play")
	}
	if err := g.CanPlay(0, base); !errors.Is(err, ErrGiganticNoPartner) {
		t.Errorf("CanPlay lone half = %v, want ErrGiganticNoPartner", err)
	}
}

// TestPlayGiganticBarredByCreatureRestriction confirms a gigantic is a creature,
// so it cannot be played while creatures are barred even through an effect-play.
func TestPlayGiganticBarredByCreatureRestriction(t *testing.T) {
	g := started(t)
	g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common, WithPower(1),
		WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	baseDef, artDef := brobnarGiganticHalves(15)
	base := g.AddToHand(baseDef, 0)
	g.AddToHand(artDef, 0)

	_, err := g.playCardFromZone(0, base,
		func() { g.State.Hand[0].remove(base) }, playCardOptions{})
	if !errors.Is(err, ErrCannotPlayCreature) {
		t.Fatalf("gigantic while creatures barred = %v, want ErrCannotPlayCreature", err)
	}
	if !g.State.Hand[0].contains(base) {
		t.Error("a barred gigantic should stay in hand")
	}
}

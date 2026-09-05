package engine

import "testing"

// TestUnderChainStitchesWhenMiddleLeaves attaches several cards under one host
// and detaches the one in the middle, exercising the tail-append in AttachUnder
// and the predecessor-stitch branch in detachUnder — mirroring
// TestUpgradeChainStitchesWhenMiddleLeaves for the Under mechanic (see ADR 0016).
func TestUnderChainStitchesWhenMiddleLeaves(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("host", Brobnar, Artifact, Common), 0)
	u1 := g.Register(NewCard("u1", Brobnar, Creature, Common, WithPower(1)), 0)
	u2 := g.Register(NewCard("u2", Brobnar, Creature, Common, WithPower(2)), 0)
	u3 := g.Register(NewCard("u3", Brobnar, Creature, Common, WithPower(3)), 0)
	g.AttachUnder(host, u1, false)
	g.AttachUnder(host, u2, false)
	g.AttachUnder(host, u3, false)

	if got := g.Under(host); len(got) != 3 || got[0] != u1 || got[1] != u2 || got[2] != u3 {
		t.Fatalf("under = %v, want [%d %d %d] in attach order", got, u1, u2, u3)
	}

	if _, ok := g.detachUnder(u2); !ok {
		t.Fatal("detachUnder should report the host it left")
	}

	if got := g.Under(host); len(got) != 2 || got[0] != u1 || got[1] != u3 {
		t.Fatalf("under after middle detached = %v, want [%d %d]", got, u1, u3)
	}
	if _, ok := g.underHostOf(u2); ok {
		t.Error("detached card should no longer report a host")
	}

	// Detaching the tail of a longer chain walks past a non-head predecessor to
	// find and re-stitch the link, covering the predecessor walk in detachUnder.
	other := g.AddArtifact(NewCard("other", Brobnar, Artifact, Common), 0)
	a := g.Register(NewCard("a", Brobnar, Creature, Common, WithPower(1)), 0)
	b := g.Register(NewCard("b", Brobnar, Creature, Common, WithPower(1)), 0)
	c := g.Register(NewCard("c", Brobnar, Creature, Common, WithPower(1)), 0)
	g.AttachUnder(other, a, false)
	g.AttachUnder(other, b, false)
	g.AttachUnder(other, c, false)

	g.detachUnder(c)

	if got := g.Under(other); len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatalf("under after tail detached = %v, want [%d %d]", got, a, b)
	}
}

// TestUnderHostOfMissing confirms a card that was never placed under anything
// reports no host.
func TestUnderHostOfMissing(t *testing.T) {
	g := started(t)
	loose := g.AddToHand(NewCard("loose", Brobnar, Creature, Common, WithPower(1)), 0)

	if _, ok := g.underHostOf(loose); ok {
		t.Error("a card never placed under anything should report no host")
	}
	if got := g.Under(loose); len(got) != 0 {
		t.Errorf("Under(loose) = %v, want none", got)
	}
	if _, ok := g.detachUnder(loose); ok {
		t.Error("detaching a card that was never placed under anything should report no host")
	}
}

// TestPutCardUnderMovesFromHand puts a card from a hand under a host and checks
// it leaves the hand, reports the host, carries its facedown flag, and logs the
// move.
func TestPutCardUnderMovesFromHand(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("host", Brobnar, Artifact, Common), 0)
	buried := g.AddToHand(NewCard("buried", Brobnar, Creature, Common, WithPower(2)), 0)

	g.PutCardUnder(0, buried, host, true)

	if g.State.Hand[0].contains(buried) {
		t.Error("the buried card should have left the hand")
	}
	if got := g.Under(host); len(got) != 1 || got[0] != buried {
		t.Errorf("under = %v, want [%d]", got, buried)
	}
	if !g.UnderFaceDown(buried) {
		t.Error("the buried card should be facedown")
	}
	last := g.Log[len(g.Log)-1].Entry
	put, ok := last.(CardPutUnder)
	if !ok {
		t.Fatalf("last log entry = %T, want CardPutUnder", last)
	}
	if put.Card != buried || put.Host != host || !put.FaceDown {
		t.Errorf("log entry = %+v, want Card=%d Host=%d FaceDown=true", put, buried, host)
	}
}

// TestPutCardUnderIgnoresACardElsewhere mirrors
// TestGamePlayFromIgnoresACardElsewhere: a card not in the given hand does
// nothing.
func TestPutCardUnderIgnoresACardElsewhere(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("host", Brobnar, Artifact, Common), 0)
	deckCard := g.AddToDeck(NewCard("Not In Hand", Brobnar, Creature, Common, WithPower(2)), 0)

	g.PutCardUnder(0, deckCard, host, false)

	if got := g.Under(host); len(got) != 0 {
		t.Errorf("under = %v, want none", got)
	}
}

// TestPlayFromUnderPlaysAndDetaches plays a card placed under a host and checks
// it lands in play and no longer reports a host.
func TestPlayFromUnderPlaysAndDetaches(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("host", Brobnar, Artifact, Common), 0)
	buried := g.Register(NewCard("buried", Brobnar, Creature, Common, WithPower(2)), 0)
	g.AttachUnder(host, buried, true)

	g.PlayFromUnder(0, buried)

	if got := g.Battleline(0); len(got) != 1 || got[0] != buried {
		t.Errorf("battleline = %v, want the played creature %d", got, buried)
	}
	if _, ok := g.underHostOf(buried); ok {
		t.Error("a played under-card should no longer report a host")
	}
}

// TestPlayFromUnderIgnoresACardNotUnderAnything confirms a card not placed
// under anything is left alone.
func TestPlayFromUnderIgnoresACardNotUnderAnything(t *testing.T) {
	g := started(t)
	loose := g.Register(NewCard("loose", Brobnar, Creature, Common, WithPower(2)), 0)

	g.PlayFromUnder(0, loose)

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none", got)
	}
}

// TestDiscardUnderSendsBuriedCardsToTheirOwnersDiscard confirms that when a host
// leaves play, every card placed under it goes to its own owner's discard pile
// (Graft's explicit rule, generalized here — see ADR 0016) and no longer reports
// a host.
func TestDiscardUnderSendsBuriedCardsToTheirOwnersDiscard(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("host", Brobnar, Artifact, Common), 0)
	mine := g.Register(NewCard("mine", Brobnar, Creature, Common, WithPower(1)), 0)
	theirs := g.Register(NewCard("theirs", Logos, Creature, Common, WithPower(1)), 1)
	g.AttachUnder(host, mine, true)
	g.AttachUnder(host, theirs, false)

	g.discardUnder(host)

	if d := g.Discard(0); len(d) != 1 || d[0] != mine {
		t.Errorf("P0 discard = %v, want only %d", d, mine)
	}
	if d := g.Discard(1); len(d) != 1 || d[0] != theirs {
		t.Errorf("P1 discard = %v, want only %d", d, theirs)
	}
	if _, ok := g.underHostOf(mine); ok {
		t.Error("a discarded under-card should no longer report a host")
	}
}

// TestPeekable confirms only a host's controller may look at its under-cards.
func TestPeekable(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("host", Brobnar, Artifact, Common), 0)

	if !g.Peekable(0, host) {
		t.Error("the controller should be able to peek")
	}
	if g.Peekable(1, host) {
		t.Error("the opponent should not be able to peek")
	}
}

package engine

import "testing"

func TestDiscardArchivesEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0 // player 0 discards player 1's (hidden) archives
	g.State.Archives[1].add(g.Register(testCreature("a", 1), 1))
	g.State.Archives[1].add(g.Register(testCreature("b", 1), 1))
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DiscardArchives{Player: Opponent}
	if e.Text() != "your opponent discards each of their archived cards" {
		t.Errorf("text = %q", e.Text())
	}
	if (DiscardArchives{Player: Controller}).Text() != "discard each of your archived cards" {
		t.Errorf("controller text = %q", (DiscardArchives{Player: Controller}).Text())
	}

	e.Resolve(ctx)
	if g.State.Archives[1].Count != 0 {
		t.Error("opponent archives should be emptied")
	}
	if len(g.Discard(1)) != 2 {
		t.Errorf("opponent discard = %d, want 2", len(g.Discard(1)))
	}

	// Empty archives: nothing to move.
	e.Resolve(ctx)
	if len(g.Discard(1)) != 2 {
		t.Error("discarding empty archives should change nothing")
	}
}

// You can see your own archives, so you choose the order the cards are discarded.
func TestDiscardYourOwnArchivesIsOrdered(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0 // player 0 discards their own archives
	c1 := g.Register(testCreature("c1", 1), 0)
	c2 := g.Register(testCreature("c2", 1), 0)
	c3 := g.Register(testCreature("c3", 1), 0)
	g.State.Archives[0].add(c1)
	g.State.Archives[0].add(c2)
	g.State.Archives[0].add(c3)
	g.SetChooser(0, orderLastChooser{}) // choose to discard last-first
	ctx := &EffectContext{Resolver: g, Controller: 0}

	DiscardArchives{Player: Controller}.Resolve(ctx) // your own archives
	if g.State.Archives[0].Count != 0 {
		t.Error("own archives should be emptied")
	}
	got := g.Discard(0)
	want := []LocalID{c3, c2, c1} // the chosen (reversed) order
	if len(got) != 3 || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("own discard order = %v, want %v (owner chooses the order)", got, want)
	}
}

func TestGainAemberPerOpponentArchivedCards(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Archives[1].add(g.Register(testCreature("x", 1), 1))
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := GainAember{
		Player: Controller,
		Amount: 1,
		Per:    CardsInZone{Zone: Archives, Player: Opponent},
	}
	if e.Text() != "for each card in your opponent's archives, gain 1 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Aember(0) != 1 {
		t.Errorf("aember = %d, want 1", g.Aember(0))
	}
}

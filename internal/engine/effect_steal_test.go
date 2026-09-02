package engine

import "testing"

func TestStealAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 1

	e := StealAember{Amount: 3}
	if e.Text() != "steal 3 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx) // opponent has only 1
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 0 {
		t.Errorf("after steal: you=%d opp=%d, want 1/0", g.State.Aember[0], g.State.Aember[1])
	}
}

func TestStealAemberBy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 10

	e := StealAember{By: AllBut(6)}
	if e.Text() != "steal all but 6 Æmber from your opponent" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 4 || g.State.Aember[1] != 6 {
		t.Errorf("after steal: you=%d opp=%d, want 4/6", g.State.Aember[0], g.State.Aember[1])
	}
	if err := validateEffect(StealAember{Amount: 1, By: AllBut(6)}); err == nil {
		t.Error("want error for both Amount and By")
	}
	if err := validateEffect(StealAember{Amount: 1}); err != nil {
		t.Errorf("valid steal rejected: %v", err)
	}
}

func TestStealAemberPer(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	g.AddToBattleline(testCreature("mate", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 5

	e := StealAember{Amount: 1, Per: InPlay{Player: Controller, Type: Creature, Ready: true}}
	if want := "for each friendly ready creature in play, steal 1 Æmber"; e.Text() != want {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 2 || g.State.Aember[1] != 3 {
		t.Errorf("after steal: you=%d opp=%d, want 2/3", g.State.Aember[0], g.State.Aember[1])
	}
}

// TestStealAemberReversed covers the theft turned around, so the opponent takes
// from the controller (Magda the Rat as she leaves play).
func TestStealAemberReversed(t *testing.T) {
	e := StealAember{Player: Opponent, Amount: 2}
	if got := e.Text(); got != "your opponent steals 2 Æmber" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.SetAember(0, 5)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Aember(0) != 3 || g.Aember(1) != 2 {
		t.Errorf("pools = %d/%d, want 3/2", g.Aember(0), g.Aember(1))
	}
}

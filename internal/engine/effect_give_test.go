package engine

import "testing"

func TestGiveAemberValidate(t *testing.T) {
	if err := (GiveAember{Amount: 1}).validate(); err != nil {
		t.Errorf("fixed Amount = %v", err)
	}
	if err := (GiveAember{All: true}).validate(); err != nil {
		t.Errorf("All = %v", err)
	}
	if (GiveAember{}).validate() == nil {
		t.Error("neither Amount nor All should be rejected")
	}
	if (GiveAember{
		Amount: 1,
		All:    true,
	}).validate() == nil {
		t.Error("both Amount and All should be rejected")
	}
}

func TestGiveAemberText(t *testing.T) {
	if got := (GiveAember{Amount: 1}).Text(); got != "your opponent gives you 1 Æmber" {
		t.Errorf("Amount text = %q", got)
	}
	if got := (GiveAember{All: true}).Text(); got != "your opponent gives you all their Æmber" {
		t.Errorf("All text = %q", got)
	}
}

func TestGiveAemberResolve(t *testing.T) {
	// A fixed Amount moves that much from the opponent's pool into the controller's.
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	g.State.Aember[1] = 3
	GiveAember{Amount: 1}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.Aember(1) != 2 || g.Aember(0) != 1 {
		t.Fatalf("fixed give: opponent=%d controller=%d, want 2/1", g.Aember(1), g.Aember(0))
	}

	// All hands over the opponent's whole pool.
	GiveAember{All: true}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.Aember(1) != 0 || g.Aember(0) != 3 {
		t.Fatalf("all give: opponent=%d controller=%d, want 0/3", g.Aember(1), g.Aember(0))
	}

	// The opponent can give only what they hold, so a larger Amount is capped and an
	// empty pool gives nothing.
	g.State.Aember[1] = 1
	g.State.Aember[0] = 0
	GiveAember{Amount: 5}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.Aember(1) != 0 || g.Aember(0) != 1 {
		t.Fatalf("capped give: opponent=%d controller=%d, want 0/1", g.Aember(1), g.Aember(0))
	}
	GiveAember{Amount: 5}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.Aember(1) != 0 || g.Aember(0) != 1 {
		t.Fatalf("empty give: opponent=%d controller=%d, want 0/1", g.Aember(1), g.Aember(0))
	}
}

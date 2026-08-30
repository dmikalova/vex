package engine

import "testing"

func TestForgeKeyEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := ForgeKey{}
	if e.Text() != "forge a key at current cost" {
		t.Errorf("text = %q", e.Text())
	}

	// Not enough Æmber: no key is forged.
	e.Resolve(ctx)
	if g.Keys(0) != 0 {
		t.Errorf("keys = %d, want 0 (could not afford)", g.Keys(0))
	}

	// Enough Æmber: one key is forged and its cost paid.
	g.State.Aember[0] = KeyCost + 2
	e.Resolve(ctx)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (paid the key cost)", g.Aember(0))
	}
}

func TestForgeKeyFreeEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	forger := g.AddToBattleline(exGiant(), 0)
	foe := g.AddToBattleline(testCreature("foe", 4), 1)
	ctx := &EffectContext{Resolver: g, Source: forger, Controller: 0}
	e := ForgeKey{FreeOfCost: true}
	if e.Text() != "forge a key at no cost" {
		t.Errorf("text = %q", e.Text())
	}

	e.Resolve(ctx)

	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 0 {
		t.Errorf("aember = %d, want 0 (no cost paid)", g.Aember(0))
	}
	if g.Damage(foe) != 2 {
		t.Errorf("after-forge ability damage = %d, want 2", g.Damage(foe))
	}
}

func TestGiveRemainingAemberAfterOpponentForgeKey(t *testing.T) {
	e := GiveRemainingAemberAfterOpponentForgeKey{}
	if e.Text() != "if an opponent forges a key on their next turn, they must give you their remaining Æmber" {
		t.Errorf("text = %q", e.Text())
	}

	g := NewGame("A", "B", 1)
	g.BeginTurn(0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.LastingCount != 1 || g.State.Lasting[0].Controller != 1 || g.State.Lasting[0].On != EventForgeKey || g.State.Lasting[0].Once {
		t.Fatal("opponent's next-turn key-forge reaction was not armed as a durable reaction")
	}

	g.State.Aember[1] = KeyCost + 4
	g.EndTurn(0) // the opponent-owned reaction survives the controller's turn end
	g.BeginTurn(1)

	if g.Keys(1) != 1 {
		t.Errorf("opponent keys = %d, want 1", g.Keys(1))
	}
	if g.Aember(1) != 0 {
		t.Errorf("opponent aember = %d, want 0", g.Aember(1))
	}
	if g.Aember(0) != 4 {
		t.Errorf("controller aember = %d, want 4", g.Aember(0))
	}
	if g.State.LastingCount != 1 {
		t.Error("reaction should persist so further forges this turn also transfer")
	}

	g.EndTurn(1)
	if g.State.LastingCount != 0 {
		t.Error("reaction should clear at the end of the opponent's next turn")
	}
}

func TestGiveRemainingAemberAfterOpponentForgeKeyEveryForge(t *testing.T) {
	// A key cheat can forge more than one key in a turn; the opponent must give
	// their remaining Æmber each time.
	g := NewGame("A", "B", 1)
	g.BeginTurn(0)
	GiveRemainingAemberAfterOpponentForgeKey{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	g.EndTurn(0)
	g.BeginTurn(1) // the start-of-turn forge is skipped (no Æmber), reaction still armed

	g.State.Aember[1] = 5
	ForgeKey{FreeOfCost: true}.Resolve(&EffectContext{Resolver: g, Controller: 1})
	if g.Aember(1) != 0 || g.Aember(0) != 5 {
		t.Fatalf("first forge: opponent=%d controller=%d, want 0/5", g.Aember(1), g.Aember(0))
	}

	g.State.Aember[1] = 3
	ForgeKey{FreeOfCost: true}.Resolve(&EffectContext{Resolver: g, Controller: 1})
	if g.Aember(1) != 0 || g.Aember(0) != 8 {
		t.Fatalf("second forge: opponent=%d controller=%d, want 0/8", g.Aember(1), g.Aember(0))
	}
}

func TestGiveRemainingAemberAfterOpponentForgeKeyExpiresIfNoForge(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.BeginTurn(0)
	GiveRemainingAemberAfterOpponentForgeKey{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	g.State.Aember[1] = KeyCost - 1
	g.EndTurn(0)
	g.BeginTurn(1)
	if g.Aember(1) != KeyCost-1 {
		t.Errorf("opponent aember = %d, want %d", g.Aember(1), KeyCost-1)
	}

	g.EndTurn(1)
	if g.State.LastingCount != 0 {
		t.Error("reaction should expire at the end of the opponent's next turn")
	}

	g.State.Aember[1] = KeyCost + 2
	g.BeginTurn(1)
	if g.Aember(1) != 2 {
		t.Errorf("later opponent aember = %d, want 2", g.Aember(1))
	}
	if g.Aember(0) != 0 {
		t.Errorf("controller aember = %d, want 0 (expired transfer)", g.Aember(0))
	}
}

func TestGiveRemainingAemberAfterOpponentForgeKeyAppliesToCardForge(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.BeginTurn(0)
	GiveRemainingAemberAfterOpponentForgeKey{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	g.EndTurn(0)
	g.BeginTurn(1)

	g.State.Aember[1] = 3
	ForgeKey{FreeOfCost: true}.Resolve(&EffectContext{Resolver: g, Controller: 1})

	if g.Aember(1) != 0 {
		t.Errorf("opponent aember = %d, want 0", g.Aember(1))
	}
	if g.Aember(0) != 3 {
		t.Errorf("controller aember = %d, want 3", g.Aember(0))
	}
}

func TestForgeKeyLastingText(t *testing.T) {
	if got := EventForgeKey.clause(); got != "after forging a key" {
		t.Errorf("clause = %q", got)
	}
	if got := actGiveRemainingAember.describe(); got != "give remaining Æmber" {
		t.Errorf("describe = %q", got)
	}
}

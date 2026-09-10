package engine

import (
	"strings"
	"testing"
)

// spendAsPoolCreature is a Senator-Bracchus-style body: while in play, every
// friendly creature's Æmber may be spent as if it were in the controller's pool.
func spendAsPoolCreature() CardDefinition {
	return NewCard("Bracchus", Brobnar, Creature, Rare, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:      Target{Kind: TargetEachFriendlyCreature},
			SpendAsPool: true,
		}))
}

// spendAsPoolUpgrade is a Callipygian-Ideal-style upgrade: the host creature's
// Æmber may be spent as if it were in the controller's pool.
func spendAsPoolUpgrade() CardDefinition {
	return NewCard("Ideal", Brobnar, Upgrade, Uncommon,
		WithStatic(StaticModifier{SpendAsPool: true}))
}

// TestSpendAsPoolForgesFromCreatureConstant covers the forge path: a friendly
// creature's Æmber, made spendable by an in-play constant ability, tops up the
// pool to pay a key.
func TestSpendAsPoolForgesFromCreatureConstant(t *testing.T) {
	g := started(t)
	g.AddToBattleline(spendAsPoolCreature(), 0)
	bank := g.AddToBattleline(testCreature("bank", 4), 0)
	g.AddAmberOn(bank, 2)

	g.State.Aember[0] = KeyCost - 2
	if got := g.spendableAember(0); got != KeyCost {
		t.Fatalf("spendableAember = %d, want %d", got, KeyCost)
	}
	g.forgeKey(0)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 0 {
		t.Errorf("pool = %d, want 0", g.Aember(0))
	}
	if g.AmberOn(bank) != 0 {
		t.Errorf("bank Æmber = %d, want 0 (spent on the key)", g.AmberOn(bank))
	}
}

// TestSpendAsPoolForgesFromHostUpgrade covers the same forge path when the
// permission comes from an attached upgrade (Callipygian Ideal) rather than a
// constant ability.
func TestSpendAsPoolForgesFromHostUpgrade(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 4), 0)
	up := g.Register(spendAsPoolUpgrade(), 0)
	g.AttachUpgrade(host, up)
	g.AddAmberOn(host, 3)

	g.State.Aember[0] = KeyCost - 3
	g.forgeKey(0)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.AmberOn(host) != 0 {
		t.Errorf("host Æmber = %d, want 0", g.AmberOn(host))
	}
}

// TestSpendAsPoolNeededToForge shows the permission is load-bearing: the same
// creature Æmber cannot fund a key without an enabler in play.
func TestSpendAsPoolNeededToForge(t *testing.T) {
	g := started(t)
	bank := g.AddToBattleline(testCreature("bank", 4), 0)
	g.AddAmberOn(bank, 2)

	g.State.Aember[0] = KeyCost - 2
	if got := g.spendableAember(0); got != KeyCost-2 {
		t.Fatalf("spendableAember = %d, want %d", got, KeyCost-2)
	}
	g.forgeKey(0)
	if g.Keys(0) != 0 {
		t.Errorf("keys = %d, want 0 (no enabler, creature Æmber is not pool)", g.Keys(0))
	}
	if g.AmberOn(bank) != 2 {
		t.Errorf("bank Æmber = %d, want 2 (untouched)", g.AmberOn(bank))
	}
}

// TestSpendAsPoolPaysPlayRequirement covers the play-cost path: a card whose play
// spends Æmber (Truebaru) can draw the shortfall from a spend-as-pool creature,
// and cannot be played without one.
func TestSpendAsPoolPaysPlayRequirement(t *testing.T) {
	truebaru := NewCard("Truebaru", Brobnar, Creature, Rare, WithPower(7),
		WithPlayRequirement(AemberCost(3)))

	t.Run("draws the shortfall from the creature", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(spendAsPoolCreature(), 0)
		bank := g.AddToBattleline(testCreature("bank", 4), 0)
		g.AddAmberOn(bank, 2)
		g.SetAember(0, 1) // 1 pool + 2 on the creature covers the cost of 3
		id := g.AddToHand(truebaru, 0)

		if err := g.CanPlay(0, id); err != nil {
			t.Fatalf("CanPlay = %v, want nil", err)
		}
		if _, err := g.PlayCreature(0, 0, false); err != nil {
			t.Fatalf("PlayCreature = %v, want nil", err)
		}
		if g.Aember(0) != 0 {
			t.Errorf("pool = %d, want 0 (pool drawn first)", g.Aember(0))
		}
		if g.AmberOn(bank) != 0 {
			t.Errorf(
				"bank Æmber = %d, want 0 (the other 2 came from the creature)",
				g.AmberOn(bank),
			)
		}
	})

	t.Run("cannot pay without an enabler", func(t *testing.T) {
		g := started(t)
		bank := g.AddToBattleline(testCreature("bank", 4), 0)
		g.AddAmberOn(bank, 2)
		g.SetAember(0, 1)
		id := g.AddToHand(truebaru, 0)

		if err := g.CanPlay(0, id); err != ErrPlayRequirement {
			t.Errorf("CanPlay = %v, want %v", err, ErrPlayRequirement)
		}
		if _, err := g.PlayCreature(0, 0, false); err != ErrPlayRequirement {
			t.Errorf("PlayCreature = %v, want %v", err, ErrPlayRequirement)
		}
	})
}

// TestMeetsPlayRequirement covers the three branches of the affordability helper:
// no requirement, a threshold that only checks the pool (creature Æmber excluded),
// and a spend that may draw from spend-as-pool creatures.
func TestMeetsPlayRequirement(t *testing.T) {
	g := started(t)
	g.AddToBattleline(spendAsPoolCreature(), 0)
	bank := g.AddToBattleline(testCreature("bank", 4), 0)
	g.AddAmberOn(bank, 5)
	g.SetAember(0, 1)

	plain := NewCard("Plain", Brobnar, Creature, Common, WithPower(3))
	if !g.meetsPlayRequirement(0, &plain) {
		t.Error("a card with no requirement is always affordable")
	}

	threshold := NewCard("Kelifi", Brobnar, Creature, Rare, WithPower(9),
		WithPlayRequirement(AemberThreshold(3)))
	if g.meetsPlayRequirement(0, &threshold) {
		t.Error("a threshold checks only the pool; creature Æmber must not count")
	}

	cost := NewCard("Truebaru", Brobnar, Creature, Rare, WithPower(7),
		WithPlayRequirement(AemberCost(3)))
	if !g.meetsPlayRequirement(0, &cost) {
		t.Error("a spend requirement may draw from spend-as-pool creatures")
	}
}

// TestSpendAsPoolCreatureNotCovered exercises the false path: a friendly creature
// no active permission reaches is not spendable.
func TestSpendAsPoolCreatureNotCovered(t *testing.T) {
	g := started(t)
	lone := g.AddToBattleline(testCreature("lone", 4), 0)
	if g.creatureSpendableAsPool(0, lone) {
		t.Error("a creature with no enabler is not spendable as pool")
	}
	if got := g.spendAsPoolCreatures(0); len(got) != 0 {
		t.Errorf("spendAsPoolCreatures = %v, want none", got)
	}
}

// TestDrawFromSpendAsPoolStopsWhenSatisfied covers the early break: with two
// spend-as-pool creatures, a small draw takes from the first and stops before the
// second.
func TestDrawFromSpendAsPoolStopsWhenSatisfied(t *testing.T) {
	g := started(t)
	g.AddToBattleline(spendAsPoolCreature(), 0)
	first := g.AddToBattleline(testCreature("first", 4), 0)
	second := g.AddToBattleline(testCreature("second", 4), 0)
	g.AddAmberOn(first, 3)
	g.AddAmberOn(second, 3)

	if got := g.spendAsPoolTotal(0); got != 6 {
		t.Fatalf("spendAsPoolTotal = %d, want 6", got)
	}
	if got := g.drawFromSpendAsPool(0, 2); got != 2 {
		t.Errorf("drew %d, want 2", got)
	}
	if g.AmberOn(first) != 1 {
		t.Errorf("first = %d, want 1", g.AmberOn(first))
	}
	if g.AmberOn(second) != 3 {
		t.Errorf("second = %d, want 3 (untouched, draw already satisfied)", g.AmberOn(second))
	}
}

// TestSpendAsPoolText covers the rendered rules text for both carriers and the
// non-friendly Target fallback of the subject phrasing.
func TestSpendAsPoolText(t *testing.T) {
	bracchus := spendAsPoolCreature()
	wantB := "You may spend Æmber on friendly creatures as if it were in your pool."
	if got := RenderCardText(&bracchus); !strings.Contains(got, wantB) {
		t.Errorf("constant text = %q, want it to contain %q", got, wantB)
	}

	ideal := spendAsPoolUpgrade()
	wantI := `This creature gains, "You may spend Æmber on this creature as if it were in your pool."`
	if got := RenderCardText(&ideal); !strings.Contains(got, wantI) {
		t.Errorf("upgrade text = %q, want it to contain %q", got, wantI)
	}

	// On a host's face the upgrade's grant renders without the "This creature gains,"
	// frame.
	hosted := spendAsPoolLines(&ideal, true)
	if len(hosted) != 1 ||
		hosted[0] != "You may spend Æmber on this creature as if it were in your pool." {
		t.Errorf("hosted lines = %v, want the unframed permission", hosted)
	}

	// A card with no permission renders no line.
	plain := NewCard("Plain", Brobnar, Creature, Common, WithPower(3))
	if got := spendAsPoolLines(&plain, false); got != nil {
		t.Errorf("spendAsPoolLines(plain) = %v, want nil", got)
	}

	// The subject falls back to the Target's own phrasing for a non-friendly reach.
	if got := spendAsPoolSubject(Target{Kind: TargetEachCreature}); got != "each creature" {
		t.Errorf("subject fallback = %q, want %q", got, "each creature")
	}
}

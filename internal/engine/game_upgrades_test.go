package engine

import "testing"

// TestUpgradeChainStitchesWhenMiddleLeaves attaches several upgrades to one host and
// destroys the one in the middle, exercising the tail-append in AttachUpgrade and the
// predecessor-stitch branch in detachUpgrade: the surviving upgrades stay on the host
// in their original order and the destroyed one goes to its owner's discard pile.
func TestUpgradeChainStitchesWhenMiddleLeaves(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up1 := g.Register(
		NewCard("up1", Brobnar, Upgrade, Common, WithStatic(StaticModifier{PowerBonus: 1})),
		0,
	)
	up2 := g.Register(
		NewCard("up2", Brobnar, Upgrade, Common, WithStatic(StaticModifier{PowerBonus: 2})),
		0,
	)
	up3 := g.Register(
		NewCard("up3", Brobnar, Upgrade, Common, WithStatic(StaticModifier{PowerBonus: 4})),
		0,
	)
	g.AttachUpgrade(host, up1)
	g.AttachUpgrade(host, up2)
	g.AttachUpgrade(host, up3)

	if got := g.Upgrades(host); len(got) != 3 || got[0] != up1 || got[1] != up2 || got[2] != up3 {
		t.Fatalf("upgrades = %v, want [%d %d %d] in attach order", got, up1, up2, up3)
	}

	g.destroyAttachedUpgrade(up2)

	if got := g.Upgrades(host); len(got) != 2 || got[0] != up1 || got[1] != up3 {
		t.Fatalf("upgrades after middle destroyed = %v, want [%d %d]", got, up1, up3)
	}
	if _, ok := g.hostOf(up2); ok {
		t.Error("destroyed upgrade should no longer report a host")
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != up2 {
		t.Errorf("discard = %v, want only the destroyed upgrade", d)
	}
	// Power reflects only the surviving upgrades (base 3 + 1 + 4).
	if p := g.Power(host); p != 8 {
		t.Errorf("host power = %d, want 8 after middle upgrade destroyed", p)
	}

	// Destroying the tail of a longer chain walks past a non-head predecessor to
	// find and re-stitch the link, covering the predecessor walk in detachUpgrade.
	other := g.AddToBattleline(testCreature("other", 3), 0)
	a := g.Register(NewCard("a", Brobnar, Upgrade, Common), 0)
	b := g.Register(NewCard("b", Brobnar, Upgrade, Common), 0)
	c := g.Register(NewCard("c", Brobnar, Upgrade, Common), 0)
	g.AttachUpgrade(other, a)
	g.AttachUpgrade(other, b)
	g.AttachUpgrade(other, c)

	g.destroyAttachedUpgrade(c)

	if got := g.Upgrades(other); len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatalf("upgrades after tail destroyed = %v, want [%d %d]", got, a, b)
	}
}

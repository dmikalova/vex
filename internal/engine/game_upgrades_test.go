package engine

import "testing"

// TestStaticModifierPerUpgrade attaches a Light-of-the-Archons-style upgrade whose
// power/armor bonus scales by the number of upgrades on the host, then adds a
// second upgrade and confirms the bonus grows with the count.
func TestStaticModifierPerUpgrade(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 4), 0)
	light := g.Register(
		NewCard("light", StarAlliance, Upgrade, Common,
			WithStatic(StaticModifier{
				PowerBonus: 1,
				ArmorBonus: 1,
				Per:        UpgradesOnIt,
			})),
		0,
	)
	g.AttachUpgrade(host, light)

	// One upgrade on the host: +1 power and +1 armor.
	if p := g.Power(host); p != 5 {
		t.Errorf("host power with one upgrade = %d, want 5", p)
	}
	if a := g.Armor(host); a != 1 {
		t.Errorf("host armor with one upgrade = %d, want 1", a)
	}

	plain := g.Register(
		NewCard("plain", StarAlliance, Upgrade, Common, WithBonus(BonusAember)),
		0,
	)
	g.AttachUpgrade(host, plain)

	// Two upgrades on the host: the scaled bonus is now +2 power and +2 armor.
	if p := g.Power(host); p != 6 {
		t.Errorf("host power with two upgrades = %d, want 6", p)
	}
	if a := g.Armor(host); a != 2 {
		t.Errorf("host armor with two upgrades = %d, want 2", a)
	}
}

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

// TestStaticModifierKeywordGrantReach attaches a Cloaking-Dongle-style upgrade to a
// creature in the middle of the battleline and confirms its KeywordGrant reaches
// the host and both of its neighbors, but no farther.
func TestStaticModifierKeywordGrantReach(t *testing.T) {
	g := started(t)
	far := g.AddToBattleline(testCreature("far", 3), 0)
	left := g.AddToBattleline(testCreature("left", 3), 0)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)

	dongle := g.Register(
		NewCard("dongle", StarAlliance, Upgrade, Common,
			WithStatic(StaticModifier{KeywordGrants: []KeywordGrant{{
				Keywords: []Keyword{Elusive}, Host: true, Neighbors: true,
			}}})),
		0,
	)
	g.AttachUpgrade(host, dongle)

	if !g.HasKeyword(host, Elusive) {
		t.Error("the host should gain elusive from its own upgrade")
	}
	if !g.HasKeyword(left, Elusive) {
		t.Error("the left neighbor should gain elusive")
	}
	if !g.HasKeyword(right, Elusive) {
		t.Error("the right neighbor should gain elusive")
	}
	if g.HasKeyword(far, Elusive) {
		t.Error("a non-neighbor should not gain elusive")
	}
}

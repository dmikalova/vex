package engine

import "testing"

// TestInvulnerableKeyword covers Ghostform's grant: an invulnerable creature
// takes no damage and cannot be destroyed, whether the keyword is printed on it
// or granted by an attached upgrade.
func TestInvulnerableKeyword(t *testing.T) {
	t.Run("takes no damage", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ghost := g.AddToBattleline(testCreature("ghost", 3, WithKeywords(Invulnerable)), 0)
		g.dealDamage(0, DamageTarget{ID: ghost, Amount: 5})
		if got := g.Damage(ghost); got != 0 {
			t.Errorf("invulnerable creature took %d damage, want 0", got)
		}
		if !g.inPlay(ghost) {
			t.Fatal("invulnerable creature should survive the damage")
		}
	})

	t.Run("cannot be destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ghost := g.AddToBattleline(testCreature("ghost", 3, WithKeywords(Invulnerable)), 0)
		g.destroyEach(0, []LocalID{ghost})
		if !g.inPlay(ghost) {
			t.Fatal("invulnerable creature should survive destruction")
		}
	})

	t.Run("granted by an attached upgrade", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		host := g.AddToBattleline(testCreature("host", 3), 0)
		attachUpgrade(g, host, NewCard("cloak", Brobnar, Upgrade, Rare,
			WithStatic(StaticModifier{Keywords: []Keyword{Invulnerable}})))
		g.destroyEach(0, []LocalID{host})
		if !g.inPlay(host) {
			t.Fatal("host granted invulnerable should survive destruction")
		}
	})
}

// TestArchiveGrantingUpgrade covers Ghostform archiving itself: the effect
// renders the {card} placeholder and sends the granting upgrade off its host to
// the owner's archives while the host stays in play.
func TestArchiveGrantingUpgrade(t *testing.T) {
	if got := (ArchiveGrantingUpgrade{}).Text(); got != "archive "+CardName {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up := attachUpgrade(
		g,
		host,
		NewCard("Ghostform", Brobnar, Upgrade, Rare, WithBonus(BonusAember)),
	)
	ArchiveGrantingUpgrade{}.Resolve(
		&EffectContext{Resolver: g, Source: host, Controller: 0, Upgrade: up},
	)
	if !containsID(g.Archives(0), up) {
		t.Errorf("archived upgrade should be in the owner's archives, got %v", g.Archives(0))
	}
	if _, ok := g.hostOf(up); ok {
		t.Error("archived upgrade should be detached from its host")
	}
	if !g.inPlay(host) {
		t.Error("the host should stay in play")
	}
}

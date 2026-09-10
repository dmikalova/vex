package engine

import (
	"strings"
	"testing"
)

func TestDestroyFriendlyCreaturesToForgeValidate(t *testing.T) {
	good := DestroyFriendlyCreaturesToForge{
		Target:        Target{Kind: TargetEachFriendlyCreature},
		MinTotalPower: 25,
		Then:          ForgeKey{FreeOfCost: true},
	}
	if err := good.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}
	if (DestroyFriendlyCreaturesToForge{MinTotalPower: 25, Then: ForgeKey{}}).validate() == nil {
		t.Error("unset target should be rejected")
	}
	if (DestroyFriendlyCreaturesToForge{Target: Target{Kind: TargetEachFriendlyCreature}, Then: ForgeKey{}}).
		validate() == nil {
		t.Error("non-positive MinTotalPower should be rejected")
	}
	if (DestroyFriendlyCreaturesToForge{Target: Target{Kind: TargetEachFriendlyCreature}, MinTotalPower: 1}).
		validate() == nil {
		t.Error("nil Then should be rejected")
	}
	if got := good.Text(); !strings.Contains(got, "total power of 25 or more") ||
		!strings.Contains(got, "forge a key at no cost") {
		t.Errorf("Text = %q", got)
	}
}

// pickThenDecline picks the given id once, then declines every later ask.
type pickThenDecline struct {
	FirstChooser
	id   LocalID
	done bool
}

func (c *pickThenDecline) ChooseCardOrDecline(_, _ string, _ []LocalID) (LocalID, bool) {
	if c.done {
		return 0, false
	}
	c.done = true
	return c.id, true
}

func TestDestroyFriendlyCreaturesToForgeResolve(t *testing.T) {
	e := DestroyFriendlyCreaturesToForge{
		Target:        Target{Kind: TargetEachFriendlyCreature},
		MinTotalPower: 25,
		Then:          ForgeKey{FreeOfCost: true},
	}

	// Destroying enough total power (the default chooser takes every creature)
	// forges a key for free.
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 13), 0)
	b := g.AddToBattleline(testCreature("b", 12), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.Keys[0] != 1 {
		t.Errorf("keys = %d, want 1 (forged)", g.State.Keys[0])
	}
	if g.inPlay(a) || g.inPlay(b) {
		t.Error("both creatures should be destroyed")
	}

	// Stopping below the threshold destroys nothing and forges no key.
	g2 := NewGame("A", "B", 1)
	c := g2.AddToBattleline(testCreature("c", 10), 0)
	g2.AddToBattleline(testCreature("d", 10), 0)
	g2.SetChooser(0, &pickThenDecline{id: c})
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.State.Keys[0] != 0 {
		t.Errorf("keys = %d, want 0 (below threshold)", g2.State.Keys[0])
	}
	if !g2.inPlay(c) {
		t.Error("creature should not be sacrificed below the threshold")
	}
}

func TestSacrificeToForge(t *testing.T) {
	e := SacrificeToForge{Target: Target{Kind: TargetEachFriendlyCreature}, Extra: 6}
	if (SacrificeToForge{Extra: 6}).validate() == nil {
		t.Error("unset target should be rejected")
	}
	if (SacrificeToForge{Target: Target{Kind: TargetEachFriendlyCreature}}).validate() == nil {
		t.Error("non-positive Extra should be rejected")
	}
	if got := e.Text(); !strings.Contains(
		got,
		"reduced by 1 Æmber for each creature destroyed this way",
	) ||
		!strings.Contains(got, "purge {self}") {
		t.Errorf("Text = %q", got)
	}

	// Sacrificing two creatures drops the +6 surcharge to +4, so 10 Æmber forges
	// the key and Obsidian Forge purges itself.
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("Obsidian", Dis, Artifact, Common), 0)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	b := g.AddToBattleline(testCreature("b", 5), 0)
	g.State.Aember[0] = 10
	e.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})
	if g.State.Keys[0] != 1 {
		t.Errorf("keys = %d, want 1 (forged)", g.State.Keys[0])
	}
	if g.inPlay(a) || g.inPlay(b) {
		t.Error("both creatures should be sacrificed")
	}
	if g.inPlay(src) {
		t.Error("Obsidian Forge should purge itself after forging")
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("aember = %d, want 0 (spent 6 + 6 - 2)", g.State.Aember[0])
	}

	if err := e.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}

	// Unable to pay the reduced cost leaves the artifact intact.
	g3 := NewGame("A", "B", 1)
	src3 := g3.AddArtifact(NewCard("Obsidian", Dis, Artifact, Common), 0)
	g3.AddToBattleline(testCreature("a", 5), 0)
	e.Resolve(&EffectContext{Resolver: g3, Controller: 0, Source: src3})
	if g3.State.Keys[0] != 0 {
		t.Errorf("keys = %d, want 0 (unaffordable)", g3.State.Keys[0])
	}
	if !g3.inPlay(src3) {
		t.Error("an unaffordable forge should leave the artifact in play")
	}
}

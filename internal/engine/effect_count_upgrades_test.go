package engine

import "testing"

// TestUpgradesOnCount covers the count Walls' Blaster scales its stun by: the
// number of upgrades on the creature its Target names, zero when the Target names
// nothing, and the "for each" noun it renders.
func TestUpgradesOnCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up1 := g.Register(NewCard("Up1", StarAlliance, Upgrade, Rare), 0)
	up2 := g.Register(NewCard("Up2", StarAlliance, Upgrade, Rare), 0)
	g.AttachUpgrade(host, up1)
	g.AttachUpgrade(host, up2)

	onHost := UpgradesOn{Target: Target{Kind: TargetThisCreature}}
	if got := onHost.Value(&EffectContext{Resolver: g, Controller: 0, Source: host}); got != 2 {
		t.Errorf("Value with two upgrades = %d, want 2", got)
	}

	// A Target that selects nothing (no attached upgrade in context) counts zero.
	empty := UpgradesOn{Target: Target{Kind: TargetAttachedHost}}
	if got := empty.Value(&EffectContext{Resolver: g, Controller: 0}); got != 0 {
		t.Errorf("Value with no selection = %d, want 0", got)
	}

	named := UpgradesOn{Target: Target{Kind: TargetAttachedHost}.Named("Chief Engineer Walls")}
	if got := named.CountText(); got != "upgrade on Chief Engineer Walls" {
		t.Errorf("CountText = %q", got)
	}
}

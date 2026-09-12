package engine

import "testing"

// blasterUpgrade builds a "blaster" upgrade whose granted Reap ability sends the
// upgrade itself to the signature creature named host — the shared Star Alliance
// blaster mechanic in miniature.
func blasterUpgrade(host string) CardDefinition {
	return NewCard("Test Blaster", StarAlliance, Upgrade, Rare,
		WithStatic(StaticModifier{Granted: []Ability{
			{Trigger: TriggerAfterReap, Effect: AttachSelfTo{Host: host}},
		}}))
}

// TestAttachSelfToMovesGrantingUpgrade covers the whole blaster path: an upgrade
// grants its host a Reap ability that moves the upgrade itself onto its named
// creature. It exercises the ctx.Upgrade plumbing for upgrade-granted abilities,
// AttachSelfTo's name match, and MoveUpgrade relocating an attached upgrade.
func TestAttachSelfToMovesGrantingUpgrade(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	signature := g.AddToBattleline(testCreature("Commander Chan", 4), 0)
	up := g.Register(blasterUpgrade("Commander Chan"), 0)
	g.AttachUpgrade(host, up)

	g.triggerAbilities(host, TriggerAfterReap, 0, false)

	if h, ok := g.hostOf(up); !ok || h != signature {
		t.Fatalf("blaster host = (%d, %v), want Commander Chan (%d)", h, ok, signature)
	}
}

// boundBlasterUpgrade builds a blaster whose granted Reap ability homes onto its
// signature creature and then wards the creature it bound to — the payoff targets
// TargetAttachedHost, the exact instance the blaster now sits on, not the printed
// name.
func boundBlasterUpgrade(host string) CardDefinition {
	return NewCard("Bound Blaster", StarAlliance, Upgrade, Rare,
		WithStatic(StaticModifier{Granted: []Ability{{
			Trigger: TriggerAfterReap,
			Effect: Sequence{Effects: []Effect{
				AttachSelfTo{Host: host},
				Ward{Target: Target{Kind: TargetAttachedHost}.Named(host)},
			}},
		}}}))
}

// TestAttachSelfToBindsHostInstance proves the blaster's payoff acts on the exact
// creature it attached to, never on a second same-named copy. With two "Commander
// Chan" creatures in play, the blaster homes onto the first and wards it; the
// other copy is untouched. Re-triggering keeps the blaster on the same instance
// (no bounce to the second copy) and wards only that one.
func TestAttachSelfToBindsHostInstance(t *testing.T) {
	g := NewGame("A", "B", 1)
	carrier := g.AddToBattleline(testCreature("carrier", 3), 0)
	chanA := g.AddToBattleline(testCreature("Commander Chan", 4), 0)
	chanB := g.AddToBattleline(testCreature("Commander Chan", 4), 0)
	up := g.Register(boundBlasterUpgrade("Commander Chan"), 0)
	g.AttachUpgrade(carrier, up)

	g.triggerAbilities(carrier, TriggerAfterReap, 0, false)

	if h, ok := g.hostOf(up); !ok || h != chanA {
		t.Fatalf("blaster host = (%d, %v), want first Commander Chan (%d)", h, ok, chanA)
	}
	if !g.Warded(chanA) {
		t.Error("bound Commander Chan should be warded")
	}
	if g.Warded(chanB) {
		t.Error("the other Commander Chan must not be warded")
	}

	// Re-trigger from the bound instance: the blaster stays on chanA rather than
	// re-homing onto chanB, and wards only chanA again.
	g.State.Cards[chanA].Warded = false
	g.triggerAbilities(chanA, TriggerAfterReap, 0, false)

	if h, ok := g.hostOf(up); !ok || h != chanA {
		t.Fatalf("blaster re-homed to (%d, %v), want it to stay on chanA (%d)", h, ok, chanA)
	}
	if !g.Warded(chanA) {
		t.Error("bound Commander Chan should be warded again")
	}
	if g.Warded(chanB) {
		t.Error("the other Commander Chan must never be warded")
	}
}

// TestAttachSelfToWithoutNamedCreatureDoesNothing covers the miss branch: when no
// friendly creature carries the host name, the upgrade stays where it is.
func TestAttachSelfToWithoutNamedCreatureDoesNothing(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up := g.Register(blasterUpgrade("Commander Chan"), 0)
	g.AttachUpgrade(host, up)

	g.triggerAbilities(host, TriggerAfterReap, 0, false)

	if h, ok := g.hostOf(up); !ok || h != host {
		t.Fatalf("blaster host = (%d, %v), want unchanged host (%d)", h, ok, host)
	}
}

// TestMoveUpgradeUnattachedIsNoOp covers MoveUpgrade's guard: an id that is not an
// attached upgrade is left alone rather than grafted onto a creature.
func TestMoveUpgradeUnattachedIsNoOp(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	loose := g.AddToBattleline(testCreature("loose", 2), 0)

	g.MoveUpgrade(loose, host)

	if _, ok := g.hostOf(loose); ok {
		t.Fatal("MoveUpgrade attached an unattached card")
	}
}

// TestTargetAttachedHostUnattached covers TargetAttachedHost when the resolving
// upgrade is not attached to any creature: it selects nothing, so the payoff acts
// on no one. It also renders the bare (unnamed) phrase.
func TestTargetAttachedHostUnattached(t *testing.T) {
	g := NewGame("A", "B", 1)
	up := g.Register(blasterUpgrade("Commander Chan"), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, Upgrade: up}

	if got := (Target{Kind: TargetAttachedHost}).Select(ctx); got != nil {
		t.Errorf("unattached AttachedHost selected %v, want nil", got)
	}
	if got := (Target{Kind: TargetAttachedHost}).Text(); got != "the attached creature" {
		t.Errorf("bare AttachedHost text = %q", got)
	}
	if got := (Target{Kind: TargetAttachedHost}).Named("Commander Chan").
		Text(); got != "Commander Chan" {
		t.Errorf("named AttachedHost text = %q, want %q", got, "Commander Chan")
	}
}

// TestAttachSelfToText renders the effect and rejects a missing host name.
func TestAttachSelfToText(t *testing.T) {
	if got := (AttachSelfTo{Host: "Commander Chan"}).Text(); got != "attach "+CardName+" to Commander Chan" {
		t.Errorf("text = %q", got)
	}
	if err := (AttachSelfTo{}).validate(); err == nil {
		t.Error("AttachSelfTo with no host should be invalid")
	}
	if err := (AttachSelfTo{Host: "x"}).validate(); err != nil {
		t.Errorf("valid AttachSelfTo rejected: %v", err)
	}
}

// TestTargetNamedChosenRendersBareName covers the proper-name rendering: a chosen
// single target narrowed to a name prints the name outright, with no article.
func TestTargetNamedChosenRendersBareName(t *testing.T) {
	tgt := Target{Kind: TargetChosenFriendlyCreature}.Named("Lieutenant Khrkhar")
	if got := tgt.Text(); got != "Lieutenant Khrkhar" {
		t.Errorf("text = %q, want %q", got, "Lieutenant Khrkhar")
	}
}

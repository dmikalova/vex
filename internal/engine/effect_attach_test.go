package engine

import "testing"

// blasterUpgrade builds a "blaster" upgrade whose granted Reap ability sends the
// upgrade itself to the signature creature named host — the shared Star Alliance
// blaster mechanic in miniature.
func blasterUpgrade(host string) CardDefinition {
	return NewCard("Test Blaster", StarAlliance, Upgrade, Rare,
		WithStatic(StaticModifier{Granted: []Ability{
			{Trigger: TriggerAfterReap, Effect: AttachSelfTo{
				Target: Target{Kind: TargetChosenFriendlyCreature}.Named(host),
			}},
		}}))
}

// TestAttachSelfToMovesGrantingUpgrade covers the whole blaster path: an upgrade
// grants its host a Reap ability that moves the upgrade itself onto its named
// creature. It exercises the ctx.Upgrade plumbing for upgrade-granted abilities,
// AttachSelfTo's name match, and MoveUpgrade relocating an attached upgrade.
func TestAttachSelfToMovesGrantingUpgrade(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, FirstChooser{})
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
				AttachSelfTo{Target: Target{Kind: TargetChosenFriendlyCreature}.Named(host)},
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
	g.SetChooser(0, idChooser{id: chanA})
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
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		Upgrade:    up,
	}

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

// TestAttachSelfToText renders the effect and rejects a missing target.
func TestAttachSelfToText(t *testing.T) {
	e := AttachSelfTo{Target: Target{Kind: TargetChosenFriendlyCreature}.Named("Commander Chan")}
	if got := e.Text(); got != "attach "+CardName+" to Commander Chan" {
		t.Errorf("text = %q", got)
	}
	if err := (AttachSelfTo{}).validate(); err == nil {
		t.Error("AttachSelfTo with no target should be invalid")
	}
	if err := e.validate(); err != nil {
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

// A "you may attach {card} to a friendly creature" (Blast Shielding) is one
// clickable creature, so May drives it by the click rather than by a Yes/No. An
// upgrade already on an admissible host stays put without asking.
func TestAttachSelfToDeclinable(t *testing.T) {
	chosen := AttachSelfTo{Target: Target{Kind: TargetChosenFriendlyCreature}}
	if !chosen.declinable() {
		t.Error("a chosen AttachSelfTo should be declinable")
	}
	if (AttachSelfTo{Target: Target{Kind: TargetEachFriendlyCreature}}).declinable() {
		t.Error("an untargeted AttachSelfTo should not be declinable")
	}

	empty := NewGame("A", "B", 1)
	if !chosen.vacuous(&EffectContext{
		Resolver:   empty,
		Controller: 0,
	}) {
		t.Error("an AttachSelfTo with no creature to attach to should be vacuous")
	}

	taken := NewGame("A", "B", 1)
	taken.SetChooser(0, &cardDecliner{})
	from := taken.AddToBattleline(testCreature("from", 3), 0)
	up := taken.Register(NewCard("shield", StarAlliance, Upgrade, Common), 0)
	taken.AttachUpgrade(from, up)
	onto := taken.AddToBattleline(testCreature("onto", 3), 0)
	ctx := &EffectContext{
		Resolver:   taken,
		Upgrade:    up,
		Controller: 0,
	}
	// The upgrade's own host is a friendly creature the target admits, so the
	// already-homed shortcut returns before anything is asked.
	if !chosen.resolveOptional(ctx) {
		t.Error("an upgrade already on an admissible host should report attached")
	}
	if h, ok := taken.hostOf(up); !ok || h != from {
		t.Errorf("host = (%d, %v), want the original host %d", h, ok, from)
	}

	// Narrowed to the neighbor by name, the click moves the upgrade there.
	toOnto := AttachSelfTo{Target: Target{Kind: TargetChosenFriendlyCreature}.Named("onto")}
	if !toOnto.resolveOptional(ctx) {
		t.Error("clicking the creature should report the upgrade attached")
	}
	if h, ok := taken.hostOf(up); !ok || h != onto {
		t.Errorf("host = (%d, %v), want the clicked creature %d", h, ok, onto)
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	stay := declined.AddToBattleline(testCreature("stay", 3), 0)
	declined.AddToBattleline(testCreature("other", 3), 0)
	held := declined.Register(NewCard("shield", StarAlliance, Upgrade, Common), 0)
	declined.AttachUpgrade(stay, held)
	toOther := AttachSelfTo{Target: Target{Kind: TargetChosenFriendlyCreature}.Named("other")}
	if toOther.resolveOptional(
		&EffectContext{
			Resolver:   declined,
			Upgrade:    held,
			Controller: 0,
		},
	) {
		t.Error("declining should report nothing attached")
	}
	if h, ok := declined.hostOf(held); !ok || h != stay {
		t.Errorf("host = (%d, %v), want the unchanged host %d", h, ok, stay)
	}
}

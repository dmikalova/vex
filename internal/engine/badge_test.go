package engine

import "testing"

// recordingBadgeChooser picks creatures from a queue and records every badge the
// engine previews, so a test can assert the badges a per-instance effect sends
// around its choose loop.
type recordingBadgeChooser struct {
	ids    []LocalID
	badges []SelectionBadge
}

func (c *recordingBadgeChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(c.ids) == 0 || len(cands) == 0 {
		return 0, false
	}
	id := c.ids[0]
	c.ids = c.ids[1:]
	return id, true
}

func (c *recordingBadgeChooser) PreviewBadge(b SelectionBadge) {
	c.badges = append(c.badges, b)
}

// TestPreviewBadgeIgnoredWithoutCapability covers PreviewBadge forwarding to a
// chooser that cannot show a badge: it is a silent no-op.
func TestPreviewBadgeIgnoredWithoutCapability(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, FirstChooser{})
	g.PreviewBadge(0, SelectionBadge{Icon: DamageIcon, Amount: 3})
}

// TestDealDamagePerInstancePreviewsDamageBadge covers the badge a per-instance
// DealDamage sends its client: the damage badge before the choose loop and the
// zero badge after it.
func TestDealDamagePerInstancePreviewsDamageBadge(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("f1", 5), 0)
	g.AddToBattleline(testCreature("f2", 5), 0)
	e1 := g.AddToBattleline(testCreature("e1", 10), 1)
	ch := &recordingBadgeChooser{ids: []LocalID{e1, e1}}
	g.SetChooser(0, ch)

	DealDamage{
		Amount: 3,
		Per:    InPlay{Player: Controller, Type: Creature},
		Target: Target{Kind: TargetChosenEnemyCreature},
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	want := []SelectionBadge{{Icon: DamageIcon, Amount: 3}, {}}
	if len(ch.badges) != len(want) || ch.badges[0] != want[0] || ch.badges[1] != want[1] {
		t.Errorf("badges = %v, want %v", ch.badges, want)
	}
}

// TestWardAmountPreviewsWardBadge covers the badge a "ward N" sends its client:
// a numberless ward badge before the choose loop and the zero badge after.
func TestWardAmountPreviewsWardBadge(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	ch := &recordingBadgeChooser{ids: []LocalID{a, c}}
	g.SetChooser(0, ch)

	Ward{Target: Target{Kind: TargetEachFriendlyCreature}, Amount: 2}.Resolve(
		&EffectContext{Resolver: g, Source: a, Controller: 0},
	)

	want := []SelectionBadge{{Icon: WardIcon}, {}}
	if len(ch.badges) != len(want) || ch.badges[0] != want[0] || ch.badges[1] != want[1] {
		t.Errorf("badges = %v, want %v", ch.badges, want)
	}
}

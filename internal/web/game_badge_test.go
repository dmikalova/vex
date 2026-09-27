package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// preview mirrors a selection badge onto the client the way a pending Request
// carrying one does: the engine stamps the badge an effect previewed onto every
// request its choose loop yields, and syncBadge is what reads it off.
func (c *client) preview(badge engine.SelectionBadge) {
	c.t.Helper()
	c.g.syncBadge(engine.Request{Badge: badge}, true)
	c.settle()
}

// A damage badge preview arms the badge, and picks accumulate the running damage
// on each creature — a creature chosen twice sums its hits.
func TestBadgePreviewAccumulatesDamage(t *testing.T) {
	c := newClient(t)
	c.preview(engine.SelectionBadge{
		Icon:   engine.DamageIcon,
		Amount: 3,
	})
	if c.g.selBadge != (engine.SelectionBadge{
		Icon:   engine.DamageIcon,
		Amount: 3,
	}) {
		t.Fatalf("preview did not arm the badge: %+v", c.g.selBadge)
	}

	c.g.recordBadge(7)
	c.g.recordBadge(7)
	c.g.recordBadge(9)
	if c.g.badgeTotals[7] != 6 || c.g.badgeTotals[9] != 3 {
		t.Errorf("totals = %v, want 7:6 9:3", c.g.badgeTotals)
	}
	if c.g.cardBadge(7) == nil {
		t.Error("a badged card should render an overlay")
	}
	if c.g.cardBadge(1) != nil {
		t.Error("an unbadged card should render no overlay")
	}
}

// A ward badge carries no number, so each pick records the creature with a zero
// total — a numberless icon — rather than a running amount.
func TestBadgePreviewWardIsNumberless(t *testing.T) {
	c := newClient(t)
	c.preview(engine.SelectionBadge{Icon: engine.WardIcon})

	c.g.recordBadge(3)
	if total, ok := c.g.badgeTotals[3]; !ok || total != 0 {
		t.Errorf("ward pick recorded %d (present %v), want 0/true", total, ok)
	}
	if c.g.cardBadge(3) == nil {
		t.Error("a warded pick should still draw its icon")
	}
}

// Ending a preview that landed badges grows and fades them, then clears the
// preview once the animation has run.
func TestBadgePreviewClears(t *testing.T) {
	c := newClient(t)
	c.preview(engine.SelectionBadge{
		Icon:   engine.DamageIcon,
		Amount: 2,
	})
	c.g.recordBadge(4)

	c.preview(engine.SelectionBadge{})
	if !c.g.badgeClearing {
		t.Fatal("ending a preview with badges should start the fade")
	}

	// The fade is the one thing in the client that is still on a clock, so this is
	// the one place a test waits rather than asserts.
	c.awaitTimer("the badges to clear", 2*badgeFadeDur, func() bool {
		return c.g.badgeTotals == nil && !c.g.badgeClearing
	})
	if c.g.selBadge.Icon != engine.NoStatusIcon {
		t.Error("a cleared preview should drop the badge")
	}
}

// Without a preview armed, a pick records no badge, so an ordinary card prompt
// draws no overlay.
func TestNoBadgeWithoutPreview(t *testing.T) {
	c := newClient(t)
	c.g.recordBadge(5)
	if len(c.g.badgeTotals) != 0 {
		t.Errorf("recorded a badge with no preview armed: %v", c.g.badgeTotals)
	}
	if c.g.cardBadge(5) != nil {
		t.Error("an unbadged card should render no overlay")
	}
}

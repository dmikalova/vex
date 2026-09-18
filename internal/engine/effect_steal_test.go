package engine

import (
	"strings"
	"testing"
)

func TestStealAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 1

	e := StealAember{Amount: 3}
	if e.Text() != "steal 3 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx) // opponent has only 1
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 0 {
		t.Errorf("after steal: you=%d opp=%d, want 1/0", g.State.Aember[0], g.State.Aember[1])
	}
}

func TestStealAemberBy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 10

	e := StealAember{By: AllBut(6)}
	if e.Text() != "steal all but 6 Æmber from your opponent" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 4 || g.State.Aember[1] != 6 {
		t.Errorf("after steal: you=%d opp=%d, want 4/6", g.State.Aember[0], g.State.Aember[1])
	}
	if err := validateEffect(StealAember{Amount: 1, By: AllBut(6)}); err == nil {
		t.Error("want error for both Amount and By")
	}
	if err := validateEffect(StealAember{Amount: 1}); err != nil {
		t.Errorf("valid steal rejected: %v", err)
	}
}

func TestStealAemberPer(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	g.AddToBattleline(testCreature("mate", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 5

	e := StealAember{Amount: 1, Per: CardsInPlay{Player: Controller, Type: Creature, Ready: true}}
	if want := "for each friendly ready creature in play, steal 1 Æmber"; e.Text() != want {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 2 || g.State.Aember[1] != 3 {
		t.Errorf("after steal: you=%d opp=%d, want 2/3", g.State.Aember[0], g.State.Aember[1])
	}
}

// TestStealAemberReversed covers the theft turned around, so the opponent takes
// from the controller (Magda the Rat as she leaves play).
func TestStealAemberReversed(t *testing.T) {
	e := StealAember{Player: Opponent, Amount: 2}
	if got := e.Text(); got != "your opponent steals 2 Æmber" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.SetAember(0, 5)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Aember(0) != 3 || g.Aember(1) != 2 {
		t.Errorf("pools = %d/%d, want 3/2", g.Aember(0), g.Aember(1))
	}
}

// TestStealAemberCapturedByRedirect covers Gargantodon's continuous replacement:
// a steal's Æmber never reaches the thief's pool — it is captured onto a creature
// the thief controls instead, while the victim's pool still drops.
func TestStealAemberCapturedByRedirect(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("garg", 16,
		WithReplaces(Instead{Of: EventAemberStolen, With: Capture})), 1)
	thief := g.AddToBattleline(testCreature("thief", 3), 0)
	g.State.Aember[1] = 5
	ctx := &EffectContext{Resolver: g, Source: thief, Controller: 0}

	StealAember{Amount: 2}.Resolve(ctx)
	if g.Aember(0) != 0 {
		t.Errorf("thief pool = %d, want 0 (captured, not pooled)", g.Aember(0))
	}
	if g.Aember(1) != 3 {
		t.Errorf("victim pool = %d, want 3", g.Aember(1))
	}
	if got := g.AmberOn(thief); got != 2 {
		t.Errorf("captured on creature = %d, want 2", got)
	}
}

// TestStealAemberRedirectNoCreatureFallsBack covers the redirect being active but
// the thief controlling no creature to hold the capture: the steal lands in the
// pool as usual.
func TestStealAemberRedirectNoCreatureFallsBack(t *testing.T) {
	g := NewGame("A", "B", 1)
	garg := g.AddToBattleline(testCreature("garg", 16,
		WithReplaces(Instead{Of: EventAemberStolen, With: Capture})), 1)
	g.State.Aember[1] = 5
	ctx := &EffectContext{Resolver: g, Source: garg, Controller: 0}

	StealAember{Amount: 2}.Resolve(ctx)
	if g.Aember(0) != 2 {
		t.Errorf("thief pool = %d, want 2 (normal steal)", g.Aember(0))
	}
	if g.Aember(1) != 3 {
		t.Errorf("victim pool = %d, want 3", g.Aember(1))
	}
}

// TestStealAemberRedirectChoosesCaptor covers the redirect with several friendly
// creatures: the thief chooses which one captures the stolen Æmber.
func TestStealAemberRedirectChoosesCaptor(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("garg", 16,
		WithReplaces(Instead{Of: EventAemberStolen, With: Capture})), 1)
	first := g.AddToBattleline(testCreature("first", 3), 0)
	second := g.AddToBattleline(testCreature("second", 3), 0)
	g.State.Aember[1] = 5
	g.SetChooser(0, idChooser{id: second})
	ctx := &EffectContext{Resolver: g, Source: first, Controller: 0}

	StealAember{Amount: 2}.Resolve(ctx)
	if got := g.AmberOn(second); got != 2 {
		t.Errorf("captured on chosen = %d, want 2", got)
	}
	if got := g.AmberOn(first); got != 0 {
		t.Errorf("unchosen holds = %d, want 0", got)
	}
}

// TestStealAemberRedirectDeclineFallsToFirst covers a declined choice among
// several captors: the redirect falls back to the first friendly creature rather
// than dropping the capture.
func TestStealAemberRedirectDeclineFallsToFirst(t *testing.T) {
	g := NewGame("A", "B", 1)
	garg := g.AddToBattleline(testCreature("garg", 16,
		WithReplaces(Instead{Of: EventAemberStolen, With: Capture})), 1)
	first := g.AddToBattleline(testCreature("first", 3), 0)
	g.AddToBattleline(testCreature("second", 3), 0)
	g.State.Aember[1] = 5
	// A chooser that never matches a candidate declines, so the fallback applies.
	g.SetChooser(0, idChooser{id: garg})
	ctx := &EffectContext{Resolver: g, Source: first, Controller: 0}

	StealAember{Amount: 2}.Resolve(ctx)
	if got := g.AmberOn(first); got != 2 {
		t.Errorf("captured on first (fallback) = %d, want 2", got)
	}
}

// TestCaptureStolenAemberText renders Gargantodon's continuous replacement line.
func TestCaptureStolenAemberText(t *testing.T) {
	def := testCreature("garg", 16,
		WithReplaces(Instead{Of: EventAemberStolen, With: Capture}))
	want := "Each Æmber that would be stolen is captured by a creature controlled by the active player instead."
	if got := RenderCardRules(&def); !strings.Contains(got, want) {
		t.Errorf("rules missing redirect line:\n%s", got)
	}
}

// A creature with StealsInsteadOfDamageWhenAttacked deals no retaliation damage;
// instead its controller steals — Shoulder Id. The attacker still takes no damage
// and Shoulder Id still takes the attacker's fight damage.
func TestStealsInsteadOfDamageWhenAttacked(t *testing.T) {
	shoulder := NewCard("Shoulder Id", Shadows, Creature, Common,
		WithPower(6), WithStealsInsteadOfDamageWhenAttacked(1))
	if got := RenderCardRules(&shoulder); !strings.Contains(got,
		"When Shoulder Id would deal damage, steal 1 Æmber instead.") {
		t.Fatalf("Shoulder Id rules = %q", got)
	}

	g := NewGame("A", "B", 1)
	attacker := g.AddToBattleline(testCreature("attacker", 5), 0)
	defender := g.AddToBattleline(shoulder, 1)
	g.State.ActivePlayer = 0
	g.SetAember(0, 3)

	if err := g.Fight(0, attacker, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if !g.inPlay(attacker) {
		t.Error("attacker should survive: Shoulder Id deals no retaliation damage")
	}
	if got := g.Damage(attacker); got != 0 {
		t.Errorf("attacker damage = %d, want 0", got)
	}
	if got := g.Damage(defender); got != 5 {
		t.Errorf("Shoulder Id damage = %d, want 5 (it still takes fight damage)", got)
	}
	if got := g.Aember(1); got != 1 {
		t.Errorf("defender controller Æmber = %d, want 1 (stolen)", got)
	}
	if got := g.Aember(0); got != 2 {
		t.Errorf("attacker controller Æmber = %d, want 2 (robbed of 1)", got)
	}
}

// When the attacker's controller has no Æmber, the steal takes nothing, but the
// retaliation damage is still replaced — the attacker takes no damage.
func TestStealsInsteadOfDamageWhenAttackedEmptyPool(t *testing.T) {
	shoulder := NewCard("Shoulder Id", Shadows, Creature, Common,
		WithPower(6), WithStealsInsteadOfDamageWhenAttacked(1))
	g := NewGame("A", "B", 1)
	attacker := g.AddToBattleline(testCreature("attacker", 5), 0)
	defender := g.AddToBattleline(shoulder, 1)
	g.State.ActivePlayer = 0

	if err := g.Fight(0, attacker, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if got := g.Damage(attacker); got != 0 {
		t.Errorf("attacker damage = %d, want 0", got)
	}
	if got := g.Aember(1); got != 0 {
		t.Errorf("defender controller Æmber = %d, want 0 (nothing to steal)", got)
	}
}

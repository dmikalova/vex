package engine

import "testing"

func TestExaltEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if got := (Exalt{Target: Target{Kind: TargetChosenFriendlyCreature}, Amount: 1}).Text(); got != "exalt a friendly creature" {
		t.Errorf("single exalt text = %q", got)
	}
	e := Exalt{Target: Target{Kind: TargetChosenEnemyCreature}, Amount: 2}
	if e.Text() != "exalt an enemy creature 2 times" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Cards[enemy].Amber != 2 {
		t.Errorf("amber on enemy = %d, want 2", g.State.Cards[enemy].Amber)
	}

	// No candidates: remove the enemy and resolve again (logs, no panic).
	g.DestroyEach(0, []LocalID{enemy})
	e.Resolve(ctx)
}

// exaltRepeater accepts the first exalt-to-repeat prompt and declines the next,
// so the preceding effect resolves exactly twice.
type exaltRepeater struct {
	FirstChooser
	calls int
}

func (c *exaltRepeater) ChooseCardOrDecline(
	_, _ string,
	candidates []LocalID,
) (LocalID, bool) {
	c.calls++
	if c.calls == 1 {
		return candidates[0], true
	}
	return 0, false
}

func TestExaltToRepeatResolvesThenStopsWhenDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	pay := g.AddToBattleline(testCreature("pay", 3), 0)
	g.SetChooser(0, &exaltRepeater{})
	ctx := &EffectContext{Resolver: g, Source: pay, Controller: 0}

	e := ExaltToRepeat{
		Do:    GainAember{Player: Controller, Amount: 1},
		Exalt: Target{Kind: TargetChosenFriendlyCreature},
	}
	if got := e.Text(); got !=
		"gain 1 \u00c6mber. You may exalt a friendly creature to repeat the preceding effect" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)

	// Do runs once up front, then once more after the single accepted exalt.
	if got := g.State.Aember[0]; got != 2 {
		t.Errorf("pool = %d, want 2 (Do resolved twice)", got)
	}
	// The accepted exalt placed 1 Æmber on the paying creature.
	if got := g.State.Cards[pay].Amber; got != 1 {
		t.Errorf("exalted amber = %d, want 1", got)
	}
}

func TestExaltToRepeatValidate(t *testing.T) {
	full := ExaltToRepeat{
		Do:    GainAember{Player: Controller, Amount: 1},
		Exalt: Target{Kind: TargetChosenFriendlyCreature},
	}
	if err := validateEffect(full); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}
	if (ExaltToRepeat{Exalt: Target{Kind: TargetChosenFriendlyCreature}}).validate() == nil {
		t.Error("missing Do should be rejected")
	}
	if (ExaltToRepeat{Do: GainAember{Player: Controller, Amount: 1}}).validate() == nil {
		t.Error("unset exalt target should be rejected")
	}
}

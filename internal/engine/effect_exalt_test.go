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

// TestExaltDistinctSpread covers Bawretchadontius: a Times exalt spreads across
// that many distinct creatures, and stops once the pool of distinct creatures
// runs out.
func TestExaltDistinctSpread(t *testing.T) {
	spread := Exalt{
		Target:   Target{Kind: TargetChosenEnemyCreature},
		Amount:   1,
		Times:    Fixed(2),
		Distinct: true,
	}
	if got := spread.Text(); got != "exalt 2 enemy creatures" {
		t.Errorf("distinct-spread text = %q, want %q", got, "exalt 2 enemy creatures")
	}
	if err := (Exalt{Target: Target{Kind: TargetChosenEnemyCreature}, Distinct: true}).validate(); err == nil {
		t.Error("Distinct without a Times should be rejected")
	}

	t.Run("exalts two distinct creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 1), 0)
		e1 := g.AddToBattleline(testCreature("e1", 1), 1)
		e2 := g.AddToBattleline(testCreature("e2", 1), 1)
		e3 := g.AddToBattleline(testCreature("e3", 1), 1)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{e1, e2}})
		spread.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
		if g.AmberOn(e1) != 1 || g.AmberOn(e2) != 1 || g.AmberOn(e3) != 0 {
			t.Errorf("amber e1/e2/e3 = %d/%d/%d, want 1/1/0",
				g.AmberOn(e1), g.AmberOn(e2), g.AmberOn(e3))
		}
	})

	t.Run("stops when distinct creatures run out", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 1), 0)
		only := g.AddToBattleline(testCreature("only", 1), 1)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{only}})
		spread.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
		if g.AmberOn(only) != 1 {
			t.Errorf("amber on only = %d, want 1", g.AmberOn(only))
		}
	})
}

// A "you may exalt <self>" is one clickable card — the source — so it is offered
// declinably (Senator Shrix): clicking the source confirms, Done declines.
func TestMayExaltSelfDeclinable(t *testing.T) {
	self := Exalt{Target: Target{Kind: TargetThisCreature}, Amount: 1}
	if !self.declinable() {
		t.Fatal("a self-exalt should be declinable")
	}
	if (Exalt{Target: Target{Kind: TargetChosenEnemyCreature}}).declinable() {
		t.Error("a chosen-target exalt is not offered as clicking the source")
	}

	t.Run("accepted exalts the clicked source", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ch := &cardDecliner{}
		g.SetChooser(0, ch)
		src := g.AddToBattleline(testCreature("src", 3), 0)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		May{Do: self}.Resolve(ctx)

		if ch.asked != 1 {
			t.Errorf("declinable prompts = %d, want 1", ch.asked)
		}
		if g.State.Cards[src].Amber != 1 {
			t.Errorf("amber on source = %d, want 1", g.State.Cards[src].Amber)
		}
	})

	t.Run("declined exalts nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetChooser(0, &cardDecliner{decline: true})
		src := g.AddToBattleline(testCreature("src", 3), 0)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		May{Do: self}.Resolve(ctx)

		if g.State.Cards[src].Amber != 0 {
			t.Errorf("a declined May should exalt nothing, amber = %d", g.State.Cards[src].Amber)
		}
	})
}

// exaltRepeater accepts the exalt-to-repeat prompt, so the preceding effect
// resolves twice: once up front and once for the single allowed repeat.
type exaltRepeater struct {
	FirstChooser
}

func (exaltRepeater) ChooseCardOrDecline(
	_, _ string,
	candidates []LocalID,
) (LocalID, bool) {
	return candidates[0], true
}

func TestRepeatByExaltingResolvesThenStopsWhenDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	pay := g.AddToBattleline(testCreature("pay", 3), 0)
	g.SetChooser(0, &exaltRepeater{})
	ctx := &EffectContext{Resolver: g, Source: pay, Controller: 0}

	e := Repeat{
		Do:   GainAember{Player: Controller, Amount: 1},
		Gate: ByExalting{Creature: Target{Kind: TargetChosenFriendlyCreature}},
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

// exaltConfirmer accepts the Yes/No exalt-to-repeat confirm, so a back-reference
// exalt repeats once.
type exaltConfirmer struct {
	FirstChooser
}

func (exaltConfirmer) ChooseOption(_, _ string, _ []string) int {
	return 0 // Yes
}

// exaltDecliner declines the Yes/No exalt-to-repeat confirm.
type exaltDecliner struct {
	FirstChooser
}

func (exaltDecliner) ChooseOption(_, _ string, _ []string) int {
	return 1 // No
}

func TestRepeatByExaltingConfirmsBackReference(t *testing.T) {
	g := NewGame("A", "B", 1)
	that := g.AddToBattleline(testCreature("that", 3), 0)
	g.SetChooser(0, &exaltConfirmer{})
	ctx := &EffectContext{Resolver: g, Source: that, Controller: 0, It: that, HasIt: true}

	e := Repeat{
		Do:   GainAember{Player: Controller, Amount: 1},
		Gate: ByExalting{Creature: Target{Kind: TargetTheChosenCreature}},
	}
	if got := e.Text(); got !=
		"gain 1 \u00c6mber. You may exalt the chosen creature to repeat the preceding effect" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)

	// Do runs once, then once more after the single confirmed exalt.
	if got := g.State.Aember[0]; got != 2 {
		t.Errorf("pool = %d, want 2 (Do resolved twice)", got)
	}
	// The confirmed exalt placed 1 Æmber on the context creature.
	if got := g.State.Cards[that].Amber; got != 1 {
		t.Errorf("exalted amber = %d, want 1", got)
	}
}

func TestRepeatByExaltingDeclinesBackReference(t *testing.T) {
	g := NewGame("A", "B", 1)
	that := g.AddToBattleline(testCreature("that", 3), 0)
	g.SetChooser(0, &exaltDecliner{})
	ctx := &EffectContext{Resolver: g, Source: that, Controller: 0, It: that, HasIt: true}

	e := Repeat{
		Do:   GainAember{Player: Controller, Amount: 1},
		Gate: ByExalting{Creature: Target{Kind: TargetTheChosenCreature}},
	}
	e.Resolve(ctx)

	// Declining the confirm resolves Do only once and exalts nothing.
	if got := g.State.Aember[0]; got != 1 {
		t.Errorf("pool = %d, want 1 (Do resolved once)", got)
	}
	if got := g.State.Cards[that].Amber; got != 0 {
		t.Errorf("exalted amber = %d, want 0", got)
	}
}

func TestRepeatByExaltingBackReferenceStopsWithoutContext(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := Repeat{
		Do:   GainAember{Player: Controller, Amount: 1},
		Gate: ByExalting{Creature: Target{Kind: TargetTheChosenCreature}},
	}
	e.Resolve(ctx)

	// With no context creature there is nothing to exalt, so Do resolves once.
	if got := g.State.Aember[0]; got != 1 {
		t.Errorf("pool = %d, want 1 (Do resolved once)", got)
	}
}

func TestRepeatByExaltingValidate(t *testing.T) {
	full := Repeat{
		Do:   GainAember{Player: Controller, Amount: 1},
		Gate: ByExalting{Creature: Target{Kind: TargetChosenFriendlyCreature}},
	}
	if err := validateEffect(full); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}
	if (ByExalting{Creature: Target{Kind: TargetChosenFriendlyCreature}}).validate() != nil {
		t.Error("valid exalt gate rejected")
	}
	if (ByExalting{}).validate() == nil {
		t.Error("unset exalt target should be rejected")
	}
}

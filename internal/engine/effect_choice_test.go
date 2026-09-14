package engine

import "testing"

// optionPicker is a Chooser that also chooses a fixed "choose one" option index.
type optionPicker struct {
	FirstChooser
	idx int
}

func (o optionPicker) ChooseOption(_, _ string, _ []string) int { return o.idx }

func TestChooseOne(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := ChooseOne{
		Options: []Effect{
			GainAember{Player: Controller, Amount: 1},
			GainAember{Player: Controller, Amount: 5},
		},
	}
	if e.Text() != "choose one:\n- Gain 1 Æmber\n- Gain 5 Æmber" {
		t.Errorf("text = %q", e.Text())
	}

	// The default chooser expresses no preference, so the first option is taken.
	e.Resolve(ctx)
	if g.Aember(0) != 1 {
		t.Errorf("default option: aember = %d, want 1", g.Aember(0))
	}

	// A chooser that selects the second option.
	g.SetChooser(0, optionPicker{idx: 1})
	e.Resolve(ctx)
	if g.Aember(0) != 6 { // 1 + 5
		t.Errorf("second option: aember = %d, want 6", g.Aember(0))
	}

	// An out-of-range choice resolves nothing.
	g.SetChooser(0, optionPicker{idx: 9})
	e.Resolve(ctx)
	if g.Aember(0) != 6 {
		t.Errorf("out-of-range option changed aember to %d, want 6", g.Aember(0))
	}
}

func TestChooseOneValidate(t *testing.T) {
	bad := ChooseOne{
		Options: []Effect{Heal{Fully: true, Amount: 1, Target: Target{Kind: TargetThisCreature}}},
	}
	if validateEffect(bad) == nil {
		t.Error("ChooseOne with an invalid option should fail validation")
	}
	good := ChooseOne{Options: []Effect{GainAember{Player: Controller, Amount: 1}, Draw{Amount: 1}}}
	if validateEffect(good) != nil {
		t.Error("ChooseOne with valid options should pass validation")
	}
}

func TestChooseHouseThen(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 1)
	sanc := g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(3)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := ChooseHouseThen{
		Then: Stun{Target: Target{Kind: TargetEachEnemyCreature}.House(chosenHouse)},
	}
	if e.Text() != "choose a house - stun each enemy creature of the chosen house" {
		t.Errorf("text = %q", e.Text())
	}
	// Mars is the fourth house (option index 3).
	g.SetChooser(0, optionPicker{idx: 3})
	e.Resolve(ctx)
	if !g.State.Cards[mars].Stunned {
		t.Error("the chosen-house (Mars) creature should be stunned")
	}
	if g.State.Cards[sanc].Stunned {
		t.Error("a creature of a different house should not be stunned")
	}
}

// TestOpponentNamesHouse covers the opponent naming a house onto the context, and
// the guard when a misbehaving chooser returns an out-of-range option.
func TestOpponentNamesHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	// The opponent of controller 0 is player 1; they name the house.
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := OpponentNamesHouse{}
	if e.Text() != "they name a house" {
		t.Errorf("text = %q", e.Text())
	}

	// Logos is the third house (option index 2); it lands on the context.
	g.SetChooser(1, optionPicker{idx: 2})
	e.Resolve(ctx)
	if ctx.ChosenHouse != Logos {
		t.Errorf("named house = %v, want Logos", ctx.ChosenHouse)
	}

	// An out-of-range choice names no house, leaving the context unchanged.
	ctx.ChosenHouse = HouseNone
	g.SetChooser(1, optionPicker{idx: 99})
	e.Resolve(ctx)
	if ctx.ChosenHouse != HouseNone {
		t.Errorf("out-of-range choice named %v, want HouseNone", ctx.ChosenHouse)
	}
}

func TestChooseHouseThenGuardsAndValidate(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// An out-of-range house choice resolves nothing.
	g.SetChooser(0, optionPicker{idx: 99})
	ChooseHouseThen{
		Then: Stun{Target: Target{Kind: TargetEachCreature}.House(chosenHouse)},
	}.Resolve(
		ctx,
	)
	if g.State.Cards[mars].Stunned {
		t.Error("an out-of-range house choice should stun nothing")
	}

	bad := ChooseHouseThen{
		Then: Heal{Fully: true, Amount: 1, Target: Target{Kind: TargetThisCreature}},
	}
	if validateEffect(bad) == nil {
		t.Error("ChooseHouseThen should surface an invalid Then via validate")
	}
	good := ChooseHouseThen{Then: Stun{Target: Target{Kind: TargetEachCreature}}}
	if validateEffect(good) != nil {
		t.Error("ChooseHouseThen with a valid Then should pass validation")
	}
}

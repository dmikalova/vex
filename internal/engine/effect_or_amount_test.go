package engine

import "testing"

// TestOrAmountStealAember covers the linear "steal 1 Æmber, or 2 if …" form
// (card-wording rule 22): the alternate amount is taken only when the guard holds,
// and the tail renders in place of a two-armed Otherwise branch (Ronnie Wristclocks).
func TestOrAmountStealAember(t *testing.T) {
	e := StealAember{
		Amount: 1,
		Or:     OrAmount{Amount: 2, When: OpponentAember{Is: AtLeast, Amount: 7}},
	}
	if want := "steal 1 Æmber, or 2 if your opponent has 7 Æmber or more"; e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	// Guard unmet: the base amount is stolen.
	g := NewGame("A", "B", 1)
	g.State.Aember[1] = 6
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 5 {
		t.Errorf("unmet: you=%d opp=%d, want 1/5", g.State.Aember[0], g.State.Aember[1])
	}

	// Guard met: the alternate amount is stolen.
	g = NewGame("A", "B", 1)
	g.State.Aember[1] = 7
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.Aember[0] != 2 || g.State.Aember[1] != 5 {
		t.Errorf("met: you=%d opp=%d, want 2/5", g.State.Aember[0], g.State.Aember[1])
	}
}

// TestOrAmountStealValidate covers the conflict and guard-validation checks.
func TestOrAmountStealValidate(t *testing.T) {
	good := StealAember{
		Amount: 1,
		Or:     OrAmount{Amount: 2, When: OpponentAember{Is: AtLeast, Amount: 7}},
	}
	if err := good.validate(); err != nil {
		t.Errorf("valid steal rejected: %v", err)
	}
	if err := (StealAember{
		Or: OrAmount{Amount: 2, When: OpponentAember{Is: AtLeast, Amount: 7}},
		By: AllBut(6),
	}).validate(); err == nil {
		t.Error("Or and By together should not validate")
	}
	if err := (StealAember{
		Amount: 1,
		Or:     OrAmount{Amount: 2, When: OpponentAember{}},
	}).validate(); err == nil {
		t.Error("an Or with an invalid guard should not validate")
	}
}

// TestOrAmountForgeKey covers the "forge a key at +6 …, or +2 if …" form
// (Key of Darkness): the surcharge drops to the alternate only when the guard holds.
func TestOrAmountForgeKey(t *testing.T) {
	e := ForgeKey{
		Extra: 6,
		Or:    OrAmount{Amount: 2, When: OpponentAember{Is: Exactly, Amount: 0}},
	}
	want := "forge a key at +6 Æmber current cost, or +2 if your opponent has no Æmber"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	// Guard unmet (opponent holds Æmber): the +6 surcharge is paid.
	g := started(t)
	g.State.Aember[0] = 100
	g.State.Aember[1] = 1
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if got := g.State.Aember[0]; got != 100-(KeyCost+6) {
		t.Errorf("Æmber = %d, want the +6 surcharge paid", got)
	}

	// Guard met (opponent has no Æmber): the +2 surcharge is paid.
	g = started(t)
	g.State.Aember[0] = 100
	g.State.Aember[1] = 0
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if got := g.State.Aember[0]; got != 100-(KeyCost+2) {
		t.Errorf("Æmber = %d, want the +2 surcharge paid", got)
	}
}

// TestOrAmountForgeKeyValidate covers the free-forge conflict and guard validation.
func TestOrAmountForgeKeyValidate(t *testing.T) {
	if err := (ForgeKey{
		FreeOfCost: true,
		Or:         OrAmount{Amount: 2, When: OpponentAember{Is: Exactly, Amount: 0}},
	}).validate(); err == nil {
		t.Error("a free forge with an Or surcharge should not validate")
	}
	if err := (ForgeKey{
		Extra: 6,
		Or:    OrAmount{Amount: 2, When: OpponentAember{}},
	}).validate(); err == nil {
		t.Error("an Or with an invalid guard should not validate")
	}
}

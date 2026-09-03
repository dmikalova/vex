package engine

import "testing"

func TestForgeKeyExtraCost(t *testing.T) {
	cases := []struct {
		name   string
		effect ForgeKey
		want   string
	}{
		{"current cost", ForgeKey{}, "forge a key at current cost"},
		{"free", ForgeKey{FreeOfCost: true}, "forge a key at no cost"},
		{"surcharge", ForgeKey{Extra: 6}, "forge a key at +6 Æmber current cost"},
		{
			"reduced surcharge",
			ForgeKey{Extra: 9, ReducedBy: CardsInHand{Player: Controller, House: AnyHouse}},
			"forge a key at +9 Æmber current cost, reduced by 1 Æmber for each card in your hand",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.effect.Text(); got != tc.want {
				t.Errorf("text = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestForgeKeyValidate(t *testing.T) {
	if err := (ForgeKey{ReducedBy: CardsInHand{}}).validate(); err == nil {
		t.Error("a reduction with no Extra should not validate")
	}
	if err := (ForgeKey{FreeOfCost: true, Extra: 2}).validate(); err == nil {
		t.Error("a free forge with an Extra should not validate")
	}
	if err := (ForgeKey{Extra: 6}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

func TestForgeKeyPaysTheSurcharge(t *testing.T) {
	g := started(t)
	g.State.Aember[0] = KeyCost + 2
	ctx := &EffectContext{Resolver: g, Controller: 0}

	ForgeKey{Extra: 6}.Resolve(ctx)
	if g.State.Keys[0] != 0 {
		t.Fatalf("keys = %d, want 0 — the surcharge is unaffordable", g.State.Keys[0])
	}

	ForgeKey{Extra: 2}.Resolve(ctx)
	if g.State.Keys[0] != 1 {
		t.Errorf("keys = %d, want 1", g.State.Keys[0])
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("Æmber = %d, want 0", g.State.Aember[0])
	}
}

func TestForgeKeyReducedBelowTheSurcharge(t *testing.T) {
	g := started(t)
	g.State.Aember[0] = KeyCost
	for i := 0; i < 12; i++ {
		g.AddToHand(NewCard("Filler", Brobnar, Tactic, Common), 0)
	}

	// A +9 surcharge reduced by 12 cards in hand cannot drop below the key cost.
	ForgeKey{
		Extra:     9,
		ReducedBy: CardsInHand{Player: Controller, House: AnyHouse},
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Keys[0] != 1 {
		t.Errorf("keys = %d, want 1", g.State.Keys[0])
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("Æmber = %d, want 0 — the cost floor is the key cost", g.State.Aember[0])
	}
}

func TestRaiseKeyCost(t *testing.T) {
	effect := RaiseKeyCost{Player: Opponent, Amount: 3, Duration: NextTurn}
	if got := effect.Text(); got != "keys cost +3 Æmber during your opponent's next turn" {
		t.Errorf("text = %q", got)
	}
	if got := (RaiseKeyCost{Player: Controller, Amount: 1, Duration: NextTurn}).Text(); got !=
		"keys cost +1 Æmber during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if err := (RaiseKeyCost{Amount: 3, Duration: NextTurn}).validate(); err == nil {
		t.Error("an unset player should not validate")
	}
	if err := (RaiseKeyCost{Player: Opponent}).validate(); err == nil {
		t.Error("a zero raise should not validate")
	}
	if err := effect.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (RaiseKeyCost{Player: Opponent, Amount: 3}).validate(); err == nil {
		t.Error("an unset duration should not validate")
	}
	if err := (RaiseKeyCost{
		Player:   Opponent,
		Amount:   3,
		Duration: Forever,
	}).validate(); err == nil {
		t.Error("a duration the surcharge cannot express should not validate")
	}
}

// TestRaiseKeyCostThisTurn checks the EndOfTurn duration bites immediately and
// lifts when the current turn ends, rather than waiting for a turn boundary.
func TestRaiseKeyCostThisTurn(t *testing.T) {
	effect := RaiseKeyCost{Player: Controller, Amount: 2, Duration: EndOfTurn}
	if got := effect.Text(); got != "your keys cost +2 Æmber for the remainder of the turn" {
		t.Errorf("text = %q", got)
	}
	if got := (RaiseKeyCost{
		Player:   Opponent,
		Amount:   2,
		Duration: EndOfTurn,
	}).Text(); got != "your opponent's keys cost +2 Æmber for the remainder of the turn" {
		t.Errorf("opponent text = %q", got)
	}

	g := started(t)
	source := g.AddToBattleline(NewCard("Lash", Dis, Creature, Common, WithPower(1)), 0)
	effect.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})

	if got := g.CurrentKeyCost(0); got != KeyCost+2 {
		t.Errorf("key cost = %d, want %d right away", got, KeyCost+2)
	}
	g.EndPlayPhase(0)
	if got := g.CurrentKeyCost(0); got != KeyCost {
		t.Errorf("key cost = %d, want %d after the turn ends", got, KeyCost)
	}
}

func TestRaiseKeyCostLandsOnTheNextTurn(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(NewCard("Lash", Dis, Creature, Common, WithPower(1)), 0)

	RaiseKeyCost{Player: Opponent, Amount: 3, Duration: NextTurn}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})

	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost = %d, want %d before the raise lands", got, KeyCost)
	}
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if got := g.CurrentKeyCost(1); got != KeyCost+3 {
		t.Errorf("key cost = %d, want %d", got, KeyCost+3)
	}
	if reasons := g.RestrictionSources(1); len(reasons) == 0 {
		t.Error("the raise should name its source as a turn restriction")
	}
	g.EndPlayPhase(1)
	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost = %d, want %d after the turn ends", got, KeyCost)
	}
}

func TestRaiseKeyCostStacks(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	RaiseKeyCost{Player: Opponent, Amount: 3, Duration: NextTurn}.Resolve(ctx)
	RaiseKeyCost{Player: Opponent, Amount: 3, Duration: NextTurn}.Resolve(ctx)
	if got := g.State.KeyCostBumpNext[1].Value; got != 6 {
		t.Errorf("armed raise = %d, want 6", got)
	}
}

func TestConditionalElse(t *testing.T) {
	effect := Conditional{
		Cond: OpponentAember{Is: Exactly},
		Then: ForgeKey{Extra: 2},
		Else: ForgeKey{Extra: 6},
	}
	want := "if your opponent has no Æmber, forge a key at +2 Æmber current cost. " +
		"Otherwise, forge a key at +6 Æmber current cost"
	if got := effect.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if err := effect.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (Conditional{
		Cond: OpponentAember{Is: Exactly},
		Then: ForgeKey{},
		Else: ForgeKey{FreeOfCost: true, Extra: 1},
	}).validate(); err == nil {
		t.Error("an invalid Else should not validate")
	}

	g := started(t)
	g.State.Aember[0] = 100
	g.State.Aember[1] = 1
	effect.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if got := g.State.Aember[0]; got != 100-(KeyCost+6) {
		t.Errorf("Æmber = %d, want the Else branch's +6 cost paid", got)
	}
}

func TestCardsInHandAnyHouse(t *testing.T) {
	g := started(t)
	g.AddToHand(NewCard("A", Brobnar, Tactic, Common), 0)
	g.AddToHand(NewCard("B", Logos, Tactic, Common), 0)
	g.AddToHand(NewCard("C", Dis, Tactic, Common), 1)

	mine := CardsInHand{Player: Controller, House: AnyHouse}
	theirs := CardsInHand{Player: Opponent, House: AnyHouse}
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := mine.Value(ctx); got != 2 {
		t.Errorf("controller hand = %d, want 2", got)
	}
	if got := theirs.Value(ctx); got != 1 {
		t.Errorf("opponent hand = %d, want 1", got)
	}
	if got := mine.CountText(); got != "card in your hand" {
		t.Errorf("controller text = %q", got)
	}
	if got := theirs.CountText(); got != "card in your opponent's hand" {
		t.Errorf("opponent text = %q", got)
	}
}

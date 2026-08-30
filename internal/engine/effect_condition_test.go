package engine

import "testing"

func TestConditionalEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	atLeast := Conditional{Cond: OpponentAemberAtLeast{Amount: 7}, Then: LoseAember{Player: Opponent, Amount: 4}}
	if atLeast.Text() != "if your opponent has 7 Æmber or more, your opponent loses 4 Æmber" {
		t.Errorf("at-least text = %q", atLeast.Text())
	}
	g.State.Aember[1] = 5
	atLeast.Resolve(ctx) // condition not met (5 < 7)
	if g.State.Aember[1] != 5 {
		t.Errorf("unmet condition should be a no-op; opp = %d", g.State.Aember[1])
	}
	g.State.Aember[1] = 8
	atLeast.Resolve(ctx) // met
	if g.State.Aember[1] != 4 {
		t.Errorf("met condition should apply; opp = %d, want 4", g.State.Aember[1])
	}

	exact := Conditional{Cond: OpponentAemberExactly{Amount: 1}, Then: StealAember{Amount: 1}}
	if exact.Text() != "if your opponent has exactly 1 Æmber, steal 1 Æmber" {
		t.Errorf("exact text = %q", exact.Text())
	}
	g.State.Aember[0], g.State.Aember[1] = 0, 2
	exact.Resolve(ctx) // 2 != 1, no-op
	if g.State.Aember[0] != 0 {
		t.Errorf("exact unmet should be a no-op; you = %d", g.State.Aember[0])
	}
	g.State.Aember[1] = 1
	exact.Resolve(ctx) // met
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 0 {
		t.Errorf("exact met: you=%d opp=%d, want 1/0", g.State.Aember[0], g.State.Aember[1])
	}
}

func TestControlsMoreCreatures(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	c := ControlsMoreCreatures{}
	if c.CondText() != "if you control more creatures than your opponent" {
		t.Errorf("CondText = %q", c.CondText())
	}
	g.AddToBattleline(testCreature("mine", 1), 0)
	if !c.Met(ctx) {
		t.Error("1 vs 0 should be met")
	}
	g.AddToBattleline(testCreature("theirs", 1), 1)
	if c.Met(ctx) {
		t.Error("1 vs 1 should not be met")
	}
}

func TestCardsPlayed(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	cond := CardsPlayed{Player: Controller, House: Sanctum, Amount: 7}
	if cond.CondText() != "if you have played 7 or more Sanctum cards this turn" {
		t.Errorf("cond text = %q", cond.CondText())
	}
	if cond.CountText() != "Sanctum card you have played this turn" {
		t.Errorf("count text = %q", cond.CountText())
	}
	g.State.CardsPlayedByHouseThisTurn[0][Sanctum] = 6
	if cond.Met(ctx) || cond.Value(ctx) != 6 {
		t.Errorf("six Sanctum cards: met=%v value=%d, want false/6", cond.Met(ctx), cond.Value(ctx))
	}
	g.State.CardsPlayedByHouseThisTurn[0][Sanctum] = 7
	if !cond.Met(ctx) {
		t.Error("seven Sanctum cards should satisfy the condition")
	}
	// The default threshold is one played card.
	if (CardsPlayed{Player: Controller, House: Mars}).Met(ctx) {
		t.Error("no Mars cards played should not meet the default threshold")
	}
}

func TestRepeatWhile(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0], g.State.Aember[1] = 0, 5 // opponent leads
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := RepeatWhile{Cond: OpponentAemberMoreThanYou{}, Do: StealAember{Amount: 1}}
	if e.Text() != "if your opponent has more Æmber than you, steal 1 Æmber -> repeat this effect" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	// 5/0 -> 4/1 -> 3/2 -> 2/3 (opponent no longer leads): 3 steals.
	if g.Aember(0) != 3 || g.Aember(1) != 2 {
		t.Errorf("after repeat: you=%d opp=%d, want 3/2", g.Aember(0), g.Aember(1))
	}

	// Condition false from the start: the loop never runs.
	g2 := NewGame("A", "B", 1)
	g2.State.Aember[0], g2.State.Aember[1] = 3, 3
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Aember(0) != 3 || g2.Aember(1) != 3 {
		t.Errorf("equal pools should not steal: you=%d opp=%d", g2.Aember(0), g2.Aember(1))
	}

	if err := validateEffect(RepeatWhile{Do: StealAember{Amount: 1}}); err != nil {
		t.Errorf("validate = %v", err)
	}
}

func TestMayRepeat(t *testing.T) {
	e := MayRepeat{Cond: OpponentAemberMoreThanYou{}, Do: StealAember{Amount: 1}}
	if got := e.Text(); got != "steal 1 Æmber -> if your opponent has more Æmber than you, you may repeat this effect" {
		t.Errorf("text = %q", got)
	}

	// Accepting every offer repeats until the condition fails.
	g := NewGame("A", "B", 1)
	g.State.Aember[0], g.State.Aember[1] = 0, 5 // opponent leads
	g.SetChooser(0, optionPicker{idx: 0})       // always "Yes"
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	// 5/0 -> 4/1 -> 3/2 -> 2/3 (opponent no longer leads): Do once plus 2 repeats.
	if g.Aember(0) != 3 || g.Aember(1) != 2 {
		t.Errorf("accepted: you=%d opp=%d, want 3/2", g.Aember(0), g.Aember(1))
	}

	// Declining the first offer runs Do exactly once even though the condition holds.
	g2 := NewGame("A", "B", 1)
	g2.State.Aember[0], g2.State.Aember[1] = 0, 5
	g2.SetChooser(0, optionPicker{idx: 1}) // "No"
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Aember(0) != 1 || g2.Aember(1) != 4 {
		t.Errorf("declined: you=%d opp=%d, want 1/4", g2.Aember(0), g2.Aember(1))
	}

	if err := validateEffect(MayRepeat{Cond: InPlay{Player: Controller, Type: Creature}, Do: StealAember{Amount: 1}}); err != nil {
		t.Errorf("validate = %v", err)
	}
}

func TestItIs(t *testing.T) {
	// CondText renders the filtered noun with the right article.
	cases := map[string]ItIs{
		"if it is a Mars creature": {House: Mars, Type: Creature},
		"if it is an artifact":     {Type: Artifact},
		"if it is a Mars card":     {House: Mars},
		"if it is a card":          {},
	}
	for want, e := range cases {
		if got := e.CondText(); got != want {
			t.Errorf("CondText() = %q, want %q", got, want)
		}
	}

	g := NewGame("Alice", "Bob", 1)
	mars := g.AddToDeck(NewCard("Martian", Mars, Creature, Common, WithPower(1)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// No card in context is never met.
	if (ItIs{Type: Creature}).Met(ctx) {
		t.Error("Met with no context card should be false")
	}
	ctx.It, ctx.HasIt = mars, true
	if !(ItIs{House: Mars, Type: Creature}).Met(ctx) {
		t.Error("a Mars creature should match a Mars-creature filter")
	}
	if (ItIs{House: Logos}).Met(ctx) {
		t.Error("a Mars creature should not match a Logos filter")
	}
	if (ItIs{Type: Artifact}).Met(ctx) {
		t.Error("a creature should not match an artifact filter")
	}
}

package engine

import "testing"

func TestInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("m1", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("b1", Brobnar, Creature, Common, WithPower(2)), 0)
	g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	g.AddToBattleline(NewCard("foe", Shadows, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Value across type and house filters, and the opposing side.
	values := []struct {
		name string
		in   InPlay
		want int
	}{
		{"friendly creatures", InPlay{Player: Controller, Type: Creature}, 3},
		{"friendly Mars creatures", InPlay{Player: Controller, Type: Creature, House: Mars}, 2},
		{"friendly artifacts", InPlay{Player: Controller, Type: Artifact}, 1},
		{"friendly cards, any type", InPlay{Player: Controller}, 4},
		{"enemy creatures", InPlay{Player: Opponent, Type: Creature}, 1},
	}
	for _, tc := range values {
		if got := tc.in.Value(ctx); got != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.name, got, tc.want)
		}
	}

	// CountText — the singular "for each" noun.
	texts := []struct {
		in   InPlay
		want string
	}{
		{InPlay{Player: Controller, Type: Creature}, "friendly creature in play"},
		{InPlay{Player: Controller, Type: Creature, House: Mars}, "friendly Mars creature"},
		{InPlay{Player: Opponent, Type: Creature}, "enemy creature in play"},
		{InPlay{Player: Controller, Type: Artifact}, "friendly artifact in play"},
		{InPlay{Player: Controller}, "friendly card in play"},
	}
	for _, tc := range texts {
		if got := tc.in.CountText(); got != tc.want {
			t.Errorf("CountText = %q, want %q", got, tc.want)
		}
	}

	// CondText — singular and plural.
	if got := (InPlay{Player: Controller, Type: Creature}).CondText(); got != "if there is a friendly creature in play" {
		t.Errorf("singular CondText = %q", got)
	}
	if got := (InPlay{Player: Controller, Type: Creature, Amount: 2}).CondText(); got != "if there are 2 friendly creatures in play" {
		t.Errorf("plural CondText = %q", got)
	}

	// Met — Amount defaults to one; a higher threshold may not be reached.
	if !(InPlay{Player: Controller, Type: Creature}).Met(ctx) {
		t.Error("default threshold should be met with 3 creatures")
	}
	if !(InPlay{Player: Controller, Type: Creature, Amount: 3}).Met(ctx) {
		t.Error("threshold 3 should be met with 3 creatures")
	}
	if (InPlay{Player: Controller, Type: Creature, Amount: 4}).Met(ctx) {
		t.Error("threshold 4 should not be met with 3 creatures")
	}
}

func TestInPlayEachPlayer(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("f", 5), 0)
	g.AddToBattleline(testCreature("e", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// EachPlayer counts both players' matching cards, with no friendly/enemy
	// qualifier in the rendered noun.
	byHouse := InPlay{Player: EachPlayer, Type: Creature, House: Brobnar}
	if got := byHouse.CountText(); got != "Brobnar creature in play" {
		t.Errorf("count text = %q, want %q", got, "Brobnar creature in play")
	}
	if got := byHouse.Value(ctx); got != 2 {
		t.Errorf("value = %d, want 2 (one creature per player)", got)
	}

	if got := (InPlay{Player: EachPlayer, Type: Creature}).CountText(); got != "creature in play" {
		t.Errorf("no-house count text = %q, want %q", got, "creature in play")
	}
}

func TestInPlayReady(t *testing.T) {
	g := NewGame("A", "B", 1)
	spent := g.AddToBattleline(marsCreature("spent", 3), 0)
	g.AddToBattleline(marsCreature("fresh", 3), 0)
	g.SetExhausted(spent, true)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	ready := InPlay{Player: Controller, Type: Creature, House: Mars, Ready: true}
	if got := ready.CountText(); got != "friendly ready Mars creature" {
		t.Errorf("count text = %q, want %q", got, "friendly ready Mars creature")
	}
	if got := ready.Value(ctx); got != 1 {
		t.Errorf("value = %d, want 1 (the exhausted creature does not count)", got)
	}
}

func TestCardsDestroyedCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ctx.Produced.Destroyed = [2]int{2, 0}

	c := CardsDestroyed{}
	if got := c.CountText(); got != "card destroyed this way" {
		t.Errorf("count text = %q", got)
	}
	if got := c.Value(ctx); got != 2 {
		t.Errorf("value = %d, want 2", got)
	}
}

func TestEachFriendlyArtifactTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddArtifact(NewCard("mine", Brobnar, Artifact, Common), 0)
	g.AddArtifact(NewCard("theirs", Brobnar, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	target := Target{Kind: TargetEachFriendlyArtifact}
	if got := target.Text(); got != "each friendly artifact" {
		t.Errorf("text = %q", got)
	}
	got := target.Select(ctx)
	if len(got) != 1 || got[0] != mine {
		t.Errorf("selected = %v, want [%v]", got, mine)
	}
}

func TestExcessCreatures(t *testing.T) {
	opponents := ExcessCreatures{Player: Opponent}
	mine := ExcessCreatures{Player: Controller}
	if opponents.CountText() != "creature your opponent controls in excess of you" {
		t.Errorf("count text = %q", opponents.CountText())
	}
	if mine.CountText() != "creature you have in excess of your opponent" {
		t.Errorf("count text = %q", mine.CountText())
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 3), 1)
	g.AddToBattleline(testCreature("o2", 3), 1)
	g.AddToBattleline(testCreature("m1", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if got := opponents.Value(ctx); got != 1 {
		t.Errorf("excess = %d, want 1 (opponent 2, you 1)", got)
	}
	if got := mine.Value(ctx); got != 0 {
		t.Errorf("own excess = %d, want 0 when your opponent controls more", got)
	}

	// When the controller controls at least as many, the opponent's excess floors
	// at zero and the controller's own excess starts counting up.
	g.AddToBattleline(testCreature("m2", 3), 0)
	g.AddToBattleline(testCreature("m3", 3), 0)
	if got := opponents.Value(ctx); got != 0 {
		t.Errorf("excess = %d, want 0 when you control more", got)
	}
	if got := mine.Value(ctx); got != 1 {
		t.Errorf("own excess = %d, want 1 (you 3, opponent 2)", got)
	}
}

func TestInPlayByName(t *testing.T) {
	g := NewGame("A", "B", 1)
	bear := g.Register(NewCard("Ancient Bear", Untamed, Creature, Common, WithPower(6)), 0)
	g.AddToBattleline(testCreature("other", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	none := InPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear", None: true}
	if want := "if there are no Ancient Bears in play"; none.CondText() != want {
		t.Errorf("cond text = %q, want %q", none.CondText(), want)
	}
	if !none.Met(ctx) {
		t.Error("None should be met while no Ancient Bear is in play")
	}

	g.State.Battleline[0].add(bear)
	if none.Met(ctx) {
		t.Error("None should not be met once an Ancient Bear is in play")
	}
	some := InPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear"}
	if n := some.Value(ctx); n != 1 {
		t.Errorf("value = %d, want 1 (only the named card counts)", n)
	}
	if want := "if there is an Ancient Bear in play"; some.CondText() != want {
		t.Errorf("cond text = %q, want %q", some.CondText(), want)
	}
	two := InPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear", Amount: 2}
	if want := "if there are 2 Ancient Bears in play"; two.CondText() != want {
		t.Errorf("cond text = %q, want %q", two.CondText(), want)
	}
}

package engine

import "testing"

// TestPurgeFromHandEachText covers the rendered phrase for each filter
// combination of an Each selection.
func TestPurgeFromHandEachText(t *testing.T) {
	cases := []struct {
		e    PurgeFromHand
		want string
	}{
		{
			PurgeFromHand{Player: Controller, Selection: Each{}},
			"purge each card from your hand",
		},
		{
			PurgeFromHand{Player: Opponent, Selection: Each{Type: Creature}},
			"purge each creature from your opponent's hand",
		},
		{
			PurgeFromHand{Player: Controller, Selection: Each{Type: Creature, ExceptHouse: Mars}},
			"purge each non-Mars creature from your hand",
		},
	}
	for _, c := range cases {
		if got := c.e.Text(); got != c.want {
			t.Errorf("Text = %q, want %q", got, c.want)
		}
		if err := c.e.validate(); err != nil {
			t.Errorf("validate = %v, want nil", err)
		}
	}
}

// TestPurgeFromHandEachPurgesEveryMatch checks the filters decide what goes,
// that nothing else is touched, and that the tally feeds a following count.
func TestPurgeFromHandEachPurgesEveryMatch(t *testing.T) {
	g := NewGame("A", "B", 1)
	martian := g.AddToHand(NewCard("martian", Mars, Creature, Common, WithPower(3)), 0)
	brobnar := g.AddToHand(NewCard("brobnar", Brobnar, Creature, Common, WithPower(3)), 0)
	untamed := g.AddToHand(NewCard("untamed", Untamed, Creature, Common, WithPower(3)), 0)
	tactic := g.AddToHand(NewCard("tactic", Shadows, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	PurgeFromHand{Player: Controller, Selection: Each{Type: Creature, ExceptHouse: Mars}}.
		Resolve(ctx)

	if got := g.Purge(0); len(got) != 2 {
		t.Errorf("purged = %v, want brobnar and untamed", got)
	}
	for _, id := range []LocalID{brobnar, untamed} {
		if containsID(g.Hand(0), id) {
			t.Errorf("%s should have been purged", g.Name(id))
		}
	}
	for _, id := range []LocalID{martian, tactic} {
		if !containsID(g.Hand(0), id) {
			t.Errorf("%s should have been spared", g.Name(id))
		}
	}
	if got := (CardsPurged{}).Value(ctx); got != 2 {
		t.Errorf("tally = %d, want 2", got)
	}
}

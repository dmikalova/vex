package engine

import "testing"

// TestPurgeEachFromDiscardText covers the rendered phrase for each filter
// combination.
func TestPurgeEachFromDiscardText(t *testing.T) {
	cases := []struct {
		e    PurgeEachFromDiscard
		want string
	}{
		{
			PurgeEachFromDiscard{},
			"purge each card from each player's discard pile",
		},
		{
			PurgeEachFromDiscard{House: Untamed, Type: Creature},
			"purge each Untamed creature from each player's discard pile",
		},
		{
			PurgeEachFromDiscard{House: Untamed, Type: Creature, GainOwnerAember: true},
			"purge each Untamed creature from each player's discard pile. " +
				"For each card purged this way, its owner gains 1 Æmber",
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

// TestPurgeEachFromDiscardPurgesBothPiles checks the filters decide what goes
// from both players' discard piles and that each owner gains 1 Æmber per card.
func TestPurgeEachFromDiscardPurgesBothPiles(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.Register(NewCard("mine", Untamed, Creature, Common, WithPower(3)), 0)
	g.State.Discard[0].add(mine)
	myBrobnar := g.Register(NewCard("myBro", Brobnar, Creature, Common, WithPower(3)), 0)
	g.State.Discard[0].add(myBrobnar)
	theirs := g.Register(NewCard("theirs", Untamed, Creature, Common, WithPower(3)), 1)
	g.State.Discard[1].add(theirs)
	theirTactic := g.Register(NewCard("theirTac", Untamed, Tactic, Common), 1)
	g.State.Discard[1].add(theirTactic)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	PurgeEachFromDiscard{House: Untamed, Type: Creature, GainOwnerAember: true}.
		Resolve(ctx)

	if !g.State.Purge[0].contains(mine) || !g.State.Purge[1].contains(theirs) {
		t.Error("both Untamed creatures should have been purged")
	}
	if !g.State.Discard[0].contains(myBrobnar) {
		t.Error("the Brobnar creature should have been spared")
	}
	if !g.State.Discard[1].contains(theirTactic) {
		t.Error("the Untamed tactic should have been spared")
	}
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 1 {
		t.Errorf("each owner should gain 1, got %d and %d",
			g.State.Aember[0], g.State.Aember[1])
	}
}

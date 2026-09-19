package engine

import "testing"

// TestTurnHistoryRollover covers the tallies the engine keeps about what a player
// did during a turn, and the rollover that turns "this turn" into "their previous
// turn" when their turn ends.
func TestTurnHistoryRollover(t *testing.T) {
	g := NewGame("A", "B", 1)
	if got := g.TurnHistory(0, KeysForgedThisTurn); got != 0 {
		t.Errorf("fresh game keys forged = %d, want 0", got)
	}

	g.SetAember(0, 6)
	g.forgeKey(0)
	if got := g.TurnHistory(0, KeysForgedThisTurn); got != 1 {
		t.Errorf("keys forged this turn = %d, want 1", got)
	}

	g.EndPlayPhase(0)
	if got := g.TurnHistory(0, KeysForgedLastTurn); got != 1 {
		t.Errorf("keys forged last turn = %d, want 1", got)
	}
	if got := g.TurnHistory(0, KeysForgedThisTurn); got != 0 {
		t.Errorf("keys forged this turn should reset, got %d", got)
	}
}

// TestForgedKeyCondition covers both windows the condition can ask about and both
// subjects it can render.
func TestForgedKeyCondition(t *testing.T) {
	if err := (ForgedKey{}).validate(); err == nil {
		t.Error("unset player should be invalid")
	}
	if err := (ForgedKey{Player: Controller}).validate(); err != nil {
		t.Errorf("validate: %v", err)
	}

	mine := ForgedKey{Player: Controller}
	theirs := ForgedKey{
		Player:   Opponent,
		Previous: true,
	}
	notMine := Not{Cond: ForgedKey{Player: Controller}}
	if got := mine.CondText(); got != "if you forged a key this turn" {
		t.Errorf("CondText = %q", got)
	}
	if got := theirs.CondText(); got != "if your opponent forged a key during their previous turn" {
		t.Errorf("CondText = %q", got)
	}
	if got := notMine.CondText(); got != "if you have not forged a key this turn" {
		t.Errorf("CondText = %q", got)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	if mine.Met(ctx) || theirs.Met(ctx) {
		t.Error("nothing forged yet, both conditions should be unmet")
	}
	if !notMine.Met(ctx) {
		t.Error("no key forged yet, the negated condition should be met")
	}

	g.SetAember(0, 6)
	g.forgeKey(0)
	if !mine.Met(ctx) {
		t.Error("a key forged this turn should meet the condition")
	}
	if notMine.Met(ctx) {
		t.Error("a key forged this turn should not meet the negated condition")
	}

	g.SetAember(1, 6)
	g.forgeKey(1)
	g.EndPlayPhase(1)
	if !theirs.Met(ctx) {
		t.Error("the opponent's key from their last turn should meet the condition")
	}
}

// TestAemberStolenFromYouCondition covers the tally a steal bumps, its rollover
// into the victim's "previous turn" window, and the condition that reads it.
func TestAemberStolenFromYouCondition(t *testing.T) {
	c := AemberStolenFromYou{}
	if got := c.CondText(); got != "if your opponent stole Æmber from you on their previous turn" {
		t.Errorf("CondText = %q", got)
	}

	g := NewGame("A", "B", 1)
	victim := &EffectContext{
		Resolver:   g,
		Controller: 1,
	}
	if c.Met(victim) {
		t.Error("nothing stolen yet, condition should be unmet")
	}

	// Player 0 steals from player 1 during player 0's turn.
	g.SetAember(1, 3)
	StealAember{Amount: 2}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if got := g.TurnHistory(1, AemberStolenFromThisTurn); got != 2 {
		t.Fatalf("stolen-from this turn = %d, want 2", got)
	}
	if c.Met(victim) {
		t.Error("the theft is only this turn, not yet the previous turn")
	}

	g.EndPlayPhase(0) // player 0's turn ends, arming player 1's previous-turn window
	if got := g.TurnHistory(1, AemberStolenFromLastTurn); got != 2 {
		t.Fatalf("stolen-from last turn = %d, want 2", got)
	}
	if !c.Met(victim) {
		t.Error("player 1 was robbed on player 0's previous turn, condition should be met")
	}
}

// TestCreatureDestroyedThisTurnCondition covers both sides of the merged
// condition, plus the validation that keeps it to the two sides it can read.
func TestCreatureDestroyedThisTurnCondition(t *testing.T) {
	for _, tc := range []struct {
		name string
		cond CreatureDestroyedThisTurn
		want string
		stat TurnStat
	}{
		{
			"enemy", CreatureDestroyedThisTurn{Player: Opponent},
			"if an enemy creature has been destroyed this turn", EnemyCreaturesDestroyed,
		},
		{
			"friendly", CreatureDestroyedThisTurn{Player: Controller},
			"if a friendly creature has been destroyed this turn", FriendlyCreaturesDestroyed,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cond.CondText(); got != tc.want {
				t.Errorf("CondText = %q, want %q", got, tc.want)
			}
			if err := tc.cond.validate(); err != nil {
				t.Errorf("valid side rejected: %v", err)
			}

			g := NewGame("A", "B", 1)
			ctx := &EffectContext{
				Resolver:   g,
				Controller: 0,
			}
			if tc.cond.Met(ctx) {
				t.Error("no creature destroyed yet, condition should be unmet")
			}

			g.State.TurnHistory[0][tc.stat] = 1
			if !tc.cond.Met(ctx) {
				t.Error("a creature destroyed this turn should meet the condition")
			}
		})
	}

	if err := (CreatureDestroyedThisTurn{}).validate(); err == nil {
		t.Error("an unset side should be rejected")
	}
}

// TestTurnCount covers the shared count over a turn-history tally, in both the
// "for each" and the "if" rendering.
func TestTurnCount(t *testing.T) {
	c := TurnCount{
		Player: Controller,
		Of:     EnemyCreaturesFightKilled,
	}
	want := "enemy creature that was destroyed in a fight this turn"
	if got := c.CountText(); got != want {
		t.Errorf("CountText = %q, want %q", got, want)
	}

	played := TurnCount{
		Player: Opponent,
		Of:     CreaturesPlayedLastTurn,
	}
	if got := played.CountClause("3 or more", true); got !=
		"your opponent played 3 or more creatures on their previous turn" {
		t.Errorf("CountClause = %q", got)
	}
	mine := TurnCount{
		Player: Controller,
		Of:     CreaturesPlayedLastTurn,
	}
	if got := mine.CountClause("exactly 1", false); got !=
		"you played exactly 1 creature on your previous turn" {
		t.Errorf("CountClause = %q", got)
	}

	destroyed := TurnCount{
		Player: Controller,
		Of:     EnemyCreaturesDestroyed,
	}
	if got := destroyed.CountClause("3 or more", true); got !=
		"3 or more enemy creatures have been destroyed this turn" {
		t.Errorf("CountClause(destroyed, plural) = %q", got)
	}
	if got := destroyed.CountClause("exactly 1", false); got !=
		"exactly 1 enemy creature has been destroyed this turn" {
		t.Errorf("CountClause(destroyed, singular) = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.State.TurnHistory[0][EnemyCreaturesFightKilled] = 2
	if got := c.Value(&EffectContext{
		Resolver:   g,
		Controller: 0,
	}); got != 2 {
		t.Errorf("Value = %d, want 2", got)
	}
}

// TestUnforgeKey covers taking a forged key back off a player, including the case
// where there is no key to take.
func TestUnforgeKey(t *testing.T) {
	if err := (UnforgeKey{}).validate(); err == nil {
		t.Error("unset player should be invalid")
	}
	if err := (UnforgeKey{Player: Opponent}).validate(); err != nil {
		t.Errorf("validate: %v", err)
	}
	if got := (UnforgeKey{Player: Opponent}).Text(); got != "unforge one of your opponent's keys" {
		t.Errorf("Text = %q", got)
	}
	if got := (UnforgeKey{Player: Controller}).Text(); got != "unforge one of your keys" {
		t.Errorf("Text = %q", got)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	if (UnforgeKey{Player: Opponent}).resolveGate(ctx) {
		t.Error("unforging with no keys should report false")
	}
	if g.Keys(1) != 0 {
		t.Error("unforging with no keys should do nothing")
	}

	g.SetAember(1, 6)
	g.forgeKey(1)
	if !(UnforgeKey{Player: Opponent}).resolveGate(ctx) {
		t.Error("unforging a forged key should report true")
	}
	if g.Keys(1) != 0 {
		t.Errorf("keys = %d, want 0", g.Keys(1))
	}

	g.SetAember(1, 6)
	g.forgeKey(1)
	UnforgeKey{Player: Opponent}.Resolve(ctx)
	if g.Keys(1) != 0 {
		t.Errorf("keys after Resolve = %d, want 0", g.Keys(1))
	}
}

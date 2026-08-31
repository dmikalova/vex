package engine

import "testing"

func TestPlayTopOfDeckEffect(t *testing.T) {
	g := started(t)
	top := g.AddToDeck(testCreature("Rando", 2), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := (PlayTopOfDeck{}).Text(); got != "play the top card of your deck" {
		t.Errorf("text = %q", got)
	}

	(PlayTopOfDeck{}).Resolve(ctx)
	if got := g.Battleline(0); len(got) != 1 || got[0] != top {
		t.Errorf("battleline = %v, want [%d]", got, top)
	}
	if g.State.Deck[0].Count != 0 {
		t.Errorf("deck count = %d, want 0", g.State.Deck[0].Count)
	}

	// An empty deck is a no-op.
	(PlayTopOfDeck{}).Resolve(ctx)
	if len(g.Battleline(0)) != 1 {
		t.Error("playing from an empty deck should do nothing")
	}
}

func TestEndOfTurnTriggerFires(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.BeginTurn(0)
	g.State.Aember[1] = 3
	g.AddToBattleline(NewCard("Shaffles", Dis, Creature, Common, WithPower(2),
		WithAbility(TriggerEndOfTurn, LoseAember{Player: Opponent, Amount: 1})), 0)

	g.EndTurn(0)

	if g.State.Aember[1] != 2 {
		t.Errorf("opponent Æmber = %d, want 2 after end-of-turn drain", g.State.Aember[1])
	}
}

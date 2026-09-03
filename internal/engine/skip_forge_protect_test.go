package engine

import (
	"strings"
	"testing"
)

func TestSkipForgeStep(t *testing.T) {
	if got := (SkipForgeStep{Player: Opponent}).Text(); got != `your opponent skips the "forge a key" step during their next turn` {
		t.Errorf("opponent text = %q", got)
	}
	if got := (SkipForgeStep{Player: Controller}).Text(); got != `you skip the "forge a key" step during your next turn` {
		t.Errorf("self text = %q", got)
	}
	if (SkipForgeStep{}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	if (SkipForgeStep{Player: Opponent}).validate() != nil {
		t.Error("a set player should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	g.State.Aember[1] = 6 // enough for player 1 to forge on their turn
	SkipForgeStep{Player: Opponent}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if !g.State.SkipForgeNext[1].Value {
		t.Fatal("the skip should arm the opponent's next turn")
	}
	g.EndPlayPhase(0)

	keysBefore := g.State.Keys[1]
	g.StartTurn(1)
	if g.State.SkipForgeNext[1].Value {
		t.Error("the skip should be consumed when the turn begins")
	}
	if g.State.Keys[1] != keysBefore {
		t.Error("the opponent should not have forged a key on the skipped turn")
	}
}

func TestAemberProtection(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[1] = 3
	g.AddToBattleline(
		NewCard("keeper", Sanctum, Creature, Rare, WithPower(4), WithAemberTheftImmunity()),
		1,
	)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (StealAember{Amount: 2}).resolveGate(ctx) {
		t.Error("stealing from a protected pool should report no movement")
	}
	if g.State.Aember[1] != 3 || g.State.Aember[0] != 0 {
		t.Errorf("Æmber = %d/%d, want 3/0 (nothing stolen)", g.State.Aember[0], g.State.Aember[1])
	}

	def := NewCard("keeper", Sanctum, Creature, Rare, WithPower(4), WithAemberTheftImmunity())
	if !strings.Contains(RenderCardRules(&def), "Your Æmber cannot be stolen.") {
		t.Error("card rules should render the theft-immunity line")
	}
}

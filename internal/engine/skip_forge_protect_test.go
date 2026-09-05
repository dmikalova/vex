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

func TestSkipsForgeConstant(t *testing.T) {
	g := NewGame("A", "B", 1)
	if g.skipsForge(0) {
		t.Error("a player with nothing in play should not skip forging")
	}
	def := NewCard("The Sting", Shadows, Artifact, Rare,
		WithRestrictions(Restrictions{SkipForge: true}))
	sting := g.AddArtifact(def, 0)
	if !g.skipsForge(0) {
		t.Error("a controller of a SkipForge card should skip forging")
	}
	if g.skipsForge(1) {
		t.Error("the opponent should be unaffected")
	}
	if !strings.Contains(RenderCardRules(&def), `You skip your "forge a key" step.`) {
		t.Error("card rules should render the skip-forge line")
	}

	g.State.Aember[0] = 6
	keysBefore := g.State.Keys[0]
	g.forgePhase(0)
	if g.State.Keys[0] != keysBefore {
		t.Error("a player skipping their forge step should not forge a key")
	}
	if g.State.Aember[0] != 6 {
		t.Error("skipping the forge step should not spend any Æmber")
	}
	_ = sting
}

func TestForgeAemberGain(t *testing.T) {
	if got := gainsForgeAemberText(&CardDefinition{}); got != "" {
		t.Errorf("gainsForgeAemberText on a plain card = %q, want empty", got)
	}
	def := NewCard("The Sting", Shadows, Artifact, Rare, WithGainsForgeAember())
	if !def.GainsForgeAember {
		t.Error("WithGainsForgeAember should set the flag")
	}
	if got := gainsForgeAemberText(
		&def,
	); got != "You gain all Æmber your opponent spends when forging a key." {
		t.Errorf("gainsForgeAemberText = %q", got)
	}
	if !strings.Contains(
		RenderCardRules(&def), "You gain all Æmber your opponent spends when forging a key.",
	) {
		t.Error("card rules should render the forge-aember-gain line")
	}

	g := NewGame("A", "B", 1)
	if _, ok := g.forgeAemberGainer(0); ok {
		t.Error("no gainer should be found with nothing in play")
	}
	sting := g.AddArtifact(def, 1)
	gainer, ok := g.forgeAemberGainer(0)
	if !ok || gainer != sting {
		t.Errorf("forgeAemberGainer = %d, %v, want %d, true", gainer, ok, sting)
	}
	if _, ok := g.forgeAemberGainer(1); ok {
		t.Error("a player should not gain their own forge spending")
	}

	g.State.Aember[0] = 6
	g.forgeKey(0)
	if g.State.Keys[0] != 1 {
		t.Fatal("the payer should still forge their key")
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("payer pool = %d, want 0", g.State.Aember[0])
	}
	if g.State.Aember[1] != KeyCost {
		t.Errorf("gainer pool = %d, want %d", g.State.Aember[1], KeyCost)
	}
}

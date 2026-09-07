package engine

import "testing"

// TestNoForgeWhileAheadOnKeys checks the constant forge bar Heart of the Forest
// imposes: a player leading on keys skips their forge phase, but a tied or trailing
// player forges normally.
func TestNoForgeWhileAheadOnKeys(t *testing.T) {
	heart := func() CardDefinition {
		return NewCard("heart", Untamed, Artifact, Rare,
			WithRestrictions(Restrictions{NoForgeWhileAheadOnKeys: true}))
	}

	t.Run("bars a player who leads on keys", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(heart(), 0)
		g.State.Keys[0] = 1
		g.State.Aember[0] = 3 * KeyCost
		g.forgePhase(0)
		if g.State.Keys[0] != 1 {
			t.Errorf("keys = %d, want 1 (barred while ahead)", g.State.Keys[0])
		}
	})

	t.Run("lets a tied player forge", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(heart(), 0)
		g.State.Aember[0] = 3 * KeyCost
		g.forgePhase(0)
		if g.State.Keys[0] != 1 {
			t.Errorf("keys = %d, want 1 (tied forges)", g.State.Keys[0])
		}
	})

	t.Run("lets a trailing player forge", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(heart(), 0)
		g.State.Keys[1] = 1
		g.State.Aember[0] = 3 * KeyCost
		g.forgePhase(0)
		if g.State.Keys[0] != 1 {
			t.Errorf("keys = %d, want 1 (trailing forges)", g.State.Keys[0])
		}
	})

	t.Run("does not bar without the restriction in play", func(t *testing.T) {
		g := started(t)
		g.State.Keys[0] = 1
		if g.forgeBarredWhileAhead(0) {
			t.Error("should not bar with no Heart in play")
		}
	})

	t.Run("renders its rule", func(t *testing.T) {
		lines := restrictionText(Restrictions{NoForgeWhileAheadOnKeys: true}, false)
		want := "Each player cannot forge keys while they have more forged keys than their opponent."
		found := false
		for _, l := range lines {
			if l == want {
				found = true
			}
		}
		if !found {
			t.Errorf("lines = %v, want to contain %q", lines, want)
		}
	})
}

package engine

import "testing"

// TestTargetGrantingCard covers the granting-card target: it renders the {card}
// placeholder so the granted-text renderer names the granting card, resolves to
// the card that granted the ability (ctx.Grantor) when one is set, and selects
// nothing when no grantor is set.
func TestTargetGrantingCard(t *testing.T) {
	if got := (Target{Kind: TargetGrantingCard}).Text(); got != CardName {
		t.Errorf("Text() = %q, want %q", got, CardName)
	}

	g := NewGame("A", "B", 1)
	artifact := g.AddArtifact(NewCard("Grantor", StarAlliance, Artifact, Rare), 0)

	withGrantor := &EffectContext{Resolver: g, Grantor: artifact, HasGrantor: true}
	got := Target{Kind: TargetGrantingCard}.selectBase(withGrantor)
	if len(got) != 1 || got[0] != artifact {
		t.Errorf("selectBase with a grantor = %v, want [%d]", got, artifact)
	}

	noGrantor := &EffectContext{Resolver: g}
	if got := (Target{Kind: TargetGrantingCard}).selectBase(noGrantor); got != nil {
		t.Errorf("selectBase without a grantor = %v, want nil", got)
	}
}

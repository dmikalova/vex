package engine

import "testing"

// TestTargetGrantingArtifact covers the granting-artifact target: it resolves to
// the card that granted the ability (ctx.Grantor) when one is set, renders its
// name when Named, falls back to a generic phrase otherwise, and selects nothing
// when no grantor is set.
func TestTargetGrantingArtifact(t *testing.T) {
	named := Target{Kind: TargetGrantingArtifact}.Named("Uncharted Lands")
	if got := named.Text(); got != "Uncharted Lands" {
		t.Errorf("named Text() = %q, want %q", got, "Uncharted Lands")
	}
	if got := (Target{Kind: TargetGrantingArtifact}).Text(); got != "the granting artifact" {
		t.Errorf("unnamed Text() = %q, want %q", got, "the granting artifact")
	}

	g := NewGame("A", "B", 1)
	artifact := g.AddArtifact(NewCard("Grantor", StarAlliance, Artifact, Rare), 0)

	withGrantor := &EffectContext{Resolver: g, Grantor: artifact, HasGrantor: true}
	got := Target{Kind: TargetGrantingArtifact}.selectBase(withGrantor)
	if len(got) != 1 || got[0] != artifact {
		t.Errorf("selectBase with a grantor = %v, want [%d]", got, artifact)
	}

	noGrantor := &EffectContext{Resolver: g}
	if got := (Target{Kind: TargetGrantingArtifact}).selectBase(noGrantor); got != nil {
		t.Errorf("selectBase without a grantor = %v, want nil", got)
	}
}

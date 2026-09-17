package engine

// ArtifactsInPlay counts every artifact in play, both players' — Blossom Drake
// gets +1 power for each.
type ArtifactsInPlay struct{}

// Value returns the number of artifacts in play across both sides.
func (ArtifactsInPlay) Value(ctx *EffectContext) int {
	return len(ctx.Resolver.Artifacts(ctx.Controller)) +
		len(ctx.Resolver.Artifacts(ctx.Opponent()))
}

// CountText renders the "for each" noun, e.g. "artifact in play".
func (ArtifactsInPlay) CountText() string { return "artifact in play" }

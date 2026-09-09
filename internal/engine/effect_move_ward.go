package engine

import "fmt"

// MoveWard takes the ward off one warded creature and places it on another. From
// names the eligible sources; only warded creatures are offered, so nothing
// happens when none carry a ward. Onto names the destination the ward moves to —
// Hunter or Hunted?'s "move a ward from a creature to another creature".
type MoveWard struct {
	// From selects the source creature; only warded creatures are offered as the
	// choice, so nothing happens when none are warded.
	From Target
	// Onto selects the destination creature the ward moves to.
	Onto Target
}

// validate requires both a source and a destination target.
func (e MoveWard) validate() error {
	if !e.From.valid() {
		return errUnsetTarget("MoveWard")
	}
	if !e.Onto.valid() {
		return errUnsetTarget("MoveWard")
	}
	return nil
}

// Text renders the effect, e.g. "move a ward from a creature to another creature".
func (e MoveWard) Text() string {
	return fmt.Sprintf("move a ward from %s to %s", e.From.Text(), e.Onto.Text())
}

// Resolve moves a ward from the chosen warded source to the chosen destination.
// Only warded creatures are offered as sources, so nothing happens when none are
// warded; the chosen source is left in context (ctx.It) so an "another creature"
// destination excludes it.
func (e MoveWard) Resolve(ctx *EffectContext) {
	sources := e.From.selectWith(ctx, false, ctx.Resolver.Warded)
	if len(sources) == 0 {
		return
	}
	from := sources[0]
	ctx.It, ctx.HasIt = from, true
	dest := e.Onto.Select(ctx)
	if len(dest) == 0 {
		return
	}
	onto := dest[0]
	ctx.Resolver.SetWarded(from, false)
	ctx.Resolver.SetWarded(onto, true)
	ctx.Resolver.Record(WardMoved{From: from, To: onto, By: ctx.Source})
}

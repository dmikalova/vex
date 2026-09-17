package engine

import "fmt"

// AttachSelfTo moves the resolving upgrade onto a target creature, keeping it in
// play. It homes a movable upgrade — a Star Alliance "blaster" homes to its
// signature creature (Target.Named), Blast Shielding homes to a chosen neighbor
// of its current host. The upgrade being moved is the one whose granted ability
// is resolving (ctx.Upgrade), so it needs no target for the upgrade itself, only
// the destination creature.
type AttachSelfTo struct {
	// Target selects the creature the upgrade attaches to.
	Target Target
}

// Text renders the effect, e.g. "attach {card} to Commander Chan".
func (e AttachSelfTo) Text() string { return "attach " + CardName + " to " + e.Target.Text() }

// Resolve moves the resolving upgrade onto the target creature.
func (e AttachSelfTo) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate moves the upgrade onto the target and reports whether it now sits
// there, so an "attach {card} to Y -> Z" gate runs Z only when the target admits
// a creature. It does nothing when the target selects none, or when the resolving
// ability was not granted by an attached upgrade.
//
// Binding is by instance, not by name: once the upgrade already sits on a
// creature the target admits, it stays on that exact creature and never re-homes
// onto a second same-named copy, so the payoff that follows always acts on the
// instance it first bound to.
func (e AttachSelfTo) resolveGate(ctx *EffectContext) bool {
	if host, ok := ctx.Resolver.HostOf(ctx.Upgrade); ok && e.Target.couldSelect(ctx, host) {
		return true
	}
	dest := e.Target.Select(ctx)
	if len(dest) == 0 {
		return false
	}
	ctx.Resolver.MoveUpgrade(ctx.Upgrade, dest[0])
	return true
}

// declinable reports that the attachment is a single clickable creature.
func (e AttachSelfTo) declinable() bool { return e.Target.isChosen() }

// vacuous reports that there is no creature to attach to, so a "you may" wrapping
// it need not ask.
func (e AttachSelfTo) vacuous(ctx *EffectContext) bool { return e.Target.empty(ctx) }

// resolveOptional asks for the host declinably, so "you may attach {card} to a
// creature" is answered by clicking that creature rather than by a separate
// Yes/No. An upgrade already sitting on an admissible host stays put, as in
// resolveGate.
func (e AttachSelfTo) resolveOptional(ctx *EffectContext) bool {
	if host, ok := ctx.Resolver.HostOf(ctx.Upgrade); ok && e.Target.couldSelect(ctx, host) {
		return true
	}
	dest := e.Target.SelectOptional(ctx)
	if len(dest) == 0 {
		return false
	}
	ctx.Resolver.MoveUpgrade(ctx.Upgrade, dest[0])
	return true
}

// validate requires a destination target.
func (e AttachSelfTo) validate() error {
	if !e.Target.valid() {
		return fmt.Errorf("AttachSelfTo: Target must be set")
	}
	return nil
}

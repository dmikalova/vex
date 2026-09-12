package engine

import "fmt"

// AttachSelfTo moves the resolving upgrade onto the friendly creature with a
// given printed name, keeping it in play. It is how a movable upgrade — a Star
// Alliance "blaster" — homes to its signature creature: the blaster grants each
// host a Fight/Reap choice to deal damage or send itself to its named creature
// for a payoff, and this node carries out that send. The upgrade being moved is
// the one whose granted ability is resolving (ctx.Upgrade), so a "blaster"
// attaching itself needs no target for the upgrade, only the host to move it to.
type AttachSelfTo struct {
	// Host is the printed name of the creature the upgrade attaches to.
	Host string
}

// Text renders the effect, e.g. "attach {card} to Commander Chan".
func (e AttachSelfTo) Text() string { return "attach " + CardName + " to " + e.Host }

// Resolve moves the resolving upgrade onto the controller's creature named Host.
func (e AttachSelfTo) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate moves the upgrade onto the host and reports whether it now sits
// there, so an "attach {card} to Y -> Z" gate runs Z only when a friendly
// creature named Host is in play. It does nothing when none is, or when the
// resolving ability was not granted by an attached upgrade.
//
// Binding is by instance, not by name: once the blaster already sits on a
// creature named Host, it stays on that exact creature and never re-homes onto a
// second same-named copy, so the payoff that follows always acts on the instance
// it first bound to. A blaster already on its named host counts as attached.
func (e AttachSelfTo) resolveGate(ctx *EffectContext) bool {
	if host, ok := ctx.Resolver.HostOf(ctx.Upgrade); ok &&
		ctx.Resolver.Name(host) == e.Host {
		return true
	}
	for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
		if ctx.Resolver.Name(id) == e.Host {
			ctx.Resolver.MoveUpgrade(ctx.Upgrade, id)
			return true
		}
	}
	return false
}

// validate requires the host creature's name.
func (e AttachSelfTo) validate() error {
	if e.Host == "" {
		return fmt.Errorf("AttachSelfTo: Host name must be set")
	}
	return nil
}

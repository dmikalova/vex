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

// Text renders the effect, e.g. "attach {self} to Commander Chan".
func (e AttachSelfTo) Text() string { return "attach " + SelfName + " to " + e.Host }

// Resolve moves the resolving upgrade onto the controller's creature named Host.
// It does nothing when no friendly creature with that name is in play, or when
// the resolving ability was not granted by an attached upgrade.
//
// Binding is by instance, not by name: once the blaster already sits on a
// creature named Host, it stays on that exact creature and never re-homes onto a
// second same-named copy, so the payoff that follows always acts on the instance
// it first bound to.
func (e AttachSelfTo) Resolve(ctx *EffectContext) {
	if host, ok := ctx.Resolver.HostOf(ctx.Upgrade); ok &&
		ctx.Resolver.Name(host) == e.Host {
		return
	}
	for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
		if ctx.Resolver.Name(id) == e.Host {
			ctx.Resolver.MoveUpgrade(ctx.Upgrade, id)
			return
		}
	}
}

// validate requires the host creature's name.
func (e AttachSelfTo) validate() error {
	if e.Host == "" {
		return fmt.Errorf("AttachSelfTo: Host name must be set")
	}
	return nil
}

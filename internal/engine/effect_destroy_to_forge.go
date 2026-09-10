package engine

import "fmt"

// DestroyFriendlyCreaturesToForge lets the controller destroy any number of their
// own creatures whose combined power meets MinTotalPower; the Then effect
// resolves only if they commit creatures totalling that much. Might Makes Right
// destroys friendly creatures totalling 25 power or more to forge a key at no
// cost.
type DestroyFriendlyCreaturesToForge struct {
	Target        Target
	MinTotalPower int
	Then          Effect
}

// validate requires a target pool, a positive threshold, and a follow-up effect.
func (e DestroyFriendlyCreaturesToForge) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DestroyFriendlyCreaturesToForge")
	}
	if e.MinTotalPower <= 0 {
		return fmt.Errorf("DestroyFriendlyCreaturesToForge: MinTotalPower must be positive")
	}
	if e.Then == nil {
		return fmt.Errorf("DestroyFriendlyCreaturesToForge: Then is required")
	}
	return nil
}

// Text renders the effect, e.g. "you may destroy any number of friendly creatures
// with total power of 25 or more - forge a key at no cost".
func (e DestroyFriendlyCreaturesToForge) Text() string {
	return fmt.Sprintf(
		"you may destroy any number of %ss with total power of %d or more - %s",
		singularNoun(e.Target.Text()),
		e.MinTotalPower,
		e.Then.Text(),
	)
}

// Resolve gathers the controller's picks one at a time, tallying their power. When
// they stop, the creatures are destroyed and Then resolves only if the tally
// reached the threshold; below it, nothing is destroyed.
func (e DestroyFriendlyCreaturesToForge) Resolve(ctx *EffectContext) {
	chosen := pickCards(ctx, "Choose a creature to destroy", 0, true, func() []LocalID {
		return e.Target.Select(ctx)
	})
	total := 0
	for _, id := range chosen {
		total += ctx.Resolver.Power(id)
	}
	if total < e.MinTotalPower {
		return
	}
	Destroy{}.destroy(ctx, chosen)
	e.Then.Resolve(ctx)
}

// SacrificeToForge is Obsidian Forge: the controller destroys any number of their
// own creatures, then forges a key at Extra surcharge reduced by one Æmber for
// each creature destroyed this way. A forge that lands purges the source artifact.
type SacrificeToForge struct {
	Target Target
	Extra  int
}

// validate requires a target pool and a positive surcharge.
func (e SacrificeToForge) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("SacrificeToForge")
	}
	if e.Extra <= 0 {
		return fmt.Errorf("SacrificeToForge: Extra must be positive")
	}
	return nil
}

// Text renders the effect across its clauses.
func (e SacrificeToForge) Text() string {
	return fmt.Sprintf(
		"destroy any number of %ss, then forge a key at +%d Æmber current cost, "+
			"reduced by 1 Æmber for each creature destroyed this way -> purge %s",
		singularNoun(e.Target.Text()),
		e.Extra,
		SelfName,
	)
}

// Resolve gathers the controller's picks one at a time, destroys them, then forges
// at the reduced cost. A forge that lands purges the source artifact.
func (e SacrificeToForge) Resolve(ctx *EffectContext) {
	chosen := pickCards(ctx, "Choose a creature to destroy", 0, true, func() []LocalID {
		return e.Target.Select(ctx)
	})
	Destroy{}.destroy(ctx, chosen)
	extra := max(e.Extra-len(chosen), 0)
	if ctx.Resolver.ForgeKeyAtExtraCost(ctx.Controller, extra) {
		PurgeSource{}.Resolve(ctx)
	}
}

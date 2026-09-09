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
	picked := map[LocalID]bool{}
	var chosen []LocalID
	total := 0
	for {
		var cands []LocalID
		for _, id := range e.Target.Select(ctx) {
			if !picked[id] {
				cands = append(cands, id)
			}
		}
		if len(cands) == 0 {
			break
		}
		pick, ok := ctx.ChooseCardOptional("Choose a creature to destroy", cands)
		if !ok {
			break
		}
		picked[pick] = true
		chosen = append(chosen, pick)
		total += ctx.Resolver.Power(pick)
	}
	if total < e.MinTotalPower {
		return
	}
	Destroy{}.destroy(ctx, chosen)
	e.Then.Resolve(ctx)
}

// SacrificeToForge is Obsidian Forge: the controller destroys any number of their
// own creatures, then may forge a key at Extra surcharge reduced by one Æmber for
// each creature destroyed this way. When a key is forged this way, the source
// artifact is destroyed.
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

// Text renders the effect across its three clauses.
func (e SacrificeToForge) Text() string {
	return fmt.Sprintf(
		"destroy any number of %ss. Then, you may forge a key at +%d Æmber current cost, "+
			"reduced by 1 Æmber for each creature destroyed this way. If you do, destroy %s",
		singularNoun(e.Target.Text()),
		e.Extra,
		SelfName,
	)
}

// Resolve gathers the controller's picks one at a time, destroys them, then offers
// the reduced-cost forge. Forging (which only happens when affordable and accepted)
// destroys the source artifact.
func (e SacrificeToForge) Resolve(ctx *EffectContext) {
	picked := map[LocalID]bool{}
	var chosen []LocalID
	for {
		var cands []LocalID
		for _, id := range e.Target.Select(ctx) {
			if !picked[id] {
				cands = append(cands, id)
			}
		}
		if len(cands) == 0 {
			break
		}
		pick, ok := ctx.ChooseCardOptional("Choose a creature to destroy", cands)
		if !ok {
			break
		}
		picked[pick] = true
		chosen = append(chosen, pick)
	}
	Destroy{}.destroy(ctx, chosen)
	extra := max(e.Extra-len(chosen), 0)
	prompt := fmt.Sprintf("You may forge a key at +%d Æmber current cost", extra)
	if ctx.ChooseOption(prompt, []string{"Yes", "No"}) != 0 {
		return
	}
	if ctx.Resolver.ForgeKeyAtExtraCostReport(ctx.Controller, extra) {
		Destroy{}.destroy(ctx, []LocalID{ctx.Source})
	}
}

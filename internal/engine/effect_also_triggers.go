package engine

import (
	"fmt"
	"slices"
)

// AN ALSO-TRIGGERS-ON rule makes a creature's abilities under one trigger also
// fire on another — its play effect fires again on reap (Kompsos Haruspex), or its
// fight and reap effects each fire on the other (Livia the Elder). Nothing about
// the ability changes; only the moment it fires does. Two sources feed the one
// gather seam in game_abilities.go:
//
//   - a friendly ConstantAbility in play, through ConstantAbility.AlsoTriggers,
//     active only while its source stays in play (Kompsos Haruspex).
//   - an active lasting rule, installed by FuseTriggersForTurn for the remainder
//     of the controller's turn (Livia the Elder).
//
// Because the game state is a flat, pointerless value, a lasting rule is a small
// comparable record (LastingAlsoTriggersOn) held in a fixed array, mirroring the
// lasting reaction registry in game_lasting.go.

// AlsoTriggersOn declares that an ability whose trigger is From also fires when Onto
// occurs. Kompsos Haruspex carries {From: Play, Onto: Reap}: a friendly creature's
// play effect fires again when it reaps.
type AlsoTriggersOn struct {
	From Trigger
	Onto Trigger
}

// valid reports whether both ends are action triggers a card can name
// as "the <verb> effect" — Play, Fight, or Reap.
func (m AlsoTriggersOn) valid() bool {
	return triggerEffectNoun(m.From) != "" && triggerEffectNoun(m.Onto) != ""
}

// LastingAlsoTriggersOn is one also-triggers rule active for the remainder of a
// player's turn: whose it is, and the From/Onto pair it adds. It is a plain
// comparable value so the flat GameState can hold a fixed array of them.
type LastingAlsoTriggersOn struct {
	Controller int8
	From       Trigger
	Onto       Trigger
}

// maxAlsoTriggers bounds how many lasting rules can be active at once — generous
// for the handful of cards that fuse triggers, and a fixed size keeps the state
// flat.
const maxAlsoTriggers = 4

// AddLastingAlsoTriggers registers an also-triggers rule owned by m.Controller for
// the rest of their turn, dropping it silently when the registry is full. It is the seam
// FuseTriggersForTurn uses instead of hardcoding a fused trigger into the reap or
// fight path.
func (g *Game) AddLastingAlsoTriggers(m LastingAlsoTriggersOn) {
	if int(g.State.AlsoTriggersCount) >= maxAlsoTriggers {
		return
	}
	g.State.AlsoTriggers[g.State.AlsoTriggersCount] = m
	g.State.AlsoTriggersCount++
}

// clearLastingAlsoTriggers drops the lasting rules a player owns, called from
// clearLasting when their turn ends so a fuse lasts only that turn.
func (g *Game) clearLastingAlsoTriggers(player int) {
	n := 0
	for i := 0; i < int(g.State.AlsoTriggersCount); i++ {
		if int(g.State.AlsoTriggers[i].Controller) != player {
			g.State.AlsoTriggers[n] = g.State.AlsoTriggers[i]
			n++
		}
	}
	for i := n; i < int(g.State.AlsoTriggersCount); i++ {
		g.State.AlsoTriggers[i] = LastingAlsoTriggersOn{}
	}
	g.State.AlsoTriggersCount = uint8(n)
}

// additionalTriggers returns the triggers whose abilities on creature src should also
// fire when firing occurs, gathered from every active rule that reaches src: a
// friendly ConstantAbility in play (Kompsos Haruspex) and a lasting rule the
// creature's controller owns (Livia the Elder). It is the one helper both cards
// feed, so game_abilities.go stays unaware of which card armed a rule.
func (g *Game) additionalTriggers(src LocalID, firing Trigger) []Trigger {
	var out []Trigger
	add := func(t Trigger) {
		if !slices.Contains(out, t) {
			out = append(out, t)
		}
	}
	for grantor, c := range g.constantAbilitiesInPlay() {
		if len(c.AlsoTriggers) == 0 || !g.constantActive(grantor, c) ||
			!g.constantAffects(grantor, c, src) {
			continue
		}
		for _, m := range c.AlsoTriggers {
			if m.Onto == firing {
				add(m.From)
			}
		}
	}
	controller := g.controller(src)
	for i := 0; i < int(g.State.AlsoTriggersCount); i++ {
		if m := g.State.AlsoTriggers[i]; int(m.Controller) == controller && m.Onto == firing {
			add(m.From)
		}
	}
	return out
}

// FuseTriggersForTurn makes each friendly creature's A abilities also fire on B and
// its B abilities also fire on A, for the remainder of the controller's turn — the
// rule Livia the Elder installs to fuse fight and reap effects. It records the two
// directions as lasting rules; the ready phase drops them.
type FuseTriggersForTurn struct {
	A Trigger
	B Trigger
}

// validate requires two distinct action triggers to fuse.
func (e FuseTriggersForTurn) validate() error {
	if triggerEffectNoun(e.A) == "" || triggerEffectNoun(e.B) == "" {
		return fmt.Errorf("FuseTriggersForTurn: A and B must be Play, Fight, or Reap")
	}
	if e.A == e.B {
		return fmt.Errorf("FuseTriggersForTurn: A and B must differ")
	}
	return nil
}

// Text renders the effect, e.g. "each friendly creature's fight effects and reap
// effects are fight/reap effects for the remainder of the turn".
func (e FuseTriggersForTurn) Text() string {
	a, b := triggerEffectNoun(e.A), triggerEffectNoun(e.B)
	return "each friendly creature's " + a + " effects and " + b +
		" effects are " + a + "/" + b + " effects for the remainder of the turn"
}

// Resolve installs both directions of the fuse for the controller's turn.
func (e FuseTriggersForTurn) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLastingAlsoTriggers(
		LastingAlsoTriggersOn{Controller: int8(ctx.Controller), From: e.A, Onto: e.B},
	)
	ctx.Resolver.AddLastingAlsoTriggers(
		LastingAlsoTriggersOn{Controller: int8(ctx.Controller), From: e.B, Onto: e.A},
	)
}

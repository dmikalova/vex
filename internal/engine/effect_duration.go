package engine

import (
	"fmt"
	"strings"
)

// ForDuration applies several timed effects that share one duration and renders
// their shared "for the remainder of the turn, ..." clause once — "for the
// remainder of the turn, it belongs to house Sanctum and cannot be dealt damage"
// rather than repeating the clause for each effect (Golden Aura). Each child is a
// timed effect whose own Text still reads standalone; ForDuration folds only the
// shared prefix and resolves the children in order, exactly as a Sequence.
type ForDuration struct {
	Duration Duration
	Effects  []Effect
}

// durationScoped is a timed effect that renders its body in two parts — the
// subject it acts on and the predicate it applies — so ForDuration can state the
// shared "for the remainder of the turn, " clause once and, when the children act
// on the same subject, name that subject once too. BelongToHouse and
// CannotBeDealtDamage implement it.
type durationScoped interface {
	durationSubject() string
	durationPredicate() string
}

// validate requires a supported duration and at least two body-renderable timed
// children to combine, then descends into the children.
func (e ForDuration) validate() error {
	if e.Duration != RemainderOfPlayerTurn {
		return fmt.Errorf("ForDuration: duration must be RemainderOfPlayerTurn")
	}
	if len(e.Effects) < 2 {
		return fmt.Errorf("ForDuration: needs at least two effects to combine")
	}
	for _, child := range e.Effects {
		if _, ok := child.(durationScoped); !ok {
			return fmt.Errorf("ForDuration: %T does not render a duration body", child)
		}
		if err := validateEffect(child); err != nil {
			return err
		}
	}
	return nil
}

// Text renders the shared duration clause once, then the child bodies. When every
// child acts on the same subject it is named once and the predicates join with
// " and " — "for the remainder of the turn, it belongs to house Sanctum and cannot
// be dealt damage"; otherwise each subject-predicate body joins with " and ".
func (e ForDuration) Text() string {
	subject := e.Effects[0].(durationScoped).durationSubject()
	shared := true
	predicates := make([]string, len(e.Effects))
	bodies := make([]string, len(e.Effects))
	for i, child := range e.Effects {
		c := child.(durationScoped)
		if c.durationSubject() != subject {
			shared = false
		}
		predicates[i] = c.durationPredicate()
		bodies[i] = c.durationSubject() + " " + c.durationPredicate()
	}
	if shared {
		return "for the remainder of the turn, " + subject + " " +
			strings.Join(predicates, " and ")
	}
	return "for the remainder of the turn, " + strings.Join(bodies, " and ")
}

// Resolve resolves each child effect in order.
func (e ForDuration) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

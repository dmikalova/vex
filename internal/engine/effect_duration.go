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

// durationClause renders the standalone duration prefix a timed effect uses when
// it is not folded into a ForDuration — "for the remainder of the turn" or "during
// your opponent's next turn".
func durationClause(d Duration) string {
	switch d {
	case OpponentNextTurn:
		return "during your opponent's next turn"
	default:
		return "for the remainder of the turn"
	}
}

// durationScoped is a timed effect that renders its body in two parts — the
// subject it acts on and the predicate it applies — so a fold (ForDuration's
// "for the remainder of the turn, ..." prefix, GainUntilNextTurn's "... until the
// start of your next turn" suffix) can state the shared duration clause once and,
// when the children act on the same subject, name that subject once too.
// BelongToHouse, CannotBeDealtDamage, GainAssault, GainKeyword, GainTrait, and
// GainAssaultUntilNextTurn implement it.
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
	subject, joined, shared := foldDurationBodies(e.Effects)
	if shared {
		return "for the remainder of the turn, " + subject + " " + joined
	}
	return "for the remainder of the turn, " + joined
}

// Resolve resolves each child effect in order.
func (e ForDuration) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

// foldDurationBodies renders a set of durationScoped children as one body. When
// every child shares a subject (shared true) it is returned once in subject and
// the predicates join with " and " in joined; otherwise subject is empty and joined
// holds each "<subject> <predicate>" body joined with " and ". The caller frames
// the result with its own duration phrasing — a prefix for ForDuration, a suffix
// for GainUntilNextTurn.
func foldDurationBodies(effects []Effect) (subject, joined string, shared bool) {
	subject = effects[0].(durationScoped).durationSubject()
	shared = true
	predicates := make([]string, len(effects))
	bodies := make([]string, len(effects))
	for i, child := range effects {
		c := child.(durationScoped)
		if c.durationSubject() != subject {
			shared = false
		}
		predicates[i] = c.durationPredicate()
		bodies[i] = c.durationSubject() + " " + c.durationPredicate()
	}
	if shared {
		return subject, strings.Join(predicates, " and "), true
	}
	return "", strings.Join(bodies, " and "), false
}

// GainUntilNextTurn applies several grants that all last until the start of the
// controller's next turn and renders their shared "... until the start of your next
// turn" clause once — "the chosen creature gains skirmish and the Mutant trait until
// the start of your next turn" rather than repeating the duration for each grant
// (the Mutation cycle grants a keyword and the Mutant trait together). Each child is
// a next-turn grant (GainKeyword, GainTrait, GainAssaultUntilNextTurn) whose own
// Text reads standalone; GainUntilNextTurn folds the shared subject and suffix and
// resolves the children in order, exactly as a Sequence.
type GainUntilNextTurn struct {
	Effects []Effect
}

// validate requires at least two body-renderable next-turn children to combine,
// then descends into the children.
func (e GainUntilNextTurn) validate() error {
	if len(e.Effects) < 2 {
		return fmt.Errorf("GainUntilNextTurn: needs at least two effects to combine")
	}
	for _, child := range e.Effects {
		if _, ok := child.(durationScoped); !ok {
			return fmt.Errorf("GainUntilNextTurn: %T does not render a duration body", child)
		}
		if err := validateEffect(child); err != nil {
			return err
		}
	}
	return nil
}

// Text renders the child bodies, then the shared "until the start of your next
// turn" suffix once. When every child acts on the same subject it is named once and
// the predicates join with " and "; otherwise each subject-predicate body joins.
func (e GainUntilNextTurn) Text() string {
	subject, joined, shared := foldDurationBodies(e.Effects)
	if shared {
		return subject + " " + joined + " until the start of your next turn"
	}
	return joined + " until the start of your next turn"
}

// Resolve resolves each child effect in order.
func (e GainUntilNextTurn) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

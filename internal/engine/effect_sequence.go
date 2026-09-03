package engine

import "strings"

// Sequence resolves several effects in order, the way a card lists several
// sentences of rules text that happen one after another. Each child resolves
// fully before the next begins, and the rendered text joins them with ", and".
type Sequence struct {
	Effects []Effect
}

// combinable is a plain "verb the target" effect (e.g. "stun this creature")
// whose text can be folded together with its neighbours in a Sequence. A run of
// combinables folds along whichever axis they share: neighbours with the same
// target fold their verbs ("stun and exhaust this creature"), and neighbours with
// the same verb fold their targets ("destroy an enemy creature and a friendly
// creature"). Either way the Sequence reads as one phrase instead of the clumsier
// "stun this creature, and exhaust this creature".
type combinable interface {
	verb() string
	targetText() string
}

// Text joins the child effect texts, folding each run of combinables that shares
// a verb or a target into a single "verb and verb ... target" or "verb target and
// target ..." phrase.
func (e Sequence) Text() string {
	parts := make([]string, 0, len(e.Effects))
	for i := 0; i < len(e.Effects); {
		c, ok := e.Effects[i].(combinable)
		if !ok {
			parts = append(parts, e.Effects[i].Text())
			i++
			continue
		}
		phrase, next := foldCombinable(e.Effects, i, c)
		parts = append(parts, phrase)
		i = next
	}
	return joinSequenceParts(parts)
}

// joinSequenceParts joins a Sequence's rendered children into one compound
// instruction: "a, and b, and c". A card whose rules are separate statements
// wants Sentences instead, which punctuates each child rather than conjoining.
func joinSequenceParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ", and ")
}

// foldCombinable folds the run of combinables starting at i into one phrase and
// reports the index just past the run. The fold axis is chosen from the first
// neighbour: a shared target folds the verbs, a shared verb folds the targets; a
// lone combinable renders as "verb target".
func foldCombinable(effects []Effect, i int, c combinable) (string, int) {
	verb, target := c.verb(), c.targetText()
	next, ok := peekCombinable(effects, i+1)
	switch {
	case ok && next.targetText() == target:
		verbs := []string{verb}
		i++
		for ; ; i++ {
			n, ok := peekCombinable(effects, i)
			if !ok || n.targetText() != target {
				break
			}
			verbs = append(verbs, n.verb())
		}
		return strings.Join(verbs, " and ") + " " + target, i
	case ok && next.verb() == verb:
		targets := []string{target}
		i++
		for ; ; i++ {
			n, ok := peekCombinable(effects, i)
			if !ok || n.verb() != verb {
				break
			}
			targets = append(targets, n.targetText())
		}
		return verb + " " + strings.Join(targets, " and "), i
	default:
		return verb + " " + target, i + 1
	}
}

// peekCombinable reports the effect at i as a combinable, if it is one and in range.
func peekCombinable(effects []Effect, i int) (combinable, bool) {
	if i >= len(effects) {
		return nil, false
	}
	c, ok := effects[i].(combinable)
	return c, ok
}

// Resolve resolves each child effect in order.
func (e Sequence) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

// validate surfaces the first configuration error among the child effects.
func (e Sequence) validate() error {
	for _, child := range e.Effects {
		if err := validateEffect(child); err != nil {
			return err
		}
	}
	return nil
}

// Sentences resolves several effects in order exactly as a Sequence does, but
// renders each as its own sentence instead of joining them with ", and". It is
// the shape a card takes when its rules are separate statements rather than one
// compound instruction: Sigil of Brotherhood reads "Destroy Sigil of Brotherhood.
// Until the end of the turn, you may use friendly Sanctum creatures", not
// "destroy Sigil of Brotherhood, and until the end of the turn ...". Nest a
// Sequence inside one child to conjoin just that part.
type Sentences struct {
	Effects []Effect
}

// Text renders each child as its own sentence. The first is left uncapitalized
// because whatever precedes it — a trigger prefix, an enclosing clause — decides
// its case; every later child opens a sentence, so it is capitalized here.
func (e Sentences) Text() string {
	if len(e.Effects) == 0 {
		return ""
	}
	text := punctuate(e.Effects[0].Text())
	for _, child := range e.Effects[1:] {
		text += " " + punctuate(capitalizeFirst(child.Text()))
	}
	return text
}

// Resolve resolves each child effect in order.
func (e Sentences) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

// validate surfaces the first configuration error among the child effects.
func (e Sentences) validate() error {
	for _, child := range e.Effects {
		if err := validateEffect(child); err != nil {
			return err
		}
	}
	return nil
}

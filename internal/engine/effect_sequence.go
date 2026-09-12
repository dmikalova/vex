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

// foldable optionally refines combinable: a combinable that only folds under some
// condition reports it here (Exalt folds a single exalt but keeps "exalt X 2
// times" standing alone). A combinable that does not implement foldable always
// folds.
type foldable interface {
	foldable() bool
}

// Text joins the child effect texts, folding each run of combinables that shares
// a verb or a target into a single "verb and verb ... target" or "verb target and
// target ..." phrase.
func (e Sequence) Text() string {
	parts := make([]string, 0, len(e.Effects))
	for i := 0; i < len(e.Effects); {
		c, ok := peekCombinable(e.Effects, i)
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
// instruction: "a", "a, and b", "a, b, and c" — a serial (Oxford) comma once
// there are three or more, never the run-on "a, and b, and c". A card whose rules
// are separate statements wants Sentences instead, which punctuates each child
// rather than conjoining.
func joinSequenceParts(parts []string) string {
	return serialJoin(parts, ", and ")
}

// serialJoin renders parts as an English list. One item stands alone; two are
// joined by two (", and " for independent clauses, " and " for a folded run of
// verbs or targets); three or more take a serial (Oxford) comma, "a, b, and c".
func serialJoin(parts []string, two string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + two + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1]
	}
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
		return serialJoin(verbs, " and ") + " " + target, i
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
		return verb + " " + serialJoin(targets, " and "), i
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
	if !ok {
		return nil, false
	}
	if f, isFoldable := c.(foldable); isFoldable && !f.foldable() {
		return nil, false
	}
	return c, true
}

// Resolve resolves each child effect in order.
func (e Sequence) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

// declinable reports that the sequence leads with a single clickable choice, so a
// May or MayRepeat wrapping it can be driven by that choice (and a Done to pass)
// rather than a separate Yes/No — the rest of the sequence then follows.
func (e Sequence) declinable() bool {
	if len(e.Effects) == 0 {
		return false
	}
	d, ok := e.Effects[0].(declinableEffect)
	return ok && d.declinable()
}

// resolveOptional asks the leading choice declinably; only when it is taken do the
// remaining effects resolve, so declining the first pick passes on the whole
// sequence.
func (e Sequence) resolveOptional(ctx *EffectContext) bool {
	if len(e.Effects) == 0 {
		return false
	}
	first, ok := e.Effects[0].(declinableEffect)
	if !ok || !first.declinable() || !first.resolveOptional(ctx) {
		return false
	}
	for _, child := range e.Effects[1:] {
		child.Resolve(ctx)
	}
	return true
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

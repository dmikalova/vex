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
		if phrase, next, ok := foldNounList(e.Effects, i); ok {
			parts = append(parts, phrase)
			i = next
			continue
		}
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

// nounListable is an effect whose text is a fixed head, an indefinite noun, and a
// fixed tail — "put a tactic from your discard pile into your hand". A run of them
// in a Sequence that shares a head and tail folds into one article-led list, "put
// a tactic, artifact, creature, and upgrade from your discard pile into your
// hand", rather than repeating the tail once per item (Look What I Found!). An
// effect reports listNoun "" when it is not in this shape, so it is not folded.
type nounListable interface {
	listHead() string
	listNoun() string
	listTail() string
}

// foldNounList folds the run of nounListables starting at i that shares a head and
// tail into one article-led noun list, reporting the phrase and the index just
// past the run. It declines (ok false) unless at least two consecutive effects
// qualify.
func foldNounList(effects []Effect, i int) (string, int, bool) {
	head, ok := effects[i].(nounListable)
	if !ok || head.listNoun() == "" {
		return "", i, false
	}
	nouns := []string{head.listNoun()}
	j := i + 1
	for ; j < len(effects); j++ {
		n, ok := effects[j].(nounListable)
		if !ok || n.listNoun() == "" ||
			n.listHead() != head.listHead() || n.listTail() != head.listTail() {
			break
		}
		nouns = append(nouns, n.listNoun())
	}
	if len(nouns) < 2 {
		return "", i, false
	}
	return head.listHead() + " " + indefinite(oxfordAnd(nouns)) + " " + head.listTail(), j, true
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
// May or a Repeat's MayWhile gate wrapping it can be driven by that choice (and a
// Done to pass) rather than a separate Yes/No — the rest of the sequence then
// follows.
func (e Sequence) declinable() bool { return leadsWithACardChoice(e.Effects) }

// resolveOptional asks the leading choice declinably; only when it is taken do the
// remaining effects resolve, so declining the first pick passes on the whole
// sequence.
func (e Sequence) resolveOptional(ctx *EffectContext) bool {
	return resolveLeadingCardChoice(ctx, e.Effects)
}

// leadsWithACardChoice reports that a run of effects opens with one clickable card
// choice. It is what makes a whole run declinable: the lead is the only decision a
// player makes before the rest follows, so clicking that card (or Done) answers
// for the run.
func leadsWithACardChoice(effects []Effect) bool {
	if len(effects) == 0 {
		return false
	}
	d, ok := effects[0].(declinableEffect)
	return ok && d.declinable()
}

// resolveLeadingCardChoice asks a run's leading choice declinably and resolves the
// rest only when it is taken, so declining the first pick passes on the whole run.
func resolveLeadingCardChoice(ctx *EffectContext, effects []Effect) bool {
	if !leadsWithACardChoice(effects) {
		return false
	}
	if !effects[0].(declinableEffect).resolveOptional(ctx) {
		return false
	}
	for _, child := range effects[1:] {
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

// Text renders each child as its own sentence, after folding any threshold ladder
// among them. The first is left uncapitalized because whatever precedes it — a
// trigger prefix, an enclosing clause — decides its case; every later child opens
// a sentence, so it is capitalized here.
func (e Sentences) Text() string {
	parts := foldLadders(e.Effects)
	if len(parts) == 0 {
		return ""
	}
	text := punctuate(parts[0])
	for _, part := range parts[1:] {
		text += " " + punctuate(capitalizeFirst(part))
	}
	return text
}

// ladderRung is a condition that compares a counted subject against a threshold.
// A run of Conditionals whose rungs count the same subject is a threshold ladder —
// Galactic Census pays again at 3, 5, and 6 houses — and naming the subject in
// every rung repeats a long clause verbatim. Folded, the ladder reads "if there
// are 3 or more houses represented among creatures in play, gain 1 Æmber. Gain 1
// more if there are 5 or more. Gain 1 more if there are 6 or more."
type ladderRung interface {
	// ladderSubject names what is counted, so rungs counting different things do
	// not fold together.
	ladderSubject() string
	// ladderThreshold renders the comparison with the subject left implicit, e.g.
	// "if there are 5 or more".
	ladderThreshold() string
}

// ladderRepeating is an effect that can render as a repetition of an identical one
// already stated, so a later rung says "gain 1 more" instead of restating "gain 1
// Æmber". It reports "" when it is not in that shape, and is then not folded.
type ladderRepeating interface {
	repeatedText() string
}

// foldLadders renders each effect's sentence, collapsing every threshold ladder it
// finds. Effects outside a ladder render unchanged.
func foldLadders(effects []Effect) []string {
	parts := make([]string, 0, len(effects))
	for i := 0; i < len(effects); {
		run, next := foldThresholdLadder(effects, i)
		parts = append(parts, run...)
		i = next
	}
	return parts
}

// foldThresholdLadder folds the ladder starting at i, reporting its sentences and
// the index just past it. The first rung keeps its full text and each later one
// becomes "<repeat> <threshold>". A lone rung, or an effect that is no rung at
// all, renders as itself.
func foldThresholdLadder(effects []Effect, i int) ([]string, int) {
	subject, _, repeat, ok := ladderRungAt(effects[i])
	if !ok {
		return []string{effects[i].Text()}, i + 1
	}
	parts := []string{effects[i].Text()}
	j := i + 1
	for ; j < len(effects); j++ {
		s, threshold, r, ok := ladderRungAt(effects[j])
		if !ok || s != subject || r != repeat {
			break
		}
		parts = append(parts, repeat+" "+threshold)
	}
	if len(parts) < 2 {
		return parts[:1], i + 1
	}
	return parts, j
}

// ladderRungAt reports eff as a rung of a threshold ladder: a Conditional with no
// Else, over a condition that compares a counted subject, gating an effect that
// can render as a repetition of itself.
func ladderRungAt(eff Effect) (subject, threshold, repeat string, ok bool) {
	c, isConditional := eff.(Conditional)
	if !isConditional || c.Else != nil {
		return "", "", "", false
	}
	rung, isRung := c.Cond.(ladderRung)
	if !isRung {
		return "", "", "", false
	}
	repeating, isRepeating := c.Then.(ladderRepeating)
	if !isRepeating {
		return "", "", "", false
	}
	if repeat = repeating.repeatedText(); repeat == "" {
		return "", "", "", false
	}
	return rung.ladderSubject(), rung.ladderThreshold(), repeat, true
}

// Resolve resolves each child effect in order.
func (e Sentences) Resolve(ctx *EffectContext) {
	for _, child := range e.Effects {
		child.Resolve(ctx)
	}
}

// declinable reports that the sentences lead with a single clickable choice, so a
// May wrapping them is driven by that click rather than a separate Yes/No.
func (e Sentences) declinable() bool { return leadsWithACardChoice(e.Effects) }

// resolveOptional asks the leading choice declinably; the later sentences resolve
// only when it is taken.
func (e Sentences) resolveOptional(ctx *EffectContext) bool {
	return resolveLeadingCardChoice(ctx, e.Effects)
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

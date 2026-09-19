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

// Text renders each clause as its own sentence — the Rules voice is one
// instruction per sentence (ADR 0019), so joining is something clauses earn by
// folding rather than the default. Folds run first, so a folded run is one
// clause: "destroy an artifact, a creature, and an upgrade" stays one sentence.
func (e Sequence) Text() string { return e.text(true) }

// gatedText renders the sequence as a Conditional's consequence, where the
// default break is suppressed so the gate visibly covers every clause:
// Shadowsaurus must not read "take control of it. That creature belongs to house
// Shadows", whose second sentence looks unconditional. The explicit breaks still
// fire, because grammar outranks house style — Triumph still starts a sentence at
// its inner Conditional. Pinned by TestConditionalKeepsItsSequenceJoined.
func (e Sequence) gatedText() string { return e.text(false) }

// text renders the children, breaking between every clause when breakByDefault
// and only at the explicit breaks otherwise.
func (e Sequence) text(breakByDefault bool) string {
	parts := make([]string, 0, len(e.Effects))
	opensSentence := make([]bool, 0, len(e.Effects))
	add := func(part string, opens bool) {
		parts, opensSentence = append(parts, part), append(opensSentence, opens)
	}
	for i := 0; i < len(e.Effects); {
		opens := opensOwnSentence(e.Effects, i) || (breakByDefault && i > 0)
		if rungs, next := foldThresholdLadder(e.Effects, i); next > i+1 {
			for j, rung := range rungs {
				add(rung, opens || j > 0)
			}
			i = next
			continue
		}
		if phrase, next, ok := foldNounList(e.Effects, i); ok {
			add(phrase, opens)
			i = next
			continue
		}
		c, ok := peekCombinable(e.Effects, i)
		if !ok {
			add(e.Effects[i].Text(), opens)
			i++
			continue
		}
		phrase, next := foldCombinable(e.Effects, i, c)
		add(phrase, opens)
		i = next
	}
	return splitBeforeConditionals(parts, opensSentence)
}

// sentenceEnder is an effect whose own text finishes a sentence, so whatever
// follows it in a Sequence opens a new one rather than being conjoined. A
// choice-led clause is the case: "choose a creature. Destroy the chosen creature,
// and gain 1 Æmber" would run a conjunction on past a full stop.
type sentenceEnder interface {
	endsSentence() bool
}

// opensOwnSentence reports whether the effect at i starts a new sentence: it is a
// Conditional (R2), or the clause before it already closed one (R3).
func opensOwnSentence(effects []Effect, i int) bool {
	if _, isCond := effects[i].(Conditional); isCond {
		return true
	}
	if i == 0 {
		return false
	}
	prev, ok := effects[i-1].(sentenceEnder)
	return ok && prev.endsSentence()
}

// leadInSentence renders a choice and its consequence as two sentences — "choose a
// creature. Destroy the chosen creature" — rather than joining them with a dash.
// The lead keeps its case for whatever prefixes it (a trigger, an enclosing
// clause); the consequence opens a sentence, so it is capitalized. The trailing
// period is left to the caller, as for any other effect fragment.
func leadInSentence(lead, body string) string {
	return punctuate(lead) + " " + capitalizeFirst(body)
}

// hedgedLeadIn is the optional capability of a choice-led effect to render under
// a May. Breaking the lead-in into its own sentence would otherwise leave the
// consequence a bare imperative that reads as mandatory, so the hedge moves onto
// it: "you may choose a house. If you do, reveal the top card of your deck."
type hedgedLeadIn interface {
	hedgedText() string
}

// hedgedLeadInSentence is leadInSentence for an optional choice, carrying the
// hedge onto the consequence.
func hedgedLeadInSentence(lead, body string) string {
	return punctuate(lead) + " If you do, " + body
}

// splitBeforeConditionals renders the folded parts, starting a new sentence at
// every part flagged as opening one. Conjoining a Conditional — "destroy <self>,
// and if your opponent has 3 Æmber or fewer, steal 3 Æmber" — makes the reader
// hold two structures at once, and a Conditional carrying an Else renders its own
// "Otherwise, …" sentence that cannot be conjoined at all. A run with nothing to
// split stays an unpunctuated fragment, so its caller still supplies the period.
func splitBeforeConditionals(parts []string, opensSentence []bool) string {
	runs := make([][]string, 0, len(parts))
	for i, part := range parts {
		if i == 0 || opensSentence[i] {
			runs = append(runs, []string{part})
			continue
		}
		runs[len(runs)-1] = append(runs[len(runs)-1], part)
	}
	sentences := make([]string, len(runs))
	for i, run := range runs {
		sentences[i] = joinSequenceParts(run)
	}
	if len(sentences) < 2 {
		return joinSequenceParts(sentences)
	}
	var text strings.Builder
	text.WriteString(punctuate(sentences[0]))
	for _, s := range sentences[1:] {
		text.WriteString(" " + punctuate(capitalizeFirst(s)))
	}
	return text.String()
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

// joinSequenceParts joins the clauses of one sentence into a compound
// instruction: "a", "a, and b", "a, b, and c" — a serial (Oxford) comma once
// there are three or more, never the run-on "a, and b, and c". Clauses only reach
// here by folding or by sitting inside a gate; every other clause is its own
// sentence.
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

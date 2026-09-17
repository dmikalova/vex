package engine

// COPYING PRINTED STATS lets a creature take on another card's printed stats until
// the copier leaves play: its power becomes the source's printed power, and it
// gains the source's printed armor, keywords, and traits. It copies neither the
// source's name, type, or house, nor its text box — Cyber-Clone copies a creature
// it purges, so the source's live buffs and counters are gone and only the printed
// stats remain to copy.
//
// Like a gained text box, a copy is stored flat as the source's LocalID+1 on the
// copier's CardCore (CopiedStatsSourcePlus): it lasts until the copier leaves play
// (resetCore clears it) and the source's printed stats are read from the immutable
// catalog, so they stay available even after the source card leaves play. The
// power override, armor gain, keyword gain, and trait gain each fold the copy in at
// their read seam — Power, armor, hasKeyword, and HasTrait.

// copiedStatsSource returns the card whose printed stats a creature copies, stored
// as a LocalID+1, so ok is false when it copies none.
func (g *Game) copiedStatsSource(id LocalID) (LocalID, bool) {
	if p := g.State.Cards[id].CopiedStatsSourcePlus; p != 0 {
		return LocalID(p - 1), true
	}
	return 0, false
}

// CopyStats records that recipient copies source's printed stats until recipient
// leaves play. It is the CreatureResolver port method CopyPrintedStats uses.
func (g *Game) CopyStats(recipient, source LocalID) {
	c := g.stateOf(recipient)
	if c == nil {
		return
	}
	c.CopiedStatsSourcePlus = uint8(source) + 1
	g.record(CreatureCopiedStats{Creature: recipient, Source: source})
}

// CopyPrintedStats makes the creature its Target selects copy the printed stats of
// the creature its Source selects until the Target leaves play: its power becomes
// that card's printed power, and it gains that card's printed armor, keywords, and
// traits (Cyber-Clone copies a creature it purges).
type CopyPrintedStats struct {
	Target Target
	Source Target
}

// validate requires both an explicit recipient and an explicit source.
func (e CopyPrintedStats) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("CopyPrintedStats")
	}
	if !e.Source.valid() {
		return errUnsetTarget("CopyPrintedStats source")
	}
	return nil
}

// Text renders the effect, e.g. "{self} has power equal to the same creature's
// printed power and gains its printed armor, keywords, and traits".
func (e CopyPrintedStats) Text() string {
	return e.Target.Text() + " has power equal to " + e.Source.Text() +
		"'s printed power and gains its printed armor, keywords, and traits"
}

// Resolve makes each recipient copy the source creature's printed stats. A source
// that selects nothing copies nothing.
func (e CopyPrintedStats) Resolve(ctx *EffectContext) {
	sources := e.Source.Select(ctx)
	if len(sources) == 0 {
		return
	}
	source := sources[0]
	for _, recipient := range e.Target.Select(ctx) {
		ctx.Resolver.CopyStats(recipient, source)
	}
}

// CreatureCopiedStats narrates a creature copying another card's printed stats.
type CreatureCopiedStats struct {
	Creature LocalID
	Source   LocalID
}

// Text renders the creature and the card whose stats it copied, e.g.
// "Cyber-Clone copies the printed stats of Troll".
func (e CreatureCopiedStats) Text(n Namer) string {
	return n.Name(e.Creature) + " copies the printed stats of " + n.Name(e.Source)
}

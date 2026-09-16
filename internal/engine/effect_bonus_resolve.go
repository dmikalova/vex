package engine

// ResolveBonusIcons resolves the bonus icons printed on the card its Target names,
// as if its controller had just played it: each icon resolves in turn, honoring a
// bar (Master of the Grey) and any substitution, but none of the card's other text
// runs and no play reactions fire. The icons are read from the card's definition, so a
// card a preceding effect revealed in hand (Ensign El-Samra), discarded (LCdr.
// Trigon), or purged (Reclaimed by Nature) still resolves them — the card need not
// be in play.
type ResolveBonusIcons struct {
	Target Target
}

// validate rejects a ResolveBonusIcons whose target was left unset.
func (e ResolveBonusIcons) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("ResolveBonusIcons")
	}
	return nil
}

// Text renders the effect. It always resolves the icons on a card a preceding
// effect just handled (revealed, discarded, or purged), so it names it "that
// card" rather than the ambiguous "it".
func (e ResolveBonusIcons) Text() string {
	return "resolve that card's bonus icons"
}

// Resolve resolves the bonus icons on each card the target names.
func (e ResolveBonusIcons) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.ResolveBonusIconsOn(ctx.Controller, id)
	}
}

// ExtraBonusIconResolution arms a one-shot boost: the next time its controller
// plays a card this turn, each of that card's bonus icons resolves an additional
// time (Wild Bounty). It attaches to the flat lasting registry rather than the play
// path, and the play-path icon loop consumes it on the next played card, resolving
// each icon a second time interleaved with the first.
type ExtraBonusIconResolution struct{}

// Text renders the effect.
func (ExtraBonusIconResolution) Text() string {
	return "the next time you play a card this turn, resolve each of its bonus " +
		"icons an additional time"
}

// Resolve arms the one-shot boost on the controller.
func (ExtraBonusIconResolution) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(LastingEffect{
		On:         EventBonusIconBoost,
		Do:         actResolveBonusAgain,
		Controller: int8(ctx.Controller),
		Once:       true,
		Source:     ctx.Source,
		HasSource:  true,
	})
}

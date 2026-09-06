package engine

// BlankEnemyText blanks the text box of every enemy creature until the
// controller's next turn — its printed keywords, abilities, and constant grants
// are ignored while its traits and stats remain (Shadow of Dis). It takes no
// fields: the affected creatures are always the opponent's, and the duration is
// always "until your next turn".
type BlankEnemyText struct{}

// Text renders the effect as its printed clause.
func (BlankEnemyText) Text() string {
	return "until your next turn, enemy creatures' text boxes are " +
		"considered blank (except for traits)"
}

// Resolve blanks the opponent's creatures until the controller's next turn.
func (BlankEnemyText) Resolve(ctx *EffectContext) {
	ctx.Resolver.BlankEnemyText(ctx.Opponent(), ctx.Source)
}

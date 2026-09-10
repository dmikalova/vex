package engine

// AemberBonusOf counts the Æmber pips printed on the card its Target names, read
// once that card has left play — Rustgnawer gains 1 Æmber for each Æmber bonus on
// the artifact it just destroyed (Target: Triggering, the destroyed card in
// context). The bonus is a printed property, so it survives the card leaving play;
// a card the effect did not remove (still in play) contributes nothing.
type AemberBonusOf struct {
	Target Target
}

// Value returns the printed Æmber bonus of the card the Target names, or 0 when it
// selects nothing or the card is still in play.
func (e AemberBonusOf) Value(ctx *EffectContext) int {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 || resolverInPlay(ctx, ids[0]) {
		return 0
	}
	return ctx.Resolver.AemberBonus(ids[0])
}

// CountText renders the singular noun the "for each" clause repeats, e.g. "Æmber
// bonus on it".
func (e AemberBonusOf) CountText() string { return "Æmber bonus on " + e.Target.Text() }

// AemberOnThis counts the Æmber sitting on the source card, so a card can grow
// with what it captures (Yxili Marauder).
type AemberOnThis struct{}

// Value returns the Æmber on the source card.
func (AemberOnThis) Value(ctx *EffectContext) int { return ctx.Resolver.AmberOn(ctx.Source) }

// CountText renders the singular noun the "for each" clause repeats.
func (AemberOnThis) CountText() string { return "Æmber on it" }

// leadingCountText names the source, for a clause whose buffed subject is some
// other creature — Primus Unguis's "for each Æmber on Primus Unguis" — where a
// trailing "it" would point at the wrong card.
func (AemberOnThis) leadingCountText() string { return "Æmber on " + SelfName }

// AemberInPool counts the Æmber in a player's pool — Sack of Coins deals a
// point of damage for each Æmber in your pool, and Marmo Swarm grows by it.
type AemberInPool struct{ Player Player }

// Value reads that player's current Æmber pool.
func (c AemberInPool) Value(ctx *EffectContext) int {
	return ctx.Resolver.Aember(ctx.PlayerFor(c.Player))
}

// CountText renders the singular noun the "for each" clause repeats.
func (c AemberInPool) CountText() string {
	who := "your"
	if c.Player == Opponent {
		who = "your opponent's"
	}
	return "Æmber in " + who + " pool"
}

// AemberOnFriendlyCreatures counts the Æmber sitting on the controller's
// creatures — Imperial Forge cuts its forge surcharge by it.
type AemberOnFriendlyCreatures struct{}

// Value sums the Æmber on each creature in the controller's battleline.
func (c AemberOnFriendlyCreatures) Value(ctx *EffectContext) int {
	total := 0
	for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
		total += ctx.Resolver.AmberOn(id)
	}
	return total
}

// CountText renders the singular noun the "for each" clause repeats.
func (c AemberOnFriendlyCreatures) CountText() string {
	return "Æmber on friendly creatures"
}

package engine

// AemberOnThis counts the Æmber sitting on the source card, so a card can grow
// with what it captures (Yxili Marauder).
type AemberOnThis struct{}

// Value returns the Æmber on the source card.
func (AemberOnThis) Value(ctx *EffectContext) int { return ctx.amberOn(ctx.Source) }

// CountText renders the singular noun the "for each" clause repeats.
func (AemberOnThis) CountText() string { return "Æmber on it" }

// leadingCountText names the source, for a clause whose buffed subject is some
// other creature — Primus Unguis's "for each Æmber on Primus Unguis" — where a
// trailing "it" would point at the wrong card.
func (AemberOnThis) leadingCountText() string { return "Æmber on " + SelfName }

// CountClause renders the clause CountIs puts after "if", e.g. "there are 4 or
// more Æmber on it". Æmber is a mass noun, so the plural flag does not change it.
func (AemberOnThis) CountClause(quantity string, _ bool) string {
	return "there are " + quantity + " Æmber on it"
}

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

// eachPlayerEqualToText renders the pool from the point of view of the player
// being paid, so Binate Rupture says "equal to the Æmber in their pool".
func (c AemberInPool) eachPlayerEqualToText() string {
	who := "their"
	if c.Player == Opponent {
		who = "their opponent's"
	}
	return "the \u00c6mber in " + who + " pool"
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

package engine

// EmitAemberStolenFrom fires each of the victim's in-play After Æmber Is Stolen
// From You abilities once Æmber has been stolen from them, carrying the amount
// taken in that single theft so an ability can scale by it (Molephin deals 1
// damage to each enemy creature for each Æmber stolen). A theft of nothing, or
// from the other player, fires nothing.
func (g *Game) EmitAemberStolenFrom(victim, amount int) {
	if amount <= 0 {
		return
	}
	for _, id := range g.allInPlay(victim) {
		pending := g.triggeredBy(id, TriggerAfterAemberStolenFromYou)
		if len(pending) == 0 || !g.inPlay(id) {
			continue
		}
		for _, t := range g.orderTriggered(victim, pending) {
			closeFrame := g.openFrame(Frame{
				Actor:      victim,
				Source:     id,
				HasSource:  true,
				Trigger:    TriggerAfterAemberStolenFromYou,
				Grantor:    t.grantor,
				HasGrantor: t.grantor != id,
			})
			ctx := &EffectContext{
				Resolver:   g,
				Source:     id,
				Controller: victim,
			}
			ctx.Produced.AemberStolen = amount
			t.ability.Effect.Resolve(ctx)
			closeFrame()
			g.settleDestroyed(victim)
		}
	}
}

// AemberStolenThisEvent is the number of Æmber taken in the single theft that
// fired an After Æmber Is Stolen From You ability, read by a "for each Æmber
// stolen" clause (Molephin). It reads the amount stashed on the context when the
// trigger fires.
type AemberStolenThisEvent struct{}

// Value returns how much Æmber was just stolen in this theft.
func (AemberStolenThisEvent) Value(ctx *EffectContext) int { return ctx.Produced.AemberStolen }

// CountText renders the singular noun a "for each" clause repeats.
func (AemberStolenThisEvent) CountText() string { return "Æmber stolen" }

// CountClause renders the clause CountIs puts after "if". Æmber is a mass noun,
// so the plural flag does not change it.
func (AemberStolenThisEvent) CountClause(quantity string, _ bool) string {
	return quantity + " Æmber was stolen"
}

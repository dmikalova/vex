package engine

// FoughtCreatureIsMostPowerfulEnemy is met, in a Before Fight ability, when the
// creature the source is about to fight (ctx.It) has no enemy creature with
// greater power — it is the most powerful enemy creature. A tie still qualifies
// (no enemy creature has strictly greater power), so there is no choice to make.
// Baldric the Bold reads it to gain Æmber for picking off the opponent's biggest
// threat.
type FoughtCreatureIsMostPowerfulEnemy struct{}

// CondText renders the condition with the "if" prefix every condition leads with.
func (FoughtCreatureIsMostPowerfulEnemy) CondText() string {
	return "if the fought creature is the most powerful enemy creature"
}

// Met reports whether the fought creature (ctx.It) ties for or exceeds the power
// of every enemy creature.
func (FoughtCreatureIsMostPowerfulEnemy) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	foughtPower := ctx.Resolver.Power(ctx.It)
	for _, id := range (Target{Kind: TargetEachEnemyCreature}).Select(ctx) {
		if ctx.Resolver.Power(id) > foughtPower {
			return false
		}
	}
	return true
}

package engine

// ForgeKey has the controller forge a key outside the normal start-of-turn step.
// By default they pay the current key cost, if they can afford it; FreeOfCost
// forges without paying. Both paths fire "after you forge a key" abilities and,
// on the final key, win the game.
type ForgeKey struct {
	FreeOfCost bool
}

// Text renders the effect.
func (e ForgeKey) Text() string {
	if e.FreeOfCost {
		return "forge a key at no cost"
	}
	return "forge a key at current cost"
}

// Resolve forges one key for the controller if affordable.
func (e ForgeKey) Resolve(ctx *EffectContext) {
	if e.FreeOfCost {
		ctx.Resolver.ForgeKeyFree(ctx.Controller)
		return
	}
	ctx.Resolver.ForgeKey(ctx.Controller)
}

// GiveRemainingAemberAfterOpponentForgeKey arms Interdimensional Graft's delayed
// forge penalty: each time the opponent forges a key during their next turn, they
// give their remaining Æmber to the controller. It is durable across this turn's
// end and fires on every forge that turn (a key cheat can forge more than one), then
// expires at the end of that opponent's next turn.
type GiveRemainingAemberAfterOpponentForgeKey struct{}

// Text renders the effect.
func (GiveRemainingAemberAfterOpponentForgeKey) Text() string {
	return "if an opponent forges a key on their next turn, they must give you their remaining Æmber"
}

// Resolve arms the transfer for the opponent's next turn by registering a reaction
// to the opponent's forge, owned by the opponent so it survives this turn, fires on
// each forge during theirs, and clears at the end of their turn.
func (GiveRemainingAemberAfterOpponentForgeKey) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(EventForgeKey, actGiveRemainingAember, ctx.Opponent(), 0)
}

// UnforgeKey takes a forged key back off a player (Key Hammer). It is the one
// effect that lowers a key count, so it is a node of its own rather than a
// negative ForgeKey: nothing is paid, nothing is refunded, and no "after you
// forge a key" ability fires.
type UnforgeKey struct {
	Player Player
}

// validate rejects an UnforgeKey whose player was left unset.
func (e UnforgeKey) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("UnforgeKey")
	}
	return nil
}

// Text renders the effect, e.g. "unforge one of your opponent's keys".
func (e UnforgeKey) Text() string {
	if e.Player == Opponent {
		return "unforge one of your opponent's keys"
	}
	return "unforge one of your keys"
}

// Resolve takes one key back off the named player.
func (e UnforgeKey) Resolve(ctx *EffectContext) {
	ctx.Resolver.UnforgeKey(ctx.PlayerFor(e.Player))
}

package engine

// Restriction effects forbid a player some action for a stretch of the game,
// rather than changing the board directly. A "cannot" rule can arrive two ways:
// as a timed effect (this file) or as a constant rule printed on a card in play
// (CardDefinition.Restricts). Both feed the same gates (Game.cannotFight,
// Game.cannotPlayCreatures), so the restriction is expressed once and honored the
// same way however it is imposed.

// CannotFight bars a player from using creatures to fight. As an effect it is a
// timed bar — Fogbank stops an opponent for the Duration of their next turn. The
// same bar can be printed on a card as a constant Restrictions.Fighting rule; the
// fight gate consults both.
type CannotFight struct {
	Player   Player
	Duration Duration
}

// validate rejects a CannotFight whose player or duration was left unset.
func (e CannotFight) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("CannotFight")
	}
	if !e.Duration.valid() {
		return errUnsetDuration("CannotFight")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent cannot use creatures to fight
// during their next turn".
func (e CannotFight) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	return who + " cannot use creatures to fight during " + whose + " next turn"
}

// Resolve applies the timed bar to the chosen player.
func (e CannotFight) Resolve(ctx *EffectContext) {
	if e.Duration == NextTurn {
		ctx.Resolver.CannotFightNextTurn(ctx.PlayerFor(e.Player))
	}
}

// GrantFightForChosenHouse lets the controller's creatures of the house picked by
// an enclosing ChooseHouseThen fight this turn even out of the active house —
// Brothers in Battle's "each friendly creature of that house may fight." The
// grant lasts only the current turn (EndTurn clears it).
type GrantFightForChosenHouse struct{}

// Text renders the effect.
func (GrantFightForChosenHouse) Text() string {
	return "for the remainder of the turn, each friendly creature of the chosen house may fight"
}

// Resolve grants the controller's chosen-house creatures the right to fight this
// turn.
func (GrantFightForChosenHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.GrantFightForHouse(ctx.Controller, ctx.ChosenHouse)
}

// ForceOpponentActiveHouse makes the opponent choose the house picked by an
// enclosing ChooseHouseThen as their active house on their next turn — Control the
// Weak's "your opponent must choose that house as their active house during their
// next turn."
type ForceOpponentActiveHouse struct{}

// Text renders the effect.
func (ForceOpponentActiveHouse) Text() string {
	return "your opponent must choose that house as their active house during their next turn"
}

// Resolve arms the forced house on the opponent's next turn.
func (ForceOpponentActiveHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.ForceActiveHouseNextTurn(ctx.Opponent(), ctx.ChosenHouse)
}

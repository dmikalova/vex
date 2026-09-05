package engine

import (
	"fmt"
	"strings"
)

// A restriction forbids a player some action for a stretch of the game — "cannot
// use creatures to fight", "cannot play creatures" — rather than changing the board
// directly. A restriction can be a timed effect that lasts through a player's next
// turn, or a constant rule printed on a card in play; while it is active the
// forbidden action simply cannot be taken. When one effect says a player "cannot"
// and another says they "must" or "may" do the same thing, "cannot" wins.
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
		ctx.Resolver.CannotFightNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
	}
}

// CannotPlay bars a player from playing cards for the Duration — Lifeward stops
// creatures and Scrambler Storm stops action cards through the affected player's
// next turn, while Treasure Map stops every card for the rest of the current turn.
// It mirrors CannotFight: a Player, a Duration, and here the card Type that is
// barred, which when left unset bars every type.
type CannotPlay struct {
	Player   Player
	Type     CardType
	Duration Duration
}

// barred is the type this bar installs: the named one, or the AnyType wildcard
// when the author left Type unset to bar everything.
func (e CannotPlay) barred() CardType {
	if e.Type == TypeUnset {
		return AnyType
	}
	return e.Type
}

// validate rejects a CannotPlay whose player or duration was left unset. An unset
// Type is allowed and means every type.
func (e CannotPlay) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("CannotPlay")
	}
	if !e.Duration.valid() {
		return errUnsetDuration("CannotPlay")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent cannot play creatures during their
// next turn". The Tactic type prints as "action cards" (rule 19).
func (e CannotPlay) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	noun := strings.ToLower(e.barred().String()) + "s"
	if e.Type == Tactic {
		noun = "action cards"
	}
	when := "during " + whose + " next turn"
	if e.Duration == EndOfTurn {
		when = "for the remainder of the turn"
	}
	return who + " cannot play " + noun + " " + when
}

// Resolve arms the play-type bar on the chosen player for the Duration.
func (e CannotPlay) Resolve(ctx *EffectContext) {
	switch e.Duration {
	case NextTurn:
		ctx.Resolver.CannotPlayTypeNextTurn(ctx.PlayerFor(e.Player), e.barred(), ctx.Source)
	case EndOfTurn:
		ctx.Resolver.CannotPlayTypeThisTurn(ctx.PlayerFor(e.Player), e.barred(), ctx.Source)
	}
}

// CannotUse bars a player from using any card — reaping, fighting, or an "Action:"
// ability — for the Duration (Skippy Timehog). It is the broadest of the bars:
// where CannotFight stops one verb and CannotPlay stops cards leaving hand, this
// stops every use of what is already in play. Playing and discarding still work.
type CannotUse struct {
	Player   Player
	Duration Duration
}

// validate rejects a CannotUse whose player or duration was left unset.
func (e CannotUse) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("CannotUse")
	}
	if !e.Duration.valid() {
		return errUnsetDuration("CannotUse")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent cannot use any cards during their
// next turn".
func (e CannotUse) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	return who + " cannot use any cards during " + whose + " next turn"
}

// Resolve arms the use bar on the chosen player.
func (e CannotUse) Resolve(ctx *EffectContext) {
	if e.Duration == NextTurn {
		ctx.Resolver.CannotUseNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
	}
}

// SkipForgeStep makes a player skip their "forge a key" step at the start of their
// next turn (Miasma).
type SkipForgeStep struct {
	Player Player
}

// validate rejects a SkipForgeStep whose player was left unset.
func (e SkipForgeStep) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("SkipForgeStep")
	}
	return nil
}

// Text renders the effect, e.g. `your opponent skips the "forge a key" step during
// their next turn`.
func (e SkipForgeStep) Text() string {
	who, whose, verb := "you", "your", "skip"
	if e.Player == Opponent {
		who, whose, verb = "your opponent", "their", "skips"
	}
	return fmt.Sprintf("%s %s the %q step during %s next turn", who, verb, "forge a key", whose)
}

// Resolve arms the skip on the chosen player's next turn.
func (e SkipForgeStep) Resolve(ctx *EffectContext) {
	ctx.Resolver.SkipForgeStepNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
}

// GrantFightForChosenHouse lets the controller's creatures of the house picked by
// an enclosing ChooseHouseThen fight this turn even out of the active house —
// Brothers in Battle's "each friendly creature of that house may fight." The
// grant lasts only the current turn (the ready phase clears it).
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

// GrantFightAnyHouse lets every creature the controller has fight this turn, whatever
// its house — Follow the Leader's "each friendly creature may fight", and Horseman
// of War's longer wording for the same rule. It is GrantFightForChosenHouse with the
// house filter dropped. The grant lasts only the current turn (the ready phase
// clears it).
type GrantFightAnyHouse struct{}

// Text renders the effect.
func (GrantFightAnyHouse) Text() string {
	return "for the remainder of the turn, each friendly creature may fight"
}

// Resolve grants the controller's creatures the right to fight this turn.
func (GrantFightAnyHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.GrantFightAnyHouse(ctx.Controller)
}

// MayUseFriendlyHouse lets the controller fully use (fight, reap, or Action:) their
// creatures of House this turn even out of the active house — Sigil of Brotherhood,
// Ritual of the Hunt. The grant lasts only the current turn (the ready phase
// clears it).
type MayUseFriendlyHouse struct {
	House House
}

// validate rejects a MayUseFriendlyHouse whose house was left unset.
func (e MayUseFriendlyHouse) validate() error {
	if e.House == HouseNone {
		return fmt.Errorf("MayUseFriendlyHouse: house must be set")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, you may use friendly
// Sanctum creatures".
func (e MayUseFriendlyHouse) Text() string {
	return fmt.Sprintf("for the remainder of the turn, you may use friendly %s creatures", e.House)
}

// Resolve grants the controller full use of their House creatures this turn.
func (e MayUseFriendlyHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.GrantUseForHouse(ctx.Controller, e.House)
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
	ctx.Resolver.ForceActiveHouseNextTurn(ctx.Opponent(), ctx.ChosenHouse, ctx.Source)
}

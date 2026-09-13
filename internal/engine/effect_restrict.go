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
	if e.Duration == OpponentNextTurn {
		ctx.Resolver.CannotFightNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
	}
}

// CannotReap bars a player from using creatures to reap. As an effect it is a
// timed bar — Inky Gloom stops an opponent for the Duration of their next turn.
// It is narrower than CannotUse: only reaping is barred, so the affected player's
// creatures can still fight and fire "Action:" abilities.
type CannotReap struct {
	Player   Player
	Duration Duration
}

// validate rejects a CannotReap whose player or duration was left unset.
func (e CannotReap) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("CannotReap")
	}
	if !e.Duration.valid() {
		return errUnsetDuration("CannotReap")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent cannot use creatures to reap
// during their next turn", or the current-turn form "you cannot use creatures to
// reap for the remainder of the turn" (Ragnarok).
func (e CannotReap) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	if e.Duration == RemainderOfPlayerTurn {
		return who + " cannot use creatures to reap for the remainder of the turn"
	}
	return who + " cannot use creatures to reap during " + whose + " next turn"
}

// Resolve applies the timed bar to the chosen player, for the current turn
// (RemainderOfPlayerTurn) or the player's next turn (OpponentNextTurn).
func (e CannotReap) Resolve(ctx *EffectContext) {
	switch e.Duration {
	case RemainderOfPlayerTurn:
		ctx.Resolver.CannotReapThisTurn(ctx.PlayerFor(e.Player), ctx.Source)
	case OpponentNextTurn:
		ctx.Resolver.CannotReapNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
	}
}

// CreaturesCannot bars every creature in play — both players' — from being used
// one way (fighting or reaping) until the start of the caster's next turn, save
// for creatures of an excepted house. Where CannotFight and CannotReap bar one
// player's use of their own creatures, this is a rule on the whole board: Into
// the Night stops non-Shadows creatures fighting, Sow Salt stops every creature
// reaping. ExceptHouse left unset (HouseNone) spares no house.
type CreaturesCannot struct {
	Action      UseKind
	ExceptHouse House
	Duration    Duration
}

// validate rejects a CreaturesCannot that bars no real action or names no
// duration. The action must be fighting or reaping; an "Action:" ability cannot
// be barred this way. An unset ExceptHouse is legal and spares no house.
func (e CreaturesCannot) validate() error {
	if e.Action != FightUse && e.Action != ReapUse {
		return fmt.Errorf("CreaturesCannot: action must be FightUse or ReapUse")
	}
	if !e.Duration.valid() {
		return errUnsetDuration("CreaturesCannot")
	}
	return nil
}

// Text renders the effect, e.g. "until the start of your next turn, non-Shadows
// creatures cannot be used to fight", or with no house exception "until the start
// of your next turn, creatures cannot be used to reap".
func (e CreaturesCannot) Text() string {
	subject := "creatures"
	if e.ExceptHouse != HouseNone {
		subject = "non-" + e.ExceptHouse.String() + " creatures"
	}
	return "until the start of your next turn, " + subject +
		" cannot be used to " + e.Action.verb()
}

// Resolve arms the board-wide bar: it stops the caster using creatures this way
// for the rest of their turn and the opponent throughout their next turn, so it
// lifts at the start of the caster's next turn.
func (e CreaturesCannot) Resolve(ctx *EffectContext) {
	if e.Duration == StartOfPlayerNextTurn {
		ctx.Resolver.CreaturesCannotUntilNextTurn(
			ctx.Controller,
			e.Action,
			e.ExceptHouse,
			ctx.Source,
		)
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
// next turn". Card text raises the type noun to its proper-noun capitalization at
// the presentation layer, so this keeps it lowercase.
func (e CannotPlay) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	noun := strings.ToLower(e.barred().String()) + "s"
	when := "during " + whose + " next turn"
	if e.Duration == RemainderOfPlayerTurn {
		when = "for the remainder of the turn"
	}
	return who + " cannot play " + noun + " " + when
}

// Resolve arms the play-type bar on the chosen player for the Duration.
func (e CannotPlay) Resolve(ctx *EffectContext) {
	switch e.Duration {
	case OpponentNextTurn:
		ctx.Resolver.CannotPlayTypeNextTurn(ctx.PlayerFor(e.Player), e.barred(), ctx.Source)
	case RemainderOfPlayerTurn:
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
// next turn", or for the RemainderOfPlayerTurn form "you cannot use cards this
// turn" (United Action).
func (e CannotUse) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	if e.Duration == RemainderOfPlayerTurn {
		return who + " cannot use cards this turn"
	}
	return who + " cannot use any cards during " + whose + " next turn"
}

// Resolve arms the use bar on the chosen player, this turn or their next.
func (e CannotUse) Resolve(ctx *EffectContext) {
	switch e.Duration {
	case OpponentNextTurn:
		ctx.Resolver.CannotUseNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
	case RemainderOfPlayerTurn:
		ctx.Resolver.CannotUseThisTurn(ctx.PlayerFor(e.Player), ctx.Source)
	}
}

// ChosenHouseCannotReapNextTurn bars a player from reaping with creatures of the
// house an enclosing ChooseHouseThen picked, throughout that player's next turn
// (Seismo-entangler's "During your opponent's next turn, creatures of the chosen
// house cannot be used to reap").
type ChosenHouseCannotReapNextTurn struct {
	Player Player
}

// validate rejects an effect whose player was left unset.
func (e ChosenHouseCannotReapNextTurn) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("ChosenHouseCannotReapNextTurn")
	}
	return nil
}

// Text renders the effect, e.g. "during your opponent's next turn, creatures of
// the chosen house cannot be used to reap".
func (e ChosenHouseCannotReapNextTurn) Text() string {
	whose := "your"
	if e.Player == Opponent {
		whose = "your opponent's"
	}
	return "during " + whose + " next turn, creatures of the chosen house cannot be used to reap"
}

// Resolve arms the reap-by-house bar on the chosen player's next turn.
func (e ChosenHouseCannotReapNextTurn) Resolve(ctx *EffectContext) {
	ctx.Resolver.CannotReapHouseNextTurn(ctx.PlayerFor(e.Player), ctx.ChosenHouse, ctx.Source)
}

// ActiveHouseSource names where an active-house constraint reads the house it
// applies (and, for JustChosen, whose next turn it binds): the house an enclosing
// ChooseHouseThen picked (Chosen), the house of the creature the source fought
// (Fought), or the house a player just chose this turn, read from the board
// (JustChosen — Snag's Mirror, keyed off the "after a player chooses a house"
// trigger).
type ActiveHouseSource uint8

const (
	activeHouseUnset ActiveHouseSource = iota
	// ChosenActiveHouse reads the house an enclosing ChooseHouseThen picked.
	ChosenActiveHouse
	// FoughtActiveHouse reads the house of the creature the source fought (ctx.It).
	FoughtActiveHouse
	// JustChosenActiveHouse reads the house a player just chose, from the board.
	JustChosenActiveHouse
)

// OpponentMustChooseHouse makes the opponent choose a house as their active house
// on their next turn — the house an enclosing ChooseHouseThen picked (Control the
// Weak, Source Chosen) or the house of the creature the source fought (Snag,
// Source Fought).
type OpponentMustChooseHouse struct {
	// Source names where the forced house is read from (Chosen or Fought).
	Source ActiveHouseSource
}

// validate requires a Chosen or Fought source.
func (e OpponentMustChooseHouse) validate() error {
	switch e.Source {
	case ChosenActiveHouse, FoughtActiveHouse:
		return nil
	default:
		return fmt.Errorf("OpponentMustChooseHouse: Source must be Chosen or Fought")
	}
}

// Text renders the effect, e.g. "your opponent must choose that house as their
// active house during their next turn".
func (e OpponentMustChooseHouse) Text() string {
	house, when := "that house", "during their next turn"
	if e.Source == FoughtActiveHouse {
		house, when = "the house of the creature "+SelfName+" fights", "on their next turn"
	}
	return "your opponent must choose " + house + " as their active house " + when
}

// Resolve arms the forced house on the opponent's next turn. A Fought source with
// no creature in context does nothing.
func (e OpponentMustChooseHouse) Resolve(ctx *EffectContext) {
	house := ctx.ChosenHouse
	if e.Source == FoughtActiveHouse {
		if !ctx.HasIt {
			return
		}
		house = ctx.Resolver.House(ctx.It)
	}
	ctx.Resolver.ForceActiveHouseNextTurn(ctx.Opponent(), house, ctx.Source)
}

// OpponentCannotChooseHouse bars a player from choosing a house as their active
// house on their next turn — the source's opponent from the chosen house (Tezmal,
// Source Chosen) or the chooser's opponent from the house a player just chose
// (Snag's Mirror, Source JustChosen, keyed off the "after a player chooses a
// house" trigger so it works whichever player chose).
type OpponentCannotChooseHouse struct {
	// Source names where the barred house is read from (Chosen or JustChosen).
	Source ActiveHouseSource
}

// validate requires a Chosen or JustChosen source.
func (e OpponentCannotChooseHouse) validate() error {
	switch e.Source {
	case ChosenActiveHouse, JustChosenActiveHouse:
		return nil
	default:
		return fmt.Errorf("OpponentCannotChooseHouse: Source must be Chosen or JustChosen")
	}
}

// Text renders the effect. A JustChosen source follows the "after a player chooses
// an active house, " trigger prefix, so it speaks from the chooser's point of view.
func (e OpponentCannotChooseHouse) Text() string {
	if e.Source == JustChosenActiveHouse {
		return "their opponent cannot choose the same house as their active house on their next turn"
	}
	return "your opponent cannot choose that house as their active house on their next turn"
}

// Resolve bars the house on the barred player's next turn. A JustChosen source
// reads the active player and house from the board, so it bars the chooser's
// opponent whichever player made the choice.
func (e OpponentCannotChooseHouse) Resolve(ctx *EffectContext) {
	if e.Source == JustChosenActiveHouse {
		chooser := ctx.Resolver.ActivePlayer()
		ctx.Resolver.ForbidActiveHouseNextTurn(1-chooser, ctx.Resolver.ActiveHouse(), ctx.Source)
		return
	}
	ctx.Resolver.ForbidActiveHouseNextTurn(ctx.Opponent(), ctx.ChosenHouse, ctx.Source)
}

// WagerOpponentChoosesChosenHouse bets on the opponent matching the house an
// enclosing ChooseHouseThen picked: if they choose it as their active house next
// turn, the controller steals Amount — Snaglet's "if your opponent chooses that
// house as their active house on their next turn, steal 2A."
type WagerOpponentChoosesChosenHouse struct {
	// Amount is how much Æmber the controller steals if the bet lands.
	Amount int
}

// Text renders the effect.
func (e WagerOpponentChoosesChosenHouse) Text() string {
	return fmt.Sprintf(
		"if your opponent chooses that house as their active house on their next turn, steal %d Æmber",
		e.Amount,
	)
}

// Resolve arms the wager on the opponent's next turn.
func (e WagerOpponentChoosesChosenHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.WagerOnHouseNextTurn(
		ctx.Opponent(), ctx.ChosenHouse, e.Amount, ctx.Controller, ctx.Source,
	)
}

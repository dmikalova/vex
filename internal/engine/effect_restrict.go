package engine

import (
	"fmt"
	"strings"
)

// A restriction forbids a player some action for a stretch of the game — "cannot
// use creatures to fight", "cannot play creatures" — rather than changing the
// board directly. A "cannot" rule can arrive two ways: as a timed effect (this
// file) that lasts through a player's next turn, or as a constant rule printed on
// a card in play (CardDefinition.Restricts). Both feed the same gates
// (Game.cannotFight, Game.cannotPlayCreatures), so the restriction is expressed
// once and honored the same way however it is imposed. When one effect says a
// player "cannot" and another says they "must" or "may" do the same thing,
// "cannot" wins.

// A RestrictKind names what a Restrict bars a player from: using creatures to
// fight (Fogbank), using creatures to reap (Inky Gloom, Ragnarok), or using any
// cards at all (Skippy Timehog, United Action). Barring fighting and barring
// reaping each leave the other verb — and "Action:" abilities — open; barring use
// stops every way of using what is already in play.
type RestrictKind uint8

const (
	// restrictUnset is the invalid zero value; a real Restrict names a kind.
	restrictUnset RestrictKind = iota
	// RestrictFighting bars using creatures to fight.
	RestrictFighting
	// RestrictReaping bars using creatures to reap.
	RestrictReaping
	// RestrictUse bars using any card — reaping, fighting, or an "Action:" ability.
	RestrictUse
	// restrictKindCount bounds the enum for valid checks.
	restrictKindCount
)

// valid reports whether the restrict kind names one of the three real bars.
func (k RestrictKind) valid() bool {
	return k > restrictUnset && k < restrictKindCount
}

// phrase renders what the kind bars, filling "cannot use ___" — "creatures to
// fight", "creatures to reap", or the broader "any cards".
func (k RestrictKind) phrase() string {
	switch k {
	case RestrictFighting:
		return "creatures to fight"
	case RestrictReaping:
		return "creatures to reap"
	default:
		return "any cards"
	}
}

// Restrict bars a player from an action for a Duration — using creatures to
// fight, using creatures to reap, or using any cards at all. As an effect it is a
// timed bar: Fogbank stops an opponent fighting throughout their next turn, Inky
// Gloom stops them reaping, Ragnarok stops the caster reaping for the rest of the
// turn, Skippy Timehog stops every use of the opponent's cards. Fighting and
// reaping bars can also be printed on a card as constant Restrictions rules; the
// same gate (Game.cannotFight, Game.cannotReap) consults both. Fighting has no
// current-turn gate, so RestrictFighting is only valid for OpponentNextTurn. House
// narrows a reaping bar to one house (Seismo-entangler bars the chosen house);
// left unset it bars every house.
type Restrict struct {
	Player   Player
	Action   RestrictKind
	House    HouseChoice
	Duration Duration
}

// validate rejects a Restrict whose player, action, or duration was left unset,
// or a RestrictFighting for anything but the player's next turn (there is no
// current-turn fight bar). A House scope is only honored for a reaping bar on the
// player's next turn — the only house-scoped gate is CannotReapHouseNextTurn.
func (e Restrict) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("Restrict")
	}
	if !e.Action.valid() {
		return fmt.Errorf("Restrict: action must be set")
	}
	if !e.Duration.valid() {
		return errUnsetDuration("Restrict")
	}
	if e.Action == RestrictFighting && e.Duration != OpponentNextTurn {
		return fmt.Errorf(
			"Restrict: fighting can only be barred during the player's next turn",
		)
	}
	if e.House.phrase() != "" &&
		(e.Action != RestrictReaping || e.Duration != OpponentNextTurn) {
		return fmt.Errorf(
			"Restrict: a house scope only bars reaping during the player's next turn",
		)
	}
	return nil
}

// Text renders the effect, e.g. "your opponent cannot use creatures to fight
// during their next turn", or the current-turn form "you cannot use creatures to
// reap for the remainder of the turn" (Ragnarok). A house scope narrows the noun:
// "your opponent cannot use creatures of the chosen house to reap during their
// next turn" (Seismo-entangler).
func (e Restrict) Text() string {
	who, whose := e.Player.secondPerson()
	when := "during " + whose + " next turn"
	if e.Duration == RemainderOfPlayerTurn {
		when = "for the remainder of the turn"
	}
	phrase := e.Action.phrase()
	if h := e.House.phrase(); h != "" {
		phrase = "creatures " + h + " to reap"
	}
	return who + " cannot use " + phrase + " " + when
}

// Resolve applies the timed bar to the chosen player, dispatching to the gate for
// the action, for the current turn (RemainderOfPlayerTurn) or the player's next
// turn (OpponentNextTurn). validate guarantees fighting only reaches next turn and
// a house scope only reaches a next-turn reaping bar.
func (e Restrict) Resolve(ctx *EffectContext) {
	player := ctx.PlayerFor(e.Player)
	switch e.Action {
	case RestrictFighting:
		ctx.Resolver.CannotFightNextTurn(player, ctx.Source)
	case RestrictReaping:
		switch {
		case e.House.phrase() != "":
			ctx.Resolver.CannotReapHouseNextTurn(player, e.House.resolveHouse(ctx), ctx.Source)
		case e.Duration == RemainderOfPlayerTurn:
			ctx.Resolver.CannotReapThisTurn(player, ctx.Source)
		default:
			ctx.Resolver.CannotReapNextTurn(player, ctx.Source)
		}
	case RestrictUse:
		if e.Duration == RemainderOfPlayerTurn {
			ctx.Resolver.CannotUseThisTurn(player, ctx.Source)
		} else {
			ctx.Resolver.CannotUseNextTurn(player, ctx.Source)
		}
	}
}

// CreaturesCannot bars every creature in play — both players' — from being used
// one way (fighting or reaping) until the start of the caster's next turn, save
// for the houses Houses excludes. Where CannotFight and CannotReap bar one
// player's use of their own creatures, this is a rule on the whole board: Into
// the Night stops non-Shadows creatures fighting, Sow Salt stops every creature
// reaping. Houses admits the barred creatures — left unset (MatchAnyHouse) it
// bars every house.
type CreaturesCannot struct {
	Action   UseKind
	Houses   HouseMatcher
	Duration Duration
}

// validate rejects a CreaturesCannot that bars no real action or names no
// duration. The action must be fighting or reaping; an "Action:" ability cannot
// be barred this way. Houses may only name a static house set (any, a named
// house, or all but a named house): a board-wide bar rides in flat state with no
// bound context, so a chosen/active/contextual matcher has nothing to resolve
// against when the gate is later checked.
func (e CreaturesCannot) validate() error {
	if e.Action != FightUse && e.Action != ReapUse {
		return fmt.Errorf("CreaturesCannot: action must be FightUse or ReapUse")
	}
	switch e.Houses.Kind {
	case MatchAnyHouse, MatchNamedHouse, MatchExceptHouse:
	default:
		return fmt.Errorf("CreaturesCannot: houses must be a static house set")
	}
	if err := e.Houses.validate(); err != nil {
		return err
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
	return "until the start of your next turn, " +
		e.Houses.qualify("creatures") +
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
			e.Houses,
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
// next turn". Card types read as lowercase common nouns in card text.
func (e CannotPlay) Text() string {
	who, whose := e.Player.secondPerson()
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

// PlayersCannotPlay bars both players from playing cards of a Type until the end
// of the caster's next turn — Stealth Mode stops either player playing Tactics.
// Where CannotPlay bars one named player, this is a board-wide rule spanning both
// sides: it arms the bar on the caster for the rest of this turn and their next
// turn, and on the opponent throughout their next turn, so it lifts once the
// caster's next turn ends. An unset Type bars every type.
type PlayersCannotPlay struct {
	Type     CardType
	Duration Duration
}

// barred is the type this bar installs: the named one, or the AnyType wildcard
// when the author left Type unset to bar everything.
func (e PlayersCannotPlay) barred() CardType {
	if e.Type == TypeUnset {
		return AnyType
	}
	return e.Type
}

// validate rejects a PlayersCannotPlay whose duration is not the only span it
// models. An unset Type is allowed and means every type.
func (e PlayersCannotPlay) validate() error {
	if e.Duration != EndOfPlayerNextTurn {
		return fmt.Errorf("PlayersCannotPlay: duration must be EndOfPlayerNextTurn")
	}
	return nil
}

// Text renders the effect, e.g. "until the end of your next turn, players cannot
// play tactics". Card types read as lowercase common nouns in card text.
func (e PlayersCannotPlay) Text() string {
	noun := strings.ToLower(e.barred().String()) + "s"
	return "until the end of your next turn, players cannot play " + noun
}

// Resolve arms the play-type bar on both players: the caster for the rest of this
// turn and their next turn, the opponent throughout their next turn.
func (e PlayersCannotPlay) Resolve(ctx *EffectContext) {
	t := e.barred()
	ctx.Resolver.CannotPlayTypeThisTurn(ctx.Controller, t, ctx.Source)
	ctx.Resolver.CannotPlayTypeNextTurn(ctx.Controller, t, ctx.Source)
	ctx.Resolver.CannotPlayTypeNextTurn(ctx.Opponent(), t, ctx.Source)
}

// HouseChoiceReference names where an active-house constraint reads the house it
// applies (and, for JustChosen, whose next turn it binds): the house an enclosing
// ChooseHouseThen picked (Chosen), the house of the creature the source fought
// (Fought), or the house a player just chose this turn, read from the board
// (JustChosen — Snag's Mirror, keyed off the "after a player chooses a house"
// trigger).
type HouseChoiceReference uint8

const (
	houseChoiceRefUnset HouseChoiceReference = iota
	// ChosenActiveHouse reads the house an enclosing ChooseHouseThen picked.
	ChosenActiveHouse
	// FoughtActiveHouse reads the house of the creature the source fought (ctx.It).
	FoughtActiveHouse
	// JustChosenActiveHouse reads the house a player just chose, from the board.
	JustChosenActiveHouse
)

// MustChooseHouse forces a player to choose a house as their active house on their
// next turn — the house an enclosing ChooseHouseThen picked (Control the Weak,
// Reference Chosen) or the house of the creature the source fought (Snag, Reference
// Fought). Player names whose next-turn choice is forced; it is source-relative, so
// a card binding its own controller reads "you must choose … on your next turn".
type MustChooseHouse struct {
	// Player is whose next-turn active-house choice is forced.
	Player Player
	// Reference names where the forced house is read from (Chosen or Fought).
	Reference HouseChoiceReference
}

// validate requires a player and a Chosen or Fought reference.
func (e MustChooseHouse) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("MustChooseHouse")
	}
	if e.Reference != ChosenActiveHouse && e.Reference != FoughtActiveHouse {
		return fmt.Errorf("MustChooseHouse: reference must be Chosen or Fought")
	}
	return nil
}

// Text renders the effect in the player's voice, e.g. "your opponent must choose
// that house as their active house during their next turn".
func (e MustChooseHouse) Text() string {
	who, whose := e.Player.secondPerson()
	if e.Reference == FoughtActiveHouse {
		return who + " must choose the house of the creature " + SelfName +
			" fights as " + whose + " active house on " + whose + " next turn"
	}
	return who + " must choose that house as " + whose +
		" active house during " + whose + " next turn"
}

// Resolve arms the must on the player's next turn. A Fought reference stores a live
// reference to the fought creature, resolved to its current house at choice time,
// and does nothing when no creature is in context.
func (e MustChooseHouse) Resolve(ctx *EffectContext) {
	player := ctx.PlayerFor(e.Player)
	if e.Reference == FoughtActiveHouse {
		if !ctx.HasIt {
			return
		}
		ctx.Resolver.MustChooseFoughtHouseNextTurn(player, ctx.It, ctx.Source)
		return
	}
	ctx.Resolver.MustChooseHouseNextTurn(player, ctx.ChosenHouse, ctx.Source)
}

// CannotChooseHouse bars a player from choosing a house as their active house on
// their next turn — the chosen house (Tezmal, Reference Chosen) or the house a
// player just chose this turn (Snag's Mirror, Reference JustChosen, keyed off the
// "after a player chooses a house" trigger). Player names whose next-turn choice is
// barred; for a Chosen reference it is source-relative, and for a JustChosen
// reference it is relative to the player who just chose (read from the board), so
// Player Opponent bars the chooser's opponent whichever player made the choice.
type CannotChooseHouse struct {
	// Player is whose next-turn active-house choice is barred.
	Player Player
	// Reference names where the barred house is read from (Chosen or JustChosen).
	Reference HouseChoiceReference
}

// validate requires a player and a Chosen or JustChosen reference.
func (e CannotChooseHouse) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("CannotChooseHouse")
	}
	if e.Reference != ChosenActiveHouse && e.Reference != JustChosenActiveHouse {
		return fmt.Errorf("CannotChooseHouse: reference must be Chosen or JustChosen")
	}
	return nil
}

// Text renders the effect. A JustChosen reference follows the "after a player
// chooses a house" trigger and speaks from the chooser's point of view.
func (e CannotChooseHouse) Text() string {
	if e.Reference == JustChosenActiveHouse {
		return "their opponent cannot choose the same house as their active house on their next turn"
	}
	who, whose := e.Player.secondPerson()
	return who + " cannot choose that house as " + whose +
		" active house on " + whose + " next turn"
}

// Resolve bars the house on the barred player's next turn. A JustChosen reference
// reads the choosing player and house from the board, applying Player relative to
// the chooser, so Player Opponent bars the chooser's opponent whichever player made
// the choice.
func (e CannotChooseHouse) Resolve(ctx *EffectContext) {
	if e.Reference == JustChosenActiveHouse {
		chooser := ctx.Resolver.ActivePlayer()
		barred := chooser
		if e.Player == Opponent {
			barred = 1 - chooser
		}
		ctx.Resolver.CannotChooseHouseNextTurn(barred, ctx.Resolver.ActiveHouse(), ctx.Source)
		return
	}
	ctx.Resolver.CannotChooseHouseNextTurn(ctx.PlayerFor(e.Player), ctx.ChosenHouse, ctx.Source)
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

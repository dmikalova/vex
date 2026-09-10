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
	if e.Duration == EndOfTurn {
		return who + " cannot use creatures to reap for the remainder of the turn"
	}
	return who + " cannot use creatures to reap during " + whose + " next turn"
}

// Resolve applies the timed bar to the chosen player, for the current turn
// (EndOfTurn) or the player's next turn (NextTurn).
func (e CannotReap) Resolve(ctx *EffectContext) {
	switch e.Duration {
	case EndOfTurn:
		ctx.Resolver.CannotReapThisTurn(ctx.PlayerFor(e.Player), ctx.Source)
	case NextTurn:
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
	if e.Duration == NextTurn {
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
// next turn". The renamed Tactic type prints capitalized as "Tactics" (rule 19).
func (e CannotPlay) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	noun := strings.ToLower(e.barred().String()) + "s"
	if e.Type == Tactic {
		noun = "Tactics"
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
// next turn", or for the EndOfTurn form "you cannot use cards this turn" (United
// Action).
func (e CannotUse) Text() string {
	who, whose := "you", "your"
	if e.Player == Opponent {
		who, whose = "your opponent", "their"
	}
	if e.Duration == EndOfTurn {
		return who + " cannot use cards this turn"
	}
	return who + " cannot use any cards during " + whose + " next turn"
}

// Resolve arms the use bar on the chosen player, this turn or their next.
func (e CannotUse) Resolve(ctx *EffectContext) {
	switch e.Duration {
	case NextTurn:
		ctx.Resolver.CannotUseNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
	case EndOfTurn:
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

// SkipForgePhase makes a player skip their "forge a key" phase at the start of their
// next turn (Miasma).
type SkipForgePhase struct {
	Player Player
}

// validate rejects a SkipForgePhase whose player was left unset.
func (e SkipForgePhase) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("SkipForgePhase")
	}
	return nil
}

// Text renders the effect, e.g. `your opponent skips the "forge a key" phase during
// their next turn`.
func (e SkipForgePhase) Text() string {
	who, whose, verb := "you", "your", "skip"
	if e.Player == Opponent {
		who, whose, verb = "your opponent", "their", "skips"
	}
	return fmt.Sprintf("%s %s the %q phase during %s next turn", who, verb, "forge a key", whose)
}

// Resolve arms the skip on the chosen player's next turn.
func (e SkipForgePhase) Resolve(ctx *EffectContext) {
	ctx.Resolver.SkipForgePhaseNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
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

// GrantFightForFriendlyHouse lets the controller's creatures of a fixed House fight
// this turn even out of the active house — Signal Fire's "friendly Brobnar creatures
// may fight as though they belonged to the active house." Unlike
// GrantFightForChosenHouse, which reads the house from an enclosing ChooseHouseThen,
// this names the house on the card. The grant lasts only the current turn (the ready
// phase clears it).
type GrantFightForFriendlyHouse struct {
	House House
}

// validate rejects a GrantFightForFriendlyHouse whose house was left unset.
func (e GrantFightForFriendlyHouse) validate() error {
	if e.House == HouseNone {
		return fmt.Errorf("GrantFightForFriendlyHouse: house must be set")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, each friendly
// Brobnar creature may fight".
func (e GrantFightForFriendlyHouse) Text() string {
	return fmt.Sprintf(
		"for the remainder of the turn, each friendly %s creature may fight",
		e.House,
	)
}

// Resolve grants the controller's House creatures the right to fight this turn.
func (e GrantFightForFriendlyHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.GrantFightForHouse(ctx.Controller, e.House)
}

// HouseGrant is a bitset of what a MayActFriendlyHouse frees for one house this
// turn: playing the house's cards from hand, using its creatures in play, or both.
type HouseGrant uint8

const (
	// GrantPlay lets the controller play the house's cards from hand.
	GrantPlay HouseGrant = 1 << iota
	// GrantUse lets the controller use (fight, reap, or Action:) the house's creatures.
	GrantUse
)

// MayActFriendlyHouse lets the controller act with a friendly house's cards this
// turn even out of the active house: Grant frees playing the house's cards from
// hand (GrantPlay), using its creatures in play (GrantUse), or both — the
// Ambassador cycle plays and uses, Sigil of Brotherhood and Ritual of the Hunt
// only use. The grant lasts only the current turn (the ready phase clears it).
type MayActFriendlyHouse struct {
	House House
	Grant HouseGrant
}

// validate rejects an unset house or an empty grant.
func (e MayActFriendlyHouse) validate() error {
	if e.House == HouseNone {
		return fmt.Errorf("MayActFriendlyHouse: house must be set")
	}
	if e.Grant == 0 {
		return fmt.Errorf("MayActFriendlyHouse: at least one grant must be set")
	}
	return nil
}

// Text builds the clause from the grant, e.g. "for the remainder of the turn, you
// may play or use a Mars card" or "... you may use friendly Sanctum creatures".
// Using frees creatures in play; playing frees cards from hand, so the object is a
// card once GrantPlay is set.
func (e MayActFriendlyHouse) Text() string {
	verb, object := "use", "friendly "+e.House.String()+" creatures"
	if e.Grant&GrantPlay != 0 {
		verb, object = "play", "a "+e.House.String()+" card"
		if e.Grant&GrantUse != 0 {
			verb = "play or use"
		}
	}
	return "for the remainder of the turn, you may " + verb + " " + object
}

// Resolve applies each grant the effect names for the controller this turn.
func (e MayActFriendlyHouse) Resolve(ctx *EffectContext) {
	if e.Grant&GrantPlay != 0 {
		ctx.Resolver.GrantPlayForHouse(ctx.Controller, e.House)
	}
	if e.Grant&GrantUse != 0 {
		ctx.Resolver.GrantUseForHouse(ctx.Controller, e.House)
	}
}

// MayUseFriendlyArtifacts lets the controller use any friendly artifact this turn
// as if it belonged to the active house — Scientifical Hack. Unlike
// MayActFriendlyHouse, which frees one named house's cards, this frees every
// friendly artifact whatever its house. The grant lasts only the current turn (the
// ready phase clears it).
type MayUseFriendlyArtifacts struct{}

// Text renders the effect.
func (MayUseFriendlyArtifacts) Text() string {
	return "for the remainder of the turn, you may use friendly artifacts as if they belonged to the active house"
}

// Resolve grants the controller use of every friendly artifact this turn.
func (MayUseFriendlyArtifacts) Resolve(ctx *EffectContext) {
	ctx.Resolver.GrantUseArtifactsAnyHouse(ctx.Controller)
}

// MayPlayOffHouse lets the controller play or use a bounded number of cards this
// turn from outside their active house — the Star Alliance "play a non-Star
// Alliance card" cycle. Where MayActFriendlyHouse frees one named friendly house,
// this frees cards by exclusion: every house but a named one (Except, a card's own
// house for "non-Star Alliance"), or every house the controller has a card in play
// for (Controlled, United Action). NotType excludes a card type (Com. Officer
// Kirby frees a non-creature). Grant frees playing from hand (GrantPlay), using in
// play (GrantUse), or both (CXO Taber). Count bounds how many cards; zero is
// unbounded (United Action). The ready phase clears the grant.
type MayPlayOffHouse struct {
	Except     House
	Controlled bool
	NotType    CardType
	Grant      HouseGrant
	Count      int
}

// validate rejects a grant that frees nothing or names a negative count.
func (e MayPlayOffHouse) validate() error {
	if e.Grant == 0 {
		return fmt.Errorf("MayPlayOffHouse: at least one grant must be set")
	}
	if e.Count < 0 {
		return fmt.Errorf("MayPlayOffHouse: count must not be negative")
	}
	return nil
}

// Text renders the grant. The Controlled form (United Action) names the whole
// permission; the exclusion form reads "you may play a non-Star Alliance ... this
// turn", playing "or use" when GrantUse is set (CXO Taber).
func (e MayPlayOffHouse) Text() string {
	if e.Controlled {
		return "for the remainder of the turn, you may play cards from any house for which you have a card in play"
	}
	verb := "play"
	if e.Grant&GrantUse != 0 {
		verb = "play or use"
	}
	return "you may " + verb + " " + e.object() + " this turn"
}

// object renders the card the grant frees: a bare "card", or the enumerated
// non-creature types Com. Officer Kirby frees. The renamed Tactic type prints
// capitalized (rule 19).
func (e MayPlayOffHouse) object() string {
	house := ""
	if e.Except != HouseNone {
		house = "non-" + e.Except.String() + " "
	}
	if e.NotType == Creature {
		return "a " + house + "artifact, upgrade, or Tactic"
	}
	return "one " + house + "card"
}

// Resolve records the this-turn off-house grant for the controller.
func (e MayPlayOffHouse) Resolve(ctx *EffectContext) {
	rem := permitUnlimited
	if e.Count > 0 {
		rem = uint8(e.Count)
	}
	ctx.Resolver.GrantOffHousePermit(ctx.Controller, OffHousePermit{
		Except:     e.Except,
		Controlled: e.Controlled,
		NotType:    e.NotType,
		Grant:      e.Grant,
		Remaining:  rem,
	})
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

// ForbidOpponentActiveHouse bars the opponent from choosing the house picked by
// an enclosing ChooseHouseThen as their active house on their next turn — Tezmal's
// "your opponent cannot choose that house as their active house on their next
// turn."
type ForbidOpponentActiveHouse struct{}

// Text renders the effect.
func (ForbidOpponentActiveHouse) Text() string {
	return "your opponent cannot choose that house as their active house on their next turn"
}

// Resolve arms the forbidden house on the opponent's next turn.
func (ForbidOpponentActiveHouse) Resolve(ctx *EffectContext) {
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

// ForceOpponentActiveHouseOfFought is Snag's Fight ability: the opponent must
// choose the house of the creature Snag fought (ctx.It) as their active house on
// their next turn.
type ForceOpponentActiveHouseOfFought struct{}

// Text renders the effect.
func (ForceOpponentActiveHouseOfFought) Text() string {
	return "your opponent must choose the house of the creature " + SelfName +
		" fights as their active house on their next turn"
}

// Resolve arms the fought creature's house on the opponent's next turn.
func (ForceOpponentActiveHouseOfFought) Resolve(ctx *EffectContext) {
	if !ctx.HasIt {
		return
	}
	ctx.Resolver.ForceActiveHouseNextTurn(ctx.Opponent(), ctx.Resolver.House(ctx.It), ctx.Source)
}

// ForbidSameActiveHouseNextTurn is Snag's Mirror: keyed off the shared "after a
// player chooses a house" window, it bars whoever chose from having their
// opponent match that house next turn. It reads the active player and house from
// the board rather than the source's point of view, so it works whichever player
// made the choice.
type ForbidSameActiveHouseNextTurn struct{}

// Text renders the effect from the chooser's point of view, to follow the "after
// a player chooses an active house, " trigger prefix.
func (ForbidSameActiveHouseNextTurn) Text() string {
	return "their opponent cannot choose the same house as their active house on their next turn"
}

// Resolve bars the chooser's opponent from the chosen house on their next turn.
func (ForbidSameActiveHouseNextTurn) Resolve(ctx *EffectContext) {
	chooser := ctx.Resolver.ActivePlayer()
	ctx.Resolver.ForbidActiveHouseNextTurn(1-chooser, ctx.Resolver.ActiveHouse(), ctx.Source)
}

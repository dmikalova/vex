package engine

import (
	"fmt"
	"strings"
)

// A conditional gates an effect behind a check on the current game state — the
// "If ..." clause a card opens with, e.g. "If your opponent has 7 or more Æmber,
// they lose 4 Æmber." The effect resolves only when the condition is met. Unlike
// a result gate (A -> B), which turns on an action succeeding, a conditional turns
// on a fact about the board.
// Condition is a boolean predicate on the live game, used by Conditional.
type Condition interface {
	CondText() string
	Met(ctx *EffectContext) bool
}

// Comparison selects how a condition compares a quantity to a threshold. It has
// no valid zero value: a condition must name one, so an unset comparison is
// caught at init rather than silently reading as "0 or more".
type Comparison int

const (
	// comparisonUnset is the invalid zero value; a condition must name a real
	// comparison, so an unset one is caught at init.
	comparisonUnset Comparison = iota
	// AtLeast is met when the quantity is at least Amount.
	AtLeast
	// AtMost is met when the quantity is at most Amount.
	AtMost
	// Exactly is met when the quantity is exactly Amount (which is why the
	// comparison is named separately from the amount: Exactly with Amount 0 is a
	// real check that a bare integer field could not tell from "unset").
	Exactly
	// MoreThanYou is met when the opponent's pool holds strictly more Æmber than the
	// controller's; it ignores Amount and applies only to PoolAember{Player: Opponent}.
	MoreThanYou
	// MoreThanOpponent is met when the controller's pool holds strictly more Æmber
	// than the opponent's; it ignores Amount and applies only to
	// PoolAember{Player: Controller}.
	MoreThanOpponent
)

// PoolAember gates on one player's Æmber pool: Player names whose pool (Controller
// or Opponent), Is the comparison, and Amount the threshold it compares against
// (unused by the relative MoreThanYou / MoreThanOpponent comparisons, which compare
// the two pools). It replaces the mirror opponent-pool / your-pool conditions with
// one node.
type PoolAember struct {
	Player Player
	Is     Comparison
	Amount int
}

// validate requires a pool-owning player and a named comparison, and ties each
// relative comparison to the side it reads from.
func (c PoolAember) validate() error {
	if c.Player != Controller && c.Player != Opponent {
		return fmt.Errorf("PoolAember: Player must be Controller or Opponent")
	}
	switch c.Is {
	case AtLeast, AtMost, Exactly:
		return nil
	case MoreThanYou:
		if c.Player != Opponent {
			return fmt.Errorf("PoolAember: MoreThanYou requires Player Opponent")
		}
		return nil
	case MoreThanOpponent:
		if c.Player != Controller {
			return fmt.Errorf("PoolAember: MoreThanOpponent requires Player Controller")
		}
		return nil
	default:
		return fmt.Errorf(
			"PoolAember: Is must be AtLeast, AtMost, Exactly, MoreThanYou, or MoreThanOpponent",
		)
	}
}

// CondText renders the condition, e.g. "if your opponent has 7 Æmber or more" or
// "if you have more Æmber than your opponent".
func (c PoolAember) CondText() string {
	switch c.Is {
	case MoreThanYou:
		return "if your opponent has more Æmber than you"
	case MoreThanOpponent:
		return "if you have more Æmber than your opponent"
	}
	if c.Player == Opponent {
		switch {
		case c.Is == Exactly && c.Amount == 0:
			return "if your opponent has no Æmber"
		case c.Is == Exactly:
			return fmt.Sprintf("if your opponent has exactly %d Æmber", c.Amount)
		case c.Is == AtMost:
			return fmt.Sprintf("if your opponent has %d Æmber or fewer", c.Amount)
		default:
			return fmt.Sprintf("if your opponent has %d Æmber or more", c.Amount)
		}
	}
	switch {
	case c.Is == Exactly && c.Amount == 0:
		return "if you have no Æmber"
	case c.Is == Exactly:
		return fmt.Sprintf("if you have exactly %d Æmber", c.Amount)
	case c.Is == AtMost:
		return fmt.Sprintf("if you have %d Æmber or fewer", c.Amount)
	default:
		return fmt.Sprintf("if you have %d Æmber or more", c.Amount)
	}
}

// Met reports whether the named player's pool satisfies the comparison.
func (c PoolAember) Met(ctx *EffectContext) bool {
	mine := ctx.Resolver.Aember(ctx.PlayerFor(c.Player))
	switch c.Is {
	case Exactly:
		return mine == c.Amount
	case AtMost:
		return mine <= c.Amount
	case MoreThanYou, MoreThanOpponent:
		other := ctx.Controller
		if c.Player == Controller {
			other = ctx.Opponent()
		}
		return mine > ctx.Resolver.Aember(other)
	default:
		return mine >= c.Amount
	}
}

// validateCondition returns any configuration error a condition reports (an unset
// PoolAember comparison, say). Conditions that cannot be misconfigured
// implement no validator and pass.
func validateCondition(c Condition) error {
	if v, ok := c.(validator); ok {
		return v.validate()
	}
	return nil
}

// ControlsMoreCreatures is met while the controller has more creatures in play
// than the opponent.
type ControlsMoreCreatures struct{}

// CondText renders the condition.
func (ControlsMoreCreatures) CondText() string {
	return "if you control more creatures than your opponent"
}

// Met reports whether the controller has more creatures in play than the opponent.
func (ControlsMoreCreatures) Met(ctx *EffectContext) bool {
	return len(
		ctx.Resolver.Battleline(ctx.Controller),
	) > len(
		ctx.Resolver.Battleline(ctx.Opponent()),
	)
}

// SourceOnFlank is met by the position of the source card: with Not false when it
// is on a flank of its battleline, with Not true when it is not — Glyxl
// Proliferator archives only while on a flank, Titan Librarian only while not.
type SourceOnFlank struct {
	Not bool
}

// CondText renders the condition naming the source card.
func (c SourceOnFlank) CondText() string {
	if c.Not {
		return "if " + SelfName + " is not on a flank"
	}
	return "if " + SelfName + " is on a flank"
}

// Met reports whether the source card's flank position matches Not.
func (c SourceOnFlank) Met(ctx *EffectContext) bool {
	return onFlank(ctx, ctx.Source) != c.Not
}

// HasOtherFriendlyCreatures is met when the source's controller has at least one
// creature in play other than the source — Reassembling Automaton replaces its own
// destruction only "if you have any other creatures in play".
type HasOtherFriendlyCreatures struct{}

// CondText renders the condition as the card prints it.
func (HasOtherFriendlyCreatures) CondText() string {
	return "if you have any other creatures in play"
}

// Met reports whether the controller has any creature in play besides the source.
func (HasOtherFriendlyCreatures) Met(ctx *EffectContext) bool {
	return InPlay{Player: Controller, Type: Creature, Other: true}.Met(ctx)
}

// CardsInDeckAtMost is met when the controller's deck holds at most Amount cards —
// Manchego steals only "if you have 5 or fewer cards in your deck".
type CardsInDeckAtMost struct {
	Amount int
}

// CondText renders the condition as the card prints it.
func (e CardsInDeckAtMost) CondText() string {
	return fmt.Sprintf("if you have %d or fewer cards in your deck", e.Amount)
}

// Met reports whether the controller's deck size is at most Amount.
func (e CardsInDeckAtMost) Met(ctx *EffectContext) bool {
	return len(ctx.Resolver.Deck(ctx.Controller)) <= e.Amount
}

// ControlsNamed is met when the controller has a card of a given printed name in
// play — Hyde draws an extra card while it controls Velum.
type ControlsNamed struct {
	Name string
}

// CondText renders the condition, e.g. "if you control Velum".
func (c ControlsNamed) CondText() string {
	return "if you control " + c.Name
}

// Met reports whether the controller has the named card in play.
func (c ControlsNamed) Met(ctx *EffectContext) bool {
	return InPlay{Player: Controller, Name: c.Name}.Met(ctx)
}

// ItIsOnFlank is met when the card in context (ctx.It — the creature a preceding
// effect put in focus) sits on a flank of its battleline. Malison moves an enemy
// creature, then, if it is on a flank, has it capture. With no context creature it
// is not met.
type ItIsOnFlank struct{}

// CondText renders the condition naming the context creature.
func (ItIsOnFlank) CondText() string { return "if it is on a flank" }

// Met reports whether the context creature is on a flank.
func (ItIsOnFlank) Met(ctx *EffectContext) bool {
	return ctx.HasIt && onFlank(ctx, ctx.It)
}

// ItIsOnNamedFlank is met when the context creature (ctx.It) sits on one specific
// flank — the left end of its controller's battleline, or the right when Right is
// set. Sinestra keys off a creature the opponent plays on their left flank, Dexus
// off the right. A lone creature is on both flanks, so it meets either side.
type ItIsOnNamedFlank struct{ Right bool }

// CondText renders the condition naming which flank.
func (c ItIsOnNamedFlank) CondText() string {
	if c.Right {
		return "if it is on the right flank"
	}
	return "if it is on the left flank"
}

// Met reports whether the context creature is the end creature of its controller's
// battleline on the named side.
func (c ItIsOnNamedFlank) Met(ctx *EffectContext) bool {
	if !ctx.HasIt || !ctx.Resolver.IsCreature(ctx.It) {
		return false
	}
	line := ctx.Resolver.Battleline(ctx.Resolver.Controller(ctx.It))
	if len(line) == 0 {
		return false
	}
	if c.Right {
		return line[len(line)-1] == ctx.It
	}
	return line[0] == ctx.It
}

// SourceInCenterOfBattleline is met while the source card sits in the center of
// its controller's battleline — the middle creature of an odd-sized line, with
// equal creatures to its left and right. An even-sized line has no center.
type SourceInCenterOfBattleline struct{}

// CondText renders the condition naming the source card.
func (SourceInCenterOfBattleline) CondText() string {
	return "if " + SelfName + " is in the center of your battleline"
}

// Met reports whether the source card sits in the center of its battleline.
func (SourceInCenterOfBattleline) Met(ctx *EffectContext) bool {
	return ctx.Resolver.InCenterOfBattleline(ctx.Source)
}

// SourceReady is met while the source card is ready (unexhausted) — Bellowing
// Patrizate damages each creature that enters play only while it is ready.
type SourceReady struct{}

// CondText renders the condition naming the source card.
func (SourceReady) CondText() string { return "if " + SelfName + " is ready" }

// Met reports whether the source card is currently ready.
func (SourceReady) Met(ctx *EffectContext) bool {
	return !ctx.Resolver.Exhausted(ctx.Source)
}

// SourceNeighborsAllOfHouse is met while every battleline neighbor of the source
// card belongs to House — Xanthyx Harvester cannot be used while it has a
// non-Mars neighbor, so its use is gated on this being met.
type SourceNeighborsAllOfHouse struct {
	House House
}

// CondText renders the condition, e.g. "if it has no non-Mars neighbor".
func (c SourceNeighborsAllOfHouse) CondText() string {
	return "if it has no non-" + c.House.String() + " neighbor"
}

// Met reports whether the source card has no neighbor off House.
func (c SourceNeighborsAllOfHouse) Met(ctx *EffectContext) bool {
	for _, n := range neighbors(ctx, ctx.Source) {
		if ctx.Resolver.House(n) != c.House {
			return false
		}
	}
	return true
}

// ArchivedCreaturesShareHouse is met when the creatures a preceding
// ArchiveFromPlay set aside (ctx.Produced.Archived) all belong to one house —
// Code Monkey gains 2 Æmber when the neighbors it archived share a house. Fewer
// than two creatures cannot share a house, so it is not met.
type ArchivedCreaturesShareHouse struct{}

// CondText renders the condition naming the just-archived creatures.
func (ArchivedCreaturesShareHouse) CondText() string {
	return "if those creatures share a house"
}

// Met reports whether every archived creature belongs to the same house.
func (ArchivedCreaturesShareHouse) Met(ctx *EffectContext) bool {
	ids := ctx.Produced.Archived
	if len(ids) < 2 {
		return false
	}
	first := ctx.Resolver.House(ids[0])
	for _, id := range ids[1:] {
		if ctx.Resolver.House(id) != first {
			return false
		}
	}
	return true
}

// ControlsCreaturesOfHouses is met while the controller's creatures in play span
// at least Amount different houses — Prince Derric, Unifier pays out when three
// houses are represented.
type ControlsCreaturesOfHouses struct {
	// Amount is the number of different houses that must be represented.
	Amount int
}

// validate requires a positive Count.
func (c ControlsCreaturesOfHouses) validate() error {
	if c.Amount <= 0 {
		return fmt.Errorf("ControlsCreaturesOfHouses: Count must be positive")
	}
	return nil
}

// CondText renders the condition, e.g. "if you control creatures from 3 or more
// houses".
func (c ControlsCreaturesOfHouses) CondText() string {
	return fmt.Sprintf("if you control creatures from %d or more houses", c.Amount)
}

// Met reports whether the controller's creatures span at least Amount houses.
func (c ControlsCreaturesOfHouses) Met(ctx *EffectContext) bool {
	seen := map[House]bool{}
	for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
		seen[ctx.Resolver.House(id)] = true
	}
	return len(seen) >= c.Amount
}

// FirstCreaturePlayedThisTurn is met when the card in context (ctx.It, the
// creature that fired the trigger) is the first creature its player played this
// turn — Speed Sigil readies it. It is a once-per-turn charge that needs no state
// of its own: the turn's play record already says whether the charge is spent, and
// the record is cleared when the next turn begins.
//
// A creature put into play by an effect rather than played never matches, so it
// neither benefits nor spends the charge.
type FirstCreaturePlayedThisTurn struct{}

// CondText renders the condition.
func (FirstCreaturePlayedThisTurn) CondText() string {
	return "if it is the first creature played this turn"
}

// Met reports whether the context card is the earliest creature in the active
// player's plays this turn.
func (FirstCreaturePlayedThisTurn) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	for _, id := range ctx.Resolver.PlayedThisTurn(ctx.Resolver.ActivePlayer()) {
		if ctx.Resolver.TypeOf(id) == Creature {
			return id == ctx.It
		}
	}
	return false
}

// NoCreaturesPlayedThisTurn is met when the controller has not played any
// creatures during the current turn — Redlock pays out at end of turn only on a
// turn its controller played no creatures. A creature put into play by an effect
// rather than played does not count, so it does not spoil the payout.
type NoCreaturesPlayedThisTurn struct{}

// CondText renders the condition.
func (NoCreaturesPlayedThisTurn) CondText() string {
	return "if you did not play any creatures this turn"
}

// Met reports whether none of the controller's plays this turn were creatures.
func (NoCreaturesPlayedThisTurn) Met(ctx *EffectContext) bool {
	for _, id := range ctx.Resolver.PlayedThisTurn(ctx.Controller) {
		if ctx.Resolver.TypeOf(id) == Creature {
			return false
		}
	}
	return true
}

// ItIsYourTurn is met when the ability's controller is the active player —
// Jargogle plays the card under it when destroyed on its controller's turn, and
// archives it otherwise.
type ItIsYourTurn struct{}

// CondText renders the condition.
func (ItIsYourTurn) CondText() string { return "if it is your turn" }

// Met reports whether the controller is the active player.
func (ItIsYourTurn) Met(ctx *EffectContext) bool {
	return ctx.Resolver.ActivePlayer() == ctx.Controller
}

// AemberOnThisAtLeast is met when at least Amount Æmber sits on the source card —
// [REDACTED] sacrifices itself once it has hoarded four or more.
type AemberOnThisAtLeast struct {
	Amount int
}

// CondText renders the condition clause.
func (c AemberOnThisAtLeast) CondText() string {
	return fmt.Sprintf("if there are %d or more Æmber on it", c.Amount)
}

// Met reports whether the source card holds at least Amount Æmber.
func (c AemberOnThisAtLeast) Met(ctx *EffectContext) bool {
	return ctx.Resolver.AmberOn(ctx.Source) >= c.Amount
}

// CardsDestroyedFewerThan is met when fewer than Amount cards were destroyed this
// way — the tally a preceding effect records on the context. Bonkers Killing
// Machine destroys itself when its house-driven destruction removed fewer than two.
type CardsDestroyedFewerThan struct {
	Amount int
}

// CondText renders the condition, e.g. "if fewer than 2 cards are destroyed this
// way".
func (c CardsDestroyedFewerThan) CondText() string {
	return fmt.Sprintf("if fewer than %d cards are destroyed this way", c.Amount)
}

// Met reports whether fewer than Amount cards were destroyed this way.
func (c CardsDestroyedFewerThan) Met(ctx *EffectContext) bool {
	return ctx.Produced.TotalDestroyed() < c.Amount
}

// HouseChoice names a house a condition compares against by reference rather than
// by a fixed value — the house chosen this turn, or the active house — so one
// condition works wherever such a house is meaningful.
type HouseChoice uint8

const (
	// houseChoiceUnset is the invalid zero value.
	houseChoiceUnset HouseChoice = iota
	// TheChosenHouse is the house picked by an enclosing ChooseHouseThen.
	TheChosenHouse
	// TheActiveHouse is the player's active house this turn.
	TheActiveHouse
	// TheContextualHouse is the house of the card in context (ctx.It) — e.g. a
	// just-discarded deck card, for "of the discarded card's house".
	TheContextualHouse
	// AnyHouse names no house at all, so a count or condition using it applies no
	// house filter — Key Abduction counts every card in hand, whatever its house.
	AnyHouse
)

// resolveHouse turns a HouseChoice into the concrete house it names in the current
// context (HouseNone when it names a card in context and none is present).
func (h HouseChoice) resolveHouse(ctx *EffectContext) House {
	switch h {
	case TheActiveHouse:
		return ctx.Resolver.ActiveHouse()
	case TheContextualHouse:
		if ctx.HasIt {
			return ctx.Resolver.House(ctx.It)
		}
		return HouseNone
	default:
		return ctx.ChosenHouse
	}
}

// ItIsOfHouse is met when the card in context (ctx.It — a revealed, discarded, or
// triggering card) belongs to a referenced house. It replaces the one-off
// "revealed card of the chosen house" (Chaos Portal) and "discarded card of the
// active house" (Evasion Sigil) with a single filter on the contextual card.
type ItIsOfHouse struct {
	House HouseChoice
}

// CondText renders the condition, e.g. "if it is of the chosen house".
func (e ItIsOfHouse) CondText() string {
	if e.House == TheActiveHouse {
		return "if it is of the active house"
	}
	return "if it is of the chosen house"
}

// Met reports whether a card is in context and belongs to the referenced house.
func (e ItIsOfHouse) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	house := ctx.Resolver.House(ctx.It)
	if e.House == TheActiveHouse {
		return house == ctx.Resolver.ActiveHouse()
	}
	return house == ctx.ChosenHouse
}

// ItIs is met when the card in context (ctx.It — a just-played, revealed, or
// discarded card) matches a concrete House and/or Type filter, e.g. "if it is a
// Mars creature" (Brain Stem Antenna reacting to a played card) or "if it is an
// artifact" (Carlo Phantom). Either filter may be left unset to match any. It is
// the concrete-value counterpart to ItIsOfHouse, which names the house by
// reference (the chosen or active house).
type ItIs struct {
	// House and Type are the filters the card in context must match; either one
	// left unset matches any.
	House House
	Type  CardType
	// Not inverts the match, so the condition is met when the card in context does
	// NOT fit the filters — Neutron Shark repeats until it discards a Logos card.
	Not bool
	// Subject names the card outright when "it" has drifted too far from the trigger
	// that set it. Unset says "it".
	Subject Subject
}

// CondText renders the condition, e.g. "if it is a Mars creature", "if it is an
// artifact", or, inverted and named, "if the discarded card is not a Logos card".
func (e ItIs) CondText() string {
	if e.Not {
		return "if " + e.Subject.noun() + " is not " + indefinite(houseTypeNoun(e.House, e.Type))
	}
	return "if " + e.Subject.noun() + " is " + indefinite(houseTypeNoun(e.House, e.Type))
}

// Met reports whether a card is in context and matches the house and type
// filters, inverting the match under Not.
func (e ItIs) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	return e.matches(ctx) != e.Not
}

// matches reports whether the card in context fits the house and type filters.
func (e ItIs) matches(ctx *EffectContext) bool {
	if e.House != HouseNone && ctx.Resolver.House(ctx.It) != e.House {
		return false
	}
	if e.Type != TypeUnset && ctx.Resolver.TypeOf(ctx.It) != e.Type {
		return false
	}
	return true
}

// ItIsOfTrait is met when the creature in context (ctx.It) has the named trait.
type ItIsOfTrait struct{ Trait Trait }

// CondText renders the condition, e.g. "if it is a Dinosaur creature".
func (c ItIsOfTrait) CondText() string {
	return "if it is a " + c.Trait.String() + " creature"
}

// Met reports whether a creature is in context and has the trait.
func (c ItIsOfTrait) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.HasTrait(ctx.It, c.Trait)
}

// ItHasAember is met when the creature in context (ctx.It) has any Æmber on it.
type ItHasAember struct{}

// CondText renders the condition.
func (ItHasAember) CondText() string { return "if it has \u00c6mber on it" }

// Met reports whether a creature is in context with Æmber on it.
func (ItHasAember) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.AmberOn(ctx.It) > 0
}

// Or is met when any one of its Conditions is met, composing conditions instead of
// baking each combination into a bespoke one — Guji Dinosaur Hunter boosts against
// a Dinosaur creature or a creature with Æmber on it.
type Or struct {
	Conditions []Condition
}

// validate requires at least two conditions and rejects any invalid one.
func (o Or) validate() error {
	if len(o.Conditions) < 2 {
		return fmt.Errorf("Or: needs at least two conditions")
	}
	for _, c := range o.Conditions {
		if err := validateCondition(c); err != nil {
			return err
		}
	}
	return nil
}

// CondText joins the sub-clauses with "or", e.g. "if it is a Dinosaur creature or
// it has Æmber on it". Each condition renders "if <clause>" (the shared
// convention), so the leading "if " is dropped before the clauses are joined.
func (o Or) CondText() string {
	clauses := make([]string, len(o.Conditions))
	for i, c := range o.Conditions {
		clauses[i] = strings.TrimPrefix(c.CondText(), "if ")
	}
	return "if " + strings.Join(clauses, " or ")
}

// Met reports whether any of the conditions is met.
func (o Or) Met(ctx *EffectContext) bool {
	for _, c := range o.Conditions {
		if c.Met(ctx) {
			return true
		}
	}
	return false
}

// ItIsOffIdentity is met when the card in context (ctx.It) belongs to none of the
// controller's identity houses — the three houses of their deck. Sneklifter uses
// it to reassign a seized enemy artifact to Shadows only when it is off your
// identity; KeyForge phrases the check negatively ("if it does not belong to one
// of your three houses").
type ItIsOffIdentity struct{}

// CondText renders the condition.
func (ItIsOffIdentity) CondText() string {
	return "if it does not belong to a house on your identity"
}

// Met reports whether a card is in context and belongs to none of the controller's
// identity houses.
func (ItIsOffIdentity) Met(ctx *EffectContext) bool {
	return ctx.HasIt && !ctx.Resolver.PlayerHasHouse(ctx.Controller, ctx.Resolver.House(ctx.It))
}

// ItIsStunned is met when the creature in context (ctx.It) is already stunned.
// 1-2 Punch uses it to destroy a chosen creature that was already stunned rather
// than stunning it.
type ItIsStunned struct{}

// CondText renders the condition.
func (ItIsStunned) CondText() string {
	return "if that creature was already stunned"
}

// Met reports whether a creature is in context and is stunned.
func (ItIsStunned) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.Stunned(ctx.It)
}

// PlayerControlsFewerHousesThan is met while the chosen player controls creatures
// from fewer than Amount distinct houses (Proclamation 346E taxes the opponent's
// keys until they field three houses). It counts houses among creatures only, not
// artifacts.
type PlayerControlsFewerHousesThan struct {
	Player Player
	Amount int
}

// validate requires a positive threshold.
func (c PlayerControlsFewerHousesThan) validate() error {
	if c.Amount <= 0 {
		return fmt.Errorf("PlayerControlsFewerHousesThan: Amount must be positive")
	}
	return nil
}

// CondText renders the condition as a "while" clause naming the player and the
// house count, e.g. "while your opponent does not control creatures from 3 or
// more different houses".
func (c PlayerControlsFewerHousesThan) CondText() string {
	whose := "you do"
	if c.Player == Opponent {
		whose = "your opponent does"
	}
	return fmt.Sprintf(
		"while %s not control creatures from %d or more different houses", whose, c.Amount)
}

// Met counts the distinct houses among the chosen player's creatures and reports
// whether that count is below the threshold.
func (c PlayerControlsFewerHousesThan) Met(ctx *EffectContext) bool {
	player := ctx.PlayerFor(c.Player)
	seen := map[House]bool{}
	for _, id := range ctx.Resolver.Battleline(player) {
		seen[ctx.Resolver.House(id)] = true
	}
	return len(seen) < c.Amount
}

// houseTypeNoun renders a card filtered by house and type as a noun, e.g. "Mars
// creature", "artifact", or the bare "card" when neither is set.
func houseTypeNoun(house House, typ CardType) string {
	n := "card"
	switch typ {
	case Creature:
		n = "creature"
	case Artifact:
		n = "artifact"
	}
	if house != HouseNone {
		n = house.String() + " " + n
	}
	return n
}

// Conditional resolves Then only when Cond is met. It renders as "<cond>, <then>",
// e.g. "if your opponent has 7 Æmber or more, your opponent loses 4 Æmber".
//
// Else, when set, resolves when Cond is not met and renders as the second
// sentence the cards use for a two-way branch: "<cond>, <then>. Otherwise,
// <else>." (Vespilon Theorist archives the revealed card or discards it).
type Conditional struct {
	Cond Condition
	Then Effect
	Else Effect
}

// Text joins the condition and the gated effect.
func (e Conditional) Text() string {
	body := e.Cond.CondText() + ", " + e.Then.Text()
	if e.Else == nil {
		return body
	}
	return body + ". Otherwise, " + e.Else.Text()
}

// Resolve runs Then when Cond is met, otherwise Else if one is set.
func (e Conditional) Resolve(ctx *EffectContext) {
	if e.Cond.Met(ctx) {
		e.Then.Resolve(ctx)
		return
	}
	if e.Else != nil {
		e.Else.Resolve(ctx)
	}
}

// validate checks the gated effect for configuration errors.
func (e Conditional) validate() error {
	if err := validateCondition(e.Cond); err != nil {
		return err
	}
	if err := validateEffect(e.Then); err != nil {
		return err
	}
	if e.Else == nil {
		return nil
	}
	return validateEffect(e.Else)
}

// OrAmount is an alternate scalar an effect switches to when When holds, so a card
// that only changes a number by a fact prints the linear "<base> …, or <alt> if
// <cond>" form instead of a two-armed "If <cond>, … Otherwise, …" branch
// (card-wording rule 22). The effect owns how the numbers read — a bare "2" for a
// steal, a "+2" surcharge for a forge — so OrAmount supplies only the alternate
// value, the guard, and the shared ", or … if <cond>" tail. Its zero value (When
// nil) means "no alternate", so an effect treats an unset OrAmount as absent.
type OrAmount struct {
	Amount int
	When   Condition
}

// set reports whether an alternate is configured.
func (o OrAmount) set() bool { return o.When != nil }

// pick returns the alternate amount when the guard holds, else base.
func (o OrAmount) pick(base int, ctx *EffectContext) int {
	if o.set() && o.When.Met(ctx) {
		return o.Amount
	}
	return base
}

// tail renders ", or <alt> if <cond>", with alt already formatted by the effect.
func (o OrAmount) tail(alt string) string {
	return ", or " + alt + " " + o.When.CondText()
}

// validate rejects an alternate whose guard is missing or invalid.
func (o OrAmount) validate() error {
	if !o.set() {
		return nil
	}
	return validateCondition(o.When)
}

// RuleOfSix is the most times a card can be played, used, or made to resolve
// again in one turn. A self-repeating effect is bounded by it: Bait and Switch's
// "steal 1 Æmber -> repeat this effect" resolves the initial steal plus at most
// five repeats, so it steals six at most however far ahead the opponent is.
const RuleOfSix = 6

// RepeatWhile resolves Do again and again for as long as Cond holds, re-checking
// after each pass — a self-looping effect such as "if your opponent has more
// Æmber than you, steal 1 Æmber -> repeat this effect". The loop also stops the
// moment Do makes no progress (its gate reports it did nothing), so an action that
// is prevented — a steal against a protected pool — ends the loop instead of
// spinning even though Cond still holds. Do is a GatingEffect for exactly that
// reason: the repeat gates on the action completing, not only on Cond.
type RepeatWhile struct {
	Cond Condition
	Do   GatingEffect
}

// Text renders the loop, leading with the condition and closing with the
// self-repeat gate.
func (e RepeatWhile) Text() string {
	return e.Cond.CondText() + ", " + e.Do.Text() + " -> repeat this effect"
}

// Resolve runs Do while Cond is met, stopping as soon as Do does nothing or the
// Rule of Six is reached.
func (e RepeatWhile) Resolve(ctx *EffectContext) {
	for range RuleOfSix {
		if !e.Cond.Met(ctx) || !e.Do.resolveGate(ctx) {
			return
		}
	}
}

// validate checks the looped effect for configuration errors.
func (e RepeatWhile) validate() error {
	return validateEffect(e.Do)
}

// Overwhelmed reports whether the controller is overwhelmed — their opponent
// controls more creatures than they do. "Overwhelmed" is the keyword form of that
// board state; Numquid the Fair repeats its destruction while overwhelmed.
type Overwhelmed struct{}

// CondText renders the condition.
func (Overwhelmed) CondText() string { return "if you are overwhelmed" }

// Met reports whether the opponent controls more creatures than the controller.
func (Overwhelmed) Met(ctx *EffectContext) bool {
	return len(
		ctx.Resolver.Battleline(ctx.Opponent()),
	) > len(
		ctx.Resolver.Battleline(ctx.Controller),
	)
}

// RepeatOnCondition performs an effect and repeats it while the effect keeps
// succeeding and a condition holds — Numquid the Fair's "destroy an enemy creature
// -> if you are overwhelmed, repeat this effect." Do runs at least once; the loop
// stops as soon as Do does nothing (its gate is false) or Cond is not met. When Do
// cannot report progress the Rule of Six alone bounds the loop.
type RepeatOnCondition struct {
	Do   Effect
	Cond Condition
}

// validate checks the repeated effect.
func (e RepeatOnCondition) validate() error {
	return validateEffect(e.Do)
}

// Text renders the effect, e.g. "destroy an enemy creature -> if you are
// overwhelmed, repeat this effect".
func (e RepeatOnCondition) Text() string {
	return e.Do.Text() + " -> " + e.Cond.CondText() + ", repeat this effect"
}

// Resolve runs Do, repeating while it keeps doing something, Cond holds, and the
// Rule of Six allows another pass.
func (e RepeatOnCondition) Resolve(ctx *EffectContext) {
	for range RuleOfSix {
		if !resolveGateOf(ctx, e.Do) || !e.Cond.Met(ctx) {
			return
		}
	}
}

// MayRepeat resolves Do once, then offers the controller the choice to resolve it
// again for as long as Cond holds and they keep accepting — the optional
// counterpart to RepeatWhile, modelling "<do>. If <cond>, you may repeat this
// effect." Do should make progress toward failing Cond so the loop can end.
type MayRepeat struct {
	Cond Condition
	Do   Effect
}

// Text renders the effect, closing with the optional self-repeat gate.
func (e MayRepeat) Text() string {
	return e.Do.Text() + " -> " + e.Cond.CondText() + ", you may repeat this effect"
}

// Resolve runs Do once, then repeats it while Cond holds and the Rule of Six
// allows another pass. When Do leads with a single clickable choice, each repeat
// is driven by that choice — the controller keeps picking to repeat, or passes
// with Done — rather than a separate Yes/No question.
func (e MayRepeat) Resolve(ctx *EffectContext) {
	e.Do.Resolve(ctx)
	d, byChoice := e.Do.(declinableEffect)
	byChoice = byChoice && d.declinable()
	for range RuleOfSix - 1 {
		if !e.Cond.Met(ctx) {
			return
		}
		if byChoice {
			if !d.resolveOptional(ctx) {
				return
			}
			continue
		}
		if ctx.ChooseOption("Repeat this effect?", []string{"Yes", "No"}) != 0 {
			return
		}
		e.Do.Resolve(ctx)
	}
}

// validate checks the repeated effect for configuration errors.
func (e MayRepeat) validate() error {
	return validateEffect(e.Do)
}

// ChoseHouse is met when the controller's active house is House. It is the
// condition behind an "After you choose <House> as your active house, ..."
// ability (Jehu the Bureaucrat): the AfterChooseHouse trigger fires for the
// active player as they pick their house, and this checks whether they picked
// the house the ability watches for.
type ChoseHouse struct {
	House House
}

// CondText renders the condition clause.
func (c ChoseHouse) CondText() string {
	return "you choose " + c.House.String() + " as your active house"
}

// Met reports whether the active house is the one the ability watches for.
func (c ChoseHouse) Met(ctx *EffectContext) bool {
	return ctx.Resolver.ActiveHouse() == c.House
}

// CountIs gates on any Count: Is names the comparison and Amount the threshold
// it compares the count's value against. It is the general "if <something>
// happened N times" condition, so a card that checks a quantity reuses the Count
// vocabulary instead of adding a bespoke condition — Stampede checks the
// creatures used this turn, Vigor checks the damage it just healed.
type CountIs struct {
	Count  Count
	Is     Comparison
	Amount int
}

// validate requires a Count that can render its own "if" clause and a comparison
// the count can answer (MoreThanYou compares two Æmber pools and means nothing
// here).
func (c CountIs) validate() error {
	if _, ok := c.Count.(countClauser); !ok {
		return fmt.Errorf("CountIs: Count must be set and render a CountClause")
	}
	switch c.Is {
	case AtLeast, Exactly:
		return nil
	default:
		return fmt.Errorf("CountIs: Is must be AtLeast or Exactly")
	}
}

// CondText renders the condition, e.g. "if you used 3 or more creatures this
// turn", asking the Count for the clause that reads naturally after "if".
func (c CountIs) CondText() string {
	quantity := fmt.Sprintf("%d or more", c.Amount)
	if c.Is == Exactly {
		quantity = fmt.Sprintf("exactly %d", c.Amount)
	}
	return "if " + c.Count.(countClauser).CountClause(quantity, c.Amount != 1 || c.Is != Exactly)
}

// Met compares the count's current value against the threshold.
func (c CountIs) Met(ctx *EffectContext) bool {
	if c.Is == Exactly {
		return c.Count.Value(ctx) == c.Amount
	}
	return c.Count.Value(ctx) >= c.Amount
}

// countClauser is the optional capability a Count implements to render the "if
// ..." clause CountIs needs. A Count's CountText is a noun ("card destroyed this
// way") that reads well after "for each" but not after "if"; a Count that has a
// natural verb phrase supplies it here, given the rendered quantity ("3 or more",
// "exactly 1") and whether that quantity takes a plural noun.
type countClauser interface {
	CountClause(quantity string, plural bool) string
}

// ForgedKey is the condition on whether a player forged a key in a given window —
// this turn (Smiling Ruth) or on their own previous turn (Tendrils of Pain, Key
// Hammer). It reads the turn history rather than the running key total, so a key
// forged several turns ago does not keep the condition true.
type ForgedKey struct {
	Player   Player
	Previous bool
	// Not inverts the condition, reading "if you have not forged a key this turn"
	// (Nightforge).
	Not bool
}

// validate requires the condition to name whose key it asks about.
func (c ForgedKey) validate() error {
	if !c.Player.valid() {
		return fmt.Errorf("ForgedKey: Player must be set")
	}
	return nil
}

// CondText renders the clause, e.g. "if your opponent forged a key on their
// previous turn".
func (c ForgedKey) CondText() string {
	subject, possessive := "you", "your"
	if c.Player == Opponent {
		subject, possessive = "your opponent", "their"
	}
	when := "this turn"
	if c.Previous {
		when = "on " + possessive + " previous turn"
	}
	if c.Not {
		return fmt.Sprintf("if %s have not forged a key %s", subject, when)
	}
	return fmt.Sprintf("if %s forged a key %s", subject, when)
}

// Met reports whether the named player forged at least one key in the window.
func (c ForgedKey) Met(ctx *EffectContext) bool {
	return (ctx.Resolver.TurnHistory(ctx.PlayerFor(c.Player), c.stat()) > 0) != c.Not
}

// stat picks the tally the window corresponds to.
func (c ForgedKey) stat() TurnStat {
	if c.Previous {
		return KeysForgedLastTurn
	}
	return KeysForgedThisTurn
}

// AemberStolenFromYou is met when the controller had Æmber stolen from them on
// their opponent's previous turn — Information Exchange steals more if it was.
type AemberStolenFromYou struct{}

// CondText renders the condition.
func (AemberStolenFromYou) CondText() string {
	return "if your opponent stole Æmber from you on their previous turn"
}

// Met reports whether any Æmber was stolen from the controller last turn.
func (AemberStolenFromYou) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, AemberStolenFromLastTurn) > 0
}

// EnemyCreatureDestroyed is met while at least one enemy creature has been
// destroyed this turn — Foozle reaps for an extra Æmber once the opponent has
// lost a creature.
type EnemyCreatureDestroyed struct{}

// CondText renders the condition.
func (EnemyCreatureDestroyed) CondText() string {
	return "if an enemy creature has been destroyed this turn"
}

// Met reports whether the controller has seen an enemy creature destroyed this
// turn.
func (EnemyCreatureDestroyed) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, EnemyCreaturesDestroyed) > 0
}

// UsedCreatureToReap is met while the controller has used a creature to reap at
// least once this turn — Bramble Lynx enters play ready once you have reaped.
type UsedCreatureToReap struct{}

// CondText renders the condition.
func (UsedCreatureToReap) CondText() string {
	return "if you have used a creature to reap this turn"
}

// Met reports whether the controller has reaped with a creature this turn.
func (UsedCreatureToReap) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, CreaturesReapedThisTurn) > 0
}

// UsedCreatureToFight is met while the controller has used a creature to fight at
// least once this turn — Alaka enters play ready once you have fought.
type UsedCreatureToFight struct{}

// CondText renders the condition.
func (UsedCreatureToFight) CondText() string {
	return "if you have used a creature to fight this turn"
}

// Met reports whether the controller has fought with a creature this turn.
func (UsedCreatureToFight) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, CreaturesFoughtThisTurn) > 0
}

// HousesRepresented is met when the distinct houses represented among a chosen
// set of in-play cards compare (Is) to Amount — Galactic Census pays out more as
// more houses share the board. Among carries no Max, so the raw house count is
// compared.
type HousesRepresented struct {
	Among  HousesAmong
	Is     Comparison
	Amount int
}

// validate requires a comparison the condition supports.
func (c HousesRepresented) validate() error {
	switch c.Is {
	case AtLeast, AtMost, Exactly:
		return nil
	default:
		return fmt.Errorf("HousesRepresented: Is must be AtLeast, AtMost, or Exactly")
	}
}

// Met compares the surveyed house count against Amount by Is.
func (c HousesRepresented) Met(ctx *EffectContext) bool {
	n := c.Among.Value(ctx)
	switch c.Is {
	case AtMost:
		return n <= c.Amount
	case Exactly:
		return n == c.Amount
	default:
		return n >= c.Amount
	}
}

// CondText renders the condition, e.g. "if there are 3 or more houses represented
// among creatures in play".
func (c HousesRepresented) CondText() string {
	qty := fmt.Sprintf("%d or more", c.Amount)
	switch c.Is {
	case AtMost:
		qty = fmt.Sprintf("%d or fewer", c.Amount)
	case Exactly:
		qty = fmt.Sprintf("exactly %d", c.Amount)
	}
	return fmt.Sprintf("if there are %s houses represented among %s", qty, c.Among.scope())
}

// FirstReapOfTurn is met when the reap in context is the first time a creature
// has reaped this turn — Aember Conduction Unit stuns only the first enemy
// creature to reap. It reads the reaping creature (ctx.It) so it asks about the
// active player's tally, which counts one once this reap has been tallied.
type FirstReapOfTurn struct{}

// CondText renders the condition.
func (FirstReapOfTurn) CondText() string {
	return "if it is the first time a creature has reaped this turn"
}

// Met reports whether exactly one creature has reaped this turn, the reaping
// creature in context being that one.
func (FirstReapOfTurn) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	return ctx.Resolver.TurnHistory(ctx.Resolver.Controller(ctx.It), CreaturesReapedThisTurn) == 1
}

// CounterInPlay is met while at least one card in play carries a generic counter
// of Kind — Wretched Doll destroys every doom-marked creature when there is one,
// and otherwise marks a fresh one.
type CounterInPlay struct {
	// Kind is the counter to look for.
	Kind CounterKind
}

// CondText renders the condition.
func (c CounterInPlay) CondText() string {
	return "if there is a " + c.Kind.noun() + " in play"
}

// Met reports whether any creature in either battleline carries the counter.
func (c CounterInPlay) Met(ctx *EffectContext) bool {
	for player := 0; player < 2; player++ {
		for _, id := range ctx.Resolver.Battleline(player) {
			if ctx.Resolver.CountersOn(id, c.Kind) > 0 {
				return true
			}
		}
	}
	return false
}

// CountersOnThisAtLeast is met when the source card carries at least N counters of
// Kind — The Big One wipes the board once ten or more fuse counters sit on it.
type CountersOnThisAtLeast struct {
	// Kind is the counter to count.
	Kind CounterKind
	// N is the threshold the count must reach.
	N int
}

// CondText renders the condition clause.
func (c CountersOnThisAtLeast) CondText() string {
	return fmt.Sprintf("if there are %d or more %ss on %s", c.N, c.Kind.noun(), SelfName)
}

// Met reports whether the source card holds at least N counters of Kind.
func (c CountersOnThisAtLeast) Met(ctx *EffectContext) bool {
	return ctx.Resolver.CountersOn(ctx.Source, c.Kind) >= c.N
}

// NamedCardPurged is met by whether a card of a given name sits in the
// controller's purge pile — Igon the Terrible destroys itself unless Igon the
// Green has already been purged (Not true reads "has not been purged"). It names
// the other card by its printed name, not the source.
type NamedCardPurged struct {
	// Name is the card name to look for in the purge pile.
	Name string
	// Not flips the sense: false is met while a copy is purged, true while none is.
	Not bool
}

// CondText renders the condition naming the card it looks for.
func (c NamedCardPurged) CondText() string {
	if c.Not {
		return "if " + c.Name + " has not been purged"
	}
	return "if " + c.Name + " has been purged"
}

// Met reports whether a card of the name is in the controller's purge pile,
// flipped by Not.
func (c NamedCardPurged) Met(ctx *EffectContext) bool {
	purged := false
	for _, id := range ctx.Resolver.Purge(ctx.Controller) {
		if ctx.Resolver.Name(id) == c.Name {
			purged = true
			break
		}
	}
	return purged != c.Not
}

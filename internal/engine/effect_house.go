package engine

import "fmt"

// BelongToHouse makes each creature its Target selects belong to House for the
// given Duration, overriding the house it counts as for active-house checks (Brain
// Stem Antenna's host counts as Mars for the rest of the turn). The change is
// per-match state, dropped when the creature leaves play; RemainderOfPlayerTurn
// also drops it at end of turn, while UntilThisLeavesPlay keeps it until the
// creature leaves play.
type BelongToHouse struct {
	Target   Target
	House    House
	Duration Duration
	// Pronoun renders the subject as a back-reference pronoun ("those creatures")
	// for a sentence whose antecedent already named these creatures — Orator
	// Hissaro readies and exalts each neighboring creature, then "those creatures
	// belong to house Saurian". Selection is unchanged; only the text differs.
	Pronoun bool
}

// validate requires a target, a house, and a duration this effect supports.
func (e BelongToHouse) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("BelongToHouse")
	}
	if e.House == HouseNone {
		return fmt.Errorf("BelongToHouse: house must be set")
	}
	if e.Duration != RemainderOfPlayerTurn && e.Duration != UntilThisLeavesPlay {
		return fmt.Errorf(
			"BelongToHouse: duration must be RemainderOfPlayerTurn or UntilThisLeavesPlay",
		)
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, this creature
// belongs to house Mars". An UntilThisLeavesPlay change omits any "until it leaves
// play" clause: every targeted effect ends when the creature leaves play, so the
// clause would only restate the default (Borrow).
func (e BelongToHouse) Text() string {
	if e.Duration == UntilThisLeavesPlay {
		return e.durationSubject() + " " + e.durationPredicate()
	}
	return durationClause(e.Duration, "") + ", " + e.durationSubject() + " " + e.durationPredicate()
}

// durationSubject and durationPredicate split the body so ForDuration can state
// the shared clause and subject once: "it" / "belongs to house Sanctum". A plural
// pronoun subject takes the plural verb ("those creatures belong").
func (e BelongToHouse) durationSubject() string {
	if e.Pronoun {
		return e.Target.pronoun()
	}
	return e.Target.Text()
}
func (e BelongToHouse) durationPredicate() string {
	verb := "belongs"
	if e.Pronoun && e.Target.plural() {
		verb = "belong"
	}
	return verb + " to house " + e.House.String()
}

// Resolve makes each selected creature belong to House for the Duration.
func (e BelongToHouse) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if e.Duration == UntilThisLeavesPlay {
			ctx.Resolver.SetLastingHouse(id, e.House)
		} else {
			ctx.Resolver.BelongToHouseForRemainderOfTurn(id, e.House)
		}
	}
}

// NameHouse remembers the house a surrounding ChooseHouseThen picked on the source
// card, where it stays for as long as that card is in play. It is the writer half
// of a HouseLock whose house is not printed but named: Restringuntus chooses a
// house on play and bars its opponent from it until it leaves play. Player names
// whose choice the lock will constrain, and must match the card's HouseLock.
type NameHouse struct {
	Player Player
}

// validate requires the player whose house choice the named house will constrain.
func (e NameHouse) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("NameHouse")
	}
	return nil
}

// Text renders the lock the named house creates, e.g. "your opponent cannot choose
// that house as their active house until {self} leaves play". The house itself is
// named by the enclosing ChooseHouseThen.
func (e NameHouse) Text() string {
	who, possessive := e.Player.secondPerson()
	return who + " cannot choose that house as " + possessive +
		" active house " + durationClause(UntilThisLeavesPlay, SelfName)
}

// Resolve stores the chosen house on the source card.
func (NameHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.SetNamedHouse(ctx.Source, ctx.ChosenHouse)
}

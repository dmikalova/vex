package engine

// HouseLock is a continuous constraint a card in play puts on one player's
// active-house choice, held for as long as the card stays in play — Pitlord makes
// its controller choose Dis, Restringuntus bars its opponent from a house. It is a
// property of the card rather than an effect because it binds on every turn the
// card is out, not only the turn it was played.
//
// House names the house up front. HouseNone means the card names one as it enters
// play (a NameHouse effect), remembered in CardCore.NamedHouse; that card's Play
// ability carries the wording, so a named-on-play lock prints no rule line of its
// own.
type HouseLock struct {
	// Player is whose choice is constrained, relative to the card's controller.
	// The unset zero value means the card carries no lock.
	Player Player
	// House is the house locked, or HouseNone when the card names one on play.
	House House
	// Bars inverts the lock: the player cannot choose the house rather than must.
	Bars bool
}

// set reports whether the card carries a lock at all.
func (l HouseLock) set() bool { return l.Player.valid() }

// text renders the lock as its printed rule line, e.g. "While {self} is in play you
// must choose Dis as your active house." A lock whose house is named on play
// renders nothing, since its Play ability's text already says what was chosen.
func (l HouseLock) text() string {
	if !l.set() || l.House == HouseNone {
		return ""
	}
	who, possessive := "you", "your"
	if l.Player == Opponent {
		who, possessive = "your opponent", "their"
	}
	verb := "must choose"
	if l.Bars {
		verb = "cannot choose"
	}
	return "While " + SelfName + " is in play " + who + " " + verb + " " +
		l.House.String() + " as " + possessive + " active house."
}

// locked reports the house this lock constrains for the card holding it: the house
// printed on the card, or the one that card named as it entered play.
func (l HouseLock) locked(named House) House {
	if l.House != HouseNone {
		return l.House
	}
	return named
}

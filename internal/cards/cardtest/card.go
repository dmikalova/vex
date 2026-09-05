package cardtest

import "github.com/dmikalova/vactrol/internal/engine"

// Zone names where a card can be, mirroring the game's zones plus Attached (an
// upgrade on a creature), Under (a card placed under a host), and Nowhere (in no
// zone at all — a destroyed card that has been cleared, or an unbound handle).
// At and Location use it so tests read like the game:
// h.Expect(tiger).At(ct.Discard).
type Zone int

// The zones a card can occupy.
const (
	Nowhere  Zone = iota // not found in any zone (e.g. destroyed and cleared)
	PlayArea             // on a battleline or in an artifact row
	Hand
	Discard
	Archives
	Deck
	Attached // an upgrade attached to a creature
	Under    // a card placed under a host (Masterplan, Jargogle)
	Purge    // set aside out of the game
)

// String returns the zone's name for use in failure messages.
func (z Zone) String() string {
	switch z {
	case PlayArea:
		return "play area"
	case Hand:
		return "hand"
	case Discard:
		return "discard"
	case Archives:
		return "archives"
	case Deck:
		return "deck"
	case Attached:
		return "attached"
	case Under:
		return "under"
	case Purge:
		return "purge"
	default:
		return "nowhere"
	}
}

// Card is a handle to one physical card in a running scenario. It wraps the
// engine LocalID with the reads and mutations a test needs, so tests refer to a
// card by handle (or by its definition) rather than by battleline index.
type Card struct {
	h   *Harness
	id  engine.LocalID
	set bool
}

// require fails the test if the handle was never bound to a real card.
func (c Card) require() {
	c.h.t.Helper()
	if !c.set {
		c.h.t.Fatal("cardtest: card handle is not bound — pass it to ct.Bind in the setup")
	}
}

// ID returns the underlying engine id.
func (c Card) ID() engine.LocalID { c.require(); return c.id }

// Name returns the card's printed name.
func (c Card) Name() string { c.require(); return c.h.g.Name(c.id) }

// Damage returns the damage currently on the creature.
func (c Card) Damage() int { c.require(); return c.h.g.Damage(c.id) }

// Power returns the creature's current power (including upgrades).
func (c Card) Power() int { c.require(); return c.h.g.Power(c.id) }

// Armor returns the creature's current armor (including upgrades).
func (c Card) Armor() int { c.require(); return c.h.g.Armor(c.id) }

// AmberOn returns the Æmber sitting on the card.
func (c Card) AmberOn() int { c.require(); return c.h.g.AmberOn(c.id) }

// Exhausted reports whether the card is exhausted.
func (c Card) Exhausted() bool { c.require(); return c.h.g.State.Cards[c.id].Exhausted }

// Stunned reports whether the creature is stunned.
func (c Card) Stunned() bool { c.require(); return c.h.g.State.Cards[c.id].Stunned }

// Location returns the zone the card is currently in.
func (c Card) Location() Zone { c.require(); return c.h.location(c.id) }

// Stun stuns the creature (setup helper for scenarios that begin mid-board).
func (c Card) Stun() { c.require(); c.h.g.State.Cards[c.id].Stunned = true }

// Exhaust exhausts the card.
func (c Card) Exhaust() { c.require(); c.h.g.State.Cards[c.id].Exhausted = true }

// Ready readies the card.
func (c Card) Ready() { c.require(); c.h.g.State.Cards[c.id].Exhausted = false }

// Damaged puts a specific amount of damage on the creature.
func (c Card) Damaged(amount int) { c.require(); c.h.g.State.Cards[c.id].Damage = int16(amount) }

// Attach attaches an upgrade to this creature, refreshing its armor so an
// armor-granting upgrade takes effect immediately.
func (c Card) Attach(up engine.CardDefinition) Card {
	c.require()
	return c.h.attach(c.id, up)
}

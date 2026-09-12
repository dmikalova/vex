package engine

import (
	"fmt"
	"strings"
)

// This file holds the log entries that narrate what happens to a creature or
// artifact in play (ADR 0011): readying, exhausting, stunning, damage, armor,
// destruction, and changes of control. An entry names the card by id, so a
// reader can link it to the card face without matching its name in prose.

// CreatureReadied narrates a card turning upright again.
type CreatureReadied struct{ Creature LocalID }

// Text renders the card that was readied.
func (e CreatureReadied) Text(n Namer) string {
	return fmt.Sprintf("%s is readied", n.Name(e.Creature))
}

// CreatureGainedKeyword narrates a creature gaining a keyword for the turn.
type CreatureGainedKeyword struct {
	Creature LocalID
	Keyword  Keyword
}

// Text renders the creature and the keyword it gained.
func (e CreatureGainedKeyword) Text(n Namer) string {
	return fmt.Sprintf("%s gains %s", n.Name(e.Creature), strings.ToLower(e.Keyword.String()))
}

// CreatureLostKeyword narrates a single creature losing a keyword for the turn.
type CreatureLostKeyword struct {
	Creature LocalID
	Keyword  Keyword
}

// Text renders the creature and the keyword it lost.
func (e CreatureLostKeyword) Text(n Namer) string {
	return fmt.Sprintf("%s loses %s", n.Name(e.Creature), strings.ToLower(e.Keyword.String()))
}

// CreatureGainedStats narrates a creature gaining power and/or armor for the turn
// (Abond the Armorsmith grants +1 armor).
type CreatureGainedStats struct {
	Creature LocalID
	Power    int
	Armor    int
}

// Text renders the creature and the stats it gained, e.g. "Card2 gains +1 armor".
func (e CreatureGainedStats) Text(n Namer) string {
	return fmt.Sprintf("%s gains %s", n.Name(e.Creature),
		staticBonuses(StaticModifier{PowerBonus: e.Power, ArmorBonus: e.Armor}))
}

// CreatureGainedAssault narrates a creature gaining Assault for the turn (Creed of
// Nature grants assault equal to its power).
type CreatureGainedAssault struct {
	Creature LocalID
	Amount   int
}

// Text renders the creature and the Assault it gained, e.g. "Card2 gains assault 3".
func (e CreatureGainedAssault) Text(n Namer) string {
	return fmt.Sprintf("%s gains assault %d", n.Name(e.Creature), e.Amount)
}

// CreatureConsideredFlank narrates a creature being treated as a flank creature
// for the turn (Spectral Tunneler).
type CreatureConsideredFlank struct{ Creature LocalID }

// Text renders the creature now considered a flank creature.
func (e CreatureConsideredFlank) Text(n Namer) string {
	return fmt.Sprintf("%s is considered a flank creature", n.Name(e.Creature))
}

// CreatureExhausted narrates a card being turned sideways.
type CreatureExhausted struct{ Creature LocalID }

// Text renders the card that was exhausted.
func (e CreatureExhausted) Text(n Namer) string {
	return fmt.Sprintf("%s is exhausted", n.Name(e.Creature))
}

// CreatureStunned narrates a card being stunned, and by what — unless the
// stunner is the card itself (Chuff Ape enters play stunned), which reads better
// left passive. AlreadyStunned marks a stun that found its target already
// stunned: the source still had to choose it, so the choice is worth a line even
// though nothing changed (the same reasoning as a steal that finds an empty
// pool: see ReapedStealing's "no Æmber to steal").
type CreatureStunned struct {
	Creature       LocalID
	By             LocalID
	AlreadyStunned bool
}

// Text renders the card that was stunned and, when it is not self-inflicted, the
// card that stunned it.
func (e CreatureStunned) Text(n Namer) string {
	if e.AlreadyStunned {
		return fmt.Sprintf("%s is already stunned", n.Name(e.Creature))
	}
	if e.By != e.Creature {
		return fmt.Sprintf("%s stunned %s", n.Name(e.By), n.Name(e.Creature))
	}
	return fmt.Sprintf("%s is stunned", n.Name(e.Creature))
}

// CreatureEnraged narrates a card being enraged, and by what — unless the source
// is the card itself, which reads better left passive. AlreadyEnraged marks an
// enrage that found its target already enraged: the source still had to choose it,
// so the choice is worth a line even though nothing changed.
type CreatureEnraged struct {
	Creature       LocalID
	By             LocalID
	AlreadyEnraged bool
}

// Text renders the card that was enraged and, when it is not self-inflicted, the
// card that enraged it.
func (e CreatureEnraged) Text(n Namer) string {
	if e.AlreadyEnraged {
		return fmt.Sprintf("%s is already enraged", n.Name(e.Creature))
	}
	if e.By != e.Creature {
		return fmt.Sprintf("%s enraged %s", n.Name(e.By), n.Name(e.Creature))
	}
	return fmt.Sprintf("%s is enraged", n.Name(e.Creature))
}

// CreatureWarded narrates a card being warded, and by what — unless the source is
// the card itself, which reads better left passive. AlreadyWarded marks a ward
// that found its target already warded: the source still had to choose it, so the
// choice is worth a line even though nothing changed.
type CreatureWarded struct {
	Creature      LocalID
	By            LocalID
	AlreadyWarded bool
}

// Text renders the card that was warded and, when it is not self-inflicted, the
// card that warded it.
func (e CreatureWarded) Text(n Namer) string {
	if e.AlreadyWarded {
		return fmt.Sprintf("%s is already warded", n.Name(e.Creature))
	}
	if e.By != e.Creature {
		return fmt.Sprintf("%s warded %s", n.Name(e.By), n.Name(e.Creature))
	}
	return fmt.Sprintf("%s is warded", n.Name(e.Creature))
}

// WardRemoved narrates a ward being taken off a creature. AlreadyUnwarded marks a
// removal that found its target carrying no ward: the source still chose it, so
// the choice is worth a line even though nothing changed.
type WardRemoved struct {
	Creature        LocalID
	By              LocalID
	AlreadyUnwarded bool
}

// Text renders the creature whose ward was removed and the card that removed it.
func (e WardRemoved) Text(n Namer) string {
	if e.AlreadyUnwarded {
		return fmt.Sprintf("%s has no ward to remove", n.Name(e.Creature))
	}
	return fmt.Sprintf("%s removes the ward from %s", n.Name(e.By), n.Name(e.Creature))
}

// wardPrevented names what a spent ward stopped, so the log can say what it saved
// the creature from: an instance of damage, a destruction, or another way of
// leaving play (returned to hand, archived, shuffled or put back on the deck,
// purged, or grafted under a host).
type wardPrevented uint8

const (
	wardLeavePlay wardPrevented = iota
	wardDestruction
	wardDamage
)

// WardAbsorbed narrates a creature's ward being spent: it absorbed an instance of
// damage or a removal from play, so the creature stays and loses its ward.
// Prevented names what it stopped; Amount is the damage it refused when the ward
// spent on damage.
type WardAbsorbed struct {
	Creature  LocalID
	Prevented wardPrevented
	Amount    int
}

// Text renders the creature whose ward was spent and what the ward stopped.
func (e WardAbsorbed) Text(n Namer) string {
	switch e.Prevented {
	case wardDamage:
		return fmt.Sprintf("%s's ward prevents the %d damage", n.Name(e.Creature), e.Amount)
	case wardDestruction:
		return fmt.Sprintf("%s's ward prevents the destruction", n.Name(e.Creature))
	default:
		return fmt.Sprintf("%s's ward keeps it in play", n.Name(e.Creature))
	}
}

// NoCreatureToFight narrates a fight that found no enemy creature to attack, so
// nothing happened.
type NoCreatureToFight struct{ Creature LocalID }

// Text renders a fight that found no enemy creature to attack.
func (e NoCreatureToFight) Text(n Namer) string {
	return fmt.Sprintf("%s has no creature to fight", n.Name(e.Creature))
}

// CardsRevealedToAll narrates cards shown to both players. It is the one entry
// allowed to name a card that is otherwise hidden, because revealing it is
// exactly what made it public (ADR 0011).
type CardsRevealedToAll struct {
	Player int
	Cards  []LocalID
}

// Text renders the cards a player revealed, each by name.
func (e CardsRevealedToAll) Text(n Namer) string {
	return fmt.Sprintf("%s reveals %s", n.PlayerName(e.Player), namedCards(n, e.Cards))
}

// PositionsSwapped narrates two creatures trading places in a battleline.
type PositionsSwapped struct{ A, B LocalID }

// Text renders the two creatures that traded places.
func (e PositionsSwapped) Text(n Namer) string {
	return fmt.Sprintf("%s swaps positions with %s", n.Name(e.A), n.Name(e.B))
}

// CardsSwapped narrates a card entering play in another card's place while that
// card leaves to the zone the entering card came from — a swap across zones
// (Gebuk). A names the card that stayed on the board; B the card that entered it.
type CardsSwapped struct {
	A, B       LocalID
	FromPlayer int
	FromZone   Zone
}

// Text renders the swap, naming the zone the entering card came from.
func (e CardsSwapped) Text(n Namer) string {
	return fmt.Sprintf("%s swaps places with %s from %s's %s",
		n.Name(e.A), n.Name(e.B), n.PlayerName(e.FromPlayer), e.FromZone.noun())
}

// MovedToFlank narrates a creature moving to a flank of its battleline.
type MovedToFlank struct {
	Creature LocalID
	Right    bool
}

// Text renders the creature and the flank it moved to.
func (e MovedToFlank) Text(n Namer) string {
	side := "left"
	if e.Right {
		side = "right"
	}
	return fmt.Sprintf("%s moves to the %s flank", n.Name(e.Creature), side)
}

// MovedWithinBattleline narrates a creature repositioned within its battleline.
type MovedWithinBattleline struct {
	Creature LocalID
}

// Text renders the creature that moved.
func (e MovedWithinBattleline) Text(n Namer) string {
	return fmt.Sprintf("%s moves within its battleline", n.Name(e.Creature))
}

// TurnedIntoCreature narrates an artifact turning itself into a creature and
// moving onto a flank of its battleline (Auto-Legionary).
type TurnedIntoCreature struct {
	Card  LocalID
	Right bool
}

// Text renders the card and the flank it entered as a creature.
func (e TurnedIntoCreature) Text(n Namer) string {
	side := "left"
	if e.Right {
		side = "right"
	}
	return fmt.Sprintf(
		"%s becomes a creature on the %s flank", n.Name(e.Card), side)
}

// ControlTaken narrates a card moving into another player's rows without
// changing owner.
type ControlTaken struct {
	Player int
	Card   LocalID
}

// Text renders the card a player took control of.
func (e ControlTaken) Text(n Namer) string {
	return fmt.Sprintf("%s takes control of %s", n.PlayerName(e.Player), n.Name(e.Card))
}

// ControlReturned narrates borrowed control lapsing when its source left play.
type ControlReturned struct {
	Card  LocalID
	Owner int
}

// Text renders borrowed control lapsing back to the card's owner.
func (e ControlReturned) Text(n Namer) string {
	return fmt.Sprintf("%s returns to %s's control", n.Name(e.Card), n.PlayerName(e.Owner))
}

// CardDestroyed narrates a card in play being destroyed.
type CardDestroyed struct{ Card LocalID }

// Text renders the card that was destroyed.
func (e CardDestroyed) Text(n Namer) string {
	return fmt.Sprintf("%s is destroyed", n.Name(e.Card))
}

// CardsDestroyedBy narrates one card's effect destroying a group of cards, so the
// destruction is credited to its agent — "Strange Gizmo destroys A, B, and C" —
// rather than a passive line per creature.
type CardsDestroyedBy struct {
	Source LocalID
	Cards  []LocalID
}

// Text renders the source card and the cards it destroyed.
func (e CardsDestroyedBy) Text(n Namer) string {
	return fmt.Sprintf("%s destroys %s", n.Name(e.Source), namedCardsAnd(n, e.Cards))
}

// DestructionReplaced narrates an upgrade or ability taking a card's destruction
// on itself.
type DestructionReplaced struct {
	Card LocalID
	By   LocalID
}

// Text renders the card that took another's destruction on itself.
func (e DestructionReplaced) Text(n Namer) string {
	return fmt.Sprintf("%s would be destroyed, so %s replaces its destruction",
		n.Name(e.Card), n.Name(e.By))
}

// AemberOnCardReleased narrates the Æmber a card was holding going to a pool
// when the card left play.
type AemberOnCardReleased struct {
	Card   LocalID
	Amount int
	To     int
}

// Text renders the Æmber a card released to a pool as it left play.
func (e AemberOnCardReleased) Text(n Namer) string {
	return fmt.Sprintf("%d Æmber on %s goes to %s's pool",
		e.Amount, n.Name(e.Card), n.PlayerName(e.To))
}

// StunRecovered narrates a use spent removing a stun counter instead of reaping,
// fighting, or acting. It reads as the controller unstunning the creature, so the
// log speaks the same "unstun" verb as the Unstun effect (Clear Mind).
type StunRecovered struct {
	Player   int
	Creature LocalID
}

// Text renders a use spent recovering from stun.
func (e StunRecovered) Text(n Namer) string {
	return fmt.Sprintf("%s unstuns %s", n.PlayerName(e.Player), n.Name(e.Creature))
}

// CardCannotBeUsed narrates a use refused because the card was exhausted.
type CardCannotBeUsed struct{ Card LocalID }

// Text renders a use refused because the card was exhausted.
func (e CardCannotBeUsed) Text(n Namer) string {
	return fmt.Sprintf("%s is exhausted and cannot be used", n.Name(e.Card))
}

package engine

// An ItNoun is a wording choice: the noun phrase a condition prints in place of
// the pronoun "it". It never changes which card the condition reads — that is the
// Subject's job — only what the printed text calls it. The pronoun reads fine
// while the condition sits right under the trigger that put the card in context
// ("After you play a card, if it is a creature…"), but drifts once an effect puts
// other cards in focus between the two — Neutron Shark destroys two cards before
// it discards, so it asks about "the discarded card", not "it".
type ItNoun uint8

const (
	// The zero value is the bare "it", which needs no name of its own.
	_ ItNoun = iota
	// DiscardedCard names the card an effect just discarded.
	DiscardedCard
	// ThatCard names the card an effect just acted on when "it" would be ambiguous
	// — Fidgit discards from one of two sources, so its follow-up says "that card".
	ThatCard
	// FoughtCreature names the creature the source is fighting — Baldric the Bold
	// asks about "the fought creature", not "it".
	FoughtCreature
	// ThatCreature names the creature an effect just acted on, when the effect
	// between the two put something else in focus that "it" would otherwise attach
	// to — Shadowsaurus moves Æmber off a creature, so "it" would read as the Æmber.
	ThatCreature
)

// noun renders the wording choice as the phrase a condition puts in front of "is".
func (n ItNoun) noun() string {
	switch n {
	case DiscardedCard:
		return "the discarded card"
	case ThatCard:
		return "that card"
	case FoughtCreature:
		return "the fought creature"
	case ThatCreature:
		return "that creature"
	default:
		return "it"
	}
}

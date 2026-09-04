package engine

// A Subject names the card a condition means by "it". The pronoun reads fine
// while the condition sits right under the trigger that put the card in context
// ("After you play a card, if it is a creature…"), but drifts once an effect puts
// other cards in focus between the two — Neutron Shark destroys two cards before
// it discards, so it asks about "the discarded card", not "it".
type Subject uint8

const (
	// The zero value is the bare "it", which needs no name of its own.
	_ Subject = iota
	// DiscardedCard names the card an effect just discarded.
	DiscardedCard
)

// noun renders the subject as the phrase a condition puts in front of "is".
func (s Subject) noun() string {
	if s == DiscardedCard {
		return "the discarded card"
	}
	return "it"
}

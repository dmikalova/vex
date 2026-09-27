package engine

// A Subject names which card a condition reads. It is a real referent, not a
// wording choice: a condition's Met resolves it to a card and asks its question
// of that card. This is what lets one condition answer the same question about
// either card — HasAember{} asks about the card in context, and
// HasAember{Subject: This} asks about the card the ability is printed on.
//
// The two referents are the only two the engine has. There is no separate "this"
// card: a card's own identity is its source, so This and the source are one
// referent, under the name a card's printed text uses for itself.
//
// Every other condition whose name starts with It or Source is deliberately left
// without a Subject, and there are more than a dozen of them. The prefix is not a
// missing field: an audit of each found no counterpart question on the other
// referent — there is no ItIsReady beside SourceReady, no ThisIsStunned beside
// ItIsStunned, no SourceIsNamed beside ItIsNamed. HasAember is the only question
// the card pool asks of both referents, which is why it alone carries a Subject.
// Adding the field to the rest now would be speculative generality: the unused
// branch is unreachable by any card, so it could only be covered by a test
// written to hold the coverage gate at 100%. Their CondText also varies more than
// the referent does — SourceIsFighting says "if fighting" so the constant-ability
// renderer can reframe it, and ItIsStunned says "that creature" rather than "it"
// — so a merged wording would not fall out of name() anyway.
//
// Add Subject to a condition when a second card actually asks that question of
// the other referent. The enum and its card/name helpers are already here, so it
// is a one-node change then.
type Subject uint8

const (
	// It is the zero value and the card in context (ctx.It) — the card a trigger or
	// a preceding effect just put in focus. A condition on It is not met when no
	// card is in context.
	It Subject = iota
	// This is the card the ability is printed on (ctx.Source), the referent a card
	// uses to ask a question about itself (Odoac the Patrician protects its pool
	// only while it holds Æmber).
	This
)

// card resolves the subject to the card it names, reporting false when the
// subject names the card in context and no card is in context.
func (s Subject) card(ctx *EffectContext) (LocalID, bool) {
	if s == This {
		return ctx.Source, true
	}
	return ctx.It, ctx.HasIt
}

// name renders the subject as the phrase printed text uses for it.
func (s Subject) name() string {
	if s == This {
		return SelfName
	}
	return "it"
}

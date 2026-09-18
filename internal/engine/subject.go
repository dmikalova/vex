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

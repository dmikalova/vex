package engine

// Card Types rulebook terms (ADR 0018): each describes itself next to the code it
// governs; the completeness test fails the build if a member of the matching
// closed catalog has no term here.
func init() {
	registerRuleSectionIntro(
		SectionCardType,
		`Every card is one of the following four types. A card's type determines where it
goes when played and how it is used.

Every card also shows the same anatomy: a **house** icon in the upper-left corner
(the faction it belongs to), the card **name**, its **type**, any **traits**
(flavor labels such as _Knight_ or _Robot_ that other cards reference but that
carry no rules of their own), and its rules text. Creatures additionally show a
**power** value and, sometimes, an **armor** value.`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section: SectionCardType,
			Title:   "Creature",
			Body: `A creature is a unit you play into your battleline. Once it is ready, it can
reap for Æmber, fight an enemy creature, or use an "Action:" ability.`,
		},
		{
			Section: SectionCardType,
			Title:   "Tactic",
			Body: `A tactic (KeyForge's "action" card type, renamed to free the word "Action"
for the ability) is a one-shot card: its effect resolves as you play it, and
it then goes straight to your discard pile.`,
		},
		{
			Section: SectionCardType,
			Title:   "Artifact",
			Body: `An artifact is a permanent card you play alongside your creatures. It stays
in play until something removes it and is typically used for its "Action:"
ability.`,
		},
		{
			Section: SectionCardType,
			Title:   "Upgrade",
			Body: `An upgrade attaches to a creature as you play it, changing that creature's
stats or granting it keywords and abilities for as long as it stays attached.`,
		},
	})
}

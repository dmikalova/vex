package engine

// Bonus-icon rulebook terms (ADR 0018): the completeness test fails the build if a
// BonusIcon has no term here.
func init() {
	registerRuleSectionIntro(
		SectionBonus,
		`Bonus icons are printed in a card's upper-left corner. After a card is played,
its icons resolve one at a time, top to bottom, before the card's "Play:" ability.
Two Vactrol rules differ from KeyForge: the card carrying an icon is the source of
that icon's effect, and if a creature leaves play before or while its icons
resolve, its remaining icons do not resolve.`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section:    SectionBonus,
			Title:      "Æmber",
			Definition: "Gain 1 Æmber.",
			Body:       `An Æmber bonus icon gains its player 1 Æmber from the common supply.`,
		},
		{
			Section:    SectionBonus,
			Title:      "Capture",
			Definition: "A friendly creature captures 1 Æmber from the opponent.",
			Body: `A Capture bonus icon has a friendly creature its controller chooses capture 1
Æmber from the opponent. It does nothing when the opponent has no Æmber or the
controller has no creature to hold it.`,
		},
		{
			Section:    SectionBonus,
			Title:      "Damage",
			Definition: "Deal 1 damage to a creature in play.",
			Body: `A Damage bonus icon deals 1 damage to a creature its controller chooses, friend
or foe. With no enemy creature in play the damage must land on a friendly one.`,
		},
		{
			Section:    SectionBonus,
			Title:      "Draw",
			Definition: "Draw 1 card.",
			Body:       `A Draw bonus icon draws its player 1 card.`,
		},
		{
			Section:    SectionBonus,
			Title:      "Enhance",
			Definition: "During deck generation, adds the listed bonus icons to random cards in the deck.",
			Body: `Enhance is a deck-generation rule, not an in-game one: when the deck is built,
the listed bonus icons are added to random cards in it (up to five icons per card,
counting printed icons), and the Enhance card itself gains nothing during play. A
card may opt out of receiving Enhance icons, a Vactrol addition for cards a bonus
would only weaken.`,
		},
	})
}

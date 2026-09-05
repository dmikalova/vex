package engine

// The Turn rulebook terms (ADR 0018): each describes itself next to the code it
// governs; the completeness test fails the build if a member of the matching
// closed catalog has no term here.
func init() {
	registerRuleSectionIntro(
		SectionTurn,
		`Each turn runs in a fixed order. You forge any keys you can afford, choose the one
house you will act with, then play and use that house's cards as you like —
playing from hand, reaping for Æmber, and fighting with your creatures — before
readying your cards and drawing back up to a full hand.`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: "2. Choose a house",
			Body: `Choose a house: pick one of your deck's three houses to be your active house
for this turn. For the rest of the turn you may play from hand and use only
cards of that house, except cards that ignore the restriction such as Versatile
ones.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: "3. Ready and draw",
			Body: `Ready and draw: at the end of your turn every card you control readies (turns
back upright, ready to act next turn) and you draw back up to a full hand of
six cards. Creatures and artifacts that entered play exhausted this turn ready
here too.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: "1. Forge a key",
			Body: `Forge a key: at the start of your turn you forge a single key if you can pay
its current cost — 6 Æmber by default. A player forges at most one key per turn.
Keys are the win condition — forge your third key and you win the game.`,
		},
	})
}

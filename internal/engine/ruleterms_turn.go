package engine

// The Turn rulebook terms (ADR 0018): each describes itself next to the code it
// governs; the completeness test fails the build if a member of the matching
// closed catalog has no term here.
func init() {
	registerRuleSectionIntro(
		SectionTurn,
		`Each turn runs through eight phases in a fixed order. Most run on their own; two
of them — choosing a house and the main phase — wait for the active player. Each
phase carries its own rules, and abilities that trigger at the start or end of the
turn resolve in the matching phase.`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhaseStartOfTurn.rulebookStep(),
			Definition: "The turn runs through eight phases in a fixed order: start of turn, " +
				"forge a key, choose a house, archives, the main phase, ready, and draw, " +
				"then end of turn.",
			Body: `The turn opens here. Abilities that trigger "at the start of your turn" resolve
now, before you forge, so an ability that changes what a key costs or adjusts your
Æmber acts before the forge phase reads it. You order your own start-of-turn
abilities when more than one triggers.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhaseForge.rulebookStep(),
			Body: `You forge a single key if you can pay its cost — 6 Æmber by default — spending
that Æmber. You forge at most one key per turn, and an effect can raise the cost
or make you skip the phase. Keys are the win condition: forge your third key and
you win the game.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhaseChooseHouse.rulebookStep(),
			Body: `You pick one of your deck's three houses as your active house for the turn. For
the rest of the turn you may play from hand and use only cards of that house,
except cards that ignore the restriction such as Versatile ones.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhaseArchives.rulebookStep(),
			Body: `With a house chosen, you are offered your archived cards. You may take all of
them into your hand at once — archived cards are set aside face-down on earlier
turns, out of your opponent's reach. You are not prompted when your archives are
empty.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhasePlay.rulebookStep(),
			Body: `The open phase, where you take your turn. In any order you may play cards from
your hand, discard from your hand, and use your ready cards of the active house: a
creature reaps for Æmber, fights an enemy creature, or takes an "Action:" ability,
and an artifact takes its "Action:". Combat happens here — fighting is one way of
using a creature. You stay in the main phase until you end your turn.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhaseReady.rulebookStep(),
			Body: `Ending your turn readies every card you control — turning your exhausted cards
upright — and refreshes each creature's armor to full for the turns to come. Cards
that entered play exhausted this turn ready here too, and the turn's own temporary
effects expire as it ends.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhaseDraw.rulebookStep(),
			Body: `You draw back up to a full hand — six cards by default, adjusted by effects.
Chains cut your draw: you draw one fewer card for every 6 chains you hold, and
shed a single chain only on a turn the reduction actually kept you from a card.`,
		},
		{
			Section:  SectionTurn,
			Title:    "Turn structure",
			Subtitle: PhaseEndOfTurn.rulebookStep(),
			Body: `The turn closes here. Abilities that trigger "at the end of your turn" resolve
now, last of all, so they see the board and hand the turn actually ends with. You
order your own end-of-turn abilities when more than one triggers. Play then passes
to your opponent.`,
		},
	})
}

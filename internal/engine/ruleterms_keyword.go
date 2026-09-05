package engine

// Keywords rulebook terms (ADR 0018): each describes itself next to the code it
// governs; the completeness test fails the build if a member of the matching
// closed catalog has no term here.
func init() {
	registerRuleSectionIntro(
		SectionKeyword,
		`Keywords are standing rules a creature always has, printed on a line of their own
at the top of the card's text box. A keyword needs no trigger — it simply changes
how the creature behaves. Some keywords take a number (for example _Assault 2_).`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section:    SectionKeyword,
			Title:      "Assault",
			Definition: "The creature deals N damage to the creature it attacks, just before combat damage is dealt.",
			Body: `A creature with Assault N deals N damage to the creature it attacks,
immediately before combat damage is dealt. Zero means the creature does not
have Assault.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Hazardous",
			Definition: "The creature deals N damage to any creature that attacks it, before that attacker deals combat damage.",
			Body: `A creature with Hazardous N deals N damage to any creature that attacks it,
before that attacker deals its combat damage. Zero means the creature does
not have Hazardous.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Skirmish",
			Definition: "The creature takes no damage back when it is used to fight.",
			Body: `A creature with Skirmish takes no damage when it is used to fight: it deals
its power to the enemy creature but takes none back.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Poison",
			Definition: "Any damage dealt to the creature destroys it, whatever power it has left.",
			Body: `Any amount of damage dealt to a creature with Poison destroys it, however
much power it has left.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Elusive",
			Definition: "The first time the creature is chosen to be fought each turn, no fight damage is dealt by or to it.",
			Body: `Elusive: the first time this creature is chosen to be fought each turn, no
pending fight damage is dealt by or to it. Later fights that same turn deal
damage normally.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Taunt",
			Definition: "The creature's neighbors cannot be chosen to be fought unless they have Taunt themselves.",
			Body: `Taunt: this creature's neighbors cannot be chosen to be fought unless they
have taunt themselves, so a Taunt creature shields the creatures beside it.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Versatile",
			Definition: "The card may be used as if it belonged to your active house, though it is still played only on its own house's turn.",
			Body: `A card with Versatile may, once in play, be used (reap/fight/action) as if
it belonged to the active house. It does not relax playing from hand — a
Versatile card is still played only when its own house is the one chosen
this turn.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Alpha",
			Definition: "The card can be played only as the first thing its player does on their turn.",
			Body: `A card with Alpha can only be played as the first thing its player does on
their turn: it cannot be played once that player has played, used, or
discarded any other card this turn (First Blood, Unlocked Gateway's
counterpart). It is a restriction on playing, so it never appears as a
granted or lost creature keyword.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Omega",
			Definition: "The card ends the current step of the turn the moment it resolves; no more cards may be played, used, or discarded that step.",
			Body: `A card with Omega ends the current step of the turn the moment it resolves:
no more cards may be played, used, or discarded for the rest of that step,
except through pending abilities and effects still resolving (Unlocked
Gateway). Play then continues to the next step, so more cards can still be
played later that turn. Like Alpha it constrains playing, not combat.`,
		},
		{
			Section:    SectionKeyword,
			Title:      "Deploy",
			Definition: "The creature may enter play at any position in its controller's battleline, not only on a flank.",
			Body: `A creature with Deploy may enter play at any position in its controller's
battleline, not only on a flank — its controller chooses the spot as it is
played ("Lion" Bautrem, Challe the Safeguard). It matters only while the
creature is being played, so it too is never granted or lost.`,
		},
	})
}

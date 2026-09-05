package engine

// Combat rulebook terms (ADR 0018): each describes itself next to the code it
// governs; the completeness test fails the build if a member of the matching
// closed catalog has no term here.
func init() {
	registerRuleTerms([]RuleTerm{
		{
			Section:    SectionCombat,
			Title:      "Combat",
			Definition: "Using one of your ready creatures during the main phase to fight an enemy creature; both deal damage equal to their power at the same time.",
			Body: `Fighting is one of the ways you use a creature during the main phase, alongside
reaping and taking an "Action:" ability; combat is not a phase of its own. Use one
of your ready creatures to fight an enemy creature. Using it to fight exhausts it.
First, any "Before Fight" abilities and the Assault and
Hazardous keywords resolve; if these destroy either creature, the fight does not
occur. Otherwise both creatures deal damage equal to their power at the same
time — armor reduces the damage a creature takes, and a Skirmish attacker takes
no damage back. A creature with damage equal to or greater than its power is
destroyed, and both deaths are resolved together, so neither fighter's
destruction changes the damage the other deals. Finally, if the attacker
survived, its "Fight" abilities resolve.`,
		},
		{
			Section:    SectionCombat,
			Title:      "Armor",
			Definition: "A value that absorbs damage — each point stops 1 — refreshing when its creature's controller readies.",
			Body: `Armor absorbs damage. A creature with armor prevents that much of the damage it
would be dealt: each point of armor stops 1 damage, and armor spent this way does
not come back until the creature's controller readies at the end of their turn.
Armor never reduces a creature's power, and healing does not restore spent armor.`,
		},
	})
}

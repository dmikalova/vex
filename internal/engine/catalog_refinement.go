package engine

// This file is the census of the Refinement family: one row per Refinement the
// package declares, each naming the rulebook term it owes or the reason it owes
// none. See catalog.go for the classification rule, and target_refinement.go for
// the refinements themselves. TestFamilyTotality keeps it complete — a
// refinement with no row, or a row naming no refinement, fails the build.

// refinementPhrase is the target phrase the census renders each refinement's
// clause against, standing in for the Target a card would pair it with.
const refinementPhrase = "each creature"

// refinementFamily is the Refinement family's census entry. A refinement is
// discovered by its refine method, which no other family declares.
func refinementFamily() Family {
	return newFamily(
		"Refinement",
		"refine",
		[]string{"*EffectContext", "[]LocalID"},
		RefinementCatalog(),
		func(r Refinement) string { return r.clause(refinementPhrase) },
	).gated()
}

// RefinementCatalog returns one representative value of every Refinement, with
// the rulebook term each owes. Refinements group by what they measure: the
// combinators that reshape another refinement's result, the power tiers and
// thresholds, and the per-battleline selections. A rule decidable one candidate
// at a time is a Filter axis, not a refinement, so it has no row here.
func RefinementCatalog() []Catalogued[Refinement] {
	return []Catalogued[Refinement]{
		// Combinators over other refinements.
		{
			Node:  Except(MostPowerful),
			Rules: plumbing("composition: keeps what its inner refinement drops"),
		},
		{
			Node:  AnyOf(LowestPower, HighestPower),
			Rules: plumbing("composition: keeps what any of its inner refinements keeps"),
		},
		{
			Node: Filter{Trait: Dinosaur},
			Rules: plumbing(
				"composition: a per-card Filter, whose axes carry their own terms as " +
					"Target filters; it is a refinement only inside a union combinator"),
		},

		// Power tiers: which end of the power scale is selected, and how many.
		{
			Node:  MostPowerful,
			Rules: bears("Power Tier"),
		},
		{
			Node:  LeastPowerful,
			Rules: bears("Power Tier"),
		},
		{
			Node:  HighestPower,
			Rules: bears("Power Tier"),
		},
		{
			Node:  LowestPower,
			Rules: bears("Power Tier"),
		},

		// Power thresholds: a power compared against a number, a count, or a card.
		{
			Node:  PowerLessThan(ForgedKeys{Player: Controller}),
			Rules: bears("Power Threshold"),
		},
		{
			Node:  SamePowerAsChosen,
			Rules: bears("Same Power as a Chosen Creature"),
		},

		// Per-battleline selections, which act on each side in turn.
		{
			Node:  KeepPerSide(3),
			Rules: bears("Per Battleline"),
		},
		{
			Node:  PortionPerSide(ThirdRoundedUp),
			Rules: bears("Per Battleline"),
		},
	}
}

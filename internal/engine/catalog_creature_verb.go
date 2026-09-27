package engine

// This file is the census of the CreatureVerb family: one row per verb an
// OnChooseCreature can apply to the creature it picked. See catalog.go for the
// classification rule, and effect_creature.go for the verbs themselves.
//
// Every verb names an action the rulebook already teaches in its turn, combat or
// keyword sections — using, reaping, fighting, stunning, exhausting — so the rows
// bind across sections to those terms rather than restating them here.

// creatureVerbFamily is the CreatureVerb family's census entry. A verb is
// discovered by its VerbText method.
func creatureVerbFamily() Family {
	return newFamily(
		"CreatureVerb",
		"VerbText",
		nil,
		CreatureVerbCatalog(),
		CreatureVerb.VerbText,
	).gated()
}

// CreatureVerbCatalog returns one representative value of every CreatureVerb,
// with the rulebook term each owes.
func CreatureVerbCatalog() []Catalogued[CreatureVerb] {
	return []Catalogued[CreatureVerb]{
		{
			Node:  UseVerb{},
			Rules: bears("Use a Card"),
		},
		{
			Node:  ReapVerb{},
			Rules: bears("Reap"),
		},
		{
			Node:  FightVerb{},
			Rules: bears("Fight"),
		},
		{
			Node:  ReadyVerb{},
			Rules: bears("Ready"),
		},
		{
			Node:  ExhaustVerb{},
			Rules: bears("Exhaust"),
		},
		{
			Node:  StunVerb{},
			Rules: bears("Stun"),
		},
		{
			Node:  GainKeywordVerb{Keyword: Skirmish},
			Rules: bears("Gain and Lose Keywords"),
		},
	}
}

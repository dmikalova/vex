package engine

// This file is the census of the Condition family: one row per condition the
// package declares, each naming the rulebook term it owes or the reason it owes
// none. See catalog.go for the classification rule, and effect_condition.go for
// the interface the conditions implement. TestFamilyTotality keeps it complete —
// a condition with no row, or a row naming no condition, fails the build.
//
// Conditions fold onto the "Conditional" umbrella term harder than any other
// family. A condition that only reads a number a card already prints — a pool
// size, a discard count, a card name — teaches a player nothing the umbrella
// does not already say, so it names "Conditional" rather than earning a title of
// its own; a per-condition term would turn the rulebook into an API listing. A
// condition that rests on a game concept with a term elsewhere in the registry —
// the tide, a flank, a counter, a keyword — names that term instead, because
// binding is by title across every section.

// conditionFamily is the Condition family's census entry. A condition is
// discovered by its CondText method, which takes no parameters and which no
// other family declares.
func conditionFamily() Family {
	return newFamily(
		"Condition",
		"CondText",
		nil,
		ConditionCatalog(),
		Condition.CondText,
	).gated()
}

// ConditionCatalog returns one representative value of every Condition, with the
// rulebook term each owes. The rows are grouped the way the effect_condition_*.go
// files group the conditions themselves: the combinators, then the checks on the
// board, the discard pile, the card in context, the source card, the tide, and
// the turn so far.
func ConditionCatalog() []Catalogued[Condition] {
	return []Catalogued[Condition]{
		// Combinators and the always-met sentinel.
		{
			Node:   AlwaysMet{},
			Rules:  plumbing("sentinel: always met, and adds no clause of its own"),
			Silent: true,
		},
		{
			Node:  And{Conditions: []Condition{ItIsFriendly{}, ItIsOfTrait{Trait: Dinosaur}}},
			Rules: plumbing("composition: met when every inner condition is met"),
		},
		{
			Node:  Or{Conditions: []Condition{ItIsStunned{}, HasAember{}}},
			Rules: plumbing("composition: met when any inner condition is met"),
		},
		{
			Node:  Not{Cond: OnFlank{}},
			Rules: plumbing("composition: met when its inner condition is not"),
		},
		{
			Node:  CountIs{Count: CreaturesUsed{Player: Controller}, Is: AtLeast, Amount: 3},
			Rules: bears("Conditional"),
		},

		// The board: what each player controls, holds, and has forged.
		{
			Node:  PoolAember{Player: Opponent, Is: AtLeast, Amount: 7},
			Rules: bears("Conditional"),
		},
		{
			Node:  ControlsMoreCreatures{},
			Rules: bears("Conditional"),
		},
		{
			Node:  ControlsNamed{Name: "Velum"},
			Rules: bears("Conditional"),
		},
		{
			Node:  Overwhelmed{},
			Rules: bears("Overwhelmed"),
		},
		{
			Node: HousesRepresented{
				Among:  HousesAmong{Player: Controller, Type: Creature},
				Is:     AtLeast,
				Amount: 3,
			},
			Rules: bears("Belong to House"),
		},
		{
			Node:  CounterInPlay{Kind: CounterDoom},
			Rules: bears("Generic Counters"),
		},
		{
			Node:  NamedCardPurged{Name: "Orbital Observation"},
			Rules: bears("Purge"),
		},
		{
			Node:  ForgedKey{Player: Controller},
			Rules: bears("Turn structure"),
		},
		{
			Node:  HasMoreForgedKeys{Player: Opponent},
			Rules: bears("Turn structure"),
		},
		{
			Node:  KeyColorForged{Player: Controller, Color: KeyColorRed},
			Rules: bears("Turn structure"),
		},
		{
			Node:  ActiveHouseMatchesNoCardsInPlay{},
			Rules: bears("Active House"),
		},

		// The discard pile.
		{
			Node: CardsInDiscardAtLeast{
				House:  HouseMatcher{Kind: MatchNamedHouse, House: Untamed},
				Type:   Creature,
				Amount: 3,
			},
			Rules: bears("Conditional"),
		},
		{
			Node:  DiscardedThisWay{Type: Creature},
			Rules: bears("Discard"),
		},
		{
			Node:  NamedCardInDiscard{Name: "Subtle Chain"},
			Rules: bears("Conditional"),
		},
		{
			Node:  Haunted{},
			Rules: bears("Haunted"),
		},
		{
			Node:  CardsDiscarded{Player: Controller, Amount: 1},
			Rules: bears("Discard"),
		},

		// The card in context: whose it is, what it is, and what it carries.
		{
			Node:  ItIsFriendly{},
			Rules: bears("Conditional"),
		},
		{
			Node:  ItIsEnemy{},
			Rules: bears("Conditional"),
		},
		{
			Node:  ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Mars}, Type: Creature},
			Rules: bears("Belong to House"),
		},
		{
			Node:  ItIsNamed{Name: "Subtle Chain"},
			Rules: bears("Conditional"),
		},
		{
			Node:  ItIsOfTrait{Trait: Dinosaur},
			Rules: bears("Trait"),
		},
		{
			Node:  ItIsNotOfNamedHouse{},
			Rules: bears("Name a House"),
		},
		{
			Node:  ItIsOffIdentity{},
			Rules: bears("Belong to House"),
		},
		{
			Node:  ItIsStunned{},
			Rules: bears("Stun"),
		},
		{
			Node:  ItHasBonusIcon{},
			Rules: bears("Resolve Bonus Icons"),
		},
		{
			Node:  HasAember{Subject: This},
			Rules: bears("Æmber"),
		},
		{
			Node:  ItAttachedToThisOrNeighbor{},
			Rules: bears("Upgrade"),
		},
		{
			Node:  ItIsAmong{Target: Target{Kind: TargetEachFriendlyCreature}.Refine(MostPowerful)},
			Rules: bears("Conditional"),
		},
		{
			Node:  CardsInPlay{Player: Controller, Type: Creature},
			Rules: bears("Conditional"),
		},

		// The source card: where it stands and how it stands.
		{
			Node:  OnFlank{},
			Rules: bears("Flank"),
		},
		{
			Node:  SourceInCenterOfBattleline{},
			Rules: bears("Flank"),
		},
		{
			Node:  SourceHasNoNeighbor{House: HouseMatcher{Kind: MatchNamedHouse, House: Mars}},
			Rules: bears("Flank"),
		},
		{
			Node:  SourceReady{},
			Rules: bears("Ready"),
		},
		{
			Node:  SourceIsFighting{},
			Rules: bears("Fight"),
		},
		{
			Node:  CountersOnThisAtLeast{Kind: CounterFuse, N: 10},
			Rules: bears("Generic Counters"),
		},

		// The tide, whose own term carries the rule the condition reads.
		{
			Node:  TideIsHigh{},
			Rules: bears("Tide"),
		},
		{
			Node:  TideIsLow{},
			Rules: bears("Tide"),
		},

		// The turn so far: what has been played, used, and destroyed.
		{
			Node:  ItIsYourTurn{},
			Rules: bears("Turn structure"),
		},
		{
			Node:  ChoseHouse{House: Dis},
			Rules: bears("Active House"),
		},
		{
			Node:  CardsPlayed{Player: Controller, Amount: 7},
			Rules: bears("Conditional"),
		},
		{
			Node:  FirstCreaturePlayedThisTurn{},
			Rules: bears("Conditional"),
		},
		{
			Node:  NoCreaturesPlayedThisTurn{},
			Rules: bears("Conditional"),
		},
		{
			Node:  CreatureDestroyedThisTurn{Player: Opponent},
			Rules: bears("Destroy"),
		},
		{
			Node:  UsedCreatureToReap{},
			Rules: bears("Reap"),
		},
		{
			Node:  UsedCreatureToFight{},
			Rules: bears("Fight"),
		},
		{
			Node:  UsedNoCreatures{},
			Rules: bears("Use a Card"),
		},
		{
			Node:  FirstReapOfTurn{},
			Rules: bears("Reap"),
		},
		{
			Node:  SourceFirstUseThisTurn{},
			Rules: bears("Use a Card"),
		},
		{
			Node:  AemberStolenFromYou{},
			Rules: bears("Steal Æmber"),
		},
		{
			Node:  MovedAnyAember{},
			Rules: bears("Move Æmber"),
		},
		{
			Node:  ArchivedCreaturesShareHouse{},
			Rules: bears("Belong to House"),
		},
	}
}

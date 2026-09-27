package engine

// This file is the census of the Count family: one row per count the package
// declares, each naming the rulebook term it owes or the reason it owes none.
// See catalog.go for the classification rule, and effect_count.go for the
// interface the counts implement. TestFamilyTotality keeps it complete — a count
// with no row, or a row naming no count, fails the build.
//
// Counts fold the way conditions do. "For each friendly creature" teaches a
// player nothing past the "For Each" umbrella term, so a count that only totals
// cards or events names that title. A count that rests on a game concept with a
// term elsewhere in the registry — the Æmber on a card, a counter kind, a purge,
// a forged key — names that term instead, because binding is by title across
// every section.

// countFamily is the Count family's census entry. A count is discovered by both
// of its methods: Value, whose *EffectContext parameter separates it from the
// other families' same-named methods, and CountText, which is what tells a count
// from a condition that merely exposes a Value for CountIs to read.
func countFamily() Family {
	return newFamily(
		"Count",
		"Value",
		[]string{"*EffectContext"},
		CountCatalog(),
		Count.CountText,
	).also(MethodSpec{Method: "CountText"}).gated()
}

// CountCatalog returns one representative value of every Count, with the
// rulebook term each owes. The rows are grouped by what they measure: the
// constant, the Æmber and key economies, the board, the turn so far, and the
// tallies an earlier effect in the same resolution left behind.
func CountCatalog() []Catalogued[Count] {
	return []Catalogued[Count]{
		// The constant, which the effect that holds it prints itself.
		{
			Node:   Fixed(3),
			Rules:  plumbing("literal: a constant the effect prints as its own number"),
			Silent: true,
		},

		// The Æmber economy.
		{
			Node:  AemberInPool{Player: Opponent},
			Rules: bears("Æmber"),
		},
		{
			Node:  AemberOnThis{},
			Rules: bears("Æmber"),
		},
		{
			Node:  AemberOnFriendlyCreatures{},
			Rules: bears("Æmber"),
		},
		{
			Node:  AemberStolenThisEvent{},
			Rules: bears("Steal Æmber"),
		},

		// Keys.
		{
			Node:  ForgedKeys{Player: Opponent},
			Rules: bears("Turn structure"),
		},
		{
			Node:  UnforgedKeys{Player: Controller},
			Rules: bears("Turn structure"),
		},

		// The board: what is in play, and how the sides compare.
		{
			Node:  CardsInPlay{Player: Controller, Type: Creature},
			Rules: bears("For Each"),
		},
		{
			Node:  ExcessCreatures{Player: Opponent},
			Rules: bears("For Each"),
		},
		{
			Node:  ArtifactsInPlay{},
			Rules: bears("Artifact"),
		},
		{
			Node:  UpgradesOn{Target: Target{Kind: TargetThisCreature}},
			Rules: bears("Upgrade"),
		},
		{
			Node:  HousesInPlay{Except: Sanctum},
			Rules: bears("Belong to House"),
		},
		{
			Node:  HousesAmong{Player: Controller, Type: Creature},
			Rules: bears("Belong to House"),
		},

		// One card's own measurements: what sits on it, and what it is.
		{
			Node:  DamageOnThis{},
			Rules: bears("Damage"),
		},
		{
			Node:  CountersOnThis{Kind: CounterFuse},
			Rules: bears("Generic Counters"),
		},
		{
			Node:  PowerCountersOnThis{},
			Rules: bears("Power Counter"),
		},
		{
			Node:  PowerOfChosen{},
			Rules: bears("Power and Armor"),
		},
		{
			Node:  TraitsOfChosen{},
			Rules: bears("Trait"),
		},
		{
			Node:  BonusIconsOf{Over: ThePurgedCards{}, Kind: BonusAember},
			Rules: bears("Resolve Bonus Icons"),
		},

		// The battleline around a card.
		{
			Node:  NeighborsOfThis{},
			Rules: bears("Battleline Position"),
		},
		{
			Node:  NeighborsMatching{House: HouseMatcher{Kind: MatchNamedHouse, House: Mars}},
			Rules: bears("Battleline Position"),
		},
		{
			Node:  CombinedPowerOfNeighborsWithout{Without: Changeling},
			Rules: bears("Battleline Position"),
		},

		// Zones: how many cards a pile holds.
		{
			Node:  CardsInHand{Player: Opponent, House: AnyHouse},
			Rules: bears("For Each"),
		},
		{
			Node:  CardsInZone{Zone: Archives, Player: Controller},
			Rules: bears("For Each"),
		},
		{
			Node:  CopiesInDiscard{},
			Rules: bears("For Each"),
		},
		{
			Node:  PurgedCards{},
			Rules: bears("Purge"),
		},

		// The turn so far.
		{
			Node:  CardsPlayed{Player: Controller},
			Rules: bears("For Each"),
		},
		{
			Node:  CreaturesUsed{Player: Controller},
			Rules: bears("Use a Card"),
		},
		{
			Node:  TurnCount{Player: Controller, Of: CreaturesReapedThisTurn},
			Rules: bears("For Each"),
		},

		// Tallies an earlier effect in the same resolution left behind.
		{
			Node:  CardsDestroyed{},
			Rules: bears("Destroy"),
		},
		{
			Node:  CreaturesDestroyed{},
			Rules: bears("Destroy"),
		},
		{
			Node:  PowerDestroyedThisWay{},
			Rules: bears("Destroy"),
		},
		{
			Node:  CardsShuffledIntoDeck{},
			Rules: bears("Shuffle"),
		},
		{
			Node:  ProducedThisWay{Tally: TallyCardsPurged, Player: Controller},
			Rules: bears("For Each"),
		},
		{
			Node:  CardsPurged{Type: Creature},
			Rules: bears("Purge"),
		},
		{
			Node:  CardsRevealed{},
			Rules: bears("Reveal a Card"),
		},
		{
			Node:  CreaturesHealed{},
			Rules: bears("Heal"),
		},
		{
			Node:  DamageHealed{},
			Rules: bears("Heal"),
		},
		{
			Node:  DamagePrevented{},
			Rules: bears("Armor"),
		},
	}
}

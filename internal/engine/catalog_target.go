package engine

// This file is the census of Target, which is the odd member of the census: a
// Target is a flag struct rather than an interface, because ADR 0005 keeps it
// comparable, so it has no implementations to scan for. Its members are instead
// its halves — the TargetKind constants that pick a base set, the filter builders
// that narrow one, and the axis values those builders write — and the census
// covers each half with the machinery that fits it: the kinds and the axis values
// as Enums, discovered as the constants of their types, and the filters as a
// Family discovered by shape, every exported method on Target returning a Target.
// An added kind, an added axis value and an added filter are all a red build
// until they have a row.
//
// An axis value and the builder that writes it are catalogued separately and
// carry the same term, because neither implies the other: a builder could be
// dropped without its axis losing a value, and an axis can gain a value no
// builder writes.
//
// Nearly every row names the one "Target" umbrella term. "Each enemy creature"
// and "a friendly creature with power 3 or lower" are the same rule about how
// card text names what an effect reaches, and a term per filter would turn the
// rulebook into an API listing. A filter that rests on a concept with a term of
// its own — a flank, a trait, a house, Æmber on a card, a counter — names that
// term instead, because binding is by title across every section.

// targetBase is the base set the census applies each filter to, standing in for
// the kind a card would pair it with. It is the widest board-scanning kind, so
// every positional and per-card filter has something to narrow.
var targetBase = Target{Kind: TargetEachCreature}

// targetFilterFamily is the Target filter builders' census entry. They are
// discovered by shape rather than by a method name: an exported method on Target
// that returns a Target is a filter a card can write.
func targetFilterFamily() Family {
	return newFamily(
		"Target",
		"",
		nil,
		TargetFilterCatalog(),
		Target.Text,
	).builds("Target").gated()
}

// TargetFilterCatalog returns one representative Target per filter builder — the
// base set with that filter applied — with the rulebook term each owes. The rows
// are grouped by what the filter reads: identity, power, condition, position, and
// the set-relative refinement.
func TargetFilterCatalog() []Catalogued[Target] {
	return []Catalogued[Target]{
		// Identity: trait, house, and printed name.
		{
			Name:  "WithTrait",
			Node:  targetBase.WithTrait(Scientist),
			Rules: bears("Trait"),
		},
		{
			Name:  "ExceptTrait",
			Node:  targetBase.ExceptTrait(Scientist),
			Rules: bears("Trait"),
		},
		{
			Name:  "SharingTrait",
			Node:  targetBase.SharingTrait(),
			Rules: bears("Trait"),
		},
		{
			Name:  "House",
			Node:  targetBase.House(HouseMatcher{Kind: MatchNamedHouse, House: Mars}),
			Rules: bears("Belong to House"),
		},
		{
			Name: "MatchingAny",
			Node: targetBase.
				House(HouseMatcher{Kind: MatchNamedHouse, House: Mars}).
				WithTrait(Robot).
				MatchingAny(),
			Rules: bears("Belong to House"),
		},
		{
			Name:  "OfHouseWithMostCreatures",
			Node:  targetBase.OfHouseWithMostCreatures(),
			Rules: bears("Belong to House"),
		},
		{
			Name:  "SharesHouseWithNeighbors",
			Node:  targetBase.SharesHouseWithNeighbors(1),
			Rules: bears("Belong to House"),
		},
		{
			Name:  "Named",
			Node:  targetBase.Named("Ancient Bear"),
			Rules: bears("Target"),
		},
		{
			Name:  "Keyword",
			Node:  targetBase.Keyword(Taunt),
			Rules: bears("Target"),
		},
		{
			Name:  "Other",
			Node:  targetBase.Other(),
			Rules: bears("Target"),
		},

		// Power.
		{
			Name:  "PowerAtMost",
			Node:  targetBase.PowerAtMost(3),
			Rules: bears("Power Threshold"),
		},
		{
			Name:  "PowerAtLeast",
			Node:  targetBase.PowerAtLeast(5),
			Rules: bears("Power Threshold"),
		},
		{
			Name:  "PowerExactly",
			Node:  targetBase.PowerExactly(4),
			Rules: bears("Power Threshold"),
		},
		{
			Name:  "OddPower",
			Node:  targetBase.OddPower(),
			Rules: bears("Power Threshold"),
		},
		{
			Name:  "EvenPower",
			Node:  targetBase.EvenPower(),
			Rules: bears("Power Threshold"),
		},

		// What the card carries or has had done to it.
		{
			Name:  "Damaged",
			Node:  targetBase.Damaged(),
			Rules: bears("Damage"),
		},
		{
			Name:  "Undamaged",
			Node:  targetBase.Undamaged(),
			Rules: bears("Damage"),
		},
		{
			Name:  "WithArmor",
			Node:  targetBase.WithArmor(),
			Rules: bears("Armor"),
		},
		{
			Name:  "WithAember",
			Node:  targetBase.WithAember(),
			Rules: bears("Æmber"),
		},
		{
			Name:  "WithoutAember",
			Node:  targetBase.WithoutAember(),
			Rules: bears("Æmber"),
		},
		{
			Name:  "WithUpgrade",
			Node:  targetBase.WithUpgrade(),
			Rules: bears("Upgrade"),
		},
		{
			Name:  "WithCounter",
			Node:  targetBase.WithCounter(CounterDoom),
			Rules: bears("Generic Counters"),
		},
		{
			Name:  "WithoutBonusIcons",
			Node:  targetBase.WithoutBonusIcons(),
			Rules: bears("Resolve Bonus Icons"),
		},
		{
			Name:  "Stunned",
			Node:  targetBase.Stunned(),
			Rules: bears("Stun"),
		},
		{
			Name:  "Ready",
			Node:  targetBase.Ready(),
			Rules: bears("Ready"),
		},

		// Position in the battleline.
		{
			Name:  "OnFlank",
			Node:  targetBase.OnFlank(),
			Rules: bears("Flank"),
		},
		{
			Name:  "NotOnFlank",
			Node:  targetBase.NotOnFlank(),
			Rules: bears("Flank"),
		},
		{
			Name:  "InCenter",
			Node:  targetBase.InCenter(),
			Rules: bears("Battleline Position"),
		},
		{
			Name:  "Neighboring",
			Node:  targetBase.Neighboring(),
			Rules: bears("Battleline Position"),
		},
		{
			Name:  "AndNeighbors",
			Node:  Target{Kind: TargetChosenCreature}.AndNeighbors(),
			Rules: bears("Battleline Position"),
		},
		{
			Name:  "NeighborsOf",
			Node:  Target{Kind: TargetChosenCreature}.NeighborsOf(),
			Rules: bears("Battleline Position"),
		},
		{
			Name:  "ToRightOfSource",
			Node:  targetBase.ToRightOfSource(),
			Rules: bears("Battleline Position"),
		},
		{
			Name:  "ToLeftOfSource",
			Node:  targetBase.ToLeftOfSource(),
			Rules: bears("Battleline Position"),
		},

		// The set-relative refinement, which carries its own census and its own term.
		{
			Name:  "Refine",
			Node:  targetBase.Refine(MostPowerful),
			Rules: plumbing("composition: defers to a Refinement, which owes its own term"),
		},
	}
}

// targetKindEnum is the TargetKind census: the base set a Target picks before any
// filter narrows it. A kind that names a card the context already holds binds to
// the term for the relation that put it there — the fight, the upgrade, the
// ability grant, the battleline — and the rest name the Target umbrella.
func targetKindEnum() Enum {
	return newEnum(
		"TargetKind",
		[]string{"targetUnset", "targetKindCount"},
		len(TargetKinds()),
		targetKindRows(),
		func(k TargetKind) string { return Target{Kind: k}.Text() },
	)
}

// targetKindRows is the TargetKind catalog, split out so the enum entry reads as
// its discovery spec and this reads as the table.
func targetKindRows() []Catalogued[TargetKind] {
	return []Catalogued[TargetKind]{
		// The source card and the card the context holds.
		{Name: "TargetThisCreature", Node: TargetThisCreature, Rules: bears("Target")},
		{
			Name:  "TargetTriggeringCreature",
			Node:  TargetTriggeringCreature,
			Rules: bears("Target"),
		},
		{Name: "TargetTheSameCreature", Node: TargetTheSameCreature, Rules: bears("Target")},
		{Name: "TargetTheChosenCreature", Node: TargetTheChosenCreature, Rules: bears("Target")},
		{Name: "TargetTheOtherCreature", Node: TargetTheOtherCreature, Rules: bears("Target")},

		// Cards the fight, the attachment, or the grant put in reach.
		{Name: "TargetCreatureFought", Node: TargetCreatureFought, Rules: bears("Fight")},
		{Name: "TargetTheFoughtCreature", Node: TargetTheFoughtCreature, Rules: bears("Fight")},
		{Name: "TargetAttachedHost", Node: TargetAttachedHost, Rules: bears("Upgrade")},
		{
			Name:  "TargetEachUpgradeOnThis",
			Node:  TargetEachUpgradeOnThis,
			Rules: bears("Upgrade"),
		},
		{Name: "TargetGrantingCard", Node: TargetGrantingCard, Rules: bears("Gain an Ability")},
		{
			Name:  "TargetEachNeighbor",
			Node:  TargetEachNeighbor,
			Rules: bears("Battleline Position"),
		},
		{
			Name:  "TargetFormerNeighbors",
			Node:  TargetFormerNeighbors,
			Rules: bears("Battleline Position"),
		},

		// Whole sets on the board.
		{Name: "TargetEachCreature", Node: TargetEachCreature, Rules: bears("Target")},
		{
			Name:  "TargetEachFriendlyCreature",
			Node:  TargetEachFriendlyCreature,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetEachEnemyCreature",
			Node:  TargetEachEnemyCreature,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetEachOtherFriendlyCreature",
			Node:  TargetEachOtherFriendlyCreature,
			Rules: bears("Target"),
		},
		{Name: "TargetEachArtifact", Node: TargetEachArtifact, Rules: bears("Target")},
		{
			Name:  "TargetEachFriendlyArtifact",
			Node:  TargetEachFriendlyArtifact,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetEachEnemyArtifact",
			Node:  TargetEachEnemyArtifact,
			Rules: bears("Target"),
		},
		{Name: "TargetEachCardInPlay", Node: TargetEachCardInPlay, Rules: bears("Target")},
		{
			Name:  "TargetEachFriendlyCardInPlay",
			Node:  TargetEachFriendlyCardInPlay,
			Rules: bears("Target"),
		},

		// Sets the controller chooses one card from.
		{Name: "TargetChosenCreature", Node: TargetChosenCreature, Rules: bears("Target")},
		{
			Name:  "TargetChosenEnemyCreature",
			Node:  TargetChosenEnemyCreature,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetChosenFriendlyCreature",
			Node:  TargetChosenFriendlyCreature,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetChosenOtherCreature",
			Node:  TargetChosenOtherCreature,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetChosenOtherFriendlyCreature",
			Node:  TargetChosenOtherFriendlyCreature,
			Rules: bears("Target"),
		},
		{Name: "TargetChosenArtifact", Node: TargetChosenArtifact, Rules: bears("Target")},
		{
			Name:  "TargetChosenEnemyArtifact",
			Node:  TargetChosenEnemyArtifact,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetChosenFriendlyArtifact",
			Node:  TargetChosenFriendlyArtifact,
			Rules: bears("Target"),
		},
		{Name: "TargetChosenUpgrade", Node: TargetChosenUpgrade, Rules: bears("Upgrade")},
		{
			Name:  "TargetChosenCreatureOrArtifact",
			Node:  TargetChosenCreatureOrArtifact,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetChosenFriendlyCreatureOrArtifact",
			Node:  TargetChosenFriendlyCreatureOrArtifact,
			Rules: bears("Target"),
		},
		{
			Name:  "TargetChosenEnemyCreatureOrArtifact",
			Node:  TargetChosenEnemyCreatureOrArtifact,
			Rules: bears("Target"),
		},
	}
}

// powerBoundKindEnum is the PowerBoundKind census: the comparison a target's
// power filter makes. Every kind names the one "Power Threshold" term — whether
// a bound reads a number or a parity, the rule it rests on is how a card names
// the powers it reaches.
func powerBoundKindEnum() Enum {
	return newEnum(
		"PowerBoundKind",
		[]string{"boundUnset"},
		len(allPowerBoundKinds()),
		[]Catalogued[PowerBoundKind]{
			{Name: "BoundAtMost", Node: BoundAtMost, Rules: bears("Power Threshold")},
			{Name: "BoundAtLeast", Node: BoundAtLeast, Rules: bears("Power Threshold")},
			{Name: "BoundExactly", Node: BoundExactly, Rules: bears("Power Threshold")},
			{Name: "BoundOdd", Node: BoundOdd, Rules: bears("Power Threshold")},
			{Name: "BoundEven", Node: BoundEven, Rules: bears("Power Threshold")},
		},
		func(k PowerBoundKind) string {
			return Target{Kind: TargetEachCreature, power: PowerBound{Kind: k, Amount: 3}}.Text()
		},
	)
}

// damagePresenceEnum is the DamagePresence census: whether a target's damage
// filter wants a creature carrying damage or one carrying none. Both name the
// damage rule that puts it there.
func damagePresenceEnum() Enum {
	return newEnum(
		"DamagePresence",
		[]string{"damageAny"},
		len(allDamagePresences()),
		[]Catalogued[DamagePresence]{
			{Name: "DamageSome", Node: DamageSome, Rules: bears("Damage")},
			{Name: "DamageNone", Node: DamageNone, Rules: bears("Damage")},
		},
		func(d DamagePresence) string {
			return Target{Kind: TargetEachCreature, damage: d}.Text()
		},
	)
}

// aemberPresenceEnum is the AemberPresence census: whether a target's Æmber
// filter wants a card with Æmber on it or one with none. Both name the term for
// Æmber sitting on a card.
func aemberPresenceEnum() Enum {
	return newEnum(
		"AemberPresence",
		[]string{"aemberAny"},
		len(allAemberPresences()),
		[]Catalogued[AemberPresence]{
			{Name: "AemberSome", Node: AemberSome, Rules: bears("Æmber")},
			{Name: "AemberNone", Node: AemberNone, Rules: bears("Æmber")},
		},
		func(a AemberPresence) string {
			return Target{Kind: TargetEachCreature, aember: a}.Text()
		},
	)
}

// positionEnum is the Position census: the battleline place a target narrows to.
// The two flank values name the flank rule; the rest name the battleline
// positions term, which is where a line's center and sides are taught.
func positionEnum() Enum {
	return newEnum(
		"Position",
		[]string{"positionAny"},
		len(allPositions()),
		[]Catalogued[Position]{
			{Name: "PositionOnFlank", Node: PositionOnFlank, Rules: bears("Flank")},
			{Name: "PositionNotOnFlank", Node: PositionNotOnFlank, Rules: bears("Flank")},
			{
				Name:  "PositionCenter",
				Node:  PositionCenter,
				Rules: bears("Battleline Position"),
			},
			{
				Name:  "PositionRightOfSource",
				Node:  PositionRightOfSource,
				Rules: bears("Battleline Position"),
			},
			{
				Name:  "PositionLeftOfSource",
				Node:  PositionLeftOfSource,
				Rules: bears("Battleline Position"),
			},
		},
		func(p Position) string {
			return Target{Kind: TargetEachCreature, position: p}.Text()
		},
	)
}

// exclusionEnum is the Exclusion census: the one card a target leaves out. Each
// value names the Target umbrella, because all three print the same "other" and
// differ only in which card the exclusion is aimed at.
func exclusionEnum() Enum {
	return newEnum(
		"Exclusion",
		[]string{"excludeNone"},
		len(allExclusions()),
		[]Catalogued[Exclusion]{
			{Name: "ExcludeSource", Node: ExcludeSource, Rules: bears("Target")},
			{Name: "ExcludeIt", Node: ExcludeIt, Rules: bears("Target")},
			{Name: "ExcludeFocus", Node: ExcludeFocus, Rules: bears("Target")},
		},
		func(e Exclusion) string {
			return Target{Kind: TargetEachCreature, exclusion: e}.Text()
		},
	)
}

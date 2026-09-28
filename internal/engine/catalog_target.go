package engine

// This file is the census of Target, which is the odd member of the census: a
// Target is a flag struct rather than an interface, because ADR 0005 keeps it
// comparable, so it has no implementations to scan for. Its members are instead
// its halves — the TargetKind constants that pick a base set, the builders that
// narrow one, and the values of the axes they narrow on — and the census covers
// each half with the machinery that fits it: the kinds and the axis values as
// Enums, discovered as the constants of their types, and the builders as a Family
// discovered by shape, every exported method on Target returning a Target. An
// added kind, an added axis value and an added builder are all a red build until
// they have a row.
//
// Target once carried one builder per filter, and those builders are gone: a card
// writes its narrowing as a Filter literal, so the axes are catalogued as the
// fields of Filter in catalog_filter.go. What the Target family still covers is
// the four methods left — With, Refine, AndNeighbors, NeighborsOf — each of which
// defers to a value with a census of its own, so their rows say which census owns
// the words rather than owning any. The family stays because the discovery is what
// makes a newly added Target method a red build.
//
// An axis value and the field that carries it are catalogued separately and carry
// the same term, because neither implies the other: a field could be dropped
// without its axis losing a value, and an axis can gain a value no field writes.
//
// Nearly every row names the one "Target" umbrella term. "Each enemy creature"
// and "a friendly creature with power 3 or lower" are the same rule about how
// card text names what an effect reaches, and a term per axis would turn the
// rulebook into an API listing. An axis that rests on a concept with a term of
// its own — a flank, a trait, a house, Æmber on a card, a counter — names that
// term instead, because binding is by title across every section.

// targetBuilderFamily is the Target builders' census entry. They are discovered by
// shape rather than by a method name: an exported method on Target that returns a
// Target is a narrowing a card can write. Every survivor defers — to a Filter, a
// Refinement, or a NeighborMode — so every row is plumbing pointing at the census
// that owns its words.
func targetBuilderFamily() Family {
	return newFamily(
		"Target",
		"",
		nil,
		TargetBuilderCatalog(),
		Target.Text,
	).builds("Target").gated()
}

// TargetBuilderCatalog returns one representative Target per surviving builder —
// the base set with that builder applied — with the census each defers to.
func TargetBuilderCatalog() []Catalogued[Target] {
	base := Target{Kind: TargetEachCreature}
	return []Catalogued[Target]{
		{
			Name: "With",
			Node: base.With(Filter{Trait: Scientist}),
			Rules: plumbing(
				"composition: writes a Filter, whose axes are catalogued as the " +
					"Filter family"),
		},
		{
			Name: "Refine",
			Node: base.Refine(MostPowerful),
			Rules: plumbing(
				"composition: defers to a Refinement, which owes its own term"),
		},
		{
			Name: "AndNeighbors",
			Node: Target{Kind: TargetChosenCreature}.AndNeighbors(),
			Rules: plumbing(
				"composition: writes a NeighborMode, which owes its own term"),
		},
		{
			Name: "NeighborsOf",
			Node: Target{Kind: TargetChosenCreature}.NeighborsOf(),
			Rules: plumbing(
				"composition: writes a NeighborMode, which owes its own term"),
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
			{
				Name:  "BoundLessThanSource",
				Node:  BoundLessThanSource,
				Rules: bears("Power Threshold"),
			},
		},
		func(k PowerBoundKind) string {
			return Target{
				Kind:   TargetEachCreature,
				Filter: Filter{Power: PowerBound{Kind: k, Amount: 3}},
			}.Text()
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
			return Target{Kind: TargetEachCreature, Filter: Filter{Damage: d}}.Text()
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
			return Target{Kind: TargetEachCreature, Filter: Filter{Aember: a}}.Text()
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
			return Target{Kind: TargetEachCreature, Filter: Filter{Position: p}}.Text()
		},
	)
}

// neighborModeEnum is the NeighborMode census: how a target reaches the
// battleline neighbors of what it selects. Both modes name the battleline
// positions term, which is where adjacency is taught.
func neighborModeEnum() Enum {
	return newEnum(
		"NeighborMode",
		[]string{"NeighborsNone"},
		len(NeighborModes()),
		[]Catalogued[NeighborMode]{
			{
				Name:  "NeighborsIncluded",
				Node:  NeighborsIncluded,
				Rules: bears("Battleline Position"),
			},
			{
				Name:  "NeighborsOnly",
				Node:  NeighborsOnly,
				Rules: bears("Battleline Position"),
			},
		},
		func(m NeighborMode) string {
			return Target{Kind: TargetChosenCreature, Neighbors: m}.Text()
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
			return Target{Kind: TargetEachCreature, Filter: Filter{Except: e}}.Text()
		},
	)
}

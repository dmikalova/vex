package engine

// This file is the census of Filter, the per-card half of a Target. A Filter is
// a comparable struct rather than an interface (ADR 0005), and a card writes its
// narrowing as a struct literal, so its members are neither types nor builder
// methods but the struct's own exported fields — one per axis. The census
// discovers them by scanning the source for those fields, so an added axis is a
// red build until it has a row.
//
// The rulebook binding follows the same policy as the Target census beside it.
// Nearly every row names the one "Target" umbrella term: "each enemy creature"
// and "a friendly creature with power 3 or lower" are the same rule about how
// card text names what an effect reaches, and a term per axis would turn the
// rulebook into an API listing. An axis that rests on a concept with a term of
// its own — a flank, a trait, a house, Æmber on a card, a counter — names that
// term instead, because binding is by title across every section.
//
// The axis values a field carries are catalogued separately, as Enums in
// catalog_target.go, because neither implies the other: an axis can gain a value
// without gaining a field, and a field could be dropped without its axis losing a
// value.

// filterFamily is the Filter axes' census entry. They are discovered as the
// exported fields of the Filter struct, each of which is one axis a card writes.
func filterFamily() Family {
	return newFamily(
		"Filter",
		"",
		nil,
		FilterCatalog(),
		func(f Filter) string {
			return Target{Kind: TargetEachCreature, Filter: f}.Text()
		},
	).fieldsOf().gated()
}

// FilterCatalog returns one representative Filter per axis — the axis set on its
// own — with the rulebook term each owes. The rows are grouped the way the struct
// is: the identity axes a card carries wherever it sits, then the in-play axes
// that read the board.
func FilterCatalog() []Catalogued[Filter] {
	return append(filterIdentityRows(), filterInPlayRows()...)
}

// filterIdentityRows are the axes decided by what a card is — its type, its
// traits, its house, its printed name — which are true of it in any zone.
func filterIdentityRows() []Catalogued[Filter] {
	return []Catalogued[Filter]{
		{
			Name:  "Type",
			Node:  Filter{Type: Artifact},
			Rules: bears("Target"),
		},
		{
			Name:  "Gigantic",
			Node:  Filter{Gigantic: true},
			Rules: bears("Gigantic creature"),
		},
		{
			Name:  "Trait",
			Node:  Filter{Trait: Scientist},
			Rules: bears("Trait"),
		},
		{
			Name:  "ExceptTrait",
			Node:  Filter{ExceptTrait: Scientist},
			Rules: bears("Trait"),
		},
		{
			Name:  "House",
			Node:  Filter{House: HouseMatcher{Kind: MatchNamedHouse, House: Mars}},
			Rules: bears("Belong to House"),
		},
		{
			Name: "MatchAny",
			Node: Filter{
				House:    HouseMatcher{Kind: MatchNamedHouse, House: Mars},
				Trait:    Robot,
				MatchAny: true,
			},
			Rules: bears("Belong to House"),
		},
		{
			Name:  "HouseWithMostCreatures",
			Node:  Filter{HouseWithMostCreatures: true},
			Rules: bears("Belong to House"),
		},
		{
			Name:  "SharesTrait",
			Node:  Filter{SharesTrait: true},
			Rules: bears("Trait"),
		},
		{
			Name:  "Name",
			Node:  Filter{Name: "Ancient Bear"},
			Rules: bears("Target"),
		},
		{
			Name:  "Except",
			Node:  Filter{Except: ExcludeSource},
			Rules: bears("Target"),
		},
	}
}

// filterInPlayRows are the axes decided by what has happened to a card and where
// it stands, which mean something only of a card in play.
func filterInPlayRows() []Catalogued[Filter] {
	return []Catalogued[Filter]{
		{
			Name:  "Power",
			Node:  Filter{Power: PowerBound{Kind: BoundAtMost, Amount: 3}},
			Rules: bears("Power Threshold"),
		},
		{
			Name:  "Damage",
			Node:  Filter{Damage: DamageSome},
			Rules: bears("Damage"),
		},
		{
			Name:  "Stunned",
			Node:  Filter{Stunned: true},
			Rules: bears("Stun"),
		},
		{
			Name:  "Ready",
			Node:  Filter{Ready: true},
			Rules: bears("Ready"),
		},
		{
			Name:  "Aember",
			Node:  Filter{Aember: AemberSome},
			Rules: bears("Æmber"),
		},
		{
			Name:  "NoBonusIcons",
			Node:  Filter{NoBonusIcons: true},
			Rules: bears("Resolve Bonus Icons"),
		},
		{
			Name:  "Counter",
			Node:  Filter{Counter: CounterDoom},
			Rules: bears("Generic Counters"),
		},
		{
			Name:  "Armor",
			Node:  Filter{Armor: true},
			Rules: bears("Armor"),
		},
		{
			Name:  "Upgrade",
			Node:  Filter{Upgrade: true},
			Rules: bears("Upgrade"),
		},
		{
			Name:  "HouseWithAtLeast",
			Node:  Filter{HouseWithAtLeast: 3},
			Rules: bears("Creatures of a House"),
		},
		{
			Name:  "WithoutSharedTrait",
			Node:  Filter{WithoutSharedTrait: true},
			Rules: bears("Trait"),
		},
		{
			Name:  "SharesHouseWithNeighbors",
			Node:  Filter{SharesHouseWithNeighbors: 1},
			Rules: bears("Belong to House"),
		},
		{
			Name:  "Keyword",
			Node:  Filter{Keyword: Taunt},
			Rules: bears("Target"),
		},
		{
			Name:  "Position",
			Node:  Filter{Position: PositionOnFlank},
			Rules: bears("Flank"),
		},
		{
			Name:  "Neighboring",
			Node:  Filter{Neighboring: true},
			Rules: bears("Battleline Position"),
		},
	}
}

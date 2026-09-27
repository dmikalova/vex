package engine

import "reflect"

// This file holds the census machinery: the row type every catalog_<family>.go
// file is written in, the rules classification each row carries, and the
// registry of families the totality tests and `mage tool:census` walk.
//
// The census exists because the rulebook is complete by construction (ADR 0018):
// every member of a closed catalog owes a rulebook term, and the build fails
// when one has none. The AST node families — Effect and the strategies beside it
// — are struct types rather than an enum, so nothing enumerated them and a node
// could ship a player-facing rule no term described. A catalog gives each family
// the enumeration the enums already had, proved complete against the package's
// own source the way LogEntrySamples is (ADR 0046).

// Catalogued is one census row: a constructed node of the family, and the
// classification saying whether it owes a rulebook term. The row carries a real
// value rather than a type name so a consumer — a totality test, a Visitor, the
// style gallery — can render and walk the node it names.
type Catalogued[T any] struct {
	// Node is a representative value of the node type, filled only as far as it
	// needs to render and resolve sensibly.
	Node T
	// Rules classifies the node: the term it owes, or the reason it owes none.
	Rules RulesBearing
}

// RulesBearing is a node's rules classification, and the record of the judgement
// that produced it. Exactly one of its two columns is set; both empty and both
// set are equally invalid, and the census tests fail on either.
//
// A node is plumbing when its text contributes no vocabulary of its own — it
// only joins, repeats, gates or re-aims its children, and everything a player
// reads comes from those children (Sequence, Then, ForEach, Repeat, the duration
// wrappers). Such a node sets NoTerm with the reason.
//
// A node is rule-bearing when it prints a word a player must be taught in
// order to predict the outcome, which is why May, Choose One and the result gate
// keep the terms they already have. Such a node sets Term to the rulebook title
// it owes. Several nodes may name one title — the Archive nodes are all
// "Archive" — and a title may live in any section, so the ward nodes point at the
// Keyword section rather than restating the rule.
//
// When in doubt, classify rule-bearing and write the term: a thin term is
// recoverable, a wrong NoTerm reopens the hole the census exists to close.
type RulesBearing struct {
	// Term is the rulebook title this node owes, matched against RuleTerms by
	// title across every section.
	Term string
	// NoTerm is the one-line reason this node owes no term, e.g.
	// "composition: resolves its children in order". It is required rather than
	// optional so the table reads as a decision record a later reader can judge.
	NoTerm string
}

// bears classifies a node as rule-bearing, owing the rulebook term titled title.
func bears(title string) RulesBearing { return RulesBearing{Term: title} }

// plumbing classifies a node as owing no rulebook term, recording why.
func plumbing(reason string) RulesBearing { return RulesBearing{NoTerm: reason} }

// Family is one catalogued node family, erased to the columns the census works
// in: how the source scan discovers the family's members, and one row per
// catalogued node. Families returns them all.
type Family struct {
	// Name is the family's interface name, e.g. "Refinement".
	Name string
	// Method is the interface method a member declares, which is how the source
	// scan discovers the family.
	Method string
	// Params are the method's parameter types as written in source, which
	// separate two families whose methods share a name.
	Params []string
	// Rows is the family's census, in catalog order.
	Rows []FamilyRow
	// Gated reports whether the family's totality test is switched on. A family
	// part-way through being catalogued is reported by `mage tool:census` but
	// does not yet fail the build; the task that completes its catalog gates it.
	Gated bool
}

// FamilyRow is a Catalogued row with its node erased to what the census asks of
// it: the type name the source scan matches against, the classification, and the
// node's own rendered text.
type FamilyRow struct {
	// Type is the concrete node type's name, e.g. "mostPowerfulN".
	Type string
	// Rules is the row's classification.
	Rules RulesBearing
	// Text renders the node the way a card prints it, so the census can check
	// that the catalogued value says something rather than standing in as an
	// empty literal.
	Text func() string
}

// Families returns every catalogued node family. Each catalog_<family>.go
// contributes one entry, so a consumer that walks the whole census — the
// totality tests, `mage tool:census` — picks up a new family without being
// edited.
func Families() []Family {
	return append([]Family{
		effectFamily(),
		refinementFamily(),
	}, uncataloguedFamilies()...)
}

// uncataloguedFamilies are the node families the census has not reached yet:
// their discovery spec with no rows, so `mage tool:census` reports the whole
// family as a gap and the task that catalogues one replaces its entry here with
// a catalog_<family>.go. The list is what keeps a family from being forgotten
// rather than merely uncatalogued.
func uncataloguedFamilies() []Family {
	return []Family{
		{Name: "Condition", Method: "CondText"},
		{Name: "Count", Method: "Value", Params: []string{"*EffectContext"}},
		{Name: "CreatureVerb", Method: "VerbText"},
		{Name: "Selection", Method: "pick", Params: []string{"*EffectContext", "[]LocalID"}},
		{Name: "Spread", Method: "hits", Params: []string{"*EffectContext"}},
		{Name: "TopAct", Method: "terminal"},
		{
			Name:   "PerTarget",
			Method: "perTargetValue",
			Params: []string{"*EffectContext", "LocalID"},
		},
		{Name: "Loss", Method: "lose", Params: []string{"int"}},
		{Name: "RepeatGate", Method: "run", Params: []string{"*EffectContext", "Effect"}},
		{Name: "Gather", Method: "gather", Params: []string{"*EffectContext"}},
		{Name: "Quantity", Method: "picks", Params: []string{"*EffectContext"}},
		{Name: "BonusIconSubject", Method: "bonusIconCards", Params: []string{"*EffectContext"}},
	}
}

// newFamily erases a family's typed catalog into the Family shape, pairing each
// row with the text its node renders. text is the family's own rendering call,
// since every family prints itself differently — an Effect renders Text(), a
// Refinement a clause on a phrase.
func newFamily[T any](
	name, method string,
	params []string,
	rows []Catalogued[T],
	text func(T) string,
) Family {
	out := make([]FamilyRow, len(rows))
	// Indexed rather than ranged over the values: ruleguard's largeloopcopy rule
	// asks go/types for the size of a ranged value, which panics on a type that
	// still carries a type parameter, so `for _, row := range rows` crashes the
	// linter here.
	for i := range rows {
		node := rows[i].Node
		out[i] = FamilyRow{
			Type:  nodeTypeName(node),
			Rules: rows[i].Rules,
			Text:  func() string { return text(node) },
		}
	}
	return Family{Name: name, Method: method, Params: params, Rows: out}
}

// gated switches the family's totality test on, which a family does once its
// catalog covers every node the source declares.
func (f Family) gated() Family {
	f.Gated = true
	return f
}

// nodeTypeName returns the name of a node's concrete type, unwrapping a pointer
// so a pointer-receiver node is named the way the source declares it.
func nodeTypeName(node any) string {
	return reflect.Indirect(reflect.ValueOf(node)).Type().Name()
}

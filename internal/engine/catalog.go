package engine

import (
	"fmt"
	"reflect"

	"github.com/dmikalova/vex/internal/census"
)

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
//
// The census has three shapes, one per shape its members take. A node family is
// an interface, discovered by the method its members declare (this file, and a
// catalog_<family>.go per family). A text-bearing value enum is a set of
// constants, discovered by the type they are declared with (catalog_enum.go) —
// which closes the same hole one level down, where an enumerating function like
// Keywords() was only as complete as its author remembered. Target is neither: it
// is a flag struct, so it is covered as both, its kinds as an enum and its filter
// builders as a family discovered by shape (catalog_target.go).

// Catalogued is one census row: a constructed node of the family, and the
// classification saying whether it owes a rulebook term. The row carries a real
// value rather than a type name so a consumer — a totality test, a Visitor, the
// style gallery — can render and walk the node it names.
type Catalogued[T any] struct {
	// Node is a representative value of the member, filled only as far as it needs
	// to render and resolve sensibly: a constructed node for a node family, the
	// constant itself for an enum, the base target with one filter applied for a
	// Target builder.
	Node T
	// Name is the member's own name in source, for a family whose members are not
	// types — an enum constant, a Target filter builder — because reflection can
	// recover a value's type but never the identifier that names it. A node family
	// leaves it empty and is keyed by its node's concrete type.
	Name string
	// Rules classifies the node: the term it owes, or the reason it owes none.
	Rules RulesBearing
	// Silent says the node renders no text by design, so the census does not read
	// its empty rendering as an underfilled literal. It is rare and deliberate: an
	// always-met condition adds no "while …" clause, and a bare quantity prints
	// nothing before the noun it counts.
	Silent bool
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
	// Also are the family's remaining interface methods, for a family one method
	// does not identify. A Count is a Value and a CountText, and CardsDiscarded is
	// a condition that exposes a Value for CountIs to read without being a Count,
	// so naming Value alone would report it as an uncatalogued member.
	Also []MethodSpec
	// Returns, set with Method left empty, discovers the family by shape rather
	// than by a method name: every exported method on the Name receiver whose only
	// result is this type. Target is the one family shaped that way — it is a flag
	// struct rather than an interface, because ADR 0005 keeps it comparable, so its
	// members are its filter builders.
	Returns string
	// Rows is the family's census, in catalog order.
	Rows []FamilyRow
	// Gated reports whether the family's totality test is switched on. A family
	// part-way through being catalogued is reported by `mage tool:census` but
	// does not yet fail the build; the task that completes its catalog gates it.
	Gated bool
}

// MethodSpec names one method a family's members declare: the method's name and
// its parameter types as written in source. It is the unit the source scan
// matches against, and a family whose interface has more than one method lists
// its remaining ones in Also.
type MethodSpec struct {
	Method string
	Params []string
}

// Declared returns every type in dir's non-test source that declares all of the
// family's methods, mapped to the file declaring the primary one — the members
// the family's catalog must cover. Both the totality test and `mage tool:census`
// read the source through here, so the two always see the same family.
func (f Family) Declared(dir string) (map[string]string, error) {
	if f.Method == "" {
		found, err := census.Builders(dir, f.Name, f.Returns)
		if err != nil {
			return nil, fmt.Errorf("scanning for %s builders: %w", f.Name, err)
		}
		return found, nil
	}
	var found map[string]string
	for _, spec := range append([]MethodSpec{{Method: f.Method, Params: f.Params}}, f.Also...) {
		declared, err := census.Implementations(dir, spec.Method, census.Params(spec.Params...))
		if err != nil {
			return nil, fmt.Errorf("scanning for %s.%s: %w", f.Name, spec.Method, err)
		}
		if found == nil {
			found = declared
			continue
		}
		for name := range found {
			if _, ok := declared[name]; !ok {
				delete(found, name)
			}
		}
	}
	return found, nil
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
	// Silent carries the row's declaration that its node prints nothing by
	// design, which is what exempts it from that check.
	Silent bool
}

// Families returns every catalogued node family. Each catalog_<family>.go
// contributes one entry, so a consumer that walks the whole census — the
// totality tests, `mage tool:census` — picks up a new family without being
// edited.
func Families() []Family {
	return []Family{
		effectFamily(),
		refinementFamily(),
		conditionFamily(),
		countFamily(),
		bonusIconSubjectFamily(),
		creatureVerbFamily(),
		selectionFamily(),
		spreadFamily(),
		topActFamily(),
		perTargetFamily(),
		lossFamily(),
		repeatGateFamily(),
		gatherFamily(),
		quantityFamily(),
		targetFilterFamily(),
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
			Type:   rowName(rows[i].Name, node),
			Rules:  rows[i].Rules,
			Text:   func() string { return text(node) },
			Silent: rows[i].Silent,
		}
	}
	return Family{Name: name, Method: method, Params: params, Rows: out}
}

// also records the family's remaining interface methods, narrowing the source
// scan to the types that declare all of them.
func (f Family) also(specs ...MethodSpec) Family {
	f.Also = specs
	return f
}

// builds discovers the family by shape instead of by a method name — the
// exported methods on its own receiver that return result — for a family that is
// a struct rather than an interface.
func (f Family) builds(result string) Family {
	f.Method, f.Params, f.Returns = "", nil, result
	return f
}

// gated switches the family's totality test on, which a family does once its
// catalog covers every node the source declares.
func (f Family) gated() Family {
	f.Gated = true
	return f
}

// rowName is the name the census keys a row by: the one the row states outright,
// or — for a node family, whose members are types — the node's concrete type,
// unwrapping a pointer so a pointer-receiver node is named the way the source
// declares it.
func rowName(stated string, node any) string {
	if stated != "" {
		return stated
	}
	return reflect.Indirect(reflect.ValueOf(node)).Type().Name()
}

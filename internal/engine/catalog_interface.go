package engine

import (
	"fmt"

	"github.com/dmikalova/vex/internal/census"
)

// This file is the census one level above the node families: every interface the
// package declares, classified as either a catalogued family or not a family at
// all.
//
// The per-node totality tests close the hole inside a family — a node with no
// census row fails the build. Nothing closed the hole around one. A family is an
// interface, and adding a strategy axis is how this codebase grows: Quantity,
// Gather and TopAct were each a new interface whose implementations would have
// sat outside every catalog with the build green. TestInterfaceTotality reads the
// interfaces out of the package's own source and fails on any that is neither
// catalogued nor given a reason it is not a family.
//
// The reasons are written as reasons rather than as category labels on purpose.
// A later reader has to be able to judge two things from the row alone: whether
// the interface was classified wrongly when it was written, and whether it has
// since grown members that print card text — at which point it has become a
// family and owes a catalog.

// logEntryCatalog is the name of the catalog covering LogEntry's variants. It is
// not one of Families(): log entries were enumerated by ADR 0046 before the node
// census existed, and LogEntrySamples remains their catalog, so the
// classification points at it rather than recording LogEntry as not a family.
const logEntryCatalog = "LogEntrySamples"

// CataloguedInterface is one interface's census row: the interface's name in
// source, and the classification saying which catalog covers its implementations
// or why it has none to cover.
type CataloguedInterface struct {
	// Name is the interface's own name as the source declares it, e.g. "Spread".
	Name string
	// Role is the row's classification.
	Role InterfaceRole
}

// InterfaceRole is an interface's classification, and the record of the
// judgement that produced it. Exactly one of its two columns is set; both empty
// and both set are equally invalid, and TestInterfaceRowsClassified fails on
// either.
//
// An interface is a node family when a card definition's meaning varies along it
// and its implementations print their own fragment of card text — Effect and the
// strategies beside it. Such a row sets Catalog to the catalog enumerating those
// implementations, and that catalog's own totality test then holds them.
//
// An interface is not a family when nothing a player reads comes out of it: a
// port, which is how the engine talks to its host (ADR 0008) and whose
// implementations are drivers rather than nodes; an optional capability, which
// only reshapes the text of a node already catalogued elsewhere; or an internal
// shape two value types share. Such a row sets NotFamily with the reason.
//
// When in doubt, classify it as a family and write the catalog: a thin catalog
// costs a file, a wrong NotFamily reopens the hole this classification exists to
// close.
type InterfaceRole struct {
	// Catalog is the name of the catalog enumerating this interface's
	// implementations — a Families() entry, or LogEntrySamples.
	Catalog string
	// NotFamily is the one-line reason this interface has no catalog, e.g.
	// "port: the whole surface an effect changes the game through". It is required
	// rather than optional so the table reads as a decision record.
	NotFamily string
}

// inCatalog classifies an interface as a node family, enumerated by the named
// catalog.
func inCatalog(catalog string) InterfaceRole { return InterfaceRole{Catalog: catalog} }

// notAFamily classifies an interface as having no catalog, recording why.
func notAFamily(reason string) InterfaceRole { return InterfaceRole{NotFamily: reason} }

// DeclaredInterfaces returns every interface declared in dir's non-test source,
// mapped to the file declaring it — the interfaces the classification must
// cover. The totality test and `mage tool:census` read the source through here,
// so the two always see the same set.
func DeclaredInterfaces(dir string) (map[string]string, error) {
	found, err := census.Interfaces(dir)
	if err != nil {
		return nil, fmt.Errorf("scanning for interfaces: %w", err)
	}
	return found, nil
}

// Interfaces returns the classification of every interface the package declares,
// in three groups: the node families, the ports, and the capability and shape
// interfaces that belong to neither.
func Interfaces() []CataloguedInterface {
	out := familyInterfaces()
	out = append(out, portInterfaces()...)
	return append(out, capabilityInterfaces()...)
}

// familyInterfaces are the interfaces a card definition's meaning varies along,
// each naming the catalog that enumerates its implementations.
func familyInterfaces() []CataloguedInterface {
	return []CataloguedInterface{
		{Name: "Effect", Role: inCatalog("Effect")},
		{Name: "Condition", Role: inCatalog("Condition")},
		{Name: "Count", Role: inCatalog("Count")},
		{Name: "Refinement", Role: inCatalog("Refinement")},
		{Name: "CreatureVerb", Role: inCatalog("CreatureVerb")},
		{Name: "Selection", Role: inCatalog("Selection")},
		{Name: "Spread", Role: inCatalog("Spread")},
		{Name: "TopAct", Role: inCatalog("TopAct")},
		{Name: "PerTarget", Role: inCatalog("PerTarget")},
		{Name: "Loss", Role: inCatalog("Loss")},
		{Name: "RepeatGate", Role: inCatalog("RepeatGate")},
		{Name: "Gather", Role: inCatalog("Gather")},
		{Name: "Quantity", Role: inCatalog("Quantity")},
		{Name: "BonusIconSubject", Role: inCatalog("BonusIconSubject")},
		{Name: "LogEntry", Role: inCatalog(logEntryCatalog)},
	}
}

// portInterfaces are the interfaces the engine talks to the outside through
// (ADR 0008). Every one of them is implemented by a host — *Game, a client, the
// bot, the simulator, a test double — and never by a node in a card's effect
// tree, so none of them puts a word on a card.
func portInterfaces() []CataloguedInterface {
	return []CataloguedInterface{
		{Name: "Resolver", Role: notAFamily(
			"port: the whole surface an effect changes the game through; *Game " +
				"implements it, no node does")},
		{Name: "StateReader", Role: notAFamily(
			"port: the read-only half of Resolver, composed from the readers below")},
		{Name: "EconomyReader", Role: notAFamily(
			"port: the reads of Æmber pools and forged keys, split out of " +
				"EconomyResolver")},
		{Name: "CreatureReader", Role: notAFamily(
			"port: the reads of the state carried on one card in play, split out of " +
				"CreatureResolver")},
		{Name: "ZoneReader", Role: notAFamily(
			"port: the reads of a player's zone contents, split out of ZoneResolver")},
		{Name: "TurnReader", Role: notAFamily(
			"port: the reads of turn-scoped state, split out of TurnResolver")},
		{Name: "EconomyResolver", Role: notAFamily(
			"port: the Resolver role changing Æmber, keys and chains")},
		{Name: "CreatureResolver", Role: notAFamily(
			"port: the Resolver role changing the state carried on one card in play")},
		{Name: "CombatResolver", Role: notAFamily(
			"port: the Resolver role resolving damage, destruction, and the uses an " +
				"ability makes")},
		{Name: "ZoneResolver", Role: notAFamily(
			"port: the Resolver role moving cards between zones")},
		{Name: "TurnResolver", Role: notAFamily(
			"port: the Resolver role installing turn-scoped grants and the lasting " +
				"registry")},
		{Name: "ChoiceResolver", Role: notAFamily(
			"port: the Resolver role asking a player to decide, so an effect can " +
				"branch on the answer")},
		{Name: "Logger", Role: notAFamily(
			"port: the Resolver role an effect narrates through (ADR 0011); the " +
				"entries it carries are catalogued by " + logEntryCatalog)},
		{Name: "Chooser", Role: notAFamily(
			"port: the decision-maker a host installs — a client, the bot, the " +
				"simulator — which answers an effect's question and prints nothing")},
		{Name: "OptionChooser", Role: notAFamily(
			"port: the optional Chooser capability picking one of several labeled " +
				"options")},
		{Name: "DeclinableChooser", Role: notAFamily(
			"port: the optional Chooser capability offering a card choice the player " +
				"may pass on")},
		{Name: "PositionChooser", Role: notAFamily(
			"port: the optional Chooser capability pointing at a battleline position " +
				"for Deploy")},
		{Name: "Orderer", Role: notAFamily(
			"port: the optional Chooser capability arranging ids into a resolution " +
				"order in one call")},
		{Name: "ReactionChooser", Role: notAFamily(
			"port: the optional Chooser capability picking which reaction in a " +
				"trigger window resolves next")},
		{Name: "BadgeChooser", Role: notAFamily(
			"port: the display-only Chooser capability previewing the status a " +
				"creature about to be chosen will receive")},
		{Name: "ActionChooser", Role: notAFamily(
			"port: the optional Chooser capability choosing the next root action " +
				"(ADR 0039), installed only by an interactive driver")},
		{Name: "FirstPlayerChooser", Role: notAFamily(
			"port: the optional Chooser capability naming which player goes first, " +
				"so a Chooser that omits it defaults to player 0")},
		{Name: "Namer", Role: notAFamily(
			"port: how a log entry resolves the ids it holds to the names its reader " +
				"sees; a client implements it")},
		{Name: "sourced", Role: notAFamily(
			"port: a Namer that also knows the source card of the frame being " +
				"rendered, so an outcome entry can name its subject")},
	}
}

// capabilityInterfaces are the interfaces that are neither a family nor a port:
// the optional capabilities a node declares to shape how it reads or how it is
// offered, and the two internal shapes a pair of value types share. A capability
// has no members of its own to catalogue — the node implementing it is already a
// row in its own family's catalog — and it adds no vocabulary, only a different
// arrangement of words the node already prints.
//
// A capability that started printing a word of its own would have become a
// family, and its row here would be wrong. That is the judgement a reader is
// meant to be able to re-make from these reasons.
func capabilityInterfaces() []CataloguedInterface {
	out := effectCapabilityInterfaces()
	out = append(out, conditionCapabilityInterfaces()...)
	out = append(out, countCapabilityInterfaces()...)
	out = append(out, targetCapabilityInterfaces()...)
	return append(out, sharedShapeInterfaces()...)
}

// effectCapabilityInterfaces are the capabilities an Effect declares: how it
// folds into a Sequence's prose, and what it can tell May before being offered.
// Every implementation is already a row in the Effect catalog.
func effectCapabilityInterfaces() []CataloguedInterface {
	return []CataloguedInterface{
		{Name: "GatingEffect", Role: notAFamily(
			"Effect capability: reports whether it did anything, so a Then can gate " +
				"on it; its implementations are catalogued as Effects")},
		{Name: "declinableEffect", Role: notAFamily(
			"Effect capability: offers May its own single-target choice instead of a " +
				"Yes/No question; its implementations are catalogued as Effects")},
		{Name: "vacuousEffect", Role: notAFamily(
			"Effect capability: reports that it would do nothing, so May skips a " +
				"question with one honest answer; catalogued as Effects")},
		{Name: "validator", Role: notAFamily(
			"Effect capability: checks its own configuration when NewCard builds the " +
				"card; it prints nothing and is catalogued as Effects")},
		{Name: "combinable", Role: notAFamily(
			"Effect text capability: folds its text with its neighbours in a " +
				"Sequence; the fold rearranges words the Effect rows already print")},
		{Name: "foldable", Role: notAFamily(
			"Effect text capability refining combinable: says when the fold applies; " +
				"its implementations are catalogued as Effects")},
		{Name: "sentenceEnder", Role: notAFamily(
			"Effect text capability: declares its text ends a sentence, so a Sequence " +
				"opens a new one rather than conjoining")},
		{Name: "nounListable", Role: notAFamily(
			"Effect text capability: exposes the varying noun in a fixed frame so a " +
				"Sequence states the frame once")},
		{Name: "hedgedLeadIn", Role: notAFamily(
			"Effect text capability: re-renders its lead-in under a May so the " +
				"consequence does not read as mandatory")},
		{Name: "ladderRepeating", Role: notAFamily(
			"Effect text capability: renders as a repetition of an identical effect " +
				"already stated ('gain 1 more')")},
		{Name: "durationScoped", Role: notAFamily(
			"Effect text capability: splits its text into subject and predicate so a " +
				"duration wrapper states the window once")},
	}
}

// conditionCapabilityInterfaces are the capabilities a Condition declares, each
// a different grammatical form of the clause its Condition row already prints.
func conditionCapabilityInterfaces() []CataloguedInterface {
	return []CataloguedInterface{
		{Name: "negatable", Role: notAFamily(
			"Condition text capability: owns its own negated wording for Not, which " +
				"is not uniform enough to derive; catalogued as Conditions")},
		{Name: "itShaped", Role: notAFamily(
			"Condition text capability: renders as an adjective on a shared noun " +
				"rather than as a whole clause")},
		{Name: "symmetricCondTexter", Role: notAFamily(
			"Condition text capability: supplies the board-wide third-person wording " +
				"when trimming 'if ' from CondText will not do")},
		{Name: "ladderRung", Role: notAFamily(
			"Condition text capability: exposes the counted subject a threshold " +
				"ladder states once instead of in every rung")},
	}
}

// countCapabilityInterfaces are the capabilities a Count or a Loss declares: the
// alternate phrasings of a quantity its own family's row already renders.
func countCapabilityInterfaces() []CataloguedInterface {
	return []CataloguedInterface{
		{Name: "countClauser", Role: notAFamily(
			"Count text capability: supplies the verb phrase an 'if' clause needs, " +
				"where CountText's bare noun reads only after 'for each'")},
		{Name: "cardinalCounter", Role: notAFamily(
			"Count text capability: renders as a cardinal 'the number of …' phrase " +
				"for a clause comparing against the count")},
		{Name: "leadingCounter", Role: notAFamily(
			"Count text capability: names its subject when its clause leads the " +
				"sentence and a trailing 'it' would read as a forward reference")},
		{Name: "eachPlayerEqualTo", Role: notAFamily(
			"Count text capability: supplies the third-person 'equal to …' form " +
				"EachPlayer re-bases onto whoever is being paid")},
		{Name: "portionPhraser", Role: notAFamily(
			"Loss text capability: phrases a fraction of a count as well as of a " +
				"pool; Fraction, its implementation, is catalogued as a Loss")},
	}
}

// targetCapabilityInterfaces are the capabilities a Refinement, a Selection, or a
// CreatureVerb declares — how the card it is aimed at is worded, and which
// candidates are offered for it.
func targetCapabilityInterfaces() []CataloguedInterface {
	return []CataloguedInterface{
		{Name: "framedClause", Role: notAFamily(
			"Refinement text capability: exposes the one varying word in a fixed " +
				"frame so a combinator states the frame once")},
		{Name: "leadingRefinement", Role: notAFamily(
			"Refinement text capability: renders its choice before the effect's verb " +
				"so the phrase runs left to right")},
		{Name: "membershipRefiner", Role: notAFamily(
			"Refinement capability: answers whether a card is among its result " +
				"without making the tie-break its refine would prompt for")},
		{Name: "declinableSelection", Role: notAFamily(
			"Selection capability: says the controller may pass, which is what lets " +
				"its verb read 'you may'; catalogued as Selections")},
		{Name: "positionalSelection", Role: notAFamily(
			"Selection capability: picks by an end of an ordered zone, which the node " +
				"validates against the zone it was paired with")},
		{Name: "qualifiableSelection", Role: notAFamily(
			"Selection text capability: takes an adjective between its determiner and " +
				"its noun when the scope is not already in the noun")},
		{Name: "ownerActsSelection", Role: notAFamily(
			"Selection text capability: says the pick is made blind, which sets whose " +
				"voice the verb is written in")},
		{Name: "narrowingVerb", Role: notAFamily(
			"CreatureVerb capability: narrows the candidates offered without changing " +
				"the printed text; catalogued as CreatureVerbs")},
	}
}

// sharedShapeInterfaces are the two interfaces that are neither a family, a
// port, nor a capability: an internal shape a handful of concrete types share so
// one piece of engine code can reach all of them.
func sharedShapeInterfaces() []CataloguedInterface {
	return []CataloguedInterface{
		{Name: "cardPile", Role: notAFamily(
			"internal shape: the two operations the five resting zones' card lists " +
				"share, so the mover reaches a deckList and a wideList alike")},
		{Name: "selfHouseResolvable", Role: notAFamily(
			"internal shape: how a value type with unexported fields resolves its own " +
				"house sentinels, which reflection can read but never write")},
	}
}

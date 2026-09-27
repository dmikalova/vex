package engine

import (
	"fmt"

	"github.com/dmikalova/vex/internal/census"
)

// This file is the census of the package's text-bearing value enums: the small
// closed types a card writes that also print words a player reads — a Duration
// ("for the remainder of the turn"), a Destination ("into their hand"), a counter
// kind ("doom counter").
//
// It closes the hole one level below the node families. An enum is enumerated by
// a hand-written function — Keywords(), Durations() — and until now nothing
// checked that function against the constants the source declares, so a new
// constant nobody added to its function was invisible to the rulebook
// completeness test and ADR 0018's "complete by construction" rested on an author
// remembering. TestEnumTotality reads the constants out of the package's own
// source and fails the build when one is neither enumerated nor named as
// excluded.
//
// Two shapes of entry. An enum a card's printed text renders carries a full
// catalog, one classified row per value, the same mandatory classification the
// node families use. The five enums that already have their own rulebook
// completeness test — Keyword, CardType, BonusIcon, Trigger, Phase — carry only
// the guard, because classifying them here would duplicate that test rather than
// close anything.
//
// Terms fold hard: a Duration owes one "Duration" term, not one per window, and a
// Destination binds to the movement and zone terms the effect catalog already
// names.

// Enum is one text-bearing value enum's census: the constants the source declares
// with its type, split into the ones its enumerating function returns and the
// ones it deliberately leaves out.
type Enum struct {
	// Type is the enum's type name as the source declares it, e.g. "Duration".
	// Destination is a struct over a zone enum, so its type here is the zone half
	// the constants are declared with.
	Type string
	// Excluded names the constants of the type the enumerating function leaves out
	// — an invalid zero value, a count sentinel bounding the enum, a value no card
	// can name. Naming each one is what lets the scan insist on every other
	// constant, so a new member cannot hide among them.
	Excluded []string
	// Enumerated is how many values the enum's enumerating function returns, which
	// must be every declared constant that is not excluded.
	Enumerated int
	// Rows is one classified row per enumerated value, for the enums a card's
	// printed text renders. It is empty for an enum carrying only the guard.
	Rows []FamilyRow
}

// Declared returns every constant in dir's non-test source declared with the
// enum's type, mapped to the file declaring it.
func (e Enum) Declared(dir string) (map[string]string, error) {
	found, err := census.Constants(dir, e.Type)
	if err != nil {
		return nil, fmt.Errorf("scanning for %s constants: %w", e.Type, err)
	}
	return found, nil
}

// Enums returns every text-bearing value enum the census covers.
func Enums() []Enum {
	return []Enum{
		durationEnum(),
		destinationEnum(),
		deckDestEnum(),
		neighborScopeEnum(),
		damageAftermathEnum(),
		comparisonEnum(),
		counterKindEnum(),
		targetKindEnum(),
		useKindEnum(),
		tollActionEnum(),

		// The five that already carry their own rulebook completeness test, here
		// only for the guard that their enumerating function is complete.
		guarded("Keyword", len(Keywords()), "keywordUnset", "keywordCount"),
		guarded("CardType", len(allCardTypes()), "TypeUnset", "AnyType"),
		guarded("BonusIcon", len(allBonusIcons()), "bonusUnset", "bonusIconCount"),
		guarded("Trigger", len(Triggers()), "triggerUnset", "triggerCount"),
		guarded("Phase", len(Phases()), "phaseUnset"),
	}
}

// newEnum builds an enum's census entry, pairing each row with the text its value
// renders. text is the enum's own rendering call, since every enum prints itself
// differently — a Duration names its window, a Destination renders a whole move.
func newEnum[T any](
	typeName string,
	excluded []string,
	enumerated int,
	rows []Catalogued[T],
	text func(T) string,
) Enum {
	out := make([]FamilyRow, len(rows))
	// Indexed rather than ranged over the values, for the same reason newFamily is.
	for i := range rows {
		value := rows[i].Node
		out[i] = FamilyRow{
			Type:   rows[i].Name,
			Rules:  rows[i].Rules,
			Text:   func() string { return text(value) },
			Silent: rows[i].Silent,
		}
	}
	return Enum{Type: typeName, Excluded: excluded, Enumerated: enumerated, Rows: out}
}

// guarded builds the entry for an enum the census only guards: its enumerating
// function is checked against the source, and its values are classified by the
// rulebook completeness test that already covers them.
func guarded(typeName string, enumerated int, excluded ...string) Enum {
	return Enum{Type: typeName, Excluded: excluded, Enumerated: enumerated}
}

// durationEnum is the Duration enum's census. Every window owes the one
// "Duration" term: how long an effect lasts is one rule, and a term per window
// would restate it six times.
func durationEnum() Enum {
	return newEnum(
		"Duration",
		[]string{"durationUnset", "durationCount"},
		len(Durations()),
		[]Catalogued[Duration]{
			{Name: "RemainderOfPlayerTurn", Node: RemainderOfPlayerTurn, Rules: bears("Duration")},
			{Name: "OpponentNextTurn", Node: OpponentNextTurn, Rules: bears("Duration")},
			{Name: "StartOfPlayerNextTurn", Node: StartOfPlayerNextTurn, Rules: bears("Duration")},
			{Name: "EndOfPlayerNextTurn", Node: EndOfPlayerNextTurn, Rules: bears("Duration")},
			{Name: "UntilThisLeavesPlay", Node: UntilThisLeavesPlay, Rules: bears("Duration")},
			{Name: "UntilCardLeavesPlay", Node: UntilCardLeavesPlay, Rules: bears("Duration")},
		},
		Duration.String,
	)
}

// deckDestEnum is the DeckDest enum's census: where a look at the top of a deck
// routes the cards it took. Each routing is the deck-side face of an ordinary
// zone movement, so each binds to that movement's term.
func deckDestEnum() Enum {
	return newEnum(
		"DeckDest",
		nil,
		len(DeckDests()),
		[]Catalogued[DeckDest]{
			{Name: "IntoHand", Node: IntoHand, Rules: bears("Put a Card into Another Zone")},
			{Name: "IntoArchives", Node: IntoArchives, Rules: bears("Archive")},
			{Name: "IntoDiscard", Node: IntoDiscard, Rules: bears("Discard")},
			{Name: "IntoPurge", Node: IntoPurge, Rules: bears("Purge")},
			{
				Name:  "IntoBottomOfDeck",
				Node:  IntoBottomOfDeck,
				Rules: bears("Put a Card into Another Zone"),
			},
		},
		func(d DeckDest) string { return d.act(destinationSubject) },
	)
}

// neighborScopeEnum is the NeighborScope enum's census: how many of a struck
// creature's neighbors a spread hits, which is the splash rule.
func neighborScopeEnum() Enum {
	return newEnum(
		"NeighborScope",
		nil,
		len(NeighborScopes()),
		[]Catalogued[NeighborScope]{
			{Name: "AllNeighbors", Node: AllNeighbors, Rules: bears("Splash-attack")},
			{Name: "OneNeighbor", Node: OneNeighbor, Rules: bears("Splash-attack")},
		},
		func(s NeighborScope) string {
			return CreatureAndNeighbors{Amount: 2, Splash: 1, Scope: s}.spreadText()
		},
	)
}

// damageAftermathEnum is the DamageAftermath enum's census: when a DealDamage's
// follow-up resolves. The two branches turn on whether the damage destroyed the
// creature, which the damage rules teach; Always is a plain join.
func damageAftermathEnum() Enum {
	return newEnum(
		"DamageAftermath",
		[]string{"aftermathUnset"},
		len(DamageAftermaths()),
		[]Catalogued[DamageAftermath]{
			{
				Name:  "Always",
				Node:  Always,
				Rules: plumbing("composition: runs the follow-up whatever the damage did"),
			},
			{Name: "IfDestroyed", Node: IfDestroyed, Rules: bears("Deal Damage")},
			{Name: "IfSurvives", Node: IfSurvives, Rules: bears("Deal Damage")},
		},
		func(a DamageAftermath) string {
			return DealDamage{
				Target: Target{Kind: TargetChosenCreature},
				Amount: 2,
				Then:   GainAember{Player: Controller, Amount: 1},
				After:  a,
			}.Text()
		},
	)
}

// comparisonEnum is the Comparison enum's census: how a condition compares a
// quantity to a threshold. Every comparison is part of the one conditional rule,
// so all seven name the "Conditional" umbrella.
func comparisonEnum() Enum {
	return newEnum(
		"Comparison",
		[]string{"comparisonUnset"},
		len(Comparisons()),
		[]Catalogued[Comparison]{
			{Name: "AtLeast", Node: AtLeast, Rules: bears("Conditional")},
			{Name: "AtMost", Node: AtMost, Rules: bears("Conditional")},
			{Name: "Exactly", Node: Exactly, Rules: bears("Conditional")},
			{Name: "MoreThanYou", Node: MoreThanYou, Rules: bears("Conditional")},
			{Name: "MoreThanOpponent", Node: MoreThanOpponent, Rules: bears("Conditional")},
			{Name: "Even", Node: Even, Rules: bears("Conditional")},
			{Name: "Odd", Node: Odd, Rules: bears("Conditional")},
		},
		comparisonClause,
	)
}

// comparisonClause renders a comparison the way a card prints it, through the
// Æmber-pool condition that is the only one reading all seven. The relative
// comparisons are each tied to one side, so the subject follows the comparison.
func comparisonClause(c Comparison) string {
	player := Opponent
	if c == MoreThanOpponent {
		player = Controller
	}
	return PoolAember{Player: player, Is: c, Amount: 2}.CondText()
}

// counterKindEnum is the CounterKind enum's census. A counter kind carries no
// rule of its own beyond the generic-counter rule (ADR 0024) — it is a marker a
// card names — so every kind binds to that one term.
func counterKindEnum() Enum {
	return newEnum(
		"CounterKind",
		[]string{"CounterNone"},
		len(CounterKinds()),
		[]Catalogued[CounterKind]{
			{Name: "CounterDoom", Node: CounterDoom, Rules: bears("Generic Counters")},
			{Name: "CounterFuse", Node: CounterFuse, Rules: bears("Generic Counters")},
			{Name: "CounterGrowth", Node: CounterGrowth, Rules: bears("Generic Counters")},
			{Name: "CounterGlory", Node: CounterGlory, Rules: bears("Generic Counters")},
			{
				Name:  "CounterDisruption",
				Node:  CounterDisruption,
				Rules: bears("Generic Counters"),
			},
			{Name: "CounterScheme", Node: CounterScheme, Rules: bears("Generic Counters")},
			{Name: "CounterWarrant", Node: CounterWarrant, Rules: bears("Generic Counters")},
		},
		CounterKind.noun,
	)
}

// useKindEnum is the UseKind enum's census: the three ways a card in play can be
// used, and so the three verbs a "cannot" restriction bars. Each names the one
// "Cannot Be Used To" term — barring a use is one rule, and the kind only says
// which use it bars.
func useKindEnum() Enum {
	return newEnum(
		"UseKind",
		[]string{"useKindUnset", "useKindCount"},
		len(allUseKinds()),
		[]Catalogued[UseKind]{
			{Name: "ReapUse", Node: ReapUse, Rules: bears("Cannot Be Used To")},
			{Name: "FightUse", Node: FightUse, Rules: bears("Cannot Be Used To")},
			{Name: "ActionUse", Node: ActionUse, Rules: bears("Cannot Be Used To")},
		},
		func(k UseKind) string { return "This creature cannot " + k.verb() + "." },
	)
}

// tollActionEnum is the TollAction enum's census: the two artifact actions a card
// in play can charge the opponent for. Both name the one "Toll" term, since the
// action only says which move the Æmber is owed on.
func tollActionEnum() Enum {
	return newEnum(
		"TollAction",
		[]string{"tollActionUnset"},
		len(allTollActions()),
		[]Catalogued[TollAction]{
			{Name: "TollPlayArtifact", Node: TollPlayArtifact, Rules: bears("Toll")},
			{Name: "TollUseArtifact", Node: TollUseArtifact, Rules: bears("Toll")},
		},
		func(a TollAction) string {
			return restrictionText(Restrictions{Toll: Toll{Action: a, Amount: 1}}, false)[0]
		},
	)
}

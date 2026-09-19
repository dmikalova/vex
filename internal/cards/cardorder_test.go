package cards

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"testing"
)

// optionRank is the canonical authoring order of the option arguments to
// card.New — the tail after the four positional arguments (name, house, type,
// rarity). A card lists its options in non-decreasing rank; options sharing a
// rank may appear in any order relative to each other, and any option may appear
// zero or many times. The order matches the card-authoring guide
// (internal/cards/AGENTS.md): origin tag, creature stats, descriptors, then the
// remaining in-play modifiers, abilities, and deck-generation sidecar options.
var optionRank = map[string]int{
	// Origin tag first.
	"Provenance":     1,
	"InSet":          2,
	"RarityWeight":   3,
	"InCluster":      4,
	"LeadsCluster":   5,
	"PullsMatching":  6,
	"Houseless":      6,
	"OneCopyPerDeck": 7,
	// Creature stats.
	"WithBonus":          8,
	"WithEnhance":        9,
	"WithoutEnhancement": 10,
	"WithPower":          11,
	"WithArmor":          12,
	// Descriptors.
	"WithTraits":       13,
	"WithKeywords":     14,
	"WithAssault":      15,
	"WithHazardous":    16,
	"WithSplashAttack": 17,
	// Top priority abilities
	"WithEntersPlay": 18,
	// Miscellaneous in-play modifiers.
	"WithAemberCannotBeStolen":              19,
	"WithAemberCost":                        19,
	"WithAemberThreshold":                   19,
	"WithAlsoTakesNeighborFightDamage":      19,
	"WithAttackDamage":                      19,
	"WithAttackIgnores":                     19,
	"WithAttackKeywords":                    19,
	"WithBonusInstead":                      19,
	"WithCannotBeDealtDamageBy":             19,
	"WithCannotBeUsedTo":                    19,
	"WithCannotBeUsedWhile":                 19,
	"WithCannotPlayWhile":                   19,
	"WithConstant":                          19,
	"WithDestroyedWhen":                     19,
	"WithDrawModifier":                      19,
	"WithDrawModifierInCenter":              19,
	"WithDrawModifierOffFlank":              19,
	"WithDrawModifierPer":                   19,
	"WithFightRestriction":                  19,
	"WithFriendlyEntersPlayReady":           19,
	"WithGainsForgeAember":                  19,
	"WithHouseLock":                         19,
	"WithKeyCost":                           19,
	"WithNoDamageWhenAttacked":              19,
	"WithPlayableAsUpgrade":                 19,
	"WithPlayPermission":                    19,
	"WithPowerX":                            19,
	"WithReplaces":                          19,
	"WithRestrictions":                      19,
	"WithSpendableAember":                   19,
	"WithStatic":                            19,
	"WithStealsInsteadOfDamageWhenAttacked": 19,
	"WithTakesDamageFor":                    19,
	"WithTauntReachingNeighborsNeighbors":   19,
	"WithTriggersFromDiscard":               19,
	// Standard abilities	.
	"WithAbility":           20,
	"WithEachPlayerAbility": 20,
	// Generative template last: it replaces the face with a materialized variant,
	// so it wraps everything the face declares above it.
	"Template": 21,
}

// TestOptionsAreInCanonicalOrder enforces a single authoring order for the
// options passed to card.New across every card, so a reader always finds a card's
// stats, keywords, and abilities in the same place. It checks the option tail of
// each card.New call against optionRank; a call whose final argument is spread
// (a family wrapper forwarding `opts...`, e.g. Master of N) is skipped because
// its order cannot be read statically.
func TestOptionsAreInCanonicalOrder(t *testing.T) {
	fset := token.NewFileSet()
	for _, pkg := range setPackages(t) {
		dir := filepath.Join(setsDir, pkg.name)
		for _, file := range pkg.cardFiles {
			f := parseFile(t, fset, filepath.Join(dir, file))
			ast.Inspect(f, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok || !isCardNewCall(call) {
					return true
				}
				checkOptionOrder(t, pkg.name, file, call)
				return true
			})
		}
	}
}

// checkOptionOrder reports any card.New option argument that falls out of
// canonical order, or any option missing from optionRank.
func checkOptionOrder(t *testing.T, pkg, file string, call *ast.CallExpr) {
	t.Helper()
	if call.Ellipsis != token.NoPos || len(call.Args) <= 4 {
		return
	}
	prevRank, prevName := 0, ""
	for _, arg := range call.Args[4:] {
		name, ok := optionName(arg)
		if !ok {
			// A non-inline option (an identifier, a helper call): its order cannot
			// be read here, so stop checking this call rather than guess.
			return
		}
		rank, known := optionRank[name]
		if !known {
			t.Errorf(
				"%s/%s: card.%s has no entry in optionRank; add it in canonical order",
				pkg,
				file,
				name,
			)
			return
		}
		if rank < prevRank {
			t.Errorf(
				"%s/%s: card.%s (rank %d) appears after card.%s (rank %d); options must be in canonical order (see optionRank)",
				pkg,
				file,
				name,
				rank,
				prevName,
				prevRank,
			)
			return
		}
		prevRank, prevName = rank, name
	}
}

// optionName returns the option constructor name of a `card.X(...)` argument
// (e.g. "WithPower" for card.WithPower(3)), or ok=false if the argument is not a
// call on the card package.
func optionName(arg ast.Expr) (string, bool) {
	call, ok := arg.(*ast.CallExpr)
	if !ok {
		return "", false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "card" {
		return "", false
	}
	return sel.Sel.Name, true
}

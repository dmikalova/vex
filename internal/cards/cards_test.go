package cards

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/deckgen"
	"github.com/dmikalova/vactrol/internal/engine"
)

// TestAllIsAValidDatabase checks the assembled card database: every set package
// is imported, every card self-registers, and each entry is well-formed with a
// unique name, a real house, and a rarity.
func TestAllIsAValidDatabase(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("All() returned no cards; are the set packages imported?")
	}

	seen := make(map[string]bool, len(all))
	for _, c := range all {
		switch {
		case c.Name == "":
			t.Error("card with empty name")
		case c.House == engine.HouseNone:
			t.Errorf("%q has no house", c.Name)
		case c.Rarity == "":
			t.Errorf("%q has no rarity", c.Name)
		case c.Type == engine.TypeUnset:
			t.Errorf("%q has no card type", c.Name)
		}
		if seen[c.Name] {
			t.Errorf("duplicate card name %q", c.Name)
		}
		seen[c.Name] = true
	}
}

// TestNoDuplicateActionTrigger enforces that a card never carries more than one
// ability on the same action trigger (Play, Fight, or Reap). A card that shares
// one effect across several of these triggers must declare it once with a
// composite trigger (Trigger.PlayReap, Trigger.FightReap, Trigger.PlayFightReap),
// which fans out into the atomic triggers; declaring both an atomic ability and a
// composite that also covers it — e.g. Trigger.Play and Trigger.PlayReap — would
// silently give the card two Play abilities that both fire.
func TestNoDuplicateActionTrigger(t *testing.T) {
	actionTriggers := map[engine.Trigger]string{
		engine.TriggerAfterPlay:  "Play",
		engine.TriggerAfterFight: "Fight",
		engine.TriggerAfterReap:  "Reap",
	}
	for _, c := range All() {
		counts := make(map[engine.Trigger]int)
		for _, ab := range c.Abilities {
			counts[ab.Trigger]++
		}
		for trig, name := range actionTriggers {
			if counts[trig] > 1 {
				t.Errorf(
					"%s has %d %s abilities; consolidate them into one (composite trigger)",
					c.Name, counts[trig], name,
				)
			}
		}
	}
}

// TestEveryCreatureAndArtifactHasTrait enforces the card-database policy that
// every creature and artifact carries at least one trait (e.g. Giant, Beast,
// Weapon). Actions and upgrades are exempt.
func TestEveryCreatureAndArtifactHasTrait(t *testing.T) {
	for _, c := range All() {
		if c.Type != engine.Creature && c.Type != engine.Artifact {
			continue
		}
		if len(c.Traits) == 0 {
			t.Errorf(
				"%s (%s) has no trait; every creature and artifact needs at least one",
				c.Name,
				c.Type,
			)
		}
	}
}

// TestAllIsSortedDeterministically verifies the database comes back in a stable
// order (house, then name), independent of package initialization order.
func TestAllIsSortedDeterministically(t *testing.T) {
	all := All()
	for i := 1; i < len(all); i++ {
		prev, cur := all[i-1], all[i]
		if prev.House > cur.House || (prev.House == cur.House && prev.Name > cur.Name) {
			t.Errorf("not sorted at %d: %s (%s) before %s (%s)",
				i, prev.Name, prev.House, cur.Name, cur.House)
		}
	}
}

// TestMaterializedNamesAreUnique extends the duplicate-name check in
// TestAllIsAValidDatabase past materialization: a template never appears in a
// deck under its own registered name, only under whatever name its Materializer
// gives the concrete variant it produces, so those variant names must be just as
// unique as any registered card's name — both against every other registered
// name and against every other template's variants. It samples many (House, seed)
// combinations per template so a randomized Materializer exercises its whole
// output range instead of just whatever it happens to produce once.
func TestMaterializedNamesAreUnique(t *testing.T) {
	regs := card.Cards()

	// owner tracks who a name belongs to: a registered card by its own name, or a
	// template by the name it was found under.
	owner := make(map[string]string, len(regs))
	for _, rc := range regs {
		owner[rc.Def.Name] = rc.Def.Name
	}

	const samplesPerHouse = 20
	for _, rc := range regs {
		if rc.Materializer == nil {
			continue
		}
		for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
			for seed := int64(0); seed < samplesPerHouse; seed++ {
				ctx := deckgen.SlotContext{House: h, Rarity: rc.Def.Rarity}
				out := rc.Materializer.Materialize(ctx, rand.New(rand.NewSource(seed)))
				if out.Name == rc.Def.Name {
					// The template face's own name never reaches a deck (materialize
					// always replaces it), so it does not need to out-compete itself.
					continue
				}
				if by, ok := owner[out.Name]; ok && by != rc.Def.Name {
					t.Errorf(
						"%s materializes a variant named %q, which collides with %s",
						rc.Def.Name, out.Name, by,
					)
					continue
				}
				owner[out.Name] = rc.Def.Name
			}
		}
	}
}

// TODO: resolve these - also bring back hidden stash
// identicalTwinTODO lists cards deliberately implemented as byte-identical to an
// existing card because they print the same effect in a different house (a
// cross-house twin). A plain reprint would keep the original's house and folding
// would drop the twin's, so these are allowed to duplicate pending a manual fold.
// Each carries a TODO(duplicate) on its definition. See internal/cards/AGENTS.md.
var identicalTwinTODO = map[string]bool{
	"Subtle Chain": true, // ≡ Mind Barb (Dis)
	"Stealth Mode": true, // ≡ Scrambler Storm (Logos)
	// Worlds Collide "Brew" variants: differently named copies of Alaka's Brew
	// (same Brobnar effect), Variant rarity pending a real-rarity mapping.
	"Chieftain's Brew":  true, // ≡ Alaka's Brew (Brobnar)
	"Cowfyne's Brew":    true, // ≡ Alaka's Brew (Brobnar)
	"Groke's Brew":      true, // ≡ Alaka's Brew (Brobnar)
	"Gron's Brew":       true, // ≡ Alaka's Brew (Brobnar)
	"Mogghunter's Brew": true, // ≡ Alaka's Brew (Brobnar)
	"Narp's Brew":       true, // ≡ Alaka's Brew (Brobnar)
	"Shorty's Brew":     true, // ≡ Alaka's Brew (Brobnar)
}

// TestNoDuplicateImplementations flags two differently named cards whose whole
// behavior — type, stats, keywords, and abilities — is otherwise identical,
// ignoring rarity and the house/name that legitimately differ between a card and
// its reprint. Cards that share an implementation must share one definition with
// multiple Provenance tags, not two copy-pasted definitions, so any collision
// here is a duplicate to fold into a single card.
func TestNoDuplicateImplementations(t *testing.T) {
	all := All()
	seenBy := make(map[string]string, len(all))
	for _, c := range all {
		sig := c
		sig.Name = ""
		sig.House = engine.HouseNone
		sig.Rarity = ""
		key := fmt.Sprintf("%#v", sig)
		if owner, ok := seenBy[key]; ok {
			if identicalTwinTODO[c.Name] || identicalTwinTODO[owner] {
				continue // deliberate cross-house twin; TODO(duplicate) tracks the fold
			}
			t.Errorf(
				"%q and %q have identical implementations (same type, stats, keywords, and abilities); fold them into one card with multiple Provenance tags",
				owner,
				c.Name,
			)
			continue
		}
		seenBy[key] = c.Name
	}
}

// TestProvenanceHasNoOverlap checks that no two differently named cards claim the
// same original source printing (set + collector number). Provenance is
// repeatable on a single card (a reprint across sets shares one implementation,
// e.g. Fogbank's four printings), but two different vactrol cards pointing at
// the same original is an authoring mistake — most likely a Provenance tag
// copied from another card's file and never updated.
func TestProvenanceHasNoOverlap(t *testing.T) {
	seenBy := make(map[provenance.Ref]string)
	for _, rc := range card.Cards() {
		for _, ref := range rc.Provenance {
			if owner, ok := seenBy[ref]; ok && owner != rc.Def.Name {
				t.Errorf(
					"%s and %s both claim provenance %s #%s",
					owner, rc.Def.Name, ref.Set.Code, ref.Number,
				)
				continue
			}
			seenBy[ref] = rc.Def.Name
		}
	}
}

// TestEveryReprintResolvesToACard checks that every reprint claim a set makes (its
// 0set.go card.Reprint entries) names a card some set actually implements. A claim
// that resolves to nothing is a stale 0set.go — a card renamed or removed without
// regenerating — which would silently shrink that set's pool; here it fails the
// gate loudly and names the offending claim instead.
func TestEveryReprintResolvesToACard(t *testing.T) {
	byName := make(map[string]bool)
	for _, rc := range card.Cards() {
		byName[normalizeName(rc.Def.Name)] = true
	}
	for _, rp := range card.ReprintRefs() {
		if !byName[normalizeName(rp.Name)] {
			t.Errorf(
				"%s reprint #%s refers to %q, which no set implements "+
					"(stale 0set.go — regenerate with `mage tool:stub %s`)",
				rp.Set.Name, rp.Number, rp.Name, rp.Set.Slug,
			)
		}
	}
}

// TestNoCardHardcodesItsOwnHouse enforces the self-house convention (see
// internal/cards/AGENTS.md): a card whose ability names its own house writes
// card.House.Self, never that house spelled out. Self is resolved to the card's
// own house once, at card.New time, so the printed house and the ability can
// never drift and a Maverick inherits the house it is printed in. Only a card
// that names a *different* house — Brobnar Ambassador, a Sanctum card, naming
// Brobnar — spells the house out, and that never matches its own house here.
//
// The check reads source because card.New erases the distinction: after
// resolution both card.House.Self and the literal are the same concrete house,
// so a test over the registered database could not tell an author who wrote Self
// from one who hardcoded the house. Parsing the card.New call sites keeps the two
// apart.
func TestNoCardHardcodesItsOwnHouse(t *testing.T) {
	fset := token.NewFileSet()
	err := filepath.WalkDir("sets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if bytes.Contains(src, []byte("//go:build todo")) {
			return nil // build-excluded stub; not part of the database
		}
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !isCardNew(call.Fun) || len(call.Args) < 2 {
				return true
			}
			own := houseName(call.Args[1])
			if own == "" {
				return true
			}
			// call.Args[1] is the card's own house arg and is left untouched; only
			// the ability options that follow are checked for naming it again.
			for _, arg := range call.Args[2:] {
				ast.Inspect(arg, func(m ast.Node) bool {
					if houseName(m) == own {
						t.Errorf(
							"%s: card.House.%s names the card's own house; write card.House.Self",
							path, own,
						)
					}
					return true
				})
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// isCardNew reports whether fun is the selector card.New.
func isCardNew(fun ast.Expr) bool {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "New" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "card"
}

// houseName returns the house named by a card.House.<House> selector, or "" when
// expr is not such a selector or names the Self sentinel (which never counts as a
// hardcoded house).
func houseName(expr ast.Node) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name == "Self" {
		return ""
	}
	house, ok := sel.X.(*ast.SelectorExpr)
	if !ok || house.Sel.Name != "House" {
		return ""
	}
	pkg, ok := house.X.(*ast.Ident)
	if !ok || pkg.Name != "card" {
		return ""
	}
	return sel.Sel.Name
}

// TestReferencedCardIsConnected enforces that a card naming another card in its
// text — Grumpus Tamer tutoring a War Grumpus, Faygin returning an Urchin — shares
// a cluster with that card, so a generated deck never deals the tutor without its
// target. Co-membership counts in either direction: a puller and its pulled
// partner (Grumpus Tamer and War Grumpus) sit in one cluster, and a card only ever
// pulled in by X always shares X's pod. References are read straight from the
// definition: card text is generated from the effect tree rather than stored, so
// the only strings that equal another card's name are genuine references (an
// effect's SearchForName, a Target's Named filter, …).
func TestReferencedCardIsConnected(t *testing.T) {
	regs := card.Cards()
	names := make(map[string]bool, len(regs))
	cluster := make(map[string]string, len(regs))
	for _, rc := range regs {
		names[rc.Def.Name] = true
		if !rc.Profile.Cluster.Empty() {
			cluster[rc.Def.Name] = rc.Profile.Cluster.Name
		}
	}
	for _, rc := range regs {
		for ref := range referencedCardNames(reflect.ValueOf(rc.Def), names) {
			sameCluster := cluster[rc.Def.Name] != "" && cluster[rc.Def.Name] == cluster[ref]
			if ref == rc.Def.Name || sameCluster {
				continue
			}
			t.Errorf(
				"%s names %q but they share no cluster; add both to one cluster "+
					"with card.InCluster (ask the author for the pull rate)",
				rc.Def.Name, ref,
			)
		}
	}
}

// TestConnectedCardIsPulled is the mirror of TestReferencedCardIsConnected: a
// card of Rarity.Connected is kept out of the pool and never rolls on its own
// (deck generation indexes it by name and only places it through a cluster), so it
// must be a member of a cluster whose validated rolling trigger guarantees a
// rollable card places it — otherwise it can never reach a deck. Unlike a named
// reference, the puller need not mention the card in its text: the three Connected
// Horsemen ride in on Horseman of Pestilence's cluster, which names none of them.
func TestConnectedCardIsPulled(t *testing.T) {
	for _, rc := range card.Cards() {
		if rc.Def.Rarity != engine.Connected || !rc.Profile.Cluster.Empty() {
			continue
		}
		t.Errorf(
			"%s is Rarity.Connected but no cluster pulls it in, so it can never reach "+
				"a deck; add it to a cluster with card.InCluster",
			rc.Def.Name,
		)
	}
}

// TestSearchIsFollowedByShuffle enforces the KeyForge rule that a deck search is
// always followed by a shuffle: whenever an ability's effect tree contains a
// search effect (SearchForName or SearchDeck), the same tree must also contain a
// shuffle effect (any Shuffle* effect — ShuffleDeck, ShuffleIntoDeck, etc.). The
// search and the shuffle are deliberately separate effects, so this lint is what
// keeps a search from silently skipping its shuffle.
func TestSearchIsFollowedByShuffle(t *testing.T) {
	searches := map[string]bool{"SearchForName": true, "SearchDeck": true}
	for _, rc := range card.Cards() {
		for _, ab := range rc.Def.Abilities {
			types := effectTypeNames(reflect.ValueOf(ab.Effect))
			searched := false
			for name := range searches {
				if types[name] {
					searched = true
				}
			}
			if !searched {
				continue
			}
			shuffled := false
			for name := range types {
				if strings.HasPrefix(name, "Shuffle") {
					shuffled = true
				}
			}
			if !shuffled {
				t.Errorf(
					"%s searches its deck but the ability has no following shuffle; "+
						"a search must be followed by a shuffle (add card.ShuffleDeck{} "+
						"or a Shuffle...IntoDeck effect)",
					rc.Def.Name,
				)
			}
		}
	}
}

// effectTypeNames collects the type name of every struct that appears anywhere in
// v's tree (effect nodes, targets, nested effects). It only reads types and
// descends into fields, slices, interfaces, pointers, and maps, so it works on
// unexported fields too.
func effectTypeNames(v reflect.Value) map[string]bool {
	found := map[string]bool{}
	var walk func(v reflect.Value)
	walk = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.Struct:
			found[v.Type().Name()] = true
			for i := range v.NumField() {
				walk(v.Field(i))
			}
		case reflect.Slice, reflect.Array:
			for i := range v.Len() {
				walk(v.Index(i))
			}
		case reflect.Interface, reflect.Pointer:
			if !v.IsNil() {
				walk(v.Elem())
			}
		case reflect.Map:
			for iter := v.MapRange(); iter.Next(); {
				walk(iter.Value())
			}
		}
	}
	walk(v)
	return found
}

// referencedCardNames collects every registered card name (a key of names) that
// appears as a string value anywhere in def's effect tree, targets, and other
// fields. It reads unexported fields too, so a name tucked inside a Target's
// filter is found; only read operations are used, never Set or Interface.
func referencedCardNames(def reflect.Value, names map[string]bool) map[string]bool {
	found := map[string]bool{}
	var walk func(v reflect.Value)
	walk = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.String:
			if names[v.String()] {
				found[v.String()] = true
			}
		case reflect.Struct:
			for i := range v.NumField() {
				walk(v.Field(i))
			}
		case reflect.Slice, reflect.Array:
			for i := range v.Len() {
				walk(v.Index(i))
			}
		case reflect.Interface, reflect.Pointer:
			if !v.IsNil() {
				walk(v.Elem())
			}
		case reflect.Map:
			for iter := v.MapRange(); iter.Next(); {
				walk(iter.Key())
				walk(iter.Value())
			}
		}
	}
	walk(def)
	return found
}

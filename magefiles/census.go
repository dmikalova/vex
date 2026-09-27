//go:build mage

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dmikalova/vex/internal/census"
	"github.com/dmikalova/vex/internal/engine"
)

// engineDir is the package the census scans: the node families all live in the
// engine, beside the effect AST they belong to.
const engineDir = "internal/engine"

// censusSections are the rulebook sections whose terms the census claims, and so
// the sections an orphan term is reported from. The Card text section joins the
// Effects section here once the rulebook carries one.
var censusSections = []engine.Section{engine.SectionEffect}

// Census reports the node census's catalog and term gaps. Per family it prints
// the node types the source scan finds that no census row covers, grouped by the
// file that declares them, and the rows that match no type. It then prints the
// rulebook terms the rows name that no term carries — the prose still to write —
// and the terms in the census's own sections that no row claims. It only
// reports: the census is meant to be read while it is half-filled, so a gap is
// never an error here. The totality tests are what fail the build, family by
// family, as each one is gated.
func (Tool) Census() error {
	families := engine.Families()
	fmt.Println("CATALOG GAPS")
	for _, family := range families {
		if err := reportCatalogGaps(family); err != nil {
			return err
		}
	}
	fmt.Println()
	fmt.Println("TERM GAPS")
	reportTermGaps(families)
	return nil
}

// reportCatalogGaps prints one family's line, and under it the nodes it still
// owes rows for (grouped by source file) and the rows that name no node.
func reportCatalogGaps(family engine.Family) error {
	declared, err := census.Implementations(
		engineDir, family.Method, census.Params(family.Params...),
	)
	if err != nil {
		return fmt.Errorf("scanning for %s implementations: %w", family.Name, err)
	}
	catalogued := map[string]bool{}
	for _, row := range family.Rows {
		catalogued[row.Type] = true
	}
	missing := map[string][]string{}
	uncatalogued := 0
	for node, file := range declared {
		if !catalogued[node] {
			missing[file] = append(missing[file], node)
			uncatalogued++
		}
	}
	var orphans []string
	for node := range catalogued {
		if _, ok := declared[node]; !ok {
			orphans = append(orphans, node)
		}
	}
	fmt.Printf("\n  %-31s %3d declared, %3d catalogued  %s\n",
		family.Name+"."+family.Method, len(declared), len(family.Rows),
		familyState(uncatalogued, family.Gated),
	)
	for _, file := range sortedKeys(missing) {
		nodes := missing[file]
		sort.Strings(nodes)
		fmt.Printf("    %-34s %s\n", file, strings.Join(nodes, ", "))
	}
	if len(orphans) > 0 {
		sort.Strings(orphans)
		fmt.Printf("    %-34s %s\n", "rows naming no node:", strings.Join(orphans, ", "))
	}
	return nil
}

// familyState renders a family's standing: complete and gated, complete but not
// yet gated, or how many nodes it still owes rows for.
func familyState(uncatalogued int, gated bool) string {
	switch {
	case uncatalogued > 0:
		return fmt.Sprintf("— %d to catalogue", uncatalogued)
	case gated:
		return "— complete, gated"
	default:
		return "— complete, not gated yet"
	}
}

// reportTermGaps prints the rulebook terms the census still owes prose for, and
// the terms in the census's sections that no row claims.
func reportTermGaps(families []engine.Family) {
	claimed := map[string][]string{}
	for _, family := range families {
		for _, row := range family.Rows {
			if row.Rules.Term == "" {
				continue
			}
			claimed[row.Rules.Term] = append(claimed[row.Rules.Term], family.Name+"."+row.Type)
		}
	}
	written := map[string]bool{}
	sectioned := map[string]bool{}
	for _, term := range engine.RuleTerms() {
		written[term.Title] = true
		for _, section := range censusSections {
			if term.Section == section {
				sectioned[term.Title] = true
			}
		}
	}

	fmt.Println("\n  terms to write (named by a row, no term carries the title):")
	printTitles(claimed, func(title string) bool { return !written[title] })
	fmt.Println("\n  orphan terms (in the census's sections, claimed by no row):")
	for _, title := range sortedKeys(sectioned) {
		if len(claimed[title]) == 0 {
			fmt.Printf("    %s\n", title)
		}
	}
}

// printTitles prints each claimed title the keep predicate selects, with the
// nodes that name it, so the prose author sees at once what one term must cover.
func printTitles(claimed map[string][]string, keep func(string) bool) {
	for _, title := range sortedKeys(claimed) {
		if !keep(title) {
			continue
		}
		nodes := claimed[title]
		sort.Strings(nodes)
		fmt.Printf("    %-34s %s\n", title, strings.Join(nodes, ", "))
	}
}

// sortedKeys returns a map's keys in order, so two runs of the report read the
// same.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

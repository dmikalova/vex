package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
)

// stub generates a build-excluded stub file for every source card in a set that
// no implemented card covers yet. Each stub starts with a `//go:build todo`
// constraint so it is left out of the build, vet, test, lint, and — crucially —
// the card registry, keeping the database and coverage numbers honest. The stub
// carries the card's printed text and a TODO marker; implementing the card means
// removing the build tag and writing the real ability. Existing files (an
// implemented card or an earlier stub) are never overwritten.
func stub(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: cardlookup stub <setSlug>")
	}
	slug := args[0]

	// Touch the aggregator so every set package registers its cards.
	_ = cards.All()
	covered := coveredNumbers()[slug]

	sets := provenance.Sets()
	var set *provenance.Set
	for i := range sets {
		if sets[i].Slug == slug {
			set = &sets[i]
			break
		}
	}
	if set == nil {
		return fmt.Errorf("unknown set slug %q", slug)
	}

	dir := filepath.Join("internal", "cards", "sets", slug)
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("set package dir %q does not exist: %w", dir, err)
	}

	var miss []provenance.Card
	for _, c := range set.Cards {
		if !covered[c.Number] {
			miss = append(miss, c)
		}
	}
	sort.Slice(miss, func(i, j int) bool { return miss[i].Number < miss[j].Number })

	written, skipped := 0, 0
	for _, c := range miss {
		path := filepath.Join(dir, fileName(c.Name)+".go")
		if _, err := os.Stat(path); err == nil {
			skipped++
			continue
		}
		src, err := format.Source([]byte(stubSource(slug, set.SourceSet, c)))
		if err != nil {
			return fmt.Errorf("formatting stub for %q: %w", c.Name, err)
		}
		if err := os.WriteFile(path, src, 0o644); err != nil {
			return fmt.Errorf("writing %q: %w", path, err)
		}
		written++
	}
	reprints, err := writeSetConfig(dir, *set)
	if err != nil {
		return err
	}
	fmt.Printf(
		"%s: wrote %d stub(s), skipped %d existing of %d unimplemented; %d reprint(s) in 0set.go\n",
		set.Name,
		written,
		skipped,
		len(miss),
		reprints,
	)
	return nil
}

// writeSetConfig (re)generates the set package's 0set.go: the set's registrar
// (`var set = card.NewSet(card.XX)`, through which every card in the package
// declares its home set) plus the catalog of the cards this set reprints from
// earlier sets, so they join its deck-generation pool as full members without this
// package importing another set. It returns the reprint count, and always writes
// 0set.go — even with no reprints it declares the set registrar.
func writeSetConfig(dir string, set provenance.Set) (int, error) {
	reprints := reprintsForSet(set)
	path := filepath.Join(dir, "0set.go")
	src, err := format.Source([]byte(setConfigSource(set, reprints)))
	if err != nil {
		return 0, fmt.Errorf("formatting 0set.go for %q: %w", set.Name, err)
	}
	if err := os.WriteFile(path, src, 0o644); err != nil {
		return 0, fmt.Errorf("writing %q: %w", path, err)
	}
	return len(reprints), nil
}

// reprintsForSet returns the cards in set's catalog that are implemented in a
// different set — the reprints it claims as pool members via 0set.go — in
// collector-number order. A source printing is a reprint when an implemented
// card built from another printing of that same card lives in a different set:
// each implemented card's provenance Refs are resolved to the names of the source
// cards it was built from, and any same-named printing in this set is that card.
// This is what lets a card tag only the printing it was built from (and a card
// renamed away from its printed name still resolve). The emitted claim carries
// the implementing card's own name so it resolves in the aggregator, and each
// implementing card appears at most once.
func reprintsForSet(set provenance.Set) []provenance.Card {
	refName := sourceNameByRef()
	type impl struct{ name, home string }
	bySourceName := map[string]impl{}
	for _, rc := range card.Cards() {
		// The card's actual home set is where it is implemented (InSet), not its
		// provenance origin — an anomaly's provenance points at Worlds Collide but it
		// lives in Anomaly Expansion, so keying home off provenance would claim it as
		// a reprint of itself in its own set.
		home := rc.Set.Name
		for _, ref := range rc.Provenance {
			if name, ok := refName[refKey(ref)]; ok {
				bySourceName[normalizeName(name)] = impl{rc.Def.Name, home}
			}
		}
		// Also key by the card's own name, so a card whose provenance Ref does not
		// resolve to a catalog printing (an anomaly's dangling A0x Ref, since the
		// source set's catalog omits its anomalies) is still recognized as the
		// implementation of a same-named printing in another set (Orb of Wonder is a
		// Worlds Collide anomaly in Anomaly Expansion and a Mass Mutation card).
		if _, ok := bySourceName[normalizeName(rc.Def.Name)]; !ok {
			bySourceName[normalizeName(rc.Def.Name)] = impl{rc.Def.Name, home}
		}
	}

	// Cluster awareness: a card pulled into a cluster by a lead (card.InCluster)
	// only forms that cluster when the lead is in the same set's pool. When a set
	// reprints a pulled member without its lead, the outcome depends on the member's
	// rarity: a rollable member (any rarity but Connected) prints in this set on its
	// own — KeyForge prints Sensor Chief Garcia in Mass Mutation without its blaster
	// — so it is reprinted and the aggregator drops its lead-less cluster for this
	// set (ownPool/reprintPoolCard). A Connected member never rolls alone, so it can
	// never appear without its lead; reprinting it would be a card that can never be
	// drawn, so skip it here (the aggregator also panics if one is forced in by hand).
	regByName := map[string]card.RegisteredCard{}
	leadOfCluster := map[string]string{}
	for _, rc := range card.Cards() {
		regByName[normalizeName(rc.Def.Name)] = rc
		if m := rc.Profile.Cluster; m.Name != "" && m.Lead {
			leadOfCluster[m.Name] = rc.Def.Name
		}
	}
	pooled := map[string]bool{}
	for _, c := range set.Cards {
		if im, ok := bySourceName[normalizeName(c.Name)]; ok {
			pooled[normalizeName(im.name)] = true
		}
	}
	leadAbsent := func(name string) bool {
		rc, ok := regByName[normalizeName(name)]
		if !ok {
			return false
		}
		m := rc.Profile.Cluster
		if m.Name == "" || m.Lead {
			return false
		}
		lead := leadOfCluster[m.Name]
		return lead == "" || !pooled[normalizeName(lead)]
	}
	// unpooledOrphan reports a reprint that cannot join this set's pool: a Connected
	// cluster member whose lead is absent. A rollable orphan is kept.
	unpooledOrphan := func(name string) bool {
		if !leadAbsent(name) {
			return false
		}
		rc, ok := regByName[normalizeName(name)]
		return ok && rc.Def.Rarity == card.Rarity.Connected
	}

	var out []provenance.Card
	seen := map[string]bool{}
	for _, c := range set.Cards {
		im, ok := bySourceName[normalizeName(c.Name)]
		if !ok || im.home == set.Name || seen[im.name] {
			continue
		}
		if unpooledOrphan(im.name) {
			continue
		}
		seen[im.name] = true
		c.Name = im.name
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out
}

// setConfigSource renders the Go source for a set package's generated 0set.go.
func setConfigSource(set provenance.Set, reprints []provenance.Card) string {
	var b strings.Builder
	fmt.Fprintf(&b, "// Code generated by \"mage tool:stub %s\"; DO NOT EDIT.\n\n", set.Slug)
	b.WriteString("package " + set.Slug + "\n\n")
	b.WriteString("import \"github.com/dmikalova/vactrol/internal/card\"\n\n")
	b.WriteString(
		"// set is " + set.Name + "'s registrar: every card in this package declares its\n",
	)
	b.WriteString(
		"// home set by registering through set.New, so deck generation groups the card\n",
	)
	b.WriteString(
		"// by this set alone and never infers it from provenance (ADR 0003).\n",
	)
	fmt.Fprintf(&b, "var set = card.NewSet(card.%s)\n", provCodeVar(set.SourceSet))
	if len(reprints) > 0 {
		b.WriteString("\n")
		b.WriteString(
			"// " + set.Name + " reprints these cards from earlier sets. Each is implemented\n",
		)
		b.WriteString(
			"// in the set that introduced it; here it is claimed as a full member of this\n",
		)
		b.WriteString(
			"// set's deck-generation pool, resolved by name in the cards aggregator so this\n",
		)
		b.WriteString(
			"// package imports no other set. The list is this set's provenance catalog minus\n",
		)
		b.WriteString("// its own new cards; regenerate with `mage tool:stub " + set.Slug + "`.\n")
		b.WriteString("func init() {\n")
		for _, c := range reprints {
			fmt.Fprintf(&b, "\tset.Reprint(%q, %s)\n", c.Number, quote(c.Name))
		}
		b.WriteString("}\n")
	}
	return b.String()
}

// stubSource renders the Go source for one build-excluded stub file.
func stubSource(pkg string, set provenance.SourceSet, c provenance.Card) string {
	var b strings.Builder
	b.WriteString("//go:build todo\n\n")
	b.WriteString("package " + pkg + "\n\n")
	b.WriteString("import \"github.com/dmikalova/vactrol/internal/card\"\n\n")

	// Doc comment: the printed text (unsure marker), so the file documents the
	// card even though gencomments skips build-excluded files.
	b.WriteString("// " + varName(c.Name) + "\n//\n")
	b.WriteString("// TODO(stub): unimplemented. Remove the //go:build todo tag and\n")
	b.WriteString("// implement the ability once the needed effect exists.\n//\n")
	b.WriteString("//\tHouse:  " + titleWord(c.House) + "\n")
	b.WriteString("//\tType:   " + cardTypeName(c.Type) + "\n")
	b.WriteString("//\tRarity: " + rarityName(c.Rarity) + "\n")
	if isCreature(c.Type) {
		fmt.Fprintf(&b, "//\tPower:  %d\n", c.Power)
		if c.Armor > 0 {
			fmt.Fprintf(&b, "//\tArmor:  %d\n", c.Armor)
		}
	}
	if c.Amber > 0 {
		fmt.Fprintf(&b, "//\tÆmber:  %d\n", c.Amber)
	}
	if len(c.Traits) > 0 {
		b.WriteString("//\tTraits: " + strings.Join(titleWords(c.Traits), " • ") + "\n")
	}
	if c.Text != "" {
		b.WriteString("//\n")
		for line := range strings.SplitSeq(c.Text, "\n") {
			b.WriteString("//\t" + line + "\n")
		}
	}

	// A vanilla card.New skeleton with stats — compiles as-is once the build tag
	// is removed, leaving only the ability to add.
	b.WriteString("var " + varName(c.Name) + " = set.New(\n")
	b.WriteString("\t" + quote(c.Name) + ",\n")
	b.WriteString("\tcard.House." + titleWord(c.House) + ",\n")
	b.WriteString("\tcard.Type." + cardTypeName(c.Type) + ",\n")
	if mapped := rarityName(c.Rarity); !strings.EqualFold(mapped, c.Rarity) {
		fmt.Fprintf(&b,
			"\t// TODO(variant): rarity relabelled from %s to %s — handle manually\n",
			c.Rarity, mapped)
	}
	b.WriteString("\tcard.Rarity." + rarityName(c.Rarity) + ",\n")
	fmt.Fprintf(&b, "\tcard.Provenance(card.%s, %q),\n", provCodeVar(set), c.Number)
	if isCreature(c.Type) {
		fmt.Fprintf(&b, "\tcard.WithPower(%d),\n", c.Power)
		if c.Armor > 0 {
			fmt.Fprintf(&b, "\tcard.WithArmor(%d),\n", c.Armor)
		}
	}
	if c.Amber > 0 {
		fmt.Fprintf(&b, "\tcard.WithBonus(%s),\n", bonusAemberArgs(c.Amber))
	}
	if len(c.Traits) > 0 {
		named := make([]string, len(c.Traits))
		for i, t := range c.Traits {
			named[i] = "card.Traits." + varName(t)
		}
		b.WriteString("\tcard.WithTraits(" + strings.Join(named, ", ") + "),\n")
	}
	b.WriteString("\t// TODO(stub): add WithKeywords / WithAbility for the printed text above.\n")
	b.WriteString(")\n")
	return b.String()
}

// bonusAemberArgs renders n Æmber bonus icons as WithBonus arguments.
func bonusAemberArgs(n int) string {
	args := make([]string, n)
	for i := range args {
		args[i] = "card.Bonus.Aember"
	}
	return strings.Join(args, ", ")
}

func isCreature(t string) bool { return strings.EqualFold(t, "creature") }

// facadeRarities is the set of rarity names the card facade defines
// (card.Rarity.*). Any source rarity outside it has no facade value of its own.
var facadeRarities = map[string]bool{
	"common": true, "uncommon": true, "rare": true, "special": true, "connected": true,
}

// rarityName maps a source rarity onto the facade Rarity namespace field. The
// facade has no value for the special-slot rarities the catalogs carry — Variant
// (ambassadors, brews, Master of N), FIXED (anomalies), Token, Evil Twin, The
// Tide — so every rarity outside facadeRarities renders as Special; the standard
// rarities pass through unchanged.
func rarityName(rarity string) string {
	if facadeRarities[strings.ToLower(rarity)] {
		return rarity
	}
	return "Special"
}

// cardTypeName maps a source card type onto the facade Type namespace field.
// KeyForge's "action" card type is named Tactic in this repo (card-wording rule 19).
func cardTypeName(t string) string {
	switch strings.ToLower(t) {
	case "action":
		return "Tactic"
	case "creature":
		return "Creature"
	case "artifact":
		return "Artifact"
	case "upgrade":
		return "Upgrade"
	default:
		return titleWord(t)
	}
}

// varName turns a card name into an exported Go identifier: it keeps only ASCII
// letters and digits, capitalizing the first letter of each whitespace-separated
// word while preserving existing capitals (so "EMP Blast" stays EMPBlast), and
// drops every other character without a word break ("Coward's End" becomes
// CowardsEnd, "Z-Y-X Researcher" becomes ZYXResearcher). Source names arriving
// from the provenance importer are already ASCII-folded, so no accented letter or
// Æ reaches here.
func varName(name string) string {
	var b strings.Builder
	upNext := true
	for _, r := range name {
		switch {
		case isASCIIAlnum(r):
			if upNext {
				b.WriteRune(unicode.ToUpper(r))
				upNext = false
			} else {
				b.WriteRune(r)
			}
		case unicode.IsSpace(r):
			upNext = true
		default:
			// drop any other character without a word break
		}
	}
	id := b.String()
	if id == "" || (id[0] >= '0' && id[0] <= '9') {
		id = "Card" + id
	}
	return id
}

// fileName turns a card name into a snake_case file base name: it keeps only
// ASCII letters and digits (lowercased), joins whitespace-separated words with
// '_', and drops every other character without a separator ("Coward's End"
// becomes cowards_end, "Z-Y-X Researcher" becomes zyx_researcher).
func fileName(name string) string {
	var b strings.Builder
	prevSep := true
	for _, r := range name {
		switch {
		case isASCIIAlnum(r):
			b.WriteRune(unicode.ToLower(r))
			prevSep = false
		case unicode.IsSpace(r):
			if !prevSep {
				b.WriteRune('_')
				prevSep = true
			}
		default:
			// drop any other character without a separator
		}
	}
	return strings.Trim(b.String(), "_")
}

// isASCIIAlnum reports whether r is a plain ASCII letter or digit — the only
// runes varName and fileName keep, since a card name reaches them already
// ASCII-folded by the provenance importer.
func isASCIIAlnum(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
}

// titleWord capitalizes the first rune of a single lowercase source token.
func titleWord(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func titleWords(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = titleWord(s)
	}
	return out
}

func quote(s string) string { return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"` }

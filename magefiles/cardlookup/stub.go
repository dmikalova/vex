package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

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
	fmt.Printf("%s: wrote %d stub(s), skipped %d existing of %d unimplemented\n",
		set.Name, written, skipped, len(miss))
	return nil
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
	b.WriteString("//\tRarity: " + c.Rarity + "\n")
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
		for _, line := range strings.Split(c.Text, "\n") {
			b.WriteString("//\t" + line + "\n")
		}
	}

	// A vanilla card.New skeleton with stats — compiles as-is once the build tag
	// is removed, leaving only the ability to add.
	b.WriteString("var " + varName(c.Name) + " = card.New(\n")
	b.WriteString("\t" + quote(c.Name) + ",\n")
	b.WriteString("\tcard.House." + titleWord(c.House) + ",\n")
	b.WriteString("\tcard.Type." + cardTypeName(c.Type) + ",\n")
	b.WriteString("\tcard.Rarity." + c.Rarity + ",\n")
	fmt.Fprintf(&b, "\tcard.Provenance(card.%s, %d),\n", provCodeVar(set), c.Number)
	if isCreature(c.Type) {
		fmt.Fprintf(&b, "\tcard.WithPower(%d),\n", c.Power)
		if c.Armor > 0 {
			fmt.Fprintf(&b, "\tcard.WithArmor(%d),\n", c.Armor)
		}
	}
	if c.Amber > 0 {
		fmt.Fprintf(&b, "\tcard.WithAemberBonus(%d),\n", c.Amber)
	}
	if len(c.Traits) > 0 {
		quoted := make([]string, len(c.Traits))
		for i, t := range c.Traits {
			quoted[i] = quote(titleWord(t))
		}
		b.WriteString("\tcard.WithTraits(" + strings.Join(quoted, ", ") + "),\n")
	}
	b.WriteString("\t// TODO(stub): add WithKeywords / WithAbility for the printed text above.\n")
	b.WriteString(")\n")
	return b.String()
}

func isCreature(t string) bool { return strings.EqualFold(t, "creature") }

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

// varName turns a card name into an exported Go identifier, dropping punctuation
// and capitalizing each word's first rune while preserving existing capitals
// (so "EMP Blast" stays EMPBlast, "Coward's End" becomes CowardsEnd).
func varName(name string) string {
	var b strings.Builder
	upNext := true
	for _, r := range name {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if upNext {
				b.WriteRune(unicode.ToUpper(r))
				upNext = false
			} else {
				b.WriteRune(r)
			}
		case r == '\'':
			// drop apostrophes without a word break (Coward's -> Cowards)
		default:
			upNext = true
		}
	}
	id := b.String()
	if id == "" || unicode.IsDigit(rune(id[0])) {
		id = "Card" + id
	}
	return id
}

// fileName turns a card name into a snake_case file base name.
func fileName(name string) string {
	var b strings.Builder
	prevSep := true
	for _, r := range name {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
			prevSep = false
		case r == '\'':
			// drop apostrophes without a separator (Coward's -> cowards)
		default:
			if !prevSep {
				b.WriteRune('_')
				prevSep = true
			}
		}
	}
	return strings.Trim(b.String(), "_")
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

package worldscollide

import (
	"hash/fnv"
	"math/rand"
	"sort"
	"strings"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Trait Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Template: its concrete card is materialized per deck at generation.
var TraitBane = card.New(
	"Trait Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	// The generative bane subsumes Worlds Collide's seven printed single-trait
	// banes; it carries their provenance so coverage stays honest without them.
	card.Provenance(card.WC, "89"),  // Giants' Bane
	card.Provenance(card.WC, "122"), // Beasts' Bane
	card.Provenance(card.WC, "123"), // Demons' Bane
	card.Provenance(card.WC, "125"), // Dinosaurs' Bane
	card.Provenance(card.WC, "126"), // Humans' Bane
	card.Provenance(card.WC, "127"), // Scientists' Bane
	card.Provenance(card.WC, "128"), // Thieves' Bane
	card.WithAemberBonus(1),
	card.Template(baneFor),
)

// baneFor materializes a Bane for a random combination of three Houses: a Dis
// Tactic that destroys a creature of each House's most common creature trait.
// Every unordered triple of the nine Houses yields one of C(9,3)=84 variants,
// named by splicing a growing window of each trait, so the card is generative and
// does not care which Houses the deck actually runs.
func baneFor(_ card.SlotContext, r *rand.Rand) card.Definition {
	houses := allHouses()
	r.Shuffle(len(houses), func(i, j int) { houses[i], houses[j] = houses[j], houses[i] })
	return baneForHouses(houses[0], houses[1], houses[2])
}

// baneForHouses builds the Bane for three specific Houses, destroying a creature
// of each House's most common creature trait. The traits are ordered by a shuffle
// seeded from their own names, so a given triple always produces the same order
// (for both the name and the destroy clause) but not the alphabetical one.
func baneForHouses(h1, h2, h3 engine.House) card.Definition {
	traits := []engine.Trait{
		card.MostCommonCreatureTrait(h1),
		card.MostCommonCreatureTrait(h2),
		card.MostCommonCreatureTrait(h3),
	}
	shuffleTraits(traits)
	return card.Build(
		baneName(traits),
		card.House.Dis,
		card.Type.Tactic,
		card.Rarity.Rare,
		card.WithAemberBonus(1),
		card.WithAbility(
			card.Trigger.Play, card.Sequence{Effects: []card.Effect{
				card.Destroy{Target: card.Target.Creature.WithTrait(traits[0])},
				card.Destroy{Target: card.Target.Creature.WithTrait(traits[1])},
				card.Destroy{Target: card.Target.Creature.WithTrait(traits[2])},
			}}),
	)
}

// shuffleTraits reorders traits in place with a rand seeded from the sorted trait
// names, so the order is deterministic and independent of the argument order.
func shuffleTraits(traits []engine.Trait) {
	sort.Slice(traits, func(i, j int) bool { return traits[i].String() < traits[j].String() })
	h := fnv.New64a()
	for _, t := range traits {
		h.Write([]byte(t.String()))
		h.Write([]byte{0})
	}
	r := rand.New(rand.NewSource(int64(h.Sum64())))
	r.Shuffle(len(traits), func(i, j int) { traits[i], traits[j] = traits[j], traits[i] })
}

// baneName splices a growing window of each of the three traits — two, three,
// then four letters — into one generative word and appends "'s Bane", e.g.
// Beast+Scientist+Thief -> "Bescithief's Bane".
func baneName(traits []engine.Trait) string {
	var b strings.Builder
	for i, t := range traits {
		b.WriteString(window(t.String(), 2+i))
	}
	word := b.String()
	return strings.ToUpper(word[:1]) + word[1:] + "'s Bane"
}

// window returns the first n letters of s, lowercased, extended by one closing
// consonant when the nth letter would otherwise split a vowel from the consonant
// after it (so a window ends on a consonant rather than a bare vowel).
func window(s string, n int) string {
	s = strings.ToLower(s)
	if n > len(s) {
		n = len(s)
	}
	if n < len(s) && isVowel(s[n-1]) && !isVowel(s[n]) {
		n++
	}
	return s[:n]
}

// isVowel reports whether c is an English vowel letter.
func isVowel(c byte) bool { return strings.IndexByte("aeiou", c) >= 0 }

// allHouses is every real House, in enum order.
func allHouses() []engine.House {
	houses := make([]engine.House, 0, engine.NumHouses)
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		houses = append(houses, h)
	}
	return houses
}

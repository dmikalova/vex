package massmutation

import (
	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/engine"
)

// The mutant cycle. Mass Mutation's mutants are house hybrids: each creature
// carries a house prefix and a house suffix, and its stats and abilities are the
// sum of the two houses' contributions. The suffix house owns the creature — it
// is the creature's house and lends its creature trait — while the prefix is any
// other house. Rather than hand-write all 42 printed combinations, this file
// models each house's fixed contribution once (mutantHouses) and composes every
// mutant from a prefix and a suffix (mutantOptions), registering the 42 that Mass
// Mutation prints (mutantProvenance).
//
// Every house — even Brobnar and Mars, which Mass Mutation never prints as a
// mutant — defines a full contribution, so a later set that mutates them has the
// data ready. TestMutantHousesComplete enforces that promise: a house added to
// the game without its mutant data fails the build.

// mutantHouse is one house's fixed contribution to a mutant creature: the name
// fragments it lends in the prefix and suffix positions, the creature trait its
// suffix stamps, the power and armor it adds, and the keywords and abilities it
// grants. A mutant's power and armor are the sum of its two houses' values; its
// keywords and abilities are their union; its trait and owning house come from
// the suffix alone.
type mutantHouse struct {
	prefix    string
	suffix    string
	trait     card.Trait
	power     int
	armor     int
	keywords  []card.KeywordValue
	abilities []card.Ability
}

// mutantHouses is every KeyForge house's mutant contribution. The seven houses
// Mass Mutation prints supply their existing signature abilities; Brobnar and
// Mars are proposals, dormant until a set prints their mutants.
var mutantHouses = map[engine.House]mutantHouse{
	// Dis lends its demons a "steal on death" bite.
	card.House.Dis: {
		prefix: "Daemo",
		suffix: "Fiend",
		trait:  card.Traits.Demon,
		power:  1,
		abilities: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.StealAember{Amount: 1},
		}},
	},
	// Logos lends its scientists the reap-cycle "discard then draw".
	card.House.Logos: {
		prefix: "Techno",
		suffix: "Bot",
		trait:  card.Traits.Scientist,
		power:  2,
		abilities: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect: card.Then{
				First: card.DiscardCard{
					Player:    card.Controller,
					Zones:     []card.Zone{card.Hand},
					Selection: card.Chosen{},
				},
				Result: card.Draw{Amount: 1},
			},
		}},
	},
	// Untamed lends its beasts Skirmish.
	card.House.Untamed: {
		prefix:   "Lyco",
		suffix:   "Beast",
		trait:    card.Traits.Beast,
		power:    2,
		keywords: []card.KeywordValue{card.Keyword.Skirmish},
	},
	// Sanctum lends its knights armor.
	card.House.Sanctum: {
		prefix: "Sacro",
		suffix: "Knight",
		trait:  card.Traits.Knight,
		power:  3,
		armor:  2,
	},
	// Shadows lends its thieves Elusive.
	card.House.Shadows: {
		prefix:   "Umbra",
		suffix:   "Thief",
		trait:    card.Traits.Thief,
		power:    1,
		keywords: []card.KeywordValue{card.Keyword.Elusive},
	},
	// Saurian lends its dinosaurs an "exalt to deal 3 damage" play.
	card.House.Saurian: {
		prefix: "Dino",
		suffix: "Saurus",
		trait:  card.Traits.Dinosaur,
		power:  3,
		abilities: []card.Ability{{
			Trigger: card.Trigger.Play,
			Effect: card.May{
				Do: card.Then{
					First: card.Exalt{
						Target: card.Target.This,
						Amount: 1,
					},
					Result: card.DealDamage{
						Amount: 3,
						Target: card.Target.Creature,
					},
				},
			},
		}},
	},
	// Star Alliance lends its aliens a "dig 3, keep 1, bottom 1" fight.
	card.House.StarAlliance: {
		prefix: "Xeno",
		suffix: "Alien",
		trait:  card.Traits.Alien,
		power:  2,
		abilities: []card.Ability{{
			Trigger: card.Trigger.Fight,
			Effect: card.LookAtTopOfDeck{
				Amount: 3,
				Then: []card.TopAct{
					card.ChooseAndMove{
						Cards: 1,
						Dest:  card.Into.Hand,
					},
					card.ChooseAndMove{
						Cards: 1,
						Dest:  card.Into.BottomOfDeck,
					},
				},
			},
		}},
	},
	// Brobnar (proposal) lends its giants a "fight with a neighbor" play.
	card.House.Brobnar: {
		prefix: "Bruto",
		suffix: "Zerker",
		trait:  card.Traits.Giant,
		power:  4,
		abilities: []card.Ability{{
			Trigger: card.Trigger.Play,
			Effect: card.OnChooseCreature{
				Target: card.Target.Creature.Neighboring(),
				Verbs:  []card.CreatureVerb{card.FightVerb{}},
			},
		}},
	},
	// Mars (proposal) lends its martians a "stun a creature" play.
	card.House.Mars: {
		prefix: "Xzyz",
		suffix: "Zxyx",
		trait:  card.Traits.Martian,
		power:  2,
		abilities: []card.Ability{{
			Trigger: card.Trigger.Play,
			Effect:  card.Stun{Target: card.Target.Creature},
		}},
	},
}

// mutantName is a mutant's card name: the prefix house's prefix form joined to the
// suffix house's suffix form, e.g. Daemo-Bot.
func mutantName(prefix, suffix engine.House) string {
	return mutantHouses[prefix].prefix + "-" + mutantHouses[suffix].suffix
}

// mutantOptions assembles the gameplay options for the mutant with the given
// prefix and suffix houses: summed power and armor, unioned keywords and
// abilities, and the Mutant trait plus the suffix house's creature trait. The
// caller appends provenance.
func mutantOptions(prefix, suffix engine.House) []card.Option {
	p, s := mutantHouses[prefix], mutantHouses[suffix]
	opts := []card.Option{
		card.WithPower(p.power + s.power),
		card.WithTraits(card.Traits.Mutant, s.trait),
	}
	if armor := p.armor + s.armor; armor > 0 {
		opts = append(opts, card.WithArmor(armor))
	}
	if kw := append(append([]card.KeywordValue{}, p.keywords...), s.keywords...); len(kw) > 0 {
		opts = append(opts, card.WithKeywords(kw...))
	}
	for _, a := range append(append([]card.Ability{}, p.abilities...), s.abilities...) {
		opts = append(opts, card.WithAbility(a.Trigger, a.Effect))
	}
	return opts
}

// mutantProvenance maps each printed mutant to its Mass Mutation collector number,
// keyed [suffix][prefix] — the suffix house owns the creature. Its 42 entries are
// exactly the mutants Mass Mutation prints: the seven native houses, each paired
// as a suffix with the six other native houses as prefixes.
var mutantProvenance = map[engine.House]map[engine.House]string{
	card.House.Dis: { // -Fiend
		card.House.Saurian:      "055",
		card.House.Untamed:      "059",
		card.House.Sanctum:      "061",
		card.House.Logos:        "016",
		card.House.Shadows:      "063",
		card.House.StarAlliance: "065",
	},
	card.House.Saurian: { // -Saurus
		card.House.Dis:          "190",
		card.House.Untamed:      "235",
		card.House.Sanctum:      "240",
		card.House.Logos:        "241",
		card.House.Shadows:      "242",
		card.House.StarAlliance: "243",
	},
	card.House.Logos: { // -Bot
		card.House.Dis:          "068",
		card.House.Saurian:      "119",
		card.House.Untamed:      "120",
		card.House.Sanctum:      "122",
		card.House.Shadows:      "123",
		card.House.StarAlliance: "124",
	},
	card.House.Untamed: { // -Beast
		card.House.Dis:          "363",
		card.House.Saurian:      "417",
		card.House.Sanctum:      "418",
		card.House.Logos:        "419",
		card.House.Shadows:      "420",
		card.House.StarAlliance: "421",
	},
	card.House.Shadows: { // -Thief
		card.House.Dis:          "246",
		card.House.Saurian:      "297",
		card.House.Untamed:      "298",
		card.House.Sanctum:      "299",
		card.House.Logos:        "300",
		card.House.StarAlliance: "301",
	},
	card.House.Sanctum: { // -Knight
		card.House.Dis:          "132",
		card.House.Saurian:      "178",
		card.House.Untamed:      "179",
		card.House.Logos:        "180",
		card.House.Shadows:      "181",
		card.House.StarAlliance: "182",
	},
	card.House.StarAlliance: { // -Alien
		card.House.Dis:     "306",
		card.House.Saurian: "357",
		card.House.Untamed: "358",
		card.House.Sanctum: "359",
		card.House.Logos:   "360",
		card.House.Shadows: "361",
	},
}

// init registers every printed mutant as an Uncommon creature of its suffix
// house. Registration is driven by mutantProvenance, so the set's pool always
// matches the printed cards exactly.
func init() {
	for suffix, byPrefix := range mutantProvenance {
		for prefix, number := range byPrefix {
			opts := append(mutantOptions(prefix, suffix), card.Provenance(card.MM, number))
			set.New(
				mutantName(prefix, suffix),
				suffix,
				card.Type.Creature,
				card.Rarity.Uncommon,
				opts...,
			)
		}
	}
}

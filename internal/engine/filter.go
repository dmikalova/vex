package engine

import (
	"fmt"
	"strings"
)

// A Filter is a predicate over one card, decidable by looking at that card alone
// however much board it reads to do so. It is the per-card half of a Target — the
// adjectives and trailing clauses that narrow "each creature" to "each damaged
// Mars creature with power 3 or lower" — held as one comparable value so a card
// names its narrowing in a struct literal rather than through a chain of builder
// methods (ADR 0005).
//
// Every axis conjoins with every other, so a Filter admits a card that satisfies
// all of them; MatchAny disjoins the identity axes instead, for the "Mars or
// Robot creature" that must be one set rather than two. The zero value narrows
// nothing.
//
// A Filter both selects (matches) and renders (qualifyNoun, clauses), which is
// what keeps printed text and behaviour from drifting. The rendering splits in
// two because the consumer owns what sits between the halves: the quantifier and
// the article ("each …", "a …", "another …") are the Target's, so it prefixes the
// adjectives, inserts its quantifier, then appends the trailing clauses.
//
// A rule that needs the whole candidate set — comparing candidates to each other,
// or prompting across them — is not a Filter but a Refinement; see
// target_refinement.go.
type Filter struct {
	// The identity axes: what a card is, printed on it and true wherever it sits.
	// A consumer may ask them of a card in any zone.

	// Type narrows to cards of this type; the zero value narrows nothing. It also
	// supplies the noun a consumer with no noun of its own prints ("upgrade" rather
	// than "card").
	Type CardType
	// Gigantic narrows to the halves of a gigantic creature, either half, and prints
	// "gigantic creature" (the tutors search for a gigantic's halves).
	Gigantic bool
	// Trait narrows to cards carrying this trait; the zero value narrows nothing.
	Trait Trait
	// ExceptTrait narrows to cards lacking this trait, rendering the "non-<trait>"
	// qualifier. It is the only negated axis.
	ExceptTrait Trait
	// House narrows to the cards the matcher admits — a named house, every house but
	// one, the chosen house, the active house, or the house of the card in context
	// (ctx.It). The zero value (any house) narrows nothing.
	House HouseMatcher
	// MatchAny makes the house and trait axes disjoin instead of conjoin, so "each
	// Mars or Robot creature" is one filter rather than two sequenced targets and a
	// Mars Robot is affected once (EMP Blast). It is a flag rather than a list of
	// alternatives because a Filter must stay comparable (ADR 0005). It governs the
	// identity axes only — power, damage, and position always conjoin, where "or" has
	// no clear meaning.
	MatchAny bool
	// HouseWithMostCreatures narrows to creatures of the house with the most
	// creatures in play, counting both players' battlelines; on a tie every tied
	// house's creatures are eligible so the chooser picks among them (Etaromme). It
	// renders "of the house with the most creatures in play".
	HouseWithMostCreatures bool
	// SharesTrait narrows to cards sharing at least one trait with the card in
	// context (ctx.It), rendering "that shares a trait with it".
	SharesTrait bool
	// The in-play axes: what has happened to a card and where it stands. They read
	// the board, so they are meaningful only of a card in play, and a consumer that
	// points at a deck or a pile rejects them in validate rather than matching
	// nothing.

	// Power narrows to creatures whose power meets a bound — at most, at least,
	// exactly, odd, or even. The zero value bounds nothing.
	Power PowerBound
	// Damage narrows to creatures that have damage on them, or to those that have
	// none. The zero value narrows nothing.
	Damage DamagePresence
	// Stunned narrows to creatures that are currently stunned.
	Stunned bool
	// Ready narrows to creatures that are not exhausted (Swap Widget's "a ready
	// friendly Mars creature").
	Ready bool
	// Aember narrows to cards that have Æmber on them, or to those that have none
	// (Draining Touch destroys a creature with no Æmber on it). The zero value
	// narrows nothing.
	Aember AemberPresence
	// NoBonusIcons narrows to cards printing no bonus icons, rendering " with
	// no bonus icons" (Wail of the Damned destroys a creature with no bonus icons).
	NoBonusIcons bool
	// Counter narrows to cards carrying a generic counter of this kind, rendering
	// " with a doom counter" and the like (Wretched Doll destroys every creature with
	// a doom counter). CounterNone leaves the axis off.
	Counter CounterKind
	// Armor narrows to creatures that have armor at all, rendering " with armor".
	// It reads the creature's armor value, not what is left of it, so a creature that
	// has already spent its armor absorbing damage still has armor.
	Armor bool
	// Upgrade narrows to creatures that have at least one upgrade attached,
	// rendering " with an upgrade" (Tachyon Pulse exhausts each creature with an
	// upgrade).
	Upgrade bool
	// HouseWithAtLeast narrows to creatures whose house has at least this many
	// creatures in play, counting that house across both players' battlelines — a
	// house is a house regardless of who controls its creatures (No Safety in
	// Numbers). Zero leaves the axis off.
	//
	// It reads the whole board but is still decided one creature at a time, which is
	// what makes it a filter rather than a Refinement: it never compares candidates
	// to each other.
	HouseWithAtLeast int
	// WithoutSharedTrait narrows to creatures that share no trait with any other
	// creature in the same controller's battleline — a "loner". A creature with no
	// traits shares no trait, so it is kept (Good of the Many).
	WithoutSharedTrait bool
	// SharesHouseWithNeighbors narrows to creatures sharing a house with at least this
	// many of their battleline neighbors, rendering "that shares a house with N of
	// its neighbors" (Groupthink Tank, Mini Groupthink Tank). Zero leaves the axis
	// off.
	SharesHouseWithNeighbors int
	// Keyword narrows to creatures that have this keyword, rendering it as an
	// adjective: "each elusive creature".
	Keyword Keyword
	// Position narrows to the cards standing in one place in a battleline — on a
	// flank, off a flank, in the center, or to one side of the source card. The zero
	// value narrows nothing.
	Position Position
	// Neighboring narrows to the source card's battleline neighbors (the creatures
	// immediately to its left and right).
	Neighboring bool
	// Except drops one card from the admitted set — the source card for the "other"
	// cards a card names. The zero value drops none.
	Except Exclusion
	// Name narrows to cards with this printed name, and replaces the rendered noun
	// with it: a card that names another card outright says "an Ancient Bear", not
	// "an Ancient Bear creature".
	Name string
}

// Narrows reports whether the filter sets any axis at all. The zero value admits
// every card and adds no word to the phrase, which is how a consumer says "no
// restriction" without a second field to mean it.
func (f Filter) Narrows() bool {
	return f != Filter{}
}

// branches returns one single-axis Filter per axis this filter sets, in the order
// those axes print. It is what MatchAny means: the filter admits a card any branch
// admits, and prints one phrase per branch. A conjoining filter never asks for it.
//
// The order is the printed order, which is why Type leads (Chief Engineer Walls
// prints "an upgrade or Robot card") and House precedes Trait (EMP Blast prints
// "each Mars or Robot creature").
func (f Filter) branches() []Filter {
	var out []Filter
	add := func(set bool, b Filter) {
		if set {
			out = append(out, b)
		}
	}
	add(f.Type != TypeUnset, Filter{Type: f.Type})
	add(f.Gigantic, Filter{Gigantic: true})
	add(f.Name != "", Filter{Name: f.Name})
	add(f.House.filters(), Filter{House: f.House})
	add(f.Trait != traitUnset, Filter{Trait: f.Trait})
	add(f.ExceptTrait != traitUnset, Filter{ExceptTrait: f.ExceptTrait})
	add(f.Power.filters(), Filter{Power: f.Power})
	add(f.Damage.filters(), Filter{Damage: f.Damage})
	add(f.Aember.filters(), Filter{Aember: f.Aember})
	add(f.Position.filters(), Filter{Position: f.Position})
	add(f.Keyword.valid(), Filter{Keyword: f.Keyword})
	add(f.Counter.valid(), Filter{Counter: f.Counter})
	add(f.Stunned, Filter{Stunned: true})
	add(f.Ready, Filter{Ready: true})
	add(f.Armor, Filter{Armor: true})
	add(f.Upgrade, Filter{Upgrade: true})
	add(f.NoBonusIcons, Filter{NoBonusIcons: true})
	add(f.Neighboring, Filter{Neighboring: true})
	add(f.SharesTrait, Filter{SharesTrait: true})
	add(f.WithoutSharedTrait, Filter{WithoutSharedTrait: true})
	add(f.HouseWithAtLeast > 0, Filter{HouseWithAtLeast: f.HouseWithAtLeast})
	add(f.SharesHouseWithNeighbors > 0,
		Filter{SharesHouseWithNeighbors: f.SharesHouseWithNeighbors})
	add(f.HouseWithMostCreatures, Filter{HouseWithMostCreatures: true})
	add(f.Except.filters(), Filter{Except: f.Except})
	return out
}

// matches reports whether one card passes the filter. A conjoining filter asks
// every axis and a card must pass all of them; a disjoining one (MatchAny) asks
// each branch and a card need pass only one.
func (f Filter) matches(ctx *EffectContext, id LocalID) bool {
	if f.MatchAny {
		for _, b := range f.branches() {
			if b.conjoined(ctx, id) {
				return true
			}
		}
		return false
	}
	return f.conjoined(ctx, id)
}

// conjoined reports whether one card passes every axis the filter sets. The axes
// are pure reads combined with AND, grouped into families so no single predicate
// carries them all; a card must pass every family to match.
func (f Filter) conjoined(ctx *EffectContext, id LocalID) bool {
	return f.matchesTypeTraitHouse(ctx, id) &&
		f.matchesPower(ctx, id) &&
		f.matchesState(ctx, id) &&
		f.matchesPosition(ctx, id) &&
		f.matchesIdentity(ctx, id)
}

// matchesTypeTraitHouse reports whether a card passes the identity axes that say
// what kind of card it is: its type, its traits, and its house.
func (f Filter) matchesTypeTraitHouse(ctx *EffectContext, id LocalID) bool {
	if f.Type != TypeUnset && ctx.Resolver.TypeOf(id) != f.Type {
		return false
	}
	if f.Gigantic && ctx.Resolver.GiganticRoleOf(id) == GiganticNone {
		return false
	}
	if f.Trait != traitUnset && !ctx.Resolver.HasTrait(id, f.Trait) {
		return false
	}
	if !f.House.matches(ctx, id) {
		return false
	}
	if f.ExceptTrait != traitUnset && ctx.Resolver.HasTrait(id, f.ExceptTrait) {
		return false
	}
	if f.HouseWithMostCreatures && !isOfMostPopulousHouse(ctx, id) {
		return false
	}
	if f.SharesTrait && (!ctx.HasIt || !ctx.Resolver.SharesTrait(ctx.It, id)) {
		return false
	}
	return true
}

// matchesPower reports whether a card's power passes the power bound — a maximum,
// minimum, exact, odd, or even power requirement. A filter that bounds no power
// reads no power, so a card without one is never asked.
func (f Filter) matchesPower(ctx *EffectContext, id LocalID) bool {
	if !f.Power.filters() {
		return true
	}
	return f.Power.admits(ctx, ctx.Resolver.Power(id))
}

// matchesState reports whether a card passes the per-card state axes: damage,
// Æmber, bonus icons, counters, armor, upgrades, house-sharing neighbors, a
// keyword, and the stunned/ready flags.
func (f Filter) matchesState(ctx *EffectContext, id LocalID) bool {
	if f.Damage.filters() && !f.Damage.admits(ctx.Resolver.Damage(id)) {
		return false
	}
	if f.Aember.filters() && !f.Aember.admits(ctx.Resolver.AmberOn(id)) {
		return false
	}
	if f.NoBonusIcons && ctx.Resolver.HasBonusIcons(id) {
		return false
	}
	if f.Counter.valid() && ctx.Resolver.CountersOn(id, f.Counter) == 0 {
		return false
	}
	if f.Armor && ctx.Resolver.Armor(id) == 0 {
		return false
	}
	if f.Upgrade && len(ctx.Resolver.Upgrades(id)) == 0 {
		return false
	}
	if f.SharesHouseWithNeighbors > 0 &&
		sharedHouseNeighbors(ctx, id) < f.SharesHouseWithNeighbors {
		return false
	}
	if f.Keyword.valid() && !ctx.Resolver.HasKeyword(id, f.Keyword) {
		return false
	}
	if f.Stunned && !ctx.Resolver.Stunned(id) {
		return false
	}
	if f.Ready && ctx.Resolver.Exhausted(id) {
		return false
	}
	if f.HouseWithAtLeast > 0 && housePopulation(ctx, id) < f.HouseWithAtLeast {
		return false
	}
	if f.WithoutSharedTrait && sharesTraitWithBattlelineMate(ctx, id) {
		return false
	}
	return true
}

// matchesPosition reports whether a card passes the battleline position axes — on
// or off a flank, in the center, neighboring the source, or to the source's right
// or left.
func (f Filter) matchesPosition(ctx *EffectContext, id LocalID) bool {
	if f.Position.filters() && !f.Position.admits(ctx, id) {
		return false
	}
	if f.Neighboring && !isNeighbor(ctx, ctx.Source, id) {
		return false
	}
	return true
}

// matchesIdentity reports whether a card passes the identity axes: an exclusion
// drops one card the filter is defined against, and a name axis keeps only a
// named card.
func (f Filter) matchesIdentity(ctx *EffectContext, id LocalID) bool {
	if !f.Except.admits(ctx, id) {
		return false
	}
	if f.Name != "" && ctx.Resolver.Name(id) != f.Name {
		return false
	}
	return true
}

// qualifyNoun prefixes the adjectives the filter contributes to a base noun, in
// printed order, and substitutes the printed name for the noun when the filter
// names one — "each creature" becomes "each stunned damaged Mars creature". It is
// the first half of the rendering; the consumer then writes its quantifier and
// calls clauses for the rest.
func (f Filter) qualifyNoun(base string) string {
	if f.MatchAny {
		return f.disjunction(base)
	}
	return f.adjectives(base)
}

// disjunction renders a MatchAny filter as one complete noun phrase per branch,
// joined by "or". Branches that are nothing but an adjective on the shared base
// noun fold into one phrase — "Mars creature" and "Robot creature" print as "Mars
// or Robot creature" (EMP Blast), not as two nouns — because a card that names one
// set of creatures should read as one noun. A branch that carries a clause or a
// noun of its own cannot fold, so it prints in full ("an upgrade or Robot card",
// Chief Engineer Walls).
//
// The whole phrase is built here, clauses and all, so clauses adds nothing
// afterwards: a disjunction cannot be split around a quantifier the way a
// conjoining filter's adjectives and trailing clauses are.
func (f Filter) disjunction(base string) string {
	branches := f.branches()
	if len(branches) == 0 {
		return base
	}
	parts := make([]string, len(branches))
	for i, b := range branches {
		parts[i] = b.trailingClauses(b.adjectives(base))
	}
	if folded, ok := foldAdjectivesOn(parts, base); ok {
		return folded
	}
	return strings.Join(parts, " or ")
}

// foldAdjectivesOn folds phrases that are each one or more adjectives on the same
// base noun into a single phrase joining the adjectives with "or", and reports
// whether every phrase had that shape.
func foldAdjectivesOn(parts []string, base string) (string, bool) {
	adjectives := make([]string, len(parts))
	for i, p := range parts {
		if !strings.HasSuffix(p, " "+base) {
			return "", false
		}
		adjectives[i] = strings.TrimSuffix(p, " "+base)
	}
	return strings.Join(adjectives, " or ") + " " + base, true
}

// adjectives prefixes the adjectives a conjoining filter contributes to a base
// noun, in printed order, substituting the printed name and the card type for the
// noun when the filter names one.
//
// An empty base means the filter's own words are the whole noun, which is how a
// consumer says a trait names a kind of card by itself: the Shards count "each
// friendly Shard", not "each friendly Shard card" (Shard of Life). A consumer
// that passes no base must set an axis that supplies one.
func (f Filter) adjectives(base string) string {
	noun := base
	if f.Type != TypeUnset {
		noun = typeWord(f.Type)
	}
	if f.Gigantic {
		noun = "gigantic creature"
	}
	if f.Name != "" {
		noun = f.Name
	}
	if f.Trait != traitUnset {
		noun = qualifyNoun(f.Trait.String(), noun)
	}
	noun = f.House.qualifyNoun(noun)
	if f.ExceptTrait != traitUnset {
		noun = "non-" + f.ExceptTrait.String() + " " + noun
	}
	if adj := f.Position.adjective(); adj != "" {
		// A flank is a battleline position, so on a noun that also reaches artifacts
		// the qualifier binds to the creature half alone — which the printed phrase
		// says by naming the artifact first (Snudge).
		if strings.Contains(noun, creatureOrArtifact) {
			noun = strings.Replace(
				noun, creatureOrArtifact, "artifact or "+adj+" creature", 1)
		} else {
			noun = adj + " " + noun
		}
	}
	if f.Neighboring {
		noun = "neighboring " + noun
	}
	if adj := f.Damage.adjective(); adj != "" {
		noun = adj + " " + noun
	}
	if f.Stunned {
		noun = "stunned " + noun
	}
	if f.Ready {
		noun = "ready " + noun
	}
	if f.Keyword.valid() {
		noun = strings.ToLower(f.Keyword.String()) + " " + noun
	}
	return noun
}

// clauses appends the trailing clauses the filter contributes to an already
// quantified phrase, in printed order — "each creature" becomes "each creature
// with power 3 or lower on a flank". It is the second half of the rendering,
// after the consumer has written its quantifier onto the noun qualifyNoun built.
func (f Filter) clauses(phrase string) string {
	if f.MatchAny {
		return phrase
	}
	return f.trailingClauses(phrase)
}

// trailingClauses appends the trailing clauses a conjoining filter contributes,
// in printed order.
func (f Filter) trailingClauses(phrase string) string {
	phrase = f.Power.clause(phrase)
	phrase = f.Aember.clause(phrase)
	if f.NoBonusIcons {
		phrase += " with no bonus icons"
	}
	if f.Counter.valid() {
		phrase += " with a " + f.Counter.noun()
	}
	if f.Armor {
		phrase += " with armor"
	}
	if f.Upgrade {
		phrase += " with an upgrade"
	}
	if f.SharesHouseWithNeighbors == 1 {
		phrase += " that shares a house with at least 1 of its neighbors"
	} else if f.SharesHouseWithNeighbors > 1 {
		phrase += fmt.Sprintf(
			" that shares a house with %d of its neighbors",
			f.SharesHouseWithNeighbors,
		)
	}
	phrase = f.Position.clause(phrase)
	phrase = f.House.qualifyPhrase(phrase)
	if f.HouseWithMostCreatures {
		phrase += " of the house with the most creatures in play"
	}
	if f.SharesTrait {
		phrase += " that shares a trait with it"
	}
	if f.WithoutSharedTrait {
		phrase += " that does not share a trait with another creature in its controller's battleline"
	}
	if f.HouseWithAtLeast > 0 {
		phrase += fmt.Sprintf(
			" that belongs to a house that has %d or more creatures in play",
			f.HouseWithAtLeast,
		)
	}
	return phrase
}

// noun renders the whole filter as one noun phrase with no article — "Mars
// creature", "upgrade or Robot card", "creature with power 3 or lower" — over the
// base noun a consumer prints when the filter names no type of its own. It is the
// unsplit rendering, for a consumer that prints its subject in one piece; a Target
// instead calls qualifyNoun and clauses around its own quantifier.
func (f Filter) noun(base string) string {
	return f.clauses(f.qualifyNoun(base))
}

// article puts the indefinite article on a phrase this filter rendered, unless
// the filter names a card outright: a proper name identifies one specific card and
// so takes no article — "Subtle Chain", never "a Subtle Chain". It is the one
// article rule every consumer that names a single card follows.
func (f Filter) article(phrase string) string {
	if f.Name != "" {
		return phrase
	}
	return indefinite(phrase)
}

// object renders the filter as a single card with its article — "a Mars
// creature", "an upgrade", "Subtle Chain" — over the base noun a consumer prints
// when the filter names no type of its own.
func (f Filter) object(base string) string {
	return f.article(f.noun(base))
}

// readsPlay reports whether the filter sets an axis that only means something for
// a card in play — its power, its damage, the Æmber and counters on it, its armor,
// its upgrades, its stun or ready state, its keywords, or where it stands in a
// battleline. The identity axes (type, house, trait, name) are true of a card
// wherever it sits and are always legal.
func (f Filter) readsPlay() bool {
	return f.Power.filters() || f.Damage.filters() || f.Aember.filters() ||
		f.Position.filters() || f.Keyword.valid() || f.Counter.valid() ||
		f.Stunned || f.Ready || f.Armor || f.Upgrade || f.NoBonusIcons ||
		f.Neighboring || f.SharesTrait || f.SharesHouseWithNeighbors > 0 ||
		f.HouseWithMostCreatures || f.Except.filters() ||
		f.HouseWithAtLeast > 0 || f.WithoutSharedTrait
}

// validateIdentityOnly rejects a filter that narrows on an in-play axis, for a
// consumer whose cards sit in a deck or a pile rather than on the board. Such an
// axis has nothing to read there, so it would silently match nothing; naming it is
// a definition error, caught at build time the way an unset source zone is.
func (f Filter) validateIdentityOnly(effect string) error {
	if f.readsPlay() {
		return fmt.Errorf(
			"%s: Filter narrows on an in-play axis, which a card outside play has none of",
			effect,
		)
	}
	return nil
}

// refine keeps the ids the filter admits, which is how a Filter serves as a
// Refinement inside a union combinator (Regrettable Meteor destroys each Dinosaur
// creature and each creature with power 6 or higher as one set). Every per-card
// test is trivially a set rule; the reverse is false, which is why the merge runs
// this way only.
func (f Filter) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if f.matches(ctx, id) {
			kept = append(kept, id)
		}
	}
	return kept
}

// clause narrows an already quantified phrase, the Refinement half of the
// rendering: the adjectives go before the phrase's last noun ("each creature"
// becomes "each Dinosaur creature") and the trailing clauses after the whole
// phrase ("each creature with power 6 or higher").
func (f Filter) clause(phrase string) string {
	head, noun := "", phrase
	if i := strings.LastIndex(phrase, " "); i >= 0 {
		head, noun = phrase[:i+1], phrase[i+1:]
	}
	return f.clauses(head + f.qualifyNoun(noun))
}

// housePopulation counts the creatures of id's house in play, across both
// players' battlelines — a house is a house regardless of who controls it.
func housePopulation(ctx *EffectContext, id LocalID) int {
	house := ctx.Resolver.House(id)
	n := 0
	for player := range 2 {
		for _, cid := range ctx.Resolver.Battleline(player) {
			if ctx.Resolver.House(cid) == house {
				n++
			}
		}
	}
	return n
}

// sharesTraitWithBattlelineMate reports whether another creature in id's own
// controller's battleline shares a trait with it. An enemy sharing a trait does
// not count: the comparison is with its own side alone.
func sharesTraitWithBattlelineMate(ctx *EffectContext, id LocalID) bool {
	for _, mate := range ctx.Resolver.Battleline(ctx.Resolver.Controller(id)) {
		if mate != id && ctx.Resolver.SharesTrait(id, mate) {
			return true
		}
	}
	return false
}

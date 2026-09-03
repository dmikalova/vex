package engine

import (
	"fmt"
	"sort"
	"strings"
)

// A Target names the cards an ability acts on. KeyForge abilities are written in
// terms of noun phrases — "this creature", "each enemy creature", "a friendly
// creature", "each Scientist trait creature", "each creature with power 3 or
// lower" — and Target captures exactly that: a base set chosen by Kind, narrowed
// by optional filters. Because the same value both renders the phrase (Text) and
// selects the cards (Select), one effect such as Destroy can express many
// different printed cards just by changing its Target.

// TargetKind enumerates the base sets a Target can select before filtering.
type TargetKind int

const (
	// targetUnset is the invalid zero value — a Target whose base set was never
	// chosen. An effect that requires a target rejects it in validation, so a card
	// must always name its target explicitly rather than leaning on a default.
	targetUnset TargetKind = iota
	// TargetThisCreature selects the source card itself.
	TargetThisCreature
	// TargetTriggeringCreature selects the creature that caused the trigger ("it").
	TargetTriggeringCreature
	// TargetCreatureFought selects the creature the source is fighting, named in
	// full ("the creature <self> fights") so a Before Fight ability that reaches
	// past it — Lord Golgotha damaging its neighbors — reads unambiguously.
	TargetCreatureFought
	// TargetEachCreature selects every creature in play.
	TargetEachCreature
	// TargetEachFriendlyCreature selects every creature the controller controls.
	TargetEachFriendlyCreature
	// TargetEachEnemyCreature selects every creature the opponent controls.
	TargetEachEnemyCreature
	// TargetEachArtifact selects every artifact in play, both players'.
	TargetEachArtifact
	// TargetEachFriendlyArtifact selects every artifact the controller controls.
	TargetEachFriendlyArtifact
	// TargetEachEnemyArtifact selects every artifact the opponent controls.
	TargetEachEnemyArtifact
	// TargetEachCardInPlay selects every card in play — every creature and artifact,
	// both players' — including the source card itself.
	TargetEachCardInPlay
	// TargetEachFriendlyCardInPlay selects the controller's cards in play — their
	// creatures and artifacts.
	TargetEachFriendlyCardInPlay
	// TargetEachOtherFriendlyCreature selects the controller's creatures except
	// the source card.
	TargetEachOtherFriendlyCreature
	// TargetChosenCreature selects a single creature the controller chooses from
	// all creatures in play (either player's).
	TargetChosenCreature
	// TargetChosenEnemyCreature selects a single enemy creature the controller
	// chooses.
	TargetChosenEnemyCreature
	// TargetChosenFriendlyCreature selects a single friendly creature the
	// controller chooses.
	TargetChosenFriendlyCreature
	// TargetChosenOtherFriendlyCreature selects a single friendly creature the
	// controller chooses, excluding the source card ("another friendly creature").
	TargetChosenOtherFriendlyCreature
	// TargetChosenOtherCreature selects a single creature the controller chooses
	// from all in play, excluding the creature in context (ctx.It) — "another
	// creature" than the one a preceding effect put in focus (Guardian Demon deals
	// to another creature than the one it healed).
	TargetChosenOtherCreature
	// TargetChosenArtifact selects a single artifact the controller chooses from
	// all artifacts in play (either player's).
	TargetChosenArtifact
	// TargetChosenEnemyArtifact selects a single enemy artifact the controller
	// chooses (Sneklifter seizes one).
	TargetChosenEnemyArtifact
	// TargetTheOtherCreature selects the creature in context (ctx.It) — the one a
	// preceding effect chose as "another" creature — and renders it as "the other
	// creature". Transposition Sandals swaps with another creature, then uses the
	// other creature.
	TargetTheOtherCreature
	// TargetChosenCreatureOrArtifact selects a single creature or artifact the controller
	// chooses from all in play (either player's), rendered "a creature or artifact".
	TargetChosenCreatureOrArtifact
	// TargetChosenFriendlyCreatureOrArtifact selects a single friendly creature or artifact the
	// controller chooses, rendered "a friendly creature or artifact".
	TargetChosenFriendlyCreatureOrArtifact
)

// Target describes which cards an effect applies to. Kind picks the base set;
// the optional filters added by WithTrait and PowerAtMost narrow that set and
// extend the rendered text.
type Target struct {
	Kind        TargetKind
	trait       Trait
	exceptTrait Trait
	house       House
	exceptHouse House
	chosenHouse bool
	// contextualHouse narrows the target to cards sharing the house of the card in
	// context (ctx.It), rendering "of that card's house" — ForEachDiscarded's Do.
	contextualHouse bool
	// sharesTrait narrows the target to cards sharing at least one trait with the
	// card in context (ctx.It), rendering "that shares a trait with it".
	sharesTrait   bool
	maxPower      int
	hasMaxPower   bool
	minPower      int
	hasMinPower   bool
	exactPower    int
	hasExactPower bool
	// variablePower renders an exact-power qualifier as the placeholder "X" — a
	// template face whose concrete variants each fix the number (Master of X).
	variablePower bool
	damaged       bool
	undamaged     bool
	stunned       bool
	withAember    bool
	// withArmor narrows the target to creatures that have armor at all, rendering
	// " with armor". It reads the creature's armor value, not what is left of it, so
	// a creature that has already spent its armor absorbing damage still has armor.
	withArmor   bool
	keyword     Keyword
	onFlank     bool
	notOnFlank  bool
	neighboring bool
	// withNeighbors expands a single chosen creature to include its battleline
	// neighbors (Tremor stuns a creature and each of its neighbors).
	withNeighbors bool
	// neighborsOf narrows the selection to the battleline neighbors of what it
	// selects, dropping the selected creature itself.
	neighborsOf bool
	// other excludes the source card from the selected set ("other" cards).
	other bool
	// named narrows the target to cards with this printed name, and replaces the
	// rendered noun with it: a card that names another card outright says "an
	// Ancient Bear", not "an Ancient Bear creature".
	named string
	// selector is a set-relative refinement applied after the per-card filters. It
	// can compare the candidates to each other (e.g. "except the most powerful")
	// and contributes a clause to the printed phrase. nil for targets that select
	// their whole filtered set.
	selector Selector
}

// WithTrait narrows the target to cards that have the given trait, e.g.
// Target{Kind: TargetEachCreature}.WithTrait("Scientist").
func (t Target) WithTrait(trait Trait) Target {
	t.trait = trait
	return t
}

// ExceptTrait narrows the target to cards that do NOT have the given trait,
// rendering the "non-<trait> trait" qualifier, e.g. a friendly Mars creature
// ExceptTrait("Agent") reads "a friendly non-Agent trait Mars creature".
func (t Target) ExceptTrait(trait Trait) Target {
	t.exceptTrait = trait
	return t
}

// OfHouse narrows the target to cards of the given house, e.g.
// Target{Kind: TargetEachCreature}.OfHouse(Mars).
func (t Target) OfHouse(h House) Target {
	t.house = h
	return t
}

// ExceptHouse narrows the target to cards NOT of the given house, rendering the
// "non-<house>" qualifier, e.g. a chosen friendly creature ExceptHouse(Sanctum)
// reads "a friendly non-Sanctum creature".
func (t Target) ExceptHouse(h House) Target {
	t.exceptHouse = h
	return t
}

// OfChosenHouse narrows the target to cards of the house picked by an enclosing
// ChooseHouseThen (read from the effect context at selection time).
func (t Target) OfChosenHouse() Target {
	t.chosenHouse = true
	return t
}

// OfContextualHouse narrows the target to cards sharing the house of the card in
// context (ctx.It) — the discarded card ForEachDiscarded puts in focus — rendering
// "of that card's house".
func (t Target) OfContextualHouse() Target {
	t.contextualHouse = true
	return t
}

// selfHouseResolved fills the card's own house in for any SelfHouse sentinel the
// target narrows on. A Target keeps its houses and selector unexported, so it
// resolves itself rather than being rewritten by reflection (see self_house.go).
func (t Target) selfHouseResolved(house House) any {
	if t.house == SelfHouse {
		t.house = house
	}
	if t.exceptHouse == SelfHouse {
		t.exceptHouse = house
	}
	if t.selector != nil {
		t.selector = resolvedIn(t.selector, house)
	}
	return t
}

// SharingTrait narrows the target to cards that share at least one trait with the
// card in context (ctx.It), rendering "that shares a trait with it" — the purged
// creature after PurgeCreatureFromHand (Custom Virus), or whatever an earlier
// effect put in context.
func (t Target) SharingTrait() Target {
	t.sharesTrait = true
	return t
}

// PowerAtMost narrows the target to creatures whose power is maxPower or lower,
// e.g. Target{Kind: TargetEachCreature}.PowerAtMost(3).
func (t Target) PowerAtMost(maxPower int) Target {
	t.maxPower = maxPower
	t.hasMaxPower = true
	return t
}

// PowerAtLeast narrows the target to creatures whose power is minPower or higher,
// e.g. Target{Kind: TargetEachCreature}.PowerAtLeast(3).
func (t Target) PowerAtLeast(minPower int) Target {
	t.minPower = minPower
	t.hasMinPower = true
	return t
}

// PowerExactly narrows the target to creatures whose power is exactly power,
// e.g. Target{Kind: TargetChosenCreature}.PowerExactly(1).
func (t Target) PowerExactly(power int) Target {
	t.exactPower = power
	t.hasExactPower = true
	return t
}

// PowerVariable narrows the target to a power the card fills in per instance,
// rendering the placeholder "X" (the Master of X template face) until a concrete
// variant sets the number with PowerExactly.
func (t Target) PowerVariable() Target {
	t.hasExactPower = true
	t.variablePower = true
	return t
}

// Damaged narrows the target to creatures that currently have damage on them.
func (t Target) Damaged() Target {
	t.damaged = true
	return t
}

// Undamaged narrows the target to creatures that currently have no damage on them.
func (t Target) Undamaged() Target {
	t.undamaged = true
	return t
}

// Named narrows the target to cards with the given printed name, e.g.
// Target{Kind: TargetChosenCreature}.Named("Ancient Bear").
func (t Target) Named(name string) Target {
	t.named = name
	return t
}

// WithAember narrows the target to creatures that have Æmber on them, rendering
// " with Æmber on it", e.g. "each creature with Æmber on it".
func (t Target) WithAember() Target {
	t.withAember = true
	return t
}

// WithArmor narrows the target to creatures that have armor, rendering " with
// armor", e.g. "each enemy creature with armor".
func (t Target) WithArmor() Target {
	t.withArmor = true
	return t
}

// Keyword narrows the target to creatures that have the given keyword (e.g.
// Elusive), rendering it as an adjective: "each elusive creature".
func (t Target) Keyword(k Keyword) Target {
	t.keyword = k
	return t
}

// Stunned narrows the target to creatures that are currently stunned.
func (t Target) Stunned() Target {
	t.stunned = true
	return t
}

// allows reports whether a single card satisfies the target's per-card filters,
// ignoring its base-set Kind. It is how a Target expresses a condition on one
// specific card (e.g. a fight restriction testing the defender).
func (t Target) allows(ctx *EffectContext, id LocalID) bool {
	return len(t.filter(ctx, []LocalID{id})) == 1
}

// OnFlank narrows the target to creatures on a flank of their battleline (its
// leftmost or rightmost creature). A flank is a battleline position, so the
// filter only constrains creatures: on a target that also reaches artifacts
// ("an artifact or flank creature", Snudge) an artifact passes it untouched.
func (t Target) OnFlank() Target {
	t.onFlank = true
	return t
}

// NotOnFlank narrows the target to creatures that are not on a flank of their
// battleline (neither its leftmost nor rightmost creature).
func (t Target) NotOnFlank() Target {
	t.notOnFlank = true
	return t
}

// Neighboring narrows the target to the source card's battleline neighbors (the
// creatures immediately to its left and right).
func (t Target) Neighboring() Target {
	t.neighboring = true
	return t
}

// AndNeighbors expands a single chosen creature to also include its battleline
// neighbors, so an effect applies to the chosen creature and each of its
// neighbors (Tremor). It is meaningful only on a chosen-creature target.
func (t Target) AndNeighbors() Target {
	t.withNeighbors = true
	return t
}

// NeighborsOf narrows a target to the battleline neighbors of what it selects,
// dropping the selected creature itself — Lord Golgotha damages each neighbor of
// the creature it fights, but not that creature.
func (t Target) NeighborsOf() Target {
	t.neighborsOf = true
	return t
}

// Other excludes the source card from the selected set, rendering the "other"
// qualifier ("each other friendly card").
func (t Target) Other() Target {
	t.other = true
	return t
}

// Selector refines the target with a set-relative rule applied after the per-card
// filters, e.g. Target{...}.Selector(ExceptMostPowerful). The Selector both picks
// the final subset and describes itself for the printed phrase, so a niche
// "relative to the rest of the set" rule composes onto any Target without adding
// a dedicated field (and future rules — least powerful, and so on — are just more
// Selector values).
func (t Target) Selector(s Selector) Target {
	t.selector = s
	return t
}

// valid reports whether the target's base set was chosen (its Kind is not the
// unset zero value). Effects that require a target check this in validation.
func (t Target) valid() bool {
	return t.Kind != targetUnset
}

// Text renders the target as an English noun phrase, e.g. "each enemy creature",
// "each Scientist trait creature", or "each creature with power 3 or lower".
func (t Target) Text() string {
	switch t.Kind {
	case TargetThisCreature:
		return SelfName
	case TargetTriggeringCreature:
		return t.decorateNeighbors("it")
	case TargetCreatureFought:
		return t.decorateNeighbors("the creature " + SelfName + " fights")
	case TargetTheOtherCreature:
		return "the other creature"
	}
	noun := "creature"
	if t.Kind == TargetEachArtifact || t.Kind == TargetChosenArtifact ||
		t.Kind == TargetChosenEnemyArtifact || t.Kind == TargetEachFriendlyArtifact ||
		t.Kind == TargetEachEnemyArtifact {
		noun = "artifact"
	}
	if t.Kind == TargetEachFriendlyCardInPlay {
		noun = "card"
	}
	orArtifact := t.Kind == TargetChosenCreatureOrArtifact ||
		t.Kind == TargetChosenFriendlyCreatureOrArtifact
	if orArtifact {
		noun = "creature or artifact"
	}
	if t.named != "" {
		noun = t.named
	}
	if t.exceptHouse != HouseNone {
		noun = "non-" + t.exceptHouse.String() + " " + noun
	}
	if t.trait != "" {
		noun = string(t.trait) + " trait " + noun
	}
	if t.house != HouseNone {
		noun = t.house.String() + " " + noun
	}
	if t.exceptTrait != "" {
		noun = "non-" + string(t.exceptTrait) + " trait " + noun
	}
	if t.onFlank {
		// A flank is a battleline position, so on a target that also reaches artifacts
		// the qualifier binds to the creature half alone — which the printed phrase
		// says by naming the artifact first (Snudge).
		if orArtifact {
			noun = strings.Replace(noun, "creature or artifact", "artifact or flank creature", 1)
		} else {
			noun = "flank " + noun
		}
	}
	if t.neighboring {
		noun = "neighboring " + noun
	}
	if t.damaged {
		noun = "damaged " + noun
	}
	if t.undamaged {
		noun = "undamaged " + noun
	}
	if t.stunned {
		noun = "stunned " + noun
	}
	if t.keyword.valid() {
		noun = strings.ToLower(t.keyword.String()) + " " + noun
	}
	var phrase string
	switch t.Kind {
	case TargetEachCardInPlay:
		phrase = "each card in play"
	case TargetEachCreature, TargetEachArtifact:
		if t.other {
			phrase = "each other " + noun
		} else {
			phrase = "each " + noun
		}
	case TargetEachFriendlyCreature:
		phrase = "each friendly " + noun
	case TargetEachFriendlyArtifact:
		phrase = "each friendly " + noun
	case TargetEachEnemyArtifact:
		phrase = "each enemy " + noun
	case TargetEachFriendlyCardInPlay:
		if t.other {
			phrase = "each other friendly " + noun
		} else {
			phrase = "each friendly " + noun
		}
	case TargetEachEnemyCreature:
		phrase = "each enemy " + noun
	case TargetEachOtherFriendlyCreature:
		phrase = "each other friendly " + noun
	case TargetChosenEnemyCreature:
		phrase = "an enemy " + noun
	case TargetChosenFriendlyCreature, TargetChosenFriendlyCreatureOrArtifact:
		phrase = "a friendly " + noun
	case TargetChosenOtherFriendlyCreature:
		phrase = "another friendly " + noun
	case TargetChosenOtherCreature:
		phrase = "another " + noun
	case TargetChosenArtifact:
		phrase = "an " + noun
	case TargetChosenEnemyArtifact:
		phrase = "an enemy " + noun
	default:
		phrase = indefinite(noun)
	}
	if t.hasMaxPower {
		phrase += fmt.Sprintf(" with power %d or lower", t.maxPower)
	}
	if t.hasMinPower {
		phrase += fmt.Sprintf(" with power %d or higher", t.minPower)
	}
	if t.hasExactPower {
		if t.variablePower {
			phrase += " with power X"
		} else {
			phrase += fmt.Sprintf(" with power %d", t.exactPower)
		}
	}
	if t.withAember {
		phrase += " with \u00c6mber on it"
	}
	if t.withArmor {
		phrase += " with armor"
	}
	if t.notOnFlank {
		phrase += " that is not on a flank"
	}
	if t.chosenHouse {
		phrase += " of the chosen house"
	}
	if t.contextualHouse {
		phrase += " of that card's house"
	}
	if t.sharesTrait {
		phrase += " that shares a trait with it"
	}
	if t.selector != nil {
		phrase = t.selector.clause(phrase)
	}
	return t.decorateNeighbors(phrase)
}

// decorateNeighbors wraps a rendered noun phrase with the neighbour builders:
// AndNeighbors reads "<phrase> and each of its neighbors", NeighborsOf reads
// "each neighbor of <phrase>".
func (t Target) decorateNeighbors(phrase string) string {
	if t.withNeighbors {
		phrase += " and each of its neighbors"
	}
	if t.neighborsOf {
		phrase = "each neighbor of " + phrase
	}
	return phrase
}

// Select resolves the target into concrete card ids, applying its filters. For a
// chosen kind it asks the controller to pick one of the filtered candidates
// (returning nil when there are none or the choice is declined).
func (t Target) Select(ctx *EffectContext) []LocalID {
	return t.selectWith(ctx, false, nil)
}

// SelectOptional is Select inside a "you may": a chosen target is asked
// declinably, so the controller clicks the card they mean or passes, instead of
// answering a Yes/No and then being handed a pick they can no longer refuse. A
// target that chooses nothing has no decision to decline and behaves like Select.
func (t Target) SelectOptional(ctx *EffectContext) []LocalID {
	return t.selectWith(ctx, true, nil)
}

// empty reports that nothing matches this target, reading only the candidates a
// selector would narrow: with nothing to narrow there is nothing to select, and
// unlike Select it asks the controller nothing.
func (t Target) empty(ctx *EffectContext) bool {
	return len(t.filter(ctx, t.selectBase(ctx))) == 0
}

// selectWith is the shared selection path; optional switches the chosen-kind
// prompt between a forced pick and a declinable one, and keep (when set) drops
// candidates the calling effect could not act on.
func (t Target) selectWith(
	ctx *EffectContext,
	optional bool,
	keep func(LocalID) bool,
) []LocalID {
	ids := t.filter(ctx, t.selectBase(ctx))
	if keep != nil {
		kept := ids[:0:0]
		for _, id := range ids {
			if keep(id) {
				kept = append(kept, id)
			}
		}
		ids = kept
	}
	if t.selector != nil {
		ids = t.selector.refine(ctx, ids)
	}
	if !t.isChosen() {
		return t.expandNeighbors(ctx, ids)
	}
	if len(ids) == 0 {
		return nil
	}
	prompt := "Choose " + t.Text()
	var id LocalID
	var ok bool
	if optional {
		id, ok = ctx.ChooseCardOptional(prompt, ids)
	} else {
		id, ok = ctx.ChooseCreature(prompt, ids)
	}
	if !ok {
		return nil
	}
	return t.expandNeighbors(ctx, []LocalID{id})
}

// expandNeighbors applies the neighbour builders to an already-selected set:
// AndNeighbors keeps each selected creature and adds its battleline neighbors,
// NeighborsOf replaces the selection with them (Lord Golgotha hits the neighbors
// of the creature it fights, not that creature).
func (t Target) expandNeighbors(ctx *EffectContext, ids []LocalID) []LocalID {
	if !t.withNeighbors && !t.neighborsOf {
		return ids
	}
	out := ids[:0:0]
	for _, id := range ids {
		if t.withNeighbors {
			out = append(out, id)
		}
		out = append(out, neighbors(ctx, id)...)
	}
	return out
}

// A Selector refines a Target's selected set relative to the whole set — a rule
// that compares the candidates to each other rather than testing each on its own,
// such as "except the most powerful creature". It both narrows the ids (refine)
// and contributes a clause to the target's printed phrase (clause), so these
// niche rules compose onto any Target without a field per rule.
type Selector interface {
	refine(ctx *EffectContext, ids []LocalID) []LocalID
	clause(phrase string) string
}

// leadingSelector is the optional capability of a Selector whose choice reads
// before the effect's verb, so the phrase runs left to right — the choice stated
// first, then its consequence (SamePowerAsChosen's "choose a creature - destroy
// …"). A Selector without it contributes only a trailing clause.
type leadingSelector interface {
	lead() string
}

// leadIn returns the leading clause a Selector contributes ahead of the effect's
// verb, and whether the target has one — so an effect can render "choose a
// creature - destroy …" left to right instead of burying the choice.
func (t Target) leadIn() (string, bool) {
	if l, ok := t.selector.(leadingSelector); ok {
		return l.lead(), true
	}
	return "", false
}

// ExceptMostPowerful is a Selector that drops the single most powerful creature
// from the set, letting the controller choose which one to keep when several tie
// for most powerful. A set of one or none keeps its (only, most powerful) member,
// so nothing is selected.
var ExceptMostPowerful Selector = exceptMostPowerful{}

// exceptMostPowerful implements the ExceptMostPowerful selector.
type exceptMostPowerful struct{}

// clause renders "<phrase> except the most powerful <noun>", e.g. "each enemy
// creature except the most powerful enemy creature".
func (exceptMostPowerful) clause(phrase string) string {
	return phrase + " except the most powerful " + strings.TrimPrefix(phrase, "each ")
}

// refine returns ids without the single most powerful creature, letting the
// controller choose which to keep when several tie. A set of one or none keeps
// its (most powerful) member, so nothing is selected.
func (exceptMostPowerful) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) <= 1 {
		return nil
	}
	highest := ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		if p := ctx.Resolver.Power(id); p > highest {
			highest = p
		}
	}
	mostPowerful := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == highest {
			mostPowerful = append(mostPowerful, id)
		}
	}
	spared := mostPowerful[0]
	if len(mostPowerful) > 1 {
		if chosen, ok := ctx.ChooseCreature(
			"Choose the most powerful creature to keep",
			mostPowerful,
		); ok {
			spared = chosen
		}
	}
	out := make([]LocalID, 0, len(ids)-1)
	for _, id := range ids {
		if id != spared {
			out = append(out, id)
		}
	}
	return out
}

// SamePowerAsChosen is a Selector that keeps every creature sharing the power of
// one the controller chooses from the set — the chosen creature included — so a
// Destroy paired with it wipes out a whole power bracket (Dance of Doom). A
// declined choice selects nothing. It leads with "choose a creature" so the
// printed phrase reads left to right, the choice before its consequence.
var SamePowerAsChosen Selector = samePowerAsChosen{}

// samePowerAsChosen implements the SamePowerAsChosen selector.
type samePowerAsChosen struct{}

// lead renders the choice ahead of the effect's verb.
func (samePowerAsChosen) lead() string { return "choose a creature" }

// clause renders "<phrase> with the same power as the chosen creature".
func (samePowerAsChosen) clause(phrase string) string {
	return phrase + " with the same power as the chosen creature"
}

// refine picks a creature and keeps every one in the set with matching power.
func (samePowerAsChosen) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	chosen, ok := ctx.ChooseCreature("Choose a creature", ids)
	if !ok {
		return nil
	}
	power := ctx.Resolver.Power(chosen)
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == power {
			out = append(out, id)
		}
	}
	return out
}

// LeastPowerful is a Selector that keeps only the single least powerful creature
// of a set, e.g. card.Target.EachCreature.Selector(card.LeastPowerful) (Horseman
// of Famine). When several tie for least powerful the controller chooses which.
var LeastPowerful Selector = leastPowerful{}

// leastPowerful implements the LeastPowerful selector.
type leastPowerful struct{}

// clause renders "the least powerful <noun>", e.g. "each creature" -> "the least
// powerful creature".
func (leastPowerful) clause(phrase string) string {
	return "the least powerful " + strings.TrimPrefix(phrase, "each ")
}

// refine returns the single least powerful creature, letting the controller
// choose which to keep when several tie. An empty set selects nothing.
func (leastPowerful) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) == 0 {
		return nil
	}
	lowest := ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		if p := ctx.Resolver.Power(id); p < lowest {
			lowest = p
		}
	}
	tied := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == lowest {
			tied = append(tied, id)
		}
	}
	pick := tied[0]
	if len(tied) > 1 {
		if chosen, ok := ctx.ChooseCreature("Choose the least powerful creature", tied); ok {
			pick = chosen
		}
	}
	return []LocalID{pick}
}

// MostPowerful returns a Selector that keeps the n most powerful creatures of a
// set, e.g. card.Target.EachCreature.Selector(card.MostPowerful(3)) (Three Fates).
// When more creatures tie at the cutoff than there are remaining slots, the
// controller chooses which of the tied creatures to include.
func MostPowerful(n int) Selector { return mostPowerfulN{n: n} }

// mostPowerfulN implements the MostPowerful selector.
type mostPowerfulN struct{ n int }

// clause renders "the N most powerful <noun>s", e.g. "the 3 most powerful creatures".
func (m mostPowerfulN) clause(phrase string) string {
	return fmt.Sprintf("the %d most powerful %ss", m.n, strings.TrimPrefix(phrase, "each "))
}

// refine keeps the n highest-power creatures, letting the controller break ties
// at the cutoff. A set no larger than n keeps all of it.
func (m mostPowerfulN) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) <= m.n {
		return ids
	}
	sorted := append([]LocalID(nil), ids...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return ctx.Resolver.Power(sorted[i]) > ctx.Resolver.Power(sorted[j])
	})
	threshold := ctx.Resolver.Power(sorted[m.n-1])
	chosen := make([]LocalID, 0, m.n)
	tied := make([]LocalID, 0)
	for _, id := range sorted {
		switch p := ctx.Resolver.Power(id); {
		case p > threshold:
			chosen = append(chosen, id)
		case p == threshold:
			tied = append(tied, id)
		}
	}
	for slots := m.n - len(chosen); slots > 0 && len(tied) > 0; slots-- {
		if len(tied) == slots {
			chosen = append(chosen, tied...)
			break
		}
		pick := tied[0]
		if c, ok := ctx.ChooseCreature("Choose one of the most powerful creatures", tied); ok {
			pick = c
		}
		chosen = append(chosen, pick)
		for i, id := range tied {
			if id == pick {
				tied = append(tied[:i], tied[i+1:]...)
				break
			}
		}
	}
	return chosen
}

// isChosen reports whether the Kind resolves to a single player-chosen creature.
func (t Target) isChosen() bool {
	return t.Kind == TargetChosenCreature || t.Kind == TargetChosenEnemyCreature ||
		t.Kind == TargetChosenFriendlyCreature || t.Kind == TargetChosenOtherFriendlyCreature ||
		t.Kind == TargetChosenOtherCreature ||
		t.Kind == TargetChosenArtifact || t.Kind == TargetChosenEnemyArtifact || t.Kind == TargetChosenCreatureOrArtifact ||
		t.Kind == TargetChosenFriendlyCreatureOrArtifact
}

// filter narrows ids to those matching the target's trait, power, damaged, and
// flank filters.
func (t Target) filter(ctx *EffectContext, ids []LocalID) []LocalID {
	if t.trait == "" &&
		t.exceptTrait == "" &&
		t.house == HouseNone &&
		t.exceptHouse == HouseNone &&
		!t.chosenHouse &&
		!t.contextualHouse &&
		!t.sharesTrait &&
		!t.hasMaxPower &&
		!t.hasMinPower &&
		!t.hasExactPower &&
		!t.damaged &&
		!t.undamaged &&
		!t.stunned &&
		!t.withAember &&
		!t.withArmor &&
		t.keyword == keywordUnset &&
		!t.onFlank &&
		!t.notOnFlank &&
		!t.neighboring &&
		!t.other &&
		t.named == "" {
		return ids
	}
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if t.trait != "" && !ctx.Resolver.HasTrait(id, t.trait) {
			continue
		}
		if t.exceptTrait != "" && ctx.Resolver.HasTrait(id, t.exceptTrait) {
			continue
		}
		if t.house != HouseNone && ctx.Resolver.House(id) != t.house {
			continue
		}
		if t.exceptHouse != HouseNone && ctx.Resolver.House(id) == t.exceptHouse {
			continue
		}
		if t.chosenHouse && ctx.Resolver.House(id) != ctx.ChosenHouse {
			continue
		}
		if t.contextualHouse &&
			(!ctx.HasIt || ctx.Resolver.House(id) != ctx.Resolver.House(ctx.It)) {
			continue
		}
		if t.sharesTrait && (!ctx.HasIt || !ctx.Resolver.SharesTrait(ctx.It, id)) {
			continue
		}
		if t.hasMaxPower && ctx.Resolver.Power(id) > t.maxPower {
			continue
		}
		if t.hasMinPower && ctx.Resolver.Power(id) < t.minPower {
			continue
		}
		if t.hasExactPower && ctx.Resolver.Power(id) != t.exactPower {
			continue
		}
		if t.damaged && ctx.Resolver.Damage(id) == 0 {
			continue
		}
		if t.undamaged && ctx.Resolver.Damage(id) != 0 {
			continue
		}
		if t.withAember && ctx.Resolver.AmberOn(id) == 0 {
			continue
		}
		if t.withArmor && ctx.Resolver.Armor(id) == 0 {
			continue
		}
		if t.keyword.valid() && !ctx.Resolver.HasKeyword(id, t.keyword) {
			continue
		}
		if t.stunned && !ctx.Resolver.Stunned(id) {
			continue
		}
		if t.onFlank && ctx.Resolver.IsCreature(id) && !onFlank(ctx, id) {
			continue
		}
		if t.notOnFlank && onFlank(ctx, id) {
			continue
		}
		if t.neighboring && !isNeighbor(ctx, ctx.Source, id) {
			continue
		}
		if t.other && id == ctx.Source {
			continue
		}
		if t.named != "" && ctx.Resolver.Name(id) != t.named {
			continue
		}
		out = append(out, id)
	}
	return out
}

// onFlank reports whether a creature is on a flank of its battleline (its
// leftmost or rightmost creature).
func onFlank(ctx *EffectContext, id LocalID) bool {
	bl := battlelineContaining(ctx, id)
	return len(bl) > 0 &&
		(bl[0] == id || bl[len(bl)-1] == id)
}

// isNeighbor reports whether id is one of src's battleline neighbors.
func isNeighbor(ctx *EffectContext, src, id LocalID) bool {
	for _, n := range neighbors(ctx, src) {
		if n == id {
			return true
		}
	}
	return false
}

// neighbors returns the creatures immediately adjacent to id in its controller's
// battleline — its left and right neighbors, when present. A card that is not in
// a battleline has no neighbors.
func neighbors(ctx *EffectContext, id LocalID) []LocalID {
	bl := battlelineContaining(ctx, id)
	i := -1
	for j, x := range bl {
		if x == id {
			i = j
			break
		}
	}
	if i < 0 {
		return nil
	}
	out := make([]LocalID, 0, 2)
	if i > 0 {
		out = append(out, bl[i-1])
	}
	if i < len(bl)-1 {
		out = append(out, bl[i+1])
	}
	return out
}

// creaturesExcept returns every creature in play except one, walking both
// players' battlelines in order. It backs effects that target "a different
// creature" or "another creature" than one already chosen.
func creaturesExcept(ctx *EffectContext, exclude LocalID) []LocalID {
	var out []LocalID
	for p := 0; p < 2; p++ {
		for _, id := range ctx.Resolver.Battleline(p) {
			if id != exclude {
				out = append(out, id)
			}
		}
	}
	return out
}

func battlelineContaining(ctx *EffectContext, id LocalID) []LocalID {
	for p := 0; p < 2; p++ {
		bl := ctx.Resolver.Battleline(p)
		for _, x := range bl {
			if x == id {
				return bl
			}
		}
	}
	return nil
}

// selectBase resolves the unfiltered base set chosen by Kind. Chosen kinds return
// the pool of candidates; Select applies filters and prompts for the choice.
func (t Target) selectBase(ctx *EffectContext) []LocalID {
	switch t.Kind {
	case TargetThisCreature:
		return []LocalID{ctx.Source}
	case TargetTriggeringCreature, TargetTheOtherCreature, TargetCreatureFought:
		if ctx.HasIt {
			return []LocalID{ctx.It}
		}
		return nil
	case TargetEachArtifact, TargetChosenArtifact:
		return append(
			ctx.Resolver.Artifacts(ctx.Controller),
			ctx.Resolver.Artifacts(ctx.Opponent())...)
	case TargetChosenEnemyArtifact:
		return ctx.Resolver.Artifacts(ctx.Opponent())
	case TargetEachEnemyArtifact:
		return ctx.Resolver.Artifacts(ctx.Opponent())
	case TargetEachFriendlyArtifact:
		return ctx.Resolver.Artifacts(ctx.Controller)
	case TargetEachCardInPlay, TargetChosenCreatureOrArtifact:
		ids := ctx.Resolver.Battleline(ctx.Controller)
		ids = append(ids, ctx.Resolver.Battleline(ctx.Opponent())...)
		ids = append(ids, ctx.Resolver.Artifacts(ctx.Controller)...)
		ids = append(ids, ctx.Resolver.Artifacts(ctx.Opponent())...)
		return ids
	case TargetEachFriendlyCardInPlay, TargetChosenFriendlyCreatureOrArtifact:
		return append(
			ctx.Resolver.Battleline(ctx.Controller),
			ctx.Resolver.Artifacts(ctx.Controller)...)
	case TargetEachCreature, TargetChosenCreature:
		return append(
			ctx.Resolver.Battleline(ctx.Controller),
			ctx.Resolver.Battleline(ctx.Opponent())...)
	case TargetEachFriendlyCreature, TargetChosenFriendlyCreature:
		return ctx.Resolver.Battleline(ctx.Controller)
	case TargetEachEnemyCreature, TargetChosenEnemyCreature:
		return ctx.Resolver.Battleline(ctx.Opponent())
	case TargetEachOtherFriendlyCreature, TargetChosenOtherFriendlyCreature:
		out := make([]LocalID, 0)
		for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
			if id != ctx.Source {
				out = append(out, id)
			}
		}
		return out
	case TargetChosenOtherCreature:
		if !ctx.HasIt {
			return append(
				ctx.Resolver.Battleline(ctx.Controller),
				ctx.Resolver.Battleline(ctx.Opponent())...)
		}
		return creaturesExcept(ctx, ctx.It)
	default:
		return nil
	}
}

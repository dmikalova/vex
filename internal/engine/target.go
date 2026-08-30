package engine

import (
	"fmt"
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
	// TargetEachCreature selects every creature in play.
	TargetEachCreature
	// TargetEachFriendlyCreature selects every creature the controller controls.
	TargetEachFriendlyCreature
	// TargetEachEnemyCreature selects every creature the opponent controls.
	TargetEachEnemyCreature
	// TargetEachArtifact selects every artifact in play, both players'.
	TargetEachArtifact
	// TargetEachInPlay selects every card in play — every creature and artifact,
	// both players' — including the source card itself.
	TargetEachInPlay
	// TargetEachFriendlyInPlay selects the controller's cards in play — their
	// creatures and artifacts.
	TargetEachFriendlyInPlay
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
	// TargetChosenArtifact selects a single artifact the controller chooses from
	// all artifacts in play (either player's).
	TargetChosenArtifact
	// TargetTheOtherCreature selects the creature in context (ctx.It) — the one a
	// preceding effect chose as "another" creature — and renders it as "the other
	// creature". Transposition Sandals swaps with another creature, then uses the
	// other creature.
	TargetTheOtherCreature
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
	maxPower    int
	hasMaxPower bool
	damaged     bool
	stunned     bool
	onFlank     bool
	notOnFlank  bool
	neighboring bool
	// other excludes the source card from the selected set ("other" cards).
	other bool
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

// PowerAtMost narrows the target to creatures whose power is max or lower, e.g.
// Target{Kind: TargetEachCreature}.PowerAtMost(3).
func (t Target) PowerAtMost(max int) Target {
	t.maxPower = max
	t.hasMaxPower = true
	return t
}

// Damaged narrows the target to creatures that currently have damage on them.
func (t Target) Damaged() Target {
	t.damaged = true
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
// leftmost or rightmost creature).
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
		return "it"
	case TargetTheOtherCreature:
		return "the other creature"
	}
	noun := "creature"
	if t.Kind == TargetEachArtifact || t.Kind == TargetChosenArtifact {
		noun = "artifact"
	}
	if t.Kind == TargetEachFriendlyInPlay {
		noun = "card"
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
		noun = "flank " + noun
	}
	if t.neighboring {
		noun = "neighboring " + noun
	}
	if t.damaged {
		noun = "damaged " + noun
	}
	if t.stunned {
		noun = "stunned " + noun
	}
	var phrase string
	switch t.Kind {
	case TargetEachInPlay:
		phrase = "each card in play"
	case TargetEachCreature, TargetEachArtifact:
		phrase = "each " + noun
	case TargetEachFriendlyCreature:
		phrase = "each friendly " + noun
	case TargetEachFriendlyInPlay:
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
	case TargetChosenFriendlyCreature:
		phrase = "a friendly " + noun
	case TargetChosenOtherFriendlyCreature:
		phrase = "another friendly " + noun
	case TargetChosenArtifact:
		phrase = "an " + noun
	default:
		phrase = "a " + noun
	}
	if t.hasMaxPower {
		phrase += fmt.Sprintf(" with power %d or lower", t.maxPower)
	}
	if t.notOnFlank {
		phrase += " that is not on a flank"
	}
	if t.chosenHouse {
		phrase += " of the chosen house"
	}
	if t.selector != nil {
		phrase = t.selector.clause(phrase)
	}
	return phrase
}

// Select resolves the target into concrete card ids, applying its filters. For a
// chosen kind it asks the controller to pick one of the filtered candidates
// (returning nil when there are none or the choice is declined).
func (t Target) Select(ctx *EffectContext) []LocalID {
	ids := t.filter(ctx, t.selectBase(ctx))
	if t.selector != nil {
		ids = t.selector.refine(ctx, ids)
	}
	if !t.isChosen() {
		return ids
	}
	if len(ids) == 0 {
		return nil
	}
	id, ok := ctx.ChooseCreature("Choose "+t.Text(), ids)
	if !ok {
		return nil
	}
	return []LocalID{id}
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
	max := ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		if p := ctx.Resolver.Power(id); p > max {
			max = p
		}
	}
	mostPowerful := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == max {
			mostPowerful = append(mostPowerful, id)
		}
	}
	spared := mostPowerful[0]
	if len(mostPowerful) > 1 {
		if chosen, ok := ctx.ChooseCreature("Choose the most powerful creature to keep", mostPowerful); ok {
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

// isChosen reports whether the Kind resolves to a single player-chosen creature.
func (t Target) isChosen() bool {
	return t.Kind == TargetChosenCreature || t.Kind == TargetChosenEnemyCreature ||
		t.Kind == TargetChosenFriendlyCreature || t.Kind == TargetChosenOtherFriendlyCreature ||
		t.Kind == TargetChosenArtifact
}

// filter narrows ids to those matching the target's trait, power, damaged, and
// flank filters.
func (t Target) filter(ctx *EffectContext, ids []LocalID) []LocalID {
	if t.trait == "" &&
		t.exceptTrait == "" &&
		t.house == HouseNone &&
		t.exceptHouse == HouseNone &&
		!t.chosenHouse &&
		!t.hasMaxPower &&
		!t.damaged &&
		!t.stunned &&
		!t.onFlank &&
		!t.notOnFlank &&
		!t.neighboring &&
		!t.other {
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
		if t.hasMaxPower && ctx.Resolver.Power(id) > t.maxPower {
			continue
		}
		if t.damaged && ctx.Resolver.Damage(id) == 0 {
			continue
		}
		if t.stunned && !ctx.Resolver.Stunned(id) {
			continue
		}
		if t.onFlank && !onFlank(ctx, id) {
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
	case TargetTriggeringCreature, TargetTheOtherCreature:
		if ctx.HasIt {
			return []LocalID{ctx.It}
		}
		return nil
	case TargetEachArtifact, TargetChosenArtifact:
		return append(ctx.Resolver.Artifacts(ctx.Controller), ctx.Resolver.Artifacts(ctx.Opponent())...)
	case TargetEachInPlay:
		ids := ctx.Resolver.Battleline(ctx.Controller)
		ids = append(ids, ctx.Resolver.Battleline(ctx.Opponent())...)
		ids = append(ids, ctx.Resolver.Artifacts(ctx.Controller)...)
		ids = append(ids, ctx.Resolver.Artifacts(ctx.Opponent())...)
		return ids
	case TargetEachFriendlyInPlay:
		return append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Artifacts(ctx.Controller)...)
	case TargetEachCreature, TargetChosenCreature:
		return append(ctx.Resolver.Battleline(ctx.Controller), ctx.Resolver.Battleline(ctx.Opponent())...)
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
	default:
		return nil
	}
}

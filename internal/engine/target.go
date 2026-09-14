package engine

import (
	"fmt"
	"strings"
)

// A Target names the cards an ability acts on. KeyForge abilities are written in
// terms of noun phrases — "this creature", "each enemy creature", "a friendly
// creature", "each Scientist creature", "each creature with power 3 or
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
	// TargetCreatureFought selects the creature the source fought, named in full.
	// A Fight: ability resolves after the fight, so the bare form reads in the past
	// ("the creature <self> fought"). A Before Fight ability that reaches past it —
	// Lord Golgotha damaging its neighbors — resolves before the fight, so the
	// neighbor-decorated form reads in the present ("the creature <self> fights").
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
	// TargetChosenFriendlyArtifact selects a single friendly artifact the
	// controller chooses (Anahita the Trader gives one away).
	TargetChosenFriendlyArtifact
	// TargetChosenUpgrade selects a single upgrade the controller chooses from all
	// upgrades attached to any creature in play (Destroy Them All! destroys one).
	TargetChosenUpgrade
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
	// TargetChosenEnemyCreatureOrArtifact selects a single enemy creature or artifact
	// the controller chooses, rendered "an enemy creature or artifact".
	TargetChosenEnemyCreatureOrArtifact
	// TargetTheChosenCreature selects the creature a preceding ChooseCreatureThen
	// put in context (ctx.It) and renders it as "the chosen creature", the standard
	// referent for a just-chosen creature (Sack of Coins deals its per-Æmber damage
	// to the chosen creature).
	TargetTheChosenCreature
	// TargetFormerNeighbors selects the battleline neighbors a preceding effect
	// snapshotted before it removed a creature (ctx.Produced.Neighbors) and renders
	// them as "each of that creature's neighbors" — Pain Reaction hits the neighbors
	// of the creature its damage just destroyed.
	TargetFormerNeighbors
	// TargetTheFoughtCreature selects the creature a preceding effect had a chosen
	// creature fight (ctx.It) and renders it as "the fought creature", naming no
	// fighter — Smite makes a friendly creature fight, then damages the fought
	// creature's neighbors, so the fight is not the source's own.
	TargetTheFoughtCreature
	// TargetAttachedHost selects the creature the resolving upgrade (ctx.Upgrade)
	// is attached to — the exact instance a blaster bound to when AttachSelfTo
	// homed it, identified by LocalID rather than by name. Its payoff acts on that
	// one creature, never on a second same-named copy, and selects nothing once the
	// bound instance has left play (its upgrade is no longer attached). Named()
	// supplies the printed name so the text still reads as the signature creature.
	TargetAttachedHost
	// TargetGrantingCard selects the in-play card that granted the resolving ability
	// (ctx.Grantor) — the exact artifact or upgrade whose constant ability or static
	// modifier handed this creature its ability, identified by LocalID rather than by
	// name. A creature reaping through Uncharted Lands' grant moves Æmber off that one
	// artifact, never off a second same-named copy in play. Its text renders the
	// {card} placeholder, which the granted-text renderer resolves to the granting
	// card's own name ("from Uncharted Lands").
	TargetGrantingCard
	// TargetEachNeighbor selects the source card's live battleline neighbors and
	// renders them as "each of <self>'s neighbors" — Ghosthawk reaps with each of
	// its neighbors, one at a time.
	TargetEachNeighbor
)

// Target describes which cards an effect applies to. Kind picks the base set;
// the optional filters added by WithTrait and PowerAtMost narrow that set and
// extend the rendered text.
type Target struct {
	Kind        TargetKind
	trait       Trait
	exceptTrait Trait
	// house narrows the target to the cards the matcher admits — a named house, every
	// house but one, the chosen house, the active house, or the house of the card in
	// context (ctx.It). The zero value (any house) narrows nothing. Because the field
	// is unexported, a SelfHouse sentinel in it resolves through selfHouseResolved
	// rather than by reflection.
	house HouseMatcher
	// houseWithMostCreatures narrows the target to creatures of the house with the
	// most creatures in play, counting both players' battlelines; on a tie every
	// tied house's creatures are eligible so the chooser picks among them
	// (Etaromme). It renders "of the house with the most creatures in play".
	houseWithMostCreatures bool
	// sharesTrait narrows the target to cards sharing at least one trait with the
	// card in context (ctx.It), rendering "that shares a trait with it".
	sharesTrait   bool
	maxPower      int
	hasMaxPower   bool
	minPower      int
	hasMinPower   bool
	exactPower    int
	hasExactPower bool
	// oddPower narrows the target to creatures whose power is odd, and evenPower to
	// those whose power is even (Onyx Knight, Opal Knight).
	oddPower  bool
	evenPower bool
	damaged   bool
	undamaged bool
	stunned   bool
	// ready narrows the target to creatures that are not exhausted (Swap Widget's
	// "a ready friendly Mars creature").
	ready      bool
	withAember bool
	// withoutAember narrows the target to creatures that have no Æmber on them,
	// rendering " with no Æmber on it" (Draining Touch destroys a creature with no
	// Æmber on it).
	withoutAember bool
	// withCounter narrows the target to cards carrying a generic counter of this
	// kind, rendering " with a doom counter" and the like (Wretched Doll destroys
	// every creature with a doom counter). CounterNone leaves the filter off.
	withCounter CounterKind
	// withArmor narrows the target to creatures that have armor at all, rendering
	// " with armor". It reads the creature's armor value, not what is left of it, so
	// a creature that has already spent its armor absorbing damage still has armor.
	withArmor bool
	// withUpgrade narrows the target to creatures that have at least one upgrade
	// attached, rendering " with an upgrade" (Tachyon Pulse exhausts each creature
	// with an upgrade).
	withUpgrade bool
	// sharesHouseNeighbors narrows the target to creatures sharing a house with at
	// least this many of their battleline neighbors, rendering "that shares a house
	// with N of its neighbors" (Groupthink Tank, Mini Groupthink Tank). Zero leaves
	// the filter off.
	sharesHouseNeighbors int
	keyword              Keyword
	onFlank              bool
	notOnFlank           bool
	neighboring          bool
	// toRightOfSource narrows the target to the creatures positioned to the right of
	// the source card in its battleline, and toLeftOfSource to those on its left —
	// the Panpacas, which buff one direction of their line.
	toRightOfSource bool
	toLeftOfSource  bool
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
	// refinement is a set-relative refinement applied after the per-card filters. It
	// can compare the candidates to each other (e.g. "except the most powerful")
	// and contributes a clause to the printed phrase. nil for targets that select
	// their whole filtered set.
	refinement Refinement
}

// WithTrait narrows the target to cards that have the given trait, e.g.
// Target{Kind: TargetEachCreature}.WithTrait(Scientist).
func (t Target) WithTrait(trait Trait) Target {
	t.trait = trait
	return t
}

// ExceptTrait narrows the target to cards that do NOT have the given trait,
// rendering the "non-<trait>" qualifier, e.g. a friendly Mars creature
// ExceptTrait(Agent) reads "a friendly non-Agent Mars creature".
func (t Target) ExceptTrait(trait Trait) Target {
	t.exceptTrait = trait
	return t
}

// House narrows the target to the cards a HouseMatcher admits — a named house,
// every house but one, the chosen house, the active house, or the house of the
// card in context. Target{Kind: TargetEachCreature}.House(namedHouse(Mars)) reads
// "each Mars creature"; .House(exceptHouse(Sanctum)) reads "each non-Sanctum
// creature"; .House(HouseMatcher{Kind: MatchChosenHouse}) reads "each creature of
// the chosen house".
func (t Target) House(m HouseMatcher) Target {
	t.house = m
	return t
}

// OfHouseWithMostCreatures narrows the target to creatures of the house with the
// most creatures in play across both battlelines, ties keeping every tied house
// eligible (Etaromme).
func (t Target) OfHouseWithMostCreatures() Target {
	t.houseWithMostCreatures = true
	return t
}

// selfHouseResolved fills the card's own house in for a SelfHouse sentinel the
// target narrows on. A Target keeps its house matcher and refinement unexported,
// so it resolves itself rather than being rewritten by reflection (see
// self_house.go).
func (t Target) selfHouseResolved(house House) any {
	if t.house.House == SelfHouse {
		t.house.House = house
	}
	if t.refinement != nil {
		t.refinement = resolvedIn(t.refinement, house)
	}
	return t
}

// SharingTrait narrows the target to cards that share at least one trait with the
// card in context (ctx.It), rendering "that shares a trait with it" — the purged
// creature after a PurgeFromHand that moves a single card (Custom Virus), or
// whatever an earlier effect put in context.
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

// OddPower narrows the target to creatures whose power is odd (Onyx Knight).
func (t Target) OddPower() Target {
	t.oddPower = true
	return t
}

// EvenPower narrows the target to creatures whose power is even (Opal Knight).
func (t Target) EvenPower() Target {
	t.evenPower = true
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

// WithoutAember narrows the target to creatures that have no Æmber on them,
// rendering " with no Æmber on it", e.g. "a creature with no Æmber on it".
func (t Target) WithoutAember() Target {
	t.withoutAember = true
	return t
}

// WithCounter narrows the target to cards carrying a generic counter of the
// given kind, rendering " with a <kind>", e.g. "each creature with a doom
// counter".
func (t Target) WithCounter(kind CounterKind) Target {
	t.withCounter = kind
	return t
}

// WithArmor narrows the target to creatures that have armor, rendering " with
// armor", e.g. "each enemy creature with armor".
func (t Target) WithArmor() Target {
	t.withArmor = true
	return t
}

// WithUpgrade narrows the target to creatures that have at least one upgrade
// attached, rendering " with an upgrade", e.g. "each creature with an upgrade"
// (Tachyon Pulse).
func (t Target) WithUpgrade() Target {
	t.withUpgrade = true
	return t
}

// SharesHouseWithNeighbors narrows the target to creatures sharing a house with
// at least the given number of their battleline neighbors, rendering "that shares
// a house with at least 1 of its neighbors" for 1 (Groupthink Tank) and "that
// shares a house with N of its neighbors" otherwise (Mini Groupthink Tank).
func (t Target) SharesHouseWithNeighbors(atLeast int) Target {
	t.sharesHouseNeighbors = atLeast
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

// Ready narrows the target to creatures that are not exhausted.
func (t Target) Ready() Target {
	t.ready = true
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

// ToRightOfSource narrows the target to the creatures positioned to the right of
// the source card in its battleline (Panpaca, Anga).
func (t Target) ToRightOfSource() Target {
	t.toRightOfSource = true
	return t
}

// ToLeftOfSource narrows the target to the creatures positioned to the left of
// the source card in its battleline (Panpaca, Jaga).
func (t Target) ToLeftOfSource() Target {
	t.toLeftOfSource = true
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

// Refine refines the target with a set-relative rule applied after the per-card
// filters, e.g. Target{...}.Refine(ExceptMostPowerful). The Refinement both picks
// the final subset and describes itself for the printed phrase, so a niche
// "relative to the rest of the set" rule composes onto any Target without adding
// a dedicated field (and future rules — least powerful, and so on — are just more
// Refinement values).
func (t Target) Refine(s Refinement) Target {
	t.refinement = s
	return t
}

// valid reports whether the target's base set was chosen (its Kind is not the
// unset zero value). Effects that require a target check this in validation.
func (t Target) valid() bool {
	return t.Kind != targetUnset
}

// plural reports whether the target names more than one card — exactly the
// "each ..." phrasing, read off Text rather than a parallel switch.
func (t Target) plural() bool {
	return strings.HasPrefix(t.Text(), "each")
}

// pronoun renders the target as a back-reference for a sentence whose antecedent
// already named these creatures — "those creatures" for a plural (each) set, "that
// creature" for a single one.
func (t Target) pronoun() string {
	if t.plural() {
		return "those creatures"
	}
	return "that creature"
}

// Text renders the target as an English noun phrase, e.g. "each enemy creature",
// "each Scientist creature", or "each creature with power 3 or lower".
func (t Target) Text() string {
	switch t.Kind {
	case TargetThisCreature:
		return SelfName
	case TargetTriggeringCreature:
		return t.decorateNeighbors("it")
	case TargetCreatureFought:
		if t.withNeighbors || t.neighborsOf {
			return t.decorateNeighbors("the creature " + SelfName + " fights")
		}
		return "the creature " + SelfName + " fought"
	case TargetTheOtherCreature:
		return "the other creature"
	case TargetTheChosenCreature:
		return "the chosen creature"
	case TargetAttachedHost:
		if t.named != "" {
			return t.named
		}
		return "the attached creature"
	case TargetGrantingCard:
		return CardName
	case TargetFormerNeighbors:
		return "each of that creature's neighbors"
	case TargetEachNeighbor:
		return "each of " + SelfName + "'s neighbors"
	case TargetTheFoughtCreature:
		return t.decorateNeighbors("the fought creature")
	}
	noun := "creature"
	if t.Kind == TargetEachArtifact || t.Kind == TargetChosenArtifact ||
		t.Kind == TargetChosenEnemyArtifact || t.Kind == TargetEachFriendlyArtifact ||
		t.Kind == TargetChosenFriendlyArtifact ||
		t.Kind == TargetEachEnemyArtifact {
		noun = "artifact"
	}
	if t.Kind == TargetChosenUpgrade {
		noun = "upgrade"
	}
	if t.Kind == TargetEachFriendlyCardInPlay {
		noun = "card"
	}
	orArtifact := t.Kind == TargetChosenCreatureOrArtifact ||
		t.Kind == TargetChosenFriendlyCreatureOrArtifact ||
		t.Kind == TargetChosenEnemyCreatureOrArtifact
	if orArtifact {
		noun = "creature or artifact"
	}
	if t.named != "" {
		noun = t.named
	}
	if t.trait != traitUnset {
		noun = t.trait.String() + " " + noun
	}
	noun = t.house.qualifyNoun(noun)
	if t.exceptTrait != traitUnset {
		noun = "non-" + t.exceptTrait.String() + " " + noun
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
	if t.ready {
		noun = "ready " + noun
	}
	if t.keyword.valid() {
		noun = strings.ToLower(t.keyword.String()) + " " + noun
	}
	if t.named != "" && t.isChosen() {
		// A proper name identifies one specific card, which takes no article: a
		// single-target phrase names it outright ("ward Lieutenant Khrkhar", not
		// "ward a friendly Lieutenant Khrkhar").
		return t.decorateNeighbors(noun)
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
	case TargetChosenUpgrade:
		phrase = "an " + noun
	case TargetChosenEnemyArtifact, TargetChosenEnemyCreatureOrArtifact:
		phrase = "an enemy " + noun
	case TargetChosenFriendlyArtifact:
		phrase = "a friendly " + noun
	case TargetChosenCreature:
		// Other() drops the source card, which reads as "another creature" — the
		// wording Replicator prints for a creature in play that is not itself.
		if t.other {
			phrase = "another " + noun
		} else {
			phrase = indefinite(noun)
		}
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
		phrase += fmt.Sprintf(" with power %d", t.exactPower)
	}
	if t.oddPower {
		phrase += " with odd power"
	}
	if t.evenPower {
		phrase += " with even power"
	}
	if t.withAember {
		phrase += " with \u00c6mber on it"
	}
	if t.withoutAember {
		phrase += " with no \u00c6mber on it"
	}
	if t.withCounter.valid() {
		phrase += " with a " + t.withCounter.noun()
	}
	if t.withArmor {
		phrase += " with armor"
	}
	if t.withUpgrade {
		phrase += " with an upgrade"
	}
	if t.sharesHouseNeighbors == 1 {
		phrase += " that shares a house with at least 1 of its neighbors"
	} else if t.sharesHouseNeighbors > 1 {
		phrase += fmt.Sprintf(
			" that shares a house with %d of its neighbors",
			t.sharesHouseNeighbors,
		)
	}
	if t.notOnFlank {
		phrase += " that is not on a flank"
	}
	if t.toRightOfSource {
		phrase += " to the right of " + SelfName
	}
	if t.toLeftOfSource {
		phrase += " to the left of " + SelfName
	}
	phrase = t.house.qualifyPhrase(phrase)
	if t.houseWithMostCreatures {
		phrase += " of the house with the most creatures in play"
	}
	if t.sharesTrait {
		phrase += " that shares a trait with it"
	}
	if t.refinement != nil {
		phrase = t.refinement.clause(phrase)
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

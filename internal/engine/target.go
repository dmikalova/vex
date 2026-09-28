package engine

import "strings"

// A Target names the cards an ability acts on. KeyForge abilities are written in
// terms of noun phrases — "this creature", "each enemy creature", "a friendly
// creature", "each Scientist creature", "each creature with power 3 or
// lower" — and Target captures exactly that: a base set chosen by Kind, narrowed
// per card by one Filter and relative to the whole set by one Refinement. Because
// the same value both renders the phrase (Text) and
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
	// It renders in the past ("the creature <self> fought") because a Fight: ability
	// resolves after the fight; the renderer puts it in the present for a Before
	// Fight: ability, which resolves before the fight (see fightTense).
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
	// TargetChosenCreature selects a single creature the controller chooses from
	// all creatures in play (either player's).
	TargetChosenCreature
	// TargetChosenEnemyCreature selects a single enemy creature the controller
	// chooses.
	TargetChosenEnemyCreature
	// TargetChosenFriendlyCreature selects a single friendly creature the
	// controller chooses.
	TargetChosenFriendlyCreature
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
	// TargetTheSameCreature selects the triggering creature (ctx.It) like
	// TargetTriggeringCreature but renders it as "the same creature" — used when the
	// trigger clause names no creature to be an antecedent for "it", as with the
	// bonus-icon triggers where Maleficorn deals damage to the creature that just
	// resolved a Damage bonus icon.
	TargetTheSameCreature
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
	// TargetEachUpgradeOnThis selects the upgrades attached to the source card and
	// renders them as "each upgrade on <self>" — Away Team archives its own upgrades
	// when it is destroyed.
	TargetEachUpgradeOnThis
	// targetKindCount bounds the enumeration. It is not a real kind.
	targetKindCount
)

// TargetKinds returns every real target kind, so a caller can enumerate them
// rather than maintain its own list that a newly added kind would fall out of.
func TargetKinds() []TargetKind {
	all := make([]TargetKind, 0, int(targetKindCount)-1)
	for k := targetUnset + 1; k < targetKindCount; k++ {
		all = append(all, k)
	}
	return all
}

// isContextReference reports whether a kind names a card the resolution context
// already holds — ctx.It, ctx.Upgrade, ctx.Grantor, the source card, a snapshot a
// preceding effect left behind — rather than a set to be found by scanning the
// board. It is the line three switches in this file draw independently: a context
// reference resolves without a search, renders through specialKindText as a
// definite phrase ("the chosen creature") instead of a quantified one ("each enemy
// creature"), and is never something the player picks. The partition is pinned by
// TestContextReferencesAreExactlyTheSpecialTextKinds.
func isContextReference(k TargetKind) bool {
	switch k {
	case TargetThisCreature, TargetTriggeringCreature, TargetCreatureFought,
		TargetTheOtherCreature, TargetTheSameCreature, TargetTheChosenCreature,
		TargetAttachedHost, TargetGrantingCard, TargetFormerNeighbors,
		TargetEachNeighbor, TargetEachUpgradeOnThis, TargetTheFoughtCreature:
		return true
	}
	return false
}

// creatureOrArtifact is the noun the kinds that name both types render. It is
// named because the flank adjective has to find it to bind to the creature half
// alone (Snudge's "an artifact or flank creature").
const creatureOrArtifact = "creature or artifact"

// Target describes which cards an effect applies to. Kind picks the base set;
// the filter narrows that set one card at a time and extends the rendered text;
// the refinement narrows it relative to the whole set.
type Target struct {
	Kind TargetKind
	// Filter is the per-card narrowing — trait, house, power, state, position,
	// identity — held as one comparable value that both selects and renders (see
	// filter.go). Write it with With; the zero value narrows nothing.
	Filter Filter
	// Refinement is a set-relative rule applied after the Filter. It can compare the
	// candidates to each other (e.g. "except the most powerful") or prompt across
	// them, and contributes a clause to the printed phrase. nil for targets that
	// select their whole filtered set.
	//
	// It is the one field that keeps a Target from being unconditionally comparable:
	// AnyOf holds a slice, so == on a Target carrying one panics. Isolating that
	// hazard to a single field is the point — every other field is a value.
	Refinement Refinement
	// Neighbors is how the target reaches the battleline neighbors of what it
	// selects. It expands the selected set rather than testing each candidate, which
	// is why it sits here rather than in the Filter.
	Neighbors NeighborMode
}

// With narrows the target by a Filter — the one method a filtered card reaches
// for, so a narrowing is written as one struct literal rather than a chain:
// card.Target.EachCreature.With(card.Filter{Power: card.Power.AtMost(3)}). It
// replaces the target's filter wholesale, since a Filter already says every axis.
func (t Target) With(f Filter) Target {
	t.Filter = f
	return t
}

// AndNeighbors expands a single chosen creature to also include its battleline
// neighbors, so an effect applies to the chosen creature and each of its
// neighbors (Tremor). It is meaningful only on a chosen-creature target.
func (t Target) AndNeighbors() Target {
	t.Neighbors = NeighborsIncluded
	return t
}

// NeighborsOf narrows a target to the battleline neighbors of what it selects,
// dropping the selected creature itself — Lord Golgotha damages each neighbor of
// the creature it fights, but not that creature.
func (t Target) NeighborsOf() Target {
	t.Neighbors = NeighborsOnly
	return t
}

// Refine refines the target with a set-relative rule applied after the per-card
// filters, e.g. Target{...}.Refine(ExceptMostPowerful). The Refinement both picks
// the final subset and describes itself for the printed phrase, so a niche
// "relative to the rest of the set" rule composes onto any Target without adding
// a dedicated field (and future rules — least powerful, and so on — are just more
// Refinement values).
//
// It panics on a bare Filter. A Filter is a Refinement so a union combinator can
// take one as a member (Regrettable Meteor unions a trait against a power floor),
// but alone it says exactly what the Filter field says, and one narrowing must
// have one spelling — write it as With. Pinned by
// TestFilterAsRefinementInAUnion.
func (t Target) Refine(s Refinement) Target {
	if _, bare := s.(Filter); bare {
		panic("Target.Refine: a Filter alone is not a refinement; pass it to With")
	}
	t.Refinement = s
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
	if s, ok := t.specialKindText(); ok {
		return s
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
	if t.Kind == TargetChosenCreatureOrArtifact ||
		t.Kind == TargetChosenFriendlyCreatureOrArtifact ||
		t.Kind == TargetChosenEnemyCreatureOrArtifact {
		noun = creatureOrArtifact
	}
	noun = t.Filter.qualifyNoun(noun)
	if t.Filter.Name != "" && t.isChosen() {
		// A proper name identifies one specific card, which takes no article: a
		// single-target phrase names it outright ("ward Lieutenant Khrkhar", not
		// "ward a friendly Lieutenant Khrkhar").
		return t.decorateNeighbors(noun)
	}
	return t.decorateNeighbors(t.quantifiedPhrase(noun))
}

// specialKindText renders the target Kinds whose phrasing is a fixed noun phrase
// rather than a noun built up from adjectives and a quantifier. The bool is false
// for a Kind that falls through to the general noun/phrase construction.
func (t Target) specialKindText() (string, bool) {
	switch t.Kind {
	case TargetThisCreature:
		return SelfName, true
	case TargetTriggeringCreature:
		return t.decorateNeighbors("it"), true
	case TargetCreatureFought:
		return t.decorateNeighbors("the creature " + SelfName + " fought"), true
	case TargetTheOtherCreature:
		return "the other creature", true
	case TargetTheSameCreature:
		return "the same creature", true
	case TargetTheChosenCreature:
		return t.decorateNeighbors("the chosen creature"), true
	case TargetAttachedHost:
		if t.Filter.Name != "" {
			return t.Filter.Name, true
		}
		return "the attached creature", true
	case TargetGrantingCard:
		return CardName, true
	case TargetFormerNeighbors:
		return "each of that creature's neighbors", true
	case TargetEachNeighbor:
		return "each of " + SelfName + "'s neighbors", true
	case TargetEachUpgradeOnThis:
		return "each upgrade on " + SelfName, true
	case TargetTheFoughtCreature:
		return t.decorateNeighbors("the fought creature"), true
	}
	return "", false
}

// quantifiedPhrase turns a constructed noun into the full phrase: it prefixes the
// article or quantifier the target's Kind calls for ("each", "a friendly",
// "another", …), then appends the power, resource, position, and refinement clauses
// that narrow which cards match. The clauses accumulate in printed order.
func (t Target) quantifiedPhrase(noun string) string {
	var phrase string
	switch t.Kind {
	case TargetEachCardInPlay:
		phrase = "each card in play"
	case TargetEachCreature, TargetEachArtifact:
		phrase = "each " + t.Filter.Except.qualifyNoun(noun)
	case TargetEachFriendlyCreature, TargetEachFriendlyArtifact,
		TargetEachFriendlyCardInPlay:
		phrase = "each " + t.Filter.Except.qualifyNoun("friendly "+noun)
	case TargetEachEnemyArtifact:
		phrase = "each enemy " + noun
	case TargetEachEnemyCreature:
		phrase = "each enemy " + noun
	case TargetChosenEnemyCreature:
		phrase = "an enemy " + noun
	case TargetChosenFriendlyCreature, TargetChosenFriendlyCreatureOrArtifact:
		phrase = t.Filter.Except.quantifyOne("friendly " + noun)
	case TargetChosenArtifact:
		phrase = "an " + noun
	case TargetChosenUpgrade:
		phrase = "an " + noun
	case TargetChosenCreatureOrArtifact:
		phrase = indefinite(noun)
	case TargetChosenEnemyArtifact, TargetChosenEnemyCreatureOrArtifact:
		phrase = "an enemy " + noun
	case TargetChosenFriendlyArtifact:
		phrase = "a friendly " + noun
	case TargetChosenCreature:
		phrase = t.Filter.Except.quantifyOne(noun)
	default:
		phrase = indefinite(noun)
	}
	return t.narrowingClauses(phrase)
}

// narrowingClauses appends the trailing clauses that narrow a quantified phrase
// to the cards that match, in printed order: the filter's own clauses first, then
// the set-relative refinement's, which frames the whole filtered phrase.
func (t Target) narrowingClauses(phrase string) string {
	phrase = t.Filter.clauses(phrase)
	if t.Refinement != nil {
		phrase = t.Refinement.clause(phrase)
	}
	return phrase
}

// decorateNeighbors wraps a rendered noun phrase in the wording its NeighborMode
// calls for.
func (t Target) decorateNeighbors(phrase string) string {
	return t.Neighbors.decorate(phrase)
}

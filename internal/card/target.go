package card

import "github.com/dmikalova/vactrol/internal/engine"

// Target groups ready-made targets, e.g. card.Target.EachEnemyCreature. Each is
// an engine.Target value, so the filter methods (WithTrait, PowerAtMost, OnFlank,
// Refine, ...) chain off them:
// card.Target.EachEnemyCreature.Refine(card.Except(card.MostPowerful)).
var Target = targets{
	This:                 engine.Target{Kind: engine.TargetThisCreature},
	Triggering:           engine.Target{Kind: engine.TargetTriggeringCreature},
	Creature:             engine.Target{Kind: engine.TargetChosenCreature},
	FriendlyCreature:     engine.Target{Kind: engine.TargetChosenFriendlyCreature},
	EnemyCreature:        engine.Target{Kind: engine.TargetChosenEnemyCreature},
	EachCreature:         engine.Target{Kind: engine.TargetEachCreature},
	EachFriendlyCreature: engine.Target{Kind: engine.TargetEachFriendlyCreature},
	EachEnemyCreature:    engine.Target{Kind: engine.TargetEachEnemyCreature},
	EachArtifact:         engine.Target{Kind: engine.TargetEachArtifact},
	EachFriendlyArtifact: engine.Target{
		Kind: engine.TargetEachFriendlyArtifact,
	},
	EachEnemyArtifact:          engine.Target{Kind: engine.TargetEachEnemyArtifact},
	EachFriendlyCardInPlay:     engine.Target{Kind: engine.TargetEachFriendlyCardInPlay},
	EachOtherFriendlyCreature:  engine.Target{Kind: engine.TargetEachOtherFriendlyCreature},
	OtherFriendlyCreature:      engine.Target{Kind: engine.TargetChosenOtherFriendlyCreature},
	OtherCreature:              engine.Target{Kind: engine.TargetChosenOtherCreature},
	TheOtherCreature:           engine.Target{Kind: engine.TargetTheOtherCreature},
	TheSameCreature:            engine.Target{Kind: engine.TargetTheSameCreature},
	TheChosenCreature:          engine.Target{Kind: engine.TargetTheChosenCreature},
	CreatureFought:             engine.Target{Kind: engine.TargetCreatureFought},
	CreatureOrArtifact:         engine.Target{Kind: engine.TargetChosenCreatureOrArtifact},
	FriendlyCreatureOrArtifact: engine.Target{Kind: engine.TargetChosenFriendlyCreatureOrArtifact},
	EnemyCreatureOrArtifact:    engine.Target{Kind: engine.TargetChosenEnemyCreatureOrArtifact},
	Artifact:                   engine.Target{Kind: engine.TargetChosenArtifact},
	Upgrade:                    engine.Target{Kind: engine.TargetChosenUpgrade},
	FriendlyArtifact:           engine.Target{Kind: engine.TargetChosenFriendlyArtifact},
	EnemyArtifact:              engine.Target{Kind: engine.TargetChosenEnemyArtifact},
	FormerNeighbors:            engine.Target{Kind: engine.TargetFormerNeighbors},
	EachNeighbor:               engine.Target{Kind: engine.TargetEachNeighbor},
	EachUpgradeOnThis:          engine.Target{Kind: engine.TargetEachUpgradeOnThis},
	TheFoughtCreature:          engine.Target{Kind: engine.TargetTheFoughtCreature},
	AttachedHost:               engine.Target{Kind: engine.TargetAttachedHost},
	GrantingCard:               engine.Target{Kind: engine.TargetGrantingCard},
}

type targets struct {
	// This selects the source card itself.
	This engine.Target
	// Triggering selects the creature that fired the trigger ("it").
	Triggering engine.Target
	// CreatureFought selects the creature the source is fighting, for a Before
	// Fight ability that names it in full.
	CreatureFought engine.Target
	// Creature is a single creature the controller chooses, either side.
	Creature engine.Target
	// FriendlyCreature is a single friendly creature the controller chooses.
	FriendlyCreature engine.Target
	// EnemyCreature is a single enemy creature the controller chooses.
	EnemyCreature engine.Target
	// EachCreature selects every creature in play.
	EachCreature engine.Target
	// EachFriendlyCreature selects every friendly creature.
	EachFriendlyCreature engine.Target
	// EachEnemyCreature selects every enemy creature.
	EachEnemyCreature engine.Target
	// EachArtifact selects every artifact in play.
	EachArtifact engine.Target
	// EachFriendlyArtifact selects every artifact the controller controls.
	EachFriendlyArtifact engine.Target
	// EachEnemyArtifact selects every artifact the opponent controls.
	EachEnemyArtifact engine.Target
	// EachFriendlyCardInPlay selects the controller's creatures and artifacts.
	EachFriendlyCardInPlay engine.Target
	// EachOtherFriendlyCreature selects the controller's creatures except the source.
	EachOtherFriendlyCreature engine.Target
	// OtherFriendlyCreature is a friendly creature the controller chooses except the source.
	OtherFriendlyCreature engine.Target
	// OtherCreature is a creature the controller chooses except the one in context (ctx.It).
	OtherCreature engine.Target
	// TheOtherCreature selects the creature in context (ctx.It), "the other creature".
	TheOtherCreature engine.Target
	// TheSameCreature selects the triggering creature (ctx.It), "the same creature".
	TheSameCreature engine.Target
	// TheChosenCreature selects the creature in context (ctx.It), "the chosen creature".
	TheChosenCreature engine.Target
	// CreatureOrArtifact is a creature or artifact the controller chooses, either side.
	CreatureOrArtifact engine.Target
	// FriendlyCreatureOrArtifact is a friendly creature or artifact the controller chooses.
	FriendlyCreatureOrArtifact engine.Target
	// EnemyCreatureOrArtifact is an enemy creature or artifact the controller chooses.
	EnemyCreatureOrArtifact engine.Target
	// Artifact is a single artifact the controller chooses, either side.
	Artifact engine.Target
	// Upgrade is a single upgrade the controller chooses from all in play.
	Upgrade engine.Target
	// FriendlyArtifact is a single friendly artifact the controller chooses.
	FriendlyArtifact engine.Target
	// EnemyArtifact is a single enemy artifact the controller chooses.
	EnemyArtifact engine.Target
	// FormerNeighbors selects the neighbors a preceding effect snapshotted before
	// removing a creature ("each of that creature's neighbors"), for a follow-up
	// that hits a destroyed creature's former neighbors (Pain Reaction).
	FormerNeighbors engine.Target
	// EachNeighbor selects the source card's live battleline neighbors ("each of
	// <self>'s neighbors") — Ghosthawk reaps with each of its neighbors.
	EachNeighbor engine.Target
	// EachUpgradeOnThis selects the upgrades attached to the source card ("each
	// upgrade on <self>") — Away Team archives its own upgrades when destroyed.
	EachUpgradeOnThis engine.Target
	// TheFoughtCreature selects the creature a preceding effect had a chosen creature
	// fight ("the fought creature"), naming no fighter — Smite makes a friendly
	// creature fight, then damages the fought creature's neighbors.
	TheFoughtCreature engine.Target
	// AttachedHost selects the creature the resolving upgrade is attached to — the
	// instance a blaster bound to when AttachSelfTo homed it, not a same-named copy.
	// Chain Named() to give it the signature creature's printed name.
	AttachedHost engine.Target
	// GrantingCard selects the in-play card whose constant ability or static modifier
	// granted the resolving ability — the exact card, not a same-named copy. Its text
	// renders the granting card's own name (Uncharted Lands).
	GrantingCard engine.Target
}

// Refinement refines a Target relative to the whole selected set (see
// MostPowerful); pass one to a target's Refine method.
type Refinement = engine.Refinement

// Power selectors come in two kinds. A tier keeps every creature tied at the
// extreme and makes no choice: HighestPower / LowestPower. A singular selector
// keeps exactly one creature and lets the controller break ties: MostPowerful /
// LeastPowerful. Compose them with the Except (complement) and AnyOf (union)
// combinators — Except(MostPowerful) spares one creature and takes the rest,
// AnyOf(LowestPower, HighestPower) takes both extremes at once.

// HighestPower is a Refinement that keeps every creature tied for the highest
// power of a set (a tier, so no choice), e.g.
// card.Target.EachCreature.Refine(card.HighestPower).
var HighestPower = engine.HighestPower

// LowestPower is a Refinement that keeps every creature tied for the lowest power
// of a set (a tier, so no choice), e.g.
// card.Target.EachCreature.Refine(card.LowestPower).
var LowestPower = engine.LowestPower

// MostPowerful is a Refinement that keeps the single most powerful creature of a
// set, the controller breaking ties, e.g.
// card.Target.EachCreature.Refine(card.MostPowerful) (Soulkeeper). For the top-N
// form use MostPowerfulN.
var MostPowerful = engine.MostPowerful

// Except returns a Refinement that keeps every creature the inner Refinement drops
// — the complement, e.g. card.Target.EachEnemyCreature.Refine(
// card.Except(card.MostPowerful)) is "each enemy creature except the most powerful"
// (Champion's Challenge).
var Except = engine.Except

// AnyOf returns a Refinement that keeps every creature any member keeps — the
// union, e.g. card.Target.EachCreature.Refine(
// card.AnyOf(card.LowestPower, card.HighestPower)) (Standardized Testing).
var AnyOf = engine.AnyOf

// SamePowerAsChosen is a Refinement that keeps every creature sharing the power of
// one the controller chooses, e.g.
// card.Target.EachCreature.Refine(card.SamePowerAsChosen) (Dance of Doom).
var SamePowerAsChosen = engine.SamePowerAsChosen

// SamePowerAsEitherChosen is a Refinement that keeps every creature sharing the
// power of a chosen friendly or enemy creature, e.g.
// card.Target.EachCreature.Refine(card.SamePowerAsEitherChosen) (Quintrino
// Flux).
var SamePowerAsEitherChosen = engine.SamePowerAsEitherChosen

// LeastPowerful is a Refinement that keeps only the single least powerful creature
// of a set, e.g. card.Target.EachCreature.Refine(card.LeastPowerful) (Horseman
// of Famine). When several tie the controller chooses which one to keep.
var LeastPowerful = engine.LeastPowerful

// KeepPerSide returns a Refinement that spares a chosen number of creatures on
// each battleline and selects every other creature, e.g.
// card.Target.EachCreature.Refine(card.KeepPerSide(3)) (Unnatural Selection).
var KeepPerSide = engine.KeepPerSide

// PortionPerSide returns a Refinement that selects a Fraction of the creatures on
// each battleline, chosen by the controller, e.g.
// card.Target.EachCreature.Refine(card.PortionPerSide(card.ThirdRoundedUp)) (Tertiate).
var PortionPerSide = engine.PortionPerSide

// MostPowerfulN returns a Refinement that keeps the n most powerful creatures of a
// set, e.g. card.Target.EachCreature.Refine(card.MostPowerfulN(3)) (Three Fates).
// When more tie at the cutoff than there are slots, the controller chooses which.
var MostPowerfulN = engine.MostPowerfulN

// HouseWithAtLeast returns a Refinement that keeps only creatures whose house has
// at least n creatures in play, counting each house across both battlelines, e.g.
// card.Target.EachCreature.Refine(card.HouseWithAtLeast(3)) (No Safety in
// Numbers).
var HouseWithAtLeast = engine.HouseWithAtLeast

// WithoutSharedTrait returns a Refinement that keeps only creatures that share no
// trait with another creature in the same controller's battleline, e.g.
// card.Target.EachCreature.Refine(card.WithoutSharedTrait()) (Good of the Many).
var WithoutSharedTrait = engine.WithoutSharedTrait

// PowerLessThan is a Refinement that keeps every creature of a set whose power is
// below a running count, e.g.
// card.Target.EachCreature.House(card.Houses.Except(card.House.Self)).Refine(card.PowerLessThan(count))
// (Exterminate! Exterminate!).
var PowerLessThan = engine.PowerLessThan

// PowerLessThanSource is a Refinement that keeps every creature whose power is
// below the source card's own power, e.g.
// card.Target.Creature.Refine(card.PowerLessThanSource()) (Dreadbone Decimus).
var PowerLessThanSource = engine.PowerLessThanSource

// PowerAtLeast is a Refinement that keeps every creature whose power reaches a
// minimum. Prefer the Target axis of the same name; reach for this one only to
// union power against another axis inside a card.AnyOf (Regrettable Meteor).
var PowerAtLeast = engine.PowerAtLeast

// OfTrait is a Refinement that keeps every creature with a trait. Prefer the
// Target axis card.Target.EachCreature.WithTrait; reach for this one only to union
// a trait against another axis inside a card.AnyOf (Regrettable Meteor).
var OfTrait = engine.OfTrait

// Stunned is the set of stunned creatures, used as a fight restriction: pass it to
// card.WithFightRestriction to limit a creature to fighting only stunned creatures
// (Bigtwig).
var Stunned = engine.Target{Kind: engine.TargetEachCreature}.Stunned()

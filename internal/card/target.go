package card

import "github.com/dmikalova/vactrol/internal/engine"

// Target groups ready-made targets, e.g. card.Target.EachEnemyCreature. Each is
// an engine.Target value, so the filter methods (WithTrait, PowerAtMost, OnFlank,
// Selector, ...) chain off them:
// card.Target.EachEnemyCreature.Selector(card.ExceptMostPowerful).
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
	CreatureFought:             engine.Target{Kind: engine.TargetCreatureFought},
	CreatureOrArtifact:         engine.Target{Kind: engine.TargetChosenCreatureOrArtifact},
	FriendlyCreatureOrArtifact: engine.Target{Kind: engine.TargetChosenFriendlyCreatureOrArtifact},
	EnemyCreatureOrArtifact:    engine.Target{Kind: engine.TargetChosenEnemyCreatureOrArtifact},
	Artifact:                   engine.Target{Kind: engine.TargetChosenArtifact},
	EnemyArtifact:              engine.Target{Kind: engine.TargetChosenEnemyArtifact},
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
	// CreatureOrArtifact is a creature or artifact the controller chooses, either side.
	CreatureOrArtifact engine.Target
	// FriendlyCreatureOrArtifact is a friendly creature or artifact the controller chooses.
	FriendlyCreatureOrArtifact engine.Target
	// EnemyCreatureOrArtifact is an enemy creature or artifact the controller chooses.
	EnemyCreatureOrArtifact engine.Target
	// Artifact is a single artifact the controller chooses, either side.
	Artifact engine.Target
	// EnemyArtifact is a single enemy artifact the controller chooses.
	EnemyArtifact engine.Target
}

// Selector refines a Target relative to the whole selected set (see
// ExceptMostPowerful); pass one to a target's Selector method.
type Selector = engine.Selector

// ExceptMostPowerful is a Selector that drops the single most powerful creature
// from a set, e.g. card.Target.EachEnemyCreature.Selector(card.ExceptMostPowerful).
// When several tie for most powerful the controller chooses which one to keep.
var ExceptMostPowerful = engine.ExceptMostPowerful

// SamePowerAsChosen is a Selector that keeps every creature sharing the power of
// one the controller chooses, e.g.
// card.Target.EachCreature.Selector(card.SamePowerAsChosen) (Dance of Doom).
var SamePowerAsChosen = engine.SamePowerAsChosen

// LeastPowerful is a Selector that keeps only the single least powerful creature
// of a set, e.g. card.Target.EachCreature.Selector(card.LeastPowerful) (Horseman
// of Famine). When several tie the controller chooses which one to keep.
var LeastPowerful = engine.LeastPowerful

// MostPowerful returns a Selector that keeps the n most powerful creatures of a
// set, e.g. card.Target.EachCreature.Selector(card.MostPowerful(3)) (Three Fates).
// When more tie at the cutoff than there are slots, the controller chooses which.
var MostPowerful = engine.MostPowerful

// Stunned is the set of stunned creatures, used as a fight restriction: pass it to
// card.WithFightRestriction to limit a creature to fighting only stunned creatures
// (Bigtwig).
var Stunned = engine.Target{Kind: engine.TargetEachCreature}.Stunned()

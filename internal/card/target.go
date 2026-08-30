package card

import "github.com/dmikalova/vactrol/internal/engine"

// Target groups ready-made targets, e.g. card.Target.EachEnemyCreature. Each is
// an engine.Target value, so the filter methods (WithTrait, PowerAtMost, OnFlank,
// Selector, ...) chain off them:
// card.Target.EachEnemyCreature.Selector(card.ExceptMostPowerful).
var Target = targets{
	This:                      engine.Target{Kind: engine.TargetThisCreature},
	Triggering:                engine.Target{Kind: engine.TargetTriggeringCreature},
	Creature:                  engine.Target{Kind: engine.TargetChosenCreature},
	FriendlyCreature:          engine.Target{Kind: engine.TargetChosenFriendlyCreature},
	EnemyCreature:             engine.Target{Kind: engine.TargetChosenEnemyCreature},
	EachCreature:              engine.Target{Kind: engine.TargetEachCreature},
	EachFriendlyCreature:      engine.Target{Kind: engine.TargetEachFriendlyCreature},
	EachEnemyCreature:         engine.Target{Kind: engine.TargetEachEnemyCreature},
	EachArtifact:              engine.Target{Kind: engine.TargetEachArtifact},
	EachFriendlyInPlay:        engine.Target{Kind: engine.TargetEachFriendlyInPlay},
	EachOtherFriendlyCreature: engine.Target{Kind: engine.TargetEachOtherFriendlyCreature},
	OtherFriendlyCreature:     engine.Target{Kind: engine.TargetChosenOtherFriendlyCreature},
	TheOtherCreature:          engine.Target{Kind: engine.TargetTheOtherCreature},
	Artifact:                  engine.Target{Kind: engine.TargetChosenArtifact},
}

type targets struct {
	This,
	Triggering,
	Creature,
	FriendlyCreature,
	EnemyCreature,
	EachCreature,
	EachFriendlyCreature,
	EachEnemyCreature,
	EachArtifact,
	EachFriendlyInPlay,
	EachOtherFriendlyCreature,
	OtherFriendlyCreature,
	TheOtherCreature,
	Artifact engine.Target
}

// Selector refines a Target relative to the whole selected set (see
// ExceptMostPowerful); pass one to a target's Selector method.
type Selector = engine.Selector

// ExceptMostPowerful is a Selector that drops the single most powerful creature
// from a set, e.g. card.Target.EachEnemyCreature.Selector(card.ExceptMostPowerful).
// When several tie for most powerful the controller chooses which one to keep.
var ExceptMostPowerful = engine.ExceptMostPowerful

// Stunned is the set of stunned creatures, used as a fight restriction: pass it to
// card.WithFightRestriction to limit a creature to fighting only stunned creatures
// (Bigtwig).
var Stunned = engine.Target{Kind: engine.TargetEachCreature}.Stunned()

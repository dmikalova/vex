package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the trigger glyph
// vocabulary that heads a glyph line, and the action-trigger grouping that merges
// adjacent Play/Fight/Reap abilities onto one line.

// isActionTrigger reports whether a trigger is one of the three action triggers
// (Play, Fight, Reap) that the Trigger.PlayFightReap composite and its pair
// variants merge onto one line.
func isActionTrigger(t engine.Trigger) bool {
	switch t {
	case engine.TriggerAfterPlay, engine.TriggerAfterFight, engine.TriggerAfterReap:
		return true
	}
	return false
}

// triggerIcon is the glyph a trigger shows at the head of its line. Every trigger
// a card uses maps to its own icon; the abstract fallback is a tripwire, not a
// shipping glyph — TestNoResidualUnknownGlyph fails if any card's strip reaches
// it, so a new trigger's icon must be added here rather than falling back.
func triggerIcon(t engine.Trigger) string {
	switch t {
	case engine.TriggerAfterPlay,
		engine.TriggerEntersPlay,
		engine.TriggerAfterCreatureEnters,
		engine.TriggerAfterCreaturePlayedAdjacent,
		engine.TriggerAfterCreaturePlayed,
		engine.TriggerAfterCardPlayed,
		engine.TriggerAfterEnemyCardPlayed,
		engine.TriggerAfterTacticPlayedBeforeResolve,
		engine.TriggerAfterUpgradeEnters:
		return "glyph-play"
	case engine.TriggerAfterReap, engine.TriggerAfterCreatureReaps:
		return "glyph-reap"
	case engine.TriggerAfterFight, engine.TriggerBeforeFight,
		engine.TriggerAfterDestroyedFighting,
		engine.TriggerAfterEnemyDestroyedFighting,
		engine.TriggerAfterCreatureFights,
		engine.TriggerAfterAssaultDestroys,
		engine.TriggerAfterNeighborFights:
		return "glyph-fight"
	case engine.TriggerAction,
		engine.TriggerAfterUse,
		engine.TriggerAfterUsedSelf:
		return "glyph-action"
	case engine.TriggerDestroyed,
		engine.TriggerLeavesPlay,
		engine.TriggerAfterCreatureDestroyed:
		return "glyph-destroyed"
	case engine.TriggerAfterForgeKey,
		engine.TriggerAfterPlayerForgesKey,
		engine.TriggerAfterOpponentForgesKey,
		engine.TriggerBeforeOpponentForgesKey:
		return "forge"
	case engine.TriggerAfterBonusDamage:
		return "damage"
	case engine.TriggerAfterBonusDraw:
		return "draw"
	case engine.TriggerAfterChooseHouse:
		return "glyph-choose"
	case engine.TriggerAfterDiscardFromHand:
		return "zone-discard"
	case engine.TriggerAfterAemberStolenFromYou:
		return "aember"
	case engine.TriggerAfterArmorPrevents:
		return "shield"
	case engine.TriggerStartOfTurn,
		engine.TriggerEndOfTurn,
		engine.TriggerEndOfReadyStep:
		return "phase-turn"
	default:
		return "glyph-unknown"
	}
}

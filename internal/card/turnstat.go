package card

import "github.com/dmikalova/vactrol/internal/engine"

// TurnStat groups the tallies the engine keeps about what a player did during a
// turn, e.g. card.TurnStat.EnemyCreaturesFightKilled (see card.TurnCount). It
// mirrors the engine's turnstat.go.
var TurnStat = turnStats{
	KeysForgedThisTurn:        engine.KeysForgedThisTurn,
	KeysForgedLastTurn:        engine.KeysForgedLastTurn,
	CreaturesPlayedLastTurn:   engine.CreaturesPlayedLastTurn,
	EnemyCreaturesFightKilled: engine.EnemyCreaturesFightKilled,
}

type turnStats struct {
	// KeysForgedThisTurn counts the keys the player forged during the current turn.
	KeysForgedThisTurn engine.TurnStat
	// KeysForgedLastTurn counts the keys the player forged on their previous turn.
	KeysForgedLastTurn engine.TurnStat
	// CreaturesPlayedLastTurn counts the creatures the player played on their
	// previous turn.
	CreaturesPlayedLastTurn engine.TurnStat
	// EnemyCreaturesFightKilled counts the player's enemies destroyed in a fight
	// this turn.
	EnemyCreaturesFightKilled engine.TurnStat
}

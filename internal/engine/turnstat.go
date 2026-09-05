package engine

// TurnStat names one of the tallies the engine keeps about what a player did
// during a turn, so a card can ask "did your opponent forge a key on their
// previous turn?" or "how many enemy creatures died fighting this turn?" without
// the engine keeping a replayable log of the game. The tallies themselves live in
// GameState.TurnHistory, which is indexed by these values; a new question about a
// turn is one more constant here plus the one place that bumps it.
type TurnStat int

const (
	// KeysForgedThisTurn counts the keys the player has forged during the current
	// turn — Smiling Ruth.
	KeysForgedThisTurn TurnStat = iota
	// KeysForgedLastTurn counts the keys the player forged during their own previous
	// turn — Tendrils of Pain, Key Hammer.
	KeysForgedLastTurn
	// CreaturesPlayedLastTurn counts the creatures the player played during their own
	// previous turn — Lifeweb.
	CreaturesPlayedLastTurn
	// EnemyCreaturesFightKilled counts the player's enemies destroyed in a fight this
	// turn — The Warchest. It is the only tally kept from the watching player's side
	// rather than the acting player's, because that is the side the card pays.
	EnemyCreaturesFightKilled
	// EnemyCreaturesDestroyed counts the player's enemy creatures destroyed by any
	// means this turn — Foozle. Like EnemyCreaturesFightKilled it is kept from the
	// watching player's side, but it counts every destruction, not only fights.
	EnemyCreaturesDestroyed
	// turnStatCount sizes GameState.TurnHistory and is not a tally itself.
	turnStatCount
)

// turnStatNoun is the singular noun each tally repeats after "for each".
var turnStatNoun = map[TurnStat]string{
	KeysForgedThisTurn:        "key you have forged this turn",
	KeysForgedLastTurn:        "key forged on the previous turn",
	CreaturesPlayedLastTurn:   "creature played on the previous turn",
	EnemyCreaturesFightKilled: "enemy creature that was destroyed in a fight this turn",
	EnemyCreaturesDestroyed:   "enemy creature that was destroyed this turn",
}

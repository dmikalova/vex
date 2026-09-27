package engine

import "errors"

// This file holds the Alpha keyword: a card with Alpha can only be played as the
// first thing its player does on their turn. The play gates in game_play.go call
// barredByAlpha before letting an Alpha card be played.

// ErrAlphaNotFirst reports an Alpha card played after its player has already acted
// this turn, so it can no longer be the first thing they do.
var ErrAlphaNotFirst = errors.New("an Alpha card must be the first card played this turn")

// actedThisTurn reports whether player has already played, used, or discarded a
// card this turn — the "anything else this step" an Alpha card must come before.
// The scan is row-only: TimesUsedThisTurn is per-card row state, and an upgrade
// is never used, so a widened scan would have nothing to ask it
// (docs/adr/0047-upgrade-in-play-not-an-ability-source.md).
func (g *Game) actedThisTurn(player int) bool {
	if g.State.PlayedThisTurn[player].Count > 0 ||
		g.State.DiscardedThisTurn[player].Count > 0 {
		return true
	}
	for _, id := range g.creaturesAndArtifacts(player) {
		if g.State.Cards[id].TimesUsedThisTurn > 0 {
			return true
		}
	}
	return false
}

// barredByAlpha reports whether def carries Alpha and player has already acted
// this turn, so the Alpha card can no longer be the first thing they do.
func (g *Game) barredByAlpha(player int, def *CardDefinition) bool {
	return def.hasKeyword(Alpha) && g.actedThisTurn(player)
}

package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the turn-and-house
// family (ADR 0047) — changing the active house, ending the turn, naming or
// choosing houses, restricting an action, chains, and the bonus-icon resolution
// a forged key's flip triggers.

// turnEffectGlyphs claims every turn-structure and house-choice effect.
func turnEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.ChangeActiveHouse:
		return []glyph{{asset: "glyph-choose"}}, true, true
	case engine.EndTurn:
		return []glyph{{asset: "phase-turn"}}, true, true
	case engine.MustChooseHouse:
		return []glyph{{asset: "glyph-choose", decor: playerDecor(v.Player)}}, true, true
	case engine.CannotChooseHouse:
		return []glyph{
			{asset: "glyph-choose", decor: playerDecor(v.Player)},
			{asset: "glyph-ban"},
		}, true, true
	case engine.NameHouse:
		// The chosen house is barred; ChooseHouseThen supplies the choose glyph.
		return []glyph{{asset: "glyph-ban"}}, true, true
	case engine.NameCard:
		// Name a card, then bar every copy of it from being played.
		return []glyph{{asset: "glyph-choose"}, {asset: "glyph-ban"}}, true, true
	case engine.OpponentNamesHouse:
		return []glyph{{asset: "glyph-choose"}}, true, true
	case engine.WagerOpponentChoosesChosenHouse:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorEnemy | decorChosen},
		}, true, true
	case engine.Restrict:
		verb := "glyph-reap"
		switch v.Action {
		case engine.RestrictFighting:
			verb = "glyph-fight"
		case engine.RestrictUse:
			verb = "glyph-action"
		}
		// Restrict bars a player, where CreaturesCannot bars the creatures on the
		// board; the subject decor is what tells the two apart (ADR 0047).
		return []glyph{
			{asset: verb, decor: playerDecor(v.Player)},
			{asset: "glyph-ban"},
		}, true, true
	case engine.CreaturesCannot:
		action := "glyph-reap"
		if v.Action == engine.FightUse {
			action = "glyph-fight"
		}
		return []glyph{
			{asset: action, decor: decorEach},
			{asset: "glyph-ban"},
		}, true, true
	case engine.GainChains:
		return []glyph{{asset: "chains", qty: v.Amount, decor: playerDecor(v.Player)}}, true, true
	case engine.ResolveBonusIcons:
		return []glyph{
			targetGlyph(v.Target),
			arrowTo(glyph{asset: "aember"}),
			{asset: "capture"},
			{asset: "damage"},
			{asset: "draw"},
		}, true, true
	case engine.ExtraBonusIconResolution:
		return []glyph{
			{asset: "glyph-play"},
			arrowTo(glyph{asset: "aember"}),
			{asset: "capture"},
			{asset: "damage"},
			{asset: "draw"},
		}, true, true
	default:
		return nil, false, false
	}
}

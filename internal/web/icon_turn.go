package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the turn-and-house
// family (ADR 0047) — changing the active house, ending the turn, naming or
// choosing houses, restricting an action, chains, and the bonus-icon resolution
// a forged key's flip triggers.

// turnEffectGlyphs claims every turn-structure and house-choice effect.
func turnEffectGlyphs(e engine.Effect) ([]glyph, bool) {
	switch v := e.(type) {
	case engine.ChangeActiveHouse:
		// The house it switches to, then the turn marker the active house is already
		// drawn by (houseMatcherGlyphs): this house becomes this turn's house. The
		// choose mark is left to the effects that ask a player to pick one.
		return append(
			houseChoiceGlyphs(v.To),
			arrowTo(glyph{asset: "phase-turn"}),
		), true
	case engine.EndTurn:
		return []glyph{{asset: "phase-turn"}}, true
	case engine.MustChooseHouse:
		return []glyph{{asset: "glyph-choose", decor: playerDecor(v.Player)}}, true
	case engine.CannotChooseHouse:
		return []glyph{
			{asset: "glyph-choose", decor: playerDecor(v.Player)},
			{asset: "glyph-ban"},
		}, true
	case engine.NameHouse:
		// The chosen house is barred; ChooseHouseThen supplies the choose glyph.
		return []glyph{{asset: "glyph-ban"}}, true
	case engine.NameCard:
		// Name a card, then bar every copy of it from being played.
		return []glyph{{asset: "glyph-choose"}, {asset: "glyph-ban"}}, true
	case engine.OpponentNamesHouse:
		// A house the opponent picks, so the house noun leads a choose mark tinted the
		// way every other enemy subject is; MustChooseHouse above is the bare tinted
		// mark for whichever player it names.
		return []glyph{
			{asset: "house-any"},
			{asset: "glyph-choose", decor: decorEnemy},
		}, true
	case engine.WagerOpponentChoosesChosenHouse:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorEnemy | decorChosen},
		}, true
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
		}, true
	case engine.CreaturesCannot:
		action := "glyph-reap"
		if v.Action == engine.FightUse {
			action = "glyph-fight"
		}
		return []glyph{
			{asset: action, decor: decorEach},
			{asset: "glyph-ban"},
		}, true
	case engine.GainChains:
		return []glyph{{asset: "chains", qty: v.Amount, decor: playerDecor(v.Player)}}, true
	case engine.ResolveBonusIcons:
		return append(
			targetGlyphs(v.Target),
			arrowTo(glyph{asset: "aember"}),
			glyph{asset: "capture"},
			glyph{asset: "damage"},
			glyph{asset: "draw"},
		), true
	case engine.ExtraBonusIconResolution:
		return []glyph{
			{asset: "glyph-play"},
			arrowTo(glyph{asset: "aember"}),
			{asset: "capture"},
			{asset: "damage"},
			{asset: "draw"},
		}, true
	default:
		return nil, false
	}
}

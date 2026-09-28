package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the continuous-rules
// renderers — a card's keywords, static bonuses, restrictions, tolls, key-cost
// and Æmber-flow replacements, and card-level feature flags — as glyph lines.

// keywordGlyphs is the one line of a card's static keywords and combat keywords —
// Skirmish, Elusive, Assault 2, Hazardous 3 — so a vanilla creature whose only
// mechanics are keywords still draws an Icon strip. Numeric combat keywords carry
// their value as the glyph's numeral.
func keywordGlyphs(def *engine.CardDefinition) []glyph {
	gs := make([]glyph, 0, len(def.Keywords)+2)
	for _, k := range def.Keywords {
		if a := keywordIcon(k); a != "" {
			gs = append(gs, glyph{asset: a})
		}
	}
	if def.Assault != 0 {
		gs = append(gs, glyph{
			asset: "kw-assault",
			qty:   def.Assault,
		})
	}
	if def.Hazardous != 0 {
		gs = append(gs, glyph{
			asset: "kw-hazardous",
			qty:   def.Hazardous,
		})
	}
	if def.SplashAttack != 0 {
		gs = append(gs, glyph{
			asset: "damage",
			qty:   def.SplashAttack,
			decor: decorEach,
		})
	}
	return gs
}

// keywordIcon is the glyph asset for a printed keyword, or "" for keywords with no
// drawn glyph yet.
func keywordIcon(k engine.Keyword) string {
	switch k {
	case engine.Skirmish:
		return "kw-skirmish"
	case engine.Poison:
		return "kw-poison"
	case engine.Elusive:
		return "kw-elusive"
	case engine.Taunt:
		return "kw-taunt"
	case engine.Versatile:
		return "kw-versatile"
	case engine.Alpha:
		return "kw-alpha"
	case engine.Omega:
		return "kw-omega"
	case engine.Deploy:
		return "kw-deploy"
	}
	return ""
}

// constantLines renders a constant ability's lines: one for its stat bonuses
// ("Each friendly creature gains +1 power") arrowed to the creatures it reaches,
// and one per triggered ability it grants ("… gains Destroyed: purge this
// creature"), so a card whose only mechanic is a granted ability still draws a
// strip.
func constantLines(ca engine.ConstantAbility) []glyphLine {
	target := ca.Target
	if target == (engine.Target{}) {
		target = engine.Target{Kind: engine.TargetEachCardInPlay}
	}
	var lines []glyphLine
	if ca.PowerBonus != 0 || ca.ArmorBonus != 0 || len(ca.Keywords) > 0 {
		gs := make([]glyph, 0, 4)
		if ca.PowerBonus != 0 {
			gs = append(gs, glyph{
				asset: "power",
				qty:   ca.PowerBonus,
			})
		}
		if ca.ArmorBonus != 0 {
			gs = append(gs, glyph{
				asset: "shield",
				qty:   ca.ArmorBonus,
			})
		}
		for _, k := range ca.Keywords {
			if a := keywordIcon(k); a != "" {
				gs = append(gs, glyph{asset: a})
			}
		}
		gs = append(gs, arrowTo(targetGlyph(target)))
		lines = append(lines, glyphLine{
			glyphs:  gs,
			covered: true,
		})
	}
	for _, gr := range ca.Granted {
		gs, covered := effectGlyphs(gr.Effect)
		lines = append(lines, glyphLine{
			triggers: []string{triggerIcon(gr.Trigger)},
			glyphs:   gs,
			covered:  covered,
		})
	}
	return lines
}

// staticLines renders an Upgrade's continuous modifier (def.Static) as glyph lines
// describing what it grants the creature it is attached to: one line of stat and
// keyword bonuses, and one line per triggered ability it grants. An Upgrade that
// only buffs its host carries every mechanic here, so without this pass such a
// card — Blood of Titans, Backup Copy, Bonerot Venom — drew an empty strip. The
// host creature is implicit, so the bonus line takes no arrowed target.
func staticLines(m engine.StaticModifier) []glyphLine {
	var lines []glyphLine
	gs := make([]glyph, 0, 6)
	if m.PowerBonus != 0 {
		gs = append(gs, glyph{
			asset: "power",
			qty:   m.PowerBonus,
		})
	}
	if m.ArmorBonus != 0 {
		gs = append(gs, glyph{
			asset: "shield",
			qty:   m.ArmorBonus,
		})
	}
	if m.AssaultBonus != 0 {
		gs = append(gs, glyph{
			asset: "kw-assault",
			qty:   m.AssaultBonus,
		})
	}
	if m.HazardousBonus != 0 {
		gs = append(gs, glyph{
			asset: "kw-hazardous",
			qty:   m.HazardousBonus,
		})
	}
	if m.SplashAttackBonus != 0 {
		gs = append(gs, glyph{
			asset: "damage",
			qty:   m.SplashAttackBonus,
			decor: decorEach,
		})
	}
	for _, k := range m.Keywords {
		if a := keywordIcon(k); a != "" {
			gs = append(gs, glyph{asset: a})
		}
	}
	if m.AemberCannotBeStolen != nil {
		gs = append(gs,
			glyph{
				asset: "aember",
				decor: decorEnemy,
			}, glyph{asset: "glyph-ban"})
	}
	if m.ProtectsFromNonFlank {
		gs = append(gs, glyph{asset: "glyph-flank"}, glyph{asset: "shield"})
	}
	if a := houseIconName(m.HouseOverride); a != "" {
		gs = append(gs, glyph{asset: a})
	}
	if len(gs) > 0 {
		lines = append(lines, glyphLine{
			glyphs:  gs,
			covered: true,
		})
	}
	for _, gr := range m.Granted {
		egs, covered := effectGlyphs(gr.Effect)
		lines = append(lines, glyphLine{
			triggers: []string{triggerIcon(gr.Trigger)},
			glyphs:   egs,
			covered:  covered,
		})
	}
	return lines
}

// restrictionLines renders the continuous "cannot" rules a card imposes while in
// play (def.Restricts) as one line of banned actions, so a card whose only mechanic
// is a restriction — Barrister Joya's "enemy creatures cannot reap", Ember Imp's
// play cap — still draws a strip. Each restriction is its action glyph struck by
// the ban glyph; the finer qualifier (which player, which condition) stays in the
// rules text.
func restrictionLines(r engine.Restrictions) []glyphLine {
	gs := make([]glyph, 0, 4)
	if r.Fighting {
		gs = append(gs, glyph{asset: "glyph-fight"}, glyph{asset: "glyph-ban"})
	}
	if playerSet(r.Reaping) {
		gs = append(gs,
			glyph{
				asset: "glyph-reap",
				decor: playerDecor(r.Reaping),
			},
			glyph{asset: "glyph-ban"})
	}
	if a := typeIconName(r.CannotPlay); a != "" {
		gs = append(gs, glyph{asset: a}, glyph{asset: "glyph-ban"})
	}
	if r.PlayCardLimit.Amount != 0 {
		gs = append(gs,
			glyph{
				asset: "glyph-play",
				qty:   r.PlayCardLimit.Amount,
				decor: playerDecor(r.PlayCardLimit.Player),
			},
			glyph{asset: "glyph-ban"})
	}
	if r.UseCondition != nil {
		gs = append(gs, glyph{asset: "glyph-action"}, glyph{asset: "glyph-ban"})
	}
	if r.SkipForge || r.NoForgeWhileAheadOnKeys {
		gs = append(gs, glyph{asset: "forge"}, glyph{asset: "glyph-ban"})
	}
	if r.NoForgeKeyNumber != 0 {
		gs = append(gs,
			glyph{
				asset: "forge",
				qty:   r.NoForgeKeyNumber,
			}, glyph{asset: "glyph-ban"})
	}
	if r.MustFightIfAble {
		gs = append(gs, glyph{asset: "glyph-fight"})
	}
	if r.Toll.Amount != 0 {
		action := "type-artifact"
		if r.Toll.Action == engine.TollUseArtifact {
			action = "glyph-action"
		}
		gs = append(gs,
			glyph{asset: action},
			glyph{
				asset: "aember",
				qty:   r.Toll.Amount,
				decor: decorEnemy,
			})
	}
	if len(gs) == 0 {
		return nil
	}
	return []glyphLine{{glyphs: gs, covered: true}}
}

// cardFeatureLines renders the card-level continuous rules that live directly on
// the definition rather than in a Static, Restrictions, or Replaces value — the
// bool flags a card sets while it is in play — so a creature whose only mechanic
// is one of them still draws a strip. AemberCannotBeStolen shields the
// controller's Æmber from the enemy; DealsNoDamageWhenAttacked bars its
// retaliation; EntersReadyGrant makes friendly cards of a type enter unexhausted.
func cardFeatureLines(def *engine.CardDefinition) []glyphLine {
	gs := make([]glyph, 0, 4)
	if def.AemberCannotBeStolen != nil {
		gs = append(gs,
			glyph{
				asset: "aember",
				decor: decorEnemy,
			}, glyph{asset: "glyph-ban"})
	}
	if def.DealsNoDamageWhenAttacked {
		gs = append(gs, glyph{asset: "damage"}, glyph{asset: "glyph-ban"})
	}
	if def.CannotBeDealtDamageBy.Narrows() {
		gs = append(gs, glyph{asset: "shield"}, glyph{asset: "damage"}, glyph{asset: "glyph-ban"})
	}
	if a := typeIconName(def.EntersReadyGrant.Type); a != "" {
		gs = append(gs,
			glyph{
				asset: "exhausted",
				decor: decorFriendly,
			},
			glyph{asset: "glyph-ban"},
			arrowTo(glyph{
				asset: a,
				decor: decorFriendly,
			}))
	}
	if len(gs) == 0 {
		return nil
	}
	return []glyphLine{{glyphs: gs, covered: true}}
}

// keyCostChangeGlyphs renders a card's continuous key-cost change (def.KeyCostChanges)
// as the forge glyph beside the Æmber glyph — "keys cost more/less Æmber". The
// amount and which player it hits are scaled and unexported on the change, so they
// stay in the rules text, as with any fine filter.
func keyCostChangeGlyphs() []glyph {
	return []glyph{{asset: "forge"}, {asset: "aember"}}
}

// replaceGlyphs renders a card's continuous replacement of an Æmber-flow event
// (def.Replaces) — Ether Spider capturing the Æmber added to its opponent's pool,
// Po's Pixies drawing a steal from the common supply instead of its own pool — as
// the affected pool's Æmber swapped for the outcome the replacement substitutes.
func replaceGlyphs(r engine.Instead) []glyph {
	src := glyph{
		asset: "aember",
		decor: playerDecor(r.Player),
	}
	var out glyph
	switch r.With {
	case engine.Capture:
		out = glyph{
			asset: "aember",
			decor: decorEnemy,
		}
	case engine.Steal:
		out = glyph{
			asset: "aember",
			decor: decorEnemy | decorChosen,
		}
	case engine.FromCommonSupply:
		out = glyph{asset: "glyph-return"}
	default:
		out = glyph{asset: "glyph-swap"}
	}
	return []glyph{src, {asset: "glyph-swap"}, arrowTo(out)}
}

// playerSet reports whether p names a real player rather than the unset zero value,
// so a restriction's relative player reads as set only when it is.
func playerSet(p engine.Player) bool {
	return p == engine.Controller || p == engine.Opponent || p == engine.EachPlayer
}

// fightRestrictionGlyphs renders a creature's fight restriction — the creatures it
// is limited to fighting — as a fight glyph arrowed to the creature glyph. Every
// fight restriction is about creatures, so the noun is fixed and the restriction
// carries no set of its own; the qualifier that narrows it (stunned, flank) stays
// in the rules text, as with any fine filter.
func fightRestrictionGlyphs() []glyph {
	return []glyph{
		{asset: "glyph-fight"},
		arrowTo(glyph{
			asset: "type-creature",
			decor: decorEach,
		}),
	}
}

// drawModifierGlyphs renders a card's continuous end-of-turn hand-refill change
// (Mother draws +1, Succubus makes the opponent draw −1) as a hand glyph carrying
// the signed amount, tinted to the player it affects. EachPlayer stays untinted.
func drawModifierGlyphs(m engine.DrawModifier) []glyph {
	g := glyph{
		asset: "zone-hand",
		qty:   m.Amount,
	}
	switch m.Player {
	case engine.Controller:
		g.decor = decorFriendly
	case engine.Opponent:
		g.decor = decorEnemy
	}
	return []glyph{g}
}

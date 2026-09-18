package web

import (
	"fmt"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is the Iconography pass (ADR 0022): the Visitor that turns a card's
// mechanics into the composed Glyphs of its Icon strip, the visual counterpart of
// rules-text generation. It type-switches over the engine's exported Effect AST
// and Targets. It lives here in the client, not in package engine, because the
// engine is held at 100% coverage and an in-engine Visitor would force an icon
// test for every one of the ~150 effect nodes; internal/web is ungated. The
// totality test (icon_test.go) still walks every card and fails loud on any
// effect that falls back to the abstract unknown glyph. There is no allowlist to
// exempt one: an unmapped mechanic is a bug, and the fix is always to draw a
// glyph, so a new mechanic cannot ship without a glyph decision.

// decor is a set of edge/overlay treatments applied to a noun glyph to carry a
// Target's shape without spending a horizontal slot: an enemy tint, a friendly
// tint, an "each" stack, a "chosen" outline, or a "this/self" marker.
type decor uint8

const (
	decorEnemy decor = 1 << iota
	decorFriendly
	decorEach
	decorChosen
	decorThis
)

// glyph is one composed icon in the strip. A glyph with an asset renders that SVG
// (tinted and decorated); a glyph with no asset renders its text as a small chip,
// which is how a not-yet-drawn mechanic still shows something readable. Qty, when
// positive, prints a numeral badge on the glyph (the "3" on 3 damage).
type glyph struct {
	asset string
	text  string
	qty   int
	decor decor
	arrow bool // render a leading result-gate arrow (→) before this glyph
}

// glyphLine is one ability's transcription: its trigger glyph(s) and the glyphs
// its effect composes into. Everything is an icon — no words — so triggers render
// as their own glyph too (ADR 0022). Adjacent abilities that share one effect and
// fire on distinct action triggers (Play/Fight/Reap) merge into one line carrying
// all their trigger glyphs, mirroring the "Play/Fight/Reap:" rules-text shorthand.
type glyphLine struct {
	triggers []string // trigger glyph asset stems, in canonical order
	glyphs   []glyph
	covered  bool // false when the effect fell back to the abstract unknown glyph
}

// cardGlyphs is the Iconography pass: it renders a card's triggered abilities and
// its continuous rules — keywords, static bonuses, restrictions, tolls, key-cost
// and Æmber-flow replacements, and card-level flags — as glyph lines, so a card
// whose only mechanic is a continuous rule still draws a strip.
func cardGlyphs(def *engine.CardDefinition) []glyphLine {
	lines := make([]glyphLine, 0, len(def.Abilities)+len(def.ConstantAbilities)+2)
	if kw := keywordGlyphs(def); len(kw) > 0 {
		lines = append(lines, glyphLine{glyphs: kw, covered: true})
	}
	if def.FightRestriction != (engine.Target{}) {
		lines = append(lines, glyphLine{
			glyphs:  fightRestrictionGlyphs(def.FightRestriction),
			covered: true,
		})
	}
	if def.DrawModifier.Amount != 0 {
		lines = append(lines, glyphLine{
			glyphs:  drawModifierGlyphs(def.DrawModifier),
			covered: true,
		})
	}
	lines = append(lines, staticLines(def.Static)...)
	lines = append(lines, restrictionLines(def.Restricts)...)
	if def.Replaces != (engine.Instead{}) {
		lines = append(lines, glyphLine{glyphs: replaceGlyphs(def.Replaces), covered: true})
	}
	if len(def.KeyCostChanges) > 0 {
		lines = append(lines, glyphLine{glyphs: keyCostChangeGlyphs(), covered: true})
	}
	lines = append(lines, cardFeatureLines(def)...)
	for i := 0; i < len(def.Abilities); {
		ab := def.Abilities[i]
		gs, covered := effectGlyphs(ab.Effect)
		triggers := []string{triggerIcon(ab.Trigger)}
		j := i + 1
		if isActionTrigger(ab.Trigger) {
			text := ab.Effect.Text()
			for j < len(def.Abilities) &&
				isActionTrigger(def.Abilities[j].Trigger) &&
				def.Abilities[j].Effect.Text() == text {
				triggers = append(triggers, triggerIcon(def.Abilities[j].Trigger))
				j++
			}
		}
		lines = append(lines, glyphLine{
			triggers: triggers,
			glyphs:   gs,
			covered:  covered,
		})
		i = j
	}
	for _, ca := range def.ConstantAbilities {
		lines = append(lines, constantLines(ca)...)
	}
	return lines
}

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
		gs = append(gs, glyph{asset: "kw-assault", qty: def.Assault})
	}
	if def.Hazardous != 0 {
		gs = append(gs, glyph{asset: "kw-hazardous", qty: def.Hazardous})
	}
	if def.SplashAttack != 0 {
		gs = append(gs, glyph{asset: "damage", qty: def.SplashAttack, decor: decorEach})
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
			gs = append(gs, glyph{asset: "power", qty: ca.PowerBonus})
		}
		if ca.ArmorBonus != 0 {
			gs = append(gs, glyph{asset: "shield", qty: ca.ArmorBonus})
		}
		for _, k := range ca.Keywords {
			if a := keywordIcon(k); a != "" {
				gs = append(gs, glyph{asset: a})
			}
		}
		gs = append(gs, arrowTo(targetGlyph(target)))
		lines = append(lines, glyphLine{glyphs: gs, covered: true})
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
		gs = append(gs, glyph{asset: "power", qty: m.PowerBonus})
	}
	if m.ArmorBonus != 0 {
		gs = append(gs, glyph{asset: "shield", qty: m.ArmorBonus})
	}
	if m.AssaultBonus != 0 {
		gs = append(gs, glyph{asset: "kw-assault", qty: m.AssaultBonus})
	}
	if m.HazardousBonus != 0 {
		gs = append(gs, glyph{asset: "kw-hazardous", qty: m.HazardousBonus})
	}
	if m.SplashAttackBonus != 0 {
		gs = append(gs, glyph{asset: "damage", qty: m.SplashAttackBonus, decor: decorEach})
	}
	for _, k := range m.Keywords {
		if a := keywordIcon(k); a != "" {
			gs = append(gs, glyph{asset: a})
		}
	}
	if m.AemberCannotBeStolen != nil {
		gs = append(gs,
			glyph{asset: "aember", decor: decorEnemy}, glyph{asset: "glyph-ban"})
	}
	if m.ProtectsFromNonFlank {
		gs = append(gs, glyph{asset: "glyph-flank"}, glyph{asset: "shield"})
	}
	if a := houseIconName(m.HouseOverride); a != "" {
		gs = append(gs, glyph{asset: a})
	}
	if len(gs) > 0 {
		lines = append(lines, glyphLine{glyphs: gs, covered: true})
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
			glyph{asset: "glyph-reap", decor: playerDecor(r.Reaping)},
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
			glyph{asset: "forge", qty: r.NoForgeKeyNumber}, glyph{asset: "glyph-ban"})
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
			glyph{asset: "aember", qty: r.Toll.Amount, decor: decorEnemy})
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
			glyph{asset: "aember", decor: decorEnemy}, glyph{asset: "glyph-ban"})
	}
	if def.DealsNoDamageWhenAttacked {
		gs = append(gs, glyph{asset: "damage"}, glyph{asset: "glyph-ban"})
	}
	if def.CannotBeDealtDamageBy.Active() {
		gs = append(gs, glyph{asset: "shield"}, glyph{asset: "damage"}, glyph{asset: "glyph-ban"})
	}
	if a := typeIconName(def.EntersReadyGrant.Type); a != "" {
		gs = append(gs,
			glyph{asset: "exhausted", decor: decorFriendly},
			glyph{asset: "glyph-ban"},
			arrowTo(glyph{asset: a, decor: decorFriendly}))
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
	src := glyph{asset: "aember", decor: playerDecor(r.Player)}
	var out glyph
	switch r.With {
	case engine.Capture:
		out = glyph{asset: "aember", decor: decorEnemy}
	case engine.Steal:
		out = glyph{asset: "aember", decor: decorEnemy | decorChosen}
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
// is limited to fighting — as a fight glyph arrowed to that noun. The qualifier
// that narrows the set (stunned, damaged) stays in the rules text, as with any
// fine target filter.
func fightRestrictionGlyphs(fr engine.Target) []glyph {
	return []glyph{{asset: "glyph-fight"}, arrowTo(targetGlyph(fr))}
}

// drawModifierGlyphs renders a card's continuous end-of-turn hand-refill change
// (Mother draws +1, Succubus makes the opponent draw −1) as a hand glyph carrying
// the signed amount, tinted to the player it affects. EachPlayer stays untinted.
func drawModifierGlyphs(m engine.DrawModifier) []glyph {
	g := glyph{asset: "zone-hand", qty: m.Amount}
	switch m.Player {
	case engine.Controller:
		g.decor = decorFriendly
	case engine.Opponent:
		g.decor = decorEnemy
	}
	return []glyph{g}
}

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
	case engine.TriggerAfterChooseHouse,
		engine.TriggerAfterAnyPlayerChoosesHouse:
		return "glyph-choose"
	case engine.TriggerAfterDiscardFromHand:
		return "zone-discard"
	case engine.TriggerAfterAemberStolenFromYou:
		return "aember"
	case engine.TriggerAfterArmorPrevents:
		return "shield"
	case engine.TriggerStartOfTurn,
		engine.TriggerEndOfTurn,
		engine.TriggerEndOfReadyStep,
		engine.TriggerAfterAnyPlayerStartOfTurn,
		engine.TriggerAfterAnyPlayerEndOfTurn:
		return "phase-turn"
	default:
		return "glyph-unknown"
	}
}

// effectGlyphs renders one effect to its glyphs, reporting whether the effect was
// covered by a real mapping (false means it fell back to a text chip). It handles
// the common leaf effects; composite and rarer effects fall back for now and are
// tracked by the totality test's allowlist.
func effectGlyphs(e engine.Effect) ([]glyph, bool) {
	switch v := e.(type) {
	case engine.DealDamage:
		// A follow-up on the damaged creature draws the two clauses joined, not the
		// bare damage glyph.
		if v.Then != nil {
			return damageThenGlyphs(v.Amount, v.Target, v.Then)
		}
		// A Spread carries its own creature targets rather than filling Target, so it
		// draws its own summary noun; the counts and neighbor split stay in the text.
		if v.Spread != nil {
			return []glyph{{asset: "damage"}, arrowTo(spreadTargetGlyph())}, true
		}
		// The per-count variants hit several creatures or scale by a board count; the
		// numeral lives in the text, so the glyph drops the qty.
		if v.Per != nil || v.PerTarget != nil || v.AmountFrom != nil {
			return []glyph{{asset: "damage"}, arrowTo(targetGlyph(v.Target))}, true
		}
		return []glyph{
			{asset: "damage", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.ForEachHouse:
		return effectGlyphs(v.Do)
	case engine.GainAember:
		if v.Per != nil || v.EqualTo != nil {
			return []glyph{{asset: "aember", decor: playerDecor(v.Player)}}, true
		}
		return []glyph{{asset: "aember", qty: v.Amount, decor: playerDecor(v.Player)}}, true
	case engine.LoseAember:
		return []glyph{{asset: "aember", text: "−", decor: playerDecor(v.Player)}}, true
	case engine.StealAember:
		return []glyph{{asset: "aember", qty: v.Amount, decor: decorEnemy | decorChosen}}, true
	case engine.CaptureAember:
		return []glyph{{asset: "aember", qty: v.Amount, decor: decorEnemy}}, true
	case engine.CaptureFromAnyPlayer:
		return []glyph{{asset: "aember", qty: v.Amount}}, true
	case engine.DistributeCapture:
		return []glyph{{asset: "aember", decor: decorEnemy}}, true
	case engine.GiveAember:
		src := glyph{asset: "aember", decor: decorEnemy}
		if !v.All {
			src.qty = v.Amount
		}
		return []glyph{src, arrowTo(glyph{asset: "aember"})}, true
	case engine.GainChains:
		return []glyph{{asset: "chains", qty: v.Amount, decor: playerDecor(v.Player)}}, true
	case engine.Draw:
		return []glyph{{asset: "zone-hand", qty: v.Amount}}, true
	case engine.Stun:
		return []glyph{{asset: "stun"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.Enrage:
		return []glyph{{asset: "glyph-fight"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.Ward:
		return []glyph{{asset: "shield"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.RemoveWard:
		return []glyph{
			{asset: "shield"},
			{asset: "glyph-ban"},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.Exhaust:
		return []glyph{{asset: "exhausted"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.Ready:
		return []glyph{
			{asset: "exhausted", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.Unstun:
		return []glyph{
			{asset: "stun", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.Destroy:
		return []glyph{{asset: "glyph-destroy"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.PurgeCreature:
		return []glyph{{asset: "zone-purge"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.Heal:
		return []glyph{
			{asset: "glyph-heal", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.Exalt:
		return []glyph{
			{asset: "aember", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.ForgeKey:
		return []glyph{{asset: "forge"}}, true
	case engine.PlaceCounter:
		return []glyph{{asset: counterAsset(v.Kind)}, arrowTo(targetGlyph(v.Target))}, true
	case engine.RemoveCounters:
		return []glyph{{asset: counterAsset(v.Kind)}, arrowTo(targetGlyph(v.Target))}, true
	case engine.BlankEnemyText:
		return []glyph{{asset: "type-creature", decor: decorEnemy | decorEach}}, true
	case engine.ArchiveCard:
		return []glyph{{asset: "zone-archives", qty: v.Amount}}, true
	case engine.ArchiveFromPlay:
		return []glyph{{asset: "zone-archives"}}, true
	case engine.ArchiveSource:
		return []glyph{{asset: "zone-archives", decor: decorThis}}, true
	case engine.ArchiveGrantingUpgrade:
		return []glyph{
			{asset: "card-back", decor: decorThis},
			arrowTo(glyph{asset: "zone-archives"}),
		}, true
	case engine.ArchivePurgedCard:
		return []glyph{{asset: "zone-purge"}, arrowTo(glyph{asset: "zone-archives"})}, true
	case engine.DiscardCard:
		return []glyph{{asset: "zone-discard"}}, true
	case engine.PurgeCard:
		return []glyph{{asset: "zone-purge"}}, true
	case engine.PurgeArchives:
		return []glyph{
			{asset: "zone-archives"},
			arrowTo(glyph{asset: "zone-purge"}),
		}, true
	case engine.PurgeArchivedCardThen:
		gs := []glyph{{asset: "zone-purge"}}
		more, _ := effectGlyphs(v.Then)
		return append(gs, more...), true
	case engine.Shuffle:
		return []glyph{{asset: "zone-deck"}}, true
	case engine.ShuffleFriendlyCardsIntoDeck:
		return []glyph{{asset: "zone-deck"}}, true
	case engine.ShuffleFromDiscard:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-deck"})}, true
	case engine.ShuffleNamedFromDiscardIntoDeck:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-deck"})}, true
	case engine.PutNamedIntoHand:
		return []glyph{{asset: "glyph-return"}}, true
	case engine.PutItIntoHand:
		return []glyph{{asset: "glyph-return"}}, true
	case engine.PlayFrom, engine.PlayTopOfDeck, engine.PutIntoPlay:
		return []glyph{{asset: "glyph-play"}}, true
	case engine.PlayFromOpponent:
		zone := "zone-deck"
		if v.From == engine.Archives {
			zone = "zone-archives"
		}
		return []glyph{
			{asset: zone, decor: decorEnemy},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true
	case engine.PlayItFromOpponentDiscard:
		return []glyph{
			{asset: "zone-discard", decor: decorEnemy},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true
	case engine.DiscardOpponentArchivesOrDeckTop:
		return []glyph{
			{asset: "zone-archives", decor: decorEnemy},
			arrowTo(glyph{asset: "zone-discard", decor: decorEnemy}),
		}, true
	case engine.PutFromPlay:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: a})}, true
		}
		return fallbackGlyphs(e), false
	case engine.ExhaustCreatures:
		return []glyph{{asset: "exhausted"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.Use:
		return []glyph{{asset: "glyph-action"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.RepeatedFight:
		return []glyph{{asset: "glyph-fight"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.GainStats:
		gs := make([]glyph, 0, 3)
		if v.Power != 0 {
			gs = append(gs, glyph{asset: "power", qty: v.Power})
		}
		if v.Armor != 0 {
			gs = append(gs, glyph{asset: "shield", qty: v.Armor})
		}
		return append(gs, arrowTo(targetGlyph(v.Target))), true
	case engine.OverrideStats:
		gs := make([]glyph, 0, 2)
		if v.HasPower {
			gs = append(gs, glyph{asset: "power", qty: v.Power})
		}
		if v.HasArmor {
			gs = append(gs, glyph{asset: "shield", qty: v.Armor})
		}
		return gs, true
	case engine.GainAssault:
		return []glyph{{asset: "kw-assault"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.GainAssaultUntilNextTurn:
		return []glyph{{asset: "kw-assault"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.GainTextBox:
		return []glyph{targetGlyph(v.Source), arrowTo(targetGlyph(v.Target))}, true
	case engine.LendTextBoxFromHand:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			arrowTo(glyph{asset: "type-creature", decor: decorChosen}),
		}, true
	case engine.FuseTriggersForTurn:
		return []glyph{{asset: "glyph-reap"}, {asset: "glyph-swap"}, {asset: "glyph-fight"}}, true
	case engine.AddPowerCounter:
		if v.Per != nil || v.Equal != nil {
			return []glyph{{asset: "power"}, arrowTo(targetGlyph(v.Target))}, true
		}
		return []glyph{{asset: "power", qty: v.Amount}, arrowTo(targetGlyph(v.Target))}, true
	case engine.Swap:
		return []glyph{
			{asset: "type-creature", decor: decorThis},
			{asset: "glyph-swap"},
			arrowTo(targetGlyph(v.With)),
		}, true
	case engine.SwapChosen:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "glyph-swap"},
			{asset: "type-creature", decor: decorChosen},
		}, true
	case engine.RearrangeBattleline:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "glyph-swap"},
			{asset: "type-creature", decor: decorChosen},
		}, true
	case engine.Repeat:
		return effectGlyphs(v.Do)
	case engine.May:
		return effectGlyphs(v.Do)
	case engine.ByActivePlayer:
		return effectGlyphs(v.Do)
	case engine.Then:
		if first, ok := v.First.(engine.Effect); ok {
			return composeGlyphs(first, v.Result)
		}
		gs, _ := effectGlyphs(v.Result)
		return gs, false
	case engine.Sequence:
		return composeGlyphs(v.Effects...)
	case engine.Sentences:
		return composeGlyphs(v.Effects...)
	case engine.ChooseOne:
		return append([]glyph{{asset: "glyph-choose"}}, mustCompose(v.Options...)...), true
	case engine.Conditional:
		if v.Else != nil {
			return composeGlyphs(v.Then, v.Else)
		}
		return effectGlyphs(v.Then)
	case engine.ChooseCreatureThen:
		gs := []glyph{targetGlyph(v.Target)}
		more, covered := effectGlyphs(v.Then)
		return append(gs, more...), covered
	case engine.ChooseHouseThen:
		return append([]glyph{{asset: "glyph-choose"}}, mustCompose(v.Then)...), true
	case engine.OnChooseCreature:
		return append([]glyph{targetGlyph(v.Target)}, verbGlyphs(v.Verbs)...), true
	case engine.OneAtATime:
		return append([]glyph{targetGlyph(v.Target)}, verbGlyphs(v.Verbs)...), true
	case engine.ForRemainderOfTurn:
		return effectGlyphs(v.Do)
	case engine.ForOpponentNextTurn:
		return effectGlyphs(v.Do)
	case engine.NextPlayed:
		return append([]glyph{{asset: "glyph-play"}}, mustCompose(v.EntersPlay)...), true
	case engine.ForEach:
		return effectGlyphs(v.Do)
	case engine.TriggerAbility:
		return []glyph{
			{asset: triggerIcon(v.Trigger)},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.MoveAember:
		return []glyph{{asset: "aember"}, arrowTo(moveAemberDest(v.Onto, v.To))}, true
	case engine.MoveAemberFromPool:
		return []glyph{{asset: "aember", qty: v.Amount}, arrowTo(targetGlyph(v.Target))}, true
	case engine.PlaceAemberOnThis:
		return []glyph{
			{asset: "aember", qty: v.Amount},
			arrowTo(glyph{asset: "card-back", decor: decorThis}),
		}, true
	case engine.MoveAemberToSupply:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorChosen},
			arrowTo(glyph{asset: "glyph-return"}),
		}, true
	case engine.RedistributeCapturedAember:
		decor := decorFriendly
		if v.Side == engine.Opponent {
			decor = decorEnemy
		}
		return []glyph{{asset: "aember", decor: decor}, {asset: "glyph-swap"}}, true
	case engine.RedistributeDamage:
		return []glyph{{asset: "damage"}, {asset: "glyph-swap"}}, true
	case engine.MayPlayOrUse:
		return mayPlayOrUseGlyphs(v), true
	case engine.PlayOrUse:
		return []glyph{{asset: "glyph-play"}, {asset: "glyph-action"}}, true
	case engine.CannotBeDealtDamage:
		return []glyph{{asset: "shield"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.RedirectFightDamage:
		return []glyph{
			{asset: "glyph-fight"},
			{asset: "damage"},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.TakeControl:
		// The host-creature form (Collar of Subordination) has no Target; it takes
		// this creature, so render the "this creature" noun rather than a blank.
		subject := targetGlyph(v.Target)
		if v.Target == (engine.Target{}) {
			subject = glyph{asset: "type-creature", decor: decorThis}
		}
		return []glyph{
			subject,
			arrowTo(glyph{asset: "type-creature", decor: decorFriendly}),
		}, true
	case engine.PutChosen:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: a})}, true
		}
		return fallbackGlyphs(e), false
	case engine.PutFromDiscard:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: a})}, true
		}
		return fallbackGlyphs(e), false
	case engine.AttachSelfTo:
		return []glyph{
			{asset: "card-back", decor: decorThis},
			arrowTo(glyph{asset: "type-creature"}),
		}, true
	case engine.PutUnderFromHand:
		return []glyph{{asset: "zone-hand"}, arrowTo(glyph{asset: "card-back"})}, true
	case engine.PutUnderIntoPlay:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true
	case engine.TriggerGraftedPlayEffect:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true
	case engine.SwapDeckAndDiscard:
		return []glyph{
			{asset: "zone-deck"},
			{asset: "glyph-swap"},
			{asset: "zone-discard"},
		}, true
	case engine.RaiseKeyCost:
		return []glyph{{asset: "forge"}, {asset: "aember", qty: v.Amount}}, true
	case engine.RaiseKeyCostPerHouseCreature:
		return []glyph{{asset: "forge"}, {asset: "aember", qty: v.Amount}}, true
	case engine.LowerKeyCost:
		return []glyph{{asset: "forge"}, {asset: "aember", qty: -v.Amount}}, true
	case engine.SkipForgePhase:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true
	case engine.Restrict:
		verb := "glyph-reap"
		switch v.Action {
		case engine.RestrictFighting:
			verb = "glyph-fight"
		case engine.RestrictUse:
			verb = "glyph-action"
		}
		return []glyph{{asset: verb}, {asset: "glyph-ban"}}, true
	case engine.CancelFight:
		return []glyph{{asset: "glyph-fight"}, {asset: "glyph-ban"}}, true
	case engine.CancelForge:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true
	case engine.CannotPlay:
		return []glyph{{asset: "glyph-play"}, {asset: "glyph-ban"}}, true
	case engine.PlayersCannotPlay:
		if a := typeIconName(v.Type); a != "" {
			return []glyph{{asset: a}, {asset: "glyph-play"}, {asset: "glyph-ban"}}, true
		}
		return []glyph{{asset: "glyph-play"}, {asset: "glyph-ban"}}, true
	case engine.CreaturesCannot:
		action := "glyph-reap"
		if v.Action == engine.FightUse {
			action = "glyph-fight"
		}
		return []glyph{{asset: action}, {asset: "glyph-ban"}}, true
	case engine.LoseKeyword:
		if a := keywordIcon(v.Keyword); a != "" {
			return []glyph{{asset: a}, {asset: "glyph-ban"}}, true
		}
		return []glyph{{asset: "glyph-unknown"}, {asset: "glyph-ban"}}, true
	case engine.LoseKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return []glyph{{asset: "glyph-unknown"}, {asset: "glyph-ban"}}, true
			}
			gs = append(gs, glyph{asset: a})
		}
		return append(gs, glyph{asset: "glyph-ban"}), true
	case engine.GainKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return []glyph{{asset: "glyph-unknown"}}, true
			}
			gs = append(gs, glyph{asset: a})
		}
		return append(gs, arrowTo(targetGlyph(v.Target))), true
	case engine.NameHouse:
		// The chosen house is barred; ChooseHouseThen supplies the choose glyph.
		return []glyph{{asset: "glyph-ban"}}, true
	case engine.NameCard:
		// Name a card, then bar every copy of it from being played.
		return []glyph{{asset: "glyph-choose"}, {asset: "glyph-ban"}}, true
	case engine.OpponentNamesHouse:
		return []glyph{{asset: "glyph-choose"}}, true
	case engine.LookAtTopOfDeck:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}, true
	case engine.RevealHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true
	case engine.RevealRandomFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true
	case engine.RevealChosenFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true
	case engine.Search:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-search"}}, true
	case engine.Instead:
		return []glyph{{asset: "glyph-swap"}}, true
	case engine.DiscardUntil:
		return []glyph{{asset: "zone-deck"}, {asset: "zone-discard"}}, true
	case engine.ArchiveDiscardedThisWay:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-archives"})}, true
	case engine.PutDiscardedIntoHand:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-hand"})}, true
	case engine.DiscardHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{h, arrowTo(glyph{asset: "zone-discard"})}, true
	case engine.RefillHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{{asset: "zone-deck"}, arrowTo(h)}, true
	case engine.DiscardTop:
		return []glyph{{asset: "zone-discard", qty: v.Amount, decor: playerDecor(v.Player)}}, true
	case engine.ResolveBonusIcons:
		return []glyph{
			targetGlyph(v.Target),
			arrowTo(glyph{asset: "aember"}),
			{asset: "capture"},
			{asset: "damage"},
			{asset: "draw"},
		}, true
	case engine.ExtraBonusIconResolution:
		return []glyph{
			{asset: "glyph-play"},
			arrowTo(glyph{asset: "aember"}),
			{asset: "capture"},
			{asset: "damage"},
			{asset: "draw"},
		}, true
	case engine.ForEachDiscarded:
		return effectGlyphs(v.Do)
	case engine.UnforgeKey:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true
	case engine.ShuffleChosenCreaturesFromZones:
		return []glyph{{asset: "type-creature"}, arrowTo(glyph{asset: "zone-deck"})}, true
	case engine.LoseArmor:
		return []glyph{
			{asset: "shield"},
			{asset: "glyph-ban"},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.Graft:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "card-back"})}, true
	case engine.ConsiderFlank:
		return []glyph{{asset: "glyph-flank"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.MoveToFlank:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "glyph-flank"})}, true
	case engine.TurnIntoCreature:
		return []glyph{
			targetGlyph(v.Target),
			arrowTo(glyph{asset: "type-creature"}),
			{asset: "glyph-flank"},
		}, true
	case engine.MoveWithinBattleline:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "glyph-swap"})}, true
	case engine.BelongToHouse:
		if a := houseIconName(v.House); a != "" {
			return []glyph{{asset: a}, arrowTo(targetGlyph(v.Target))}, true
		}
		return []glyph{{asset: "glyph-swap"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.GainAbility:
		gs := []glyph{targetGlyph(v.Target), {asset: triggerIcon(v.Ability.Trigger)}}
		return append(gs, mustCompose(v.Ability.Effect)...), true
	case engine.TakesExtraDamage:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "damage", qty: v.Amount})}, true
	case engine.ReadyCreatures:
		return []glyph{
			{asset: "exhausted", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true
	case engine.PlayCardUnder:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true
	case engine.ArchiveCardUnder:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "zone-archives"})}, true
	case engine.PlayRevealedCard:
		return []glyph{{asset: "glyph-play"}}, true
	case engine.RevealTopOfDeck:
		g := []glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}
		for _, act := range v.Then {
			if m, ok := act.(engine.ChooseAndMove); ok && m.Dest == engine.IntoPurge {
				g = append(g, arrowTo(glyph{asset: "zone-purge"}))
			}
		}
		return g, true
	case engine.PutRevealedCard:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: deckDestZone(v.To)})}, true
	case engine.ChangeActiveHouse:
		return []glyph{{asset: "glyph-choose"}}, true
	case engine.EndTurn:
		return []glyph{{asset: "phase-turn"}}, true
	case engine.PurgeFromHand:
		hand := glyph{asset: "zone-hand"}
		if _, each := v.Selection.(engine.Each); each {
			hand.decor = decorEach
		}
		return []glyph{hand, arrowTo(glyph{asset: "zone-purge"})}, true
	case engine.MustChooseHouse:
		return []glyph{{asset: "glyph-choose", decor: playerDecor(v.Player)}}, true
	case engine.CannotChooseHouse:
		return []glyph{
			{asset: "glyph-choose", decor: playerDecor(v.Player)},
			{asset: "glyph-ban"},
		}, true
	case engine.WagerOpponentChoosesChosenHouse:
		return []glyph{{asset: "aember", qty: v.Amount, decor: decorEnemy | decorChosen}}, true
	case engine.ForDuration:
		return composeGlyphs(v.Effects...)
	case engine.GainUntilNextTurn:
		return composeGlyphs(v.Effects...)
	case engine.GainTrait:
		// A trait has no icon in the strip's vocabulary — traits render as the
		// card's text, not glyphs. It only ever folds beside a keyword grant that
		// carries the line's glyph, so it renders nothing yet counts as covered.
		return nil, true
	case engine.DestroyChosen:
		return []glyph{{asset: "glyph-destroy"}, arrowTo(targetGlyph(v.Target))}, true
	case engine.BatchDestroy:
		return []glyph{
			{asset: "glyph-destroy"},
			arrowTo(glyph{asset: "type-creature", decor: decorEach}),
		}, true
	case engine.DestroyEachCreatureAtEndOfTurn:
		return []glyph{
			{asset: "glyph-destroy"},
			arrowTo(glyph{asset: "type-creature", decor: decorEach}),
		}, true
	case engine.PurgeSource:
		return []glyph{{asset: "zone-purge", decor: decorThis}}, true
	case engine.DiscardArchives:
		return []glyph{{asset: "zone-archives"}, arrowTo(glyph{asset: "zone-discard"})}, true
	case engine.PutFromHand:
		return []glyph{{asset: "zone-hand"}, arrowTo(glyph{asset: "glyph-play"})}, true
	case engine.CopyPrintedStats:
		return []glyph{
			targetGlyph(v.Source),
			arrowTo(glyph{asset: "type-creature", decor: decorThis}),
		}, true
	case engine.ScheduleOnLeave:
		return append(
			[]glyph{{asset: "type-creature", decor: decorThis}},
			mustCompose(v.Do)...), true
	case engine.EachPlayerPutsHandCreaturesIntoPlay:
		return []glyph{
			{asset: "zone-hand", decor: decorEach},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true
	case engine.PutNextTacticIntoHand:
		return []glyph{{asset: "type-action"}, arrowTo(glyph{asset: "zone-hand"})}, true
	case engine.DamageOthersAfterUsingTrait:
		return []glyph{
			{asset: "glyph-action"},
			{asset: "damage", qty: v.Amount},
			arrowTo(glyph{asset: "type-creature", decor: decorEach}),
		}, true
	case engine.PutDiscardedIntoPlay:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "glyph-play"})}, true
	default:
		return fallbackGlyphs(e), false
	}
}

// mustCompose renders sub-effects like composeGlyphs but discards the covered
// flag: the caller has already prefixed a real glyph, so the line reads even when
// a nested effect falls back to the abstract glyph.
func mustCompose(effects ...engine.Effect) []glyph {
	gs, _ := composeGlyphs(effects...)
	return gs
}

// mayPlayOrUseGlyphs renders an out-of-house permission grant, narrowing to the
// verbs and houses its axes select: a fight grant to a fight glyph, an
// artifacts-any-house grant to an artifact-and-action pair, a named-house grant to
// play/action, and an exclusion or controlled grant to play (plus action when it
// also frees use).
func mayPlayOrUseGlyphs(e engine.MayPlayOrUse) []glyph {
	if e.Houses.Controlled {
		if e.Grant&engine.GrantUse != 0 {
			return []glyph{{asset: "glyph-play"}, {asset: "glyph-action"}}
		}
		return []glyph{{asset: "glyph-play"}}
	}
	switch e.Houses.Match.Kind {
	case engine.MatchExceptHouse:
		if e.Grant&engine.GrantUse != 0 {
			return []glyph{{asset: "glyph-play"}, {asset: "glyph-action"}}
		}
		return []glyph{{asset: "glyph-play"}}
	case engine.MatchAnyHouse:
		if e.Grant&engine.GrantFight != 0 {
			return []glyph{{asset: "glyph-fight", decor: decorFriendly | decorEach}}
		}
		return []glyph{
			{asset: "type-artifact", decor: decorFriendly},
			{asset: "glyph-action"},
		}
	default: // MatchNamedHouse, MatchChosenHouse
		if e.Grant == engine.GrantFight {
			d := decorEach
			if e.Houses.Match.House != engine.HouseNone {
				d |= decorFriendly
			}
			return []glyph{{asset: "glyph-fight", decor: d}}
		}
		if e.Grant&engine.GrantPlay != 0 {
			return []glyph{{asset: "glyph-play"}, {asset: "glyph-action"}}
		}
		return []glyph{{asset: "glyph-action", decor: decorFriendly}}
	}
}

// verbGlyphs renders the verbs a chosen-creature effect applies in order — ready,
// fight, use, stun, exhaust, gain-keyword — reusing the same glyphs those actions
// draw on their own.
func verbGlyphs(verbs []engine.CreatureVerb) []glyph {
	gs := make([]glyph, 0, len(verbs))
	for _, verb := range verbs {
		switch vv := verb.(type) {
		case engine.ReadyVerb:
			gs = append(gs, glyph{asset: "exhausted", decor: decorFriendly})
		case engine.ReapVerb:
			gs = append(gs, glyph{asset: "glyph-reap"})
		case engine.FightVerb:
			gs = append(gs, glyph{asset: "glyph-fight"})
		case engine.UseVerb:
			gs = append(gs, glyph{asset: "glyph-action"})
		case engine.StunVerb:
			gs = append(gs, glyph{asset: "stun"})
		case engine.ExhaustVerb:
			gs = append(gs, glyph{asset: "exhausted"})
		case engine.GainKeywordVerb:
			if a := keywordIcon(vv.Keyword); a != "" {
				gs = append(gs, glyph{asset: a})
			}
		}
	}
	return gs
}

// moveAemberDest is where moved Æmber lands: onto a chosen card, or into a pool.
func moveAemberDest(onto engine.Target, to engine.Player) glyph {
	if onto != (engine.Target{}) {
		return targetGlyph(onto)
	}
	return glyph{asset: "aember", decor: playerDecor(to)}
}

// damageThenGlyphs renders a "deal N damage to <target>, then <follow-up>" effect
// as the damage glyph arrowed to its target, followed by the follow-up's glyphs.
func damageThenGlyphs(amount int, target engine.Target, then engine.Effect) ([]glyph, bool) {
	gs := []glyph{{asset: "damage", qty: amount}, arrowTo(targetGlyph(target))}
	more, covered := effectGlyphs(then)
	return append(gs, more...), covered
}

// composeGlyphs renders a run of sub-effects one after another, reporting covered
// only when every sub-effect maps to a real glyph.
func composeGlyphs(effects ...engine.Effect) ([]glyph, bool) {
	out := make([]glyph, 0, len(effects)*2)
	covered := true
	for _, e := range effects {
		gs, c := effectGlyphs(e)
		out = append(out, gs...)
		covered = covered && c
	}
	return out, covered
}

// destinationGlyph is the zone glyph a card is moved to, or "" for a destination
// with no drawn zone yet. The three deck spots share the one deck glyph.
func destinationGlyph(d engine.Destination) string {
	switch d {
	case engine.ToHand:
		return "zone-hand"
	case engine.ToTopOfDeck, engine.ToBottomOfDeck, engine.ToDeckShuffled:
		return "zone-deck"
	case engine.ToArchives, engine.ToArchives.Yours():
		return "zone-archives"
	}
	return ""
}

// deckDestZone is the zone glyph a revealed deck card is moved to.
func deckDestZone(d engine.DeckDest) string {
	switch d {
	case engine.IntoHand:
		return "zone-hand"
	case engine.IntoArchives:
		return "zone-archives"
	case engine.IntoDiscard:
		return "zone-discard"
	case engine.IntoPurge:
		return "zone-purge"
	}
	return ""
}

// arrowTo marks a glyph as the result of a result-gate arrow, so a target reads
// as "→ <target>".
func arrowTo(g glyph) glyph {
	g.arrow = true
	return g
}

// counterAsset maps a generic counter kind to its icon-strip asset, or "" for a
// kind with no icon yet — TestCounterIconNamesHaveAssets fails on the "" so every
// counter ships its own unique SVG.
func counterAsset(kind engine.CounterKind) string {
	switch kind {
	case engine.CounterDoom:
		return "generic-counter-doom"
	case engine.CounterFuse:
		return "generic-counter-fuse"
	case engine.CounterGrowth:
		return "generic-counter-growth"
	case engine.CounterGlory:
		return "generic-counter-glory"
	case engine.CounterDisruption:
		return "generic-counter-disruption"
	case engine.CounterScheme:
		return "generic-counter-scheme"
	case engine.CounterWarrant:
		return "generic-counter-warrant"
	default:
		return ""
	}
}

// spreadTargetGlyph is the creature glyph a DealDamage Spread hits. Every spread
// chooses one or more creatures; the amounts and neighbor split stay in the rules
// text, so the strip summarises the spread as its chosen-creature noun.
func spreadTargetGlyph() glyph {
	return glyph{asset: "type-creature", decor: decorChosen}
}

// targetGlyph renders a Target as its noun glyph plus the decorations that carry
// its enemy/friendly/each/chosen shape. Fine filters (power, house, trait) stay
// in the rules text; the strip summarises the noun.
func targetGlyph(t engine.Target) glyph {
	switch t.Kind {
	case engine.TargetThisCreature:
		return glyph{asset: "type-creature", decor: decorThis}
	case engine.TargetTriggeringCreature, engine.TargetTheOtherCreature,
		engine.TargetTheChosenCreature, engine.TargetCreatureFought,
		engine.TargetTheFoughtCreature:
		return glyph{asset: "type-creature", decor: decorChosen}
	case engine.TargetEachCreature:
		return glyph{asset: "type-creature", decor: decorEach}
	case engine.TargetEachFriendlyCreature, engine.TargetEachOtherFriendlyCreature:
		return glyph{asset: "type-creature", decor: decorEach | decorFriendly}
	case engine.TargetEachEnemyCreature:
		return glyph{asset: "type-creature", decor: decorEach | decorEnemy}
	case engine.TargetChosenCreature, engine.TargetChosenOtherCreature:
		return glyph{asset: "type-creature", decor: decorChosen}
	case engine.TargetChosenFriendlyCreature, engine.TargetChosenOtherFriendlyCreature:
		return glyph{asset: "type-creature", decor: decorChosen | decorFriendly}
	case engine.TargetChosenEnemyCreature:
		return glyph{asset: "type-creature", decor: decorChosen | decorEnemy}
	case engine.TargetEachArtifact:
		return glyph{asset: "type-artifact", decor: decorEach}
	case engine.TargetEachFriendlyArtifact:
		return glyph{asset: "type-artifact", decor: decorEach | decorFriendly}
	case engine.TargetEachEnemyArtifact:
		return glyph{asset: "type-artifact", decor: decorEach | decorEnemy}
	case engine.TargetChosenArtifact:
		return glyph{asset: "type-artifact", decor: decorChosen}
	case engine.TargetChosenFriendlyArtifact:
		return glyph{asset: "type-artifact", decor: decorChosen | decorFriendly}
	case engine.TargetChosenEnemyArtifact:
		return glyph{asset: "type-artifact", decor: decorChosen | decorEnemy}
	case engine.TargetChosenUpgrade:
		return glyph{asset: "type-upgrade", decor: decorChosen}
	case engine.TargetEachCardInPlay, engine.TargetEachFriendlyCardInPlay,
		engine.TargetChosenCreatureOrArtifact,
		engine.TargetChosenFriendlyCreatureOrArtifact,
		engine.TargetChosenEnemyCreatureOrArtifact:
		return glyph{asset: "card-back", decor: cardInPlayDecor(t.Kind)}
	default:
		return glyph{text: t.Text()}
	}
}

// cardInPlayDecor picks the enemy/friendly/each shading for the "card in play"
// target kinds that share one noun glyph.
func cardInPlayDecor(k engine.TargetKind) decor {
	switch k {
	case engine.TargetEachCardInPlay:
		return decorEach
	case engine.TargetEachFriendlyCardInPlay:
		return decorEach | decorFriendly
	case engine.TargetChosenFriendlyCreatureOrArtifact:
		return decorChosen | decorFriendly
	case engine.TargetChosenEnemyCreatureOrArtifact:
		return decorChosen | decorEnemy
	default:
		return decorChosen
	}
}

// playerDecor tints an economy glyph by who it acts on: an opponent-facing effect
// reads as enemy, an each-player effect stacks.
func playerDecor(p engine.Player) decor {
	switch p {
	case engine.Opponent:
		return decorEnemy
	case engine.EachPlayer:
		return decorEach
	default:
		return 0
	}
}

// fallbackGlyphs renders an effect the pass does not yet map as the single
// abstract "unknown" glyph. The strip is pure icons — no words — so an unmapped
// mechanic shows a placeholder sigil rather than its printed text (ADR 0022).
func fallbackGlyphs(engine.Effect) []glyph {
	return []glyph{{asset: "glyph-unknown"}}
}

// effectTypeName is the effect's concrete Go type name (e.g. "engine.DealDamage"),
// used by the totality test to name an uncovered effect.
func effectTypeName(e engine.Effect) string {
	return fmt.Sprintf("%T", e)
}

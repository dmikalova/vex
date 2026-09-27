package web

import (
	"reflect"

	"github.com/dmikalova/vex/internal/engine"
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
		lines = append(lines, glyphLine{
			glyphs:  kw,
			covered: true,
		})
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
		lines = append(lines, glyphLine{
			glyphs:  replaceGlyphs(def.Replaces),
			covered: true,
		})
	}
	if len(def.KeyCostChanges) > 0 {
		lines = append(lines, glyphLine{
			glyphs:  keyCostChangeGlyphs(),
			covered: true,
		})
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

// glyphFamilies is the dispatch chain the Iconography pass runs an effect down:
// each family claims the effect types of one mechanic domain, coarsely mirroring
// the engine's effect_<mechanic>.go split, and the first family to claim an effect
// wins (ADR 0047). A family returns (glyphs, covered, ok).
//
// The third value is load-bearing and must not be collapsed into covered, because
// covered = false already means something else: it says "this drew the abstract
// unknown glyph", and four cases return real glyphs with covered = false — a Then
// whose First is not an Effect, and PutFromPlay / PutChosen / PutCard with a
// destination that has no glyph. A family therefore cannot signal "not mine" by
// returning covered = false; ok is how it declines.
//
// The compiler cannot see across the chain, so two families claiming the same
// effect type would silently make the second dead code — drawing a wrong glyph
// rather than an unknown one. TestGlyphFamiliesAreDisjoint polices that.
//
// The chain is filled in init rather than at its declaration because every family
// recurses through effectGlyphs, which reads the slice, and Go rejects that as an
// initialization cycle.
var glyphFamilies []func(engine.Effect) (gs []glyph, covered, ok bool)

func init() {
	glyphFamilies = []func(engine.Effect) (gs []glyph, covered, ok bool){
		residualEffectGlyphs,
	}
}

// effectGlyphs renders one effect to its glyphs, reporting whether the effect was
// covered by a real mapping (false means it fell back to the abstract unknown
// glyph). It runs the effect down glyphFamilies and takes the first claim.
func effectGlyphs(e engine.Effect) ([]glyph, bool) {
	for _, family := range glyphFamilies {
		if gs, covered, ok := family(e); ok {
			return gs, covered
		}
	}
	return fallbackGlyphs(e), false
}

// claim adapts a two-valued glyph rendering — a recursive effectGlyphs, a
// composeGlyphs over sub-effects — into a family's three-valued return, saying
// the family claimed the effect.
func claim(gs []glyph, covered bool) ([]glyph, bool, bool) {
	return gs, covered, true
}

// residualEffectGlyphs is the trailing family of the chain: every effect type not
// yet lifted into a mechanic-domain family. It declines (ok = false) only from its
// default, where effectGlyphs falls back to the abstract unknown glyph.
func residualEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.DealDamage:
		// A follow-up on the damaged creature draws the two clauses joined, not the
		// bare damage glyph.
		if v.Then != nil {
			return claim(damageThenGlyphs(v.Amount, v.Target, v.Then))
		}
		// A Spread carries its own creature targets rather than filling Target, so it
		// draws its own summary noun; the counts and neighbor split stay in the text.
		if v.Spread != nil {
			return []glyph{{asset: "damage"}, arrowTo(spreadTargetGlyph())}, true, true
		}
		// The per-count variants hit several creatures or scale by a board count; the
		// numeral lives in the text, so the glyph drops the qty.
		if v.Per != nil || v.PerTarget != nil || v.AmountFrom != nil {
			return []glyph{{asset: "damage"}, arrowTo(targetGlyph(v.Target))}, true, true
		}
		return []glyph{
			{asset: "damage", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.ForEachHouse:
		return claim(effectGlyphs(v.Do))
	case engine.GainAember:
		if v.Per != nil || v.EqualTo != nil {
			return []glyph{{asset: "aember", decor: playerDecor(v.Player)}}, true, true
		}
		return []glyph{{asset: "aember", qty: v.Amount, decor: playerDecor(v.Player)}}, true, true
	case engine.LoseAember:
		return []glyph{{asset: "aember", text: "−", decor: playerDecor(v.Player)}}, true, true
	case engine.StealAember:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorEnemy | decorChosen},
		}, true, true
	case engine.CaptureAember:
		return []glyph{{asset: "aember", qty: v.Amount, decor: decorEnemy}}, true, true
	case engine.CaptureFromAnyPlayer:
		return []glyph{{asset: "aember", qty: v.Amount}}, true, true
	case engine.DistributeCapture:
		return []glyph{{asset: "aember", decor: decorEnemy}}, true, true
	case engine.GiveAember:
		src := glyph{
			asset: "aember",
			decor: decorEnemy,
		}
		if !v.All {
			src.qty = v.Amount
		}
		return []glyph{src, arrowTo(glyph{asset: "aember"})}, true, true
	case engine.GainChains:
		return []glyph{{asset: "chains", qty: v.Amount, decor: playerDecor(v.Player)}}, true, true
	case engine.Draw:
		return []glyph{{asset: "zone-hand", qty: v.Amount}}, true, true
	case engine.Stun:
		return []glyph{{asset: "stun"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Enrage:
		return []glyph{{asset: "glyph-fight"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Ward:
		return []glyph{{asset: "shield"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.RemoveWard:
		return []glyph{
			{asset: "shield"},
			{asset: "glyph-ban"},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Exhaust:
		return []glyph{{asset: "exhausted"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Ready:
		return []glyph{
			{asset: "exhausted", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Unstun:
		return []glyph{
			{asset: "stun", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Destroy:
		return []glyph{{asset: "glyph-destroy"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.PurgeCreature:
		return []glyph{{asset: "zone-purge"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Heal:
		return []glyph{
			{asset: "glyph-heal", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Exalt:
		return []glyph{
			{asset: "aember", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.ForgeKey:
		return []glyph{{asset: "forge"}}, true, true
	case engine.PlaceCounter:
		return []glyph{{asset: counterAsset(v.Kind)}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.RemoveCounters:
		return []glyph{{asset: counterAsset(v.Kind)}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.BlankEnemyText:
		return []glyph{{asset: "type-creature", decor: decorEnemy | decorEach}}, true, true
	case engine.ArchiveCard:
		if v.Zone == engine.Purged {
			return []glyph{
				{asset: "zone-purge"},
				arrowTo(glyph{asset: "zone-archives"}),
			}, true, true
		}
		return []glyph{{asset: "zone-archives", qty: engine.FixedCardCount(v.Quantity)}}, true, true
	case engine.ArchiveFromPlay:
		return []glyph{{asset: "zone-archives"}}, true, true
	case engine.ArchiveSource:
		return []glyph{{asset: "zone-archives", decor: decorThis}}, true, true
	case engine.ArchiveGrantingUpgrade:
		return []glyph{
			{asset: "card-back", decor: decorThis},
			arrowTo(glyph{asset: "zone-archives"}),
		}, true, true
	case engine.DiscardCard:
		return []glyph{{asset: "zone-discard"}}, true, true
	case engine.PurgeCard:
		src := glyph{asset: "zone-hand"}
		switch v.Zones[0] {
		case engine.Discard:
			src.asset = "zone-discard"
		case engine.Archives:
			src.asset = "zone-archives"
		}
		if _, each := v.Selection.(engine.Each); each {
			src.decor = decorEach
		}
		return []glyph{src, arrowTo(glyph{asset: "zone-purge"})}, true, true
	case engine.PurgeArchivedCardThen:
		gs := []glyph{{asset: "zone-purge"}}
		more, _ := effectGlyphs(v.Then)
		return append(gs, more...), true, true
	case engine.Shuffle:
		return []glyph{{asset: "zone-deck"}}, true, true
	case engine.ShuffleIntoDeck:
		// A multi-zone shuffle has no single source glyph, so it shows what it takes.
		src := glyph{asset: "zone-discard"}
		if len(v.From) > 1 {
			src = glyph{asset: "type-creature"}
		}
		return []glyph{src, arrowTo(glyph{asset: "zone-deck"})}, true, true
	case engine.PutItIntoHand:
		return []glyph{{asset: "glyph-return"}}, true, true
	case engine.PlayFrom, engine.PlayTopOfDeck, engine.PutIntoPlay:
		return []glyph{{asset: "glyph-play"}}, true, true
	case engine.PlayFromOpponent:
		zone := "zone-deck"
		if v.From == engine.Archives {
			zone = "zone-archives"
		}
		return []glyph{
			{asset: zone, decor: decorEnemy},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true, true
	case engine.PlayItFromOpponentDiscard:
		return []glyph{
			{asset: "zone-discard", decor: decorEnemy},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true, true
	case engine.DiscardFromOpponent:
		zone := "zone-deck"
		if len(v.Sources) > 0 && v.Sources[0] == engine.Archives {
			zone = "zone-archives"
		}
		return []glyph{
			{asset: zone, decor: decorEnemy},
			arrowTo(glyph{
				asset: "zone-discard",
				decor: decorEnemy,
			}),
		}, true, true
	case engine.PutFromPlay:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: a})}, true, true
		}
		return fallbackGlyphs(e), false, true
	case engine.ExhaustCreatures:
		return []glyph{{asset: "exhausted"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Use:
		return []glyph{{asset: "glyph-action"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.RepeatedFight:
		return []glyph{{asset: "glyph-fight"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.GainStats:
		gs := make([]glyph, 0, 3)
		if v.Power != 0 {
			gs = append(gs, glyph{
				asset: "power",
				qty:   v.Power,
			})
		}
		if v.Armor != 0 {
			gs = append(gs, glyph{
				asset: "shield",
				qty:   v.Armor,
			})
		}
		return append(gs, arrowTo(targetGlyph(v.Target))), true, true
	case engine.OverrideStats:
		gs := make([]glyph, 0, 2)
		if v.HasPower {
			gs = append(gs, glyph{
				asset: "power",
				qty:   v.Power,
			})
		}
		if v.HasArmor {
			gs = append(gs, glyph{
				asset: "shield",
				qty:   v.Armor,
			})
		}
		return gs, true, true
	case engine.GainAssault:
		return []glyph{{asset: "kw-assault"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.GainAssaultUntilNextTurn:
		return []glyph{{asset: "kw-assault"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.GainTextBox:
		return []glyph{targetGlyph(v.Source), arrowTo(targetGlyph(v.Target))}, true, true
	case engine.LendTextBoxFromHand:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorChosen,
			}),
		}, true, true
	case engine.FuseTriggersForTurn:
		return []glyph{
			{asset: "glyph-reap"},
			{asset: "glyph-swap"},
			{asset: "glyph-fight"},
		}, true, true
	case engine.AddPowerCounter:
		if v.Per != nil || v.Equal != nil {
			return []glyph{{asset: "power"}, arrowTo(targetGlyph(v.Target))}, true, true
		}
		return []glyph{{asset: "power", qty: v.Amount}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Swap:
		return []glyph{
			{asset: "type-creature", decor: decorThis},
			{asset: "glyph-swap"},
			arrowTo(targetGlyph(v.With)),
		}, true, true
	case engine.SwapChosen:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "glyph-swap"},
			{asset: "type-creature", decor: decorChosen},
		}, true, true
	case engine.RearrangeBattleline:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "glyph-swap"},
			{asset: "type-creature", decor: decorChosen},
		}, true, true
	case engine.Repeat:
		return claim(effectGlyphs(v.Do))
	case engine.May:
		return claim(effectGlyphs(v.Do))
	case engine.ByActivePlayer:
		return claim(effectGlyphs(v.Do))
	case engine.Then:
		if first, ok := v.First.(engine.Effect); ok {
			return claim(composeGlyphs(first, v.Result))
		}
		gs, _ := effectGlyphs(v.Result)
		return gs, false, true
	case engine.Sequence:
		return claim(composeGlyphs(v.Effects...))
	case engine.ChooseOne:
		return append([]glyph{{asset: "glyph-choose"}}, mustCompose(v.Options...)...), true, true
	case engine.Conditional:
		if v.Else != nil {
			return claim(composeGlyphs(v.Then, v.Else))
		}
		return claim(effectGlyphs(v.Then))
	case engine.ChooseCreatureThen:
		gs := []glyph{targetGlyph(v.Target)}
		more, covered := effectGlyphs(v.Then)
		return append(gs, more...), covered, true
	case engine.ChooseHouseThen:
		return append([]glyph{{asset: "glyph-choose"}}, mustCompose(v.Then)...), true, true
	case engine.OnChooseCreature:
		return append([]glyph{targetGlyph(v.Target)}, verbGlyphs(v.Verbs)...), true, true
	case engine.OneAtATime:
		return append([]glyph{targetGlyph(v.Target)}, verbGlyphs(v.Verbs)...), true, true
	case engine.ForRemainderOfTurn:
		return claim(effectGlyphs(v.Do))
	case engine.ForOpponentNextTurn:
		return claim(effectGlyphs(v.Do))
	case engine.NextPlayed:
		return append([]glyph{{asset: "glyph-play"}}, mustCompose(v.EntersPlay)...), true, true
	case engine.ForEach:
		return claim(effectGlyphs(v.Do))
	case engine.TriggerAbility:
		return []glyph{
			{asset: triggerIcon(v.Trigger)},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.MoveAember:
		return []glyph{{asset: "aember"}, arrowTo(moveAemberDest(v.Onto, v.To))}, true, true
	case engine.MoveAemberFromPool:
		return []glyph{{asset: "aember", qty: v.Amount}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.PlaceAemberOnThis:
		return []glyph{
			{asset: "aember", qty: v.Amount},
			arrowTo(glyph{
				asset: "card-back",
				decor: decorThis,
			}),
		}, true, true
	case engine.MoveAemberToSupply:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorChosen},
			arrowTo(glyph{asset: "glyph-return"}),
		}, true, true
	case engine.RedistributeCapturedAember:
		decor := decorFriendly
		if v.Side == engine.Opponent {
			decor = decorEnemy
		}
		return []glyph{{asset: "aember", decor: decor}, {asset: "glyph-swap"}}, true, true
	case engine.RedistributeDamage:
		return []glyph{{asset: "damage"}, {asset: "glyph-swap"}}, true, true
	case engine.MayPlayOrUse:
		return mayPlayOrUseGlyphs(v), true, true
	case engine.PlayOrUse:
		var gs []glyph
		if v.Grant == 0 || v.Grant&engine.GrantPlay != 0 {
			gs = append(gs, glyph{asset: "glyph-play"})
		}
		if v.Grant == 0 || v.Grant&engine.GrantUse != 0 {
			gs = append(gs, glyph{asset: "glyph-action"})
		}
		return gs, true, true
	case engine.CannotBeDealtDamage:
		return []glyph{{asset: "shield"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.RedirectFightDamage:
		return []glyph{
			{asset: "glyph-fight"},
			{asset: "damage"},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.TakeControl:
		// The host-creature form (Collar of Subordination) has no Target; it takes
		// this creature, so render the "this creature" noun rather than a blank.
		subject := targetGlyph(v.Target)
		if v.Target == (engine.Target{}) {
			subject = glyph{
				asset: "type-creature",
				decor: decorThis,
			}
		}
		return []glyph{
			subject,
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorFriendly,
			}),
		}, true, true
	case engine.PutChosen:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: a})}, true, true
		}
		return fallbackGlyphs(e), false, true
	case engine.PutCard:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: a})}, true, true
		}
		return fallbackGlyphs(e), false, true
	case engine.AttachSelfTo:
		return []glyph{
			{asset: "card-back", decor: decorThis},
			arrowTo(glyph{asset: "type-creature"}),
		}, true, true
	case engine.PutUnderFromHand:
		return []glyph{{asset: "zone-hand"}, arrowTo(glyph{asset: "card-back"})}, true, true
	case engine.PutUnderIntoPlay:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.TriggerGraftedPlayEffect:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.SwapDeckAndDiscard:
		return []glyph{
			{asset: "zone-deck"},
			{asset: "glyph-swap"},
			{asset: "zone-discard"},
		}, true, true
	case engine.RaiseKeyCost:
		return []glyph{{asset: "forge"}, {asset: "aember", qty: v.Amount}}, true, true
	case engine.LowerKeyCost:
		return []glyph{{asset: "forge"}, {asset: "aember", qty: -v.Amount}}, true, true
	case engine.SkipForgePhase:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true, true
	case engine.Restrict:
		verb := "glyph-reap"
		switch v.Action {
		case engine.RestrictFighting:
			verb = "glyph-fight"
		case engine.RestrictUse:
			verb = "glyph-action"
		}
		return []glyph{{asset: verb}, {asset: "glyph-ban"}}, true, true
	case engine.CancelFight:
		return []glyph{{asset: "glyph-fight"}, {asset: "glyph-ban"}}, true, true
	case engine.CancelForge:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true, true
	case engine.CannotPlay:
		return []glyph{{asset: "glyph-play"}, {asset: "glyph-ban"}}, true, true
	case engine.PlayersCannotPlay:
		if a := typeIconName(v.Type); a != "" {
			return []glyph{{asset: a}, {asset: "glyph-play"}, {asset: "glyph-ban"}}, true, true
		}
		return []glyph{{asset: "glyph-play"}, {asset: "glyph-ban"}}, true, true
	case engine.CreaturesCannot:
		action := "glyph-reap"
		if v.Action == engine.FightUse {
			action = "glyph-fight"
		}
		return []glyph{{asset: action}, {asset: "glyph-ban"}}, true, true
	case engine.LoseKeyword:
		if a := keywordIcon(v.Keyword); a != "" {
			return []glyph{{asset: a}, {asset: "glyph-ban"}}, true, true
		}
		return []glyph{{asset: "glyph-unknown"}, {asset: "glyph-ban"}}, true, true
	case engine.LoseKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return []glyph{{asset: "glyph-unknown"}, {asset: "glyph-ban"}}, true, true
			}
			gs = append(gs, glyph{asset: a})
		}
		return append(gs, glyph{asset: "glyph-ban"}), true, true
	case engine.GainKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return []glyph{{asset: "glyph-unknown"}}, true, true
			}
			gs = append(gs, glyph{asset: a})
		}
		return append(gs, arrowTo(targetGlyph(v.Target))), true, true
	case engine.NameHouse:
		// The chosen house is barred; ChooseHouseThen supplies the choose glyph.
		return []glyph{{asset: "glyph-ban"}}, true, true
	case engine.NameCard:
		// Name a card, then bar every copy of it from being played.
		return []glyph{{asset: "glyph-choose"}, {asset: "glyph-ban"}}, true, true
	case engine.OpponentNamesHouse:
		return []glyph{{asset: "glyph-choose"}}, true, true
	case engine.LookAtTopOfDeck:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}, true, true
	case engine.RevealHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true, true
	case engine.RevealRandomFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true, true
	case engine.RevealChosenFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true, true
	case engine.Search:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-search"}}, true, true
	case engine.Instead:
		return []glyph{{asset: "glyph-swap"}}, true, true
	case engine.DiscardUntil:
		return []glyph{{asset: "zone-deck"}, {asset: "zone-discard"}}, true, true
	case engine.ArchiveDiscardedThisWay:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-archives"})}, true, true
	case engine.PutDiscardedIntoHand:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-hand"})}, true, true
	case engine.DiscardHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{h, arrowTo(glyph{asset: "zone-discard"})}, true, true
	case engine.RefillHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{{asset: "zone-deck"}, arrowTo(h)}, true, true
	case engine.DiscardTop:
		return []glyph{
			{asset: "zone-discard", qty: v.Amount, decor: playerDecor(v.Player)},
		}, true, true
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
	case engine.ForEachDiscarded:
		return claim(effectGlyphs(v.Do))
	case engine.UnforgeKey:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true, true
	case engine.LoseArmor:
		return []glyph{
			{asset: "shield"},
			{asset: "glyph-ban"},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Graft:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "card-back"})}, true, true
	case engine.ConsiderFlank:
		return []glyph{{asset: "glyph-flank"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.MoveToFlank:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "glyph-flank"})}, true, true
	case engine.TurnIntoCreature:
		return []glyph{
			targetGlyph(v.Target),
			arrowTo(glyph{asset: "type-creature"}),
			{asset: "glyph-flank"},
		}, true, true
	case engine.MoveWithinBattleline:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "glyph-swap"})}, true, true
	case engine.BelongToHouse:
		if a := houseIconName(v.House); a != "" {
			return []glyph{{asset: a}, arrowTo(targetGlyph(v.Target))}, true, true
		}
		return []glyph{{asset: "glyph-swap"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.GainAbility:
		gs := []glyph{targetGlyph(v.Target), {asset: triggerIcon(v.Ability.Trigger)}}
		return append(gs, mustCompose(v.Ability.Effect)...), true, true
	case engine.TakesExtraDamage:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{
			asset: "damage",
			qty:   v.Amount,
		})}, true, true
	case engine.ReadyCreatures:
		return []glyph{
			{asset: "exhausted", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.PlayCardUnder:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.ArchiveCardUnder:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "zone-archives"})}, true, true
	case engine.PlayRevealedCard:
		return []glyph{{asset: "glyph-play"}}, true, true
	case engine.RevealTopOfDeck:
		g := []glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}
		for _, act := range v.Then {
			if m, ok := act.(engine.ChooseAndMove); ok && m.Dest == engine.IntoPurge {
				g = append(g, arrowTo(glyph{asset: "zone-purge"}))
			}
		}
		return g, true, true
	case engine.PutRevealedCard:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: deckDestZone(v.To)})}, true, true
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
	case engine.WagerOpponentChoosesChosenHouse:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorEnemy | decorChosen},
		}, true, true
	case engine.ForDuration:
		return claim(composeGlyphs(v.Effects...))
	case engine.GainUntilNextTurn:
		return claim(composeGlyphs(v.Effects...))
	case engine.GainTrait:
		// A trait has no icon in the strip's vocabulary — traits render as the
		// card's text, not glyphs. It only ever folds beside a keyword grant that
		// carries the line's glyph, so it renders nothing yet counts as covered.
		return nil, true, true
	case engine.DestroyChosen:
		return []glyph{{asset: "glyph-destroy"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.BatchDestroy:
		return []glyph{
			{asset: "glyph-destroy"},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorEach,
			}),
		}, true, true
	case engine.DestroyEachCreatureAtEndOfTurn:
		return []glyph{
			{asset: "glyph-destroy"},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorEach,
			}),
		}, true, true
	case engine.PurgeSource:
		return []glyph{{asset: "zone-purge", decor: decorThis}}, true, true
	case engine.DiscardArchives:
		return []glyph{{asset: "zone-archives"}, arrowTo(glyph{asset: "zone-discard"})}, true, true
	case engine.PutFromHand:
		return []glyph{{asset: "zone-hand"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.CopyPrintedStats:
		return []glyph{
			targetGlyph(v.Source),
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorThis,
			}),
		}, true, true
	case engine.ScheduleOnLeave:
		return append(
			[]glyph{{asset: "type-creature", decor: decorThis}},
			mustCompose(v.Do)...), true, true
	case engine.EachPlayerPutsHandCreaturesIntoPlay:
		return []glyph{
			{asset: "zone-hand", decor: decorEach},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true, true
	case engine.PutNextTacticIntoHand:
		return []glyph{{asset: "type-tactic"}, arrowTo(glyph{asset: "zone-hand"})}, true, true
	case engine.DamageOthersAfterUsingTrait:
		return []glyph{
			{asset: "glyph-action"},
			{asset: "damage", qty: v.Amount},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorEach,
			}),
		}, true, true
	case engine.PutDiscardedIntoPlay:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	default:
		return nil, false, false
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
			gs = append(gs, glyph{
				asset: "exhausted",
				decor: decorFriendly,
			})
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
	return glyph{
		asset: "aember",
		decor: playerDecor(to),
	}
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

// fallbackGlyphs renders an effect the pass does not yet map as the single
// abstract "unknown" glyph. The strip is pure icons — no words — so an unmapped
// mechanic shows a placeholder sigil rather than its printed text (ADR 0022).
func fallbackGlyphs(engine.Effect) []glyph {
	return []glyph{{asset: "glyph-unknown"}}
}

// effectTypeName is the effect's concrete Go type name (e.g. "engine.DealDamage"),
// used by the totality test to name an uncovered effect.
func effectTypeName(e engine.Effect) string {
	return reflect.TypeOf(e).String()
}

package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the shared noun and
// decor vocabulary — a Target's noun glyph, the zone glyphs a card moves
// between, the result-gate arrow, and the enemy/friendly/each/chosen tint an
// economy glyph carries.

// targetGlyph renders a Target as its noun glyph plus the decorations that carry
// its enemy/friendly/each/chosen shape. Fine filters (power, house, trait) stay
// in the rules text; the strip summarises the noun.
func targetGlyph(t engine.Target) glyph {
	switch t.Kind {
	case engine.TargetThisCreature:
		return glyph{
			asset: "type-creature",
			decor: decorThis,
		}
	case engine.TargetTriggeringCreature, engine.TargetTheOtherCreature,
		engine.TargetTheChosenCreature, engine.TargetCreatureFought,
		engine.TargetTheFoughtCreature, engine.TargetTheSameCreature,
		engine.TargetAttachedHost:
		return glyph{
			asset: "type-creature",
			decor: decorChosen,
		}
	case engine.TargetEachCreature:
		return glyph{
			asset: "type-creature",
			decor: decorEach,
		}
	case engine.TargetEachNeighbor, engine.TargetFormerNeighbors:
		// The strip carries "each creature"; which creatures are neighbors stays in
		// the rules text, the way the trigger pass folds AfterNeighborFights into the
		// plain fight glyph.
		return glyph{
			asset: "type-creature",
			decor: decorEach,
		}
	case engine.TargetEachUpgradeOnThis:
		return glyph{
			asset: "type-upgrade",
			decor: decorEach | decorThis,
		}
	case engine.TargetGrantingCard:
		return glyph{
			asset: "card-back",
			decor: decorChosen,
		}
	case engine.TargetEachFriendlyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorEach | decorFriendly,
		}
	case engine.TargetEachEnemyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorEach | decorEnemy,
		}
	case engine.TargetChosenCreature:
		return glyph{
			asset: "type-creature",
			decor: decorChosen,
		}
	case engine.TargetChosenFriendlyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorChosen | decorFriendly,
		}
	case engine.TargetChosenEnemyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorChosen | decorEnemy,
		}
	case engine.TargetEachArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorEach,
		}
	case engine.TargetEachFriendlyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorEach | decorFriendly,
		}
	case engine.TargetEachEnemyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorEach | decorEnemy,
		}
	case engine.TargetChosenArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorChosen,
		}
	case engine.TargetChosenFriendlyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorChosen | decorFriendly,
		}
	case engine.TargetChosenEnemyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorChosen | decorEnemy,
		}
	case engine.TargetChosenUpgrade:
		return glyph{
			asset: "type-upgrade",
			decor: decorChosen,
		}
	case engine.TargetEachCardInPlay, engine.TargetEachFriendlyCardInPlay,
		engine.TargetChosenCreatureOrArtifact,
		engine.TargetChosenFriendlyCreatureOrArtifact,
		engine.TargetChosenEnemyCreatureOrArtifact:
		return glyph{
			asset: "card-back",
			decor: cardInPlayDecor(t.Kind),
		}
	default:
		return glyph{text: t.Text()}
	}
}

// spreadTargetGlyph is the creature glyph a DealDamage Spread hits. Every spread
// chooses one or more creatures; the amounts and neighbor split stay in the rules
// text, so the strip summarises the spread as its chosen-creature noun.
func spreadTargetGlyph() glyph {
	return glyph{
		asset: "type-creature",
		decor: decorChosen,
	}
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

// verbGlyphs renders the verbs a chosen-creature effect applies in order — ready,
// fight, use, stun, exhaust, gain-keyword — reusing the same glyphs those actions
// draw on their own. It stays here rather than in a family file because more than
// one family (creature-state, composition) reads it.
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

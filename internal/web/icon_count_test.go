package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// countRows returns the Count family's catalogued rows, so a test reads the
// census rather than a second list of counts that could drift from it.
func countRows(t *testing.T) []iconRow {
	t.Helper()
	for _, family := range iconCatalog() {
		if family.name == "Count" {
			return family.rows
		}
	}
	t.Fatal("the icon catalog walks no Count family")
	return nil
}

// TestCountGlyphAssetsExist walks every catalogued count and checks each glyph it
// draws names a real asset. A glyph with a missing stem renders as a broken image
// rather than failing anything, so without this a typo in a noun ships silently.
func TestCountGlyphAssetsExist(t *testing.T) {
	for _, row := range countRows(t) {
		gs, _ := row.draw()
		for _, g := range gs {
			if !assetExists(t, g.asset) {
				t.Errorf("count %s draws %q, which has no web/assets/%s.svg",
					row.name, g.asset, g.asset)
			}
		}
	}
}

// TestFixedCountFolds pins the one plumbing row of the family: a literal draws
// nothing, because the effect holding it already prints its own number, and a
// noun before that number would say the count measured something.
func TestFixedCountFolds(t *testing.T) {
	gs, covered := countGlyphs(engine.Fixed(3))
	if !covered {
		t.Error("a fixed count should be covered, not drawn as the unknown glyph")
	}
	if len(gs) != 0 {
		t.Errorf("a fixed count draws %+v, want nothing", gs)
	}
}

// TestCountGlyphsTranscription binds the count grammar to worked examples so it
// cannot drift silently: a count over a set of cards stacks its noun and tints it
// with the side it surveys, a count reading the card holding it marks that noun
// decorThis, and a threshold inside a count rides the numeral badge with its
// comparison mark rather than drawing a glyph of its own (ADR 0047).
func TestCountGlyphsTranscription(t *testing.T) {
	// "for each enemy creature in play" → creature[each, enemy].
	gs, covered := countGlyphs(engine.CardsInPlay{
		Player: engine.Opponent,
		Type:   engine.Creature,
	})
	if !covered || len(gs) != 1 {
		t.Fatalf("enemy creatures in play = %+v (covered=%v), want one glyph", gs, covered)
	}
	if gs[0].asset != "type-creature" ||
		gs[0].decor&decorEach == 0 || gs[0].decor&decorEnemy == 0 {
		t.Errorf("glyph = %+v, want a stacked enemy creature", gs[0])
	}

	// "for each friendly creature with power 5 or higher" → creature[each] ≥5 power.
	gs, _ = countGlyphs(engine.CardsInPlay{
		Player:   engine.Controller,
		Type:     engine.Creature,
		MinPower: 5,
	})
	if len(gs) != 2 {
		t.Fatalf("power-floored count = %+v, want the noun and its threshold", gs)
	}
	if gs[1].asset != "power" || gs[1].qty != 5 || gs[1].comparison != engine.AtLeast {
		t.Errorf("threshold glyph = %+v, want power ≥5", gs[1])
	}

	// "for each Æmber on it" → aember[this], not a bare Æmber glyph.
	gs, _ = countGlyphs(engine.AemberOnThis{})
	if len(gs) != 1 || gs[0].asset != "aember" || gs[0].decor&decorThis == 0 {
		t.Errorf("Æmber on this = %+v, want the Æmber noun marked as this card", gs)
	}
}

// TestEveryCountDraws holds the family's gate from the other side: every
// catalogued count draws a real strip rather than the abstract unknown glyph.
// TestIconTotality already fails on a rule-bearing row that does not, but it
// skips the plumbing rows, and this says the same thing in the family's own
// terms — a count with no drawing says nothing about what it measures.
func TestEveryCountDraws(t *testing.T) {
	for _, row := range countRows(t) {
		gs, covered := row.draw()
		if !covered {
			t.Errorf("count %s draws no glyph of its own (ADR 0022)", row.name)
			continue
		}
		for _, g := range gs {
			if g.asset == "glyph-unknown" {
				t.Errorf("count %s draws the unknown glyph inside its strip", row.name)
			}
		}
	}
}

// turnStats lists every tally engine.TurnStat declares. It is a local list
// because TurnStat is not one of the engine's catalogued enums (engine.Enums()),
// so there is no enumerating function to walk; a tally added to the engine has to
// be added here and to turnStatGlyphs together.
func turnStats() []engine.TurnStat {
	return []engine.TurnStat{
		engine.KeysForgedThisTurn,
		engine.KeysForgedLastTurn,
		engine.CreaturesPlayedLastTurn,
		engine.EnemyCreaturesFightKilled,
		engine.EnemyCreaturesDestroyed,
		engine.FriendlyCreaturesDestroyed,
		engine.CreaturesReapedThisTurn,
		engine.CreaturesFoughtThisTurn,
		engine.CreaturesUsedThisTurn,
		engine.AemberStolenFromThisTurn,
		engine.AemberStolenFromLastTurn,
	}
}

// TestTurnStatGlyphsCoverEveryTally walks the turn tallies and fails on one that
// draws nothing. A TurnCount over an undrawn tally renders the unknown glyph, and
// only one tally is catalogued, so without this the other ten could go undrawn
// behind a passing gate.
func TestTurnStatGlyphsCoverEveryTally(t *testing.T) {
	for _, stat := range turnStats() {
		gs := turnStatGlyphs(stat)
		if len(gs) == 0 {
			t.Errorf("turn tally %d draws nothing; give it a noun and a verb (ADR 0022)", stat)
			continue
		}
		for _, g := range gs {
			if !assetExists(t, g.asset) {
				t.Errorf("turn tally %d draws %q, which has no asset", stat, g.asset)
			}
		}
	}
}

// TestProducedTallyGlyphsCoverEveryTally walks the "... this way" tallies a
// ProducedThisWay can name and fails on one with no drawing, for the same reason
// the turn tallies are walked: the census catalogues a single representative, so
// the other four would otherwise be undrawn behind a passing gate.
func TestProducedTallyGlyphsCoverEveryTally(t *testing.T) {
	tallies := []engine.ProducedTally{
		engine.TallyCreaturesDestroyed,
		engine.TallyCreaturesShuffledIntoDeck,
		engine.TallyAemberLost,
		engine.TallyCardsReturned,
		engine.TallyCardsPurged,
	}
	for _, tally := range tallies {
		gs, covered := producedTallyGlyphs(engine.ProducedThisWay{
			Tally:  tally,
			Player: engine.Controller,
		})
		if !covered {
			t.Errorf("produced tally %d draws no glyph of its own (ADR 0022)", tally)
			continue
		}
		for _, g := range gs {
			if !assetExists(t, g.asset) {
				t.Errorf("produced tally %d draws %q, which has no asset", tally, g.asset)
			}
		}
	}
}

// TestExcessCreaturesDrawsItsComparison pins the one count whose whole meaning is
// a comparison: more creatures than the other player has. The ">" rides the
// noun's badge and the side holding the excess is carried by the tint, so the
// strip says which side is ahead once rather than twice (ADR 0047).
func TestExcessCreaturesDrawsItsComparison(t *testing.T) {
	enemy, _ := countGlyphs(engine.ExcessCreatures{Player: engine.Opponent})
	if len(enemy) != 1 {
		t.Fatalf("excess creatures = %+v, want one glyph", enemy)
	}
	if enemy[0].decor&decorEnemy == 0 {
		t.Errorf("glyph = %+v, want the enemy side tinted", enemy[0])
	}
	if comparisonMark(enemy[0].comparison) != ">" {
		t.Errorf("glyph = %+v, want the more-than mark on its badge", enemy[0])
	}
	mine, _ := countGlyphs(engine.ExcessCreatures{Player: engine.Controller})
	if len(mine) != 1 || mine[0].decor&decorEnemy != 0 {
		t.Errorf("your own excess = %+v, want an untinted creature", mine)
	}
}

// TestHouseCountsShareOneDrawing pins what the two house counts are: one
// question — how many houses are represented — asked of two different sets of
// cards. They draw the same stacked house and differ only in the subject beside
// it, rather than each inventing a picture of its own.
func TestHouseCountsShareOneDrawing(t *testing.T) {
	inPlay, _ := countGlyphs(engine.HousesInPlay{})
	among, _ := countGlyphs(engine.HousesAmong{
		Player: engine.Controller,
		Type:   engine.Creature,
	})
	if len(inPlay) != 2 || len(among) != 2 {
		t.Fatalf("house counts = %+v and %+v, want a house and its subject", inPlay, among)
	}
	if inPlay[0] != among[0] || inPlay[0].asset != "house-any" {
		t.Errorf("house nouns differ: %+v and %+v", inPlay[0], among[0])
	}
	if inPlay[1].asset == among[1].asset {
		t.Errorf("both house counts survey %q; the subject is what separates them",
			among[1].asset)
	}
}

// TestPurgeCountsDrawApart pins the pair the task of drawing them flagged: the
// cards sitting in the purge pile and the cards one purge removed this
// resolution. The pile draws as the zone, the tally as the cards it took, so a
// count of what was purged reads differently from a count of what is purged.
func TestPurgeCountsDrawApart(t *testing.T) {
	pile, _ := countGlyphs(engine.PurgedCards{})
	tally, _ := countGlyphs(engine.CardsPurged{})
	if len(pile) != 1 || pile[0].asset != "zone-purge" {
		t.Errorf("the purge pile = %+v, want the zone alone", pile)
	}
	if len(tally) != 2 || tally[0].asset != "card-back" {
		t.Errorf("cards purged this way = %+v, want the cards and the zone", tally)
	}
}

// TestKeyCountsDrawApart pins the pair that would otherwise be one picture: the
// keys a player has forged and the keys they have still to forge are opposite
// halves of the same economy, and the strip separates them the way the player bar
// does — the forge for a key already forged, the bare key for one still unforged.
func TestKeyCountsDrawApart(t *testing.T) {
	forged, _ := countGlyphs(engine.ForgedKeys{Player: engine.Controller})
	unforged, _ := countGlyphs(engine.UnforgedKeys{Player: engine.Controller})
	if len(forged) != 1 || len(unforged) != 1 {
		t.Fatalf("key counts = %+v and %+v, want one glyph each", forged, unforged)
	}
	if forged[0].asset == unforged[0].asset {
		t.Errorf("forged and unforged keys both draw %q", forged[0].asset)
	}
}

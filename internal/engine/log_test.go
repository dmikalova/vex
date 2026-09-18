package engine

import (
	"fmt"
	"testing"
)

// stubNamer names cards and players from their ids alone, so an entry's own
// wording is what a test asserts on rather than a match's card pool.
type stubNamer struct{}

func (stubNamer) Name(id LocalID) string       { return fmt.Sprintf("Card%d", id) }
func (stubNamer) PlayerName(player int) string { return fmt.Sprintf("P%d", player) }

// TestLogEntryText pins the wording of every log entry variant. The log is the
// public narration of a match (ADR 0011), so each entry states a bound, past
// tense outcome, and a change to that wording should have to be made on purpose.
func TestLogEntryText(t *testing.T) {
	n := stubNamer{}
	samples := LogEntrySamples()
	// want holds the pinned wording of each sample in LogEntrySamples, index for
	// index. A new variant appends a sample there and its wording here; the length
	// check below fails loudly if the two ever drift out of alignment (ADR 0046).
	want := []string{
		// Turn shape.
		"P0 begins turn 3",
		"P0 takes the first turn",
		"P1 mulligans, drawing a new hand of 5",
		"Ready phase",
		"P0 readies Card4, Card7",
		"P1 draws 3 cards, up to 6 in hand",
		"P1 draws 1 card, up to 6 in hand",
		"P0 draws nothing, holding 2",
		"P0 draws 1 card",
		"P1 chooses house Brobnar",
		"P0 skips their forge a key phase",
		"Card5 gains the 6 Æmber P1 spends forging a key",
		"P0 forges a Red key (1/3)",
		"P1 forges a Colorless key (4/3)",
		"P0 unforges a key (1/3)",
		"P1 sheds a chain (4 remaining)",
		"P0 wins the game!",
		"P1 concedes",
		"P0 has 4 Æmber and 1 key",

		// Æmber. The source-card subject is exercised in
		// TestRecordTextSubjectsToSourceCard; here the entries render unframed, so
		// they name the player.
		"P0 gains 2 Æmber",
		"P1 loses 1 Æmber",
		"P0 steals 2 Æmber from P1",
		"Card9 has P0 steal 2 Æmber from the common supply, instead of from P1's pool",
		"Card7 captures 3 Æmber",
		"Card3 captures 1 Æmber onto Card7",
		"Card9 has Card7 capture 3 Æmber from the common supply, " +
			"instead of from P0's pool",
		"Card9 has Card3 capture 1 Æmber onto Card7 from the common supply, " +
			"instead of from P0's pool",
		"Card7 moves 1 Æmber to the common supply",
		"Card7 captures 1 Æmber, instead of P1 gaining it",
		"Card9 has Card7 capture 2 Æmber, instead of P0 stealing it",
		"Card9 has Card7 capture 2 Æmber from the common supply, " +
			"instead of P0 stealing it",
		"P0 exalts 2 Æmber onto Card4",
		"P0 moves 2 Æmber from Card4 to P1's pool",
		"P0 moves 1 Æmber from Card4 to Card5",
		"Maximum amount reached, 2 Æmber is lost",

		// Creatures and cards in play.
		"Card2 is readied",
		"Card2 gains skirmish",
		"Card2 loses elusive",
		"Card2 gains +1 armor",
		"Card2 gains +2 power and +2 armor",
		"Card2 gains assault 3",
		"Card2 gains the Mutant trait",
		"Card2 is considered a flank creature",
		"Card2 is exhausted",
		"Card2 is stunned",
		"Card5 stunned Card2",
		"Card2 is already stunned",
		"P0 unstuns Card2 and Card5",
		"Card2 is enraged",
		"Card5 enraged Card2",
		"Card2 is already enraged",
		"Card2 is warded",
		"Card5 warded Card2",
		"Card2 is already warded",
		"Card5 removes the ward from Card2",
		"Card2 has no ward to remove",
		"Card2's ward keeps it in play",
		"Card2's ward absorbs the destruction",
		"Card2's ward absorbs 5 damage",
		"Card2 has no creature to fight",
		"P0 reveals Card1, Card2",
		"P1's forge a key is prevented by Card3",
		"Card1 swaps positions with Card2",
		"Card1 swaps places with Card2 from P1's discard pile",
		"Card2 moves to the right flank",
		"Card2 moves to the left flank",
		"Card2 moves within its battleline",
		"Card2 becomes a creature on the right flank",
		"Card2 becomes a creature on the left flank",
		"Card2 reverts to an artifact",
		"P1 takes control of Card3",
		"Card3 returns to P0's control",
		"Card3 is destroyed",
		"Card5 destroys Card3",
		"Card5 destroys Card3 and Card8",
		"Card5 destroys Card3, Card8, and Card9",
		"Card3 would be destroyed, so Card8 replaces its destruction",
		"2 Æmber on Card3 goes to P1's pool",
		"P1 unstuns Card3",
		"Card3 is exhausted and cannot be used",

		// Copied stats and text box.
		"Card2 copies the printed stats of Card5",
		"Card2 gains the text box of Card5",

		// Combat.
		"Card1's fight does not occur",
		"Card1 (4 power) fights Card2 (3 power)",
		"Card1 (6 power, skirmish) fights Card2 (3 power, elusive, poison)",
		"Card1's poison is lethal to Card2",
		"Card2 cannot be dealt damage",
		"Card2's armor absorbs 1 damage",
		"Card2 takes 3 damage (4 total)",
		"Card1 deals 2 assault damage to Card2",
		"Card2 deals 5 hazardous damage to Card1",
		"Card2 takes 4 damage",

		// Zones.
		"P0 takes 2 cards from their archives into hand",
		"P0 archives a card",
		"P1 archives Card6 from their discard pile",
		"P0 archives a card from the top of their deck",
		"P1 discards 3 archived cards",
		"P0 discards Card6 from the top of their deck",
		"P0 discards Card6 from their deck",
		"P1 swaps their deck and discard pile",
		"P0 shuffles Card6 from their discard pile and 2 cards from their hand into their deck",
		"P1 shuffles Card4 and Card6 from their discard pile into their deck",
		"P0 shuffles 1 card from their archives into their deck",
		"P1 shuffles their deck",
		"P0 discards Card6",
		"P0 discards Card6 from their archives",
		"P0 purges Card6 from a discard pile",
		"P1 purges Card6 from a hand",
		"Card6 is purged from P0's hand",
		"P1 purges Card6 from archives",
		"P1 purges Card6 from a deck",
		"Card6 is purged",
		"Card6 is put on top of P0's deck",
		"Card6 is put into P1's hand",
		"Card6 is put into P0's archives",
		"P0 archives Card6 from their purge pile",
		"Card6 is shuffled into P1's deck",
		"P1's deck is shuffled",
		"P1's discard pile is shuffled into their deck",
		"P1 shuffles Card3 and Card8 into their deck",
		"P0 abducts Card6 (owned by P1) into their archives",
		"P0 puts Card6 from their discard pile into hand",
		"P1 puts a card from their deck into hand",
		"P0 puts Card6 from their discard pile on top of their deck",
		"P0 puts Card6 faceup under Card9",
		"P0 puts a card facedown under Card9",
		"Card6 is grafted onto Card9",

		// Playing and using.
		"P0 plays Card9 on their right flank",
		"P0 plays Card9 on their left flank",
		"P0 plays Card9 into their battleline",
		"P1 plays artifact Card9",
		"P0 plays tactic Card9",
		"P0 attaches Card9 to Card2",
		"Card9 is discarded from Card2",
		"P1 puts Card9 into play under their control",
		"P0 plays Card9 from the top of P0's deck",
		"Card9 has P0 gain 2 bonus Æmber",
		"Card7 captures 2 bonus Æmber from Card9, instead of P0 gaining it",
		"Card9 has Card7 bonus capture 1 Æmber",
		"Card9 deals 1 bonus damage to Card2",
		"Card9 has P0 draw 1 bonus card",
		"P0 loses 1 Æmber to play Card9",
		"P0 reaps with Card2 (+1 Æmber)",
		"Card9 has P0 steal 1 Æmber reaping with Card2, instead of gaining it",
		"Card9 has P0 steal 0 Æmber reaping with Card2, instead of gaining it",
		"Card7 captures 1 Æmber reaping with Card2, instead of P0 gaining it",
		"P1 uses Card2's action ability",

		// Lasting effects. A lasting Æmber gain records the plain AemberGained/
		// AemberCapturedInsteadOfGain entry under its source-card frame, so only the
		// draw keeps its own lasting entry (and its event context).
		"P1 draws 2 cards (each time a friendly creature fights)",
		"P0 gives 3 Æmber to P1 after forging a key",
		"P0 gives 1 Æmber to P1",
		"P0 gives 1 Æmber to P1 to use an artifact",

		// Grants, chains, and manual mode.
		"P0's Brobnar creatures may fight this turn",
		"P1's creatures may all fight this turn",
		"P0 may use Dis creatures this turn",
		"P0 may use friendly artifacts this turn",
		"P0 may play Mars cards this turn",
		"P0 may play or use Mars cards this turn",
		"P0 may play cards from other houses this turn",
		"P0 may play cards from other houses this turn",
		"P0 may use friendly Mutant creatures this turn",
		"P1 must choose house Logos next turn",
		"P1 cannot choose house Logos next turn",
		"P0 steals 2 Æmber if P1 chooses house Logos next turn",
		"each creature loses elusive for the remainder of the turn",
		"P0 gains 2 chains (5 total)",
		"P0 manually moves Card3 to purge",
		"Card3 is manually exhausted",
		"Card3 is manually readied",
		"P0 manually puts Card3 into play",
		"P1 cannot add a card: this match is full",
		"P0 manually adds Card3 to hand",
		"P0 now has 7 Æmber (manual)",
		"P1 now has 1 chain (manual)",
		"P0 manually chooses Mars as their active house",
		"P0 manually forges a Blue key (2/3)",
		"P1 manually unforges a key (0/3)",

		// A restored entry reads back exactly as it was narrated.
		"P0 gains 1 Æmber",
	}
	if len(want) != len(samples) {
		t.Fatalf("want has %d entries, LogEntrySamples has %d; keep them aligned",
			len(want), len(samples))
	}
	for i, e := range samples {
		if got := e.Text(n); got != want[i] {
			t.Errorf("sample %d: %T.Text() = %q, want %q", i, e, got, want[i])
		}
	}
}

// TestRecordTextSubjectsToSourceCard pins the wording of the entries whose
// subject is the source card the record's frame carries. Under a card ability the
// card is the subject ("Card7 deals 4 damage to Card2"); when the ability acts on
// the other player's pool the affected player is named too ("Card7 has P1 gain 2
// Æmber").
func TestRecordTextSubjectsToSourceCard(t *testing.T) {
	n := stubNamer{}
	// A frame opened for a card ability P0 controls, sourced to Card7.
	fr := Frame{Actor: 0, Source: 7, HasSource: true}
	cases := []struct {
		entry LogEntry
		want  string
	}{
		{AemberGained{Player: 0, Amount: 1}, "Card7 has P0 gain 1 Æmber"},
		{AemberGained{Player: 1, Amount: 2}, "Card7 has P1 gain 2 Æmber"},
		{AemberLost{Player: 1, Amount: 1}, "Card7 has P1 lose 1 Æmber"},
		{AemberStolen{Player: 0, From: 1, Amount: 2}, "Card7 steals 2 Æmber from P1"},
		// A redirected steal is narrated by the card that redirected it, not by the
		// frame's source: the replacement is what made the line worth printing.
		{
			AemberStolen{Player: 0, From: 1, Amount: 2, FromSupply: true, Cause: 9},
			"Card9 has P0 steal 2 Æmber from the common supply, instead of from P1's pool",
		},
		{AbilityDamageDealt{Amount: 4, Target: 2}, "Card7 deals 4 damage to Card2"},
		{
			TopOfDeckDiscarded{Player: 1, Card: 6},
			"Card7 discards Card6 from the top of P1's deck",
		},
		{
			PlayedFromTopOfDeck{Card: 9, Player: 0},
			"Card7 plays Card9 from the top of P0's deck",
		},
		{
			ChainsGained{Player: 1, Amount: 2, Total: 2},
			"Card7 has P1 gain 2 chains (2 total)",
		},
		{
			CardMoved{Player: 0, Card: 6, From: Discard, To: Archives},
			"Card7 archives Card6 from P0's discard pile",
		},
		{
			CardMoved{Player: 1, Card: 6, From: Deck, To: Discard},
			"Card7 discards Card6 from P1's deck",
		},
		{
			TopOfDeckArchived{Player: 0, Card: 6},
			"Card7 archives a card from the top of P0's deck",
		},
		{DeckAndDiscardSwapped{Player: 1}, "Card7 swaps P1's deck and discard pile"},
		{
			ShuffledIntoDeck{Player: 0, DiscardCards: []LocalID{6}, HandCount: 2},
			"Card7 shuffles Card6 from P0's discard pile and 2 cards from P0's hand into P0's deck",
		},
		{ShuffledIntoDeck{Player: 1}, "Card7 shuffles P1's deck"},
		{
			CardPutFromDiscardIntoHand{Player: 0, Card: 6},
			"Card7 puts Card6 from P0's discard pile into hand",
		},
		{
			CardPutFromDeckIntoHand{Player: 1, Card: 6},
			"Card7 puts a card from P1's deck into hand",
		},
		{
			CardPutFromDiscardOnTopOfDeck{Player: 0, Card: 6},
			"Card7 puts Card6 from P0's discard pile on top of P0's deck",
		},
		{
			CardArchivedFromPurge{Player: 0, Card: 6},
			"Card7 archives Card6 from P0's purge pile",
		},
		{
			CardPlayedToBattleline{Player: 0, Card: 9},
			"Card7 plays Card9 on P0's right flank",
		},
		{
			CardPlayedToBattleline{Player: 0, Card: 9, Interior: true},
			"Card7 plays Card9 into P0's battleline",
		},
		{
			CardPutIntoPlay{Player: 1, Card: 9},
			"Card7 puts Card9 into play under P1's control",
		},
		{
			CreaturesUnstunned{Player: 0, Creatures: []LocalID{2, 5}},
			"Card7 unstuns Card2 and Card5",
		},
		{CardsDrawnBy{Player: 0, Count: 1}, "Card7 has P0 draw 1 card"},
		{
			AemberGiven{Giver: 0, Receiver: 1, Amount: 1, Reason: TollUseArtifact},
			"Card7 has P0 give 1 Æmber to P1 to use an artifact",
		},
		{
			CardDiscarded{Player: 0, Card: 6},
			"Card7 discards Card6",
		},
		{
			CardPurgedFromHand{Card: 6, Owner: 0},
			"Card7 purges Card6 from P0's hand",
		},
		{
			CardsShuffledIntoDeckBy{Owner: 1, Cards: []LocalID{3, 8}},
			"Card7 shuffles Card3 and Card8 into P1's deck",
		},
	}
	for _, c := range cases {
		if got := (Record{Frame: fr, Entry: c.entry}).Text(n); got != c.want {
			t.Errorf("%T framed Text() = %q, want %q", c.entry, got, c.want)
		}
	}
}

// TestRenderRecordSubjectsToSourceCard checks that rendering a record under a
// frame that carries a source card splits the source card out as the subject
// segment, so a client links it (ADR 0011).
func TestRenderRecordSubjectsToSourceCard(t *testing.T) {
	rec := Record{
		Frame: Frame{Actor: 0, Source: 7, HasSource: true},
		Entry: AemberGained{Player: 0, Amount: 1},
	}
	segs := RenderRecord(rec, stubNamer{})
	want := []LogSegment{
		{Text: "Card7", Card: 7, HasCard: true},
		{Text: " has "},
		{Text: "P0", Player: 0, HasPlayer: true},
		{Text: " gain 1 "},
		{Text: "Æmber", Icon: "aember"},
	}
	if len(segs) != len(want) {
		t.Fatalf("segments = %+v, want %+v", segs, want)
	}
	for i := range want {
		if segs[i] != want[i] {
			t.Errorf("segment %d = %+v, want %+v", i, segs[i], want[i])
		}
	}
}

// TestRenderEntrySplitsOutCardNames checks that a client can find the cards an
// entry names from the ids the entry carries, rather than by matching its prose
// against a card index (ADR 0011).
func TestRenderEntrySplitsOutCardNames(t *testing.T) {
	segs := RenderEntry(PositionsSwapped{A: 1, B: 2}, stubNamer{})
	want := []LogSegment{
		{Text: "Card1", Card: 1, HasCard: true},
		{Text: " swaps positions with "},
		{Text: "Card2", Card: 2, HasCard: true},
	}
	if len(segs) != len(want) {
		t.Fatalf("segments = %+v, want %+v", segs, want)
	}
	for i := range want {
		if segs[i] != want[i] {
			t.Errorf("segment %d = %+v, want %+v", i, segs[i], want[i])
		}
	}
	// An entry that names no card still marks the keywords it contains, but "key
	// phase" names the turn phase rather than an actual key, so it stays plain.
	segs = RenderEntry(ForgeSkipped{Player: 0}, stubNamer{})
	wantPlain := []LogSegment{
		{Text: "P0", Player: 0, HasPlayer: true},
		{Text: " skips their forge a key phase"},
	}
	if len(segs) != len(wantPlain) {
		t.Fatalf("segments = %+v, want %+v", segs, wantPlain)
	}
	for i := range wantPlain {
		if segs[i] != wantPlain[i] {
			t.Errorf("segment %d = %+v, want %+v", i, segs[i], wantPlain[i])
		}
	}
}

// prefixNamer names every card the same, so a card name that is only part of a
// longer word is not mistaken for a whole one.
type prefixNamer struct{}

func (prefixNamer) Name(id LocalID) string {
	if id == 1 {
		return "Troll"
	}
	return "Trollkin"
}

func (prefixNamer) PlayerName(int) string { return "Trollkin" }

// TestRenderEntryMatchesWholeNamesOnly checks that a card name is only linked
// where it stands as a whole word: not when a longer name starts at the same
// place, and not when it is merely the start of some longer word.
func TestRenderEntryMatchesWholeNamesOnly(t *testing.T) {
	segs := RenderEntry(PositionsSwapped{A: 2, B: 1}, prefixNamer{})
	if len(segs) != 3 || segs[0].Text != "Trollkin" || segs[2].Text != "Troll" {
		t.Fatalf("segments = %+v, want Trollkin then Troll", segs)
	}
	// "Trollkin discards Troll": the player's name only starts with the card's,
	// so the card is not linked until the card itself.
	segs = RenderEntry(CardDiscarded{Player: 0, Card: 1}, prefixNamer{})
	want := []LogSegment{
		{Text: "Trollkin", Player: 0, HasPlayer: true},
		{Text: " discards "},
		{Text: "Troll", Card: 1, HasCard: true},
	}
	if len(segs) != len(want) || segs[0] != want[0] || segs[1] != want[1] ||
		segs[2] != want[2] {
		t.Fatalf("segments = %+v, want %+v", segs, want)
	}
	// A card with no name contributes no span.
	segs = RenderEntry(CardPurged{Card: 3}, namelessNamer{})
	if len(segs) != 1 || segs[0].HasCard {
		t.Errorf("segments = %+v, want one plain run for a nameless card", segs)
	}
}

// TestRenderEntryMarksIconKeywords checks that the closed keyword vocabulary —
// Æmber and the houses — comes back as icon segments, and that a name the entry
// asked for wins over a keyword sitting in the same place.
func TestRenderEntryMarksIconKeywords(t *testing.T) {
	segs := RenderEntry(HouseChosen{Player: 1, House: Brobnar}, stubNamer{})
	want := []LogSegment{
		{Text: "P1", Player: 1, HasPlayer: true},
		{Text: " chooses house "},
		{Text: "Brobnar", Icon: "house-brobnar"},
	}
	if len(segs) != len(want) || segs[0] != want[0] || segs[1] != want[1] ||
		segs[2] != want[2] {
		t.Fatalf("segments = %+v, want %+v", segs, want)
	}
	segs = RenderEntry(AemberGained{Player: 0, Amount: 2}, stubNamer{})
	var icons int
	for _, s := range segs {
		if s.Icon == "aember" && s.Text == "Æmber" {
			icons++
		}
	}
	if icons != 1 {
		t.Errorf("segments = %+v, want one Æmber icon segment", segs)
	}
	// A card named "Brobnar" is a card link, not a house emblem.
	segs = RenderEntry(CardPurged{Card: 1}, houseNamer{})
	if len(segs) != 2 || !segs[0].HasCard || segs[0].Icon != "" {
		t.Errorf("segments = %+v, want the card name to beat the house keyword", segs)
	}
	// Stunning, chains, and keys all mark their own icon.
	iconCases := []struct {
		entry LogEntry
		icon  string
		word  string
	}{
		{CreatureStunned{Creature: 2, By: 5}, "stun", "stunned"},
		{ChainsGained{Player: 0, Amount: 2, Total: 5}, "chains", "chains"},
		{ChainShed{Player: 1, Remaining: 4}, "chains", "chain"},
		{KeyForged{Player: 0, Color: KeyColorColorless, Keys: 4, Needed: 3}, "key", "key"},
		{
			KeyForged{Player: 0, Color: KeyColorRed, Keys: 1, Needed: 3},
			"key-red", "key",
		},
		{PlayerStanding{Player: 0, Aember: 4, KeyColors: []KeyColor{KeyColorRed}}, "key", "key"},
		{CardShuffledIntoDeck{Card: 6, Owner: 1}, "zone-deck", "deck"},
	}
	for _, c := range iconCases {
		segs := RenderEntry(c.entry, stubNamer{})
		var found bool
		for _, s := range segs {
			if s.Icon == c.icon && s.Text == c.word {
				found = true
			}
		}
		if !found {
			t.Errorf("segments = %+v, want a %q icon segment for %q", segs, c.icon, c.word)
		}
	}
}

// houseNamer names a card after a house, so a keyword and a card name compete
// for the same span.
type houseNamer struct{ stubNamer }

func (houseNamer) Name(LocalID) string { return "Brobnar" }

// TestZoneIconAt checks each zone noun is marked with its zone emblem, and that a
// verb sharing a stem — "discards", "purges" — stays plain text.
func TestZoneIconAt(t *testing.T) {
	cases := []struct {
		word string
		icon string
	}{
		{"hand", "zone-hand"},
		{"deck", "zone-deck"},
		{"discard", "zone-discard"},
		{"archives", "zone-archives"},
		{"purge", "zone-purge"},
	}
	for _, c := range cases {
		seg, ok := zoneIconAt(c.word, 0)
		if !ok || seg.Icon != c.icon || seg.Text != c.word {
			t.Errorf("zoneIconAt(%q) = %+v, %v; want icon %q", c.word, seg, ok, c.icon)
		}
	}
	if _, ok := zoneIconAt("discards a card", 0); ok {
		t.Error("the verb \"discards\" should not be marked as the discard zone")
	}
}

// TestWordAt checks the whole-word guard directly, including a keyword that
// begins in the middle of a longer word, so it is not marked as its own word.
func TestWordAt(t *testing.T) {
	if !wordAt("has chains now", 4, "chains") {
		t.Error("a standalone word should match")
	}
	if wordAt("unchained", 2, "chain") {
		t.Error("a word preceded by a letter should not match")
	}
	if wordAt("chained", 0, "chain") {
		t.Error("a word followed by a letter should not match")
	}
	if wordAt("chain", 0, "") {
		t.Error("an empty word never matches")
	}
}

// namelessNamer names nothing, standing in for a card the reader may not see.
type namelessNamer struct{ stubNamer }

func (namelessNamer) Name(LocalID) string { return "" }

// TestFramesNestAndPop checks that entries inherit the attribution in force and
// that a closed frame stops applying.
func TestFramesNestAndPop(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.record(ForgeSkipped{Player: 0})
	outer := g.openFrame(Frame{Actor: 1, Source: 5, HasSource: true, Trigger: TriggerAfterReap})
	g.record(AemberGained{Player: 1, Amount: 1})
	inner := g.openFrame(Frame{Actor: 1, Source: 6, HasSource: true, Trigger: TriggerDestroyed})
	g.record(CardDestroyed{Card: 6})
	inner()
	g.record(AemberGained{Player: 1, Amount: 2})
	outer()
	g.record(TurnBegan{Player: 1, Turn: 2})

	wantFrames := []Frame{
		{Actor: 0},
		{Actor: 1, Source: 5, HasSource: true, Trigger: TriggerAfterReap},
		{Actor: 1, Source: 6, HasSource: true, Trigger: TriggerDestroyed, Depth: 1},
		{Actor: 1, Source: 5, HasSource: true, Trigger: TriggerAfterReap},
		{Actor: 0},
	}
	if len(g.Log) != len(wantFrames) {
		t.Fatalf("log = %v, want %d entries", g.LogText(), len(wantFrames))
	}
	for i, want := range wantFrames {
		if g.Log[i].Frame != want {
			t.Errorf("entry %d frame = %+v, want %+v", i, g.Log[i].Frame, want)
		}
	}
	// A record renders through to its entry.
	if got := g.Log[0].Text(g); got != "A skips their forge a key phase" {
		t.Errorf("record text = %q", got)
	}
}

// TestRestoreReplacesTheLog checks the seam a frontend uses to put a persisted
// log back after a reload.
func TestRestoreReplacesTheLog(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.record(ForgeSkipped{Player: 0})
	g.Restore([]Record{{Entry: RestoredEntry{Line: "restored"}}})
	if got := g.LogText(); len(got) != 1 || got[0] != "restored" {
		t.Errorf("log = %v, want the restored line", got)
	}
}

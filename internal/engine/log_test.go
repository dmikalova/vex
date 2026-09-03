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
	cases := []struct {
		entry LogEntry
		want  string
	}{
		// Turn shape.
		{TurnBegan{Player: 0, Turn: 3}, "P0 begins turn 3"},
		{PhaseBegan{Player: 1, Phase: PhaseReady}, "Ready phase"},
		{CardsReadied{Player: 0, Cards: []LocalID{4, 7}}, "P0 readies Card4, Card7"},
		{CardsDrawn{Player: 1, Count: 3, Hand: 6}, "P1 draws 3 cards, up to 6 in hand"},
		{CardsDrawn{Player: 1, Count: 1, Hand: 6}, "P1 draws 1 card, up to 6 in hand"},
		{CardsDrawn{Player: 0, Count: 0, Hand: 2}, "P0 draws nothing, holding 2"},
		{HouseChosen{Player: 1, House: Brobnar}, "P1 chooses house Brobnar"},
		{ForgeSkipped{Player: 0}, "P0 skips their forge a key phase"},
		{
			KeyForged{Player: 0, Color: KeyColorRed, HasColor: true, Keys: 1, Needed: 3},
			"P0 forges a Red key (1/3)",
		},
		{KeyForged{Player: 1, Keys: 2, Needed: 3}, "P1 forges a key (2/3)"},
		{KeyUnforged{Player: 0, Keys: 1, Needed: 3}, "P0 unforges a key (1/3)"},
		{ChainShed{Player: 1, Remaining: 4}, "P1 sheds a chain (4 remaining)"},
		{GameWon{Player: 0}, "P0 wins the game!"},
		{PlayerStanding{Player: 0, Aember: 4, Keys: 1}, "P0 has 4 Æmber and 1 keys"},

		// Æmber.
		{AemberGained{Player: 0, Amount: 2}, "P0 gains 2 Æmber"},
		{AemberLost{Player: 1, Amount: 1}, "P1 loses 1 Æmber"},
		{AemberStolen{Player: 0, From: 1, Amount: 2}, "P0 steals 2 Æmber from P1"},
		{AemberCaptured{Creature: 7, Amount: 3}, "Card7 captures 3 Æmber"},
		{
			AemberCapturedInsteadOfGain{Creature: 7, Player: 1, Amount: 1},
			"Card7 captures 1 Æmber instead of P1 gaining it",
		},
		{AemberExalted{Creature: 4, Amount: 2}, "Card4 is exalted (2 Æmber placed)"},
		{
			AemberMovedToPool{Player: 0, From: 4, To: 1, Amount: 2},
			"P0 moves 2 Æmber from Card4 to P1's pool",
		},
		{
			AemberMovedToCard{Player: 0, From: 4, To: 5, Amount: 1},
			"P0 moves 1 Æmber from Card4 to Card5",
		},
		{
			AemberLostToCeiling{Card: 4, Amount: 2},
			"Card4 can hold no more Æmber; 2 is lost to the ceiling",
		},

		// Creatures and cards in play.
		{CreatureReadied{Creature: 2}, "Card2 is readied"},
		{CreatureExhausted{Creature: 2}, "Card2 is exhausted"},
		{NoCreatureToFight{Creature: 2}, "Card2 has no creature to fight"},
		{CardsRevealedToAll{Player: 0, Cards: []LocalID{1, 2}}, "P0 reveals Card1, Card2"},
		{PositionsSwapped{A: 1, B: 2}, "Card1 swaps positions with Card2"},
		{ControlTaken{Player: 1, Card: 3}, "P1 takes control of Card3"},
		{ControlReturned{Card: 3, Owner: 0}, "Card3 returns to P0's control"},
		{CardDestroyed{Card: 3}, "Card3 is destroyed"},
		{
			DestructionReplaced{Card: 3, By: 8},
			"Card3 would be destroyed, so Card8 replaces its destruction",
		},
		{
			AemberOnCardReleased{Card: 3, Amount: 2, To: 1},
			"2 Æmber on Card3 goes to P1's pool",
		},
		{StunRecovered{Creature: 3}, "Card3 recovers from stun instead of acting"},
		{CardCannotBeUsed{Card: 3}, "Card3 is exhausted and cannot be used"},

		// Combat.
		{FightCancelled{Attacker: 1}, "Card1's fight does not occur"},
		{
			Fought{Attacker: 1, AttackerPower: 4, Defender: 2, DefenderPower: 3},
			"Card1 (4 power) fights Card2 (3 power)",
		},
		{ElusiveAvoidedFight{Defender: 2}, "Card2 is elusive — no fight damage is dealt"},
		{DamageRefused{Creature: 2}, "Card2 cannot be dealt damage"},
		{ArmorAbsorbed{Creature: 2, Amount: 1}, "Card2's armor absorbs 1 damage"},
		{DamageTaken{Creature: 2, Amount: 3, Total: 4}, "Card2 takes 3 damage (4 total)"},

		// Zones.
		{
			ArchivesTakenIntoHand{Player: 0, Count: 2},
			"P0 takes 2 cards from their archives into hand",
		},
		{CardArchivedFromHand{Player: 0, Card: 6}, "P0 archives a card"},
		{
			CardArchivedFromDiscard{Player: 1, Card: 6},
			"P1 archives Card6 from their discard pile",
		},
		{TopOfDeckArchived{Player: 0, Card: 6}, "P0 archives a card from the top of their deck"},
		{ArchivesDiscarded{Player: 1, Count: 3}, "P1 discards 3 archived cards"},
		{
			TopOfDeckDiscarded{Player: 0, Card: 6},
			"P0 discards Card6 from the top of their deck",
		},
		{DeckAndDiscardSwapped{Player: 1}, "P1 swaps their deck and discard pile"},
		{CardDiscarded{Player: 0, Card: 6}, "P0 discards Card6"},
		{CardPurgedFromDiscard{Player: 0, Card: 6}, "P0 purges Card6 from a discard pile"},
		{CardPurgedFromHand{Player: 1, Card: 6}, "P1 purges Card6 from a hand"},
		{CardPurged{Card: 6}, "Card6 is purged"},
		{CardPutOnTopOfDeck{Card: 6, Owner: 0}, "Card6 is put on top of P0's deck"},
		{CardReturnedToHand{Card: 6, Owner: 1}, "Card6 is returned to P1's hand"},
		{CardPutIntoArchives{Card: 6, Owner: 0}, "Card6 is put into P0's archives"},
		{CardShuffledIntoDeck{Card: 6, Owner: 1}, "Card6 is shuffled into P1's deck"},
		{
			CardAbducted{Player: 0, Card: 6, Owner: 1},
			"P0 abducts Card6 (owned by P1) into their archives",
		},
		{
			CardReturnedFromDiscardToHand{Player: 0, Card: 6},
			"P0 returns Card6 from their discard pile to hand",
		},
		{
			CardPutFromDeckIntoHand{Player: 1, Card: 6},
			"P1 puts a card from their deck into hand",
		},
		{
			CardPutFromDiscardOnTopOfDeck{Player: 0, Card: 6},
			"P0 puts Card6 from their discard pile on top of their deck",
		},

		// Playing and using.
		{CardPlayedToBattleline{Player: 0, Card: 9}, "P0 plays Card9 to the battleline"},
		{ArtifactPlayed{Player: 1, Card: 9}, "P1 plays artifact Card9"},
		{ActionPlayed{Player: 0, Card: 9}, "P0 plays action Card9"},
		{UpgradeAttached{Player: 0, Upgrade: 9, Host: 2}, "P0 attaches Card9 to Card2"},
		{
			CardPutIntoPlay{Player: 1, Card: 9},
			"P1 puts Card9 into play under their control",
		},
		{AemberBonusGained{Player: 0, Card: 9, Amount: 2}, "P0 gains 2 Æmber from Card9"},
		{
			AemberBonusCaptured{Creature: 7, Card: 9, Amount: 2},
			"Card7 captures 2 Æmber from Card9's bonus",
		},
		{AemberSpentToPlay{Player: 0, Card: 9, Amount: 1}, "P0 loses 1 Æmber to play Card9"},
		{
			TollPaid{Player: 0, Payee: 1, Amount: 1, Action: TollUseArtifact},
			"P0 pays 1 Æmber to P1 to use an artifact",
		},
		{Reaped{Player: 0, Card: 2}, "P0 reaps with Card2 (+1 Æmber)"},
		{
			ReapedStealing{Player: 0, Card: 2, Amount: 1},
			"P0 reaps with Card2, stealing 1 Æmber",
		},
		{
			ReapedStealing{Player: 0, Card: 2},
			"P0 reaps with Card2 (no Æmber to steal)",
		},
		{
			ReapedCaptured{Player: 0, Card: 2, Creature: 7},
			"P0 reaps with Card2, but Card7 captures the Æmber",
		},
		{ActionAbilityUsed{Player: 1, Card: 2}, "P1 uses Card2's action ability"},

		// Lasting effects.
		{
			LastingAemberGained{Player: 0, Amount: 1, On: EventReap},
			"P0 gains 1 Æmber (after a creature reaps)",
		},
		{
			LastingAemberCaptured{Creature: 7, Player: 0, Amount: 1, On: EventForgeKey},
			"Card7 captures 1 Æmber instead of P0 gaining it (after forging a key)",
		},
		{
			LastingDraw{Player: 1, Amount: 2, On: EventFight},
			"P1 draws 2 cards (each time a friendly creature fights)",
		},
		{
			AemberGivenAfterForging{Player: 0, To: 1, Amount: 3},
			"P0 gives 3 Æmber to P1 after forging a key",
		},

		// Grants, chains, and manual mode.
		{
			FightGrantedForHouse{Player: 0, House: Brobnar},
			"P0's Brobnar creatures may fight this turn",
		},
		{FightGrantedAnyHouse{Player: 1}, "P1's creatures may all fight this turn"},
		{UseGrantedForHouse{Player: 0, House: Dis}, "P0 may use Dis creatures this turn"},
		{
			HouseForcedNextTurn{Player: 1, House: Logos},
			"P1 must choose house Logos next turn",
		},
		{
			KeywordLostByAll{Keyword: Elusive},
			"each creature loses elusive for the remainder of the turn",
		},
		{ChainsGained{Player: 0, Amount: 2, Total: 5}, "P0 gains 2 chains (5 total)"},
		{
			ManualCardMoved{Player: 0, Card: 3, To: ManualPurge},
			"P0 manually moves Card3 to purge",
		},
		{ManualExhaustSet{Card: 3, Exhausted: true}, "Card3 is manually exhausted"},
		{ManualExhaustSet{Card: 3}, "Card3 is manually readied"},
		{ManualMatchFull{Player: 1}, "P1 cannot add a card: this match is full"},
		{ManualCardAdded{Player: 0, Card: 3}, "P0 manually adds Card3 to hand"},
		{ManualAemberSet{Player: 0, Amount: 7}, "P0 now has 7 Æmber (manual)"},
		{ManualChainsSet{Player: 1, Amount: 1}, "P1 now has 1 chain (manual)"},
		{
			ManualHouseChosen{Player: 0, House: Mars},
			"P0 manually chooses Mars as their active house",
		},
		{
			ManualKeyForged{Player: 0, Color: KeyColorBlue, Keys: 2, Needed: 3},
			"P0 manually forges a Blue key (2/3)",
		},
		{ManualKeyUnforged{Player: 1, Keys: 0, Needed: 3}, "P1 manually unforges a key (0/3)"},

		// A restored entry reads back exactly as it was narrated.
		{RestoredEntry{Line: "P0 gains 1 Æmber"}, "P0 gains 1 Æmber"},
	}
	for _, c := range cases {
		if got := c.entry.Text(n); got != c.want {
			t.Errorf("%T.Text() = %q, want %q", c.entry, got, c.want)
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
	// An entry that names no card is one plain run around the player it names.
	segs = RenderEntry(ForgeSkipped{Player: 0}, stubNamer{})
	wantPlain := []LogSegment{
		{Text: "P0", Player: 0, HasPlayer: true},
		{Text: " skips their forge a key phase"},
	}
	if len(segs) != len(wantPlain) || segs[0] != wantPlain[0] || segs[1] != wantPlain[1] {
		t.Errorf("segments = %+v, want %+v", segs, wantPlain)
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
}

// houseNamer names a card after a house, so a keyword and a card name compete
// for the same span.
type houseNamer struct{ stubNamer }

func (houseNamer) Name(LocalID) string { return "Brobnar" }

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

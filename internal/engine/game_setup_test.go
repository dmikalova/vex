package engine

import (
	"errors"
	"testing"
)

// stockDeck fills both players' decks with plain Brobnar creatures so StartGame
// has cards to deal.
func stockDeck(g *Game) {
	for range 20 {
		g.AddToDeck(testCreature("Deck", 3), 0)
		g.AddToDeck(testCreature("Deck", 3), 1)
	}
}

// AddToDiscard and AddToArchives register a card and place it in the right pile.
func TestAddToDiscardAndArchives(t *testing.T) {
	g := NewGame("A", "B", 1)
	d := g.AddToDiscard(testCreature("Buried", 3), 0)
	if !g.State.Discard[0].contains(d) {
		t.Errorf("AddToDiscard did not place the card in the discard pile")
	}
	a := g.AddToArchives(testCreature("Stashed", 4), 1)
	if !g.State.Archives[1].contains(a) {
		t.Errorf("AddToArchives did not place the card in the archives")
	}
}

// StartGame deals the first player one more card than the second, then begins the
// first player's turn with the first-turn rule armed for them only.
func TestStartGameDealsAndStartsFirstTurn(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	stockDeck(g)
	g.StartGame(0)
	if got := int(g.State.Hand[0].Count); got != HandSize+FirstPlayerBonusCards {
		t.Errorf("first player hand = %d, want %d", got, HandSize+FirstPlayerBonusCards)
	}
	if got := int(g.State.Hand[1].Count); got != HandSize {
		t.Errorf("second player hand = %d, want %d", got, HandSize)
	}
	if g.State.ActivePlayer != 0 || g.State.Turn != 1 {
		t.Errorf("after StartGame: active=%d turn=%d, want active=0 turn=1",
			g.State.ActivePlayer, g.State.Turn)
	}
	if !g.State.FirstTurnPlayLimit[0] {
		t.Errorf("first-turn rule not armed for the first player")
	}
	if g.State.FirstTurnPlayLimit[1] {
		t.Errorf("first-turn rule wrongly armed for the second player")
	}
	if got := (GameStarted{FirstPlayer: 0}).Text(g); got != "Alice takes the first turn" {
		t.Errorf("GameStarted text = %q", got)
	}
}

// On the first player's first turn they may play or discard only one card; a
// second play or a discard is barred, while a card an effect plays bypasses it.
func TestFirstTurnRuleLimitsToOneCard(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	stockDeck(g)
	g.StartGame(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if _, err := g.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("first play should be allowed: %v", err)
	}
	if _, err := g.PlayCreature(0, 0, false); !errors.Is(err, ErrFirstTurnOneCard) {
		t.Errorf("second play err = %v, want %v", err, ErrFirstTurnOneCard)
	}
	if err := g.DiscardFromHand(0, 0); !errors.Is(err, ErrFirstTurnOneCard) {
		t.Errorf("discard after a play err = %v, want %v", err, ErrFirstTurnOneCard)
	}
	if err := g.CanPlay(0, g.State.Hand[0].IDs[0]); !errors.Is(err, ErrFirstTurnOneCard) {
		t.Errorf("CanPlay err = %v, want %v", err, ErrFirstTurnOneCard)
	}
	if _, err := g.PlayUpgrade(0, 0); !errors.Is(err, ErrFirstTurnOneCard) {
		t.Errorf("upgrade play after a play err = %v, want %v", err, ErrFirstTurnOneCard)
	}
	// A card an effect plays never passes through the volitional gates.
	before := g.State.Battleline[0].Count
	g.PlayFromHand(0, g.State.Hand[0].IDs[0])
	if g.State.Battleline[0].Count != before+1 {
		t.Errorf("effect-played card was wrongly barred by the first-turn rule")
	}
}

// The first-turn allowance is spent by a discard just as by a play: after one
// discard, a further play is barred.
func TestFirstTurnRuleSpentByDiscard(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	stockDeck(g)
	g.StartGame(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if err := g.DiscardFromHand(0, 0); err != nil {
		t.Fatalf("first discard should be allowed: %v", err)
	}
	if _, err := g.PlayCreature(0, 0, false); !errors.Is(err, ErrFirstTurnOneCard) {
		t.Errorf("play after a discard err = %v, want %v", err, ErrFirstTurnOneCard)
	}
	if err := g.CanDiscard(0, g.State.Hand[0].IDs[0]); !errors.Is(err, ErrFirstTurnOneCard) {
		t.Errorf("CanDiscard err = %v, want %v", err, ErrFirstTurnOneCard)
	}
}

// The rule lapses after the first turn: the second player's turn plays freely.
func TestFirstTurnRuleClearsAfterFirstTurn(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	stockDeck(g)
	g.StartGame(0)
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if g.State.FirstTurnPlayLimit[1] {
		t.Errorf("second player should never carry the first-turn rule")
	}
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if _, err := g.PlayCreature(1, 0, false); err != nil {
		t.Fatalf("first play: %v", err)
	}
	if _, err := g.PlayCreature(1, 0, false); err != nil {
		t.Errorf("second player second play should be allowed: %v", err)
	}
}

// The first-turn allowance binds only a play the first player makes of their own
// volition: after spending it, a second in-house play is barred, but an off-house
// play (one only a card effect could make legal) is not barred by the first-turn
// rule — it fails for the missing grant, not the one-card limit. This is what lets
// a grant like Captain Val Jericho or Subject Kirby play a further card on turn 1.
func TestFirstTurnRuleBindsOnlyVolitionalPlays(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	stockDeck(g)
	g.StartGame(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if _, err := g.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("first play should be allowed: %v", err)
	}

	// A second bare in-house play is barred by the first-turn rule, whatever its
	// entry point.
	inHouse := g.AddToHand(NewCard("bro", Brobnar, Creature, Common, WithPower(1)), 0)
	inHouseUp := g.AddToHand(NewCard("bro-up", Brobnar, Upgrade, Common), 0)
	if err := g.CanPlay(0, inHouse); !errors.Is(err, ErrFirstTurnOneCard) {
		t.Errorf("CanPlay in-house = %v, want %v", err, ErrFirstTurnOneCard)
	}
	if _, err := g.PlayCreature(
		0,
		handIdxByID(g, 0, inHouse),
		false,
	); !errors.Is(
		err,
		ErrFirstTurnOneCard,
	) {
		t.Errorf("PlayCreature in-house = %v, want %v", err, ErrFirstTurnOneCard)
	}
	if _, err := g.PlayUpgrade(
		0,
		handIdxByID(g, 0, inHouseUp),
	); !errors.Is(
		err,
		ErrFirstTurnOneCard,
	) {
		t.Errorf("PlayUpgrade in-house = %v, want %v", err, ErrFirstTurnOneCard)
	}

	// An off-house play is past the first-turn rule: it can only ever be legal via
	// a card effect, so it reports the missing grant (ErrWrongHouse), not the
	// first-turn limit.
	off := g.AddToHand(NewCard("off", Sanctum, Creature, Common, WithPower(1)), 0)
	offUp := g.AddToHand(NewCard("off-up", Sanctum, Upgrade, Common), 0)
	if err := g.CanPlay(0, off); !errors.Is(err, ErrWrongHouse) {
		t.Errorf("CanPlay off-house = %v, want %v", err, ErrWrongHouse)
	}
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, off), false); !errors.Is(err, ErrWrongHouse) {
		t.Errorf("PlayCreature off-house = %v, want %v", err, ErrWrongHouse)
	}
	if _, err := g.PlayUpgrade(0, handIdxByID(g, 0, offUp)); !errors.Is(err, ErrWrongHouse) {
		t.Errorf("PlayUpgrade off-house = %v, want %v", err, ErrWrongHouse)
	}
}

// A player who mulligans shuffles their hand back and draws one fewer card.
func TestMulliganDrawsOneFewer(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	stockDeck(g)
	g.SetChooser(0, optionPicker{idx: 1}) // mulligan
	g.SetChooser(1, FirstChooser{})       // keep
	g.StartGame(0)
	if got := int(g.State.Hand[0].Count); got != HandSize+FirstPlayerBonusCards-1 {
		t.Errorf("mulliganed first player hand = %d, want %d",
			got, HandSize+FirstPlayerBonusCards-1)
	}
	if got := int(g.State.Hand[1].Count); got != HandSize {
		t.Errorf("kept second player hand = %d, want %d", got, HandSize)
	}
	if got := (Mulliganed{
		Player: 0,
		Hand:   6,
	}).Text(
		g,
	); got != "Alice mulligans, drawing a new hand of 6" {
		t.Errorf("Mulliganed text = %q", got)
	}
}

// Chains cut the opening draw and shed one chain, and a mulligan sheds no further
// chain — the player is down the same single chain whether or not they mulligan.
func TestOpeningDealAndMulliganWithChains(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	stockDeck(g)
	g.State.Chains[0] = 3 // reduces the draw by one and sheds one on the opening draw
	g.SetChooser(0, optionPicker{idx: 1})
	g.SetChooser(1, FirstChooser{})
	g.StartGame(0)
	// Opening: 7 - (3+5)/6 = 6 dealt; then the mulligan redraws 6-1 = 5, no reduction.
	if got := int(g.State.Hand[0].Count); got != 5 {
		t.Errorf("chained mulligan hand = %d, want 5", got)
	}
	if g.State.Chains[0] != 2 {
		t.Errorf("chains = %d, want 2 (one shed on the opening draw, none on mulligan)",
			g.State.Chains[0])
	}
}

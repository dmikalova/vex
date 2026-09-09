package web

import (
	"slices"
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests cover the bridge between the engine's synchronous chooser and the
// browser's single UI goroutine: an effect that needs an answer parks its
// goroutine, the client renders the question, and a click sends the answer back.
//
// A prompt is raised the way the engine raises one — from a goroutine other than
// the one answering it — so what is being tested is the handoff, not a rehearsal
// of it.

// ask raises a card prompt from a background goroutine, as an effect mid-
// resolution does, and hands back the channel its answer arrives on.
func (c *client) ask(prompt string, declinable bool, cands []engine.LocalID) chan chooseReply {
	c.t.Helper()
	out := make(chan chooseReply, 1)
	go func() {
		var id engine.LocalID
		var ok bool
		if declinable {
			id, ok = c.g.chooser.ChooseCardOrDecline("A Card", prompt, cands)
		} else {
			id, ok = c.g.chooser.ChooseCreature("A Card", prompt, cands)
		}
		out <- chooseReply{id: id, ok: ok}
	}()
	return out
}

func TestAnsweringACardPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	first := c.deal(testCreature)
	second := c.deal(testCreature)
	c.playFromHand(first)
	c.playFromHand(second)
	cands := c.board()

	answer := c.ask("Choose a creature", false, cands)
	c.await("the prompt to go up", func() bool { return c.g.choosing })

	if c.g.chooserPrompt != "Choose a creature" {
		t.Errorf("the prompt reads %q, want %q", c.g.chooserPrompt, "Choose a creature")
	}
	if c.g.promptSource != "A Card" {
		t.Errorf("the prompt is attributed to %q, want %q", c.g.promptSource, "A Card")
	}
	if c.g.chooserDeclinable {
		t.Error("a mandatory prompt offered a way out")
	}

	c.g.chooseCandidate(c.ctx, cands[1])
	got := <-answer
	if !got.ok || got.id != cands[1] {
		t.Errorf("the effect was answered %v, want card %d", got, cands[1])
	}
	c.await("the prompt to come down", func() bool { return !c.g.choosing })
	if c.g.chooserPrompt != "" || c.g.chooserCandidates != nil || c.g.promptSource != "" {
		t.Error("the answered prompt left its question on screen")
	}
}

// An optional prompt can be passed on, which is what its Done button and Escape
// both mean.
func TestDecliningACardPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	answer := c.ask("You may destroy a creature", true, c.board())
	c.await("the prompt to go up", func() bool { return c.g.choosing })
	if !c.g.chooserDeclinable {
		t.Error("an optional prompt did not offer a way out")
	}

	c.press("n")
	got := <-answer
	if got.ok {
		t.Errorf("the declined prompt answered %v, want a pass", got)
	}
}

// A mandatory prompt has no way out, so n and Escape leave it up rather than
// letting the player walk away from an answer the effect needs.
func TestAMandatoryPromptCannotBeWalkedAwayFrom(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	cands := c.board()
	// Out of manual mode: a real match's prompts are the engine's to insist on.
	c.do(c.g.toggleManual)

	answer := c.ask("Choose a creature", false, cands)
	c.await("the prompt to go up", func() bool { return c.g.choosing })

	c.press("n")
	c.press("Escape")
	if !c.g.choosing {
		t.Fatal("a mandatory prompt was dismissed")
	}

	c.g.chooseCandidate(c.ctx, cands[0])
	if got := <-answer; !got.ok {
		t.Error("the prompt was not answered by the click that followed")
	}
}

// Manual mode is the exception: a playtester who has arranged a board a prompt
// cannot be answered on needs a way out of it.
func TestManualModeCanCancelAMandatoryPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	answer := c.ask("Choose a creature", false, c.board())
	c.await("the prompt to go up", func() bool { return c.g.choosing })
	c.press("Escape")
	if got := <-answer; got.ok {
		t.Errorf("the cancelled prompt answered %v, want no choice", got)
	}
}

// A prompt with no candidates has nothing to ask, so it answers itself rather
// than parking the effect on a question with no buttons.
func TestAPromptWithNoCandidatesAnswersItself(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	if _, ok := c.g.chooser.ChooseCreature("A Card", "Choose a creature", nil); ok {
		t.Error("a prompt with no candidates reported a choice")
	}
	if c.g.choosing {
		t.Error("a prompt with no candidates went up on screen")
	}
}

// order raises an ordering prompt from a background goroutine, as the engine's
// orderByChoice does through the Orderer, and hands back the channel the arranged
// order arrives on.
func (c *client) order(prompt string, ids []engine.LocalID) chan []engine.LocalID {
	c.t.Helper()
	out := make(chan []engine.LocalID, 1)
	go func() { out <- c.g.chooser.OrderCreatures("A Card", prompt, ids) }()
	return out
}

// samePermutation reports whether got is a rearrangement of want.
func samePermutation(got, want []engine.LocalID) bool {
	if len(got) != len(want) {
		return false
	}
	seen := map[engine.LocalID]int{}
	for _, id := range want {
		seen[id]++
	}
	for _, id := range got {
		seen[id]--
	}
	for _, n := range seen {
		if n != 0 {
			return false
		}
	}
	return true
}

// TestAutoResolveOrder checks that the Auto-resolve button answers an ordering
// prompt with a full random order in one click, ending the window without picking
// each ability in turn.
func TestAutoResolveOrder(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	first := c.deal(testCreature)
	second := c.deal(testCreature)
	c.playFromHand(first)
	c.playFromHand(second)
	cands := c.board()

	out := c.order("Order them", cands)
	c.await(
		"the ordering prompt to go up",
		func() bool { return c.g.choosing && c.g.chooserOrdering },
	)

	c.g.autoResolveOrder(c.ctx, nullEvent())
	got := <-out
	c.await("the prompt to come down", func() bool { return !c.g.choosing })

	if !samePermutation(got, cands) {
		t.Errorf("auto-resolve returned %v, want a permutation of %v", got, cands)
	}
}

// TestOrderByPickingEach checks that ordering a window by clicking cards still
// works through the Orderer: each pick resolves next and the last is forced.
func TestOrderByPickingEach(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	first := c.deal(testCreature)
	second := c.deal(testCreature)
	c.playFromHand(first)
	c.playFromHand(second)
	cands := c.board()

	out := c.order("Order them", cands)
	c.await("the ordering prompt to go up", func() bool { return c.g.chooserOrdering })

	c.g.chooseCandidate(c.ctx, cands[1])
	got := <-out
	c.await("the prompt to come down", func() bool { return !c.g.choosing })

	if want := []engine.LocalID{cands[1], cands[0]}; !slices.Equal(got, want) {
		t.Errorf("ordering picked %v, want %v", got, want)
	}
}

// TestAutoResolveOnlyForOrdering checks that Auto-resolve is inert unless an
// ordering prompt is up, so a stray click cannot answer another kind of prompt.
func TestAutoResolveOnlyForOrdering(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	only := c.deal(testCreature)
	next := c.deal(testCreature)
	c.playFromHand(only)
	c.playFromHand(next)
	cands := c.board()

	answer := c.ask("Choose a creature", false, cands)
	c.await("the prompt to go up", func() bool { return c.g.choosing })

	c.g.autoResolveOrder(c.ctx, nullEvent()) // no ordering prompt: ignored
	if c.g.chooserOrdering {
		t.Error("a plain card prompt was marked as ordering")
	}
	c.g.chooseCandidate(c.ctx, cands[0])
	if got := <-answer; !got.ok {
		t.Error("the plain prompt should still answer with a normal pick")
	}
}

func TestAnsweringAnOptionPrompt(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	out := make(chan int, 1)
	go func() { out <- c.g.chooser.ChooseOption("A Card", "Take them?", []string{"Yes", "No"}) }()
	c.await("the prompt to go up", func() bool { return c.g.choosingOption })

	if c.g.optionPrompt != "Take them?" {
		t.Errorf("the prompt reads %q, want %q", c.g.optionPrompt, "Take them?")
	}
	if c.g.promptButtons() != 2 {
		t.Errorf("the prompt offers %d buttons, want 2", c.g.promptButtons())
	}

	c.press("n")
	if got := <-out; got != 1 {
		t.Errorf("n answered with option %d, want No at 1", got)
	}
	c.await("the prompt to come down", func() bool { return !c.g.choosingOption })
	if c.g.optionLabels != nil || c.g.optionPrompt != "" {
		t.Error("the answered prompt left its question on screen")
	}
}

// Space affirms a yes/no prompt, which is the one option list that has an
// affirmative answer at all.
func TestSpaceAnswersYes(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	out := make(chan int, 1)
	go func() { out <- c.g.chooser.ChooseOption("A Card", "Take them?", []string{"Yes", "No"}) }()
	c.await("the prompt to go up", func() bool { return c.g.choosingOption })

	c.press(" ")
	if got := <-out; got != 0 {
		t.Errorf("Space answered with option %d, want Yes at 0", got)
	}
}

// A list of alternatives has no option that "yes" could mean, so Space leaves it
// alone rather than picking the first thing on the list.
func TestSpaceLeavesAListOfAlternativesAlone(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	out := make(chan int, 1)
	labels := []string{"Red", "Blue", "Yellow"}
	go func() { out <- c.g.chooser.ChooseOption("A Card", "Forge which key?", labels) }()
	c.await("the prompt to go up", func() bool { return c.g.choosingOption })

	c.press(" ")
	c.press("n")
	if !c.g.choosingOption {
		t.Fatal("a list of alternatives was answered by a key that means neither")
	}

	c.press("b")
	if got := <-out; got != 1 {
		t.Errorf("b answered with option %d, want Blue at 1", got)
	}
}

// A prompt over a pile the board does not draw opens the zone viewer, which
// makes the pile clickable — and closes it again once the prompt is answered,
// because the viewer belonged to the prompt.
func TestAPromptOverAPileOpensTheViewer(t *testing.T) {
	c := newClient(t)
	c.manual()
	me := c.g.active()
	id := c.deal(testCreature)
	c.g.g.ManualMove(id, engine.ManualDiscard)

	answer := c.ask("Choose a card in your discard pile", false, []engine.LocalID{id})
	c.await("the prompt to go up", func() bool { return c.g.choosing })

	if c.g.zonesPlayer != me || c.g.promptZone != "Discard" {
		t.Fatalf("the viewer opened at player %d zone %q, want player %d Discard",
			c.g.zonesPlayer, c.g.promptZone, me)
	}
	// A viewer the prompt opened is the only place its candidates are clickable,
	// so a click on the backdrop does not take it away.
	c.do(c.g.closeZones)
	if c.g.zonesPlayer != me {
		t.Error("the prompt's viewer was closed out from under it")
	}

	c.g.chooseCandidate(c.ctx, id)
	<-answer
	c.await("the viewer to close", func() bool { return c.g.zonesPlayer == -1 })
	if c.g.promptZone != "" {
		t.Errorf("the prompt's zone is still %q", c.g.promptZone)
	}
}

// A bounded pick from the top of the deck (a "look at the top N cards" reveal —
// Navigator Ali, Lay of the Land) is offered as a short list of action-bar buttons
// rather than opening the zone viewer, and a reorder prompt notes that the last
// pick ends up on top.
func TestABoundedDeckPromptOffersButtons(t *testing.T) {
	c := newClient(t)
	c.manual()
	id := c.deal(testCreature)
	c.g.g.ManualMove(id, engine.ManualDeckTop)

	answer := c.ask("Choose the next card to place on top of your deck",
		false, []engine.LocalID{id})
	c.await("the prompt to go up", func() bool { return c.g.choosing })

	if !c.g.promptAsButtons {
		t.Fatal("a bounded deck prompt did not switch to action-bar buttons")
	}
	if c.g.promptZone != "" {
		t.Errorf("a bounded deck prompt opened the viewer at %q instead of buttons",
			c.g.promptZone)
	}
	html := app.HTMLString(c.g.Render())
	if !strings.Contains(html, "prompt-pick") {
		t.Error("the prompt rendered no candidate buttons")
	}
	if !strings.Contains(html, "ends up on top") {
		t.Error("the reorder prompt is missing the last-is-on-top note")
	}
	c.g.chooseCandidate(c.ctx, id)
	if got := <-answer; !got.ok || got.id != id {
		t.Errorf("the button answered with %+v, want id %d ok", got, id)
	}
}

// A declinable pile prompt (Not Finished with You — shuffle any number, including
// zero) keeps the zone viewer, but the viewer can be finished from a Done
// affordance, and closing it submits the current selection by declining.
func TestADeclinablePilePromptIsFinishedFromTheViewer(t *testing.T) {
	c := newClient(t)
	c.manual()
	me := c.g.active()
	id := c.deal(testCreature)
	c.g.g.ManualMove(id, engine.ManualDiscard)

	answer := c.ask("Choose a creature to shuffle into your deck",
		true, []engine.LocalID{id})
	c.await("the prompt to go up", func() bool { return c.g.choosing })

	if c.g.zonesPlayer != me || c.g.promptZone != "Discard" {
		t.Fatalf("a declinable pile prompt did not open the discard viewer")
	}
	if !c.g.chooserDeclinable {
		t.Fatal("Not Finished with You's prompt is not marked declinable")
	}
	html := app.HTMLString(c.g.Render())
	if !strings.Contains(html, "zones-done") {
		t.Error("the declinable viewer has no Done affordance")
	}
	// Closing the viewer submits the current (empty) selection by declining.
	c.do(c.g.closeZones)
	if got := <-answer; got.ok {
		t.Errorf("closing the declinable viewer answered %+v, want a decline", got)
	}
	c.await("the viewer to close", func() bool { return c.g.zonesPlayer == -1 })
}

// A prompt over cards on the board leaves the viewer alone: the board already
// draws them.
func TestAPromptOverTheBoardDoesNotOpenTheViewer(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	answer := c.ask("Choose a creature", false, c.board())
	c.await("the prompt to go up", func() bool { return c.g.choosing })
	if c.g.zonesPlayer != -1 {
		t.Errorf("the viewer opened at %d for a prompt over the board", c.g.zonesPlayer)
	}
	c.g.chooseCandidate(c.ctx, c.board()[0])
	<-answer
}

// A viewer the player opened themselves is theirs to close.
func TestClosingAViewerThePlayerOpened(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.zonesPlayer = 0
	c.do(c.g.closeZones)
	if c.g.zonesPlayer != -1 {
		t.Error("the player's own viewer would not close")
	}
}

// Answering a prompt that is not up is a click that arrived too late — from a
// double click on the prompt before it — and is dropped rather than answering
// the next one.
func TestAnsweringWhenNoPromptIsUp(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	c.g.chooseCandidate(c.ctx, c.board()[0])
	c.do(c.g.declineChooser)
	c.do(c.g.cancelChooser)
	c.do(c.g.chooseOptionIdx(0))

	answer := c.ask("Choose a creature", false, c.board())
	c.await("the prompt to go up", func() bool { return c.g.choosing })
	select {
	case got := <-answer:
		t.Fatalf("the prompt was answered %v by a click that came before it", got)
	default:
	}
	c.g.chooseCandidate(c.ctx, c.board()[0])
	<-answer
}

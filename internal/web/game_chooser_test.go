package web

import (
	"slices"
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
)

// These tests cover the prompt seam: an effect that needs an answer suspends the
// session on a Request, the client renders the question off that Request, and a
// click applies the Command that answers it.
//
// A prompt is raised the way the engine raises one — from a driving action
// suspended mid-resolution (c.script) — so what is being tested is the handoff,
// not a rehearsal of it.

// pick is the answer a card prompt came back with.
type pick struct {
	id engine.LocalID
	ok bool
}

// ask raises a card prompt from a scripted action, as an effect mid-resolution
// does, and hands back the answer it eventually returns.
func (c *client) ask(
	source engine.LocalID,
	prompt string,
	declinable bool,
	cands []engine.LocalID,
) *pick {
	c.t.Helper()
	got := &pick{}
	player := c.g.active()
	c.script(func(eg *engine.Game) {
		if declinable {
			got.id, got.ok = eg.ChooseCardOptional(player, source, prompt, cands)
			return
		}
		got.id, got.ok = eg.ChooseCreature(player, source, prompt, cands)
	})
	return got
}

func TestAnsweringACardPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	first := c.deal(testCreature)
	second := c.deal(testCreature)
	c.playFromHand(first)
	c.playFromHand(second)
	cands := c.board()

	answer := c.ask(cands[0], "Choose a creature", false, cands)
	c.await("the prompt to go up", c.g.choosing)

	if c.g.chooserPrompt() != "Choose a creature" {
		t.Errorf("the prompt reads %q, want %q", c.g.chooserPrompt(), "Choose a creature")
	}
	wantSource := engine.PromptSource{Card: cands[0], HasCard: true}
	if c.g.promptSource() != wantSource {
		t.Errorf("the prompt is attributed to %+v, want %+v", c.g.promptSource(), wantSource)
	}
	if c.g.chooserDeclinable() {
		t.Error("a mandatory prompt offered a way out")
	}

	c.g.chooseCandidate(c.ctx, cands[1])
	c.settle()
	if !answer.ok || answer.id != cands[1] {
		t.Errorf("the effect was answered %v, want card %d", answer, cands[1])
	}
	c.await("the prompt to come down", func() bool { return !c.g.choosing() })
	if c.g.chooserPrompt() != "" || c.g.chooserCandidates() != nil ||
		c.g.promptSource() != (engine.PromptSource{}) {
		t.Error("the answered prompt left its question on screen")
	}
}

// A prompt taking over the board clears any card the player had selected, so no
// stale selection ring lingers on a non-candidate card behind the prompt — the
// artifact whose end-of-turn ability (Fangtooth Cavern) raised the prompt should
// not read as still selected while the board dims around the candidates.
func TestOpeningAPromptClearsTheSelection(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	first := c.deal(testCreature)
	second := c.deal(testCreature)
	c.playFromHand(first)
	c.playFromHand(second)
	cands := c.board()

	c.g.hasSel, c.g.sel = true, cands[0]

	c.ask(cands[0], "Choose a creature", false, cands)
	c.await("the prompt to go up", c.g.choosing)

	if c.g.hasSel {
		t.Error("opening a prompt left a stale selection on the board")
	}

	c.g.chooseCandidate(c.ctx, cands[1])
	c.await("the prompt to come down", func() bool { return !c.g.choosing() })
}

// An optional prompt can be passed on, which is what its Done button and Escape
// both mean.
func TestDecliningACardPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	board := c.board()

	answer := c.ask(board[0], "You may destroy a creature", true, board)
	c.await("the prompt to go up", c.g.choosing)
	if !c.g.chooserDeclinable() {
		t.Error("an optional prompt did not offer a way out")
	}

	c.press("n")
	if answer.ok {
		t.Errorf("the declined prompt answered %v, want a pass", answer)
	}
}

// A mandatory prompt has no way out, so n and Escape leave it up rather than
// letting the player walk away from an answer the effect needs.
func TestAMandatoryPromptCannotBeWalkedAwayFrom(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	cands := c.board()
	// Out of manual mode: a real match's prompts are the engine's to insist on.
	c.do(c.g.toggleManual)

	answer := c.ask(cands[0], "Choose a creature", false, cands)
	c.await("the prompt to go up", c.g.choosing)

	c.press("n")
	c.press("Escape")
	if !c.g.choosing() {
		t.Fatal("a mandatory prompt was dismissed")
	}

	c.g.chooseCandidate(c.ctx, cands[0])
	c.settle()
	if !answer.ok {
		t.Error("the prompt was not answered by the click that followed")
	}
}

// Manual mode is the exception: a playtester who has arranged a board a prompt
// cannot be answered on needs a way out of it. Cancelling is undo — the whole
// root action that raised the prompt is rewound — so it is driven here through a
// real use of a card, which is the action there is to back out of.
func TestEscapeCancelsAPromptInManualMode(t *testing.T) {
	c := newClient(t)
	c.manualTurn(promptHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	cannon := c.stagePromptArtifact()

	c.g.selectBoardID(c.ctx, cannon)
	c.do(c.g.useAction)
	c.await("the damage prompt", c.g.choosing)

	c.press("Escape")
	if c.g.choosing() {
		t.Fatal("Escape did not back out of the prompt")
	}
	if c.g.eng().State.Cards[cannon].Exhausted {
		t.Error("the cancelled action left the artifact spent")
	}
}

// A prompt with no candidates has nothing to ask, so it answers itself rather
// than suspending the action on a question with no buttons.
func TestAPromptWithNoCandidatesAnswersItself(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	answer := c.ask(c.hand()[0], "Choose a creature", false, nil)
	if answer.ok {
		t.Error("a prompt with no candidates reported a choice")
	}
	if c.g.choosing() {
		t.Error("a prompt with no candidates went up on screen")
	}
}

// noSource is a LocalID no card in these tests is ever dealt, so Game.promptSource
// resolves it to PromptSource{} — an unattributed prompt — the same as the archive
// discard's own call to orderByChoice, which names no source card.
const noSource = engine.LocalID(255)

// order raises an ordering window from a scripted action, as the engine's
// orderByChoice does, and hands back the arranged order.
func (c *client) order(prompt string, ids []engine.LocalID) *[]engine.LocalID {
	c.t.Helper()
	got := new([]engine.LocalID)
	player := c.g.active()
	c.script(func(eg *engine.Game) { *got = eg.OrderByChoice(player, noSource, prompt, ids) })
	return got
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
// window with a full random order in one click, ending the window without picking
// each card in turn.
func TestAutoResolveOrder(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	cands := c.board()

	got := c.order("Order them", cands)
	c.await(
		"the ordering prompt to go up",
		func() bool { return c.g.choosing() && c.g.chooserOrdering() },
	)

	c.g.autoResolveOrder(c.ctx, nullEvent())
	c.await("the prompt to come down", func() bool { return !c.g.choosing() })

	if !samePermutation(*got, cands) {
		t.Errorf("auto-resolve returned %v, want a permutation of %v", *got, cands)
	}
}

// TestOrderByPickingEach checks that ordering a window by clicking cards still
// works: each pick resolves next and the last is forced.
func TestOrderByPickingEach(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	cands := c.board()

	got := c.order("Order them", cands)
	c.await("the ordering prompt to go up", c.g.chooserOrdering)

	c.g.chooseCandidate(c.ctx, cands[1])
	c.await("the prompt to come down", func() bool { return !c.g.choosing() })

	if want := []engine.LocalID{cands[1], cands[0]}; !slices.Equal(*got, want) {
		t.Errorf("ordering picked %v, want %v", *got, want)
	}
}

// TestAutoResolveOnlyForOrdering checks that Auto-resolve is inert unless an
// ordering window is up, so a stray click cannot answer another kind of prompt.
// An ordering step is the one mandatory card pick with no source card behind it,
// which is what tells the two apart.
func TestAutoResolveOnlyForOrdering(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	cands := c.board()

	answer := c.ask(cands[0], "Choose a creature", false, cands)
	c.await("the prompt to go up", c.g.choosing)

	c.g.autoResolveOrder(c.ctx, nullEvent()) // no ordering window: ignored
	if c.g.chooserOrdering() {
		t.Error("a plain card prompt was marked as ordering")
	}
	if c.scriptDone() {
		t.Fatal("Auto-resolve answered a plain card prompt")
	}
	c.g.chooseCandidate(c.ctx, cands[0])
	c.settle()
	if !answer.ok {
		t.Error("the plain prompt should still answer with a normal pick")
	}
}

// option raises a labeled option prompt from a scripted action and hands back the
// index it is answered with.
func (c *client) option(source engine.LocalID, prompt string, labels []string) *int {
	c.t.Helper()
	got := new(int)
	*got = -1
	player := c.g.active()
	c.script(func(eg *engine.Game) { *got = eg.ChooseOption(player, source, prompt, labels) })
	return got
}

func TestAnsweringAnOptionPrompt(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	src := c.hand()[0]
	got := c.option(src, "Take them?", []string{"Yes", "No"})
	c.await("the prompt to go up", c.g.choosingOption)

	if c.g.optionPrompt() != "Take them?" {
		t.Errorf("the prompt reads %q, want %q", c.g.optionPrompt(), "Take them?")
	}
	if c.g.promptButtons() != 2 {
		t.Errorf("the prompt offers %d buttons, want 2", c.g.promptButtons())
	}

	c.press("n")
	if *got != 1 {
		t.Errorf("n answered with option %d, want No at 1", *got)
	}
	c.await("the prompt to come down", func() bool { return !c.g.choosingOption() })
	if c.g.optionLabels() != nil || c.g.optionPrompt() != "" {
		t.Error("the answered prompt left its question on screen")
	}
}

// Space affirms a yes/no prompt, which is the one option list that has an
// affirmative answer at all.
func TestSpaceAnswersYes(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	got := c.option(c.hand()[0], "Take them?", []string{"Yes", "No"})
	c.await("the prompt to go up", c.g.choosingOption)

	c.press(" ")
	if *got != 0 {
		t.Errorf("Space answered with option %d, want Yes at 0", *got)
	}
}

// A list of alternatives has no option that "yes" could mean, so Space leaves it
// alone rather than picking the first thing on the list.
func TestSpaceLeavesAListOfAlternativesAlone(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	got := c.option(c.hand()[0], "Forge which key?", []string{"Red", "Blue", "Yellow"})
	c.await("the prompt to go up", c.g.choosingOption)

	c.press(" ")
	c.press("n")
	if !c.g.choosingOption() {
		t.Fatal("a list of alternatives was answered by a key that means neither")
	}

	c.press("b")
	if *got != 1 {
		t.Errorf("b answered with option %d, want Blue at 1", *got)
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
	other := c.deal(testCreature)
	c.g.eng().ManualMove(id, engine.ManualDiscard)
	c.g.eng().ManualMove(other, engine.ManualDiscard)

	answer := c.ask(id, "Choose a card in your discard pile", false,
		[]engine.LocalID{id, other})
	c.await("the prompt to go up", c.g.choosing)

	if c.g.zonesPlayer != me || c.g.promptZone != "Discard" {
		t.Fatalf("the viewer opened at player %d zone %q, want player %d Discard",
			c.g.zonesPlayer, c.g.promptZone, me)
	}

	c.g.chooseCandidate(c.ctx, id)
	c.settle()
	if !answer.ok {
		t.Error("the pile prompt was not answered")
	}
	c.await("the viewer to close", func() bool { return c.g.zonesPlayer == -1 })
	if c.g.promptZone != "" {
		t.Errorf("the prompt's zone is still %q", c.g.promptZone)
	}
}

// A bounded pick from the top of the deck (a "look at the top N cards" reveal —
// Navigator Ali, Lay of the Land) is offered as a short list of action-bar buttons
// rather than opening the zone viewer, and a reorder prompt notes that the first
// pick ends up on the bottom.
func TestABoundedDeckPromptOffersButtons(t *testing.T) {
	c := newClient(t)
	c.manual()
	id := c.deal(testCreature)
	other := c.deal(testCreature)
	c.g.eng().ManualMove(id, engine.ManualDeckTop)
	c.g.eng().ManualMove(other, engine.ManualDeckTop)

	answer := c.ask(id, "Choose the next card to place on top of your deck",
		false, []engine.LocalID{id, other})
	c.await("the prompt to go up", c.g.choosing)

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
	if !strings.Contains(html, "ends up on the bottom") {
		t.Error("the reorder prompt is missing the first-is-on-the-bottom note")
	}
	c.g.chooseCandidate(c.ctx, id)
	c.settle()
	if !answer.ok || answer.id != id {
		t.Errorf("the button answered with %+v, want id %d ok", answer, id)
	}
}

// A declinable pile prompt (Not Finished with You — shuffle any number, including
// zero) keeps the zone viewer, and is finished from the viewer's Done affordance.
// Closing the viewer is not an answer: it only gets the modal out of the way.
func TestADeclinablePilePromptIsFinishedFromTheViewer(t *testing.T) {
	c := newClient(t)
	c.manual()
	me := c.g.active()
	id := c.deal(testCreature)
	c.g.eng().ManualMove(id, engine.ManualDiscard)

	answer := c.ask(id, "Choose a creature to shuffle into your deck",
		true, []engine.LocalID{id})
	c.await("the prompt to go up", c.g.choosing)

	if c.g.zonesPlayer != me || c.g.promptZone != "Discard" {
		t.Fatalf("a declinable pile prompt did not open the discard viewer")
	}
	if !c.g.chooserDeclinable() {
		t.Fatal("Not Finished with You's prompt is not marked declinable")
	}
	html := app.HTMLString(c.g.Render())
	if !strings.Contains(html, "zones-done") {
		t.Error("the declinable viewer has no Done affordance")
	}
	// Closing the viewer is not an answer — the prompt is still waiting.
	c.do(c.g.closeZones)
	if !c.g.choosing() {
		t.Error("closing the viewer answered the prompt; it should only get out of the way")
	}
	c.do(c.g.declineChooser)
	if answer.ok {
		t.Errorf("Done answered %+v, want a decline", answer)
	}
	c.await("the viewer to close", func() bool { return c.g.zonesPlayer == -1 })
}

// The zone viewer is always dismissible, even under a mandatory prompt whose only
// candidates are in the pile: the player may need to read the board underneath to
// decide. The prompt stays up, and reopening the viewer returns to the same row.
func TestAMandatoryPileViewerCanBeClosedWithoutAnswering(t *testing.T) {
	c := newClient(t)
	c.manual()
	me := c.g.active()
	id := c.deal(testCreature)
	other := c.deal(testCreature)
	c.g.eng().ManualMove(id, engine.ManualDiscard)
	c.g.eng().ManualMove(other, engine.ManualDiscard)
	// Out of manual mode: a real match's prompts are the engine's to insist on.
	c.do(c.g.toggleManual)

	answer := c.ask(id, "Choose a creature to put into play", false,
		[]engine.LocalID{id, other})
	c.await("the prompt to go up", c.g.choosing)

	c.do(c.g.closeZones)
	if c.g.zonesPlayer != -1 {
		t.Error("a mandatory pile viewer refused to close")
	}
	if !c.g.choosing() || c.g.promptZone != "Discard" {
		t.Error("closing the viewer dropped the prompt; it should still be waiting")
	}

	c.g.zonesPlayer = me
	c.g.chooseCandidate(c.ctx, id)
	c.settle()
	if !answer.ok || answer.id != id {
		t.Errorf("the reopened viewer answered %+v, want card %d", answer, id)
	}
}

// A prompt over cards on the board leaves the viewer alone: the board already
// draws them.
func TestAPromptOverTheBoardDoesNotOpenTheViewer(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))

	c.ask(c.board()[0], "Choose a creature", false, c.board())
	c.await("the prompt to go up", c.g.choosing)
	if c.g.zonesPlayer != -1 {
		t.Errorf("the viewer opened at %d for a prompt over the board", c.g.zonesPlayer)
	}
	c.g.chooseCandidate(c.ctx, c.board()[0])
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
	c.playFromHand(c.deal(testCreature))
	cands := c.board()

	c.g.chooseCandidate(c.ctx, cands[0])
	c.do(c.g.declineChooser)
	c.do(c.g.chooseOptionIdx(0))

	answer := c.ask(cands[0], "Choose a creature", false, cands)
	c.await("the prompt to go up", c.g.choosing)
	if c.scriptDone() {
		t.Fatal("the prompt was answered by a click that came before it")
	}
	c.g.chooseCandidate(c.ctx, cands[0])
	c.settle()
	if !answer.ok {
		t.Error("the prompt was not answered by the click that followed")
	}
}

// askReaction raises a trigger window from a scripted action, the way the
// engine's reaction ordering does, and hands back the index it is answered with.
func (c *client) askReaction(prompt string, reactions []engine.OrderableReaction) *int {
	c.t.Helper()
	got := new(int)
	*got = -1
	player := c.g.active()
	c.script(func(eg *engine.Game) { *got = eg.ChooseReaction(player, prompt, reactions) })
	return got
}

// A trigger window over cards on the board is answered by clicking the card whose
// ability resolves next, not by reading a menu of ability text.
func TestOrderingReactionsByClickingTheSourceCard(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	board := c.board()

	got := c.askReaction("Resolve destroyed abilities", []engine.OrderableReaction{
		{Card: board[0], HasCard: true, Label: "First: gain 1 Æmber"},
		{Card: board[1], HasCard: true, Label: "Second: draw a card"},
	})
	c.await("the reaction window to go up", c.g.choosing)

	if c.g.choosingOption() {
		t.Error("the window fell back to a button list; it should be clickable")
	}
	c.g.chooseCandidate(c.ctx, board[1])
	c.settle()
	if *got != 1 {
		t.Errorf("clicking the second source answered %d, want 1", *got)
	}
}

// One card carrying two pending abilities cannot be told apart by a click alone,
// so the click picks the card and a button list then picks the ability. It is
// never resolved top-down behind the player's back.
func TestOneCardWithTwoReactionsAsksWhichAbility(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	board := c.board()

	got := c.askReaction("Resolve destroyed abilities", []engine.OrderableReaction{
		{Card: board[0], HasCard: true, Label: "First: gain 1 Æmber"},
		{Card: board[1], HasCard: true, Label: "Second: draw a card"},
		{Card: board[1], HasCard: true, Label: "Second: deal 1 damage"},
	})
	c.await("the reaction window to go up", c.g.choosing)

	c.g.chooseCandidate(c.ctx, board[1])
	c.await("the follow-up ability prompt", c.g.choosingOption)
	if c.g.optionPrompt() != whichAbilityPrompt {
		t.Errorf("follow-up prompt = %q, want %q", c.g.optionPrompt(), whichAbilityPrompt)
	}

	c.do(c.g.chooseOptionIdx(1))
	if *got != 2 {
		t.Errorf("the second ability of the clicked card answered %d, want 2", *got)
	}
}

// A window carrying a reaction with no card behind it (a lasting duration
// reaction) has nothing to click, so the whole window stays on the labeled list
// rather than hiding the cardless entry.
func TestAReactionWithoutACardStaysOnTheLabeledList(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	board := c.board()

	got := c.askReaction("Choose which card's ability resolves next",
		[]engine.OrderableReaction{
			{Card: board[0], HasCard: true, Label: "First: gain 1 Æmber"},
			{Label: "gain 1 Æmber"},
		})
	c.await("the labeled list to go up", c.g.choosingOption)

	c.do(c.g.chooseOptionIdx(1))
	if *got != 1 {
		t.Errorf("the labeled list answered %d, want 1", *got)
	}
}

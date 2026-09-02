package cardtest

import (
	"strings"
	"time"

	"github.com/dmikalova/vactrol/internal/engine"
)

// actionTimeout is a backstop: normal misuse (a wrong click, or a prompt left
// unanswered at the end of a test) is caught immediately with a helpful message,
// so this only trips if the engine genuinely wedges — which would otherwise hang
// the whole test binary.
const actionTimeout = 5 * time.Second

// promptReq is one decision the engine is waiting on, carried from the action
// goroutine to the test goroutine. reply is how the test answers: a creature id,
// an option index, or -1 to decline.
type promptReq struct {
	player   int
	source   string // name of the card whose ability raised the prompt, or ""
	text     string
	isOption bool
	// declinable marks a card prompt the player may pass on, answered with
	// Player.ClickDone. A forced card prompt has no Done to click.
	declinable bool
	candidates []engine.LocalID
	options    []string
	reply      chan int
}

// actionResult is the outcome of a completed engine action.
type actionResult struct{ err error }

// bridgeChooser is the per-player Chooser installed on the game. When the engine
// asks it to make a choice, it forwards the prompt to the test goroutine and
// blocks until the test clicks a reply — reproducing the interactive
// play-then-click flow. A choice with a single candidate is resolved
// automatically, since there is nothing to decide.
type bridgeChooser struct {
	h      *Harness
	player int
}

// ChooseCreature forwards a creature choice to the test and waits for the click.
func (b bridgeChooser) ChooseCreature(
	source, prompt string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	if len(candidates) == 1 {
		return candidates[0], true
	}
	reply := make(chan int)
	b.h.prompt <- promptReq{
		player:     b.player,
		source:     source,
		text:       prompt,
		candidates: append([]engine.LocalID(nil), candidates...),
		reply:      reply,
	}
	r := <-reply
	if r < 0 {
		return 0, false
	}
	return engine.LocalID(r), true
}

// ChooseCardOrDecline forwards an optional card choice to the test and waits for
// the click. Unlike ChooseCreature a sole candidate is still offered, because
// declining it is a legal answer the test must be able to script.
func (b bridgeChooser) ChooseCardOrDecline(
	source, prompt string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	reply := make(chan int)
	b.h.prompt <- promptReq{
		player:     b.player,
		source:     source,
		text:       prompt,
		declinable: true,
		candidates: append([]engine.LocalID(nil), candidates...),
		reply:      reply,
	}
	r := <-reply
	if r < 0 {
		return 0, false
	}
	return engine.LocalID(r), true
}

// ChooseOption forwards a labeled option choice to the test and waits for it.
func (b bridgeChooser) ChooseOption(source, prompt string, options []string) int {
	if len(options) <= 1 {
		return 0
	}
	// Forging always asks which key colour to forge; the harness answers it
	// automatically (first colour) so it never interrupts a card test.
	if prompt == engine.KeyColorPrompt {
		return 0
	}
	reply := make(chan int)
	b.h.prompt <- promptReq{
		player:   b.player,
		source:   source,
		text:     prompt,
		isOption: true,
		options:  append([]string(nil), options...),
		reply:    reply,
	}
	r := <-reply
	if r < 0 {
		return 0
	}
	return r
}

// OrderCreatures arranges a multi-target resolution order. By default it keeps the
// engine's order (no prompt, so ordinary AoE effects do not interrupt a test); a
// test can take control for the next ordering with Player.Order.
func (b bridgeChooser) OrderCreatures(_, _ string, ids []engine.LocalID) []engine.LocalID {
	script := b.h.orderScript[b.player]
	b.h.orderScript[b.player] = nil
	if len(script) == 0 {
		return ids
	}
	return b.h.reorder(ids, script)
}

// run starts an engine action on a goroutine and advances to the first stop
// point (a prompt, or completion). The action and the test never touch the game
// concurrently: while a prompt is pending the goroutine is parked on its reply
// channel, and the test only reads the game between stop points.
func (h *Harness) run(label string, fn func() error) {
	h.t.Helper()
	h.label = label
	h.running = true
	go func() { h.done <- actionResult{err: fn()} }()
	h.advance()
}

// advance waits for the running action to raise a prompt or finish. A finished
// action that returned an error fails the test immediately, naming the action.
func (h *Harness) advance() {
	h.t.Helper()
	select {
	case req := <-h.prompt:
		h.current = &req
	case res := <-h.done:
		h.current = nil
		h.running = false
		if res.err != nil {
			h.t.Fatalf("%s: %v", h.label, res.err)
		}
	case <-time.After(actionTimeout):
		h.t.Fatalf("%s: the engine did not finish or ask for a choice (harness bug)", h.label)
	}
}

// requirePrompt returns the pending prompt, failing if none is pending or it is
// addressed to the other player.
func (h *Harness) requirePrompt(player int) promptReq {
	h.t.Helper()
	if h.current == nil {
		h.t.Fatalf("no prompt is pending for %s", playerName(player))
	}
	if h.current.player != player {
		h.t.Fatalf("the pending prompt %q is for %s, not %s",
			h.current.text, playerName(h.current.player), playerName(player))
	}
	return *h.current
}

// click answers a pending creature prompt for player. A target that is not among
// the live candidates fails immediately, listing what is clickable — so a
// mis-scripted click never needs a timeout to surface.
func (h *Harness) click(player int, target any) {
	h.t.Helper()
	req := h.requirePrompt(player)
	if req.isOption {
		h.t.Fatalf(
			"ClickCard: the pending prompt %q for %s expects an option — use ClickOption(%v)",
			req.text,
			playerName(player),
			req.options,
		)
	}
	id, ok := h.matchCandidate(target, req.candidates)
	if !ok {
		h.t.Fatalf("ClickCard(%s): not a legal target for %q\nclickable: %s",
			h.describe(target), req.text, h.nameList(req.candidates))
	}
	reply := req.reply
	h.current = nil
	reply <- int(id)
	h.advance()
}

// clickOption answers a pending option prompt for player by label. It matches an
// option exactly, or by a unique case-insensitive substring — so a long "choose
// one" effect line can be selected by a short, distinctive fragment.
func (h *Harness) clickOption(player int, label string) {
	h.t.Helper()
	req := h.requirePrompt(player)
	if !req.isOption {
		h.t.Fatalf("ClickOption: the pending prompt %q for %s expects a card — use ClickCard",
			req.text, playerName(player))
	}
	for i, opt := range req.options {
		if opt == label {
			h.answerOption(req, i)
			return
		}
	}
	match, n := -1, 0
	for i, opt := range req.options {
		if strings.Contains(strings.ToLower(opt), strings.ToLower(label)) {
			match, n = i, n+1
		}
	}
	if n == 1 {
		h.answerOption(req, match)
		return
	}
	if n > 1 {
		h.t.Fatalf(
			"ClickOption(%q): matches %d options for %q\noptions: %v",
			label,
			n,
			req.text,
			req.options,
		)
	}
	h.t.Fatalf("ClickOption(%q): not an option for %q\noptions: %v", label, req.text, req.options)
}

// answerOption sends an option index reply and advances to the next stop point.
func (h *Harness) answerOption(req promptReq, i int) {
	h.current = nil
	req.reply <- i
	h.advance()
}

// clickDone declines a pending declinable card prompt — the Done button next to
// the highlighted cards. A forced prompt has no Done, so declining one fails.
func (h *Harness) clickDone(player int) {
	h.t.Helper()
	req := h.requirePrompt(player)
	if !req.declinable {
		h.t.Fatalf("ClickDone: the pending prompt %q for %s cannot be declined",
			req.text, playerName(player))
	}
	h.current = nil
	req.reply <- -1
	h.advance()
}

// matchCandidate resolves a def-or-handle target to one of the candidate ids.
func (h *Harness) matchCandidate(target any, candidates []engine.LocalID) (engine.LocalID, bool) {
	h.t.Helper()
	switch v := target.(type) {
	case Card:
		v.require()
		if containsID(candidates, v.id) {
			return v.id, true
		}
		return 0, false
	case engine.CardDefinition:
		found, n := engine.LocalID(0), 0
		for _, id := range candidates {
			if h.g.Name(id) == v.Name {
				found, n = id, n+1
			}
		}
		if n == 1 {
			return found, true
		}
		if n > 1 {
			h.t.Fatalf(
				"ClickCard(%s): %d matching targets — use ct.Bind to name the copy you mean",
				v.Name,
				n,
			)
		}
		return 0, false
	default:
		h.t.Fatalf("ClickCard: cannot click a %T", target)
		return 0, false
	}
}

// checkReady runs as a t.Cleanup for every scenario: a test that ends with a
// prompt still pending failed to answer a choice, which is reported here with the
// pending prompt and what was clickable — the out-of-the-box "ready to take
// action" guard. It then drains the parked action goroutine so it does not leak.
func (h *Harness) checkReady() {
	if h.current == nil {
		return
	}
	req := *h.current
	if req.isOption {
		h.t.Errorf("test ended with an unanswered prompt %q for %s\noptions: %v",
			req.text, playerName(req.player), req.options)
	} else {
		h.t.Errorf("test ended with an unanswered prompt %q for %s\nclickable: %s",
			req.text, playerName(req.player), h.nameList(req.candidates))
	}
	for h.current != nil {
		reply := h.current.reply
		h.current = nil
		reply <- -1 // decline, letting the effect abort and the goroutine finish
		select {
		case req := <-h.prompt:
			h.current = &req
		case <-h.done:
			h.current = nil
		case <-time.After(actionTimeout):
			return
		}
	}
}

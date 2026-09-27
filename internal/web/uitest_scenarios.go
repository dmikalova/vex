package web

import "fmt"

// This file is the browser-scenario registry: the one definition of every
// scenario, which both the page a human drives (uitest.go) and the driver that
// reads its status element run. See uitest.go for the host that runs them, and
// internal/web/AGENTS.md for what belongs here.
//
// A scenario is a JOURNEY: several steps ending in a state change a player would
// describe. "Deal, mulligan, choose a house, play a creature, answer its prompt,
// undo it" is a scenario. "The reap button is disabled when the creature is
// exhausted" is not — that is a single-widget assertion and belongs in
// client_test.go, where 300-odd host tests already run in a second.
//
// HOW A SCENARIO GETS A KNOWN BOARD. Not from its seed. A deckgen.Set is built
// from the implemented card list, so the deal for a given seed changes every time
// a card is implemented — which is weekly here. A scenario that needs a named
// card therefore STAGES it: the manual preamble (manualMode plus addToHand) turns
// manual mode on through the real menu and adds that exact card through the real
// picker, so the card the steps act on is the one the scenario chose rather than
// the one the deal happened to deal. The seed is still fixed per scenario, so a
// failure reproduces exactly.
//
// A scenario or two deliberately SKIPS the preamble (opening, mulligan), so the
// real deal, the mulligan, and the house picker stay covered end to end. Their
// steps are written as predicates over whatever was dealt ("a house to choose",
// "a card in hand"), because which cards they get is genuinely not the point.
//
// NEVER NAME A PHYSICAL SIDE. Every step is written against the active player —
// the hand row is the active player's, the manual picker adds to the active
// player's hand, and no step asserts which row acted. Which player goes first is
// the engine's to decide and is changing (ADR 0040 makes it a seed-derived roll),
// so a scenario that named a row would fail on a switchover that changed nothing
// a player would notice.

// uiScenario is one browser journey: a name for the human, a slug for the URL, a
// seed that fixes the deal so each pass plays the same cards, and the ordered
// steps.
type uiScenario struct {
	Name  string
	Slug  string
	Seed  int64
	Steps []uiStep
}

// uiStep is one step of a journey: what it does, in the words the panel shows,
// and the do/check itself against the live DOM. Run returns nil when the step has
// happened. It returns an error — naming what it could not find — when it has
// not, so a scenario whose target is missing fails saying so instead of quietly
// passing over a board that never appeared.
type uiStep struct {
	Desc string
	Run  func(p *uiPage) error
}

// uiScenarios is every registered scenario, in the order /ui-test lists them.
var uiScenarios = []uiScenario{
	openingScenario(),
	mulliganScenario(),
	playAndUndoScenario(),
	cardPromptScenario(),
	reloadScenario(),
}

// UITestScenario is one registered scenario as a driver sees it: the name to
// report it under and the slug that addresses its page. It carries no steps,
// because a driver navigates to /ui-test/<slug> and reads the status element the
// page writes — the journey is defined here once and is never re-described
// outside this file.
type UITestScenario struct {
	Name string
	Slug string
}

// UITestScenarios lists every registered scenario, in the order /ui-test shows
// them. It exists so internal/web/uitest's headless driver enumerates the same
// registry the page does: adding a scenario to uiScenarios above adds it to both
// the page a human watches and the suite mage uiTest runs, with no second edit.
func UITestScenarios() []UITestScenario {
	out := make([]UITestScenario, len(uiScenarios))
	for i, s := range uiScenarios {
		out[i] = UITestScenario{Name: s.Name, Slug: s.Slug}
	}
	return out
}

// uiScenarioBySlug returns the registered scenario for a URL slug.
func uiScenarioBySlug(slug string) (uiScenario, bool) {
	for _, s := range uiScenarios {
		if s.Slug == slug {
			return s, true
		}
	}
	return uiScenario{}, false
}

// The cards the staged scenarios add by name. Each is chosen for what it does to
// the journey, not for what it is: the creature has no prompt of its own, so the
// play/undo journey is about playing and undoing; the Tactic's Play ability is a
// card decision over the hand, so the prompt journey has a prompt with real
// candidates whichever cards the deal dealt.
const (
	// stagedCreature is a plain creature to play: its Play ability is conditional
	// Æmber and asks nothing, so nothing interrupts the play it is there to prove.
	stagedCreature = "Flaxia"
	// stagedPromptCard is "Play: Archive a card from your hand." — a card prompt
	// whose candidates are the cards in hand, so it never auto-resolves the way a
	// single-creature board would (pickCreature answers a lone candidate itself).
	stagedPromptCard = "Labwork"
)

// openingScenario walks the opening of a match: both players keep the hands they
// were dealt, the first player picks a house, and the turn is under way. It is
// the shortest complete journey there is — every match begins with it — so it is
// also the one that proves the page itself runs. It stages nothing: the deal, the
// mulligan, and the house picker are exactly what it is covering.
func openingScenario() uiScenario {
	return uiScenario{
		Name: "Opening: keep both hands and choose a house",
		Slug: "opening",
		Seed: 1,
		Steps: []uiStep{
			{
				Desc: "Player 1 keeps the hand they were dealt",
				Run:  func(p *uiPage) error { return p.click("a Keep button", keepSel) },
			},
			{
				Desc: "Player 2 keeps the hand they were dealt",
				Run:  func(p *uiPage) error { return p.click("a Keep button", keepSel) },
			},
			{
				Desc: "the first player is asked to choose a house",
				Run: func(p *uiPage) error {
					if n := p.count(housePickSel); n < 1 {
						return fmt.Errorf(
							"no house to choose from: %s matched nothing",
							housePickSel,
						)
					}
					return nil
				},
			},
			{
				Desc: "choosing a house begins the turn",
				Run:  func(p *uiPage) error { return p.click("a house button", housePickSel) },
			},
			{
				Desc: "the turn is under way, and can be ended",
				Run: func(p *uiPage) error {
					if err := p.absent("the house picker", housePickSel); err != nil {
						return err
					}
					_, err := p.find("the End turn button", actSel(actEndTurn))
					return err
				},
			},
		},
	}
}

// mulliganScenario is the other opening: the player offered the first mulligan
// takes it, and the match carries on into the turn from the hand they drew
// instead. Like the opening it stages nothing — the point is that a real deal
// survives a real mulligan — so its last step is a predicate over whatever hand
// came back, not a named card.
func mulliganScenario() uiScenario {
	return uiScenario{
		Name: "Mulligan: shed an opening hand and take the turn anyway",
		Slug: "mulligan",
		Seed: 2,
		Steps: []uiStep{
			{
				Desc: "the player offered the first mulligan takes it",
				Run:  func(p *uiPage) error { return p.click("a Mulligan button", mulliganSel) },
			},
			{
				Desc: "the other player keeps the hand they were dealt",
				Run:  func(p *uiPage) error { return p.click("a Keep button", keepSel) },
			},
			{
				Desc: "choosing a house begins the turn",
				Run:  func(p *uiPage) error { return p.click("a house button", housePickSel) },
			},
			{
				Desc: "the turn is under way over a redrawn hand",
				Run: func(p *uiPage) error {
					if err := p.absent("the house picker", housePickSel); err != nil {
						return err
					}
					if _, err := p.find("a card in hand", anyHandCardSel); err != nil {
						return err
					}
					_, err := p.find("the End turn button", actSel(actEndTurn))
					return err
				},
			},
		},
	}
}

// playAndUndoScenario stages a creature by name, plays it, and takes it back.
// It is the journey behind every turn — pick a card up, put it down, change your
// mind — and it covers the manual edit too, since staging the creature is one:
// the card is added through the real picker, which records a command the reload
// journey then replays.
func playAndUndoScenario() uiScenario {
	steps := manualPreamble(stagedCreature)
	steps = append(steps,
		chooseAHouse(),
		playStagedCreature(),
		uiStep{
			Desc: stagedCreature + " is on the battleline and gone from hand",
			Run: func(p *uiPage) error {
				if _, err := p.find(stagedCreature+" in play",
					boardCardSel(stagedCreature)); err != nil {
					return err
				}
				return p.absent(stagedCreature+" in hand", handCardSel(stagedCreature))
			},
		},
		uiStep{
			Desc: "undo the play",
			Run:  func(p *uiPage) error { return p.click("the Undo button", actSel(actUndo)) },
		},
		uiStep{
			Desc: stagedCreature + " is back in hand and off the battleline",
			Run: func(p *uiPage) error {
				if _, err := p.find(stagedCreature+" in hand",
					handCardSel(stagedCreature)); err != nil {
					return err
				}
				return p.absent(stagedCreature+" in play", boardCardSel(stagedCreature))
			},
		},
	)
	return uiScenario{
		Name:  "Play and undo: stage a creature, play it, take it back",
		Slug:  "play-undo",
		Seed:  3,
		Steps: steps,
	}
}

// cardPromptScenario answers a card decision a card raised. It stages both sides
// of the decision — the Tactic that asks and the card the answer names — so the
// step that answers it clicks a card the scenario put there rather than whichever
// card the deal left at the front of the hand.
func cardPromptScenario() uiScenario {
	steps := manualPreamble(stagedPromptCard)
	steps = append(steps, addToHand(stagedCreature)...)
	steps = append(steps,
		chooseAHouse(),
		uiStep{
			Desc: "select " + stagedPromptCard + " in hand",
			Run: func(p *uiPage) error {
				return p.click(stagedPromptCard+" in hand", handCardSel(stagedPromptCard))
			},
		},
		uiStep{
			Desc: "play " + stagedPromptCard,
			Run:  func(p *uiPage) error { return p.click("its Play button", actSel(actPlay)) },
		},
		uiStep{
			Desc: "the prompt offers " + stagedCreature + " as a card to archive",
			Run: func(p *uiPage) error {
				_, err := p.find(stagedCreature+" highlighted as a candidate",
					candidateSel(handCardSel(stagedCreature)))
				return err
			},
		},
		uiStep{
			Desc: "answer the prompt with " + stagedCreature,
			Run: func(p *uiPage) error {
				return p.click(stagedCreature+" in hand", handCardSel(stagedCreature))
			},
		},
		uiStep{
			Desc: stagedCreature + " has left the hand for the archives",
			Run: func(p *uiPage) error {
				return p.absent(stagedCreature+" in hand", handCardSel(stagedCreature))
			},
		},
	)
	return uiScenario{
		Name:  "Card prompt: play a card that asks, and answer it",
		Slug:  "card-prompt",
		Seed:  4,
		Steps: steps,
	}
}

// reloadScenario proves a match outlives the page it was played on: a creature
// staged and played, then the client torn out of the tree and a fresh one stood
// up over the same storage slot, which rebuilds the match from the command log it
// saved. The teardown and rebuild are the host's, not a click's — a real browser
// reload would restart the run rather than the match — but what they drive is the
// client's own mount and resume, which is the path a reload takes.
//
// The staged creature is what makes the check sharp, twice over. A client that
// had dropped the log and dealt afresh would be back at the mulligan with no
// Flaxia anywhere, so finding it on the battleline says the log really replayed.
// And the play was legal only because manual mode was lifting the active-house
// restriction, so it comes back only if the mode itself replayed too — the
// journey walks the whole path that toggleManual's inSetManual record exists for
// (TestManualModeSurvivesAReload pins the same thing off-browser).
func reloadScenario() uiScenario {
	steps := manualPreamble(stagedCreature)
	steps = append(steps,
		chooseAHouse(),
		playStagedCreature(),
		uiStep{
			Desc: stagedCreature + " is on the battleline",
			Run: func(p *uiPage) error {
				_, err := p.find(stagedCreature+" in play", boardCardSel(stagedCreature))
				return err
			},
		},
		uiStep{
			Desc: "close the page on the match",
			Run:  func(p *uiPage) error { return p.dropClient() },
		},
		uiStep{
			Desc: "the board is gone with it",
			Run:  func(p *uiPage) error { return p.absent("a card", anyCardSel) },
		},
		uiStep{
			Desc: "open the page again",
			Run:  func(p *uiPage) error { return p.mountClient() },
		},
		uiStep{
			Desc: "the match resumes mid-turn with " + stagedCreature + " still in play",
			Run: func(p *uiPage) error {
				if err := p.absent("the mulligan question", keepSel); err != nil {
					return err
				}
				_, err := p.find(stagedCreature+" in play", boardCardSel(stagedCreature))
				return err
			},
		},
	)
	return uiScenario{
		Name:  "Reload: a match played, closed, and resumed",
		Slug:  "reload-resume",
		Seed:  5,
		Steps: steps,
	}
}

// manualPreamble is how a staged scenario gets a known board: both players keep,
// manual mode goes on through the real menu, and the named card is added to the
// active player's hand through the real picker. Everything after it acts on that
// card, so the journey is the same whatever the week's card pool deals.
//
// Manual mode stays on for the rest of the journey. It lifts the active-house and
// first-turn restrictions (Game.inActiveHouse, barredByFirstTurn), which is what
// lets a staged card be played on the turn it was added whatever house it belongs
// to — the alternative, staging only cards in a house the deal happened to give
// the active player, would put the deal back in charge of the journey.
func manualPreamble(card string) []uiStep {
	steps := []uiStep{
		{
			Desc: "one player keeps the hand they were dealt",
			Run:  func(p *uiPage) error { return p.click("a Keep button", keepSel) },
		},
		{
			Desc: "the other player keeps theirs",
			Run:  func(p *uiPage) error { return p.click("a Keep button", keepSel) },
		},
		{
			Desc: "open the game menu",
			Run:  func(p *uiPage) error { return p.click("the menu button", actSel(actMenu)) },
		},
		{
			Desc: "turn manual mode on",
			Run:  func(p *uiPage) error { return p.click("the Manual mode item", actSel(actManual)) },
		},
	}
	return append(steps, addToHand(card)...)
}

// addToHand is the manual edit itself, driven as a player drives it: open the
// card picker, search the pool by name, and click the row. The add is a recorded
// command, so a journey that reloads replays it and the rebuilt catalog hands the
// card back the same id.
func addToHand(card string) []uiStep {
	return []uiStep{
		{
			Desc: "open the card picker",
			Run: func(p *uiPage) error {
				return p.click("the Add card button", actSel(actManualAddCard))
			},
		},
		{
			Desc: "search the card pool for " + card,
			Run: func(p *uiPage) error {
				return p.fill("the picker's search box", pickerInputSel, card)
			},
		},
		{
			Desc: "add " + card + " to hand",
			Run: func(p *uiPage) error {
				return p.click(card+" in the picker", pickerRowSel(card))
			},
		},
		{
			Desc: card + " is in hand",
			Run: func(p *uiPage) error {
				_, err := p.find(card+" in hand", handCardSel(card))
				return err
			},
		},
	}
}

// chooseAHouse takes the active player out of the house picker and into their
// turn. Which house is taken is not the journey's business — in manual mode every
// deck house is offered — so it clicks the first one rather than naming one.
func chooseAHouse() uiStep {
	return uiStep{
		Desc: "choose a house to begin the turn",
		Run:  func(p *uiPage) error { return p.click("a house button", housePickSel) },
	}
}

// playStagedCreature selects the staged creature in hand and plays it. The
// battleline is empty at this point, so the play runs straight away instead of
// asking which flank — the flank question is a separate journey's business.
func playStagedCreature() uiStep {
	return uiStep{
		Desc: "select " + stagedCreature + " in hand and play it",
		Run: func(p *uiPage) error {
			if err := p.click(stagedCreature+" in hand", handCardSel(stagedCreature)); err != nil {
				return err
			}
			return p.click("its Play button", actSel(actPlay))
		},
	}
}

// The selectors every scenario is written from. They are the client's own hooks —
// the data-act values the controls carry, the data-card name every face carries,
// and the element-id prefix that says which zone a face is drawn in — so a
// scenario survives a relabel or a restyle the way a match on rendered text would
// not. A card's own element id cannot be used: a staged card is given the next
// free LocalID, which the scenario has no way to know.
var (
	// keepSel and mulliganSel are the two answers to the opening mulligan, and
	// housePickSel any button of the start-of-turn house picker.
	keepSel      = actSel(optionActID("Keep"))
	mulliganSel  = actSel(optionActID("Mulligan"))
	housePickSel = `.house-pick button[data-act^="house-"]`

	// pickerInputSel is the manual card picker's search box.
	pickerInputSel = "#" + pickerInputID

	// anyCardSel matches any card face on the page, and anyHandCardSel any face
	// drawn in the hand row — the predicates the unstaged scenarios read the board
	// with, since which card is there is not their point.
	anyCardSel     = ".card"
	anyHandCardSel = `.card[id^="hand-"]`
)

// pickerRowSel matches the card picker's row for a named card.
func pickerRowSel(card string) string { return `.picker-item[data-card="` + card + `"]` }

// handCardSel and boardCardSel match a named card's face in the hand row and on
// the board. The zone is read off the element id's prefix (handCardID,
// boardCardID), which is what tells a play from a card still in hand — and what
// keeps a step off the copies of the same card that carry a name but no id: the
// lifted card, the hover preview, and a prompt's source face.
func handCardSel(card string) string  { return namedCardSel("hand-", card) }
func boardCardSel(card string) string { return namedCardSel("card-", card) }

func namedCardSel(idPrefix, card string) string {
	return `.card[id^="` + idPrefix + `"][data-card="` + card + `"]`
}

// candidateSel narrows a card selector to a card a prompt is offering, so a step
// can say "the prompt offers this card" rather than only "this card is drawn".
func candidateSel(cardSel string) string { return cardSel + ".card--targetable" }

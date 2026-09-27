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

// openingScenario walks the opening of a match: both players keep the hands they
// were dealt, the first player picks a house, and the turn is under way. It is
// the shortest complete journey there is — every match begins with it — so it is
// also the one that proves the page itself runs.
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

// keepSel matches the Keep answer to the opening mulligan, and housePickSel any
// button of the start-of-turn house picker. Both are written from the data-act
// hooks the controls carry (optionActID, houseActID), so they survive a relabel
// or a restyle the way a text match would not.
var (
	keepSel      = actSel(optionActID("Keep"))
	housePickSel = `.house-pick button[data-act^="house-"]`
)

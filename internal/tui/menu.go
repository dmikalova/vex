package tui

import tea "github.com/charmbracelet/bubbletea"

// menuModel is the top-level menu that routes to the other screens.
type menuModel struct {
	choices []menuChoice
	cursor  int
	width   int
	height  int
}

type menuChoice struct {
	label string
	desc  string
	to    screen
	quit  bool
}

func newMenuModel() menuModel {
	return menuModel{
		choices: []menuChoice{
			{label: "Play game", desc: "two-player hotseat match", to: screenGame},
			{
				label: "Card explorer",
				desc:  "browse every card and read its text",
				to:    screenExplorer,
			},
			{
				label: "Card statistics",
				desc:  "cards per house by set, in total and by type",
				to:    screenStats,
			},
			{
				label: "Provenance",
				desc:  "how many original cards are implemented, by set",
				to:    screenProvenance,
			},
			{label: "Quit", desc: "exit vactrol", quit: true},
		},
	}
}

func (m menuModel) resize(w, h int) menuModel { m.width, m.height = w, h; return m }

func (m menuModel) Update(msg tea.Msg) (menuModel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "l":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "q":
		return m, tea.Quit
	case "enter", " ":
		c := m.choices[m.cursor]
		if c.quit {
			return m, tea.Quit
		}
		return m, gotoScreen(c.to)
	}
	return m, nil
}

func (m menuModel) View() string {
	s := titleStyle.Render(
		"VACTROL",
	) + "\n" + faintStyle.Render(
		"a KeyForge-like card game",
	) + "\n\n"
	for i, c := range m.choices {
		if i == m.cursor {
			s += selectedStyle.Render(
				cursor(true)+c.label,
			) + "  " + faintStyle.Render(
				c.desc,
			) + "\n"
		} else {
			s += cursor(false) + c.label + "\n"
		}
	}
	s += "\n" + helpStyle.Render("↑/↓ move · enter select · q quit")
	return s
}

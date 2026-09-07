package engine

// This file holds the Deploy keyword: a creature with Deploy may enter play at any
// position in its controller's battleline, not only on a flank. playCreatureCard
// in game_play.go calls deployPosition to place the creature.

// deployPosition decides where a creature being played enters its battleline. A
// plain creature lands on the flank flankLeft names. A Deploy creature's
// controller instead chooses any of the line's positions — the left flank, the
// right flank, or between any two creatures — so it can enter mid-line. interior
// reports a between-two-creatures landing, which the log narrates differently. An
// empty line has one position, so Deploy prompts nothing there.
func (g *Game) deployPosition(player int, id LocalID, flankLeft bool) (pos int, interior bool) {
	line := g.State.Battleline[player].slice()
	n := len(line)
	if !g.cat.def(id).hasKeyword(Deploy) || n == 0 {
		if flankLeft {
			return 0, false
		}
		return n, false
	}
	choice := g.choosePosition(player, id, "Choose where to deploy "+g.Name(id), line)
	return choice, choice > 0 && choice < n
}

// choosePosition asks a player where a Deploy creature enters the battleline.
// A chooser that speaks the battleline directly (PositionChooser) is pointed at
// the line; one that does not is offered the same positions as labeled gaps
// through the OptionChooser channel — "Left flank", "Right flank", or "Between X
// and Y". Both return the same position index: before line[i], 0..len(line).
func (g *Game) choosePosition(player int, source LocalID, prompt string, line []LocalID) int {
	if pc, ok := g.chooserFor(player).(PositionChooser); ok {
		name := g.sourceName(source)
		return pc.ChoosePosition(name, renderPrompt(name, prompt), line)
	}
	n := len(line)
	options := make([]string, n+1)
	options[0] = "Left flank"
	options[n] = "Right flank"
	for i := 1; i < n; i++ {
		options[i] = "Between " + g.Name(line[i-1]) + " and " + g.Name(line[i])
	}
	return g.ChooseOption(player, source, prompt, options)
}

package engine

// This file holds the Deploy keyword: a creature with Deploy may enter play at any
// position in its controller's battleline, not only on a flank. playCreatureCard
// in game_play.go calls deployPosition to place the creature.

// deployPosition decides where a creature entering a battleline lands. On an empty
// line there is one position and nothing to ask. When canDeploy is set and the
// creature has Deploy, its controller chooses any of the line's positions — the
// left flank, the right flank, or between any two creatures — so it can enter
// mid-line; interior reports a between-two-creatures landing, which the log
// narrates differently. Otherwise the creature lands on a flank: the flank fl
// names when the play dictated one, or the flank the active player is prompted for
// when fl is flankUnset. The engine never assumes a flank — an effect that puts a
// creature into play without naming a side always asks (ADR 0010's invalid-zero
// discipline, applied to placement).
//
// Every position is measured against the line as it stands *after* the chooser
// answers, never the line the prompt was drawn from. Asking crosses a resolution
// boundary that settles destroyed creatures (ADR 0029), so a line of two can be a
// line of one by the time the answer arrives, and the right flank the prompt
// offered is then one slot past the end
// (TestDeployPositionMeasuresTheLineAfterThePrompt).
func (g *Game) deployPosition(
	player int,
	id LocalID,
	fl flank,
	canDeploy bool,
) (pos int, interior bool) {
	line := g.State.Battleline[player].slice()
	if len(line) == 0 {
		return 0, false
	}
	if canDeploy && g.cat.def(id).hasKeyword(Deploy) {
		choice := g.choosePosition(player, id, "Choose where to deploy "+g.Name(id), line)
		n := int(g.State.Battleline[player].Count)
		choice = min(choice, n)
		return choice, choice > 0 && choice < n
	}
	if fl == flankUnset {
		// flankUnset: the play did not dictate a side, so the active player chooses.
		fl = g.chooseFlank(id)
	}
	if fl == flankLeftmost {
		return 0, false
	}
	return int(g.State.Battleline[player].Count), false
}

// chooseFlank asks the active player which flank a creature entering a non-empty
// battleline takes, when the play did not dictate one. The active player always
// makes this call, even when the creature enters the opponent's line (a creature
// put into play under the opponent). A non-Deploy creature can only land on a
// flank, so it offers just the two ends rather than the Deploy chooser's full set
// of interior gaps.
func (g *Game) chooseFlank(id LocalID) flank {
	prompt := FlankPromptPrefix + g.Name(id)
	if g.ChooseOption(
		g.State.ActivePlayer,
		id,
		prompt,
		[]string{FlankLeftLabel, FlankRightLabel},
	) == 0 {
		return flankLeftmost
	}
	return flankRightmost
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
	options[0] = FlankLeftLabel
	options[n] = FlankRightLabel
	for i := 1; i < n; i++ {
		options[i] = "Between " + g.Name(line[i-1]) + " and " + g.Name(line[i])
	}
	return g.ChooseOption(player, source, prompt, options)
}

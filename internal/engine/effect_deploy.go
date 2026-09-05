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
	options := make([]string, n+1)
	options[0] = "Left flank"
	options[n] = "Right flank"
	for i := 1; i < n; i++ {
		options[i] = "Between " + g.Name(line[i-1]) + " and " + g.Name(line[i])
	}
	choice := g.ChooseOption(player, id, "Choose where to deploy "+g.Name(id), options)
	return choice, choice > 0 && choice < n
}

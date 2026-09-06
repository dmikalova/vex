package engine

import (
	"strings"
	"testing"
)

// A WhileOffFlank constant self-buff applies only while the source sits off a
// flank (in the interior of its battleline), covering Gub.
func TestConstantAbilityWhileOffFlank(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Gub", Dis, Creature, Common, WithPower(1),
		WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetThisCreature},
			PowerBonus:    5,
			Keywords:      []Keyword{Taunt},
			WhileOffFlank: true,
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Gub gains +5 power and taunt while it is not on a flank.") {
		t.Error("card rules should render the while-off-flank line")
	}

	left := g.AddToBattleline(testCreature("left", 3), 0)
	gub := g.AddToBattleline(def, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	_ = left
	_ = right

	// Gub is the middle creature: off a flank, so the buff is active.
	if got := g.Power(gub); got != 6 {
		t.Errorf("off-flank power = %d, want 6 (1 + 5)", got)
	}
	if !g.hasKeyword(gub, Taunt) {
		t.Error("off-flank Gub should gain taunt")
	}

	// Remove the right creature so Gub becomes the right flank; the buff lifts.
	g.State.Battleline[0].remove(right)
	if got := g.Power(gub); got != 1 {
		t.Errorf("on-flank power = %d, want 1 (buff suspended)", got)
	}
	if g.hasKeyword(gub, Taunt) {
		t.Error("on-flank Gub should not have taunt")
	}
}

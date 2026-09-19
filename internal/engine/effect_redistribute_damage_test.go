package engine

import "testing"

// redistChooser answers the option prompts (player choice, yes/no) from a scripted
// queue and the placement prompts from an id queue.
type redistChooser struct {
	options []int
	ids     []LocalID
}

func (c *redistChooser) ChooseOption(_, _ string, _ []string) int {
	if len(c.options) > 0 {
		o := c.options[0]
		c.options = c.options[1:]
		return o
	}
	return 0
}

func (c *redistChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(c.ids) > 0 {
		id := c.ids[0]
		c.ids = c.ids[1:]
		return id, true
	}
	return cands[0], true
}

func TestRedistributeDamageText(t *testing.T) {
	want := "Redistribute the damage among a player's creatures"
	if got := (RedistributeDamage{}).Text(); got != want {
		t.Errorf("text = %q", got)
	}
}

func TestRedistributeDamageValidate(t *testing.T) {
	if err := validateEffect(RedistributeDamage{}); err != nil {
		t.Errorf("validate = %v", err)
	}
}

func TestRedistributeDamageMovesAllOntoOne(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	b := g.AddToBattleline(testCreature("b", 5), 0)
	g.SetDamage(a, 2)
	g.SetDamage(b, 1)
	// Choose own player, accept, then pile all 3 damage onto a.
	g.SetChooser(0, &redistChooser{
		options: []int{0, 0},
		ids:     []LocalID{a, a, a},
	})

	ctx := &EffectContext{
		Resolver:   g,
		Source:     a,
		Controller: 0,
	}
	RedistributeDamage{}.Resolve(ctx)

	if got := g.Damage(a); got != 3 {
		t.Errorf("a damage = %d, want 3", got)
	}
	if got := g.Damage(b); got != 0 {
		t.Errorf("b damage = %d, want 0", got)
	}
}

func TestRedistributeDamageDeclineLeavesDamage(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	g.SetDamage(a, 2)
	// Choose own player, then decline.
	g.SetChooser(0, &redistChooser{options: []int{0, 1}})

	ctx := &EffectContext{
		Resolver:   g,
		Source:     a,
		Controller: 0,
	}
	RedistributeDamage{}.Resolve(ctx)

	if got := g.Damage(a); got != 2 {
		t.Errorf("a damage = %d, want 2", got)
	}
}

func TestRedistributeDamageDestroysWhenPiledPastPower(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 2), 0)
	b := g.AddToBattleline(testCreature("b", 2), 0)
	g.SetDamage(a, 1)
	g.SetDamage(b, 1)
	// Choose own player, accept, pile both damage onto a to lethal.
	g.SetChooser(0, &redistChooser{
		options: []int{0, 0},
		ids:     []LocalID{a, a},
	})

	ctx := &EffectContext{
		Resolver:   g,
		Source:     a,
		Controller: 0,
	}
	RedistributeDamage{}.Resolve(ctx)

	if g.inPlay(a) {
		t.Error("a should be destroyed by lethal damage")
	}
	if !g.inPlay(b) {
		t.Error("b should survive")
	}
}

func TestRedistributeDamageNoDamageDoesNothing(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	g.SetChooser(0, &redistChooser{options: []int{0}})

	ctx := &EffectContext{
		Resolver:   g,
		Source:     a,
		Controller: 0,
	}
	RedistributeDamage{}.Resolve(ctx)

	if got := g.Damage(a); got != 0 {
		t.Errorf("a damage = %d, want 0", got)
	}
}

func TestRedistributeDamageChoosesOpponent(t *testing.T) {
	g := NewGame("A", "B", 1)
	own := g.AddToBattleline(testCreature("own", 5), 0)
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	g.SetDamage(own, 3)
	g.SetDamage(foe, 2)
	// Choose the opponent, accept; the sole enemy creature keeps its own damage.
	g.SetChooser(0, &redistChooser{options: []int{1, 0}})

	ctx := &EffectContext{
		Resolver:   g,
		Source:     own,
		Controller: 0,
	}
	RedistributeDamage{}.Resolve(ctx)

	if got := g.Damage(foe); got != 2 {
		t.Errorf("foe damage = %d, want 2", got)
	}
	if got := g.Damage(own); got != 3 {
		t.Errorf("own damage = %d, want 3 (untouched)", got)
	}
}

// redistNoPick accepts the option prompts but never names a placement creature,
// exercising Resolve's fallback to the first creature.
type redistNoPick struct{ options []int }

func (c *redistNoPick) ChooseOption(_, _ string, _ []string) int {
	if len(c.options) > 0 {
		o := c.options[0]
		c.options = c.options[1:]
		return o
	}
	return 0
}

func (redistNoPick) ChooseCreature(_, _ string, _ []LocalID) (LocalID, bool) {
	return 0, false
}

func TestRedistributeDamageFallsBackToFirstCreature(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	b := g.AddToBattleline(testCreature("b", 5), 0)
	g.SetDamage(b, 2)
	// Accept, but never pick a placement target: both units fall onto creatures[0].
	g.SetChooser(0, &redistNoPick{options: []int{0, 0}})

	ctx := &EffectContext{
		Resolver:   g,
		Source:     a,
		Controller: 0,
	}
	RedistributeDamage{}.Resolve(ctx)

	if got := g.Damage(a); got != 2 {
		t.Errorf("a damage = %d, want 2 (fallback pile onto first creature)", got)
	}
	if got := g.Damage(b); got != 0 {
		t.Errorf("b damage = %d, want 0", got)
	}
}

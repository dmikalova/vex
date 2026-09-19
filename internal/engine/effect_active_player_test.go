package engine

import "testing"

// recordController captures the controller it is resolved with and renders a
// single word, exercising ByActivePlayer's no-space Text branch.
type recordController struct {
	got  *int
	text string
}

func (r recordController) Text() string               { return r.text }
func (r recordController) Resolve(ctx *EffectContext) { *r.got = ctx.Controller }
func (r recordController) validate() error            { return nil }

func TestByActivePlayer(t *testing.T) {
	if got := (ByActivePlayer{
		Do: Destroy{Target: Target{Kind: TargetChosenArtifact}.House(activeHouse)},
	}).Text(); got != "that player destroys an artifact of that house" {
		t.Errorf("text = %q", got)
	}

	var ignore int
	if got := (ByActivePlayer{Do: recordController{
		got:  &ignore,
		text: "reap",
	}}).Text(); got != "that player reaps" {
		t.Errorf("single-word text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 1
	var seen int
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	ByActivePlayer{Do: recordController{got: &seen}}.Resolve(ctx)
	if seen != 1 {
		t.Errorf("inner resolved as controller %d, want the active player 1", seen)
	}
	if ctx.Controller != 0 {
		t.Errorf("controller not restored: got %d, want 0", ctx.Controller)
	}

	if (ByActivePlayer{Do: Destroy{Target: Target{Kind: TargetChosenArtifact}}}).validate() != nil {
		t.Error("a set target should validate")
	}
	if (ByActivePlayer{Do: Destroy{}}).validate() == nil {
		t.Error("an unset inner target should be invalid")
	}
}

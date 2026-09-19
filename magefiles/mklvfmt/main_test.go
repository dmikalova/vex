package main

import (
	"strings"
	"testing"
)

func TestFormatKeyedStructLiteral(t *testing.T) {
	src := `package massmutation

func f() {
	_ = ct.Setup{
		P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(ct.Bind(&bonesaw, Bonesaw))},
	}
}
`

	got, _, _, err := formatSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("formatSource returned error: %v", err)
	}
	if !strings.Contains(string(got), "P1: ct.Side{") ||
		!strings.Contains(string(got), "House: card.House.Dis,") ||
		!strings.Contains(string(got), "Hand:  ct.Cards(ct.Bind(&bonesaw, Bonesaw)),") {
		t.Fatalf("formatted output missing multiline keyed struct literal:\n%s", got)
	}
}

func TestFormatSourceLeavesMultilineStructsAlone(t *testing.T) {
	src := `package massmutation

func f() {
	_ = ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(ct.Bind(&bonesaw, Bonesaw)),
		},
	}
}
`

	got, _, _, err := formatSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("formatSource returned error: %v", err)
	}
	if !strings.Contains(string(got), "P1: ct.Side{") ||
		!strings.Contains(string(got), "House: card.House.Dis,") ||
		!strings.Contains(string(got), "Hand:  ct.Cards(ct.Bind(&bonesaw, Bonesaw)),") {
		t.Fatalf("formatted output changed already-correct multiline literal:\n%s", got)
	}
}

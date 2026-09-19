package main

import (
	"bytes"
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
	if !bytes.Contains(got, []byte("P1: ct.Side{")) ||
		!bytes.Contains(got, []byte("House: card.House.Dis,")) ||
		!bytes.Contains(got, []byte("Hand:  ct.Cards(ct.Bind(&bonesaw, Bonesaw)),")) {
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
	if !bytes.Contains(got, []byte("P1: ct.Side{")) ||
		!bytes.Contains(got, []byte("House: card.House.Dis,")) ||
		!bytes.Contains(got, []byte("Hand:  ct.Cards(ct.Bind(&bonesaw, Bonesaw)),")) {
		t.Fatalf("formatted output changed already-correct multiline literal:\n%s", got)
	}
}

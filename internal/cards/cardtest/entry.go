package cardtest

import (
	"fmt"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Entry is one card to place during setup: a definition, an optional handle to
// bind it to, and any upgrades to attach. Bare definitions passed to Cards become
// plain entries; Bind and Upgraded produce richer ones.
type Entry struct {
	def      engine.CardDefinition
	bind     *Card
	upgrades []upgrade
}

// upgrade is one upgrade to attach to an entry's card, with an optional handle to
// bind the upgrade itself to.
type upgrade struct {
	def  engine.CardDefinition
	bind *Card
}

// Cards normalizes a mix of card definitions and Entry values into a slice for a
// Side's zone, so a scenario can list bare cards and bound cards together:
//
//	InPlay: ct.Cards(AncientBear, ct.Bind(&troll, Troll), GiantSloth)
func Cards(items ...any) []Entry {
	out := make([]Entry, len(items))
	for i, it := range items {
		out[i] = toEntry(it)
	}
	return out
}

// Bind places a card and stores a handle to it in dst, filled in when the game is
// built. It is the named-handle alternative to positional lookup, and the only
// way to name a specific copy when a scenario has two of the same card.
func Bind(dst *Card, def engine.CardDefinition) Entry {
	return Entry{def: def, bind: dst}
}

// Upgraded places a host card with the given upgrades already attached. The host
// may be a bare definition or a Bind; each upgrade may likewise be a definition
// or a Bind.
func Upgraded(host any, upgrades ...any) Entry {
	e := toEntry(host)
	for _, u := range upgrades {
		e.upgrades = append(e.upgrades, toUpgrade(u))
	}
	return e
}

// toEntry converts a card definition or an Entry into an Entry.
func toEntry(x any) Entry {
	switch v := x.(type) {
	case Entry:
		return v
	case engine.CardDefinition:
		return Entry{def: v}
	default:
		panic(
			fmt.Sprintf(
				"cardtest: cannot use %T as a card entry (want a card definition or ct.Bind/ct.Upgraded)",
				x,
			),
		)
	}
}

// toUpgrade converts a card definition or a Bind into an upgrade spec.
func toUpgrade(x any) upgrade {
	switch v := x.(type) {
	case Entry:
		return upgrade{def: v.def, bind: v.bind}
	case engine.CardDefinition:
		return upgrade{def: v}
	default:
		panic(
			fmt.Sprintf(
				"cardtest: cannot use %T as an upgrade (want a card definition or ct.Bind)",
				x,
			),
		)
	}
}

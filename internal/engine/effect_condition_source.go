package engine

import "fmt"

// FlankPosition names which flank an OnFlank predicate asks about.
type FlankPosition uint8

const (
	// AnyFlank is met on either end of the battleline; Not inverts it.
	AnyFlank FlankPosition = iota
	// LeftFlank is met only on the left end.
	LeftFlank
	// RightFlank is met only on the right end.
	RightFlank
)

// OnFlank is met by a creature's position in its battleline. OfIt names the
// subject: the source card (false) or the context creature ctx.It (true). Where
// names which flank: AnyFlank (either end, invertible with Not — Glyxl Proliferator
// on a flank, Titan Librarian not, Malison's ctx.It on a flank), or a specific
// LeftFlank / RightFlank (Sinestra left, Dexus right). With OfIt set and no context
// creature it is not met.
type OnFlank struct {
	OfIt  bool
	Where FlankPosition
	Not   bool
}

// CondText renders the condition naming the subject and flank.
func (c OnFlank) CondText() string {
	subject := SelfName
	if c.OfIt {
		subject = "it"
	}
	switch c.Where {
	case LeftFlank:
		return "if " + subject + " is on the left flank"
	case RightFlank:
		return "if " + subject + " is on the right flank"
	default:
		if c.Not {
			return "if " + subject + " is not on a flank"
		}
		return "if " + subject + " is on a flank"
	}
}

// Met reports whether the subject sits on the named flank.
func (c OnFlank) Met(ctx *EffectContext) bool {
	subject := ctx.Source
	if c.OfIt {
		if !ctx.HasIt {
			return false
		}
		subject = ctx.It
	}
	if c.Where == AnyFlank {
		return onFlank(ctx, subject) != c.Not
	}
	if !ctx.Resolver.IsCreature(subject) {
		return false
	}
	line := ctx.Resolver.Battleline(ctx.Resolver.Controller(subject))
	if len(line) == 0 {
		return false
	}
	if c.Where == RightFlank {
		return line[len(line)-1] == subject
	}
	return line[0] == subject
}

// SourceInCenterOfBattleline is met while the source card sits in the center of
// its controller's battleline — the middle creature of an odd-sized line, with
// equal creatures to its left and right. An even-sized line has no center.
type SourceInCenterOfBattleline struct{}

// CondText renders the condition naming the source card.
func (SourceInCenterOfBattleline) CondText() string {
	return "if " + SelfName + " is in the center of your battleline"
}

// Met reports whether the source card sits in the center of its battleline.
func (SourceInCenterOfBattleline) Met(ctx *EffectContext) bool {
	return ctx.Resolver.InCenterOfBattleline(ctx.Source)
}

// SourceReady is met while the source card is ready (unexhausted) — Bellowing
// Patrizate damages each creature that enters play only while it is ready.
type SourceReady struct{}

// CondText renders the condition naming the source card.
func (SourceReady) CondText() string { return "if " + SelfName + " is ready" }

// Met reports whether the source card is currently ready.
func (SourceReady) Met(ctx *EffectContext) bool {
	return !ctx.Resolver.Exhausted(ctx.Source)
}

// SourceNeighborsAllOfHouse is met while every battleline neighbor of the source
// card belongs to House — Xanthyx Harvester cannot be used while it has a
// non-Mars neighbor, so its use is gated on this being met.
type SourceNeighborsAllOfHouse struct {
	House House
}

// CondText renders the condition, e.g. "if it has no non-Mars neighbor".
func (c SourceNeighborsAllOfHouse) CondText() string {
	return "if it has no non-" + c.House.String() + " neighbor"
}

// Met reports whether the source card has no neighbor off House.
func (c SourceNeighborsAllOfHouse) Met(ctx *EffectContext) bool {
	for _, n := range neighbors(ctx, ctx.Source) {
		if ctx.Resolver.House(n) != c.House {
			return false
		}
	}
	return true
}

// AemberOnThisAtLeast is met when at least Amount Æmber sits on the source card —
// [REDACTED] sacrifices itself once it has hoarded four or more. Not flips the
// sense to "fewer than Amount", so Crassosaurus purges itself when it captured too
// little.
type AemberOnThisAtLeast struct {
	Amount int
	// Not flips the sense: false is met at or above the threshold, true below it.
	Not bool
}

// CondText renders the condition clause, "fewer than" when Not flips the sense.
func (c AemberOnThisAtLeast) CondText() string {
	if c.Not {
		return fmt.Sprintf("if there are fewer than %d Æmber on it", c.Amount)
	}
	return fmt.Sprintf("if there are %d or more Æmber on it", c.Amount)
}

// Met reports whether the source card holds at least Amount Æmber, flipped by Not.
func (c AemberOnThisAtLeast) Met(ctx *EffectContext) bool {
	return (ctx.Resolver.AmberOn(ctx.Source) >= c.Amount) != c.Not
}

// CountersOnThisAtLeast is met when the source card carries at least N counters of
// Kind — The Big One wipes the board once ten or more fuse counters sit on it.
type CountersOnThisAtLeast struct {
	// Kind is the counter to count.
	Kind CounterKind
	// N is the threshold the count must reach.
	N int
}

// CondText renders the condition clause.
func (c CountersOnThisAtLeast) CondText() string {
	return fmt.Sprintf("if there are %d or more %ss on %s", c.N, c.Kind.noun(), SelfName)
}

// Met reports whether the source card holds at least N counters of Kind.
func (c CountersOnThisAtLeast) Met(ctx *EffectContext) bool {
	return ctx.Resolver.CountersOn(ctx.Source, c.Kind) >= c.N
}

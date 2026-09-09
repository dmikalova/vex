package engine

// This file holds the counts that read a "... this way" tally an earlier effect
// recorded in ctx.Produced during the same resolution — the producer/consumer
// channel a card scales a later amount by. Split out of effect_count.go.

// CardsDestroyed counts the cards the most recent destruction in this resolution
// removed from play — the "for each card destroyed this way" tally (Oath of
// Poverty gains 2 Æmber for each artifact it destroyed).
type CardsDestroyed struct{}

// Value returns how many cards the preceding Destroy actually removed.
func (CardsDestroyed) Value(ctx *EffectContext) int { return ctx.Produced.TotalDestroyed() }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsDestroyed) CountText() string { return "card destroyed this way" }

// CreaturesDestroyed counts the cards the most recent destruction in this
// resolution removed from play, rendered as creatures — the "for each creature
// destroyed this way" tally (Martyr's End gains 1 Æmber for each friendly
// creature it destroyed). Use it when the destruction removes only creatures;
// use CardsDestroyed when it can also remove artifacts.
type CreaturesDestroyed struct{}

// Value returns how many cards the preceding Destroy actually removed.
func (CreaturesDestroyed) Value(ctx *EffectContext) int { return ctx.Produced.TotalDestroyed() }

// CountText renders the singular noun the "for each" clause repeats.
func (CreaturesDestroyed) CountText() string { return "creature destroyed this way" }

// ProducedTally names a "... this way" tally an earlier effect in the same
// resolution records in ctx.Produced for a following ProducedThisWay count.
type ProducedTally uint8

const (
	// producedTallyUnset is the invalid zero value; a ProducedThisWay must name one.
	producedTallyUnset ProducedTally = iota
	// TallyCreaturesDestroyed is a player's share of ctx.Produced.Destroyed.
	TallyCreaturesDestroyed
	// TallyCreaturesShuffledIntoDeck is a player's share of ctx.Produced.Moved.
	TallyCreaturesShuffledIntoDeck
	// TallyAemberLost is a player's share of ctx.Produced.AemberLost.
	TallyAemberLost
	// TallyCardsReturned is ctx.Produced.Returned, a whole tally not split by player.
	TallyCardsReturned
)

// ProducedThisWay counts a "... this way" tally an earlier effect in the same
// resolution recorded — creatures destroyed or shuffled home, Æmber lost, or cards
// returned. Tally names which; Player names whose share. The per-player tallies
// flip to each player under a GainAember{Player: EachPlayer}, so Controller there
// means each player counting their own. TallyCardsReturned is a whole tally and
// ignores Player.
type ProducedThisWay struct {
	Tally  ProducedTally
	Player Player
}

// Value reads the named tally — a player's share of the per-player ones, the whole
// of TallyCardsReturned.
func (c ProducedThisWay) Value(ctx *EffectContext) int {
	switch c.Tally {
	case TallyCreaturesDestroyed:
		return ctx.Produced.Destroyed[ctx.PlayerFor(c.Player)]
	case TallyCreaturesShuffledIntoDeck:
		return ctx.Produced.Moved[ctx.PlayerFor(c.Player)]
	case TallyAemberLost:
		return ctx.Produced.AemberLost[ctx.PlayerFor(c.Player)]
	default:
		return ctx.Produced.Returned
	}
}

// CountText renders the singular noun the "for each" clause repeats.
func (c ProducedThisWay) CountText() string {
	switch c.Tally {
	case TallyCreaturesDestroyed:
		who := "they"
		if c.Player == Opponent {
			who = "your opponent"
		}
		return "creature " + who + " controlled that was destroyed this way"
	case TallyCreaturesShuffledIntoDeck:
		whose := "their"
		if c.Player == Opponent {
			whose = "your opponent's"
		}
		return "creature shuffled into " + whose + " deck this way"
	case TallyAemberLost:
		who := "you"
		if c.Player == Opponent {
			who = "your opponent"
		}
		return "Æmber " + who + " lost this way"
	default:
		return "card put into your hand this way"
	}
}

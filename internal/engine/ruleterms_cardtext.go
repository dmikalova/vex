package engine

// Card text rulebook terms (ADR 0018): the parts of a card's sentence that are
// not the change itself — what an effect reaches, how many, whether it happens,
// and how long it lasts. An effect is the verb; these are the words around it.
// Each term is an umbrella: many nodes share one title, so a reader learns the
// shape once instead of once per node.
func init() {
	registerRuleSectionIntro(
		SectionCardText,
		`A card's sentence is built from an effect and the words around it. Those words
say which cards the effect reaches, how large it is, whether it happens at all,
and how long it lasts. They change nothing on their own: they shape the effect
they are attached to. Each term below covers a whole family of such words, so a
card that says "each enemy creature", "for each Æmber in your pool", or "for the
remainder of the turn" is read the same way wherever it appears.`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section:    SectionCardText,
			Title:      "Target",
			Definition: "The part of a card's text that says which cards an effect reaches.",
			Body: `A target is how card text names the cards an effect reaches. It says whose
cards are eligible (friendly, enemy, or either), what they are (creature,
artifact, or any card in play), and how they are picked — "a creature" is one
card the controller chooses, "each creature" is every eligible card at once, and
"this creature" is the card the ability is printed on. A random pick, the top
card of a deck, and a card named by name are targets too.

Extra clauses narrow the same target further: a house ("a Dis creature"), a
trait ("a Beast"), a keyword, a power threshold, or a position in the
battleline. Every clause narrows; none widens. An effect with no eligible card
does nothing, and a choice the controller may decline stays declinable even when
one eligible card is left.`,
		},
		{
			Section:    SectionCardText,
			Title:      "For Each",
			Definition: "A count that scales an effect: the effect's size is multiplied by how many things the count finds.",
			Body: `A "for each" clause scales an effect by a count taken when the effect
resolves — "gain 1 Æmber for each friendly Mars creature", "draw a card for each
card in your discard pile". The count is read at that moment, so a creature that
left play before the effect resolved does not count and one that entered does.

A count can measure the board (cards in play, cards in a hand, Æmber in a pool,
creatures used this turn) or what the ability itself just produced ("for each
card destroyed this way"). A count of zero makes the effect do nothing rather
than something.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Duration",
			Definition: "How long a lasting effect stays in force before it lifts.",
			Body: `A lasting effect names the window it holds for. The windows are:

- **Rest of this turn** — in force now, lifting at the end of the current turn.
- **Opponent's next turn** — dormant now, biting only during that player's next
  turn and lifting when it ends.
- **Until your next turn** — in force now, through the opponent's turn, lifting
  at the start of the controller's next turn.
- **Through your next turn** — in force now and through the whole of that
  player's own next turn, lifting when it ends.
- **Until this leaves play** — held by the card whose ability established it,
  lifting when that card leaves play.
- **Until it leaves play** — held by the affected card, lifting when that card
  leaves play.

"Next turn" always means that player's next turn, not the next turn anyone
takes, so an extra turn taken by someone else never consumes the window.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Quantity",
			Definition: "How many cards a selection takes: exactly that many, up to that many, or any number.",
			Body: `A card that moves or picks several cards says how many. "2 cards" takes exactly
that many, as many as are available. "Up to 2 cards" and "any number of cards"
let the controller stop early: they pick one at a time and may decline at any
point, including at the first pick, so taking none is allowed.

The count itself can be a fixed number or a count read from the board.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Portion",
			Definition: "A share of a quantity — a half or a third — with its rounding always stated.",
			Body: `A portion is a share of a quantity: half of it, or a third of it. The rounding is
always printed, so "half your Æmber, rounded up" and "half your Æmber, rounded
down" are different instructions and neither is assumed.

The same portion applies wherever a share is taken — of a player's Æmber pool, of
a creature's power, of the creatures on a battleline. The share is measured when
the effect resolves.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Power Threshold",
			Definition: "A clause that admits only creatures whose power compares to a number, or to this creature's power.",
			Body: `Card text can narrow what an effect reaches by power: power 3 or lower, power 3
or higher, power exactly 1, odd power, even power, or less power than the card
the ability is printed on.

Power is read when the effect resolves, and it is the creature's current power —
its counters, its upgrades, and the constant abilities in play included — not its
printed power. A creature whose power changes before the effect resolves is
admitted or excluded by its new power.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Power Tier",
			Definition: "A clause that keeps the most or least powerful of a set, either one creature or every creature tied at that power.",
			Body: `Card text can pick creatures by where their power stands in a set rather than
against a number. Two shapes exist, and the difference is what a tie does.

One shape keeps a fixed number of creatures — the most powerful creature, the
least powerful creature, the 3 most powerful creatures. When creatures tie for
the last place, the controller chooses which of the tied creatures is kept.

The other shape keeps a tier: every creature tied for the highest power, or every
creature tied for the lowest. It makes no choice, and it keeps all of them.

Either way an empty set keeps nothing.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Same Power as a Chosen Creature",
			Definition: "A clause that keeps every creature whose power matches that of a creature the controller chooses.",
			Body: `Card text can reach every creature whose power matches a creature the controller
chooses — "choose a creature; destroy each creature with that power". An effect
may ask for more than one such creature, in which case every matching power
counts.

The choices are all made before the set narrows, so each is made against the
board as it stands. A choice the controller declines contributes no power.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Creatures of a House",
			Definition: "A clause that admits creatures belonging to a house with at least that many creatures in play.",
			Body: `Card text can reach every creature that belongs to a house with a given number of
creatures in play — "each creature that belongs to a house with 3 or more
creatures in play".

A house is counted across both battlelines. Who controls a creature does not
matter to the count: a house is a house on either side of the board.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Per Battleline",
			Definition: "A clause applied to each battleline separately, so each player's creatures are counted and picked on their own.",
			Body: `Some clauses apply to each battleline separately rather than to the board as a
whole. "Each player keeps 2 creatures and destroys the rest" and "destroy a third
of the creatures on each battleline" are both read side by side: each side is
counted on its own, and the controller of the effect makes the picks.

Both sides are read from the board as it stands before anything is removed. A
side with no more creatures than the clause spares is left alone, with no choice
to make.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Tide",
			Definition: "A game-wide state that is high for one player and low for their opponent, and neutral until a card raises it.",
			Body: `The tide is a single game-wide state. It starts neutral. A card can raise it,
which makes the tide high for the player who raised it and low for their
opponent at the same time; it is never high or low for both.

While the tide is neutral it is neither high nor low for either player, so a
card that acts "while the tide is high" and one that acts "while the tide is low"
both do nothing.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Overwhelmed",
			Definition: "A player is overwhelmed while their opponent controls more creatures than they do.",
			Body: `A player is overwhelmed while their opponent controls more creatures than they
control. Equal numbers are not overwhelmed.

Being overwhelmed restricts nothing on its own. It is a board state cards name,
checked when the card that names it resolves.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Haunted",
			Definition: "A player is haunted while their discard pile holds 10 or more cards.",
			Body: `A player is haunted while their own discard pile holds 10 or more cards.

Being haunted restricts nothing on its own. It is a board state cards name,
checked when the card that names it resolves, so a player can become haunted and
stop being haunted several times in a game.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Flank",
			Definition: "A creature at either end of a battleline; card text can reach flank creatures, or spare them.",
			Body: `The creatures at the two ends of a battleline are on a flank. A battleline of
one creature has that creature on both flanks, and a creature with no neighbor
at all is on a flank.

Card text reads the flank both ways: an effect can reach a creature on a flank
or a creature not on a flank, and an ability can turn on whether the card it is
printed on is on a flank, on one named flank, or in the center of its
battleline — the middle creature of an odd-sized line, which an even-sized line
does not have. An ability can also turn on having no neighbor of a named house.

An effect can also make a creature count as a flank creature for the remainder
of the turn wherever it actually sits. That is a lasting override, not a move:
the creature does not change position, but every flank check treats it as being
on one until the override lifts at the ready step.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Trait",
			Definition: "A word printed under a card's name, such as Beast or Scientist, that card text can name.",
			Body: `A trait is a word printed on a card beneath its name — Beast, Human, Scientist,
Mutant. A trait does nothing on its own. It exists so other cards can name it.

Card text uses traits to narrow what an effect reaches ("a Beast", "a creature
that shares a trait with this creature", "a creature with no trait in common
with it"), to test the card in context, and to count ("for each trait on the
chosen creature"). An effect can also grant a creature a trait it does not have
printed, and a granted trait reads the same as a printed one.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Conditional",
			Definition: "Gate an effect behind an \"If ...\" check on the current board; it resolves only when the condition is met.",
			Body: `A conditional gates an effect behind a check on the current game state — the
"If ..." clause a card opens with, e.g. "If your opponent has 7 or more Æmber,
they lose 4 Æmber." The effect resolves only when the condition is met. Unlike
a result gate (A -> B), which turns on an action succeeding, a conditional turns
on a fact about the board.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Result Gate",
			Definition: "Resolve one action, then a follow-up (A -> B) only when the first action actually happened.",
			Body: `A result gate resolves one action and then a follow-up, but only when the first
action actually happened — written A -> B (destroy a creature -> steal 1 Æmber;
purge a creature -> give a +1 power counter). The follow-up never runs when the
gate does nothing: no valid target, an empty zone, or a declined choice. It is
distinct from a conditional, which turns on a fact about the board rather than an
action succeeding.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Repeat",
			Definition: "Resolve an effect again while its stated condition holds, bounded by the six uses a card name has each turn.",
			Body: `An effect that repeats resolves again while its stated condition holds. The
condition is a fact about the board, so the effect may do nothing and the repeat
still happens (Numquid the Fair comes back while you are overwhelmed, even when
a ward absorbs its destroy). When the repeat clause is written with -> it is a
result gate: the effect also stops repeating as soon as it does nothing (Bait and
Switch stops the moment a steal moves no Æmber). A repeat is always bounded.
Every resolution past the first counts against the six uses a card name has each
turn, so a condition the effect cannot change still ends the loop.`,
		},
		{
			Section:    SectionCardText,
			Title:      "May",
			Definition: "Offer the controller the choice to resolve the inner effect or decline it entirely.",
			Body: `A "you may" effect is optional: it offers the controller the choice to resolve
its inner effect or to decline it entirely. It models KeyForge's "You may <do
X>", where passing is always allowed even when a legal target exists — the
distinction that keeps Chuff Ape's "you may destroy another friendly creature"
from ever being forced.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Choose One",
			Definition: "Offer the controller a set of alternative effects and resolve only the one they pick.",
			Body: `A "choose one" ability offers its controller a set of alternative effects and
resolves only the one they pick; the options not chosen do nothing.`,
		},
		{
			Section:    SectionCardText,
			Title:      "Restriction",
			Definition: "Forbid a player an action for a stretch of the game; while active, \"cannot\" beats any \"must\" or \"may\".",
			Body: `A restriction forbids a player some action for a stretch of the game — "cannot
use creatures to fight", "cannot play creatures" — rather than changing the board
directly. A restriction can be a timed effect that lasts through a player's next
turn, or a constant rule printed on a card in play; while it is active the
forbidden action simply cannot be taken. When one effect says a player "cannot"
and another says they "must" or "may" do the same thing, "cannot" wins.`,
		},
	})
}

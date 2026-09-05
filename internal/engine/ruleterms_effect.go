package engine

// Effects rulebook terms (ADR 0018): each describes itself next to the code it
// governs; the completeness test fails the build if a member of the matching
// closed catalog has no term here.
func init() {
	registerRuleSectionIntro(
		SectionEffect,
		`An effect is the actual change an ability makes to the game when it resolves.
Card text is built by composing the effects below; each one both prints its own
rules text and carries itself out, so a card always does exactly what it says.`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section:    SectionEffect,
			Title:      "Gain Æmber",
			Definition: "A player moves that many Æmber from the common supply into their pool.",
			Body: `To gain Æmber, a player moves that many Æmber from the common supply into
their pool — the ability's controller by default, or their opponent when the
card says so. A "for each" clause multiplies the amount by a running count.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Lose Æmber",
			Definition: "A player returns that many Æmber from their pool to the common supply, never below zero.",
			Body: `To lose Æmber, a player returns that many Æmber from their pool to the common
supply. A pool can never go below zero, so a player told to lose more Æmber than
they have simply loses all of it. Player may be EachPlayer, so both players lose.
The amount lost is either a fixed Amount or a By loss of the pool (By: Half,
By: AllBut(5)) — set one, not both.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Move Æmber to a Card",
			Definition: "Take that many Æmber out of your pool and set it on a card, where it waits until something moves it off.",
			Body: `Moving Æmber from your pool to a card takes that many Æmber out of your pool
and sets it on the card, where it stays until something moves it off again.
Parked there it only matters to a card that can spend the Æmber sitting on it —
Safe Place, Pocket Universe.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Archive",
			Definition: "Set cards aside face-down in your archives, out of the opponent's reach; take them to hand after picking a house on a later turn.",
			Body: `Archiving moves cards into your archives: they are set aside face-down, out of
the opponent's reach, and you may take them into your hand after picking a
house on a later turn. Archiving from hand lets you choose which cards to set
aside.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Lose Armor",
			Definition: "Strip all remaining armor off the chosen creatures; it returns when their controller readies.",
			Body: `LoseArmor takes all the remaining armor off each creature its Target selects,
and tallies what it took so a following effect can scale with it (Red-Hot Armor
strips armor, then deals damage for each point stripped). The armor comes back
when its controller readies, the same way armor spent absorbing damage does.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Capture Æmber",
			Definition: "Move Æmber from a pool onto a creature, where it counts for no one until that creature leaves play, then goes to the capturer's opponent.",
			Body: `Capturing Æmber moves it from a player's pool onto a capturing creature, where
it counts for no player until that creature leaves play, at which point it goes
to the pool of the capturing creature's controller's opponent. A creature can
only capture what the Source pool holds. Target is the creature that captures
(this creature by default); Source is the pool the Æmber comes from; Per repeats
the capture, choosing a fresh Target each time (Hypnotic Command captures once
for each friendly Mars creature).`,
		},
		{
			Section:    SectionEffect,
			Title:      "Gain Chains",
			Definition: "Take on chains, which reduce your card draw by one for every six held until they are shed.",
			Body: `A chain is a penalty a card can inflict on its controller: while a player holds
chains they draw fewer cards each turn — one fewer for every 6 chains — until the
chains are shed, one on each turn the reduction blocks a draw. Gaining a chain is
the cost some strong effects charge, so a card's power is paid for by a slower
hand refill (see Game.drawStep).`,
		},
		{
			Section:    SectionEffect,
			Title:      "Choose One",
			Definition: "Offer the controller a set of alternative effects and resolve only the one they pick.",
			Body: `A "choose one" ability offers its controller a set of alternative effects and
resolves only the one they pick; the options not chosen do nothing.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Conditional",
			Definition: "Gate an effect behind an \"If ...\" check on the current board; it resolves only when the condition is met.",
			Body: `A conditional gates an effect behind a check on the current game state — the
"If ..." clause a card opens with, e.g. "If your opponent has 7 or more Æmber,
they lose 4 Æmber." The effect resolves only when the condition is met. Unlike
a result gate (A -> B), which turns on an action succeeding, a conditional turns
on a fact about the board.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Power Counter",
			Definition: "A permanent token that raises a creature's power by one (+1) or lowers it by one (-1) while it stays in play.",
			Body: `A +1 power counter is a permanent token placed on a creature that raises its
power by one for as long as it stays in play; a -1 power counter lowers it. A
creature can hold any number of counters, and they are shed when it leaves play.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Deal Damage",
			Definition: "Put pending damage on the targeted creatures; armor stops some, and a creature whose damage reaches its power is destroyed.",
			Body: `Dealing damage puts that much pending damage on each creature the effect
targets. Armor prevents pending damage first — each point stops 1, and armor
spent this way stays spent for the rest of the turn — and whatever is not
prevented lands as damage tokens. A creature whose total damage reaches or
exceeds its power is destroyed. When one ability deals damage to several
creatures they are damaged simultaneously and any that died are destroyed
together, so no creature's destruction changes another's.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Reveal Top of Deck",
			Definition: "Reveal the top card of your deck so a following effect can inspect or play it; the card does not move.",
			Body: `RevealTopOfDeck reveals the top card of the controller's deck — logging it and
putting it in context (ctx.It) so a following effect can inspect or play it (Chaos
Portal plays it when it is of the chosen house). Revealing does not move the card;
an empty deck reveals nothing.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Discard",
			Definition: "Dig through the top of your deck, discarding as you go, until you turn up a card the filters admit or the deck runs out.",
			Body: `DiscardDeckUntil digs through the top of your deck, discarding as it goes,
until it turns up a card the filters admit or the deck runs out. The card it
finds stays in the discard pile and goes into context (ctx.It), so what happens
to it is a separate effect gated on the dig succeeding — Sound the Horns and
Invasion Portal both pair it with PutDiscardedIntoHand.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Return",
			Definition: "Take the card just found in the discard pile into its owner's hand.",
			Body: `PutDiscardedIntoHand takes the card in context out of the discard pile and
into its owner's hand. It is the tail of a dig through the deck (DiscardDeckUntil)
that just discarded the card. Type names what the dig stopped on so the tail
reads "put the discarded creature into your hand" rather than a bare "it".`,
		},
		{
			Section:    SectionEffect,
			Title:      "Destroy",
			Definition: "Remove a creature from play; its Destroyed abilities resolve first, then it and its upgrades go to the discard pile.",
			Body: `Destroying a creature removes it from play. When an effect destroys several
creatures they are destroyed simultaneously: every one is tagged for
destruction and stays in play while their "Destroyed:" abilities resolve, in an
order the controller chooses, so each ability sees the others still present;
only then does each creature still in play move to the discard pile, along with
its upgrades. A destroy effect can target every creature or only those matching
a filter, such as "each creature with power 3 or lower".`,
		},
		{
			Section:    SectionEffect,
			Title:      "Draw",
			Definition: "Put the top card of your deck into your hand, reshuffling your discard pile into a new deck first if the deck is empty.",
			Body: `Drawing puts the top card of your deck into your hand. If your deck is empty
when you must draw, your discard pile is shuffled to form a new deck first, so
you only fail to draw when both deck and discard are empty.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Exalt",
			Definition: "Place 1 Æmber from the common supply onto a creature, where it waits until that creature leaves play, then goes to the owner's opponent.",
			Body: `To exalt a creature is to place 1 Æmber from the common supply onto a chosen
friendly or enemy creature. The Æmber sits on the creature, belonging to no
pool, until it leaves play, then goes to the owner's opponent's pool. Exalting
N times places N Æmber.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Exhaust",
			Definition: "Turn a creature sideways so it cannot be used again until it readies at the end of its controller's turn.",
			Body: `Exhausting a creature turns it sideways so it cannot be used again until it
readies at the end of its controller's turn. It exhausts each creature the
effect targets.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Ready",
			Definition: "Turn a creature upright again so it can be used, the opposite of exhausting.",
			Body: `Readying a creature turns it upright again so it can be used, the opposite of
exhausting. It readies each creature the effect targets.`,
		},
		{
			Section:    SectionEffect,
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
			Section:    SectionEffect,
			Title:      "Heal",
			Definition: "Take damage tokens off a creature; it never removes more than is there and never changes power.",
			Body: `Healing takes damage tokens off a creature — a fixed amount, or all of them at
once. It can never remove more damage than is on the creature (a creature with
no damage is unaffected), and it never changes a creature's power.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Belong to House",
			Definition: "Make the chosen creatures count as a named house for the given duration, overriding their printed house.",
			Body: `BelongToHouse makes each creature its Target selects belong to House for the
given Duration, overriding the house it counts as for active-house checks (Brain
Stem Antenna's host counts as Mars for the rest of the turn). The change is
per-match state, dropped when the creature leaves play; EndOfTurn also drops it at
end of turn, while UntilThisLeavesPlay keeps it until the creature leaves play.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Name a House",
			Definition: "Record the house chosen on this card so a HouseLock can bar the opponent from it while the card stays in play.",
			Body: `NameHouse remembers the house a surrounding ChooseHouseThen picked on the source
card, where it stays for as long as that card is in play. It is the writer half
of a HouseLock whose house is not printed but named: Restringuntus chooses a
house on play and bars its opponent from it until it leaves play. Player names
whose choice the lock will constrain, and must match the card's HouseLock.`,
		},
		{
			Section:    SectionEffect,
			Title:      "May",
			Definition: "Offer the controller the choice to resolve the inner effect or decline it entirely.",
			Body: `A "you may" effect is optional: it offers the controller the choice to resolve
its inner effect or to decline it entirely. It models KeyForge's "You may <do
X>", where passing is always allowed even when a legal target exists — the
distinction that keeps Chuff Ape's "you may destroy another friendly creature"
from ever being forced.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Play a Card from Hand or Discard Pile",
			Definition: "The controller plays a card from their hand or discard pile right now, ignoring the active-house gate.",
			Body: `PlayFrom has the controller play a card out of their own hand or discard pile
right now, ignoring the active-house gate — Phase Shift's off-house card,
Sacrificial Altar's creature back from the discard pile. From names the source
pile; House and Type narrow which cards may be chosen; Except inverts the house
filter, so House names the house that may *not* be played ("a non-Logos card").

KeyForge prints the from-hand form as a permission held open for the rest of the
turn ("you may play one non-Logos card this turn"); it is rendered and resolved
as an immediate play instead, which needs no turn-scoped memory of an unspent
allowance (see card-wording-rules.md rule 21). With no legal card in the source
pile it does nothing.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Purge",
			Definition: "Set a card aside out of the game in the purge pile, where no ability can reach it; it never returns.",
			Body: `Purging a card sets it aside out of the game entirely, in the purge pile, where
no ability can reach it unless that ability names the purge pile. It is the most
permanent way a card leaves play: a purged card never enters a discard pile and
can never be drawn, played, or destroyed again.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Put from Hand",
			Definition: "The controller puts a card chosen from their hand directly into play.",
			Body: `PutFromHand puts a card the controller chooses from their own hand directly
into play — Swap Widget swapping in a replacement creature. Type restricts the
choice to cards of that type; House restricts it to that house; either's zero
value allows any. ExceptSameName excludes a card sharing the name of the card
currently in context (ctx.It), the "with a different name" clause — meant to
follow a gate that left a card in context, e.g. Then{PutFromPlay, PutFromHand}.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Redirect Fight Damage",
			Definition: "Before a fight, the attacker deals its outgoing fight damage to a chosen creature instead of the one it fights.",
			Body: `RedirectFightDamage is a "Before Fight" effect: the controller chooses a
creature (its Target), and the attacker deals its own fight damage to that
creature instead of to the creature it is fighting (Gabos Longarms). It only
redirects the attacker's outgoing fight damage — the attacker still takes
damage back from the creature it fights. The chosen creature is stored on the
game state for the fight in progress; the combat step reads and clears it.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Restriction",
			Definition: "Forbid a player an action for a stretch of the game; while active, \"cannot\" beats any \"must\" or \"may\".",
			Body: `A restriction forbids a player some action for a stretch of the game — "cannot
use creatures to fight", "cannot play creatures" — rather than changing the board
directly. A restriction can be a timed effect that lasts through a player's next
turn, or a constant rule printed on a card in play; while it is active the
forbidden action simply cannot be taken. When one effect says a player "cannot"
and another says they "must" or "may" do the same thing, "cannot" wins.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Put from Play",
			Definition: "Take the chosen cards out of play to a named zone — deck top, hand, or archives — shedding their in-play state.",
			Body: `PutFromPlay takes each card its Target selects out of play and puts it in a
destination zone — the top of its owner's deck, their hand, or their archives —
shedding the per-match state the card built up in play (damage, spent armor,
Æmber on it, upgrades). The destination is required. Moving a card out of play
this way is how a "Destroyed:" ability can save its own creature: the creature
leaves for the named zone as it is destroyed, so it never reaches the discard
pile. When several cards move to the top of the deck at once the controller
chooses the order they stack.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Return Named Card to Hand",
			Definition: "Put a card of a chosen name into your hand, from a friendly creature in play or your discard pile.",
			Body: `ReturnNamedToHand puts a card with a specific name that the controller chooses
into their hand, taken either from a friendly creature in play or from their
discard pile — Faygin recovering an Urchin. The controller chooses among both
zones at once; an in-play creature returns to hand (shedding its in-play state)
and a discard card is recovered.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Search for Named Card",
			Definition: "Search your deck and discard pile for a named card, reveal it, and put it into your hand.",
			Body: `SearchForName lets the controller search their deck and discard pile for a card
with a specific name, reveal it, and put it into their hand — Help from Future
Self tutoring a Timetraveller. Nothing happens if no matching card is found.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Shuffle Into Deck",
			Definition: "Move the named zones' cards into your deck, then shuffle once.",
			Body: `ShuffleIntoDeck shuffles the controller's named Zones into their deck — the
discard pile (Help from Future Self), the hand and discard pile (Screaming
Cave), or the archives and discard pile. It moves every named zone's cards into
the deck, then shuffles once.`,
		},
		{
			Section: SectionEffect,
			Title:   "Shuffle Into Deck",
			Body: `SwapDeckAndDiscard exchanges the controller's deck with their discard pile and
shuffles the new deck — Reverse Time turns a spent deck back into a fresh one.
It differs from ShuffleIntoDeck{Discard} in that the old deck goes away into
the discard pile rather than surviving underneath it.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Steal Æmber",
			Definition: "Move Æmber from the opponent's pool into your own, up to what they have.",
			Body: `Stealing Æmber moves it from the opponent's pool into your own. You can only
steal as much Æmber as the opponent actually has. How much is stolen is either
a fixed Amount — optionally multiplied by a Per count — or a By share of the
opponent's pool (By: AllBut(6) leaves them exactly six).`,
		},
		{
			Section:    SectionEffect,
			Title:      "Stun",
			Definition: "Place a status on a creature; the next time it reaps, fights, or uses an Action, it exhausts and the stun is removed instead.",
			Body: `A stun is a status placed on a creature. A stunned creature must shake off the
stun before it can do anything else: the next time it is used to reap, fight,
or use an "Action:" ability, it is exhausted and the stun is removed instead of
that action happening. Its constant abilities and any effect that does not
require using it keep working while it is stunned. Stunning applies this status
to each creature the effect targets.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Unstun",
			Definition: "Remove the stun status from the chosen creatures, freeing them to act normally.",
			Body: `Unstunning a creature removes the stun status from each creature the effect
targets, freeing it to act normally instead of having to shake the stun off.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Trigger Another Card's Ability",
			Definition: "Resolve the abilities under another card's named trigger without using that card, for the player who reached for them.",
			Body: `Triggering another card's ability is not using that card. The card does not
exhaust, its use is not recorded, and nothing that watches for a card being
used fires — only the abilities printed under the named trigger resolve, and
they resolve for the player whose effect reached for them, as if that player
controlled the card.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Put a Card From Hand Under a Card",
			Definition: "The controller places a card from their hand under the resolving card, face up or face down.",
			Body: `PutUnderFromHand has the controller choose a card from their hand and place it
under the resolving card, face up or face down. Masterplan and Jargogle place
theirs facedown; Graft always places its card faceup. It does nothing with an
empty hand.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Play the Card Under a Card",
			Definition: "Play a card placed under the resolving card, putting the one played into context.",
			Body: `PlayCardUnder plays the card placed under the resolving card, putting the one
played in context (It) — Masterplan's and Jargogle's own "play the card under
me." With more than one card underneath, the controller chooses which; it does
nothing with none.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Graft a Card",
			Definition: "Move a card in play faceup under the resolving card, out of play, until that host leaves play.",
			Body: `Graft moves a target card in play faceup under the resolving card, out of play
(rulebook: Graft). The grafted card leaves play — firing its Leaves Play
abilities, not Destroyed — and waits under its new host until that host leaves
play. Spangler Box grafts a chosen creature onto itself.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Put the Cards Under a Card Into Play",
			Definition: "Put every card under the resolving card into play under its owner's control.",
			Body: `PutUnderIntoPlay puts every card placed under the resolving card into play
under its owner's control — Spangler Box's Destroyed ability returns the
creatures grafted onto it. It does nothing with nothing underneath.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Abduct",
			Definition: "Route an enemy card you take to your battleline, artifact line, or archives; any other destination sends it to its owner's matching zone.",
			Body: `Only three zones of yours may hold a card your opponent owns: your battleline,
your artifact line, and your archives. A card that would move to any other zone
of yours — your hand, your discard pile, your deck — goes to its owner's
matching zone instead. So an enemy creature abducted into your archives goes to
your opponent's hand when you take your archives up, and to their discard pile
if those archives are discarded.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Toll",
			Definition: "Æmber a card makes the opponent give to play an artifact or use an artifact's ability, going to the toll card's controller.",
			Body: `A toll is Æmber a card in play makes its controller's opponent give in order to
take an action with an artifact — playing an artifact, or using an artifact's
ability. The opponent cannot take the action unless they can pay the toll, and
the Æmber they give goes to the toll card's controller. (The mechanic keeps the
name Toll, but its printed text always reads "give", never "pay".)`,
		},
		{
			Section:    SectionEffect,
			Title:      "Cannot Be Used To",
			Definition: "Bar a creature from one way of using it — reap, fight, or Action — while every other way stays open.",
			Body: `A card that "cannot reap" is barred from one way of using it while every other
way stays open — Tireless Crocag fights and uses its Action: normally. That is
narrower than the timed, player-wide CannotUse/CannotFight restrictions in
effect_restrict.go, so it lives on the card definition rather than on state.`,
		},
	})
}

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
The amount lost is either a fixed Amount or a By loss of the pool (By: HalfRoundedDown,
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
			Title:      "Power Counter",
			Definition: "A permanent token that raises a creature's power by one (+1) or lowers it by one (-1) while it stays in play.",
			Body: `A +1 power counter is a permanent token placed on a creature that raises its
power by one for as long as it stays in play; a -1 power counter lowers it. A
creature can hold any number of counters, and they are shed when it leaves play.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Generic Counters",
			Definition: "Card-placed markers that do nothing on their own and matter only to the cards that read them.",
			Body: `A generic counter is a token placed on a card in play whose meaning is defined
entirely by the card that reads it. It does nothing on its own; a card can hold
any number of counters of a kind, and they are shed when it leaves play. The
doom counter is one: Wretched Doll destroys every creature carrying a doom
counter when there is one in play, and otherwise places a fresh one.`,
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
			Definition: "Reveal the top cards of a deck to both players so a following effect can inspect or play the top one, or route the revealed cards.",
			Body: `RevealTopOfDeck reveals the top Amount cards of a deck to both players and binds
the top one in context (ctx.It) so a following effect can inspect or play it (Chaos
Portal plays it when it is of the chosen house). Set ChooseWhoseDeck to have the
controller pick whose deck. Ordered routing steps then send chosen cards to a
destination — hand, archives, discard pile, or purge — and a final step either
reorders whatever the earlier steps left, or shuffles that deck (the revealed cards
are still in it). It reveals as many as remain when the deck
holds fewer than Amount, and does nothing on an empty deck.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Discard",
			Definition: "Dig through the top of your deck, discarding as you go, until you turn up a card the filters admit or the deck runs out.",
			Body: `DiscardUntil digs through the top of your deck, discarding as it goes,
until it turns up a card the filters admit or the deck runs out. With MayStop the
controller may stop before a match. Every discarded card is recorded on the context
and the matching card is left in context (ctx.It), so a following effect can act on
the found card or the whole discarded run — Sound the Horns and Invasion Portal pair
it with PutDiscardedIntoHand; Old Boomy archives the run with ArchiveDiscardedThisWay.`,
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
per-match state, dropped when the creature leaves play; RemainderOfPlayerTurn
also drops it at end of turn, while UntilThisLeavesPlay keeps it until the
creature leaves play.`,
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
			Title:      "Search",
			Definition: "Search one or more of your zones for cards matching a filter, reveal what you take, and put it into your hand or archives.",
			Body: `Search is the KeyForge "search" keyword: the controller looks through one or
more of their own zones — the deck, or the deck and discard pile — for cards
matching a filter (any card, a named card, a trait, a type, or a house), reveals
what they take, and moves it to their hand or archives. It can take one chosen
card or every match. It does not shuffle: a search is always followed by a
separate Shuffle, enforced by a card lint — except a search that puts its finds
on top of the deck, which shuffles between finding and placing so the finds land
on top of an already-shuffled deck (Digging Up the Monster).`,
		},
		{
			Section:    SectionEffect,
			Title:      "Shuffle",
			Definition: "Shuffle your deck, optionally folding whole zones or your cards in play into it first.",
			Body: `Shuffle shuffles the controller's deck. With no source it is the plain "shuffle
your deck" that follows a deck search (Orb of Wonder, Saurus Rex, Grumpus Tamer),
kept a standalone effect so a search never bundles its own shuffle. Zones folds
whole piles in first — the discard pile (Help from Future Self), the hand and
discard pile (Screaming Cave), the archives and discard pile. FromPlay folds every
friendly card in play and its upgrades in, then draws a card for each card
shuffled this way (Timequake).`,
		},
		{
			Section:    SectionEffect,
			Title:      "Swap Deck And Discard",
			Definition: "Exchange your deck with your discard pile, then shuffle the new deck.",
			Body: `SwapDeckAndDiscard exchanges the controller's deck with their discard pile and
shuffles the new deck — Reverse Time turns a spent deck back into a fresh one.
It differs from Shuffle{Zones: []Zone{Discard}} in that the old deck goes away
into the discard pile rather than surviving underneath it.`,
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
			Title:      "Enrage",
			Definition: "Place a status on a creature; while enraged, its controller must use it to fight on their turn if it is able to.",
			Body: `An enrage is a status placed on a creature. While a creature is enraged, its
controller must use it to fight on their turn whenever it is able to — it cannot
reap or use an "Action:" ability while there is an enemy creature it can fight. If
it cannot fight (nothing to fight, or an effect stops it), it is free to reap or
act. After a creature is used to fight, its enrage is removed; this holds even when
the defender is Elusive and takes no damage, because the fight still happened.
Reaping does not remove enrage. Enrage otherwise persists until an effect removes
it. Enraging applies this status to each creature the effect targets.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Ward",
			Definition: "Place a one-shot shield on a creature; it absorbs the next instance of damage or the next time the creature would leave play, then is spent.",
			Body: `A ward is a one-shot shield placed on a creature. The first time a warded
creature would be dealt damage or would leave play — destroyed, purged, returned
to hand, archived, put on or shuffled into a deck, or placed under a card — that
damage or removal is absorbed: the creature stays, and its ward is spent. Even a
single point of damage spends the whole ward. Ward intercepts every removal, even
the controller's own. It covers only damage and leaving play; it does not stop a
stun, an enrage, a capture, a change of control, or a loss of power. Warding
applies this status to each creature the effect targets.

Each removal is absorbed separately. A creature whose ward absorbs a destruction
and is then warded again absorbs the next removal with the new ward.

A ward gained after a creature is already destroyed does not undo the
destruction. That creature still reaches its discard pile once the "Destroyed:"
abilities finish. The new ward does absorb any other removal in the meantime, so
a "Destroyed:" ability that tries to purge the creature is absorbed.`,
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
narrower than the timed, player-wide Restrict restrictions in
effect_restrict.go, so it lives on the card definition rather than on state.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Put a Card into Another Zone",
			Definition: "Move a card from where it is to a named zone — a hand, a deck, the archives, or play — without playing it.",
			Body: `Moving a card between zones is one mechanism. Every such effect says three
things: where the cards come from, which of them move, and where they go. The
destinations are a hand, the top or the bottom of a deck, the archives, the
discard pile, the purge pile, and play itself.

Which cards move is a selection: the card the controller chooses, every card
that matches, a fixed number, up to a number, a random card, or the card a
preceding effect left in context. Filters narrow the choice by type, trait,
house, or name. An effect that can move several cards moves them one at a time.

A card moving out of play sheds everything it built up there — its damage, its
spent armor, the Æmber on it, and its upgrades. This is how a "Destroyed:"
ability can save its own creature: the creature leaves for the named zone as it
is destroyed, so it never reaches the discard pile. When several cards go to the
top of a deck at once, their controller chooses the order they stack. A card
that has already left play is not moved a second time.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Put a Card into Play",
			Definition: "Put a card into play without playing it: no bonus icons, no Play abilities, only enters-play reactions.",
			Body: `Putting a card into play is not playing it. The card's bonus icons do not
resolve and its "Play:" abilities do not fire; only reactions to a card entering
play see it. That is what lets one player put an opponent's card into play
without making that player's decisions for them.

A creature enters play exhausted unless the effect says it enters ready, and its
controller chooses the flank it enters on. The card enters under its owner's
control unless the effect gives it to the resolving player; ownership never
changes either way. The card can come from anywhere the effect names — a hand, a
discard pile, a deck.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Reveal a Card",
			Definition: "Show a card to both players, turning hidden information public, and leave it in context for what follows.",
			Body: `Revealing a card shows it to both players and records it in the log. The card
does not move: a card revealed from a hand stays in that hand.

A card is revealed so that what follows can be trusted — the cards you draw for,
the house of the card you turned up, the text box you lend. The revealed card is
left in context, so a following effect can act on that card, and a count can
measure how many cards were revealed.

An effect reveals a whole hand, a card its owner chooses, or a random card. A
reveal narrowed to a house lets its player pick which of those cards to show,
one at a time, and "any number" includes none. An empty hand reveals nothing.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Name a Card",
			Definition: "Name any card in the game; while the naming card stays in play, no card of that name may be played by either player.",
			Body: `Naming a card has its controller name any card in the game, which is recorded on
the card that named it. While that card stays in play, no card of the named name
may be played, by either player. The bar lifts when the naming card leaves play.

The names offered are every card in the game, not only the cards in this match —
offering only the cards in the match would show a player their opponent's deck.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Attach an Upgrade",
			Definition: "Move an upgrade in play onto another creature; it stays in play and keeps what it grants.",
			Body: `Attaching moves an upgrade already in play onto another creature. The upgrade
stays in play throughout — it does not leave and return — so nothing that
watches a card entering or leaving play fires. Its host creature changes, and
everything the upgrade grants now applies to the new host.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Battleline Position",
			Definition: "Where a creature sits among its controller's creatures, and the effects that move it there.",
			Body: `A player's creatures sit in a battleline, in order. A creature's position
matters: its neighbors are the creatures immediately to its left and right, and
the creatures at the two ends are on a flank.

An effect can change a position. It can move a creature to a flank of its own
battleline, move it anywhere within that battleline, swap two creatures'
positions, or let a player reorder a whole battleline by swapping pairs until
they stop. Only positions move — a creature keeps its damage, its Æmber, its
upgrades, and its exhaustion, and it does not leave play.

Card text also reads positions: an effect can reach a creature's neighbors, the
creature to one side of it, the creature in the center of a battleline, or the
creatures that were a card's neighbors before it left its position. A count can
measure how many neighbors a creature has, or their combined power.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Become a Creature",
			Definition: "Turn a card in play into a creature on a flank of its controller's battleline.",
			Body: `An effect can turn a card in play into a creature — an artifact becoming a
creature is the usual case. The card moves to a flank of its controller's
battleline, the effect's controller choosing the flank, and reads as a creature
from then on.

The card keeps what it already had: its exhaustion, the Æmber on it, and its
power counters, so counters placed on it before it became a creature now count
toward its power. The change lasts until the card leaves play, unless the effect
names the rest of the turn, in which case the card reverts at the end of the
turn and keeps its counters.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Cancel a Fight",
			Definition: "Make the fight in progress not happen; the attacker was still used and stays exhausted.",
			Body: `Cancelling a fight makes the fight in progress not occur. It is a "Before
Fight:" effect, so it resolves while the fight is still pending. No damage is
dealt either way, and Assault, Hazardous, and "Fight:" abilities are all
skipped.

The attacker was still used to fight, so it stays exhausted and anything
watching for a creature being used still fired.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Cannot Be Dealt Damage",
			Definition: "For a named window, a creature cannot be dealt damage at all; damage aimed at it does nothing.",
			Body: `A creature that cannot be dealt damage takes none for the window the effect
names. This is checked as the damage is dealt, so it covers every source —
fight damage, an ability, a keyword — and it covers damage from its own
controller.

It is not armor: nothing is absorbed and nothing is spent. Damage already on the
creature stays there.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Fused Triggers",
			Definition: "For the rest of the turn, each friendly creature's abilities under one trigger also fire under another.",
			Body: `An effect can fuse two triggers for the remainder of its controller's turn: each
friendly creature's abilities printed under the first trigger also fire under the
second, and its abilities under the second also fire under the first. A creature
used to reap then resolves its "Fight:" abilities as well, and a creature used to
fight resolves its "Reap:" abilities.

The fusion is a rule held by the turn, not by the creatures, so a creature played
after it resolves is fused too. It lifts at the ready step.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Gain a Trait",
			Definition: "Give a creature a trait until the start of its controller's next turn.",
			Body: `An effect can give a creature a trait it does not have printed. The creature
counts as having that trait for everything that reads traits, its own card text
included.

The grant lasts until the start of the granting player's next turn, lifting
before any start-of-turn ability resolves. It therefore survives the opponent's
whole turn, where an enemy card may read the granted trait.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Gain an Ability",
			Definition: "Give a creature a triggered ability for a window, as if that ability were printed on it.",
			Body: `An effect can give a creature a triggered ability — "Reap: Draw a card" — for a
window. Only that creature's own trigger fires it, and it resolves for that
creature's controller.

By default the grant lasts the remainder of the granting player's turn. A grant
made until the start of the granting player's next turn holds through the
opponent's turn, so a creature the opponent controls can fire it on their own
turn. Card text can name the card that granted the ability, so a granted ability
can reach back to its grantor.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Gain and Lose Keywords",
			Definition: "Grant a keyword to creatures for a window, or strip a keyword from every creature for the turn.",
			Body: `An effect can grant a creature one or more keywords. The creature has them for
the window the effect names, and they work exactly as printed keywords do. A
defensive keyword is usually granted until the start of the granting player's
next turn, so it holds through the opponent's turn; an offensive one is usually
granted for the remainder of the turn only.

An effect can also take a keyword away. A loss reaches every creature in play
for the remainder of the turn, and it is held by the turn rather than by the
creatures, so a creature played later in that turn loses the keyword too.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Power and Armor",
			Definition: "Raise, set, or copy a creature's power and armor; power is recomputed as those changes come and go.",
			Body: `A creature's power and armor are computed, not stored. Its printed values, its
power counters, its upgrades, the constant abilities in play, and any temporary
bonus are all read together each time the game asks.

An effect can add power or armor for the remainder of the turn, which the ready
step clears. It can mask power and armor to fixed values for a window, which
ignores every other modifier while it holds and reveals the stored values again
when it lifts. And it can make one creature copy another's printed power, armor,
keywords, and traits until it leaves play. A count can read the power of a
chosen creature, which reads its current power, not its printed one.

Because power is recomputed, a creature can become destroyed with nothing
touching it — the card buffing it leaves play, a counter comes off, its text is
blanked. Whenever a change can lower a creature's power, damage at or above its
new power destroys it at once.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Redistribute Damage",
			Definition: "Take all the damage off one player's creatures and place it back among that player's creatures however the controller chooses.",
			Body: `Redistributing damage takes every damage token off one player's creatures and
places them back, one at a time, on creatures of that same player. The
controller of the effect chooses the player, then chooses where each token
lands; they may also decline to redistribute at all.

The total is conserved: no damage is healed and none is added. Tokens may be
piled on one creature past its power, which destroys it.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Replacement",
			Definition: "Change the outcome of an event before it happens, so the replaced outcome never occurs.",
			Body: `A replacement changes what an event does before the event happens. The original
outcome never occurs — reaping steals Æmber instead of gaining it, Æmber added to
a pool is captured before it lands, the next Tactic resolved returns to its
controller's hand instead of their discard pile.

A replacement is either a lasting effect for the remainder of the controller's
turn, or a standing rule printed on a card that holds while that card is in
play. Every replaced outcome is narrated by naming the card that caused it, so a
player can always see why an event did something other than what it says.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Take Control",
			Definition: "Become the controller of a card in play; ownership never changes, so the card still returns to its owner's piles.",
			Body: `Taking control of a card moves it to the taking player's play area and makes
that player its controller — a creature joins their battleline, an artifact joins
their artifact row. The card does not leave play, so nothing that watches for a
card leaving or entering play fires (TestTakeControlIsNotLeavingPlay). An effect
can also give a card away, handing a card the controller has to their opponent.

Ownership never changes. The card still goes to its owner's discard pile, hand,
or deck whenever it leaves play.

Control ends when the window the effect named ends: when the card whose ability
took control leaves play, or when the seized card itself leaves play. Takes
stack, so when the newest one ends, the previous one — if still in force — takes
over, and only when the last ends does the card go back to its owner
(TestControlStackIsLIFO).

A creature's power is read under its new controller, so a creature that loses a
buff by changing sides can be destroyed the moment control changes.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Text Box",
			Definition: "A card's printed traits, keywords, and abilities, which an effect can blank, copy, or lend.",
			Body: `A card's text box is its printed traits, keywords, and abilities. Its power,
armor, house, and type are not in the text box.

Blanking a text box makes it be ignored: the creature keeps its stats and its
traits, but its printed keywords, abilities, and constant grants do nothing
(TestBlankEnemyTextSuppressesAbilitiesAndConstants). A creature whose blanked
text was holding its own power up can be destroyed at once by the blanking
(TestBlankEnemyTextSettlesLethal).

An effect can instead give a creature another card's text box, as if that text
were printed on it. The source can be a creature in play, or a creature revealed
from the controller's hand, which stays in hand. The gain lasts until the
recipient leaves play, or only for the turn when the effect says so.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Move Æmber",
			Definition: "Take Æmber off a card and put it somewhere else — into a pool, onto another card, or back to the common supply.",
			Body: `Moving Æmber takes Æmber that is sitting on a card and puts it somewhere else:
into a player's pool, onto another card, or back to the common supply. Only a
card that carries Æmber can be the source, so an effect with no such card moves
nothing. A source holding less than the amount named gives up all it has.

The move is not a steal and not a capture: the Æmber is already out of both
pools, and moving it only changes where it waits. An effect can move a fixed
amount, a share of what the source carries, or all of it. A following effect can
turn on whether any Æmber actually moved.

"Your opponent gives you Æmber" is the pool-to-pool form: Æmber changes hands
without being stolen, capped at what the opponent's pool holds.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Key Cost",
			Definition: "A change to what a player's keys cost to forge; raises and lowers sum, and the total never falls below zero.",
			Body: `A key costs 6 Æmber by default. An effect can raise that cost or lower it, for
one player or for both. Raises and lowers in force at the same time sum, and the
total is floored at zero before a key is forged.

The cost is read at the moment a player forges, not when the effect resolves, so
a raise that scales with the board ("+1 for each Dis creature in play") counts
what is in play at the forge. A change that should hold for as long as its card
stays in play is printed on that card instead of resolved as an effect, and the
forge reads it the same way.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Refill Hand",
			Definition: "A player refills their hand now, as if their turn were ending, honoring their chains.",
			Body: `Refilling a hand draws a player back up to a full hand at once, as if it were
their draw step — honoring their chains and any effect that changes how many
cards they draw. It can refill both players' hands, the active player's first.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Use a Card",
			Definition: "Use a card as its controller would: a creature reaps, fights, or takes its Action, and an artifact takes its Action.",
			Body: `An effect that uses a card makes that card do what using it does: a creature
reaps, fights, or takes its "Action:", and an artifact takes its "Action:". The
card exhausts, its triggered abilities fire, and everything that watches for a
card being used sees it — this is a real use, not a copy of one. Only a ready,
usable card can be used this way.

When an effect uses several cards, they are used one at a time and each use
resolves fully before the next card is chosen, so a card that exhausted itself
is no longer a candidate. Counts and conditions that read "used this turn" count
these uses like any other.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Play or Use a Card",
			Definition: "Free a player from the active-house restriction, either as a permission held for the turn or as one play-or-use that happens now.",
			Body: `A player may normally play and use only cards of their active house. An effect
can free them from that: it names which cards it frees (a house, any house, every
house but one, or a trait), which verbs it frees (play, use, or fight), which
card types, and how many cards.

The permission form lasts the current turn and lifts at the ready step. The
immediate form spends itself at once instead: the controller plays or uses one
matching card right now, and nothing is held open. Freeing a card from the house
restriction changes nothing else about playing or using it.`,
		},
		{
			Section:    SectionEffect,
			Title:      "Resolve Bonus Icons",
			Definition: "Resolve the bonus icons printed on a card as if it had just been played, without running the rest of its text.",
			Body: `Resolving a card's bonus icons resolves each icon printed on it in turn, as if
its controller had just played that card. None of the card's other text runs and
nothing that watches for a card being played fires.

The icons are read from the card itself, so a card that was revealed, discarded,
or purged still resolves them; the card need not be in play. An effect that bars
bonus icons, or substitutes one icon for another, applies here too. An effect
can also arm an extra resolution: the next card its controller plays this turn
resolves each of its icons one additional time.`,
		},
		{
			Section:    SectionEffect,
			Title:      "End the Turn",
			Definition: "Stop the active player's turn immediately, with no ready step, no draw, and no end-of-turn abilities.",
			Body: `Ending the turn stops it the moment the effect resolves. The turn simply stops:
there is no ready step, no draw step, and no "at the end of your turn" abilities.
The active player's cards stay exhausted and their hand is not refilled.

This is not the Omega keyword. Omega ends the main phase but runs the rest of the
turn out, so its player still readies, draws, and resolves their end-of-turn
abilities.`,
		},
	})
}

# Vactrol

A KeyForge-style card game engine and (procedural) deck generator. This glossary
fixes the vocabulary shared across the engine, the card database, and deck
generation. It is a glossary only — no implementation detail.

## The game

The engine implements a KeyForge-style game; **`docs/rulebook.md` is the
authoritative, comprehensive source** for how each mechanic works. The entries
here only fix the shared vocabulary — especially where Vactrol diverges from
KeyForge.

**Æmber**:
The game's currency. A player gains Æmber into a pool and spends it to forge keys;
Æmber can also sit on cards (captured, exalted, or as a printed bonus).

**Key** / **Forge**:
Forging a key spends Æmber from the pool; the first player to forge three keys
wins.

**Key cost**:
The Æmber required to forge a key (6 by default), which cards may raise or lower.

**Check**:
The state of holding enough Æmber to afford a key; announced in the tabletop
game so the opponent knows a key is imminent. The client shows it as a glow on
the Æmber count rather than a callout.

**Key colour**:
The colour of a forged key — Red, Blue, or Yellow, the same colours as KeyForge.
Currently cosmetic; planned to matter to house Skyborn (set 8).
_Avoid_: Key type.

**Chains**:
A penalty that reduces how far a player refills their hand, shed over successive
turns.

**Active house**:
The single House a player chooses at the start of their turn; only cards of the
active House may be played or used that turn, barring specific grants.

**Phase**:
One of the eight ordered parts of a turn: start of turn, forge, choose a house,
archives, play, ready, draw, end of turn. Each phase is entered and left
explicitly, and abilities that resolve "at the start of your turn" or "at the end
of your turn" belong to the phase of that name. An effect may end the current
phase early, skipping whatever remained in it. KeyForge's rulebook calls these
divisions "steps"; Vactrol calls them phases everywhere — engine, rulebook, card
text, and game log.
_Avoid_: step, turn step, main phase.

**Reveal**:
To make a card in a hidden zone publicly known. Revealing is the only way a card
in a hidden zone is ever named — to the opponent, and in the game log.

**Reap**:
Using a ready creature to gain 1 Æmber, exhausting it.

**Fight**:
Using a ready creature to attack an enemy creature, exhausting it.

**Tactic** (vs the **Action** ability):
A **Tactic** is the one-shot card type — KeyForge's "action card", renamed so the
word "Action" is free for the "Action:" ability, which is used directly from a
card already in play. See ADR 0009.
_Avoid_: calling the card type "Action".

**Omni**:
An ability usable on any turn, not only the active House's. Vactrol has no Omni
trigger — an Omni card is authored as the **Versatile** keyword plus an Action
ability (ADR 0009).

**Battleline**:
The ordered row of a player's creatures in play.

**Flank**:
The leftmost or rightmost creature in a battleline.

**Archives**:
A set-aside zone a player can later draw back into hand; may hold either player's
cards.

**Purge**:
To set a card aside out of the game; purged cards never return.

**Capture** / **Steal** / **Exalt**:
The Æmber verbs — capture moves Æmber onto a creature (off a player's pool), steal
takes Æmber from the opponent's pool into yours, exalt places Æmber from the
supply onto a card.

**Toll**:
Æmber the opponent must pay a card's controller in order to play or use an
artifact.

**Rule of 6**:
A player cannot play and/or use the same card — or other copies of that card _by
name_ — more than six times during a given turn. The count is keyed by card name
across the whole turn, **regardless of who owns or controls the card**: six plays
or uses of "Bumpsy" in a turn, not six per copy and not six per player. Only the
active player can play or use cards during a turn, so there is no opponent
contribution to track; keying by name rather than by player is what stops a stolen
or seized copy from buying a seventh use. A self-repeating effect is bounded by
the same six: Bait and Switch resolves its steal once and repeats it at most five
more times, so one copy steals six Æmber at most. The count resets when the turn
does. It exists to stop an unbounded loop from hanging the game.

**Invulnerable**:
A KeyForge keyword: an invulnerable creature cannot be destroyed or dealt damage.
The engine does not model the full keyword yet — only the damage-prevention half,
as the `DamageImmune` flag set by the `PreventDamage` effect ("cannot be dealt
damage this turn"; Shield of Justice, Protectrix). Those cards prevent damage
only, so they are not truly invulnerable.

**Hidden zone** / **Public zone**:
A zone whose contents are not known to both players (deck, hand, archives) versus
one whose contents are (play, discard, purged). The distinction decides whether a
card can be named.

**Upgrade**:
A card type attached to a creature rather than played onto the battleline itself;
it grants its host a bonus (power, armor, a keyword) for as long as it stays
attached, and leaves play alongside its host or on its own.

**Under** / **Under-card**:
A card placed under another card rather than played into a zone of its own —
Masterplan and Jargogle bury a card from hand this way, Graft always faceup.
Unlike an Upgrade, an under-card is out of play: it does not fight, reap, or
count toward anything power does. It may be placed face up or facedown; a
facedown under-card is exactly as hidden as a card in hand.
_Avoid_: beneath.

**Peek** / **Peekable**:
Looking at a facedown under-card without revealing it to the opponent. Only the
controller of the host a card is placed under may peek at it.

**Constant ability**:
An ability with no boldfaced trigger, which applies continuously while its card
stays in play — the power and armor bonuses one card hands its neighbors, a
keyword or quoted ability it grants every creature, a restriction it imposes.
It is not a triggered ability and applying it does not exhaust the card. This is
the only continuous-effect concept in the engine; there is no separate notion of
an "aura".

**Take control** / **latest ability wins**:
Taking control moves a card to your play area and makes you its controller;
ownership never changes and still decides where the card returns when it leaves
play. When two effects change the same thing on a card (its controller, its
house), the most recently applied one wins. A `Forever` control (Sneklifter's
seized artifact) lasts the rest of the game; if a later, timed effect overrides
it, the `Forever` effect governs again once that timed effect expires.

## The game log

**Game log**:
The running account of what has happened in a match. It records outcomes — the
state changes that actually occurred — never a card's printed text, never what an
effect attempted, and never a hint about what a player may do next. It is
identical for both players and names no card held in a hidden zone, so a card is
named only once it has been revealed or has reached a public zone.
_Avoid_: chat, message, narration.

**Log entry**:
One line of the game log: one thing that happened, under one attribution. An
entry knows how to render itself, the same way an effect renders its own card
text; the two share a vocabulary of phrasings but neither is derived from the
other (ADR 0011).
_Avoid_: message, narration record.

**Frame**:
The scope a log entry was emitted inside, carrying who acted, which card the
ability came from, which ability category it was, and which card granted it.
Frames nest — playing a creature opens a frame, its Play ability opens a child
frame — and entries inherit their attribution from the frame they sit in. The
client groups a top-level frame's entries into one visual bubble.
_Avoid_: use, useId, bubble (as a code term), log group, log mark.

## Cards and sets

**House**:
A card's intrinsic allegiance (Brobnar, Dis, Logos, …). A Set groups its cards
into a small, per-set number of Houses (7 in a classic set, but adjustable).
_Avoid_: Faction (aspirational rename for later; today it is House).

**Set**:
The pool of cards a Deck is generated from, defined by the cards declared in one
`internal/cards/sets/<slug>` package. Membership is derived from the card
database, never from Provenance.
_Avoid_: Expansion, edition.

**Native set**:
The Set in whose package a card is declared. A card's home.

**Reprint**:
A card that appears in a Set other than its Native set. A card may belong to
several Sets at once (each new Set is ~half reprints).

**Provenance**:
Historical metadata tagging a card with the original KeyForge card it was derived
from (source set + collector number). It exists only to track which original card
each implementation is based on, so the author can confirm every original KeyForge
card is eventually covered. It is never consulted by the engine or by deck
generation, and a card's behavior never depends on it.

## Deck generation

**Deck**:
The 36-card result of generating from one Set with one seed: 3 House pods of 12
cards. Reproducible only within a single version of its Set's pool.

**House pod**:
One of a Deck's 3 Houses together with its 12 Slots. A Deck is exactly three House
pods.
_Avoid_: House deck, deck-House.

**Slot**:
One of a House pod's 12 positions (36 per Deck). Carries an intrinsic Rarity and
independent provenance flags (Maverick, Legacy), the card chosen to fill it, and
any enhancements/distortions applied to it.

**Rarity**:
A card's intrinsic scarcity, which drives how often it is drawn into a Slot:
Common, Uncommon, Rare, Special, Reference.

**Special**:
A Rarity of houseless, very-rare cards. A Special has no House until it fills a
Slot, at which point it is stamped with that Slot's House.

**Reference card**:
An optional, Reference-rarity card filling a Set's Reference slot to support a
set-specific mechanic (e.g. a prophecy card). Chosen by the Set's own rules,
typically by card type. Most Sets have none.
_Avoid_: Fixed, token.

**Maverick**:
A Slot property: the card originates from a House other than the Slot's House
(same Set). For play it adopts the Slot's House (it is rehoused); the Maverick
flag records only that the card is a maverick, not its origin House.

**Legacy**:
A Slot property: its card is drawn from a Set other than the Deck's Set, same
House. Combines freely with Maverick ("Legacy Maverick").

**Legacy pool**:
For a Deck of Set X and House H, the cards of House H belonging to any Set ≠ X.

**Legacy house**:
A very rare House pod whose 12 Slots draw from the same House in a different Set
(the House matches the Set, but its pool does not).

**Maverick house**:
An even rarer House pod that is a House not present in the Deck's Set at all, drawn
wholesale from another Set.

## Card modification

**Enhancement**:
A bonus (Æmber, capture, draw, …) redistributed during the deck-wide finishing
pass. A card may be an enhancement **source** — it declares bonuses it contributes
to the deck's pool — and any card may be a **recipient**, gaining landed bonuses
(up to a cap of 5 per card, including printed bonus icons). A source does not keep
what it contributes. Because bonuses land per Slot, two Slots of the same card can
differ.

**Distortion**:
A generation-time transfer, resolved in the finishing pass, that weakens one
card's property to strengthen the same property on another card, capped between -2
and +2. Landed per Slot, like an Enhancement. Not yet in the engine.

## Scoring

**Bot**:
Vactrol's automated player: a procedural Monte Carlo Tree Search (MCTS) engine
that clones game states and self-plays. It is heuristics-and-search based, not a
generative or LLM model. It shares the scoring model with Deck Rating.
_Avoid_: AI (ambiguous with generative LLMs).

**Rating**:
A Deck's aggregate score, derived from the scoring model. Optional band-targeted
generation aims for a chosen Rating range. Because generation is deterministic
chaos, band-targeting is rejection sampling from a Set's intrinsic Rating
distribution, so that distribution (mean, spread) is exposed and reported.

**Scoring model**:
A learned table of per-card capability vectors and synergy tags, refined by MCTS
self-play (optionally seeded from a card's effect analysis). Consumed both by Deck
Rating and by the Bot's in-game play decisions, so the two share one notion of
card value.

**Synergy tag**:
A named axis a card _provides_ or _consumes_ with a weight. Deck synergy is the
weighted match of providers to consumers across the Deck; antisynergy is a
negative weight.

**Connection**:
A relationship where selecting one card pulls additional cards into the same
House. Parameterized by a pool of candidates, a count (fixed or variable), and
whether duplicates are allowed: Plague Rat pulls a variable number of copies from
a one-card pool (duplicates); Timetraveller pulls exactly one specific partner; a
sin slot pulls a variable number from the seven sins without duplicates;
Groundbreaking Discovery pulls exactly one of each of its three partners.

## Templates

**Card template**:
A parameterized card that is not itself playable and occupies a single pool
entry. Deck generation binds its free parameters to produce a concrete,
engine-ready card, including its name and boilerplate. Parameters may be
combinatorial (a mutant composed from two Houses, 7 home × 6 other = 42 outcomes),
enumerable (Master of 1/2/3 as one entry), or contextual (a self-house card whose
text names its own House).
_Avoid_: Generator, factory.

**Materialize**:
To resolve a Slot into a concrete, engine-ready card at generation time — binding
any template parameters, rehousing a Maverick to the Slot's House, and binding
Home-house references. The engine only ever sees materialized cards.

**Home house**:
A template parameter for a card that references its own House in its text or
effect. Bound at generation to the Slot's House, so a Maverick reads and plays
correctly.

## The client

The names of the screen's regions, so a request or a bug report can point at one.
These are the names the browser client's regions go by.

**Board area**:
Everything but the Sidebar: the two Player bars, the Play zone between them, and
the Hand row at the bottom.

**Player bar**:
One player's summary strip — their name, Æmber, key cost, keys, deck Houses,
chains, and Zone counts. The opponent's is at the top of the Board area, the
active player's below the Play zone.
_Avoid_: score pill (the older name, still the CSS class), status bar.

**Zone counts**:
The out-of-play card counts at the right end of a Player bar — hand, deck,
discard, archives, purge — which open the Zone viewer.

**House strip**:
The three deck Houses shown in a Player bar, with the non-active ones lowlighted.

**Play zone**:
The four Board rows of cards in play, artifacts outside and battlelines inside,
split by the Midline. It is also the drop target for a card dragged from hand.

**Board row**:
One zone of one player rendered as a Row label and a strip of cards: an
_artifact row_, a _battleline row_, or the _Hand row_.

**Row label**:
The rotated caption on the left of a Board row — owner, zone, and count. It drops
the zone word for its icon, then the owner's name, as the row gets shorter.

**Midline**:
The dashed rule between the two battlelines; the board is a mirror about it.

**Sidebar**:
The right-hand column: Brand bar, Game log, Turn HUD, and the Control dock. It
can be collapsed to give the Board area the whole window, which sets the Control
dock loose to float over the board.

**Control dock**:
Everything the player answers with, as one panel: the Prompt, the Action bar, the
House picker, the Flank buttons, and End turn. It docks into the Sidebar while
the Sidebar is open and floats over the Board area once it is hidden, so hiding
the Sidebar costs the player the Game log but never the game.

**Brand bar**:
The title row at the top of the Sidebar, with the navigation buttons: undo,
redo, manual mode, new game, and hiding the Sidebar.

**Game log**:
The running record of the match. A _log line_ is one entry; a _log group_ is the
lines of a single action, drawn as one bubble; a _card mention_ is a card name in
a line, hoverable for its Card preview.

**Turn HUD**:
The at-a-glance state of the current turn above the Prompt — whose turn it is,
the turn number, the step they are in, and their active House.

**Prompt**:
A question the engine is waiting on, shown with the card that asked it (the
_prompt source_) and its answers as option buttons.

**Action bar**:
The buttons for what the selected card can do right now — play, discard, reap,
use, fight, end turn.

**House picker**:
The choice of which House to call for the turn, offered at the turn's start.
Distinct from the House strip, which only reports the three a deck has.

**Flank buttons**:
The arrow-headed buttons that place a played creature on the left or right flank.

**Manual panel**:
The Sidebar's manual controls, available only in manual mode: stat steppers,
moves between zones, and the Card picker.

**Card picker**:
The searchable list of every card in the database, for putting one into play or
hand in manual mode.

**Zone viewer**:
The overlay listing a player's out-of-play zones as rows of cards.

**Card preview**:
The enlarged face of the card under the cursor, whether on the board, in hand, or
named in the Game log.

**Result panel**:
The end-of-game result, shown in the Action bar's place once a player has forged
their third key.

**Status banner**:
The transient message, usually a rejected play, that fades in above the Action
bar.

**Flash** / **Flight**:
The two feedback animations. A _flash_ pulses a card or counter that just
changed; a _flight_ is a card that left play arcing into the zone it went to.

**Tip**:
The small label a bare icon shows on hover.

**Style gallery**:
The page at `/style` showing every piece of the client's visual vocabulary at
once — color tokens, icons, fonts, card faces, a Player bar, and the animations.
It is a development surface, served only by `mage web`, and is not part of a
game. Its regions:

- **Specimen**: one displayed example, captioned with what selected it.
- **House grid**: the specimens laid out as House by Card type, whose gaps show
  which combinations the implemented sets have no card for.
- **Font compare**: the strip that renders one specimen once per loaded font, to
  choose between faces for the same House.
- **Style header**: the sticky controls for the fonts in use.

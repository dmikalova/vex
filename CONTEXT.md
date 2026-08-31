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

**Invulnerable**:
A creature status that prevents any damage from being dealt to it (Shield of
Justice, Protectrix). It lasts until end of turn. In the engine this is the
`Invulnerable` flag on a card's state.
_Avoid_: damage-immune, damage immunity.

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

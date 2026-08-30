# Vactrol

A KeyForge-style card game engine and (procedural) deck generator. This glossary
fixes the vocabulary shared across the engine, the card database, and deck
generation. It is a glossary only — no implementation detail.

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
from (source set + collector number), used only for the author's coverage
tracking. Irrelevant to deck generation.

## Deck generation

**Deck**:
The 36-card result of generating from one Set with one seed: 3 Houses × 12 cards.
Reproducible only within a single version of its Set's pool.

**Slot**:
One of a Deck's 36 positions. Carries an intrinsic Rarity and independent
provenance flags (Maverick, Legacy), the card chosen to fill it, and any
enhancements/distortions applied to it.

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

## Card modification

**Enhancement**:
A bonus (Æmber, capture, draw, …) assigned to a card during a post-generation
pass, up to a cap of 5 bonuses per card including printed bonus icons.

**Distortion**:
A generation-time transfer that weakens one card's property to strengthen the
same property on another card, capped between -2 and +2. Not yet in the engine.

**Connection**:
A relationship where selecting one card pulls additional specific cards into the
Deck (e.g. a card that pulls its partner, a card that pulls multiple copies, or a
"sin" slot that pulls a fixed set).

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

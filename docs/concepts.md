# Concepts

This document outlines the core concepts of Vactrol, a KeyForge-style card game.

## Core Concepts

- A set is a collection of about 400 cards.
  - A set has 7 factions, each of which has a playstyle
- A deck is a collection of 36 cards, drawn from a set using procedural generation.
  - A deck has 3 factions with 12 cards in a deck
  - The rarities are:
    - Common - about 18 per deck - these establish the flavor of the set and faction
    - Uncommon - about 12 per deck - these give specific flavor to the deck
    - Rare - about 6 per deck - these are cards that do unique things
    - Legacy - about 6 per deck - these are cards from other sets
    - Maverick - about 3 per deck - these are cards from other factions
    - Legacy Mavericks - the odds of this are maverick crossed with legacy
- A card has a faction, name, rarity, type, traits, keywords, and abilities
- Factions:
  - Brobner: Board presence and large fighters that benefit from fighting
  - Dis: Destruction of creatures on both sides for benefit
  - Ekwidon: Exchange this for that, in their favor
  - Enlightened: Builds a board presence that works up to a big payoff and must be disrupted
  - Geistoid: Use the discard pile as a resource
  - Logos: Efficiency and card draw
  - Keyraken: Centered on large monsters that everything else revolves around
  - Mars: Insular benefits and synergy, at the expensive of friendly and enemy non-Martians
  - Ouboros: Stay exhausted for benefit and hoard aember
  - Redemption: Redeem the others houses - suck distortion
  - Sanctum: Protect the board and neighbors
  - Saurian: Risk reward by putting aember on the board for benefit
  - Shadows: Small but stealy
  - Skyborn: Care about board placement for benefit
  - Star Alliance: Cooperate with other factions for benefit
  - Unfathomable: Disrupt opponent's hand and aember pool
  - Untamed: Aember rush
- Keywords:
  - Skirmish
  - Taunt
  - Poison
  - Hazardous
  - Haunted
  - Ward
  - Enrage
  - Armor
  - Elusive
  - Treachery
  - Splash
  - Invulnerable
  - Versatile
- Potential Keywords:
  - Ignore Elusive (MTG: reach)
  - First Strike
  - Defender
  - Haste
  - Change control
  - Abduct
  - Scry/Surveil/Mill - deals with top of deck
- Actions:
  - Exhaust
  - Ready
  - Capture
  - Steal
  - Destroy
  - Draw
  - Exalt
- Zones:
  - Play Area
    - Battleline
    - Artifact line
  - Hand
  - Deck
  - Archives
  - Discard
  - Card under card
  - Attached upgrades
- Counters:
  - Aember
  - Damage
  - +1/-1 power
  - stun
  - ward
  - enrage
  - +1/-1 armor
  - generic counters
- Abilities:
  - Reap:
  - Fight:
  - Play:
  - Destroyed:
  - Scrap:
  - Action:
  - Omni:

## Differences from KeyForge

- A refocus on the board over actions
- Cards with tradeoffs and being situational rather than just being good
- More upgrades can help boards be more dynamic and flexible
- Enhancements at this point are pretty flat. Instead, I want each card to know its traits (eg deals damage, gives ward, etc), and then have a distortion system where some cards give distortions to the deck, so that one card would get -1 of its trait, but that same trait would then be given +1 on a different card.
  - Maybe some distortions can change the traits being changed.
  - Distortions would be capped from +2 to -2
  - When a distortion takes so much that a trait would hit 0 or go negative, it turns into a softer trait:
    - 3 steal <-> 2 steal <-> 1 steal <-> capture 2 <-> capture 1
- Minimize simultaneous effects - things happen one at a time, reflecting the physical reality of the game.
- Not sure what to do with bonus icons. Distortion covers a lot of their ground, but having bonus aember still seems like a good thing.
- Action type cards, and the Action: ability need to be renamed so they're not the same word.

## Keeping from KeyForge

- Decide everything on your turn without interrupts

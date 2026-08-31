# TODOs & design notes

A grab-bag of things to build and open design questions. Canonical vocabulary is
in [../CONTEXT.md](../CONTEXT.md), the rules in [rulebook.md](rulebook.md), and the
long-term vision in [roadmap.md](roadmap.md).

## UI / gameplay to build

- Drag and drop to flank or deploy
- When sidebar is collapsed, can still access button actions
- Display multiple houses
- Turn steps
- Start of game setup
- More keyword icons - how much is too much?
- Elusive. Skirmish. Why. Periods.
- Why so much card.Sentence - can we just do this at the right place at the right time.
- The logging prints what the card does, but not necessarily what happened. Fix that up.

## Design directions (deliberate divergences from KeyForge)

- Refocus on the board over one-shot actions.
- Cards with tradeoffs / situational value rather than being strictly good.
- Lean on upgrades to make boards more dynamic and flexible.
- Distortion system instead of flat enhancements — see the Enhancement and
  Distortion entries in [../CONTEXT.md](../CONTEXT.md) and
  [deck-generation.md](deck-generation.md).
- Minimize simultaneous effects — resolve one at a time, matching the physical
  game.
- Open question: what to do with bonus icons once distortions exist (bonus Æmber
  still seems worth keeping).

## Design principles (kept from KeyForge)

- You decide everything on your own turn — no interrupts.

## Houses & intended playstyles

Brobnar, Dis, Logos, Mars, Sanctum, Shadows, and Untamed are implemented; the rest
are planned.

- Brobnar — board presence and large fighters that benefit from fighting
- Dis — destruction of creatures on both sides for benefit
- Ekwidon — exchange this for that, in their favor
- Enlightened — build a board presence toward a big payoff that must be disrupted
- Geistoid — use the discard pile as a resource
- Logos — efficiency and card draw
- Keyraken — large monsters that everything else revolves around
- Mars — insular synergy, at the expense of friendly and enemy non-Martians
- Ouroboros — stay exhausted for benefit and hoard Æmber
- Redemption — redeem the other houses; soak distortion
- Sanctum — protect the board and neighbors
- Saurian — risk/reward by putting Æmber on the board for benefit
- Shadows — small but stealy
- Skyborn — care about board placement for benefit (and, from set 8, key colour)
- Star Alliance — cooperate with other houses for benefit
- Unfathomable — disrupt the opponent's hand and Æmber pool
- Untamed — Æmber rush

## Keywords to consider (not yet built)

- Ignore Elusive (MTG: reach), First Strike, Defender, Haste, change control,
  Abduct, Scry/Surveil/Mill (top-of-deck manipulation).

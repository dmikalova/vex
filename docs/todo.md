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
- custom keyboard shortcuts saved to player profiles
- manual mode needs to prompt for confirmation
- resolution zone

## Open wording questions

- Dropping `as if it were yours` needs care. The phrase is load-bearing on cards
  that use or play a card they do not control (Sneklifter, Ghostly Hand, Nexus,
  Talent Scout): the rulebook says you resolve the card as if you controlled it
  without ever gaining control, which decides who its "friendly"/"enemy" reads
  from and who a `Destroy this card` falls on. If the phrase goes, the rule has to
  be stated once elsewhere — e.g. using a card always grants, for that use, the
  reading that you control it — and the cards whose text turns on that reading
  need to be checked one by one rather than sed-replaced.

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

## Brain spam

- If a card is drawn/hidden data revealed then no undo. No undo across turn boundaries
- Log system should mark what is done, not what the card says/assist
- Mars -> Venusian (or Cytherian)
- Font page or better yet a design page
- icons should contain both house colors, and should be roughly roundish. Brobnar - flame, sanctum cross in shield should be the yellow, dis
- grill me on image generation
- being able to navigate completely by keyboard - need a position arrow
- Asynchronous matches
- minimize main page load to base css, icon, and wasm load by hash. index.html has no cache. wasm has json manifest of everything that caches for a while
- The action panel (context.md could have wording for this) could be the actual card and text, and then play/reap/ etc buttons within
- Upgrades could add the text that they add to the card
- No way to see upgrade visually... adding card top/bottom takes up valuable space, maybe right/left
- Squeeze text onto card title
- If low on vertical space, creatures and artifacts could shrink their bottom label while in play. Want to be able to see status most of all
- one click bug report with full logs, state, actions taken, and comments
- in mobile, have the action bar under all the cards - or on the card preview?
- stadiums
- future/ancient cards set like evil twins
- Establish full turn phases - start of turn, forge, house, archives, general play, ready, draw, end of turn
- Redo logging like the planned TCO logging
- drag and drop creature directly into battleline flank (or deploy, with dynamic moving as you go across), upgrade onto creature, artifact into artifact line
- For things like steal, capture - they both have an amount, By, Max, etc, similar structure. Is it possible to reuse these interfaces eg for Economy types, and other types?
- Skill to iterate on implementation repeatedly
- Re-run the GameState layout check once the later sets land. Adding at least
  four more card types will widen `CardType`, `Bar[CardType]`, and anything else
  keyed by type, and new mechanics tend to add fields. Measure with
  `unsafe.Sizeof(GameState{})` and a `reflect` field/offset dump, then re-decide
  the two levers left on the table: `maxCards = 128` (68% of the state, but the
  headroom is load-bearing for the sandbox's `game_manual.go` card creation) and
  packing `CardCore`'s four bools into a bitfield (~512 bytes, at the cost of
  read-modify-write bugs and debuggability). History: 4232 -> 4112 (per-turn play
  permissions to uint8) -> 4024 (CardType string to enum).

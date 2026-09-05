# TODOs & design notes

A grab-bag of things to fix, build, and answer. Canonical vocabulary is
in [../CONTEXT.md](../CONTEXT.md), the rules in the engine's rulebook term registry
(the `/rulebook` page), and the long-term vision in [roadmap.md](roadmap.md).

## Grill me

### Current focus

Typed rulebook registry (ADR 0018):

- DEFERRED — classifying effects (RuleBearing/Composition/Implementation) and
  enforcing effect completeness (ADR 0018 exempts effects "until that
  classification exists"; assigned to ADR 0019). Needs an effect-type enumeration
  that does not yet exist.
- DEFERRED — extending completeness to turn/combat steps. Needs first-class
  turn/combat step values, not a parallel list invented for the test; the current
  `Phase` enum mixes implementation-only phases.

Rules-voice text passes (no behavior change):

- Code-comments migration (ADR 0020) — staged package-by-package. The resource-spelling sweep is DONE: the only prose `Aember` in a comment (effect_aember.go's resolveGate) is now `Æmber`; the remaining `Aember` in comments are Go symbol names (`Aember` method/field), card names (`Irradiated Aember`), or generated stub source-text, which correctly keep the ASCII form. The broader plain-voice comment rewrite continues per package.

### Next focus

- grill me on image generation. image generation should adapt with upgrades and other constant abilities
- event sourcing
- decklists
- Start of game setup - p1 plays 1 cards, mulligan
- drag and drop creature directly into battleline flank (or deploy, with dynamic moving as you go across), upgrade onto creature, artifact into artifact line
- minimize main page load to base css, icon, and wasm load by hash. index.html has no cache. wasm has json manifest of everything that caches for a while. Also css should be minified and compressed, what about everything else?
- The action panel (context.md could have wording for this) could be the actual card and text, and then play/reap/ etc buttons within
- profiling - eg running property tests and outputting the profiled usage for hot paths, and then optimizing those paths as a skill
- using property testing to find unused code paths and then force specific tests there
- Is there a way to validate that the UI handles and presents all possible game states/prompts? eg if I add a new prompt route, can the UI then automatically fail bc its not handled?
- On the style page add a section with all of the Log and Text usages rendered out. The easiest wayt to do this might be to create a dedicated preview area that dynamically displays these elements as they are used in the engine (eg show a set of cards that covers every rendering element, and a log that does the same for all log entries)
- card gallery (and search). Gallery links to cards, and cards can pull in all the relevant rules onto that page
- For things like steal, capture - they both have an amount, By, Max, etc, similar structure. Is it possible to reuse these interfaces eg for Economy types, and other types?
- Improve master of x materialization
- House Ambassador (eg Brobnar Amassador) as a materialization - make it work as a legacy/maverick to swap with a card in another house
- In the rulebook have an Accuracy example-binding ratchet — let terms cite a real engine test, then require it for subtle rules over time so that players can interact with the examples and understand the evolving rules context.

## Things that can be done now

- remove abduct / simplify to archive targets - the rules already naturally handle how archiving your opponent's cards works

## UI finesse

- animation library and overhaul
- playing an action card should have an animation - eg go to center, get big, go to discard
- destroy animations are going under
- steal and capture animation
- refine being able to navigate by keyboard
- peeking opponent's facedown cards should show the card back for the hover - can do after token creatures
- More keyword icons - how much is too much?
- Styles page should automatically add new animations to the list
- discard from hand and other zones animation
- when selecting cards like for mothergun it should get a checkmark, not dim, and also be able to click again to uncheck
- simplify s curve
- toggle animations
- a whole ass settings panel
- house icons should contain both house colors, and should be roughly roundish. Brobnar - flame, sanctum cross in shield should be the yellow, dis
- Cannot act dialogue on cards is not necessary

## Game finesse

- after implementing all cards, identify cards that have unique effects and decide if they can be reworded for simplicity - is it possibility to introspect and see how many times each card facet is used?

## Full two-player support

- base58 for deck IDs
- import from MV
- manual mode needs to prompt for confirmation
- custom keyboard shortcuts saved to player profiles
- set your own primary/secondary player color
- If a card is drawn/hidden data revealed then no undo. No undo across turn boundaries
- Asynchronous matches
- one click bug report with full logs, state, actions taken, and comments. Also a feedback form
- single player mode (current) and vs bot mode
- alliance
- custom deck builder
- /demo route
- toggle keyboard shortcuts
- ability to pin players to an engine version, and then when they go to play their game they just load that engine for that game even if its an async game

## Wild ideas

- tool to extract cards from MV
- generate 10k decks, score them, and graph their scores with average, mean, std dev, and 95/99/99.9%iles
- translations
- Display multiple houses
- resolution zone
- stadiums
- future/ancient cards set like evil twins
- Change enters play ready/stunned/enraged to Play: Stun X - would change timing for dominator etc
- MM mutants - have a common, uncommon, and rare variant
- rockatiel - the concept of really good cards that mean you have to hold answers against them for archon, vs not having complete blowout surprises that you have to hold against in sealed
- If a maverick has a fate, it should pull in prophecies - how to balance prophecies so they could be in any deck?

## Bot support

- Re-run the GameState layout check once the later sets land. Adding at least
  four more card types will widen `CardType`, `Bar[CardType]`, and anything else
  keyed by type, and new mechanics tend to add fields. Measure with
  `unsafe.Sizeof(GameState{})` and a `reflect` field/offset dump, then re-decide
  the two levers left on the table: `maxCards = 128` (68% of the state, but the
  headroom is load-bearing for the sandbox's `game_manual.go` card creation) and
  packing `CardCore`'s four bools into a bitfield (~512 bytes, at the cost of
  read-modify-write bugs and debuggability). History: 4232 -> 4112 (per-turn play
  permissions to uint8) -> 4024 (CardType string to enum).
- Monte Carlo Tree Search, minimax, reinforcement learning
- Method B: Surrogate Regression (The Recommended Approach)
  You let a state-of-the-art Deep RL agent (or an AlphaZero-style hybrid of RL + MCTS) play hundreds of thousands of matches to generate a massive dataset of deck compositions and their actual win rates.

Once you have this raw data, you apply a standard, human-readable machine learning algorithm (like Ridge or Lasso Regression) over the dataset to predict the RL agent's win rates.

This regression will naturally spit out the coefficients for individual cards and pairwise interactions. This effectively reverse-engineers the RL’s "black box" brain into a highly accurate, DoK-style spreadsheet.

Which should you use for parameter tuning?

If you are currently tuning parameters by using MCTS as an evaluator (e.g., MCTS plays 1,000 games -> outputs win rate -> you adjust synergy weights -> repeat), you are likely facing a massive computational bottleneck. MCTS is simply too slow to run the millions of simulations required to tune an exhaustive matrix of CCG synergies.

The ideal pipeline: Use an AlphaZero-style architecture. Use a neural network to evaluate board states, and use a lightweight MCTS to look just 1-2 turns ahead to choose the actual play. Let this AI play millions of games to generate a dataset of deck match-ups, and run a linear regression on those match-ups to extract your human-readable synergy and anti-synergy parameters.

- [building a rating engine with alphazero](https://gemini.google.com/app/24b5499fc76c5fc1)
- Should be able to transfer the rating system to a KF rating system as long as I don't drastically change the rules - eg prophecies or the tide :/
- Ask the system who has better odds - P1 vs P2, what about mulligan? What is the line for mulliganing?
- Can the bot identify under rated cards and have bot play them more - although seems like this would be at the mechanics level?

## Design principles (kept from KeyForge)

- You decide everything on your own turn — no interrupts.

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

### House changes

- Mars -> Venusian (or Cytherian)

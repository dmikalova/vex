# TODOs & design notes

A grab-bag of things to fix, build, and answer. Canonical vocabulary is
in [../CONTEXT.md](../CONTEXT.md), the rules in the engine's rulebook term registry
(the `/rulebook` page), and the long-term vision in [roadmap.md](roadmap.md).

## Grill me

### Current focus

- tool to extract cards from MV

### Next focus

- shards should pull in shards for the other houses
- House Ambassador (eg Brobnar Amassador) as a materialization - make it work as a legacy/maverick to swap with a card in another house
- bane, brew (common), plant, and blaster variant
- Way to always settle damage anytime power could change, instead of having to have settles strewn about the codebase. Similarly, way to settle that a card is no longer in play, so its abilities don't proc, and things that it may have triggered can no longer target it consistently instead of having to know all the call sites - eg redacted strange gizmo forge a key was putting amber back on redacted
- event sourcing
- drag and drop creature directly into battleline flank (or deploy, with dynamic moving as you go across), upgrade onto creature, artifact into artifact line
- The action panel (context.md could have wording for this) could be the actual card and text, and then play/reap/ etc buttons within
- using property testing to find unused code paths and then force specific tests there
- Is there a way to validate that the UI handles and presents all possible game states/prompts? eg if I add a new prompt route, can the UI then automatically fail bc its not handled?
- On the style page add a section with all of the Log and Text usages rendered out. The easiest wayt to do this might be to create a dedicated preview area that dynamically displays these elements as they are used in the engine (eg show a set of cards that covers every rendering element, and a log that does the same for all log entries)
- card gallery (and search). Gallery links to cards, and cards can pull in all the relevant rules onto that page
- In the rulebook have an Accuracy example-binding ratchet — let terms cite a real engine test, then require it for subtle rules over time so that players can interact with the examples and understand the evolving rules context.
- Be able to set up situation and then run it in the engine UI for playwright
- rename to Vex
- remove abduct / simplify to archive targets - the rules already naturally handle how archiving your opponent's cards works
- can splash and splashattack be combined?
- enemy creature should be indicated in archives and even under my control
- Update card.New to be all opts
- sequence vs sentences wording - eg sequence is obviously game, and sentences is textual, but they're both textual and game
- Improve mega creatures
- instead of having to manually bump the state version would it be possible to hash changes to how the state is written so it automatically bumps on such changes, but also not on irrelevant changes? re event sourcing. If the hash was based on the action signature rather than overall engine you could check when loading the event sourcing if any of the used actions changed
- Maverick houses
- Move the prompt generation and options into engine rather than web (eg when playing an upgrade, am prompted to "Choose a creature to attach Stunner onto")

## Things that can be done now

- ? bdq (see screenshot) is doing the action bar with title cards thing
- ? they're everywhere sequence
- auteresolve button
- gigantic, tide
- have to double click to activate preview from logs
- facedown cards should be facedown for both players - its just that the controller can hover over to peek at the card - its important to be able to visually distinguish what is a facedown or faceup card
- test to make sure there are no unused assets
- Instead of "OnIt" should we use "OnTarget"
- Anomaly provenance
- For house select on mobile, star alliance goes onto 2 lines instead of squeezing onto the house select button. In general the text for buttons should squeeze onto one line. Similarly, why does the title for Yshi shrink by a standard amount rather than just shrink to the right amount? Are we still doing predetermined shrink amounts for cards? Is it not possible to have the text dynamically shrink to the right amount? Can you explain to me what's going on and why this can't work smoothly?
- creeping oblivion prompt
- Igon the green can just play terrible - also connected 1 to 1
- OneCopyPerDeck should be ordered higher up, check the rest of the ordering eg constant should be above play, then reap, then fight
- festering touch

## UI finesse

- should rigged lottery log everything together
- if cards are in action bar buttons - just have clickable preview toggle on right
- split the zone dialog into each zone
- the back button for house choice is under - could be in line with choose a house top right? Need to overall decide where the undo button goes on mobile
- s curve fix
- center card name and traits?
- Manual mode should allow you to move deck card to hand etc
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
- Renaming the draw pile to reserve so that deck list, the full deck itself, and the deck pile are distinct and clearly named
- Choose one: rewrites

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

- aember on artifacts goes to opponent?
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
- Find the 100 longest card tests in keyteki and digest them down to what the test is trying to capture
- manual mode - change card house, edit bonus icons/distortions - only on manual mode cards
- Bonus icons don't resolve if the creature dies while resolving them, and they count as the creature dealing the effect, not the game
- enhancements across CotA/AoA/WC
- non-aember default bonus enhancements

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

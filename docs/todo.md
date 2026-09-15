# TODOs & design notes

A grab-bag of things to fix, build, and answer. Canonical vocabulary is
in [../CONTEXT.md](../CONTEXT.md), the rules in the engine's rulebook term registry
(the `/rulebook` page), and the long-term vision in [roadmap.md](roadmap.md).

## Grill me

### Current focus

### Next focus

- using property testing to find unused code paths and then force specific tests there
- Is there a way to validate that the UI handles and presents all possible game states/prompts? eg if I add a new prompt route, can the UI then automatically fail bc its not handled?
- On the style page add a section with all of the Log and Text usages rendered out. The easiest wayt to do this might be to create a dedicated preview area that dynamically displays these elements as they are used in the engine (eg show a set of cards that covers every rendering element, and a log that does the same for all log entries)
- card gallery (and search). Gallery links to cards, and cards can pull in all the relevant rules onto that page
- In the rulebook have an Accuracy example-binding ratchet — let terms cite a real engine test, then require it for subtle rules over time so that players can interact with the examples and understand the evolving rules context.
- Be able to set up situation and then run it in the engine UI for playwright
- rename to Vex
- remove abduct / simplify to archive targets - the rules already naturally handle how archiving your opponent's cards works
- Update card.New to be all opts
- sequence vs sentences wording - eg sequence is obviously game, and sentences is textual, but they're both textual and game
- instead of having to manually bump the state version would it be possible to hash changes to how the state is written so it automatically bumps on such changes, but also not on irrelevant changes? re event sourcing. If the hash was based on the action signature rather than overall engine you could check when loading the event sourcing if any of the used actions changed
- Move the prompt generation and options into engine rather than web (eg when playing an upgrade, am prompted to "Choose a creature to attach Stunner onto")
- mage tool to view connected cards and amounts etc
- capture and bonus aember are both in the creature status area
- Distortion system instead of flat enhancements — see the Enhancement and
  Distortion entries in [../CONTEXT.md](../CONTEXT.md) and
  [deck-generation.md](deck-generation.md).
- gigantic, mimic gel
- Instead of "OnIt" should we use "OnTarget"
- Tool to open 50 random cards for me to review - and then record which ones I've seen how many times, so the next time it selects a different 50 with the least amount of reviews
- I really like this form: Grant: card.GrantPlay | card.GrantUse - where can we use it more?
- A tool that can detect card.X usage across all cards, to help identify where specific effects or abilities are being underutilized as a sign of an overly specific method.
- The way a lot of effects work is there is implied chaining between one effect to the next - is there a reasonable way to make this more explicit?
- Be able to load a test situation from a saved state or scenario file
- I've noticed that there are some UI sugars in the engine - I was wondering if it makes sense for there to be an intermediate layer - eg the engine handles state changes, the wrapper handles relevant trackers for the UI, and then the UI on top imports the wrapper and renders what it gives. For example, there are badges for counting how much damage is about to be dealt to each creature in a selection like gargantes scrapper. That seems purely UI, but also makes sense near the engine. My concern is performance when there is no UI - eg for MCTS - if MCTS is calculating the badges and never using them then that's potentially lost performance.
- Consolidate Destination and DeckDest - apparently the voicing would be a whole thing to add into this
- using shared dictionaries for wasm compression
- run a million games and then get stats on memory usage in the state and see where estimates are overly conservative and could be pulled back to save space
- Reordering the state in Go
- Minimizing the state by using bitfields more aggressively - tradeoff with having the interpret that in Go, but we are no cpu bound

## Things that can be done now

- [WithoutBonus](https://discord.com/channels/802313100485197855/802313100987990053/1549192438978191372)
- Decomposables:
  - OpponentForgedKeys
  - AfterFriendlyCreatureFights

- Capitalizing card types in text - eg Creature in Floomf
- Changing card.X to instead be e.X eg for engine - is the facade really providing value, or is there anything else we could do to organize the repo better instead of one mega engine?
- Lumilu - could card.InPlay be better represented by filters or refinements?
- mutants: technofiend, daemosaurus - you should be able to find the rest of the house suffixes/prefixes from these two
- shard of unity prompt doesn't lift creature for use
- rows have excess scroll space and don't hide the scroll bar by default

- WithAemberCost and Toll could be combined into MustPay
- decompose all the neighbor stuff
- Granted: card.FightReap(card.ArchiveGrantingUpgrade{}), should be card.Archive{Target: GrantingUpgrade}
- Why is DamageThen and ChooseCreatureThen needed? Why can't these just be sequences that pass along the effect context?
- Get rid of bar.go

- Livia and Fidgit could go further

## Sites of all the things

- effects.go
- options.go
- target.go
- types.go

## UI finesse

- Move the manual mode dialog into the zones modal
- Make the zones modal not a modal - just have it be a full screen panel
- enemy creature should be indicated in archives and even under my control
- The action panel (context.md could have wording for this) could be the actual card and text, and then play/reap/ etc buttons within
- drag and drop creature directly into battleline flank (or deploy, with dynamic moving as you go across), upgrade onto creature, artifact into artifact line
- molina's blaster and glyph spacing
- creeping oblivion prompt - currently has a done button at the top of the zone modal - should be at bottom outside the modal
- Eliminate iconFallbackAllowed for glyphs
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
- Stilt-kin needs to pull in 2 giants - in general a rule is needed for trait specific cards to pull in 2 of those cards
- VM 25 anomalies? Omega TT etc
- Special cards have a special treatment (skybeasts, revenants, dragonscale)

## Game finesse

- after implementing all cards, identify cards that have unique effects and decide if they can be reworded for simplicity - is it possibility to introspect and see how many times each card facet is used?
- Renaming the draw pile to reserve so that deck list, the full deck itself, and the deck pile are distinct and clearly named
- Choose one: rewrites
- After implementing all cards - pull 20 decks of each card from DoK and see which cards cannot be in multiples (eg tmtp)
- "Play a card from your archives"

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

- can splash and splashattack be combined?
- aember on artifacts goes to opponent?
- generate 10k decks, score them, and graph their scores with average, mean, std dev, and 95/99/99.9%iles
- translations
- Display multiple houses
- resolution zone
- stadiums/arenas - terrains - similar to stadiums but not exclusive, bonus and penalty
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
- Austin's house of silicates
- A full house worth of enhancements (eg you have a normal 36 card deck, but then there are 12 house enhancements of a fourth house that are randomly assigned)
- Weather effects - eg 4 sided reference card that turns

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
- looking for cards that have a wide range of value across different decks, vs a spike in always being good or bad
- Refocus on the board over one-shot actions.
- Cards with tradeoffs / situational value rather than being strictly good.
- Lean on upgrades to make boards more dynamic and flexible.

## Design directions (deliberate divergences from KeyForge)

- Minimize simultaneous effects — resolve one at a time, matching the physical
  game.

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
- Ouboros — stay exhausted for benefit and hoard Æmber
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

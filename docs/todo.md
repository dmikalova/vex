# TODOs & design notes

A grab-bag of things to fix, build, and answer. Canonical vocabulary is
in [../CONTEXT.md](../CONTEXT.md), the rules in the engine's rulebook term registry
(the `/rulebook` page), and the long-term vision in [roadmap.md](roadmap.md).

## Things that can be done now

- rename to Vex

- Library of polli should let you choose any creature, even with 0 aember. Could be changed to "choose a creature."
- Lost in the woods should be a vacuous choice for any side that has 0,1, or 2 creatures
- Loot the bodies should have player gain all aember simultaneously
- exile - creature appears on other side before flank choice
- Techno-saurus - wording should be Play: then Reap:
- Defense initiative - click target to exalt and ward, or click done
- Vaults blessing should cluster with mutants
- Niffle ape missing glyphs
- Decompose WithEachPlayerAbility
- golden spiral prompt buttons are on any card I click - should be stuck on mack
- axiom/troop call - generic choose a  creature for ward and bonus damage in prompt
- Update card.New to be all opts
- The way a lot of effects work is there is implied chaining between one effect to the next - is there a reasonable way to make this more explicit?
- Consolidate Destination and DeckDest - apparently the voicing would be a whole thing to add into this
- Changing card.X to instead be e.X eg for engine - is the facade really providing value, or is there anything else we could do to organize the repo better instead of one mega engine?
- Split out glyphs more in icon.go
- Instead of "OnIt" should we use "OnTarget"
- I really like this form: Grant: card.GrantPlay | card.GrantUse - where can we use it more?
- remove abduct / simplify to archive targets - the rules already naturally handle how archiving your opponent's cards works - this is fine as a label in tests, but want to remove it from rules and comments in the engine
- Livia and Fidgit could go further in decomposing (eg kompsos)
- WithAemberCost and Toll could be combined into MustPay
- decompose all the neighbor stuff
- Granted: card.FightReap(card.ArchiveGrantingUpgrade{}), should be card.Archive{Target: GrantingUpgrade}
- Why is DealDamage and ChooseCreatureThen needed? Why can't these just be sequences that pass along the effect context?
- Get rid of bar.go
- /cards view cuts off side icons - why isn't this rendering like in the engine?
- Do a sweep for defaults and missing explicit - eg destination.go
- Be able to load a test situation from a saved state or scenario file
- I've noticed that there are some UI sugars in the engine - I was wondering if it makes sense for there to be an intermediate layer - eg the engine handles state changes, the wrapper handles relevant trackers for the UI, and then the UI on top imports the wrapper and renders what it gives. For example, there are badges for counting how much damage is about to be dealt to each creature in a selection like gargantes scrapper. That seems purely UI, but also makes sense near the engine. My concern is performance when there is no UI - eg for MCTS - if MCTS is calculating the badges and never using them then that's potentially lost performance.
- using shared dictionaries for wasm compression
- using property testing to find unused code paths and then force specific tests there
- card gallery (and search). Gallery links to cards, and cards can pull in all the relevant rules onto that page
- Be able to set up situation and then run it in the engine UI for playwright


- Auto vac : choose one:
- how am I supposed to pick from opponents hand? hidden stash
- nexus should be able to target artifact with no action abilities
- change house should be yellow outline instead of circle
- sow salt doesn't show up as restriction
- universal translator is not right - let's me select creature for grant, and then need to follow up with use
- consulmprimis - should be able to choose empty creature
- essence scale does printed house instead of academy training house
- novu dynamo should be click on novu to destroy
- alpha doesn't work with picking up archives??
- rename artifact line to field/domain
- card title eg sanctum's cards don't get dimmed
- remove grey riders may
- rockatiel should be play: 3 creatures
- are taps and clicks distinguishable? Can zones be double tappable
- deck list is opening both players again - add a test for regression?
- chronus may - click on chronus or done, should be able to just click on cars in hand
- is chronus resolving for any bonus icon?
- one stood against many chooses 3 friendlies to fight - choose a creature
- tap to show tooltip, tap again to open zone
- tap/move does not highlight player bar
- copy logs - plaintext and structural
- random sets or pick sets
- new game set selection cancel should be red
- if I go new game, select a set, cancel, and then go again it gets stuck
- zone modal x highlights on tap
- fidgit prompt - text is not aligned, buttons should just say discard or archives
- times gigantics glyph
- WTF is going on with tomes gigantica text
- is there a way to test over web behaviors without having rigid tests but also not being so loose that I only notice it IRL - eg a way to loosely define design principles that can be used for automated tests?
- one of the best things KF does is present cards and situations with tradeoffs which you then have to navigate the value proposition - look to remove/edit cards that provide no trade-off (AAF), or boost cards that tend to just provide no value
- The engine has a lot of duals - eg an ability that does this, vs an effect? Not sure of the verbiage. I understand the duals are necessary, but is it possible to unify their underlying engine? Even to the point that defining one also creates the dual automatically?
- Allison did not change houses
- the icon overlay for dealing damage should count how many are left
- while in key selection can't tap outside on nothing to lower cards or remove hover
- AST analysis for code complexity - https://share.gemini.google/hJBkse4DHYEi
- graph out the mechanics split of each house
- one page design docs for vactrol
- sorting hand by houses not happening
- tomes gigantica without a gigantic
- pandemonium should let you choose and should not count 0
- hidden stash - can't select from opponent's hand
- duration markers for more than turn
- how to deal with restrictions like kaupe?
- How affected is MCTS if we do turn granularity instead of every command?
- center line and zones need to be clearer
- Austin and Cody discussion from sep 19
- tapping again on a lifted card or preview hover should dismiss it
- after forging third key still prompted for quiet anvil and forge compiler
- Questions for change: what is present, what is past, and what are the important questions
- rockatiel - selection happens and then cards move
- radiant Truth and Commandeer - order by name, have long press preview or hover
- Auto resolve if there are more than 2 options (random) - don't step prng, just use original seed?
- eyegor the remaining choices are not vacuous - test for that
- Lord Golgotha - splash attack 3
- guji - 2 damage then 4
- make set choice random
- scout Pete should just be action bar dialogue instead of open zone modal
- Scout - text until the end of turn for skirmish
- dark harbinger and scout - no way to resolve scout first
- scout select and fights at the same time
- it's coming glyph
- lost in the woods vacuous choice
- it's coming search done, and select from zones
- change the chosen creature to it
- protect the weak does not give fresh armor on play
- adding a critic pattern to the agent
- dominator bauble a vinda gives no choice on who to damage
- lights out doesn't need a may
- warrant counters not showing on book of malefaction
- remove openspec

### Automatic linters

[https://lobehub.com/skills/saifoelloh-golang-best-practices-skill-design-patterns?activeTab=installation](lint skill)

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

- [WithoutBonus](https://discord.com/channels/802313100485197855/802313100987990053/1549192438978191372)
- after implementing all cards, identify cards that have unique effects and decide if they can be reworded for simplicity - is it possibility to introspect and see how many times each card facet is used?
- Renaming the draw pile to reserve so that deck list, the full deck itself, and the deck pile are distinct and clearly named
- Choose one: rewrites
- After implementing all cards - pull 20 decks of each card from DoK and see which cards cannot be in multiples (eg tmtp)
- "Play a card from your archives"

## Internal tooling

- In the rulebook have an Accuracy example-binding ratchet — let terms cite a real engine test, then require it for subtle rules over time so that players can interact with the examples and understand the evolving rules context.

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

- Tide:
  - A dial from -3 to +3
  - Tide turns on its own (high neutral low) - fixes Lacus
  - Raising the tide pays opponent aember
  - Shoreline - discarded cards to change the tide go here, raising the tide lets you play one
  - The tide height turns off creatures lower than its height for both players. Submerged creatures get their tide bonus
- can splash and splashattack be combined?
- aember on artifacts goes to opponent?
- translations
- Display multiple houses
- resolution zone
- Change enters play ready/stunned/enraged to Play: Stun X - would change timing for dominator etc
- MM mutants - have a common, uncommon, and rare variant
- rockatiel - the concept of really good cards that mean you have to hold answers against them for archon, vs not having complete blowout surprises that you have to hold against in sealed
- If a maverick has a fate, it should pull in prophecies - how to balance prophecies so they could be in any deck?
- Find the 100 longest card tests in keyteki and digest them down to what the test is trying to capture
- manual mode - change card house, edit bonus icons/distortions - only on manual mode cards
- Distortion system instead of flat enhancements — see the Enhancement and
  Distortion entries in [../CONTEXT.md](../CONTEXT.md) and
  [deck-generation.md](deck-generation.md).
- Bonus icons don't resolve if the creature dies while resolving them, and they count as the creature dealing the effect, not the game
- enhancements across CotA/AoA/WC

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
- run a million games and then get stats on memory usage in the state and see where estimates are overly conservative and could be pulled back to save space
- Reordering the state in Go
- Minimizing the state by using bitfields more aggressively - tradeoff with having the interpret that in Go, but we are no cpu bound
- generate 10k decks, score them, and graph their scores with average, mean, std dev, and 95/99/99.9%iles

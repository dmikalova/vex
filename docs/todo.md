# TODOs & design notes

A grab-bag of things to fix, build, and answer. Canonical vocabulary is
in [../CONTEXT.md](../CONTEXT.md), the rules in the engine's rulebook term registry
(the `/rulebook` page), and the long-term vision in [roadmap.md](roadmap.md).

## Grill me

### Current focus

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
- Legacy houses
- Move the prompt generation and options into engine rather than web (eg when playing an upgrade, am prompted to "Choose a creature to attach Stunner onto")
- mage tool to view connected cards and amounts etc
- capture and bonus aember are both in the creature status area
- Anomaly provenance - eg move them to their home set, or put them in an anomaly set
- Distortion system instead of flat enhancements — see the Enhancement and
  Distortion entries in [../CONTEXT.md](../CONTEXT.md) and
  [deck-generation.md](deck-generation.md).
- Before implementing totally new mechanics, must attempt to fit them into existing mechanics, including by expanding them. If its determined to not be possible or reasonable, must ask and explain why and get confirmation before proceeding.
- gigantic, mimic gel
- Instead of "OnIt" should we use "OnTarget"
- Tool to open 50 random cards for me to review - and then record which ones I've seen how many times, so the next time it selects a different 50 with the least amount of reviews
- I really like this form: Grant: card.GrantPlay | card.GrantUse - where can we use it more?
- A tool that can detect card.X usage across all cards, to help identify where specific effects or abilities are being underutilized as a sign of an overly specific method.

## Things that can be done now

- Hunter or Hunted could be changed to "Remove a ward from a creature and ward a creature." Note - any creature could be targeted for removing the ward, even if it doesn't have a ward. This effectively makes it the same effect but a lot simpler.
- Can WithPlayFightReap, WithFightOrReap, and WithPlayReap be consolidated into WithAbility(card.Trigger.PlayFightReap) etc?
- Philophosaurus implementation - currently Philophosaurus uses "LookAtTopSort" which hard codes its behavior. However, there are a bunch of other cards that do a very similar thing - eyegor, lay of the land, vandalize in an upcoming set - I feel like all of these could be generalized into a "Look at the top x cards of deck, then do x (and optionally y, z, etc)"
- counters should be named "generic-counter-$NAME" so that they all get organized together. eg effect_generic_counter_growth_test.go and generic-counter-growth.svg - also move effect_counter.svg to effect_generic_counter.svg
- On the player bar, if I try swipe left/right on touch screen it does not scroll the player bar if I end up touching an icon
- The tooltips on the player bar are now in the player bar rather than over it
- phalanx strike repeats endlessly - it should be able to repeat one time (hence the "preceding effect" text.)
- shadow self should not deal damage in a fight - eg when it fights or when something fights it
- creature as upgrade prompts for flank then asks if you want to play it as a creature or upgrade. The play button itself above the lifted card could say "play creature" and "play upgrade" instead of just "play"
- Is it possible for things like "AttachSelfTo" to take a type of card.Name instead of string? Similar for GrantingArtifact.Named() and anything else that currently accepts a stringly typed card name.
- uncharted lands should be card.this or card.target.source to reference itself instead of a stringly typed card name.
- Chief eng walls currently does Type: Type, OrTrait: Trait which is weird and inflexible. Can the zone movement methods just take a generic Target: card.Type.Upgrade | card.Trait(card.Trait.Robot) instead?
- rename SpendAsPool (bracchus, calypigean) - eg "SpendAemberOnCard"
- when I'm prompted to select a card from zone, after some short time the modal jumps to the top
- In the sidebar, the stealth mode and garcia restriction warnings are on the same line - warnings should be one per line
- Data forge doesn't need may
- creature next to narp shouldn't have option to reap through universal translator
- when universal translator selects a creature, have the normal use buttons appear above it (eg fight/reap, if an action is available, action)
- The duration names can be more explicit - eg duration.NextTurn could be duration.UntilStartOfYourNextTurn or duration.UntilEndOfYourNextTurn.
- Update the agents file so that when it writes tests it tests the positive and negative of a card's abilities. Simple rote abilities like keywords don't need to be tested, but combined conditionals should have both positive and negative test cases. Cards that handle numeric values should be tested at the 0, 1, n-1, n, and n+1 cases
- Cards like virtuous works with no text should still have an empty black textbox that fills the card space
- manual mode option to place card under/return card under to hand on a card that's selected/lifted

- Cloaking Dongle: Target: Target and neighbors then gives the bonus
- Kompsos Haurspex and Livia the elder can be atomized into each other
- Gebuk can be simplified
- Nizak should be "While in a fight,"
- How could Encounter suit can be simplified? "This creature is invulnerable while resolving a Tactic card"?
- KeyForgery...

## UI finesse

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
- looking for cards that have a wide range of value across different decks, vs a spike in always being good or bad
- Refocus on the board over one-shot actions.
- Cards with tradeoffs / situational value rather than being strictly good.
- Lean on upgrades to make boards more dynamic and flexible.

## Design directions (deliberate divergences from KeyForge)

- Minimize simultaneous effects — resolve one at a time, matching the physical
  game.
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

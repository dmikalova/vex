# Vactrol⇄KeyForge divergence register

This is the single record of where Vactrol departs from KeyForge, and why. It is
one half of the split of the old card-wording document
([ADR 0019](adr/0019-controlled-rules-voice.md)): the wording conventions —
surface-independent house style — live in
[card-wording-rules.md](card-wording-rules.md); the deliberate departures live
here.

## Precedence

The precedence rule is fixed. Vactrol wins where Vactrol has decided. The
[KeyForge Master Rulebook](keyforge-master-rulebook.md) is the wording authority
only for what Vactrol has not decided. A later "match KeyForge" never silently
overwrites a divergence recorded here.

Matching the printed KeyForge text is not a goal in itself, and is never a
blocker. **Template consistency and simplicity across Vactrol outrank exact
KeyForge wording.** A meaning-preserving reword that puts a card in the same voice
and template as the rest of Vactrol — or that lets a mechanic decompose into
shared nodes rather than a bespoke one — is welcome, and does **not** earn a row
below: this register catalogs departures that change a **rule or a name**, not
pure house-voice alignment. Record a reword here only when it changes what a card
does or what something is called.

## Wording divergences

Each of these is a wording convention that changes a rule or a name, not just
phrasing. The convention itself — with its examples and affected cards — lives in
the numbered rule cited below in
[card-wording-rules.md](card-wording-rules.md). This register is the index that
answers "where does Vactrol diverge from KeyForge, and why".

| Divergence                                                    | Rule    | Reference                                             |
| ------------------------------------------------------------- | ------- | ----------------------------------------------------- |
| `Sacrifice` collapses into `Destroy`                          | rule 3  | one destruction verb                                  |
| `Return X` becomes `Put X …`                                  | rule 4  | one movement verb over every destination              |
| `Omni:` becomes `Versatile` plus an `Action:`                 | rule 12 | [ADR 0009](adr/0009-tactic-type-omni-as-versatile.md) |
| One card is renamed (Crazy → Bonkers)                         | rule 14 | avoids a trademark collision                          |
| `pay` becomes `give` for player-to-player Æmber               | rule 18 | one transfer verb                                     |
| The `action` card **type** is renamed `Tactic`                | rule 19 | disambiguates the type from the `Action:` ability     |
| `while under your control` becomes a one-time swap            | rule 20 | avoids continuous re-checking; sticks with the card   |
| A deferred play permission becomes an immediate play          | rule 21 | avoids turn-scoped unused-permission memory           |
| A number-only `Otherwise` branch becomes `or <alt> if <cond>` | rule 22 | one linear sentence, no fork                          |
| A turn `step` is named a `phase`                              | rule 28 | [ADR 0012](adr/0012-first-class-turn-phases.md)       |
| Fight timing is named `in a fight with`                       | rule 29 | one phrase for the fight timing window                |
| A count cap is dropped — `(to a maximum of N)` is removed     | rule 30 | Vactrol has no count cap; the count is uncapped       |
| Timed "cannot use" restrictions share one phrasing            | rule 32 | one `Restrict` node, one duration phrase              |
| House-scoped reap bar reads in the shared "cannot use" voice  | rule 32 | Seismo-entangler folds into `Restrict`                |

## Per-card rule changes

A few cards were changed in ways that affect the rules, not just phrasing. Each
change simplifies the card toward base-rules text, makes it slightly more
interesting, or brings it in line with modern errata.

- **Charge!** buffs all creatures, not just ones played this turn. The `you play`
  clause is dropped.
- **Imperial Traitor** reads `Reveal`, not `Look at`. This is the modern wording.
- **Ganger Chieftain** and **Biomatrix Backup** are mandatory. The `you may`
  clause is dropped.
- **Ghosthawk** is mandatory. The `you may` clause is dropped, so it reaps with
  each of its neighbors, one at a time. Its target is also reworded from KeyForge's
  `each neighboring creature` to `each of Ghosthawk's neighbors`, naming the
  source's own neighbors rather than the general neighboring-creature set.
- **Malison** is mandatory. The `you may` clause is dropped, so its Fight moves an
  enemy creature every time (the flank capture still only fires when the moved
  creature ends on a flank).
- **Hypnotic Command** leans on the base rule that the active player makes all
  decisions, so `an enemy creature captures …` needs no explicit `choose`.
- **Phase Shift**, **Kirby**, and **Taber** play (or, for Taber, play or use)
  their off-house card immediately rather than granting a permission for later in
  the turn (rule 21).
- **Trust No One** is a `Choose one:` rather than a forced `If … Otherwise …`. The
  conditional branch ("if there are no friendly creatures in play, steal 1 Æmber
  per house among enemy creatures") still gates on the empty board, so it does
  nothing when you control a creature — a rational player picks the flat "steal 1
  Æmber" then, reproducing the original outcome, while the choice frame reads
  cleaner. Only a forced branch whose gated arm is a strict bonus over a safe
  fallback converts this way; most `If … Otherwise …` cards (random reveals,
  target-dependent or whose-turn conditions) do not.
- **Encounter Suit** is an upgrade that grants its host the reaction "After a
  Tactic is played but before it resolves, ward this creature." KeyForge phrases
  the reaction as the upgrade's own text and says "action card"; Vactrol renames
  that type to `Tactic` (rule 19) and fires every granted upgrade ability through
  the host, so it renders with the standard `This creature gains, "…"` wrapper
  like every other granted-ability upgrade.
- **Keyforgery** drops the trailing `(no Æmber is spent)` clarifier. Vactrol
  prevents the forge before any Æmber leaves the pool, so the clause states a
  consequence the mechanic already guarantees; the Rules voice omits such
  parenthetical asides.
- **Tantadlin** reads `Your opponent discards a random card from their archives`,
  not KeyForge's imperative `Discard a random card from your opponent's archives`.
  A random discard is the discarding player's own act, so Vactrol renders it in
  the actor's voice — the same voice Mind Barb already uses for a random hand
  discard. The effect is identical; only the voice changes.
- **Hunter or Hunted?** reads `Remove a ward from a creature, and ward a creature`
  rather than KeyForge's `Choose one: Ward a creature / Move a ward from a creature
to another creature`. The two branches collapse into one linear sequence:
  removing a ward and then placing a fresh one reproduces the "move a ward" branch
  when the same source and destination are chosen, and the "ward only" branch when
  the removal finds no ward. Vactrol keeps the two atomic effects — `RemoveWard`
  (any creature, warded or not) then `Ward` — instead of a bespoke `MoveWard` node,
  so there is one fewer one-off mechanic to carry.
- **Bait and Switch** reads `Steal 1 Æmber -> if your opponent has more Æmber than
you, repeat this effect`, not KeyForge's `If your opponent has more Æmber than
you, steal 1 Æmber. Repeat this effect`. Vactrol uniformly writes a self-repeat as
  `<do> -> if <cond>, repeat this effect` (the same shape Numquid the Fair and
  Neutron Shark use), so the steal leads and the condition gates the repeat. The
  first steal is therefore unconditional: with equal pools KeyForge steals nothing
  while Vactrol steals 1, then stops. In every case where the opponent already
  leads the two are identical.- **Gebuk** swaps the discarded creature into play immediately rather than waiting
  until it has left play. KeyForge reads "Destroyed: ... **after Gebuk leaves
  play**, put that creature into play in Gebuk's position"; Vactrol reads
  "Destroyed: Discard the top card of your deck. If it is a creature, swap it with
  Gebuk." The discarded creature and Gebuk exchange places in one step during
  Gebuk's Destroyed ability — Gebuk leaves to the discard pile (still counting as
  destroyed, since the destruction window already enrolled it) and the creature
  enters play in Gebuk's slot. This drops the deferred "after ... leaves play"
  ability memory in favor of an immediate, self-contained swap, and reuses the
  shared swap mechanic (SwapCards) instead of a bespoke delayed put-into-play
  registry. Because the swap resolves mid-window, a second Destroyed ability on
  Gebuk would fizzle (Gebuk is already out of play), and the swapped-in creature
  is on the board in time to be caught by the same destruction window (an
  enters-play "deal damage" can destroy it, and its own Destroyed ability then
  resolves in that window).
- **Harvest Time** reads `Choose a creature. Purge each creature that shares a
trait with the chosen creature`, not KeyForge's `Choose a trait. Purge each card
with that trait`. Vactrol reuses the shared choose-a-creature-then-fold-on-a-
  shared-trait mechanic (the same `ChooseCreatureThen` + `SharingTrait` pair
  Extinction uses) instead of a bespoke choose-a-trait purge, so the purge is
  anchored to a creature on the board and hits creatures only (artifacts are no
  longer swept). The chosen creature shares every trait with itself, so it is
  always among the purged. The per-player payout is unchanged: each player gains 1
  Æmber for each card they controlled that was purged this way.
- **Old Boomy** reads `Discard cards from the top of your deck until you discard a
Brobnar card or choose to stop -> deal 2 damage to Old Boomy. Archive each card
discarded this way`, not KeyForge's `Reveal cards from the top of your deck until
you reveal a Brobnar card or choose to stop. Deal 2 damage to Old Boomy if a
Brobnar card was revealed. Archive each card revealed this way`. Vactrol folds the
  card into the shared deck-dig family (`DiscardUntil`, the same node Sound
  the Horns and Invasion Portal use) instead of a bespoke reveal-and-archive loop:
  each card is discarded as the dig walks the deck, then the whole run is archived as
  a distinct step (`ArchiveDiscardedThisWay`). The end state is identical — every
  walked card ends in archives and Old Boomy takes 2 damage only when a Brobnar card
  is turned up — but the cards pass through the discard pile en route to archives
  rather than being archived directly on reveal.
- **The Worlds Collide "Brews"** are each redesigned to be unique to the giant they
  are brewed for, rather than sharing KeyForge's single `Play: Give a creature two
+1 power counters` action. Each brew is a Brobnar **Upgrade** that leads a Pull
  cluster bringing its Mega giant into the pod, and grants its host an ability tied
  to that giant: **Chieftain's Brew** grants `Fight: Ready and fight with a
neighboring Creature` (Mega Ganger Chieftain's play ability), **Cowfyne's Brew**
  grants `+2 splash-attack` (Mega Cowfyne's keyword), and **Alaka's Brew** grants
  `Fight: Play a Creature -> ready it`. This makes each brew a distinct card and
  keeps the Mega giants — which are `Connected` — reachable through their brew's
  cluster.
- **Mega Cowfyne** has **Splash-attack 2** rather than KeyForge's `Before Fight:
Deal 2 damage to each neighbor of the creature Mega Cowfyne fights`. The keyword
  reaches the same neighbours for the same 2 damage while folding the card into the
  shared Splash-attack mechanic (the same keyword base Cowfyne already carries),
  and **Cowfyne's Brew** grants that keyword instead of the bespoke reaction.
- **Mega Ganger Chieftain** and **Chieftain's Brew** are mandatory: the `you may`
  is dropped from Mega Ganger Chieftain's `Play` and from the `Fight` ability
  Chieftain's Brew grants, matching the same drop already made on their base card
  **Ganger Chieftain**.
- **Mind Barb** reads `Play: Discard a card from your hand. Your opponent discards
a random card from their hand`, adding a self-discard before KeyForge's lone
  `Play: Your opponent discards a random card from their hand`. The extra clause
  makes the card distinct from its cross-house counterpart **Subtle Chain**
  (Shadows), which keeps the original opponent-only discard.
- **Speed Sigil** is limited to one copy per deck. KeyForge sets no such limit;
  Vactrol caps it at one to bound its first-creature-of-the-turn ready loop.
- **Toad** is `Connected` rather than KeyForge's `Special`: it is kept out of the
  pool and instead pulled into **Xenos Bloodshadow**'s pod one for one, so a Toad
  only ever reaches a deck alongside the Bloodshadow it rides in with.

## Invented cards

A few cards exist only in Vactrol — no KeyForge printing. They fill a structural
gap the real game left open and carry provisional abilities that may be retuned.

- **Shard of Glory** (Saurian) and **Shard of Unity** (Star Alliance) complete the
  nine-House Shard cycle. Age of Ascension printed seven Shards, one per its Houses,
  but never a Saurian or Star Alliance Shard, so an errant pod of either House (or a
  legacy-drawn Shard in a deck of those Houses) could not complete the deck-wide
  `OnePerHouse` cluster ([ADR 0036](adr/0036-clusters-place-card-families-by-strategy.md)).
  Both are `Connected` Artifacts in the **Anomaly Expansion** reservoir set,
  matching the AoA Shards' shape (Item • Shard, one copy per deck, `Action: For each
friendly Shard, …`). Their effects — Shard of Glory exalts a friendly creature,
  Shard of Unity readies a friendly creature — are Vactrol inventions in the spirit
  of the existing Shards and are provisional; only their existence and cross-set
  cycle role are load-bearing.

## Mechanic rule changes

- **Key cheats purge themselves.** Any card that forges a key outside the normal
  start-of-turn step spends itself when it succeeds: a forge that actually lands
  purges the card that made it, gated with the `->` result arrow (rule 5) as
  `… forge a key … -> purge <self>`. A forge that is barred or unaffordable does
  not purge — the card stays. This keeps a single key cheat from looping a body
  back into play to forge again and again in one turn (Chota Hazri regrowth), while
  leaving its ordinary tactical use intact. The rule lives on the `ForgeKey` effect
  node, so every forge card carries it uniformly (Data Forge, Imperial Forge,
  Forging an Alliance, Key of Darkness, Key Charge, Chota Hazri, The Colosseum,
  Nightforge, Key Abduction, Triumph, Might Makes Right, [REDACTED], Epic Quest,
  Obsidian Forge). Because the forge is now its own cost, the optional wrapper is
  dropped from the two key cheats that carried one — **Nightforge** and **Obsidian
  Forge** are mandatory when affordable (the `you may` is gone) — and the cards
  that used to **destroy** themselves on forging (**Epic Quest**, **[REDACTED]**,
  **Obsidian Forge**) now **purge** instead, removing them from the game rather
  than sending them to the discard pile where they could return.

- **A bonus icon's source is the card that carries it.** When a bonus icon
  resolves (Æmber, Capture, Damage, Draw), KeyForge treats the game itself as the
  source of the effect; Vactrol treats the card the icon is printed on as the
  source. This only matters for the rare card that reads the source of a bonus
  icon's damage, and it makes the log read naturally ("Splinter deals 1 bonus
  damage to …"). See [ADR 0041](adr/0041-bonus-icons-are-the-primitive.md).

- **A creature stops resolving its bonus icons when it leaves play.** KeyForge
  resolves every bonus icon on a played card even after the card has left play (a
  Damage icon that destroys its own creature still resolves the icons below it).
  In Vactrol, once a creature leaves play — destroyed by a constant ability before
  its icons resolve, or by its own Damage icon aimed at itself mid-resolution — its
  remaining icons do not resolve. Only creatures are gated this way: an artifact
  deals its bonus damage to other creatures, an upgrade resolves after it attaches,
  and an action is never in play, so their icons always all resolve.

- **A card may bar specific bonus-icon kinds from Enhance.** KeyForge has no such
  rule. Vactrol lets a card carry `WithoutEnhancement(kinds…)` so deck generation
  never lands those bonus-icon kinds on it, while other kinds may still land — for
  a bonus a card would only be weakened by (Effervescent Principle bars Capture).
  Rolled in per card as they are found.

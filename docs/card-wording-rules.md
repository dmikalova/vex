# Card wording rules

These are the house wording conventions for printed card text — one surface of the
controlled Rules voice ([ADR 0019](adr/0019-controlled-rules-voice.md); see the
root `AGENTS.md`). Every card's rendered text obeys the conventions below, so the
printed card maps cleanly onto the effect AST (each clause is a node) and can never
desync from behavior.

The conventions were originally derived by diffing a curated implementation set
against the original printed KeyForge text, but they now stand on their own as
house style — they apply to every set, not just the one they were first drawn
from. The pristine printed originals still live in
[callofthearchons.json](../internal/cards/provenance/callofthearchons.json), a flat
JSON array of every Call of the Archons card, for reference.

Some conventions are deliberate **divergences** from KeyForge — a renamed card
type, retired verbs, restructured abilities. Where Vactrol diverges, Vactrol wins.
Each divergence is called out inline below and cataloged in the
[Vactrol⇄KeyForge divergence register](keyforge-divergences.md), which holds the
precedence rule (fall back to the KeyForge Master Rulebook only for what Vactrol
has not decided).

Matching the printed KeyForge text is never a goal in itself and never a blocker.
Template consistency and simplicity across Vactrol outrank exact KeyForge wording:
a **meaning-preserving** reword that puts a card in the same voice and template as
the rest of Vactrol — or that lets a mechanic decompose into shared nodes — is
welcome, and needs no divergence-register entry (the register catalogs changes to
a rule or a name, not house-voice alignment).

The goal of the rewording is **one printed sentence per engine operation**: text
that maps cleanly onto the effect AST (each clause is a node), so the printed
card and the implementation never desync. Each rule below notes the intended AST
shape where relevant.

---

## 1. Keywords stand alone — reminder text is removed

Parenthetical reminder text is stripped; the keyword is the whole statement.

| Original                                                                                            | Curated     |
| --------------------------------------------------------------------------------------------------- | ----------- |
| `Elusive. (The first time this creature is attacked each turn, no damage is dealt.)`                | `Elusive.`  |
| `Skirmish. (When you use this creature to fight, it is dealt no damage in return.)`                 | `Skirmish.` |
| `Taunt. (This creature's neighbors cannot be attacked unless they have taunt.)`                     | `Taunt.`    |
| `Poison. (Any damage dealt by this creature's power during a fight destroys the damaged creature.)` | `Poison.`   |

Reminder text is a rendering concern (a keyword glossary), not part of the
ability tree. Keywords are modeled as flags/ratings on the card, so their meaning
lives in the engine, not the string. Several keywords on one card render as a
single comma-separated sentence — `Skirmish, Poison.`, not `Skirmish. Poison.` —
because they are one list of properties, not a sequence of statements.
(Affected: Drumble, Batdrone, "John Smyth", Bulleteye, Deipno Spymaster, Faygin,
Macis Asp, Dew Faerie, Chuff Ape, Champion Anaphiel, Ancient Bear, Briar
Grubbling, Way of the Bear, Way of the Wolf.)

### Articles in front of a computed noun

Whenever the noun in a phrase comes from a helper — a `Selection`'s `noun()`,
`e.noun()`, a `Target`'s noun form — build the article with `indefinite(noun)`
rather than concatenating `"a " + noun`. The card types include _artifact_ and
_action_, so a hardcoded `"a "` prints "a artifact" the moment someone points an
existing effect at a new type. `indefinite` is in `internal/engine/text.go`
alongside `plural` and `countNoun`, which exist for the same reason on the
number axis: never write `noun + "s"` or a `card(s)` placeholder by hand.

### Two wordings that are standardized, not left to the printed card

- **A leading duration clause takes a comma**: "for the remainder of the turn,
  this creature belongs to house Mars". The printed cards are split on this —
  Brain Stem Antenna omits the comma mid-sentence — so we always include it
  rather than making punctuation depend on sentence position.
- **`Damage` casing follows rule 15, not the individual card**: a number of
  damage dealt is always capitalized (`deals 5 Damage`, `deals +2 Damage`), and
  `deals no damage` stays lowercase because no icon is printed there.
- **Capitalize icon/magnitude terms and named keywords; card types are plain
  nouns.** The terms card text capitalizes are the ones backed by a printed
  icon or magnitude — `Æmber`, `Damage`, `Power` — and the named keywords
  (`Skirmish`, `Elusive`, …). Card types are not: `creature`, `artifact`,
  `upgrade`, and `tactic` (and their plurals) read as lowercase common nouns,
  exactly like the generic `card`. A card type is only capitalized when it opens
  a sentence (via the ordinary first-letter capitalization every line already
  gets — "Creatures cannot fight this creature"), never because it is a card
  type. This scope is card text only: the doc-comment `Type:` label, the rulebook
  term headwords (`Creature`, `Artifact`), and prose that discusses the types
  stay capitalized as the proper glossary/label names they are.

## 2. `gains` for everything grantable — retire `gets`

Every acquired stat modifier or keyword uses **`gains`**. `gets` is retired.

| Original                                       | Curated                                       |
| ---------------------------------------------- | --------------------------------------------- |
| `Each friendly creature gets +1 power.`        | `Each friendly creature gains +1 power.`      |
| `This creature gets +5 power.`                 | `This creature gains +5 power.`               |
| `Each of Bulwark's neighbors gets +2 armor.`   | `Each of Bulwark's neighbors gains +2 armor.` |
| `This creature gets +1 armor and gains taunt.` | `This creature gains +1 armor and taunt.`     |

A single `gains` covers both numeric modifiers and keywords, so one grant clause
can list both (`gains +1 armor and taunt`). Maps to a `Grant{...}` node whose
payload is a stat delta and/or keyword.

**Punctuation split (unchanged from prior convention):**

- Granting a **stat/keyword**: bare, no quotes — `gains +5 power`, `gains taunt`.
- Granting a **quoted ability**: comma + quotes — `gains, "Destroyed: …"`.

## 3. Self-destruction is `Destroy <self>`, never `Sacrifice`

| Original                                             | Curated                                            |
| ---------------------------------------------------- | -------------------------------------------------- |
| `Omni: Sacrifice Key to Dis. Destroy each creature.` | `Omni: Destroy Key to Dis. Destroy each creature.` |
| `Omni: Sacrifice Gorm of Omm. Destroy an artifact.`  | `Omni: Destroy Gorm of Omm. Destroy an artifact.`  |

`sacrifice` is collapsed into `destroy`; there is one destruction verb.
(Affected: Key to Dis, Combat Pheromones, Gorm of Omm, Epic Quest, Longfused
Mines, Nepenthe Seed.)

## 4. Card movement is `Put X …`, never `Return X …`

| Original                                                            | Curated                                                      |
| ------------------------------------------------------------------- | ------------------------------------------------------------ |
| `Return an enemy creature to its owner's hand.`                     | `Put an enemy creature into its owner's hand.`               |
| `Return Bad Penny to your hand.`                                    | `Put Bad Penny into your hand.`                              |
| `Return a creature from your discard pile to the top of your deck.` | `Put a creature from your discard pile on top of your deck.` |

One movement verb (`Put`) covers every destination: `into … hand(s)`,
`into … archives`, `on top of … deck`. Maps to a `Move{card, zone}` node.
(Affected: Fear, Arise!, Bad Penny, Faygin, Phoenix Heart, Grasping Vines, World
Tree, Nepenthe Seed.)

## 5. Result gates use `->`, replacing `If you do`

A conditional consequence that depends on the previous clause succeeding is
written with an arrow, not a follow-up sentence.

| Original                                                                                        | Curated                                                                              |
| ----------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| `Purge a creature from a discard pile. If you do, put a +1 power counter on Eater of the Dead.` | `Purge a creature from a discard pile -> give Eater of the Dead a +1 power counter.` |
| `Destroy a damaged creature. If you do, steal 1 Æmber.`                                         | `Destroy a damaged creature -> steal 1 Æmber.`                                       |
| `You may sacrifice another friendly creature. If you do, fully heal Chuff Ape.`                 | `You may destroy another friendly creature -> fully heal Chuff Ape.`                 |
| `Lose 1 Æmber. If you do, you may forge a key at current cost.`                                 | `Lose 1 Æmber -> forge a key at current cost.`                                       |

`A -> B` is a **gate**: attempt A, and only run B if A actually happened. Maps to
`Gate{attempt: A, then: B}`. This is distinct from an unconditional sequence
(`A. B.`) and from a state branch (`If <fact>, …`). A gate never takes an
`otherwise`.

Power counters use the verb `give`: `give <creature> a +1 power counter` (not
`put a +1 power counter on <creature>`). The holder is named — the source by name
for a self-counter (`give Eater of the Dead a +1 power counter`), or `that
creature` when a creature was just chosen.

## 6. Replacement effects front-load `instead`

When an effect replaces what _would_ happen, the `would` cue comes first and
`instead` heads the replacement clause — the twist is never buried at the end.

| Original                                                                                              | Curated                                                                                         |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `…any Æmber you would gain from reaping is stolen from your opponent instead.`                        | `…instead of gaining Æmber from reaping, steal the same amount.`                                |
| `Each Æmber that would be added to your opponent's pool is captured by Ether Spider instead.`         | `If Æmber would be added to your opponent's pool, instead Ether Spider captures it.`            |
| `Destroyed: Fully heal this creature and destroy Armageddon Cloak instead.`                           | `If this creature would be destroyed, instead fully heal it and destroy Armageddon Cloak.`      |
| `Destroyed: If you have any other creatures in play, instead of destroying Reassembling Automaton, …` | `If this creature would be destroyed and there is another friendly creature in play, instead …` |

Two shapes:

- **Passive event:** `If X would …, instead Y` (Ether Spider, Armageddon Cloak).
- **Named action:** `Instead of X, Y` (Dimension Door).

Maps to a `Replace{when, cond, with}` interceptor, distinct from a trigger. Note
Armageddon Cloak drops the `Destroyed:` trigger label — a destruction _save_ is a
replacement, not a post-destruction trigger. An Upgrade carries the `Replace` to
save its host; a creature carries it (in its own `Static`) to save itself
(Reassembling Automaton, Self-Bolstering Automata), reading in its own voice
rather than through the `This creature gains, "…"` framing. A conditional save
folds its condition into the `would` clause (`If this creature would be destroyed
**and** there is another friendly creature in play, instead …`).

Because the replacement stands in _before_ the destruction is enrolled, a saved
creature never counts as destroyed — it does not bump the enemy-creatures-destroyed
tally, the creatures-destroyed-this-way count, or fire "after … destroyed"
reactions. (Authoring a destruction save as a `Destroyed:` trigger — as the Automata
once were — wrongly counted the creature as destroyed; the `Replace` path fixes
that.)

## 7. Self-reference by card name

A creature refers to itself by name, not "this creature" (outside granted
abilities).

| Original                 | Curated                            |
| ------------------------ | ---------------------------------- |
| `Play: Capture 3 Æmber.` | `Play: Charette captures 3 Æmber.` |

The rule by context:

- **On a creature, where the creature is the subject/holder of the effect** — it
  names itself: `Charette captures 3 Æmber`, `Valdr deals +2 Damage …`,
  `fully heal Chuff Ape`. Capture and static modifiers bind to that specific
  creature, so the name carries meaning.
  - **Capture always names its target:** the
    target is now always specified (an unspecified target used to default to the
    source creature), so the capturing creature is stated explicitly rather than
    left implicit.
- **On a creature, for a one-shot action with no lingering attachment to the
  source** — imperative, no name: `Deal 1 Damage to each enemy creature`,
  `Ready and fight with a friendly creature`.
- **On an upgrade or inside a granted quoted ability** — use `this creature` /
  `its`, since the host can vary: `Fully heal this creature`, `Put this creature
into its owner's archives`.

## 8. `Choose one` renders as a bulleted list

| Original                                                                         | Curated                                                                                 |
| -------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `Choose one: Archive a card, or, for each archived card you have, gain 1 Æmber.` | `Choose one:`<br>`- Archive a card`<br>`- Gain 1 Æmber for each card in your archives.` |
| `Choose one: destroy each Dis creature, or gain 1 Æmber.`                        | `Choose one:`<br>`- Destroy each Dis creature`<br>`- Gain 1 Æmber.`                     |

Each option is a list item; maps to `ChooseOne{options: [...]}`.

## 9. Front-load `for each` / iteration — subject → effect

Counting and iteration clauses lead the sentence, so a card reads **forward**:
the subject (what you count) first, then the effect it drives — the reader never
has to run to the end for the multiplier and jump back.

| Original                                                          | Curated                                                            |
| ----------------------------------------------------------------- | ------------------------------------------------------------------ |
| `Gain 1 Æmber for each forged key your opponent has.`             | `For each forged key your opponent has, gain 1 Æmber.`             |
| `gain 1 Æmber each time you play a creature.`                     | `each time you play a creature, gain 1 Æmber.`                     |
| `Gain 1 Æmber for each creature healed this way.`                 | `For each creature healed this way, gain 1 Æmber.`                 |
| `Deal 1 Damage to a creature for each friendly creature in play.` | `For each friendly creature in play, deal 1 damage to a creature.` |

The loop header precedes the body, matching `ForEach{source, body}`. In the
engine this is the `Count` interface: any effect carrying a `Per Count` renders
it as a leading `for each <CountText>, <body>` clause (via the shared `forEach`
helper), so `GainAember`, `DealDamage`, and every future `Per`-scaled effect
front-load automatically.

## 10. Global effects are granted abilities

A "for every creature" rule is expressed as each creature _gaining_ the ability,
rather than a free-floating static/trigger.

| Original                                                                                                                                                       | Curated                                                                                                                                                       |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `When a creature would enter a discard pile from play, it is purged instead.`                                                                                  | `Each creature gains, "Destroyed: purge this creature."`                                                                                                      |
| `Before a creature fights, discard the top card of its controller's deck. If the discarded card is of the active house, exhaust that creature with no effect.` | `Each creature gains, "Before Fight: Discard the top card of its controller's deck. If the discarded card is of the active house, the fight does not occur."` |

Reuses the grant + quoted-ability machinery instead of a bespoke global rule.

## 11. `this way` for back-references

Effects that reference the cards/creatures they just touched use `this way`
(e.g. `discarded this way`, `healed this way`, `destroyed this way`), and
threshold self-checks use `fewer than N … this way`
(Bonkers Killing Machine). This is the standard back-reference idiom.

## 12. Terminology adopted from KeyForge keywords

Longhand descriptions are replaced by the canonical keyword/verb.

| Original                                                                                  | Curated                                                             | Keyword          |
| ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------- | ---------------- |
| `Place 2 Æmber from the common supply on an enemy creature.`                              | `Exalt an enemy creature 2 times.`                                  | Exalt            |
| `This creature gains, "You may use this creature as if it belonged to the active house."` | `This creature gains versatile.`                                    | versatile        |
| `This creature belongs to all houses.`                                                    | `This creature gains versatile.`                                    | versatile        |
| `Omni: Destroy Gorm of Omm. Destroy an artifact.`                                         | `Versatile.`<br>`Action: Destroy Gorm of Omm. Destroy an artifact.` | Omni → versatile |

**Omni → Versatile + `Action:`.** This is the one keyword swap that changes a
card's _structure_, so it is worth spelling out. An `Omni:` ability may be used on
your turn **whether or not the card is in your active house**; a plain `Action:`
ability may only be used when the card **is** in the active house. **Versatile** is
exactly the keyword that lets a card be used as if it belonged to the active house.
So an Omni card is re-expressed as two pieces:

1. the card **gains `Versatile`** (a keyword on the card), and
2. its `Omni:` ability is rewritten verbatim as an `Action:` ability.

For almost every card these two are behaviorally identical — using the card
exhausts it, and the Versatile grant removes the active-house restriction — so the
engine has **no separate Omni trigger**. Author it with
`card.WithKeywords(card.Keyword.Versatile)` plus a `card.WithAbility(card.Trigger.Action, …)`.

Concretely, `Omni: Destroy Gorm of Omm. Destroy an artifact.` becomes:

```text
Versatile.
Action: Destroy Gorm of Omm. Destroy an artifact.
```

(Affected: Key to Dis, Combat Pheromones, Epic Quest, Gorm of Omm, Deipno
Spymaster, Longfused Mines, Nepenthe Seed.)

## 13. Cancellation primitive: `the fight does not occur`

Fight prevention is a terse cancel, not a longhand "exhaust with no effect".

| Original                                 | Curated                      |
| ---------------------------------------- | ---------------------------- |
| `…exhaust that creature with no effect.` | `…the fight does not occur.` |

Maps to a `CancelFight` result; relies on the creature already being exhausted at
use-time.

## 14. Renames

| Original name         | Curated name            |
| --------------------- | ----------------------- |
| Crazy Killing Machine | Bonkers Killing Machine |

## 15. `Damage` vs `damage` (damage-icon casing)

Capitalize **`Damage`** only where the game _deals_ it — the printed card shows a
damage icon there: `Deal 2 Damage`, `+2 Damage`. Every other
use is lowercase.

| Context                       | Casing   | Example                                                  |
| ----------------------------- | -------- | -------------------------------------------------------- |
| Dealing damage (icon)         | `Damage` | `Deal 3 Damage to each creature.`                        |
| Healing / referring to damage | `damage` | `Heal 3 damage from a creature.` / `a damaged creature.` |

The number is what carries the icon, so a creature's fight damage is capitalized
when it names one (`Bruiser deals 5 Damage when fighting.`) and lowercase when it
does not (`Spider deals no damage when fighting.`).

## 16. Spelling, qualifiers, and referents

- **`Æmber`** — never `Aember`. The engine renders the resource as `Æmber`
  everywhere (card text, log, rulebook), so that is the one spelling in
  player-facing prose. Go identifiers keep the ASCII `Aember` by convention:
  `Æ` is a legal Go identifier rune (a Unicode uppercase letter), so `Æmber`
  would compile, but ASCII keeps identifiers typeable and greppable. This rule
  governs text, not code.
- **`<Trait>` references read like house references** — a trait filter names the
  trait as a plain adjective on the noun (`Scientist creature`, `Knight
creature`), exactly as houses do (`Mars creature`, `Sanctum creature`). Neither
  carries a `trait` qualifier.
- **`during their next turn`** — not `on their next turn`.
- **`the chosen creature`** — the standard referent for a creature just chosen
  (`Choose a creature. … the chosen creature`).

## 17. Name the source zone for card-movement effects

Effects that move a card between zones name the zone explicitly, so the text is
unambiguous about _where from_.

| Original          | Curated                          |
| ----------------- | -------------------------------- |
| `Archive a card.` | `Archive a card from your hand.` |

Applies to archiving (`from your hand`, `the top card of your deck`) and any
future zone-to-zone movement (hand, deck, discard, archives). `Archive` alone
never implies the hand.

The same holds for the Æmber **pool** a capture draws from: a capture always
names it, even where the printed card leaves it implied.

| Original                               | Curated                                                   |
| -------------------------------------- | --------------------------------------------------------- |
| `a friendly creature captures 1 Æmber` | `a friendly creature captures 1 Æmber from your opponent` |

## 18. Æmber a player owes another player is `give`, never `pay`

One verb covers Æmber changing hands between players: **`give`**. `pay` (the
printed wording on toll cards) is retired so every player-to-player transfer reads
the same way, matching Interdimensional Graft's "they must give you their
remaining Æmber".

| Original                                                           | Curated                                                              |
| ------------------------------------------------------------------ | -------------------------------------------------------------------- |
| `Your opponent must pay you 1 Æmber in order to play an artifact.` | `In order to play an artifact, your opponent must give you 1 Æmber.` |
| `…they must pay you their remaining Æmber.`                        | `…they must give you their remaining Æmber.`                         |

The engine keeps the mechanic named `Toll` (the thing a card charges), but its
rendered text says `give`. (Affected: Customs Office, Tentacus,
Interdimensional Graft.)

An `in order to` requirement states the requirement **up front**: `In order to do
X, you must do Y` — the price comes after the thing it gates, so the reader learns
what is restricted before what it costs, rather than having to backtrack. This
applies to every toll (`In order to use an artifact, your opponent must give you 1
Æmber.`) and to a play cost (`In order to play Truebaru, you must lose 3 Æmber.`).

## 19. The `action` card type is renamed `Tactic`

KeyForge overloads "Action": it is both a **card type** (a one-shot card you play
from hand) and an **ability** (`Action:` on a ready creature or artifact). To keep
the two distinct, the card **type** is renamed **`Tactic`** (rendered `Type:
Tactic`), while the **`Action:` ability** keeps its printed wording untouched.

| Context                        | Wording          |
| ------------------------------ | ---------------- |
| Card type (was "action")       | `Tactic`         |
| Ability on a creature/artifact | `Action:` (kept) |

This is a naming choice, not a text rewrite: the `Action:` prefix on every card's
ability line is unchanged.

## 20. One-time house reassignment, not `while under your control`

A few cards take control of a card and then say it belongs to a house for as long
as you control it — e.g. Sneklifter's printed _"Take control of an enemy artifact.
While under your control, if it does not belong to one of your three houses, it is
considered to be of house Shadows."_ The `while under your control` clause is
**dropped** in favor of a one-time reassignment, **evaluated once, when the effect resolves**:

> Take control of an enemy artifact. If it does not belong to a house on your
> identity then it belongs to house Shadows.

The continuous "while under your control" form needs a kind of dynamic memory that
re-checks the card's houses every time control changes; the one-time form does not.
It is also **slightly more interesting**: the reassignment sticks with the card
(until it leaves play), so if an opponent later takes it back they can be stuck
with a card of a house they may not have. "A house on your identity" is one of the
controller's three deck houses. (Affected: Sneklifter.)

## 21. A deferred play permission becomes an immediate play

Several cards grant a play that may be taken any time later in the turn — e.g.
Phase Shift's printed _"Play: You may play one non-Logos card this turn."_ The
`this turn` window is **dropped** and the play happens **immediately, as part of
resolving the card**:

> Play: Play a non-Logos card.

The deferred form needs turn-scoped memory of an unused permission — a counter
that must be armed, spent by an unrelated later play, and cleared at end of turn —
plus a prompt on every subsequent play asking which allowance to spend. The
immediate form needs none of that: the effect resolves, the player picks a card
from hand, and it is played there and then.

The general rule: **`... you may play an X card this turn` renders as `Play an X
card`.** Both `you may` and `this turn` go — a play with no legal card in hand
simply does nothing, so the permission needs no explicit opt-out. The rule extends
to a grant that may be spent on **playing or using** a card. Taber's _"You may
play or use a non-Star Alliance card this turn"_ becomes:

> Fight/Reap: Play or use a non-Star Alliance card.

The play-or-use resolves now against the cards in hand and in play. (Affected:
Phase Shift, Kirby, Taber.)

## 22. A number-only branch collapses to `or <alt> if <cond>`

A two-armed `If <cond>, <A>. Otherwise, <B>.` whose arms are the **same verb**
differing only in a number is a disguised parameter, not a real fork. It collapses
to one verb that states the base value and appends the exception as a trailing
rider — a linear sentence read in order, with no branch.

| Original                                                                                                                | Curated                                                                      |
| ----------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `If your opponent has 7 Æmber or more, steal 2 Æmber. Otherwise, steal 1 Æmber.`                                        | `Steal 1 Æmber, or 2 if your opponent has 7 Æmber or more.`                  |
| `If your opponent has no Æmber, forge a key at +2 Æmber current cost. Otherwise, forge a key at +6 Æmber current cost.` | `Forge a key at +6 Æmber current cost, or +2 if your opponent has no Æmber.` |

The base value is the one that applies when the condition is **absent**; the rider
names the value the condition switches to. Maps to a single effect node whose
amount carries an `OrAmount{Amount, When}` — the alternate value and its guard —
which renders the `, or <alt> if <cond>` tail and picks the value at resolve time.
No `Else` arm exists.

`Otherwise` survives only for a genuine two-way branch — arms that are **different
actions**, not the same verb with a different number. Vespilon Theorist keeps it
(archive the revealed card and gain 1 Æmber, or discard it), because archive and
discard are different actions on different zones. (Affected: Ronnie Wristclocks,
Key of Darkness.)

---

## 23. "That creature and each creature that shares a trait" folds into the trait sweep

A card that chooses a creature and then destroys "that creature **and** each
creature that shares a trait with it" names the chosen creature twice: once as the
subject and again, implicitly, in the trait sweep — a creature always shares a
trait with itself. The curated text drops the redundant "that creature and" and
states only the sweep, so `ChooseCreatureThen{Then: Destroy{... .SharingTrait()}}`
renders `choose a creature - destroy each creature that shares a trait with it`.
The chosen creature is still destroyed, because it matches its own trait filter.

| Original                                                                                  | Curated                                                                  |
| ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| `Choose a creature. Destroy that creature and each creature that shares a trait with it.` | `Choose a creature - destroy each creature that shares a trait with it.` |

(Affected: Extinction.)

---

## 24. "Do X. If <cond>, do Y instead" becomes a choose-then conditional

A card that names a default action on a creature and then overrides it with an
exception ("Stun an enemy creature. If that creature was already stunned, destroy
it instead.") is authored as a `ChooseCreatureThen` whose `Then` is a
`Conditional`: the exception is the `Then` branch and the default is the `Else`.
The curated text therefore leads with the choice and the condition rather than the
default action, because the condition must be read against the creature's state
_before_ the default would change it — stunning first would make "already stunned"
always true.

| Original                                                                            | Curated                                                                                            |
| ----------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `Stun an enemy creature. If that creature was already stunned, destroy it instead.` | `Choose an enemy creature - if that creature was already stunned, destroy it. Otherwise, stun it.` |

(Affected: 1-2 Punch.)

---

## 25. A trailing `if <cond>` condition moves to the front

A card that states an effect and then qualifies it with a trailing condition
("Gain 3A if you control creatures from 3 different houses.") is authored as a
`Conditional`, which always renders the condition first. The curated text leads
with `If <cond>, <effect>` so every gated effect reads the same way, matching the
condition-leading form KeyForge already uses elsewhere.

| Original                                                    | Curated                                                           |
| ----------------------------------------------------------- | ----------------------------------------------------------------- |
| `Gain 3A if you control creatures from 3 different houses.` | `If you control creatures from 3 different houses, gain 3 Æmber.` |

(Affected: Prince Derric, Unifier.)

---

## 26. Use-conditions read from the controller, and an upgrade names its host

A card that may not be used until its controller has met a condition is authored
as a `Restrictions.UseCondition`. On a creature or artifact it renders `You
cannot use this card unless <condition>.`; on an upgrade — where "this card" is
the upgrade, not the thing being used — the subject becomes the host creature:
`This creature cannot be used unless <condition>.` The condition itself always
reads from the controller's point of view (`you have discarded a card from your
hand this turn`), replacing KeyForge's third-person `its controller has
discarded a card this turn`.

| Original (upgrade)                                                                   | Curated                                                                                   |
| ------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------- |
| `This creature cannot be used unless its controller has discarded a card this turn.` | `This creature cannot be used unless you have discarded a card from your hand this turn.` |

(Affected: Earthbind. Giant Sloth is the creature form.)

---

## 27. "Gain that much Æmber" becomes a `for each Æmber bonus` clause

A card that destroys a card and then gains Æmber equal to that card's Æmber bonus
is authored as a `Destroy` followed by a `GainAember{Per: AemberBonusOf{Target:
Triggering}}`, which the engine renders as an explicit per-pip clause reading the
destroyed card in context (`it`) rather than the printed back-reference `you gain
that much Æmber`. The rendered form names what is being counted, so behavior and
text cannot drift.

| Original                                                                              | Curated                                                          |
| ------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| `Destroy an artifact. If that artifact had an Æmber bonus, you gain that much Æmber.` | `Destroy an artifact. For each Æmber bonus on it, gain 1 Æmber.` |

(Affected: Rustgnawer.)

---

## 28. A turn `step` is named a `phase`

A named part of the turn is a **phase**, never a **step** — Vactrol has one word
for a segment of the turn ([ADR 0012](adr/0012-first-class-turn-phases.md) makes
the turn phases first-class). Card and rules text that refers to the `"forge a
key"` or `"draw cards"` part of the turn writes `phase`, where KeyForge writes
`step`.

| Original                                | Curated                                  |
| --------------------------------------- | ---------------------------------------- |
| `During your "draw cards" step, …`      | `During your "draw cards" phase, …`      |
| `skips the "forge a key" step during …` | `skips the "forge a key" phase during …` |
| `You skip your "forge a key" step.`     | `You skip your "forge a key" phase.`     |

(Affected: Streke, Mother, Succubus, The Howling Pit, Miasma.)

---

## 29. Fight timing is named `in a fight with`

The moment a fight exchanges power damage between two combatants is written as
`in a fight with <creature>`, Vactrol's one phrase for that timing window.
KeyForge names it several ways — `fighting`, `while fighting`, `during a fight` —
and Vactrol collapses them to one, so a trigger that fires on the destruction of
the other combatant reads the same on every card.

| Original                                          | Curated                                                  |
| ------------------------------------------------- | -------------------------------------------------------- |
| `After a creature is destroyed fighting Krump, …` | `After a creature is destroyed in a fight with Krump, …` |

(Affected: Krump.)

---

## 30. A count cap is dropped — no `(to a maximum of N)` clause

A `for each` count that KeyForge caps prints the cap in parentheses — `gain 1
Æmber (to a maximum of 6) for each house …`. Vactrol drops the cap entirely: the
count is uncapped and the clause is not rendered, so the "for each" sentence
reads clean. This is a deliberate rules divergence, not just a wording change
(see `keyforge-divergences.md`).

| Original                                                                             | Curated                                                          |
| ------------------------------------------------------------------------------------ | ---------------------------------------------------------------- |
| `for each house represented among cards in play, gain 1 Æmber, to a maximum of 6`    | `for each house represented among cards in play, gain 1 Æmber`   |
| `steal 1 Æmber for each house represented among enemy creatures (to a maximum of 3)` | `steal 1 Æmber for each house represented among enemy creatures` |

(Affected: Free Markets, Forging an Alliance, Trust No One, and the reprints
that share their text.)

---

## 31. "Deal X, deal Y instead if …" is a choose-then-conditional

A card that deals one amount of damage to a chosen creature and a larger amount
instead when that creature meets a condition is authored as a `ChooseCreatureThen`
whose `Then` is a `Conditional` — the condition picks which amount is dealt —
rather than a bespoke damage-boost strategy on `DealDamage`. Compound conditions
compose (`Or{ItIsOfTrait{…}, ItHasAember{}}`) instead of baking each combination
into a one-off condition, and the rendered form names both amounts and the branch
plainly.

| Original                                                                                              | Curated                                                                                                                        |
| ----------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `Deal 2 damage to a creature. Deal 6 damage instead if it is a Dinosaur creature or has Æmber on it.` | `Choose a creature - if it is a Dinosaur creature or it has Æmber on it, deal 6 damage to it. Otherwise, deal 2 damage to it.` |

(Affected: Guji Dinosaur Hunter.)

---

## 32. Timed "cannot use" restrictions share one phrasing

A timed effect that bars a player from acting is one `Restrict` node —
`cannot use creatures to fight`, `cannot use creatures to reap`, or the broadest
`cannot use any cards` — and every form ends with the same duration phrase:
`during <their/your> next turn` (OpponentNextTurn) or `for the remainder of the
turn` (RemainderOfPlayerTurn). KeyForge writes the broad current-turn form as "You
cannot use cards this turn"; Vactrol harmonizes it to "you cannot use **any** cards
**for the remainder of the turn**" so the noun ("any cards") matches the
next-turn form and the duration matches the reap and play bars.

A house scope narrows the reaping noun the same way, so a house-limited reap bar
reads in the shared voice rather than a subject-first sentence of its own. KeyForge
writes Seismo-entangler as "During your opponent's next turn, creatures of the
chosen house cannot be used to reap"; Vactrol harmonizes it to "your opponent
cannot use **creatures of the chosen house** to reap during their next turn".

| Original                                                                                  | Curated                                                                                  |
| ----------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `You cannot use cards this turn.`                                                         | `You cannot use any cards for the remainder of the turn.`                                |
| `During your opponent's next turn, creatures of the chosen house cannot be used to reap.` | `Your opponent cannot use creatures of the chosen house to reap during their next turn.` |

(Affected: United Action, Seismo-entangler.)

---

## Deliberate rule changes (not just wording)

A few cards diverge from KeyForge in ways that affect the rules, not just
phrasing. Those changes, and the precedence rule that governs every divergence,
live in the [Vactrol⇄KeyForge divergence register](keyforge-divergences.md).

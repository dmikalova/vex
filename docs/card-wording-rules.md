# Card wording rules

This document characterizes the wording transformations applied when curating the
minimal implementation set. It is derived by diffing the curated set against the
original printed card text.

- **Source of truth for originals:** [callofthearchons.json](../internal/cards/provenance/callofthearchons.json) — the pristine
  printed text of every Call of the Archons card (flat JSON array).
- **Curated set:** [cota-top.jsonc](cota-top.jsonc) — the minimal set selected for
  implementation, keyed by house, with reworded text.
- **Regenerate the diff:** `python3 scripts/diff-cota.py` (prints every
  old→new text change; used to build and maintain this document).

Of 124 curated cards, 77 had their text reworded and 1 was renamed. No stats
(power, armor, Æmber, house, type, rarity) were changed.

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
lives in the engine, not the string. (Affected: Drumble, Batdrone, "John Smyth",
Bulleteye, Deipno Spymaster, Faygin, Macis Asp, Dew Faerie, Chuff Ape, Champion
Anaphiel, Ancient Bear, Briar Grubbling, Way of the Bear, Way of the Wolf.)

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
| `Destroy a damaged creature. If you do, steal 1 Aember.`                                        | `Destroy a damaged creature -> steal 1 Aember.`                                      |
| `You may sacrifice another friendly creature. If you do, fully heal Chuff Ape.`                 | `You may destroy another friendly creature -> fully heal Chuff Ape.`                 |
| `Lose 1 Aember. If you do, you may forge a key at current cost.`                                | `Lose 1 Aember -> forge a key at current cost.`                                      |

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

| Original                                                                                       | Curated                                                                                    |
| ---------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| `…any Aember you would gain from reaping is stolen from your opponent instead.`                | `…instead of gaining Æmber from reaping, steal the same amount.`                           |
| `Each Aember that would be added to your opponent's pool is captured by Ether Spider instead.` | `If Aember would be added to your opponent's pool, instead Ether Spider captures it.`      |
| `Destroyed: Fully heal this creature and destroy Armageddon Cloak instead.`                    | `If this creature would be destroyed, instead fully heal it and destroy Armageddon Cloak.` |

Two shapes:

- **Passive event:** `If X would …, instead Y` (Ether Spider, Armageddon Cloak).
- **Named action:** `Instead of X, Y` (Dimension Door).

Maps to a `Replace{when, with}` interceptor, distinct from a trigger. Note
Armageddon Cloak drops the `Destroyed:` trigger label — a destruction _save_ is a
replacement, not a post-destruction trigger.

## 7. Self-reference by card name

A creature refers to itself by name, not "this creature" (outside granted
abilities).

| Original                  | Curated                             |
| ------------------------- | ----------------------------------- |
| `Play: Capture 3 Aember.` | `Play: Charette captures 3 Aember.` |

The rule by context:

- **On a creature, where the creature is the subject/holder of the effect** — it
  names itself: `Charette captures 3 Aember`, `Valdr deals +2 Damage …`,
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

| Original                                                                          | Curated                                                                                  |
| --------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `Choose one: Archive a card, or, for each archived card you have, gain 1 Aember.` | `Choose one:`<br>`- Archive a card`<br>`- Gain 1 Aember for each card in your archives.` |
| `Choose one: destroy each Dis creature, or gain 1 Aember.`                        | `Choose one:`<br>`- Destroy each Dis creature`<br>`- Gain 1 Aember.`                     |

Each option is a list item; maps to `ChooseOne{options: [...]}`.

## 9. Front-load `for each` / iteration — subject → effect

Counting and iteration clauses lead the sentence, so a card reads **forward**:
the subject (what you count) first, then the effect it drives — the reader never
has to run to the end for the multiplier and jump back.

| Original                                                          | Curated                                                             |
| ----------------------------------------------------------------- | ------------------------------------------------------------------- |
| `Gain 1 Aember for each forged key your opponent has.`            | `For each forged key your opponent has, gain 1 Aember.`             |
| `gain 1 Aember each time you play a creature.`                    | `each time you play a creature, gain 1 Aember.`                     |
| `Gain 1 Aember for each creature healed this way.`                | `For each creature healed this way, gain 1 Aember.`                 |
| `Deal 1 Damage to a creature for each friendly creature in play.` | `For each friendly creature in play, deal 1 damage to a creature.`  |

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
| `Place 2 Aember from the common supply on an enemy creature.`                             | `Exalt an enemy creature 2 times.`                                  | Exalt            |
| `This creature gains, "You may use this creature as if it belonged to the active house."` | `This creature gains versatile.`                                    | versatile        |
| `This creature belongs to all houses.`                                                    | `This creature gains versatile.`                                    | versatile        |
| `Omni: Destroy Gorm of Omm. Destroy an artifact.`                                         | `Versatile.`<br>`Action: Destroy Gorm of Omm. Destroy an artifact.` | Omni → versatile |

**Omni** abilities are re-expressed as **Versatile** plus an `Action:` ability.
Under the simplified rule, Omni just means _usable on your turn regardless of the
active house_ — exactly what Versatile grants — so the card gains `Versatile` and
its ability is written as `Action:`. (Affected: Key to Dis, Combat Pheromones,
Epic Quest, Gorm of Omm, Deipno Spymaster, Longfused Mines, Nepenthe Seed.)

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

## 16. Spelling, qualifiers, and referents

- **`Aember`** — never `Æmber`.
- **`<Trait> trait`** — trait references carry the `trait` qualifier
  (`Scientist trait creature`, `Knight trait creature`); house references never do
  (`Mars creature`, `Sanctum creature`).
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

## 18. Æmber a player owes another player is `give`, never `pay`

One verb covers Æmber changing hands between players: **`give`**. `pay` (the
printed wording on toll cards) is retired so every player-to-player transfer reads
the same way, matching Interdimensional Graft's "they must give you their
remaining Æmber".

| Original                                                          | Curated                                                          |
| ----------------------------------------------------------------- | ---------------------------------------------------------------- |
| `Your opponent must pay you 1 Aember in order to play an artifact.` | `Your opponent must give you 1 Aember in order to play an artifact.` |
| `…they must pay you their remaining Aember.`                      | `…they must give you their remaining Aember.`                    |

The engine keeps the mechanic named `Toll` (the thing a card charges), but its
rendered text says `give`. (Affected: Customs Office, Tentacus,
Interdimensional Graft.)

## 19. The `action` card type is renamed `Tactic`

KeyForge overloads "Action": it is both a **card type** (a one-shot card you play
from hand) and an **ability** (`Action:` on a ready creature or artifact). To keep
the two distinct, the card **type** is renamed **`Tactic`** (rendered `Type:
Tactic`), while the **`Action:` ability** keeps its printed wording untouched.

| Context                              | Wording           |
| ------------------------------------ | ----------------- |
| Card type (was "action")             | `Tactic`          |
| Ability on a creature/artifact       | `Action:` (kept)  |

This is a naming choice, not a text rewrite: the `Action:` prefix on every card's
ability line is unchanged.

---

## Deliberate rule changes (not just wording)

A few cards were changed in ways that affect the rules, not just phrasing — to
**simplify** the card down to base rules text, to make it **slightly more
interesting**, or to bring it in line with **modern errata**:

- **Charge!** buffs all creatures, not just ones played this turn (dropped
  `you play`).
- **Imperial Traitor** `Reveal` (not `Look at`) — modern wording.
- **Ganger Chieftain** / **Biomatrix Backup** made mandatory (dropped `you may`).
- **Hypnotic Command** leans on the base rule that the **active player makes all
  decisions**, so `an enemy creature captures …` needs no explicit `choose`.

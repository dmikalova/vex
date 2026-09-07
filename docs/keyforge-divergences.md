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

## Per-card rule changes

A few cards were changed in ways that affect the rules, not just phrasing. Each
change simplifies the card toward base-rules text, makes it slightly more
interesting, or brings it in line with modern errata.

- **Charge!** buffs all creatures, not just ones played this turn. The `you play`
  clause is dropped.
- **Imperial Traitor** reads `Reveal`, not `Look at`. This is the modern wording.
- **Ganger Chieftain** and **Biomatrix Backup** are mandatory. The `you may`
  clause is dropped.
- **Hypnotic Command** leans on the base rule that the active player makes all
  decisions, so `an enemy creature captures …` needs no explicit `choose`.
- **Phase Shift** and **Kirby** play their off-house card immediately rather than
  granting a permission for later in the turn (rule 21).
- **Trust No One** is a `Choose one:` rather than a forced `If … Otherwise …`. The
  conditional branch ("if there are no friendly creatures in play, steal 1 Æmber
  per house among enemy creatures") still gates on the empty board, so it does
  nothing when you control a creature — a rational player picks the flat "steal 1
  Æmber" then, reproducing the original outcome, while the choice frame reads
  cleaner. Only a forced branch whose gated arm is a strict bonus over a safe
  fallback converts this way; most `If … Otherwise …` cards (random reveals,
  target-dependent or whose-turn conditions) do not.

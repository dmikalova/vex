# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Cite the ADR or doc that decided an item where one exists.

## Zone-movement consolidation (ADR 0031)

Fold the Archive/Discard/Purge/Shuffle/Put families onto one unexported movement
mechanism (source zone + selection + destination), with the KeyForge-verb sugar
delegating to it.

- Fuse the `DamageThen` / `DamageThenIfDestroyed` / `DamageThenIfSurvives` trio
  into `DealDamage{Then, After}` (`After` ∈ Always / IfDestroyed / IfSurvives).
  ~11 card uses across all sets; preserve the pre-damage neighbor snapshot the
  IfDestroyed variant relies on.
- Absorb `PurgeFromHand` / `PurgeCreatureFromHand` / `PurgeEachFromHand` into
  `PurgeCard` (House / ExceptHouse / All / bind-`it` fields).
- Collapse `ArchiveTopOfDeck` / `ArchiveTopOfDiscard` into `ArchiveTop{From}`, and
  the top/random discard variants into `Discard{Zone, Player, Amount, Random}`.

## Naming / modeling migrations (ADR 0031)

- `PoolAember{Player, Is, Amount}` replacing the `OpponentAember` / `YourAember`
  mirror pair (~21 uses across all sets; keep the relative `MoreThanYou` /
  `MoreThanOpponent` comparisons, which compare two pools).
- Invert `ItIsOffIdentity` (Sneklifter) and `AemberBonusDestroyed` (Rustgnawer)
  into positive conditions reading `ctx.It` — "does the creature it fought have
  Æmber bonus icons?" / "is it of one of your identity's houses?".
- Standardize the `…ThisWay` produced-tally count names onto one `Produced`-keyed
  convention.
- Damage modification as a replacement on a "damage about to be dealt" event,
  retiring `TakesExtraDamage` (extends the `Instead`/`Replace` spine like Po's
  Pixies). Lower priority: `EventCreatureTakesDamage` + `lastingExtraDamage`
  already add to the pending amount, so this is a re-homing, not a behavior fix.

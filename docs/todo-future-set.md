# Future-set backlog

Work that is real and decided, but that has **no consuming card in an implemented
set yet** — so it waits for the set that first introduces a card that needs it.
This is the sibling of [todo-agent.md](todo-agent.md): that file is work to do
now; this file is work parked against a set that has not been stood up.

Rules:

- **When you begin implementing (or stubbing) a set, scan this file first.** If an
  item names that set — or a card that set introduces — fold it into the run and
  build the primitive alongside its first real consumer, so a live card pins its
  shape and a card test covers it. The `implement-cards` and `stub-cards` skills
  both point here for exactly this reason.
- Each item carries the **decision, not just the task** — the chosen behavior, the
  cards affected, and the design already settled — so a future agent who was not in
  the conversation can build it without re-deriving it.
- **When the work lands, delete the item** (do not mark it done), same as
  [todo-agent.md](todo-agent.md). This file only ever shows parked work.

## Result-set reference — "choose one of the cards this effect just moved"

**Trigger set:** whichever set first stands up a consuming card — Junk Restoration
(Æmber Skies / Grim Reminders / Menagerie) or Haunting Measures (Draconian
Measures). Build it with that card, not before.

**What it is.** A reference to the exact cards a preceding sub-effect just moved —
KeyForge's "a card discarded **this way**" — where the player **chooses one** of
those cards. The producer half already exists: `DiscardTop` /
`DiscardTopOfEachDeck` record what they discard into `ctx.Produced.Discarded` (a
`[]LocalID` result set) in `internal/engine/effect_deck.go`. Two consumers of that
set already exist — `ForEachDiscarded` (act on **all** of them) and
`DiscardDeckUntil` + `PutDiscardedIntoHand` (find **the first match** by digging).
The missing consumer is **"the player picks one from the produced set."**

Why the result set is its own thing and not "the discard pile": it can be a strict
subset of the pile (discard 3 onto an already-full pile → only those 3 are
eligible), and it can be empty even when the pile is not.

**Cards (both unimplemented, in unbuilt sets):**

- **Junk Restoration** — "Discard the top 3 cards of your deck. You may put a card
  discarded this way into your hand." Discard 3, then optionally choose one of
  those 3 to keep.
- **Haunting Measures** — "Discard the top 6 cards of your deck. You may put a
  non-Geistoid card discarded this way into your hand." Same, with a house
  exclusion (non-Geistoid) and 6 cards.

**Design decided:**

- Build a **reusable "choose one from the produced result set" selection** — it
  picks one card from `ctx.Produced.Discarded` into `ctx.It`, then reuses the
  existing tail (`PutDiscardedIntoHand`, or a future purge-/archive-from-discard).
  Prefer this over a single fused `PutOneDiscardedIntoHand` node so hand / purge /
  archive all share one reference. This matches the "reference primitive" framing
  and covers the purge/archive-top cards below for free.
- The selection is **optional** ("you may") and carries a **filter**: a house
  include/exclude (Haunting Measures excludes Geistoid; Junk Restoration has no
  filter) plus a card-type refinement (for the purge/archive-top cards). One small
  filter on the chooser covers all three shapes.
- The **"purge/archive the top card"** scrap-adjacent cards are the same family
  with the tail swapped (purge or archive instead of put-into-hand) and no choice
  when only one card was produced — the `Amount: 1` case. Build them on the same
  reference when their set arrives.

This is deliberately distinct from the top-of-deck **zone slice**
(`LookAtTopOfDeck`, already built): a zone slice reads deck **positions**; a result
set reads **what an effect produced**.

## Live "considered an artifact/creature" type grant (ADR 0033 live route)

**Trigger set:** whichever set first stands up a card that converts a type on a
lifetime **other than the converted card's own** — Deanimator is the model. Build
it with that card, not before.

**What it is.** The **live route** of ADR 0033's two type-conversion routes. The
stored route (`LastingType` on `CardCore`) is built and covers permanent
self-conversions (Auto-Legionary; Effigy of Melerukh and The Mysticeti are its
still-unimplemented siblings, needing no new primitive). The **live** route is not
built: a card is a given type only while some other condition holds — a counter on
it, a source in play, a while-condition — and the change must **revert for free**
when that condition ends.

**Card (unimplemented, in an unbuilt set):**

- **Deanimator** — "Each card that has a mineralize counter on it is considered an
  artifact." The grant ends when Deanimator leaves play or the counter is removed,
  not when the converted card leaves — so it cannot be stored.

**Design decided (per ADR 0033):**

- Author it as a `ConstantAbility` composed over a `CounterInPlay{Kind, Target:
This}` read (per ADR 0024), computed **live at read time** the way `Power`/`Armor`
  fold `constantBonus` — **never** a `LastingType` write. It reverts for free: the
  next `TypeOf` read simply no longer sees the grant.
- Grow `TypeOf` one step **after** the `HostPlus`/`LastingType` checks: ask whether
  any active constant ability converts this card's type, exactly as the stat reads
  scan for bonuses. Do not store the result. Two Deanimators compose the same way,
  each grant recomputed every read.
- Only if a **stored** and a **live** override could ever land on one card at once
  does the single `LastingType` field become insufficient — then fold it into a
  flat comparable `TypeEntry{Card, Type, Source}` LIFO side-table (like the control
  stack, ADR 0028) and record the last-applied-wins precedence. No such card exists
  today; do not build the stack pre-emptively.

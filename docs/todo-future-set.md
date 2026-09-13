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

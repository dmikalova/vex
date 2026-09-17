# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Work that has **no consuming card in an implemented set yet** does not belong
  here — park it in [todo-future-set.md](todo-future-set.md), keyed to the set that
  first needs it.
- Cite the ADR or doc that decided an item where one exists.
- **When you need a decision from the human, ask in the reply itself, grill-me
  style** — a numbered list of `❓ **Q1** - **title**: <question>` with a `➡️`
  recommended answer under each — at the end of the turn, then stop. Do **not**
  reach for an interactive question tool: under autopilot it is auto-answered with
  "work autonomously" and the human never sees it. Plain-text questions at the end
  of the turn are the channel the human actually reads.

## Rule of Six — destruction replacement is not yet an orderable window entry

**Decided (2026-09-16, human) — the per-card-name usage ledger is BUILT.** The
Rule of Six is now one cap of 6 usages per card **name** per player per turn,
tracked in `GameState.UsagesThisTurn [maxCards]uint8` and summed by name+owner
(`game_ruleofsix.go`). Every usage records against it — play, discard, use
(reap/fight/`Action:`), each repeat loop past the first, each `Destroyed:`
resolution, each Replicator-style trigger — and reaching 6 bars every further
usage of that name (play, use, repeat, trigger, and `Destroyed:` resolution; the
creature still leaves play). This is how Reassembling Automaton at 0 power
terminates: its `Destroyed:` replacement stops standing in once the name's pool is
spent, so the tag is never lifted and it is destroyed for real.

**Remaining follow-up (architectural, not yet done).** The human's Flint's Stash
worked example — a creature carrying its own destruction _replacement_ **and** a
granted `Destroyed:` (gain 2 Æmber), where each destruction opens a window of two
_orderable_ `Destroyed:` abilities and choosing the replacement first forfeits the
other for that window — is **not** achievable yet. The destruction replacement is
currently a pre-window static `Replace` applied in
[game_leaves_play.go](../internal/engine/game_leaves_play.go) `filterUndestroyed`,
_before_ `destroyedAbilities` enrolls the window, so a replaced creature never
enters the orderable Destroyed window at all. Making the replacement an orderable
entry **inside** the Destroyed window (so it interleaves with other `Destroyed:`
abilities and the player orders them) is a separate, larger reshaping of the
leave-play path. Grill the desired ordering semantics before building it.

## Mass Mutation — deferred mechanics (grill before building)

The `massmutation` set is stubbed and its easy cards are being implemented. Two
mechanics are **deferred for a grilling session** because each needs a design
decision (a new node, a deckgen-time hook, or a cluster shape) before any of its
cards can be authored. Each group's cards are `//go:build todo` stubs in
`internal/cards/sets/massmutation/`. **Grill each group, decide the shape, then
implement the whole cluster together** (implement-cards: shape for the cluster,
not the first card).

### 1. Mass Mutation cards still to implement (normal cards, grill one by one)

The small-primitive bottleneck is cleared: every remaining non-gigantic Mass
Mutation card is a **normal card** that needs no new effect node, so they are
tackled one at a time or in small themed groups, grilling each for its wording and
the smallest engine seam it needs. Run `mage tool:missing -set=massmutation` for
the live list; `mage tool:nextCard -set=massmutation` names the next stub.

### 2. Remaining cards

051 Turnkey
100 Animator
102 Cyber-Clone
113 The Archivist
131 Commandeer
143 Angry Mob
153 Purify
159 Book of Malefaction
165 Lady Loreena
210 Saurian Egg
222 High Priest Torvus
224 Legion's March
257 Shoulder Id
292 Shadowsaurus
303 Blast Shielding
337 Commander Dhrxgar
347 J.O.N. Cargo
394 Aemberlution

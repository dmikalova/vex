# 36. Clusters place card families by strategy; one-per-house is complete by construction

## Context

Four card groups want deck generation to place a *family* of related cards, and
today's connection mechanism (ADR-less, documented in `docs/deck-generation.md`
§5) covers only one of the four shapes:

- **Sins** — seven cards in one House; drawing *any* one should top the pod up to a
  random 3–7 distinct sins. No single card leads it.
- **Horsemen** — four cards; a lead member (Horseman of Pestilence) pulls the other
  three.
- **Shards** — one card per House; drawing *any* Shard should place that House's
  Shard in *every* House pod of the deck. This is **deck-wide**, not pod-local, and
  needs a member for every real House or it cannot complete.
- Existing pullers — Timetraveller (`Pull(HelpFromFutureSelf, 1)`), Troop Call
  (`Pull(NiffleApe, 2)` + `PullSometimes(NiffleQueen, 0.15)`) — a fixed count from a
  named partner list.

Today's `Connection`/`Connects`/`Pull`/`PullSometimes` express only the last shape.
It is **pod-local** (`Generate` runs `expandConnections` per pod), the trigger is
implicit ("this slot's card carries a `Connection`") so there is no *any-member*
trigger, and the count is a fixed number or a single chance — no whole-pool, no
random count, no one-per-House. Adding each of the four as its own special case
would fork the fill path four ways.

Separately, none of this reaches the second axis these cards also need — a card
whose *identity* (not its companions) depends on the deck's other Houses. That is
the Template/Materialize seam (ADR 0004), extended here only by adding the deck's
Houses to `SlotContext`.

## Decision

**Generalize connections into clusters.** A **cluster** is a named family of member
cards. A card registers as a member; the cluster carries a **strategy** (how it
fills out) and a **trigger mode** (what fires it). Both are small strategies that,
per the engine's Strategy pattern, own their own behavior — the fill path branches
on neither a card name nor a family.

Strategies:

- **`WholePool`** — place every member (Horsemen).
- **`RandomCount(min, max)`** — place a random number of distinct members in the
  range (sins, 3–7).
- **`SelfPull(min, mean)`** — place a random number of *copies of the single
  triggering member itself*, not distinct members (Plague Rat pulls more Plague
  Rats). The count is `min + Poisson(mean − min)`, capped at `PodSize`: at least
  `min`, averaging about `mean`, with a thin tail that reaches a whole pod only very
  rarely. Its cluster has exactly one member, and `validateClusters` enforces
  `1 ≤ min ≤ mean ≤ PodSize`.
- **`PullExact`** — place one copy of each non-lead member per lead instance in the
  pod: two Timetravellers pull two Help from Future Self, never one. Always
  `ByLead`; the lead is the puller, the other members its exact partners. This is
  the migration target for the old exact-1 connections (Timetraveller, Hyde,
  Igon the Green), and needs no per-partner count — the count is the lead count.
- **`Pull`** — place a **per-partner** random count of each non-lead member when the
  lead rolls in: `min + Poisson(mean − min)` copies of that partner, capped at
  `PodSize`. Always `ByLead`. Unlike the family-wide strategies each pulled partner
  carries **its own** rate (set with `card.Pulled(cluster, min, mean)`), so one lead
  can pull two partners at different rates — Troop Call pulls a couple of Niffle Apes
  (`min 2, mean 3`) and, much less often, a Niffle Queen (`min 0, mean 0.85`). This
  is the migration target for the old fixed-count `Pull`/`PullSometimes` connections
  (Grumpus Tamer, Ortannu, Bear Flute, Faygin, Troop Call, Chain Gang, the nine ship
  blasters).
- **`OnePerHouse`** — place one member in each of the deck's Houses. This strategy
  alone is **deck-wide** and carries a **complete-by-construction gate**: every real
  House must have a member or the build fails, exactly the discipline of the
  rulebook term registry (ADR 0018). A missing House is satisfied explicitly by a
  `//go:build todo` stub. Shards use it.

Trigger modes:

- **`ByLead`** — a designated member fires the cluster (Horseman of Pestilence).
- **`ByAnyMember`** — any member fires it (any sin, any Shard). This is what the old
  connection model could not express.

`OneCopyPerDeck` on a member bounds duplicates independently of the strategy (each
sin, each Shard, max one per deck).

**A cluster must be able to fire.** `validateClusters` gates that a cluster has a
triggering member that actually rolls in the pool: for `ByLead` the lead must not be
`Rarity.Connected`, and for `ByAnyMember` at least one member must not be. A cluster
of only `Connected` cards — placed solely by pulling, never drawn — could never
trigger itself, so it fails the build. This is what makes the three `Connected`
Horsemen provably reachable: their `ByLead` lead (Pestilence) rolls normally and
rides them in, so the Connected-card-is-pulled invariant holds for cluster members.

**The pod-local fixpoint generalizes to deck scope.** `OnePerHouse` needs all three
pods filled first, then a deck-wide expansion pass places one member per House,
retrying another slot in a House on restriction and falling back to the existing
reseed-on-deadlock (a derived sub-seed) when a House cannot be satisfied. Pod-local
strategies keep running per pod as before.

**`SlotContext` gains `DeckHouses [3]House`** (the three resolved pod Houses) so the
Template/Materialize seam (ADR 0004) can bind a **partner house** — an Ambassador
(Sanctum) or Plant (Shadows) bound to one of the deck's other two Houses. This is
the second axis and is orthogonal to clusters: clusters decide *which other cards*
are placed; templates decide *what one drawn card becomes*.

## Consequences

- One mechanism replaces four special cases. A new family is a cluster registration
  plus a strategy, not a branch in `fillSlot`/`expandConnections`.
- `Connection`/`Connects`/`Pull`/`PullSometimes` become the `Pull` strategy of a
  cluster; existing pullers (Timetraveller, Troop Call, Horsemen) migrate onto the
  cluster API. `docs/deck-generation.md` §5 is rewritten from "connections" to
  "clusters" as part of the migration.
- Only `OnePerHouse` is deck-wide and gated; the other strategies stay pod-local and
  ungated, so most clusters cost nothing new. The deck-wide pass is contained to the
  fixpoint's scope, not the card API.
- `DeckHouses` in `SlotContext` is the only engine-facing change the template users
  (ambassadors, plants, banes) need; the concrete templates are the first real users
  of ADR 0004's seam.
- The boosted combinatorial Bane (choose 3 Houses, C(9,3)=84 variants, `1/84·Rare`
  weighting, generative name-splicing) is deliberately **out of scope** here and
  deferred to its own ADR; the plain per-House Bane proves the trait-table mechanism
  first.

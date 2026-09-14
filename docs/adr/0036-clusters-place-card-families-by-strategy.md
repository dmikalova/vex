# 36. Clusters place card families by strategy; one-per-house is complete by construction

## Context

Four card groups want deck generation to place a _family_ of related cards, and
today's connection mechanism (ADR-less, documented in `docs/deck-generation.md`
§5) covers only one of the four shapes:

- **Sins** — seven cards in one House; drawing _any_ one should top the pod up to a
  random 3–7 distinct sins. No single card leads it.
- **Horsemen** — four cards; a lead member (Horseman of Pestilence) pulls the other
  three.
- **Shards** — one card per House; drawing _any_ Shard should place that House's
  Shard in _every_ House pod of the deck. This is **deck-wide**, not pod-local, and
  needs a member for every real House or it cannot complete.
- Existing pullers — Timetraveller (`Pull(HelpFromFutureSelf, 1)`), Troop Call
  (`Pull(NiffleApe, 2)` + `PullSometimes(NiffleQueen, 0.15)`) — a fixed count from a
  named partner list.

Today's `Connection`/`Connects`/`Pull`/`PullSometimes` express only the last shape.
It is **pod-local** (`Generate` runs `expandConnections` per pod), the trigger is
implicit ("this slot's card carries a `Connection`") so there is no _any-member_
trigger, and the count is a fixed number or a single chance — no whole-pool, no
random count, no one-per-House. Adding each of the four as its own special case
would fork the fill path four ways.

Separately, none of this reaches the second axis these cards also need — a card
whose _identity_ (not its companions) depends on the deck's other Houses. That is
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
- **`SelfPull(min, mean)`** — place a random number of _copies of the single
  triggering member itself_, not distinct members (Plague Rat pulls more Plague
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
the second axis and is orthogonal to clusters: clusters decide _which other cards_
are placed; templates decide _what one drawn card becomes_.

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

## Cross-set clusters and errant pods (extension)

`OnePerHouse` was first built per-set: `NewSet` indexed a set's own cards, and its
gate required a member for each of that set's Houses. Two forces broke that scope:

- **The Shards are a cross-set family.** Age of Ascension prints seven (one per its
  Houses); the last two Houses (Saurian, Star Alliance) never got a Shard. A per-set
  index can never complete the nine-House cycle, and a Shard drawn as a legacy card
  into a set that prints none would not fire the cycle at all.
- **The errant pod** (`docs/deck-generation.md` §1) lets a pod take a **foreign
  House** — one not native to the deck's set — drawn wholesale from the cross-set
  legacy pool. A deck can therefore reach a House its set never printed a Shard for,
  which the per-set gate cannot see.

Decisions:

- **The Shard cluster declaration is hoisted to a shared package**
  (`internal/cards/clusters`), imported by every set that prints a Shard, so its
  strategy and trigger are declared once and cannot drift between members printed in
  different sets. Single-set clusters stay in their own set package.
- **Deck-wide clusters resolve from a catalog-wide `ClusterPool`.** Because cluster
  identity is the `Name` string and cards self-register globally, `NewClusterPool`
  builds one index over the whole registered catalog, with **no set→set import**.
  It is attached to each base set with `WithClusters`; `expandClusters` resolves
  `OnePerHouse` from it when present, so a legacy-drawn or errant-House Shard
  completes the cycle across every House the deck reaches.
- **The gate is derived from _deckable_ Houses.** `validateCrossClusters` requires
  an `OnePerHouse` member for every House a set can deck — its native Houses **plus**
  its errant Houses (Houses present in the legacy pool but not native,
  `Set.errantHouses`, computed in `WithLegacy`). Still complete-by-construction
  (ADR 0018), now widened to the Houses an errant pod introduces, not hardcoded to
  nine.
- **The two missing Shards are real `Connected` cards in a reservoir set.** Shard of
  Glory (Saurian) and Shard of Unity (Star Alliance) are Vactrol-invented cards in
  `internal/cards/sets/anomalyexpansion`, declaring their home set with
  `card.InSet(card.AE)` (no provenance — set membership is decoupled from provenance).
  Being `Connected`, they stay out of the draw pool and the legacy pool and enter a
  deck **only** through the Shard cluster. Anomaly Expansion is a **reservoir set**:
  its cards register and feed the catalog cluster pool, but it has no draw pool of
  its own, so `buildDeckSets` deals no deck from it (a set with no `Draftable` card
  is skipped). Reservoir does not mean "out of every draw pool", though: the set
  also holds the housed, non-`Connected` anomaly cards, and those remain eligible
  for the legacy pool, so other sets' legacy and legacy-maverick slots can draw an
  anomaly even though the Anomaly Expansion is never itself draftable. Their
  abilities are provisional Vactrol inventions (see
  `docs/keyforge-divergences.md`).
- **The Ambassador and Plant templates already cover errant Houses for free** — they
  materialize per partner House from `DeckHouses`, so an errant House gets its
  Ambassador / Plant with no new card. Only the Shard cycle needed the two new cards.

## Filtered pulls: a lead card guarantees a floor of predicate-matching cards (extension)

Clusters place _named members_. A second, orthogonal shape appeared with Chief
Engineer Walls: a card whose deck-building payoff is a whole open-ended _category_
of cards — "any Upgrade or Robot card" — that no fixed member list can name. Walls
retrieves Upgrades and Robots from the discard pile, so a deck that runs Walls
wants a couple guaranteed to retrieve, drawn from **any** House, not a specific
partner.

Decisions:

- **A filtered pull is a `FilteredCluster`, kept in a separate field from
  `ClusterMembership`.** It carries a `Name`, a `Lead` (stamped from the declaring
  card at `NewSet`), a `Floor`, and a `Match func(CardDefinition) bool` predicate.
  It names no members — the pulled cards are whatever the pool offers that `Match` —
  so it cannot be a `ClusterStrategy` (those enumerate members). Because a card
  carries only one `ClusterMembership`, the filtered pull lives in its own
  `GenerationProfile.Leads *FilteredCluster` field, so **a card can lead a filtered
  pull and still belong to a named cluster** (Walls both partners its blaster's
  `Pull` cluster and leads the Upgrade/Robot filtered pull). Authored with
  `card.PullsMatching(name, floor, match)`.
- **The floor counts cards already in the deck.** `expandFilteredClusters` runs
  after `expandClusters` (so it never disturbs the OnePerHouse cycle), counts the
  slots already matching the predicate, and tops up only the shortfall from the
  pool — a deck that rolls enough Upgrades on its own pulls nothing, and the pull is
  a floor, not a fixed add.
- **Top-up is any-House, rehousing off-House pulls as mavericks.** Candidates are
  the set's draftable pool cards that `Match`, minus those already in the deck and
  minus placed one-copy-per-deck cards, shuffled by the deck seed. Each lands in a
  pod of its own House when the deck has one (a natural card) else any pod as a
  maverick, overwriting a slot that is not the lead, not a cluster/OnePerHouse
  member (`protectedNames`), and does not already match — so the pull never reduces
  the match count, displaces a Shard, or evicts the lead.
- **Complete-by-construction, like every other cluster.** `validateFilteredClusters`
  fails the build on a nil predicate, a `Floor < 1`, or a pool that cannot satisfy
  the `Floor` — the same gate discipline as `validateClusters` (ADR 0018), so a deck
  can never silently come up short.
- **Zero blast radius when unused.** A set with no filtered pull makes no extra RNG
  draws (`filteredNames` is empty), so existing decks are byte-for-byte unchanged;
  only a deck that actually contains the lead and is short of the floor pulls.

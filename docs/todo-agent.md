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

## Mass Mutation — deferred mechanics (grill before building)

The `massmutation` set is stubbed and its easy cards are being implemented. Four
mechanics are **deferred for a grilling session** because each needs a design
decision (a new node, a deckgen-time hook, or a cluster shape) before any of its
cards can be authored. Each group's cards are `//go:build todo` stubs in
`internal/cards/sets/massmutation/`. **Grill each group, decide the shape, then
implement the whole cluster together** (implement-cards: shape for the cluster,
not the first card).

### 1. Enhancements — subsystem built; unstub the remaining cards

The bonus-icon + Enhance subsystem is **implemented** ([ADR 0041](adr/0041-bonus-icons-are-the-primitive.md)):
`card.WithBonus(card.Bonus.X…)` (printed icons, resolve on play, ordered),
`card.WithEnhance(card.Bonus.X…)` (the deckgen Enhance source), the deck-wide
finishing pass in `internal/deckgen/enhance.go`, and `card.WithoutEnhancement(kinds…)`
(the opt-out, on Effervescent Principle). Æmber, Capture, Damage, and Draw icons
resolve on play, and bonus icons render on the deck list to the right of each
card's name. Unstubbed so far: **Splinter**, **Mutant Cutpurse**, **Infomorph**,
**Gloriana's Attendant**, **General Xalvador** (pure-Enhance), plus **Dark Minion**
(Destroyed), **Hystricog** (Action: destroy a damaged creature), **Miasma Bomb**
(Action: skip forge), and **Sagittarii's Gaze** (Play: exalt).

**What remains:** unstub each Enhance card's own ability + add `WithEnhance(…)`.
The deferred manipulator cards need the bonus-icon resolution step's
interception, built next against the first-class step (ADR 0041): Wild Bounty
(resolve each icon again), Amphora Captura (resolve Æmber as Capture/Damage/Draw),
Scrivener Favian (resolve a Capture as a Steal), Master of the Grey (prevent the
opponent's icons), Ensign El-Samra (reveal top card, resolve its icons as if
played), and Adaptoid, Dark Queen Gloriana, and Mutagenesis Researcher (react to
playing a card with a bonus icon).
**Still stubbed (23):** Adaptoid, Amphora Captura, Armory Officer Nel, Boss
Zarek, Bring Low, Burning Glare, Chronus, Consul Primus, Crewman Jorg, Dark
Centurion, Dark Queen Gloriana, Ensign El-Samra, Fission Bloom,
Maleficorn, Mutagenesis Researcher, Opposition Research, Resurgence,
Scrivener Favian, Survey, Tempting Offer, Wail of the
Damned, Waking Nightmare, Wild Bounty. Icons beyond the four (Discard, House,
+1-power) and the distinct graphical element for landed-vs-printed icons are not
modelled yet.

### 2. Gigantics (two cards, one big creature) + their tutors

A **gigantic** creature is a single large creature split across **two cards** (a
"base" half carrying the text/stats and an "art" half); both must be in play
together to form the creature. This needs new engine state (a creature spanning two
card slots) and deckgen guarantees (both halves present). The stubs are
comment-only (`//go:build todo`, no `card.New`) because the stub generator has no
`gigantic` card type. **Decisions needed:** how a two-card creature is modelled in
flat state, how it is played, how it occupies the battleline, and how deckgen
pairs the halves. **Gigantic halves (14):** Ultra Gravitron, Deusillus, Niffle
Kong (MM); Tormax, Wretched Anathema, Horizon Saber, Ascendant Hester, Sirs
Colossus, Bawretchadontius, Boosted B4-RRY, Dodger's 10, Cadet Allison, J43G3R V,
Titanic Bumblebird (MoMu). **Tutor cards that fetch gigantic halves (4):** It's
Coming… (MM #117), Build Your Champion, Digging Up the Monster, Tomes Gigantica
(MoMu #002-004) — build these _after_ gigantics exist, since they search for
"halves of a gigantic creature."

### 3. Mutant cycle (each house pair contributes two properties)

The **mutant cycle** (Lyco-Thief, Daemo-Knight, the seven sins, and every
`<prefix>-<house>` creature) builds a creature from a combination of properties:
a mutation prefix plus a house suffix each contribute stats/traits/abilities.
These carry the source rarity **Variant** (48 cards) — which this repo has no
engine value for and **must be translated to a real rarity by the maintainer**
(repo memory: never add an `engine.Variant` enum; ask which rarity). Gather the
whole cycle as one cluster. **Decisions needed:** the combination model (how two
property sources merge into one card), the Variant→rarity mapping, and the cluster
shape. **Cycle cards (42 Variant):** the Variant-rarity
`Dino-/Daemo-/Lyco-/Sacro-/Techno-/Umbra-/Xeno-` creatures across all houses,
plus Desire, Envy, Gluttony, Greed, Pride, Sloth, Wrath. (The seven **Common**
cycle leads — Techno-Fiend, Daemo-Bot, Daemo-Knight, Daemo-Saurus, Daemo-Thief,
Daemo-Alien, Daemo-Beast — are ordinary implementable cards, not part of this
deferred group.)

### 4. Dark Æmber Vault + cluster specials

**Dark Æmber Vault** (MM #001, Sanctum artifact, location, Special) — "After you
play a Mutant creature, draw a card. Each friendly Mutant creature gets +2 power."
— is a **cluster card** needing special deckgen handling (it wants a Mutant-heavy
deck around it). The other Special-slot cluster cards to decide alongside it: the
`Monument to …` artifacts. **Decisions needed:** the cluster shape and pull rates
(prompt the human — never guess pull rates, per `internal/cards/AGENTS.md`).

#### 4c. Z-Force Agent 14 cluster — DECIDED, blocked on card implementation

**Decision (from the human):** Z-Force Agent 14 (MM #353, Star Alliance Creature,
Rare — a rollable lead) leads a cluster that pulls in all three `Z-` upgrades —
**Z-Particle Tracker** (#354), **Z-Ray Blaster** (#355), **Z-Wave Emitter** (#356)
— **at exactly 1:1: whenever Z-Force Agent 14 is placed, all three come, one of
each.** That is a `WholePool` / `ByLead` cluster (the whole member pool is placed
with the lead, like `horsemenCluster`). The three upgrades become
`card.Rarity.Connected` — they never roll on their own and enter a deck only when
Z-Force Agent 14 is placed.

Wiring, once the cards are real: declare one shared `card.Cluster{Name: "Z-Force
Agent 14", Strategy: card.ClusterStrategy.WholePool, Trigger:
card.ClusterTrigger.ByLead}` beside the lead; give Z-Force Agent 14
`card.LeadsCluster(...)` and each Z- upgrade `card.InCluster(...)` plus
`card.Rarity.Connected` (dropping their current `Rarity.Special`).

**Blocked on** the three Z- upgrades, which are currently `//go:build todo` stubs
— a `ByLead` cluster panics `validateClusters` until its members are registered.
Each needs an effect the engine may not have yet: Z-Particle Tracker searches the
deck for an upgrade to hand and shuffles; Z-Ray Blaster grants +3 power and a
"before fight, deal 3 damage to each neighbor of the creature it fights" ability;
Z-Wave Emitter wards the host at the start of your turn. Implement the three
upgrades, un-stub them, then wire the cluster as above.

#### 4a. Dark Harbinger cluster — DECIDED, blocked on card implementation

**Decision (from the human):** Dark Harbinger (MM #381, Untamed Creature,
Uncommon — a rollable lead) leads a `ByLead` cluster whose members are the three
Untamed Tactics **Mutation of Cunning** (#413), **Mutation of Fury** (#414), and
**Mutation of Instinct** (#415). The three mutations are **`card.Rarity.Connected`**
— they never roll on their own and enter a deck only when Dark Harbinger is placed.
When Dark Harbinger is in a pod, pull a random subset of the three: at least 1,
up to all 3, averaging 2, one copy of each. That is a `RandomCount` / `ByLead`
cluster with Min 1, Max 3 (uniform over {1,2,3} → mean 2).

Wiring, once the cards are real: declare one shared `card.Cluster{Name: "Dark
Harbinger", Strategy: card.ClusterStrategy.RandomCount, Trigger:
card.ClusterTrigger.ByLead, Min: 1, Max: 3}` beside the lead; give Dark Harbinger
`card.LeadsCluster(...)` and each mutation `card.InCluster(...)` plus
`card.Rarity.Connected`.

**Blocked on** the three mutations, which are currently `//go:build todo` stubs —
a `ByLead` cluster panics `validateClusters` until its members are registered, so
Dark Harbinger cannot carry `LeadsCluster` before they exist. Their abilities need
two lasting mechanics the engine does not have yet: (1) **grant a single trait
(Mutant) until the start of your next turn** (only whole-text-box trait copying
exists today — `HasTrait` folds `grantedTextBoxSources`, so add a granted-trait
store mirroring `KeywordsUntilNextTurn`, cleared at the ready phase); and (2)
**assault N until the start of your next turn** (`GainAssault`/`TempAssaultBonus`
are remainder-of-turn only — Fury wants a next-turn duration). Both target **a
single chosen creature**, and the Play text combines the two grants: "a creature
gains elusive/skirmish/assault 3 and the Mutant trait" — render it under one shared
"… gains … and the Mutant trait until the start of your next turn" clause the way
`GainAssault.durationSubject`/`durationPredicate` share theirs.

**Deckgen support also needed:** `randomMembers` currently shuffles all of
`ci.members`, which for a `ByLead` `RandomCount` would let the lead itself be
picked as a placed copy; exclude `ci.lead` from the random pick and validate `Max`
against the **non-lead** member count. Add a `podcluster_test.go` case for a
`RandomCount` `ByLead` cluster (lead planted, subset of the Connected partners
placed, never the lead again).

### 5. Reprinted cluster members whose lead is not in Mass Mutation

Mass Mutation's catalog reprints two cards that are **pulled cluster members**
implemented in Worlds Collide — **Commander Chan** and **Sensor Chief Garcia** —
whose cluster _leads_ (Chan's Blaster, Garcia's Blaster, both Worlds Collide
signature upgrades) are **not** in Mass Mutation. A pulled member without its lead
would leave deck generation with a lead-less `ByLead` cluster, which panics by
design (ADR 0036). As a stopgap, the reprint generator (`reprintsForSet` in
`magefiles/cardlookup/stub.go`) now **skips** an orphaned pulled-member reprint, so
those two cards are currently **absent from Mass Mutation's deck-generation pool**.
**Decision needed:** the right long-term behavior — (a) let a reprinted cluster
member whose lead is absent join the set as a plain (non-cluster) pool card, which
means growing `deckgen` to drop a lead-less `ByLead` cluster for a set _only when
the lead exists elsewhere in the catalog_ (so the "forgot `LeadsCluster`" gate
still fires); or (b) keep excluding them. Option (a) restores Commander Chan and
Sensor Chief Garcia to Mass Mutation as real cards and is the KeyForge-accurate
outcome (they print in Mass Mutation without a blaster). Grill, then implement.

### 6. Small primitives Mass Mutation cards still need (the current bottleneck)

The easy cards are largely done; most remaining Mass Mutation cards are blocked on
a small, clean engine extension (a target filter, a count, a trigger, a condition,
or a field — **not** a new effect node). This is now the throughput lever: build a
cluster of these primitives, then a card wave clears the cards they unblock. Each
is listed with the card(s) it releases.

- **Pool-threshold Æmber-theft protection** ("while you have N Æmber, your Æmber
  cannot be stolen") → Cephaloist (#362).
- **Combined hand-or-archives discard source.** `DiscardCard.Zone` is a single
  zone → Munchling (#076), Novu Dynamo (#093).
- **Off-flank key-cost modifier** (a key-cost analogue of `WithDrawModifierOffFlank`)
  → Titan Engineer (#081).
- **Trait axis on `MayPlayOrUse`** (grant scoped to a trait, not a house/type) →
  Mutagenic Serum (#091); pairs with a trait-scoped use grant.
- **"the creature it fights is the most powerful enemy" combat condition** →
  Baldric the Bold (#144).
- **Count of friendly creatures with Æmber on them** (an `InPlay.WithAember`
  filter; distinct from `AemberOnFriendlyCreatures`, which counts pips) → Faust the
  Great (#192).
- **"Lower power than this creature" refinement** → Dreadbone Decimus (#204). A
  `PowerOfThis` count feeding `PowerLessThan` renders awkwardly ("power less than
  Dreadbone Decimus"); it wants a dedicated refinement rendering "with lower power
  than `<self>`", not just a count.
- **"player who controls the most powerful creature" selector** for a gain's
  `Player` → Forum of Giants (#219).

- **Fractional / Count amount on `MoveAember`** ("move half, rounded up") →
  Patronage (#227).
- **Private peek at the opponent's deck top** (`LookAtTopOfDeck` is hardcoded to
  your deck) → Vandalize (#260).
- **Artifact self-destroy-when-no-creatures** rule (the poison grant already works)
  → Doom Sigil (#277).
- **A "scheme" generic counter kind** (or reuse an existing kind) → Mastermindy
  (#285).
- **Per-house partition `TopAct`** (archive each of a chosen house, discard the
  rest) → New Frontiers (#326).
- **Trait-scoped friendly-vs-enemy count comparison** → Pismire (#372).
- **Pool-gated + house-filtered enters-ready** (`WithFriendlyEntersPlayReady` takes
  only a `CardType`) → Fandangle (#365).
- **Resolve-the-bonus-icons-of-a-purged-artifact** effect → Reclaimed by Nature
  (#374).
- **Constant, type-scoped play permission** ("may play upgrades as if in the active
  house") → Matter Maker (#349).
- **Each-player end-of-turn trigger** → Pincerator (#289).
- **Self-scoped keyword loss with a duration** ("loses elusive until your next
  turn") → Reckless Rizzo (#273).
- **Optional-discard `TopAct`** ("look at the top card, you may discard it") →
  Scout Pete (#311).
- **`Exhaust` as a gating effect that binds its target** ("exhaust a creature; if
  you do, …") → Humble (#208).
- **Bonus-icon resolution suppression** ("opponent cannot resolve bonus icons") →
  Master of the Grey (#169). Larger; may warrant its own design pass.

**Save the Pack is a reprint that Mass Mutation cannot yet claim.** It is a Call
of the Archons card that More Mutation reprints (MoMu #431), but it is not in the
`massmutation.json` catalog, so `mage tool:stub massmutation` cannot claim it as a
reprint in `0set.go`. Its hand-written stub was deleted, so Mass Mutation's pool
currently lacks it. This is part of the **More Mutation merge**: the MoMu-only
cards Mass Mutation should include split into genuinely-new cards (implemented as
stubs already) and _reprints_ like Save the Pack (already implemented elsewhere,
needing a reprint claim). Settle how MoMu-only cards fold into Mass Mutation's
catalog/pool, then Save the Pack joins as a reprint. (Orb of Wonder, formerly in
this bucket, is no longer a reprint at all: it previewed Mass Mutation as a Worlds
Collide anomaly, so now that Mass Mutation is implemented it was moved out of the
Anomaly Expansion into Mass Mutation as a real Sanctum Rare artifact — see the
anomaly two-phase rule in `internal/cards/AGENTS.md`.)

### 7. Keyfrog forges _and purges itself_ (divergence to confirm)

Keyfrog (MM #369) is authored with `card.ForgeKey{}`, whose shared node always
renders "forge a key at current cost -> purge {self}" (every forge-from-ability
card in Vactrol self-purges). The printed Keyfrog says only "Destroyed: Forge a
key at current cost." — it discards normally, it does not purge. The card is
implemented and forges correctly; the extra self-purge is a rule divergence
(purge vs discard) inherited from the shared node. Confirm whether to accept it
(and record it in `docs/keyforge-divergences.md`) or add a forge-without-self-purge
option to `ForgeKey`.

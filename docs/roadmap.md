# Vactrol — Development Roadmap

Vactrol is a digital card game written in **Go**, designed as a clone of **KeyForge**
(complexity roughly between Magic and Pokémon). This document is the high-level path
distilled from [outline.md](outline.md). Each step is a phase of work with the critical
concepts and named techniques **bolded** so they can be researched further later.

---

## 1. Lock the game design & theme

Vactrol keeps KeyForge's core loop but swaps humanoid creatures for **abstract object
constructs** so that art can be procedurally generated at scale. Each **House** is a
distinct visual/mechanical architecture — **Alchemia** (volatile liquid chemists / vials),
**Mechanica** (clockwork gears with adjacency synergy), and **Geometria** (crystalline
prisms with line-of-sight effects). A single **unified vocabulary** is shared across all
houses so interactions stay instantly readable: **Flux** (the resource, replacing Æmber),
**Distill** (exhaust for 1 Flux, replacing Reap), **Interfere** (clash two ready cards,
each dealing Damage equal to Power, replacing Fight), **Damage** (structural degradation),
and **Stabilize an Anchor** (the win condition, replacing "forge a key" — stabilize 3).
Card types map one-for-one: creatures → units/potions, upgrades → **Reagents/Catalysts**,
artifacts → **Lab Equipment**, actions → one-shots.

## 2. Design the card DSL (idiomatic Go + AST)

Cards are authored in **idiomatic Go** rather than raw JSON, using **struct literals** for
straightforward definitions and the **Functional Options pattern** (`NewCard(name,
WithHouse(...), WithPower(...), WithTrigger(...))`) where defaults or optional configuration
are needed. This keeps definitions declarative, type-safe, and instantly readable to any Go
developer, and it sidesteps the error-handling and non-idiomatic pitfalls of pointer-receiver
method chaining (the fluent/builder style common in TypeScript). Under the hood each
definition builds an **Abstract Syntax Tree (AST)** of intents/conditions/targets that acts
as the single source of truth. Two interpreters walk the same tree: the **game-engine
interpreter** (mutates live state) and the **text/localization interpreter** (emits the
English card face). This guarantees **zero desync** between what a card says and what it
does, and makes translation a matter of adding a new text backend.

## 3. Build the core engine — deterministic & pointerless

The engine must be **deterministic** (seed in → identical game out) and built for speed.
Separate a **Card Blueprint** (parameterized generative template) from a **Card Instance**
(a concrete card in a match). Runtime state lives in a flat, **pointerless `GameState`
struct** using **fixed-capacity arrays** (e.g. `[2][10]CardInstance`) so it can be copied
by value in nanoseconds with **zero heap allocation / no GC pressure** — critical for later
simulation work. Dynamic/temporary effects (e.g. "gains a trigger until end of turn") are
handled by a **Modifier system** layered over the static definition plus an end-of-turn
**cleanup phase**, never by mutating the base card. A read-only **card catalog**
(`map[string]CardDefinition`) provides O(1) lookups.

## 4. Procedurally generate cards (Director + power budget)

Layer a **Curated Randomizer / Director pattern** on top of the card constructors so
generated cards stay coherent. The director rolls parameters within **safety rails** (a
keyword from a pool, a damage value from a range) and passes them through the same struct
literals / functional options, so variants are always valid and type-checked. Every modifier
contributes to a **Balance Score / power budget**, which
maps to an **Appearance Weight** so stronger rolls become exponentially rarer. A **print-run
generator** (CLI) can sample thousands of rolls and emit a manifest of each variant with its
`powerScore` and `spawnChance`.

## 5. Build the procedural art pipeline

Avoid AI art and per-card commissions by **compositing layered SVG assets** driven by the
card's own data. Art is a back-to-front **layer stack** (background/biome → action → target
→ frame), where each mechanic contributes a visual element (e.g. `dealDamage` stamps an
energy asset that scales with the amount). Cohesion comes from three rules: **abstract vector
shapes**, an **algorithmic anchor/layout grid**, and **dynamic palettization** (assets drawn
in grayscale, colors injected at runtime from card traits). A realistic budget path is a
**modular vector asset kit** (~60 grayscale components) rather than 50 illustrations —
roughly $3k–$7k, kept low by delivering a grey-box layout and grayscale specs to the artist.

## 6. Design deck generation (target-score + synergy)

Not all cards should be equally good — **chaff** makes **bombs** feel exciting. Assign each
blueprint an **Independent Power Score** and give each deck a **Target Deck Score**. Generation
uses a **"rubber band" algorithm**: as cards are drafted it tracks the running average and
dynamically reweights draw probabilities to pull the total back toward target (rolling bombs
forces in weak cards to "pay" for them). A **synergy matrix** discovered by simulation (see
step 10) applies **contextual multipliers** so combo pairs are treated as bombs and padded
accordingly. This distinguishes **Independent Power** (value in a vacuum — bad cards should
have *negative* Drawn Win Rate) from **Contextual Power** (value alongside synergies).

## 7. Define a set (release scope) & bootstrap content

A **release set** (Magic/Pokémon/KeyForge sense, not the math sense) is ~400 blueprints:
~50% new to the set, ~50% reprints from earlier sets. The total pool may reach ~2,000
blueprints across ~50 **action primitives**, yielding tens of thousands of possible cards —
which is a storage problem, not a compute problem (blueprints compile to only a few MB).
Bootstrap by **cloning an existing, already-balanced KeyForge set** to skip the cold-start
problem, then let simulation estimate baseline values from there. Get a **minimum playable
set of cards** working end-to-end before scaling.

## 8. Build the MCTS / ISMCTS engine

The AI uses **Monte Carlo Tree Search (MCTS)** — selection, expansion, simulation,
back-propagation — running on `GameState.FastCopy()` with zero allocations. Because a match
is a **closed world of only 72 cards** (two 36-card decks), optimize with **Local IDs**
(`uint8` per card), **Transposition Tables** (Zobrist/xxHash to prune duplicate states —
duplicates are common and collapse the tree), and **action masking via bitmasks** (`uint64`)
so legal-move lookup is a single bitwise op. Hidden information is handled with
**Information Set MCTS (ISMCTS)**, specifically **Single-Observer ISMCTS (SO-ISMCTS)**: at
the start of every rollout, **determinize** the opponent's unseen cards into a random valid
hand, then update **one shared tree** whose nodes are *information sets*. This avoids the
**strategy fusion / clairvoyance bias** of running separate perfect-information trees
(**PIMCTS**). Watch for **non-locality** — smoothed out by 2,000–10,000 rollouts per move.

## 9. Compile to WASM & build the client

Compile the same Go engine to **WebAssembly** (`GOOS=js GOARCH=wasm` or **TinyGo**) so it
runs in the browser. The architecture is a **hybrid server-authoritative "thick client"**:
the **authoritative server** holds full state, validates committed action batches through
its own copy of the engine, and broadcasts **state deltas scrubbed of hidden info** over
**WebSockets** (consider **Protocol Buffers** for compact payloads). The client runs the
engine locally for **instant legality checks**, **turn staging / sandboxing with undo**, and
**free client-side ISMCTS suggestions** in a **Web Worker** (zero server compute). Reject
**peer-to-peer** play: hidden-info P2P requires **Mental Poker** cryptography (verifiable
shuffles, threshold key exchange) that clashes with a fast simulation engine and breaks on
disconnects. A cheap VPS can referee thousands of connections since all heavy compute lives
on clients.

## 10. Stand up the deck-rating system

Turn simulation output into an automated rating system analogous to KeyForge's community
**AERC** (base traits like expected Flux, board control, effective power) and **SAS**
(**S**ynergy **A**and antisynergy **S**ystem). The engine measures each card's contribution
to the win rate to assign base trait values, then mines for **synergy pairs** (positive
multipliers) and **antisynergies** (negative modifiers) that feed back into deck generation
(step 6). The headline metric is **Drawn Win Rate (DWR)** — win rate when a card is drawn
minus when it isn't — which is what flags outlier cards.

## 11. Define AI playstyle heuristics as weight vectors

Playstyles are **interpretable linear weight vectors** (not black-box neural nets) applied
as **rollout evaluation biases** in MCTS. A bot's "genome" is a small struct of multipliers
over human-readable concepts (`FluxAdvantage`, `AnchorAdvantage`, `EnemyBoardClear`,
`SelfBoardPresence`, `SynergyTrigger`, …). Reading the weights reveals the personality (high
`FluxAdvantage` + low `EnemyBoardClear` = a greedy race deck). Aggro/Control/Combo emerge
from different weightings of the *same* metrics.

## 12. Evolve AI archetypes (Quality Diversity)

Tune those weights automatically with **evolutionary algorithms** — a proven, interpretable
approach demonstrated in games like Hearthstone via **competitive coevolution**. A basic
**Genetic Algorithm** loop (tournament fitness → selection → crossover → mutation) works;
**CMA-ES** (Covariance Matrix Adaptation Evolution Strategy) is the state-of-the-art
continuous optimizer if faster convergence is needed. Critically, **personality lives in the
weights, not the deck** — so to grow *distinct* archetypes simultaneously, use **Quality
Diversity**, specifically **MAP-Elites**: define a **behavior space** grid (e.g. game length
× unspent Flux) and keep the single best-performing **elite** per niche. One run then yields
the best Aggro, Control, and Combo bots at once, usable as fixed **sparring partners**.

## 13. Run the balance-testing pipeline at scale

Balance is measured by mass simulation. Rough math: ~140 decisions/game × ~2,000 rollouts ≈
280k rollouts/game ≈ ~3 CPU-seconds at ~100k rollouts/sec/core. A full 400-card set needs
~1–1.5M games (~$45–$75 in AWS **Graviton spot** compute; a 2,000-blueprint sweep is ~10M
games). Scale horizontally with **stateless Docker binaries** taking a JSON seed, an
**auto-scaling spot fleet / preemptible cluster** fed by a queue (SQS), sinking telemetry to
a **columnar store** (ClickHouse / Parquet-on-S3 + Athena / BigQuery). Cut cost without
touching the engine via **Multi-Armed Bandit early stopping** (halt obviously broken/fine
cards early), **mercy-rule early halts**, and a **Gauntlet** of fixed baseline decks instead
of full random-vs-random matrices.

## 14. Keep sets aligned & continuously rebalanced

Avoid an O(N²) set-vs-set explosion using a **Historical Gauntlet** and the **transitive
property**: measure every new set against the same frozen gauntlet; if Set 5 and Set 1 each
sit at ~50% vs the gauntlet, they're balanced against each other without ever playing
directly. Each release only needs ~1.5–2M sims (internal draft balance + global constructed
alignment). For ongoing tweaks, treat balance like **CI/CD**: compute a change's **blast
radius** (only re-sim decks that use the changed card), run **differential/regression tests**
with small sample gates that escalate on variance, apply **Bayesian updating** with the prior
run as the statistical prior, and **version the meta** (Gauntlet v4, etc.) so old sets don't
shift underfoot. Store only compact win/loss telemetry (~100–200 bytes/game) and reconstruct
full games on demand via **deterministic replay** from stored seeds.

## 15. Add undo & post-game replay tools

Because the engine is deterministic, undo/replay is **Event Sourcing**: record the initial
seed plus an **append-only event log**; recreate any state by replaying events, optimized
with per-turn **snapshots** so undo only replays a few recent events. The one hard part is
the **hidden-information trap** — an action that reveals a hidden card (RNG draw) grants
knowledge that can't be un-known. Solve it by separating **Staged Actions** (public-only,
freely undoable locally) from **Committed Actions** (trigger RNG / reveal state, permanently
lock undo history). **Post-game replay** ignores this entirely: load the seed + full event
log and scrub a timeline freely since the match is already resolved.

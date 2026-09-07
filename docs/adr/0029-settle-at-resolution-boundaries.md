# Settle destruction at resolution boundaries, not after every state change

## Context

The most common bug our property tests find is a creature that holds lethal
damage — `Damage >= Power`, or `Power <= 0` — yet is still in play. Power is
derived from the whole board: printed power plus upgrade bonuses, power counters,
temporary bonuses, constant-ability bonuses (some scaled `Per` a live count), and
the Æmber sitting on the creature. So a change _anywhere_ — a buff card leaving,
a counter moving, Æmber spent off a creature that draws power from it, a key
forged — can make a distant creature lethal, or revive one that was about to die.

`settleDestroyed` already answers this correctly: it is an idempotent fixpoint
sweep of both battlelines that repeats until no creature is newly lethal, guarded
by a `settling` re-entrancy flag so a sweep triggered mid-batch does not break the
simultaneous "tagged for destruction" semantics (§537). The defect is not the
sweep; it is that the sweep is hand-called at roughly twenty sites. Miss one new
site and a creature lingers with lethal damage.

Keyteki solves the same class of bug by re-validating the whole state after every
single change. We reject that, and the reason is **not** speed — a sweep is
`O(board)` (about twenty-four creatures), trivial next to the work MCTS already
does per node. The reason is **correctness**. Validating between the two mutations
of one logical batch breaks KeyForge's simultaneity: when several creatures leave
together (Epic Quest archives "Lion" Bautrem and the neighbour it was buffing at
once), a mid-batch check would destroy the neighbour for the power it just lost
before the batch finishes archiving it. The `settling` flag exists precisely to
hold that window open; per-change validation would defeat it.

## Decision

Settle at **resolution boundaries**, not after each mutation. A boundary is:

1. after each triggered ability's `Effect` resolves in `triggerAbilitiesAs`;
2. after each top-level player action resolves (`PlayFromHand`, reap, fight, use
   an `Action:`, combat resolution, forging a key);
3. before any `Chooser` decision is presented, so a player never chooses among
   creatures one of which is already dead;
4. at a phase transition that expires this-turn state, chiefly the ready phase
   dropping turn-long power buffs and lasting effects — the one boundary a
   triggered ability does not stand behind, since the expiry is bookkeeping, not
   an ability resolving.

The `settling` re-entrancy flag stays: nested resolutions collapse into the outer
sweep, and a batch that removes several cards at once holds the flag across the
batch and settles once at the end, so a batch still leaves together (ADR 0013,
and the "leave together" rule in `internal/engine/AGENTS.md`). A boundary reached
while a batch already holds the flag collapses into that batch's single closing
sweep.

The low-level writes stop settling themselves — `addAmberOn`, the Æmber and
power-counter writes, the battleline swap and flank moves, taking and releasing
control, and blanking enemy text no longer call `settleDestroyed`; the boundary
owns it. The scattered `settleDestroyed` calls those low-level writes carried are
removed, leaving only the boundaries above (and the batch-closing sweep) as the
sites that settle.

**Contract for new code.** A new mechanic that changes power or moves the board
relies on the boundary and must **not** hand-call `settleDestroyed`. The only code
that settles explicitly is the boundary itself and a batch opening its own
`settling` window. Reaching for a `settleDestroyed` call inside an effect, or
inside a `Game` method that is not one of the boundaries above, is the smell this
decision exists to remove — the framework already settled, or is about to.

## Consequences

- A creature holding lethal damage becomes impossible by construction rather than
  by remembering every call site — the bug class the property tests keep finding
  is closed structurally.
- The `settling` flag remains load-bearing, and the batch-leaves-together rule is
  unchanged; per-change validation is rejected for the reasons above.
- Boundary sweeps run more often in aggregate than the old scattered calls, but
  each is `O(board)` and dwarfed by per-node MCTS cost. The advisory `mage profile`
  baseline is the guard against a surprise regression; it is expected to stay flat.
- Effects and low-level writes get shorter: the settle bookkeeping leaves the hot
  path and lives at the boundary that already knows a resolution just finished.

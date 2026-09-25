---
name: stub-cards
description: Stub a KeyForge set's entire unimplemented backlog in this repo — scaffold a build-excluded //go:build todo file for every card no implemented card covers yet, and regenerate the set's reprint catalog. Use when the user wants to stub a set ("stub the rest of Age of Ascension", "scaffold the backlog", "stub all the cards") before implementing them.
---

Stubbing scaffolds a placeholder for every unimplemented card in a set so the
backlog is a concrete list of files to work through. Each stub is
**build-excluded** — it starts with `//go:build todo`, so it is left out of the
build, vet, test, lint, gencomments, and the card registry. The card database and
`mage tool:coverage` numbers stay honest until a card is actually implemented; to
implement one, remove the build tag and write the real ability (the
**implement-cards** skill).

This is a one-shot setup step. Once a set is stubbed, `mage tool:nextCard` can
hand out the cards one at a time.

**Scan [docs/todo-future-set.md](../../../docs/todo-future-set.md) when stubbing a
set.** It holds decided work parked against a future set — primitives with no
consumer in an implemented set yet. If an item names the set you are stubbing (or
a card it introduces), flag it so the `implement-cards` run builds the primitive
alongside its first real consumer.

## Run it

```sh
mage tool:coverage        # see how many cards each set still has to implement
mage tool:stub <setSlug>  # scaffold a stub for every unimplemented card in the set
```

Set slugs match the files in `internal/cards/provenance/` minus `.json` (e.g.
`callofthearchons`, `ageofascension`). `mage tool:stub` with no working set will
error — name the slug.

`mage tool:stub`:

- Writes `internal/cards/sets/<slug>/<snake>.go` for every source card no
  implemented card covers yet. Each file starts with `//go:build todo` and carries
  the card's printed text, a `// TODO(stub)` marker, and a vanilla `card.New(...)`
  skeleton with the right house, type, rarity, and `card.Provenance(...)` call.
- **Never overwrites** an existing file — an implemented card or an earlier stub is
  left untouched — so it is safe to re-run freely as cards land.
- (Re)generates the set package's `0set.go`, cataloging the cards this set reprints
  from earlier sets so they join its deck-generation pool as full members
  (ADR 0021).

## Verify

```sh
mage ci:build      # the //go:build todo stubs are excluded, so this stays green
mage tool:coverage # the set's total is unchanged; stubs do not count as covered
```

Stubbing changes no coverage number and no registered card — it only lays down the
files. If `mage ci:build` fails after stubbing, it is from something other than the
stubs (they do not compile into the build); investigate that separately.

Hand back once the set is stubbed. Implementing the stubs is the separate
**implement-cards** skill.

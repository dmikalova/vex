# Anomalies

An **anomaly** is a rare card (Worlds Collide onward) that shipped outside a
house's normal pool as a **preview of a future set**: it is really a member of
that later set, seeded early into Worlds Collide packs. So it has two homes
over its life, depending on whether its real set is implemented yet.

## Phase 1 — home set not implemented: park it in Anomaly Expansion

Author it as a **housed** (currently Brobnar) `card.Rarity.Special` card in the
**Anomaly Expansion** set (`internal/cards/sets/anomalyexpansion`,
`card.ReservoirSet(card.AE)`), keeping `card.Provenance(card.WC, "A0x")` so it
counts toward Worlds Collide's coverage. Preserve both properties:

- **The set is a reservoir, never offered for deck generation.** It builds no
  draft pool (`Draftable` is false for every reservoir card), so no deck is
  generated "from" it and it never seeds an interloper or errant house.
- **The cards stay housed, so they remain legacy-drawable.** The legacy pool
  keeps every housed, non-Connected card regardless of the reservoir flag. An
  anomaly authored houseless would drop out of it; never do that.

## Phase 2 — home set implemented: move it there

The anomaly becomes a real member of its home set, authored with the **home
printing's own stats**, not the Brobnar/Special shape:

- Create it in the home set package (`set.New`, so it declares
  `InSet(<home>)`) with the home printing's real **house, type, rarity,
  traits, and wording**, and `card.Provenance(<home>, "<num>")` with its
  home-set collector number, not the `A0x` ref. Look it up in that set's
  provenance JSON.
- Delete the Anomaly Expansion file and its test, and remove any
  `set.Reprint("<num>", "<name>")` the home set's `0set.go` claimed for it; it
  is now a full member. Regenerate with `mage tool:stub <homeSlug>`.

Orb of Wonder is the worked example: it previewed Mass Mutation (Sanctum •
Rare • `Omni:`), so when Mass Mutation was implemented it left Anomaly
Expansion (Brobnar/Special) and became a Sanctum Rare Mass Mutation artifact.
The Shards (Shard of Glory, Shard of Unity) are Vex-invented `Connected` cards
with no future set, so they stay in Anomaly Expansion permanently.

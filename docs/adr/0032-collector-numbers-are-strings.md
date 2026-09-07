# Collector numbers are strings, so lettered reference cards fit

## Context

A provenance `Ref` and a catalog `Card` each carry a **collector number** — the
number printed on the original KeyForge card, used only for coverage bookkeeping
(ADR is bookkeeping-only; the engine and deck generation never read it). The
first model stored that number as an `int`, and the catalog loader parsed the
JSON `"number"` string with `strconv.Atoi`.

Most cards number `001`–`364`, which an `int` holds fine. But a set's **reference
cards** do not use plain integers:

- **anomalies** are numbered `S01`, `S02`, … (they belong to no house's normal run);
- **Worlds Collide** prints its non-deck cards as `A01`–`A21`;
- **prophecies** (Winds of Exchange) are `P01`–`P22`;
- other reference cards (The Tide, archon powers, token creatures) carry their own
  printed numbers.

Under the `int` model every one of these deserialized to `0` — `strconv.Atoi("S01")`
fails and the loader silently kept the zero. Ten Worlds Collide anomalies were
tagged `card.Provenance(card.WC, 0)` because there was no way to name their real
number, and 100+ reference cards across seven sets collapsed onto number `0`,
making them indistinguishable for coverage. Implementing anomalies made the loss
concrete: the anomaly a card was built from could not be pointed at.

## Decision

A collector number is a **string** everywhere: `provenance.Ref.Number`,
`provenance.Card.Number`, `card.Provenance(set, number)`, and
`card.ReprintRef.Number` / `card.Reprint(set, number, name)`. Numeric numbers are
written in their printed, zero-padded form (`"004"`, `"151"`); lettered reference
numbers are written verbatim (`"S01"`, `"A21"`, `"P07"`).

The catalog loader keeps the JSON string as-is — no `Atoi`, no silent zeroing.
Because catalog numbers are zero-padded to a fixed width and letters sort after
digits in ASCII, a plain lexical `<` still orders a set correctly (`001` < `364` <
`A01` < `P01` < `S01`), so no custom comparator was needed.

Where a call-site number (which authors write without padding, `"4"`) must match a
catalog number (`"004"`), the cardlookup tool canonicalizes both through
`normNumber`: an all-digit number drops its leading zeros, a lettered number is
left alone. So `"4"` and `"004"` resolve to the same source card, and `"S01"`
stays `"S01"`.

## Consequences

- Anomalies, prophecies, and other reference cards can be tagged with their real
  printed number; the ten Worlds Collide anomalies now carry `A01`–`A10` instead
  of `0`.
- `card.Provenance` and `card.Reprint` take a quoted number; a bare integer no
  longer compiles. Existing call sites were rewritten to quoted strings.
- The provenance importer (`mage tool:importProvenance`, which replaced the
  removed `mage generateProvenance`) writes numeric numbers zero-padded and
  reference numbers verbatim, so the two forms never drift.

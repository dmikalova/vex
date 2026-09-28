# Debugging a simulator invariant violation

The fixed-seed property batch (`TestSimulateSeeds`, run by `mage ci:test`) and
the fuzz and soak runs report an invariant violation with the script that
reproduces it. How the simulator works is in [testing.md](testing.md)
(Option 4).

## Replay it

- `mage debug` replays the first failing game in the fixed-seed batch and
  prints both players' full deck lists, the log tail, and the violation.
  `mage debug -script=<hex>` replays the script a failure printed;
  `-tail=<n>` widens the log (`mage debug -script=<hex> -tail=200`).
- `mage trace` writes whole games to `tmp/sim/trace.log` (`-count`, `-out`),
  for reading a game end to end.

## Find the cause, not the victim

- **The violation is a symptom.** When an invariant names a card, the cause is
  usually an earlier line and a card no longer in the frame. Read the named
  card's whole lifecycle: widen `-tail` (or `mage trace` the game to a file)
  and grep the log for the card by name to see when it entered, what damage and
  power it showed, and which other card was buffing, blanking, capturing, or
  neighboring it. A creature that dies "for no reason" almost always lost a
  buff a card that has since left play was granting, so identify the cards that
  were in play around it, not only the card the invariant printed.
- **Read the whole deck, not only the cards in the log.** The log names the
  victim, not the culprit, which may never appear in the tail. Scan every card
  in both deck lists for the mechanic that could produce the bad state — a
  power reducer for a 0-power creature, a blanker for a creature that lost its
  ability, an attachment for a stat that drifted. Suspect the mechanic first,
  then find which card in the deck carries it.

## After the fix

Pin the rule with a focused engine or card test, then run `mage corpusPrune`
if the find came from fuzzing or a soak, so the corpus holds only open
findings.

# Browser scenarios (`/ui-test`)

`/ui-test` (`internal/web/uitest.go`, `uitest_scenarios.go`) runs the client in
a real browser against the real DOM. A scenario is data — a name, a slug, a
fixed seed, and ordered steps, each a description plus a do/check — and one
definition serves both the page a human watches and the driver that reads its
status element. How to run the suite and what it costs is in
[testing.md](../testing.md) ("Browser scenarios").

## A scenario is a journey

A scenario is several steps ending in a state change a player would describe:
"deal, mulligan, choose a house, play a creature, answer its prompt, undo it".
**"The reap button is disabled when the creature is exhausted" is not a
scenario**; a single-widget assertion stays a host test in `client_test.go`.
The host tests are the fast fine-grained gate and the place a behaviour is
pinned; the browser suite is the coarse proof the whole thing is wired up.
Every scenario costs seconds of wall clock, so keep them few and keep them
whole journeys.

## The driver

`mage uiTest` runs the same scenarios headlessly. The driver is
`internal/web/uitest`, every file behind the `uitest` build tag: it builds
`web/app.wasm`, serves `cmd/web` on a free port with `VEX_UITEST=1`, launches
headless Chrome through [go-rod](https://go-rod.dev), and opens
`/ui-test/<slug>?once=1` for each entry in `web.UITestScenarios()`, polling
`#ui-test-status` until its `data-state` leaves `running`. **The driver
re-describes no scenario**: the registry is the single definition, so a new
journey is one edit. It is deliberately **not** in `mage ci:check` or
`ci:test`; the build tag keeps the package out of `./...`.

## Rules the surface is built on

- **A step that cannot find its target fails saying so.** `uiPage.find`/`click`
  name what they looked for ("no playable creature in hand"), because on a
  board that never rendered every selector comes back empty, and a scenario
  treating "not there" as "nothing to do" passes without testing anything.
- **Select by hook, never by label.** Scenarios click `data-act` values
  (`actSel(actEndTurn)`, `houseActID`, `optionActID`) and card element ids
  (`boardCardID`, `handCardID`). A control a scenario needs gets a hook in
  `view_card.go`'s `act*` block, beside its neighbours. A card is reached by
  its `data-card` name plus the id prefix naming its zone (`handCardSel`,
  `boardCardSel`); the id alone is no use, because a staged card gets the next
  free `LocalID`, which the scenario cannot know.
- **A scenario stages the board it needs; it does not read the deal.** The
  deal for a fixed seed changes whenever a card is implemented. A journey that
  needs a named card turns manual mode on through the real menu and adds the
  card through the real picker (`manualPreamble`). Only journeys whose subject
  *is* the deal (`opening`, `mulligan`) read the board with predicates. **No
  step names a physical side**: steps are written against the active player.
- **Closing and reopening the page is the host's, not a click's.** A real
  reload would restart the run, so `uiPage.dropClient` takes the client out of
  the tree and `mountClient` stands a fresh one up over the same storage slot,
  driving the client's own mount and `resume`.
- **Each pass resets.** The run clears the ui-test storage namespace and
  re-deals from the scenario's seed before step 1.
- **The run is isolated from a real match.** The driven client gets an
  injected `storeKey` (`uiTestStoreKey`) and `fixedSeed`; all persistence goes
  through `matchKey()`, so a run cannot overwrite a playtester's open game.
  Both are empty/zero in normal play; do not reach around them with a second
  fixed key.
- The routes are gated on `VEX_UITEST` (`UITestEnabled`, `UITestRoutes`) and
  registered on both sides of the build, exactly as `/style` is gated on
  `VEX_STYLE`: go-app routes on the client, so a build tag would drop the
  scenarios from the bundle and they would rot uncompiled (ADR 0014).
  `mage web` sets both.

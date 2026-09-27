# AGENTS.md — `internal/web`

The browser client: a [go-app](https://go-app.dev) v11 WASM front end over
`internal/engine`. Read the repo-root `AGENTS.md` first; this file covers only
what is specific to the web client.

Build it the same way as everything else: `mage ci:build` compiles every package
for the host **and** the client for js/wasm (vex's `ci.ExtraBuilds`), so a change
that only breaks in the browser fails the same gate as anything else. To compile
just the client into `web/app.wasm` after an edit, run `mage webWasm`.

## File split

Each file has one job; put new code in the one that matches. The package splits
the same way `internal/engine` does: `game_*.go` holds the component's behaviour
grouped by area, `view_*.go` holds rendering grouped by screen region.

- `card.go` — `cardView`, the presentational card-face component, plus the `cx`
  and `ifCls` class helpers.
- `game.go` — the root `game` component: the `phase`/`selKind` enums, the state
  struct and its supporting value types, and the smallest readers over them.
- `game_lifecycle.go` — mount/dismount, post-render scrolling, keyboard
  shortcuts, and the hot-reload hand-off.
- `game_persist.go` — saving a match to local storage, resuming it, and dealing a
  new one when there is nothing to resume.
- `game_action.go` — the action plumbing: applying a Command to the session,
  settling what it yields, undo/redo over the command log, and the flash/flight
  bookkeeping after each action.
- `game_play.go` — taking a turn: selection, house choice, play (click and drag),
  reap, use, fight, end turn.
- `game_chooser.go` — the prompt seam: the readers over the session's pending
  `Request` and the handlers that answer one.
- `game_manual.go` — manual mode: manual moves, stat adjustments, card picker.
- `game_ui.go` — client-only state no rule touches: hover preview, restart
  confirmation, sidebar toggle.
- `view.go` — the outermost frame: page layout, brand bar, status banner.
- `view_board.go` — the battlelines, the score pills, and the cards in play.
- `view_hand.go` — the active player's hand.
- `view_controls.go` — the sidebar: prompts, the action bar, the manual-mode panel.
- `view_log.go` — the game log and its card-name links.
- `view_overlay.go` — what covers the board: the zone viewer and the result panel.
- `view_card.go` — the shared face helpers (labels, stat lines, rules text) and
  the small odds and ends the views share.
- `icons.go` — SVG asset lookup (`icon`, `houseIconName`, `typeIconName`, …) and
  the injected `#icon-outline` filter.
- `palette.go` — house → CSS class mapping only. Colours live in `web/app.css`.
- `uitest.go` — the `/ui-test` scenario host: the run loop, the step probe over the
  real DOM, and the route switch.
- `uitest_scenarios.go` — the scenarios themselves, as data. A new journey is an
  entry here, not a new page. `UITestScenarios` is the exported view of that
  registry the headless driver enumerates.

The headless driver is the sibling package `internal/web/uitest`, every file
behind the `uitest` build tag. It is only a driver.

Styles live in `web/app.css`; assets are `web/assets/<stem>.svg`, referenced by
stem through `icon(name, extra…)` and served from `/web/assets/`. The dev server
reads both from disk, so a CSS or SVG change needs no rebuild.

## The engine owns the rules; the client only draws them

Never reimplement a rule in the client. If the client needs to know whether
something is legal, ask the engine (`g.g.FightTargets`, `g.g.CanPlay`,
`g.g.RestrictionSources`, …). If the reader does not exist, add it to
`internal/engine` rather than reaching into `GameState` and inferring the answer
— an inference here silently drifts from the rules the engine enforces.

## `cardView` is presentational

`cardView` carries no game logic. The parent hands it already-rendered strings
(`Rules`, `Trait`, `Kind`), visual flags (`Selected`, `Targetable`, `Dimmed`,
`Exhausted`, `Enter`, `Fight`), and handlers. That is why the same component
renders a hand card, a creature in play, an artifact, a zone card, and a prompt
preview. When a card needs to show something new, add a flag to `cardView` and
compute it in the `view_*.go` file that builds that face — do not give `cardView`
a `*engine.Game`.

Compose classes with `cx(...)` and `ifCls(cond, "class")` rather than string
concatenation, so an unset modifier contributes nothing.

## Handlers are stable methods, and take the id as an argument

go-app compares event handlers by function pointer, so a per-card closure is
re-created on every render and the diff cannot keep it bound. Two consequences:

- Event handlers on a component are **methods** (`c.onClick`), not closures.
- The clicked card's identity is passed to the parent through the component's own
  field (`OnActivate(ctx, c.ID)`), never captured in the closure.

Where a handler genuinely needs a value that is not on the component, read it back
out of the DOM (`ctx.JSSrc().Get("dataset")…`) — that is what `onScorePillClick`
and the log's card mentions do — rather than minting a closure per element.

## `LocalID` 0 is a real card, so no zero sentinels

The engine hands out `LocalID`s from 0, so the first card dealt has id 0 and a
zero id cannot mean "nothing". Client state that may hold no card carries its own
flag beside the id (`hasSel`/`sel`, `hasHover`/`hoverID`, `hasCursor`/
`promptCursor`). A `!= 0` test silently drops that one card — which is how the
hover preview came to never show for it.

## One-shot animations use the `-a`/`-b` parity trick

A CSS animation only replays when its `animation-name` changes, and go-app patches
the existing element instead of replacing it. So every one-shot effect is defined
twice in `web/app.css` with identical keyframes under two names (`cardEnterA` /
`cardEnterB`), and the Go markup alternates between `--enter-a` and `--enter-b`
using a parity bit that flips on each flash. `cardFlash.odd` (per card),
`poolParity`, `keyParity`, and `discardParity` are those bits.

Flashes are **derived, not emitted**: `computeFlashes` diffs the pre-action
snapshot `beginAction` took against the resolved state. Add a new
animation by diffing the state there, not by sprinkling calls at action sites.
The exception is information the state does not carry — a fight's two
combatants, which card reaped, which card used an action ability — which the
handler arms on `g.fighters`/`g.reapID`/`g.actID` for `computeFlashes` to
consume.

A card that has left play cannot pulse, so it **flies** instead: `computeFlights`
finds the zone it landed in and `flightsInto` renders a ghost face parented to
that zone's pill, which arcs in and shrinks onto the count (`.card-flight`).

## Every click is a Command; the session owns the turn loop

The client drives one `*session.Session` (ADR 0039, 0040) whose action is the
engine's own `RunMatch`. **The client holds no `Chooser`** — the engine never
calls into it. So the client takes no turn of its own and starts no goroutine: a
click builds an `engine.Command` and applies it, which resolves the engine as far
as the next decision **before the handler returns**. Everything is synchronous,
and there is no in-flight action to guard against.

- A root action goes through `applyRoot`, a prompt answer through `answer`, a
  manual force-edit through `applyManual`. Nothing else may call the engine to
  change state — a mutation that skips the session is a mutation the command log
  does not hold, so it does not survive a reload or an undo.
- **Undo is `Session.Undo(n)`**, which rewinds by replaying the first `n` recorded
  commands from a fresh deal — not a client-side stack of `GameState` copies. The
  client keeps only the marks that say where each root action started
  (`rootMarks`) and the commands it truncated, so `redoAction` can re-apply them;
  the session has no `Redo` of its own.
- **The root boundary is "`Pending()` is a `RequestAction`."** Between two of them
  lies exactly one root action and every prompt it raised. That one signal drives
  the undo marks, the flash baseline, and the log groups: `beginAction` opens an
  action and `settleAfterApply` closes it when the next `RequestAction` arrives.
- `rootMarks` and `logGroups` are appended together and only by `beginAction`, so
  the i-th entry of each describes the same action and one undo peels one entry
  off each.
- `atPrompt()` is the single input guard: a pending request that is not a
  `RequestAction`. Most handlers begin with `if g.atPrompt() … { return }`.
- **Cancelling a prompt and undoing at one are the same call**, `undoAction`,
  which rewinds to the mark the raising action began at. The Cancel button is a
  manual-mode-only second trigger for it.
- What a prompt asks is **read** from `session.Pending()`, never pushed at the
  client: `choosing()`, `optionLabels()`, `positionLine()` and friends in
  `game_chooser.go` derive from the live `Request`, and the badge preview from
  `Request.Badge`.
- **Anything the UI needs that is not an answer travels as display-only context on
  `Request`** (ADR 0047): a new field the engine populates at the suspension point
  and the client reads off `Pending()`, with no mirror in `Command`, so it cannot
  enter the command log or change what answers are legal. Do **not** add a callback
  field on `Session` or an interface the client implements for the session to call
  — either one rebuilds the live chooser this architecture exists to remove.

## Prompts: cards are clicked, options are buttons

- A card decision is a **card prompt**: the candidates highlight on the board and
  the controls become the prompt text. Mandatory ones (`ChooseCreature`) have no
  way out; optional ones (a `RequestPickCardOrDecline`, which is what
  `chooserDeclinable()` reads) get a **Done** button, and Escape declines them.
- A genuine yes/no or "choose one" stays an **option prompt** with buttons.
- **A trigger window is ordered by clicking too** (a `RequestReaction`): the sources
  of the pending abilities become the candidates, and the click says whose ability
  resolves next. A clicked card carrying **two** pending abilities then asks which
  of them, by buttons — never silently top-down. A source sitting in a pile
  (`WithTriggersFromDiscard`) is a candidate like any other: `presentPrompt` opens
  that pile's viewer so the reaction is not missed. Only a window holding an entry
  that belongs to **no card** — a duration reaction — falls back to the flat
  labeled list of rendered ability text.
- When a prompt's candidates are not on the board, `presentPrompt` routes it:
  a bounded, mandatory "look at the top N cards of your deck" pick (Navigator Ali,
  Lay of the Land) becomes a short list of **action-bar buttons**
  (`promptCardButtons`), each naming its card and previewing it on hover; every
  other out-of-play pick — a visible pile (World Tree, Witch of the Eye's discard),
  or an unbounded one (declinable — Not Finished with You shuffles any number) —
  opens that player's **zone viewer**, which makes only the candidates clickable,
  dims the rest, and scrolls to the row. Candidates split between a pile and the
  board open the pile (`firstPileCandidate`).
- **The zone viewer is always closable**, prompt or no prompt: the player may need
  to read the board underneath before deciding. Closing answers nothing — the
  prompt stays up, `promptZone` stays recorded, and reopening the viewer returns to
  the same row. Passing on a declinable prompt is the **Done** affordance's job (in
  the viewer's header and in the dock), never a side effect of closing the modal.
- Picking a fight target is not a card prompt (the client asks it before it emits
  the `CommandFight`, so the engine never sees a decision), but it shares the same
  Tab cursor: `tabCandidates` hands
  both a chooser's candidates and `FightTargets` to `tabCandidate`/`isSelected`/
  `confirmPrompt` through one seam, so a fix to how the cursor lands or draws
  covers both instead of needing a matching fix in a second, fight-shaped copy.

## Phases and Escape layering

`phase` is the client's own interaction state, not a rules concept:
`phaseHouse` → `phaseMain` → (`phaseFlank` | `phaseFightTarget`) → `phaseOver`.

A Deploy creature skips `phaseFlank`: playing it runs straight away, and the
engine raises its `ChoosePosition` prompt (`choosingPosition`). The creature is
lifted while the prompt is up and its placement verbs sit on the lifted card
(`deployActions`), like the flank question but extended for a creature that may
enter anywhere: ordered left-to-right `[left flank] [deploy left] [deploy right]
[right flank]`, the ends answer at once and the interior`deploy left`/`deploy
right` arm which side of a clicked battleline creature it lands on (the line
lights its creatures as the click targets). With no other friendly creatures in
play there is only one placement, so `ChoosePosition` answers itself and no
prompt is shown. A plain creature on a non-empty line still asks the which-flank
question in `phaseFlank` first.

Escape calls `dismiss`, which backs out **exactly one layer**, innermost first:
picker → zone viewer → restart confirmation → key-forge picker → an escapable
prompt → mid-action targeting → end-turn confirmation → the selection. Add a new
overlay to that chain rather than giving it its own Escape handling.

Keyboard shortcuts live in `onKey`. They are all single keys and all no-ops while
a prompt is up (`atPrompt()`), except Escape, `r`, and `n`, which are handled first so they
work while a prompt blocks every other key: Escape backs a prompt out, `n`
answers a declinable one "no", and `r` reaps with the selected creature (or
takes the right flank while placing, or forges a Red key when that prompt is
up) — it is not a general "yes"; only Space's fallback (`affirm`) plays the
selected hand card, uses an artifact's action, or answers a yes/no prompt.

## Persistence

The in-progress match is stored in local storage under `persistKey`, tagged with
`snapshotVersion`. What is written is the session's `Record` — version, seed,
sets, and the ordered command log — and nothing else: state, the typed log and
the log bubbles are all replayed from it (`replayRecord`). **Bump
`snapshotVersion` whenever an engine change makes an older record unreplayable**
— a stale snapshot is dropped rather than restored into a mismatched engine.

That includes rewording a log entry. The log persists as the prose each entry was
narrated with (a typed entry does not survive JSON), so an old snapshot keeps
restoring the old wording long after the engine stopped producing it, and the
change looks like it did not take. What a bump costs, and why that means log
fixes are never batched, is on the `snapshotVersion` doc comment.

That also includes any change to how a recorded command resolves, because a
resume replays the log rather than deserializing state: the same commands would
rebuild a different match. `session.Version` guards the command format itself,
and `replayRecord` refuses a mismatch the same way `session.Load` does.

Manual force-edits are in the record **on purpose**: a force-edited board that
hits a bug is the reproduction worth keeping.

A failed write is not swallowed. `writeSnapshot` frees the storage the client can
spare — the style gallery's scroll memo and the snapshot the write is replacing —
retries once, and only then raises a standing notice that the match is no longer
being saved.

## Two banners: transient status, standing notice

The control dock carries two messages, and they are not interchangeable:

- `setStatus` is a **transient** message that fades after 5s. Use it for a
  rejected click the player can simply retry.
- `setNotice` is a **standing** message that stays up until `clearNotice` takes it
  down. Use it for a fault the player has to act on — storage being full, a set
  the deck generator does not know — where the 5s fade would let it pass unread.

The banner holds one notice, so a code path that clears a notice clears only its
own (see `clearStorageNotice`).

## Tests drive the real client off-browser

`client_test.go` is the harness. It plays the client the way a person does — the
same handlers a click or a key press calls, over a real engine game — so add a
test there rather than reimplementing a slice of the client to assert against.

Two go-app facts make it work, and neither is guessable:

- `app.Context` is a **struct**, not an interface, so it cannot be faked. One has
  to be borrowed. go-app only fires `OnMount` when `app.IsClient` is true, which
  it is not in a host test, but it fires `OnPreRender` when `app.IsServer` is
  true, which it is. So `ctxProbe` implements `OnPreRender` purely to be handed a
  live context, complete with working in-memory local storage and dispatch queue.
  `e.ConsumeAll()` then drains dispatches, asyncs, and deferred work, so the
  renders and timers a handler queued have run when it returns. The game itself is
  never waited for: applying a Command resolves the engine to its next decision
  before the handler returns, so `c.settle`/`c.await` assert rather than wait.
- A zero `app.Event{}` nil-derefs on `PreventDefault()`. Use `nullEvent()`.

`app.HTMLString(g.Render())` draws the whole screen, nested components included,
without mounting anything — which is how `view_test.go` asserts on markup.

A test that needs one exact prompt raises it with `c.script`, which stands the
client on a session whose driving action is that single question instead of the
whole turn loop. The client then sees a genuine pending `Request` and answers it
with a genuine `Command`, so what is under test is the handoff rather than a
rehearsal of it. A test whose subject is backing OUT of a prompt uses a real play
instead (`c.stagePromptArtifact`), because there has to be a root action to rewind.

**Never assert a fixed attribute order on the drawn markup.** go-app writes an
element's attributes by ranging its `attrs()` **map**, so their order is
randomised per render: a fragment that fixes `class` before `src` on an icon
matches only some of the time and is a flaky test waiting to happen. Match a
single attribute (`src="/web/assets/shield.svg"`), or bind a value to its icon
with a regex whose `[^>]*` skips the attributes without fixing their order, as
`TestArmorShowsWhatIsLeftToAbsorb` does.

What is out of reach is the DOM: `app.Window()` reads back empty off-browser, so
the pieces that measure or scroll elements (the fly-into-play animation, the log
auto-scroll, the picker's focus, `ctx.JSSrc()`) no-op rather than assert. Every
one of them is written to tolerate a render with no page behind it, which is what
makes the rest of the client testable here at all — keep it that way.

Coverage here is **deliberately ungated** (absent from `ci.CoverGates` in
`magefiles/build.go`), because the last stretch is the DOM-bound code above.

## Browser scenarios are journeys, not widget assertions

`/ui-test` (`uitest.go`, `uitest_scenarios.go`) runs the client in a real browser
against the real DOM: a scenario is data — a name, a slug, a fixed seed, and
ordered steps, each a description plus a do/check — and one definition serves both
the page a human watches and the driver that reads its status element.

A scenario is a **journey**: several steps ending in a state change a player would
describe. "Deal, mulligan, choose a house, play a creature, answer its prompt,
undo it" is a scenario. **"The reap button is disabled when the creature is
exhausted" is not** — a single-widget assertion stays a host test in
`client_test.go`. The host tests are the fast fine-grained gate and stay the place
a behaviour is pinned; the browser suite is the coarse proof that the whole thing
is wired up in a browser, and every scenario added there costs seconds of wall
clock, so keep them few and keep them whole journeys.

`mage uiTest` runs the same scenarios headlessly. The driver is
`internal/web/uitest`, every file behind the `uitest` build tag: it builds
`web/app.wasm`, serves `cmd/web` on a free port with `VEX_UITEST=1` (`PORT` is
the env var `cmd/web` already reads), launches headless Chrome through
[go-rod](https://go-rod.dev) — pure Go, no Node — and opens
`/ui-test/<slug>?once=1` for each entry in `web.UITestScenarios()`, polling
`#ui-test-status` until its `data-state` leaves `running`. **The driver
re-describes no scenario**: the registry here is the single definition, so a new
journey lands in the page and in the suite with no second edit.

It is deliberately **not** in `mage ci:check` or `ci:test` — the build tag keeps
the package out of `./...`, so the shared CI workflow needs no browser — and it
sits beside `mage profile` and `mage trace` as a real-but-ungated target. Costs,
measured warm: the js/wasm build ~3 s (28 MB) and free when nothing changed,
Chrome launch ~1 s, the first page load of the bundle ~1-2 s and ~0.5 s after.
The steps run **in** the page rather than over CDP, so a step is a click plus a
render (single-digit ms): today's suite is ~5 s warm (~13 s cold, including
go-rod downloading its Chromium), a minimal one ~10-15 s, and a ~50-journey
suite ~30-60 s, or ~15-25 s if passes share a page load. That budget is the
reason a single-widget assertion stays a host test. See
[docs/testing.md](../../docs/testing.md).

Four rules the surface is built on:

- **A step that cannot find its target fails saying so.** `uiPage.find`/`click`
  name the thing they were looking for ("no playable creature in hand"), because
  off-browser — and on a board that never rendered — every selector comes back
  empty, and a scenario that treats "not there" as "nothing to do" passes without
  testing anything.
- **Select by hook, never by label.** Scenarios click `data-act` values
  (`actSel(actEndTurn)`, `houseActID`, `optionActID`) and card element ids
  (`boardCardID`, `handCardID`). A control a scenario needs to reach gets a hook
  in `view_card.go`'s `act*` block, beside its neighbours. A card is reached by
  the `data-card` name every face carries plus the id prefix that says which zone
  drew it (`handCardSel`, `boardCardSel`) — the id alone is no use, because a card
  a scenario staged is given the next free `LocalID`, which the scenario cannot
  know.
- **A scenario stages the board it needs; it does not read the deal.** A
  `deckgen.Set` is built from the implemented card list, so the deal for a fixed
  seed changes every time a card is implemented. A journey that needs a named card
  turns manual mode on through the real menu and adds that card through the real
  picker (`manualPreamble`), then acts on the card it chose. Only the journeys
  whose subject *is* the deal — `opening`, `mulligan` — skip that and read the
  board with predicates. And **no step names a physical side**: every step is
  written against the active player, so which player goes first stays the engine's
  to decide.
- **Closing and reopening the page is the host's, not a click's.** A real browser
  reload would restart the run rather than the match, so `uiPage.dropClient` takes
  the client out of the tree and `mountClient` stands a fresh one up over the same
  storage slot. What that drives is the client's own mount and `resume`, which is
  the path a reload takes.
- **Each pass resets.** The run clears the ui-test storage namespace and re-deals
  from the scenario's seed before step 1, because a scenario that plays a creature
  cannot run again on the board it left behind.
- **The run is isolated from a real match.** The client the scenario drives gets an
  injected `storeKey` (`uiTestStoreKey`) and `fixedSeed`; everything persistence
  touches goes through `matchKey()`, so a run in a browser cannot overwrite a
  playtester's open game. Both fields are empty/zero in normal play — do not reach
  around them with a second fixed key.

The routes are gated on `VEX_UITEST` (`UITestEnabled`, `UITestRoutes`) and
registered on both sides of the build, exactly as `/style` is gated on
`VEX_STYLE`, and for the same reason: go-app routes on the client, so a build tag
would drop the scenarios from the bundle and they would rot uncompiled
(ADR 0014). `mage web` sets both.

## CSS conventions (`web/app.css`)

- BEM-ish: a block (`.card`), and modifiers as `--modifier` classes
  (`.card--dimmed`, `.log-group--p0`). No inline styles from Go.
- House colours are custom properties (`--nm`, `--tp`, `--edge`) supplied by the
  `.card-<house>` class from `palette.go`; markup only ever carries class names.
- The one exception is a **measurement**: where a particular card sits on screen
  is a runtime fact no class can name, so `view_focus.go` hands the lifted card's
  rect over as custom properties (`--focus-x`, `--focus-w`, …) and `app.css` owns
  everything done with them. Values, never styling — a rule that needs a new
  declaration in Go is a rule that belongs in the stylesheet.
- Keep every animation's `-a`/`-b` pair in sync — they must have identical
  keyframes.
- **Enlarge a card by resizing its box, not by `transform: scale()`.** The reason
  to enlarge a card is to read the text the board was clipping, and a transform
  does not reflow text. Resizing one means setting three things together, the way
  `.card-preview` and `.card-focus` both do: the box and the inner font sizes,
  which are hardcoded px and so do not scale with the parent. A
  card whose height follows its content also has to undo the
  `flex: 1 1 0%; min-height: 0` on `.card-body`/`.card-rules`, which exist to fill
  and clip a fixed slot.
- **The card frame is a flat two-tone gradient, not a clip or a mask.** `.card`
  paints a hard-stop `linear-gradient` — the name-banner colour (`--nm`) over the
  upper-left, the type band (`--tp`) over the lower-right, split by a straight
  diagonal. A gradient scales with the card at any size and respects
  `border-radius`, so the card needs no `overflow: hidden` and no SVG mask — which
  is what lets the left-edge house and bonus icons (`.card-bonuses`) hang off the
  card. The keybar and `.card-kind` round their own corners, since the parent no
  longer clips them. The `125deg` angle and `48%` stop in `.card` are the two knobs.
- **A card that overhangs its row has to leave the board's coordinate space**
  (`position: fixed`): `.card-strip` is `overflow-x: auto`, which per spec forces
  `overflow-y: hidden`, and `.board-area` is `overflow: hidden` — anything inside
  either one is clipped. The same trap catches a tooltip: `.score-pill` is
  `overflow-x: auto` (it scrolls sideways), so a per-icon `::after` bubble is
  clipped inside the bar. Every `data-tip` shares **one** floating label
  (`#tip-float`, `installTips`) placed with `position: fixed` from a measurement of
  the element under the pointer — one fixed element the pointer fills, not a bubble
  per icon.
- **Place a content-sized overlay from the edge it is nearest**, not from its own
  centre. Its height is not known until it has been laid out, and measuring it
  needs a second render pass that a frozen tab (or any dropped frame) will not
  give you. `.card-focus` anchors `top` for a card in the opponent's half and
  `bottom` for one in the player's, so however tall the copy turns out to be it
  grows away from the near edge instead of through it.
- **An overlay placed from a measurement must track everything that moves it.**
  `.card-focus` re-measures from go-app's `Resizer` (`OnResize`) and from a
  document-level `scroll` listener registered in the **capture** phase — a scroll
  event does not bubble, and every card strip scrolls on its own.
- **A copy of a card that covers its neighbours has to take the pointer.**
  `pointer-events: none` on `.card-focus` looked like a free way to let the wheel
  through, but it also let a drag through: the grab landed on whichever card the
  enlarged face happened to lie over, not on the card being enlarged. The face is
  `pointer-events: auto` and is its own drag source, and the wheel is handed on to
  the strip underneath in Go (`wheelOverFocus`).
- **Keep `filter` off the ancestor of a drag source.** A filtered ancestor stops
  the dragged element being its own layer, and the browser then cuts the drag
  image out of the ancestor — pulling in a sliver of whatever sits next to it.
  `.card-focus` gives its drop-shadow to the face and the verbs separately.

## Never shout: no all-caps

Nothing in the client is set in capitals. No `text-transform: uppercase`, no
`font-variant-caps: all-small-caps`, and no string typed in caps in Go to be
read as a label. Card types, row labels, zone names, and log headers all render
in the casing they are written in. Size, weight, and colour already separate a
label from prose; capitals only make it harder to read and louder than what it
labels.

The wide `letter-spacing` that usually accompanies caps goes with them — it
exists to make capitals legible and is vestigial without them.

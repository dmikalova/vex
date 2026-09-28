# AGENTS.md — `internal/web`

The browser client: a [go-app](https://go-app.dev) v11 WASM front end over
`internal/engine`. The gate's build compiles it for js/wasm too, so a
browser-only break fails the gate; `mage webWasm` compiles just the client into
`web/app.wasm`.

- Before editing `web/app.css` or how a card is laid out on screen, read
  [docs/web/css.md](../../docs/web/css.md).
- Before adding or changing a `/ui-test` browser scenario, read
  [docs/web/browser-scenarios.md](../../docs/web/browser-scenarios.md).

## File split

One job per file; put new code in the matching one. As in `internal/engine`,
`game_*.go` holds behaviour by area and `view_*.go` rendering by screen region.

- `card.go` — `cardView`, the presentational card face, plus the `cx` and
  `ifCls` class helpers.
- `game.go` — the root `game` component: the `phase`/`selKind` enums, the state
  struct and its value types, the smallest readers.
- `game_lifecycle.go` — mount/dismount, post-render scrolling, keyboard
  shortcuts, the hot-reload hand-off.
- `game_persist.go` — saving to local storage, resuming, dealing a new match
  when there is nothing to resume.
- `capture.go` — writing a replay failure to a file an agent can replay: the
  `Capture` record, its staleness stamps, the dev-server endpoint, `Replay`.
- `dev.go` — `DevEnv`/`DevEnabled`, the one switch the wasm client and the
  native server both read for development-only surfaces.
- `game_action.go` — applying a Command to the session, settling what it
  yields, undo/redo over the command log, flash/flight bookkeeping.
- `game_play.go` — taking a turn: selection, house choice, play (click and
  drag), reap, use, fight, end turn.
- `game_chooser.go` — the prompt seam: readers over the pending `Request` and
  the handlers that answer it.
- `game_manual.go` — manual mode: manual moves, stat adjustments, card picker.
- `game_ui.go` — client-only state no rule touches: hover preview, restart
  confirmation, sidebar toggle.
- `view.go` — page layout, brand bar, status banner.
- `view_board.go` — the battlelines, score pills, cards in play.
- `view_hand.go` — the active player's hand.
- `view_controls.go` — the sidebar: prompts, action bar, manual-mode panel.
- `view_log.go` — the game log and its card-name links.
- `view_overlay.go` — the zone viewer and the result panel.
- `view_card.go` — shared face helpers (labels, stat lines, rules text) and
  small shared odds and ends.
- `icons.go` — SVG asset lookup (`icon`, `houseIconName`, `typeIconName`, …)
  and the injected `#icon-outline` filter.
- `palette.go` — house → CSS class mapping only; colours live in `web/app.css`.
- `uitest.go` / `uitest_scenarios.go` — the `/ui-test` host and its scenarios
  as data (`UITestScenarios`); the headless driver is `internal/web/uitest`,
  behind the `uitest` build tag.

Styles live in `web/app.css`; assets are `web/assets/<stem>.svg`, referenced by
stem through `icon(name, extra…)` and served from `/web/assets/`. The dev server
reads both from disk, so a CSS or SVG change needs no rebuild.

The client's UI areas have names (**player bar**, **play zone**, **board row**,
**row label**, **midline**, **hand row**, **zone counts**, **sidebar**,
**game log**, **turn HUD**, **prompt**, **action bar**, **zone viewer**,
**card preview**); use them, and see `CONTEXT.md` for the rest.

## The engine owns the rules; the client only draws them

Never reimplement a rule in the client. Ask the engine whether something is
legal (`g.g.FightTargets`, `g.g.CanPlay`, `g.g.RestrictionSources`, …); if the
reader does not exist, add it to `internal/engine` rather than inferring the
answer from `GameState`.

## `cardView` is presentational

`cardView` carries no game logic: the parent hands it rendered strings
(`Rules`, `Trait`, `Kind`), visual flags (`Selected`, `Targetable`, `Dimmed`,
`Exhausted`, `Enter`, `Fight`), and handlers, so one component renders every
face. To show something new, add a flag and compute it in the `view_*.go` file
that builds that face; never give `cardView` a `*engine.Game`. Compose classes
with `cx(...)` and `ifCls(cond, "class")`, not string concatenation.

## Handlers are stable methods, and take the id as an argument

go-app compares handlers by function pointer, so a per-card closure is
re-created every render and cannot stay bound.

- Handlers on a component are **methods** (`c.onClick`), not closures.
- The clicked card's identity reaches the parent through the component's own
  field (`OnActivate(ctx, c.ID)`), never a captured variable.
- A value not on the component is read back from the DOM
  (`ctx.JSSrc().Get("dataset")…`), as `onScorePillClick` and the log's card
  mentions do, rather than minting a closure per element.

## `LocalID` 0 is a real card, so no zero sentinels

Client state that may hold no card carries its own flag beside the id
(`hasSel`/`sel`, `hasHover`/`hoverID`, `hasCursor`/`promptCursor`). A `!= 0`
test silently drops card 0 (the hover preview once never showed for it).

## One-shot animations use the `-a`/`-b` parity trick

A CSS animation replays only when its `animation-name` changes, and go-app
patches elements in place. So every one-shot effect is defined twice in
`web/app.css` with identical keyframes (`cardEnterA` / `cardEnterB`), and the
markup alternates `--enter-a` / `--enter-b` on a parity bit that flips per
flash (`cardFlash.odd`, `poolParity`, `keyParity`, `discardParity`).

- Flashes are **derived, not emitted**: `computeFlashes` diffs the snapshot
  `beginAction` took against the resolved state. Add an animation by diffing
  there, not by calls at action sites. Only what the state does not carry (a
  fight's combatants, the reaping card, the card using an action) is armed by
  the handler on `g.fighters`/`g.reapID`/`g.actID`.
- A card that has left play **flies** instead: `computeFlights` finds the zone
  it landed in and `flightsInto` renders a ghost face on that zone's pill
  (`.card-flight`).

## Every click is a Command; the session owns the turn loop

The client drives one `*session.Session` (ADR 0039, 0040) whose action is the
engine's `RunMatch`. **The client holds no `Chooser`**, takes no turn of its
own, and starts no goroutine: a click builds an `engine.Command` and applies
it, resolving the engine to its next decision **before the handler returns**.

- A root action goes through `applyRoot`, a prompt answer through `answer`, a
  manual force-edit through `applyManual`. Nothing else calls the engine to
  change state; a mutation that skips the session is missing from the command
  log and survives neither reload nor undo.
- **Undo is `Session.Undo(n)`**, replaying the first `n` commands from a fresh
  deal, not a stack of `GameState` copies. The client keeps only where each
  root action started (`rootMarks`) and the commands it truncated, for
  `redoAction`; the session has no `Redo`.
- **The root boundary is "`Pending()` is a `RequestAction`."** That signal
  drives the undo marks, the flash baseline, and the log groups: `beginAction`
  opens an action, `settleAfterApply` closes it at the next `RequestAction`.
  `rootMarks` and `logGroups` are appended together, only by `beginAction`, so
  one undo peels one entry off each.
- `atPrompt()` (a pending request that is not a `RequestAction`) is the single
  input guard; most handlers begin with `if g.atPrompt() … { return }`.
- **Cancelling a prompt and undoing at one are the same call**, `undoAction`;
  the Cancel button is a manual-mode-only second trigger.
- What a prompt asks is **read** from `session.Pending()` (`choosing()`,
  `optionLabels()`, `positionLine()` in `game_chooser.go`, the badge preview
  from `Request.Badge`), never pushed at the client.
- **Anything the UI needs that is not an answer is display-only context on
  `Request`** (ADR 0047), with no mirror in `Command`. Never add a callback on
  `Session` or an interface the client implements for the session to call;
  either rebuilds the live chooser this design removes.

## Prompts: cards are clicked, options are buttons

- A card decision is a **card prompt**: candidates highlight on the board and
  the controls show the prompt text. Mandatory ones (`ChooseCreature`) have no
  way out; optional ones (`RequestPickCardOrDecline`, read by
  `chooserDeclinable()`) get a **Done** button, and Escape declines them.
- A genuine yes/no or "choose one" is an **option prompt** with buttons.
- **A trigger window (`RequestReaction`) is ordered by clicking**: the sources
  of pending abilities are the candidates. A source with **two** pending
  abilities then asks which by buttons, never silently top-down. A source in a
  pile (`WithTriggersFromDiscard`) is a candidate too: `presentPrompt` opens
  that pile's viewer. Only a window holding an entry that belongs to **no
  card** (a duration reaction) falls back to a labeled list of ability text.
- Off-board candidates route through `presentPrompt`: a bounded, mandatory
  "look at the top N cards of your deck" pick (Navigator Ali, Lay of the Land)
  becomes **action-bar buttons** (`promptCardButtons`) that preview on hover;
  any other out-of-play pick (a visible pile, or an unbounded declinable one)
  opens that player's **zone viewer** with only the candidates clickable, the
  rest dimmed, scrolled to the row. Candidates split between a pile and the
  board open the pile (`firstPileCandidate`).
- **The zone viewer is always closable**; closing answers nothing (the prompt
  and `promptZone` stay, and reopening returns to the row). Declining is the
  **Done** affordance's job, never a side effect of closing.
- A fight target is not a card prompt (the client asks it before emitting
  `CommandFight`) but shares the Tab cursor: `tabCandidates` feeds both a
  chooser's candidates and `FightTargets` to
  `tabCandidate`/`isSelected`/`confirmPrompt`, so one fix covers both; keep it
  one seam.

## Phases and Escape layering

`phase` is client interaction state, not a rules concept:
`phaseHouse` → `phaseMain` → (`phaseFlank` | `phaseFightTarget`) → `phaseOver`.

- A Deploy creature skips `phaseFlank`: playing it runs at once and the engine
  raises `ChoosePosition` (`choosingPosition`). The lifted creature carries the
  placement verbs (`deployActions`), left to right `[left flank] [deploy left]
  [deploy right] [right flank]`: the ends answer at once; the interior two arm
  which side of a clicked battleline creature it lands on. With no other
  friendly creature in play `ChoosePosition` answers itself. A plain creature
  on a non-empty line asks the flank question in `phaseFlank`.
- Escape calls `dismiss`, which backs out **exactly one layer**, innermost
  first: picker → zone viewer → restart confirmation → key-forge picker → an
  escapable prompt → mid-action targeting → end-turn confirmation → the
  selection. Add a new overlay to that chain, not its own Escape handling.
- Keyboard shortcuts live in `onKey`: single keys, all no-ops while a prompt is
  up (`atPrompt()`), except Escape, `r`, and `n`, handled first: Escape backs a
  prompt out, `n` answers a declinable one "no", and `r` reaps with the
  selected creature (or takes the right flank while placing, or forges a Red
  key at that prompt). `r` is not a general "yes"; only Space's fallback
  (`affirm`) plays the selected hand card, uses an artifact's action, or
  answers a yes/no prompt.

## Persistence

The match is stored in local storage under `persistKey`, tagged with
`snapshotVersion`. Only the session's `Record` is written (version, seed, sets,
the ordered command log); state, the typed log, and the log bubbles are
replayed from it (`replayRecord`).

- **Bump `snapshotVersion` whenever an engine change makes an older record
  unreplayable**, so a stale snapshot is dropped rather than restored into a
  mismatched engine. That includes **rewording a log entry** (the log persists
  as narrated prose, so an old snapshot keeps the old wording) and any change
  to how a recorded command resolves. What a bump costs, and why log fixes are
  never batched, is on the `snapshotVersion` doc comment. `session.Version`
  guards the command format; `replayRecord` refuses a mismatch as
  `session.Load` does.
- Manual force-edits are in the record **on purpose**: a force-edited board
  that hits a bug is the reproduction worth keeping.
- **A snapshot that fails to _replay_ is kept, not deleted.** `resume` checks
  version and decode **before** replaying, so a later divergence or panic is a
  current-version save this build cannot replay: an engine regression and its
  only reproduction. Both paths (the `recover`, an error from `replayRecord`)
  call `quarantine`, which moves the snapshot to `quarantineKey`, raises
  `replayFailedNotice`, and clears the live slot for a fresh deal; one slot
  holds the latest finding. A wrong version or bad decode is still
  **deleted**. (`TestAPanickingReplayIsQuarantined`,
  `TestUnusableSnapshotsAreDropped`.)
- A failed write is not swallowed: `writeSnapshot` frees what it can spare (the
  style gallery's scroll memo, the snapshot being replaced), retries once, then
  raises a standing notice that the match is no longer saved.

## A replay failure is captured to disk

The browser cannot write a file, so on a replay failure `capture.go` POSTs the
record, panic, and stack to `POST /debug/capture`, and the **dev server** writes
it under `internal/web/testdata/capture/`. The endpoint exists only when
`web.DevEnabled()` (the switch `/style` uses); elsewhere the POST fails and the
quarantine key is the fallback.

- **A capture on disk is an open finding.** `TestCaptures` replays every one
  and fails, so a committed capture is a lapse: fix the fault, then
  `mage capturePrune`, the same contract as the `FuzzPlay` corpus.
- **A stale capture is skipped with a note, never failed.** It stamps the
  command-log version, the dev server's commit, and a **card-pool digest** (a
  hash over every implemented card's name, type, house, power, armor, and
  rendered text), because the same seed deals different cards from a different
  pool. The commit is for investigators and deliberately **not** a staleness
  key (ancestry is wrong across branches).
- **The filename is a content hash**, so re-hitting a bug overwrites its file.

## Two banners: transient status, standing notice

- `setStatus` is **transient** (fades after 5s): a rejected click the player
  can retry.
- `setNotice` is **standing** until `clearNotice`: a fault the player must act
  on (storage full, a set the deck generator does not know).
- The banner holds one notice, so a path clears only its own
  (`clearStorageNotice`).

## Tests drive the real client off-browser

`client_test.go` is the harness: it calls the same handlers a click or key
press does, over a real engine game. Add tests there rather than reimplementing
a slice of the client.

- `app.Context` is a **struct** and cannot be faked. go-app fires `OnMount`
  only when `app.IsClient` is true, which it is not in a host test, but fires
  `OnPreRender` when `app.IsServer` is true, which it is, so `ctxProbe`
  implements it purely to borrow a live context with working in-memory local
  storage and dispatch queue. `e.ConsumeAll()` drains dispatches, asyncs, and
  deferred work. The game is never waited for, so `c.settle`/`c.await` assert
  rather than wait.
- A zero `app.Event{}` nil-derefs on `PreventDefault()`; use `nullEvent()`.
- `app.HTMLString(g.Render())` draws the whole screen without mounting; that is
  how `view_test.go` asserts on markup.
- A test needing one exact prompt raises it with `c.script`, a session whose
  action is that single question, answered with a genuine `Command`. A test
  about backing OUT of a prompt uses a real play (`c.stagePromptArtifact`),
  since it needs a root action to rewind.
- **Never assert a fixed attribute order on drawn markup**: go-app ranges an
  attrs **map**, so order varies per render. Match one attribute
  (`src="/web/assets/shield.svg"`), or use a regex whose `[^>]*` skips
  attributes without fixing their order (`TestArmorShowsWhatIsLeftToAbsorb`).
- `app.Window()` reads back empty off-browser, so DOM-measuring code (the
  fly-into-play animation, log auto-scroll, picker focus, `ctx.JSSrc()`)
  no-ops. Keep every such piece tolerant of a render with no page behind it.
- Coverage here is **deliberately ungated** because of that DOM-bound code.

## Never shout: no all-caps

Nothing in the client is set in capitals: no `text-transform: uppercase`, no
`font-variant-caps: all-small-caps`, no string typed in caps in Go as a label.
Card types, row labels, zone names, and log headers render in the casing they
are written in; size, weight, and colour separate a label from prose. Drop the
wide `letter-spacing` that accompanies caps along with them.

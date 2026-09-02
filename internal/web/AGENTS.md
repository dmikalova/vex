# AGENTS.md — `internal/web`

The browser client: a [go-app](https://go-app.dev) v11 WASM front end over
`internal/engine`. Read the repo-root `AGENTS.md` first; this file covers only
what is specific to the web client.

Build it the same way as everything else: `mage build` compiles every package for
the host **and** the client for js/wasm, so a change that only breaks in the
browser fails the same gate as anything else. To compile just the client after an
edit, run `mage webWasm`.

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
- `game_action.go` — the action plumbing: undo snapshots, running an engine
  mutation off the UI goroutine, and the flash/flight bookkeeping after it.
- `game_play.go` — taking a turn: selection, house choice, play (click and drag),
  reap, use, fight, end turn.
- `game_chooser.go` — the `webChooser` bridge and the handlers that answer a
  prompt.
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
snapshot (the top of the undo stack) against the resolved state. Add a new
animation by diffing the state there, not by sprinkling calls at action sites.
The exception is information the state does not carry — a fight's two combatants
— which the handler arms on `g.fighters` for `computeFlashes` to consume.

A card that has left play cannot pulse, so it **flies** instead: `computeFlights`
finds the zone it landed in and `flightsInto` renders a ghost face parented to
that zone's pill, which arcs in and shrinks onto the count (`.card-flight`).

## Actions run off the UI goroutine; the chooser bridges back

An engine action can block on a player decision, so `runAction` resolves it on a
background goroutine and `webChooser` bridges the two: it posts prompt state via
`g.dispatch` (which runs on the UI goroutine) and blocks on a reply channel until
a card is clicked. Therefore:

- **Only touch `game` fields from the UI goroutine.** Inside `webChooser`, every
  mutation goes in a `g.dispatch(func(app.Context) { … })` closure — including any
  read of engine state used to set up the prompt.
- Reply channels are buffered and drained of stale values before each prompt, so
  a double click on the previous prompt cannot silently answer the next one.
- `g.busy` gates input while an action is in flight; most handlers begin with a
  `if g.busy || g.choosing … { return }` guard.

## Prompts: cards are clicked, options are buttons

- A card decision is a **card prompt**: the candidates highlight on the board and
  the controls become the prompt text. Mandatory ones (`ChooseCreature`) have no
  way out; optional ones (`ChooseCardOrDecline`, set by `chooserDeclinable`)
  get a **Done** button, and Escape declines them.
- A genuine yes/no or "choose one" stays an **option prompt** with buttons.
- When a prompt's candidates are not on the board (a discard pile — World Tree,
  Witch of the Eye), `openZoneForPrompt` opens that player's zone viewer, makes
  only the candidates clickable, dims the rest, and scrolls to the row. The viewer
  cannot be dismissed while it is the only place the prompt can be answered.

## Phases and Escape layering

`phase` is the client's own interaction state, not a rules concept:
`phaseHouse` → `phaseMain` → (`phaseFlank` | `phaseFightTarget`) → `phaseOver`.

Escape calls `dismiss`, which backs out **exactly one layer**, innermost first:
picker → zone viewer → restart confirmation → key-forge picker → an escapable
prompt → mid-action targeting → end-turn confirmation → the selection. Add a new
overlay to that chain rather than giving it its own Escape handling.

Keyboard shortcuts live in `onKey`. They are all single keys and all no-ops while
`busy`/`choosing`, except Escape and `r`, which are handled first: Escape backs a
prompt out and `r` answers it. `r` is the affirmative key (`affirm`) — yes to a
yes/no prompt, the right flank while placing, else the selected card's main use
(play, an artifact's action, otherwise reap).

## Persistence

The in-progress match is stored in local storage under `persistKey`, tagged with
`snapshotVersion`. **Bump `snapshotVersion` whenever an engine change makes an
older persisted state unloadable** — a stale snapshot is dropped rather than
restored into a mismatched engine.

## CSS conventions (`web/app.css`)

- BEM-ish: a block (`.card`), and modifiers as `--modifier` classes
  (`.card--dimmed`, `.log-line--p0`). No inline styles from Go.
- House colours are custom properties (`--nm`, `--tp`, `--edge`) supplied by the
  `.card-<house>` class from `palette.go`; markup only ever carries class names.
- Keep every animation's `-a`/`-b` pair in sync — they must have identical
  keyframes.

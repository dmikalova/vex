# The `Stepper` suspends an action

`Stepper` (`suspend.go`) is the interactive path of ADR 0040: it runs one
action on a goroutine with `suspendChooser` installed for both players; `Start`
yields the first `Request`, and `Advance(cmd)` answers the current one and
yields the next. `internal/session` drives it; the web client drives the
session and holds no `Chooser` (ADR 0047).

- **A panic inside the action does not escape.** The goroutine recovers
  **before** closing `requests`, never after: a panic unwinding past the close
  would leave every later `Start`/`Advance` blocked forever, and under wasm
  would take the program down (`TestStepperContainsAPanickingAction`). The
  recovered value and the stack become a `*PanicError`; capture the stack
  inside the deferred recover, or the panicking frames are gone.
- **`Err()` says why the action stopped.** `Start` and `Advance` report only
  that the action is **done**; the caller then asks `Err()` whether it finished
  or broke, the `bufio.Scanner` / `sql.Rows` iterator convention. Do not widen
  `Start`/`Advance` to return an error.
- `Close()` releases an action parked at a decision, so an abandoned `Stepper`
  (an undo, a replay that deals a fresh game) leaks neither the goroutine nor
  the `Game` it captured. It is idempotent
  (`TestStepperCloseReleasesGoroutine`).

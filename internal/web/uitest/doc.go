//go:build uitest

// Package uitest drives the browser scenarios headlessly.
//
// The scenarios themselves live in internal/web (uitest_scenarios.go) and run
// inside the page at /ui-test/<slug>. This package is ONLY the driver: it builds
// the client, serves it, points a headless Chrome at one page per registered
// scenario, and reads the one status element the page writes. It re-describes no
// journey and clicks nothing, so adding a scenario to the registry adds it to
// both the page a human watches and the suite `mage uiTest` runs, with no second
// edit here.
//
// Every file is behind the `uitest` build tag, which keeps the package out of
// `./...`: `mage ci:check` never builds a multi-megabyte wasm bundle or boots a
// browser, and the shared project-standards CI workflow needs neither. Run it
// with `mage uiTest`.
package uitest

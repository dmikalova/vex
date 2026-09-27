package web

import "github.com/maxence-charriere/go-app/v11/pkg/app"

// DevEnv is the variable that turns the client's development-only surfaces on:
// the Style gallery's and cluster map's routes, and the capture endpoint a replay
// failure posts to. `mage web` sets it, so they are there whenever the client is
// being developed and absent from every other deployment.
const DevEnv = "VEX_STYLE"

// DevEnabled reports whether the development-only surfaces should exist. The
// check has to run on both sides of the build and agree: the server must register
// a route or it serves no page at all (go-app 404s an unregistered path), and the
// wasm client must register it or the served page renders nothing. app.Getenv
// bridges the two — it reads the process environment on the server and the Env map
// the server passed down on the client — so one variable decides both.
//
// It is an environment switch rather than a build tag because go-app routes on
// the client: a tag would have to exclude the gallery from the wasm bundle every
// player downloads, which means the gallery would not be compiled by default and
// would rot exactly as the //go:build todo card stubs do. See
// docs/adr/0014-style-gallery-on-real-components.md.
func DevEnabled() bool { return app.Getenv(DevEnv) == "1" }

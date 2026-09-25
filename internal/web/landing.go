package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is the landing page served at "/". It is a static splash that
// introduces Vex and links into the game (/play) and the reference pages.
// The interactive client now lives at /play; "/" is a lightweight entry point
// that loads without spinning up a match.

// NewLanding returns the root component for the "/" landing page.
func NewLanding() app.Composer { return &landing{} }

// landing is the "/" landing page component. It holds no state.
type landing struct {
	app.Compo
}

// OnAppUpdate reloads the page onto a freshly built wasm bundle, matching the
// hot-reload hand-off the other pages use.
func (l *landing) OnAppUpdate(ctx app.Context) { ctx.Reload() }

// Render draws the splash: brand, tagline, a primary Play call to action, and
// links to the reference pages.
func (l *landing) Render() app.UI {
	return app.Div().Class("landing").Body(
		app.Div().Class("landing-hero").Body(
			app.H1().Class("landing-title").Text("Vex"),
			app.P().Class("landing-tagline").
				Text("A KeyForge-style card game, playable in your browser."),
			app.A().Class("landing-play").Href("/play").Text("Play"),
			app.Nav().Class("landing-nav").Body(
				app.A().Class("landing-link").Href("/cards").Text("Cards"),
				app.A().Class("landing-link").Href("/rulebook").Text("Rulebook"),
				app.A().Class("landing-link").Href("/glossary").Text("Glossary"),
			),
		),
	)
}

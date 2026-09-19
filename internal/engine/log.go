package engine

import (
	"fmt"
	"strings"
)

// This file holds the game log's spine (ADR 0011): the LogEntry interface every
// narrated outcome implements, the attribution frame entries inherit, and the
// recorder on Game. The entry variants themselves live one verb family per
// log_<family>.go file.
//
// The log narrates RESOLVED OUTCOMES, not card text. A card's Text() renders an
// unbound, present-tense imperative ("Deal 2 damage to a creature"); a log line
// states a bound, past-tense outcome ("Troll takes 2 damage and now has 2
// damage"), which is only knowable after resolution. The two renderers share a
// vocabulary but neither derives from the other.
//
// Four rules govern how a new entry is worded. They are the log's own voice, and
// they are where the log deliberately departs from card text:
//
//   - MECHANICS AND KEYWORDS STAY LOWERCASE. A log line is a sentence, not a
//     title, so it reads "Card2 gains skirmish" and "uses Card2's action
//     ability". Card text capitalizes a keyword; the log does not, which also
//     keeps it clear of the casing drift card text has to police.
//   - THE SOURCE CARD COMES FIRST. "Nexus has P0 gain 2 Æmber" puts the causing
//     card at the head of the sentence rather than burying it in a trailing
//     clause, so a reader scanning the left edge sees what acted.
//   - MINIMIZE LEFT-TO-RIGHT BACKTRACKING. Prefer the phrasing where no noun has
//     to be re-resolved once read. "X exalts 2 Æmber onto Y" reads straight
//     through; "X exalts Y with 2 Æmber" makes the reader go back to Y.
//   - NO PARENTHETICALS. A running total is spelled out ("takes 3 damage and now
//     has 4 damage", "forges a Red key and now has 1 of 3 keys"), and a suffix
//     that qualified the whole line is folded into the verb ("manually sets their
//     pool to 7 Æmber"). Suppress the tail when it would restate the amount, so a
//     first gain reads "gains 2 chains" rather than "gains 2 chains and now has 2
//     chains". FightKeywords.annotate is the one exemption, argued there.
//
// A rendered entry is persisted prose: the web client stores the log it has
// already rendered, so changing what an entry renders — or adding and removing
// entries — requires bumping snapshotVersion in internal/web/game.go.

// Namer resolves the ids an entry carries to the names its reader sees. Entries
// hold LocalIDs and player numbers rather than baked strings, so a client can
// turn a named card into a link and an Æmber mention into an icon without
// pattern-matching prose.
type Namer interface {
	// Name is the display name of a card.
	Name(id LocalID) string
	// PlayerName is the display name of a player.
	PlayerName(player int) string
}

// LogEntry is one narrated outcome. Each variant renders itself, the same
// Interpreter shape as Effect. Extrinsic passes over entries — the client's icon
// and card-link rendering, a replay analyzer — are standalone type-switch
// functions over LogEntry, never another method on every variant (ADR 0006).
type LogEntry interface {
	// Text renders the outcome as a past-tense sentence.
	Text(n Namer) string
}

// nameMoved names a card that moved from one zone to another. It is the single
// place the log decides whether a card may be named: it is named once either end
// of the move is public, and is otherwise just "a card", so no entry can leak a
// hand or a deck on its own initiative (ADR 0011). A card in a hidden zone is
// named only by a Reveal, which records its own entry.
func nameMoved(n Namer, id LocalID, from, to Zone) string {
	if from.public() || to.public() {
		return n.Name(id)
	}
	return "a card"
}

// namedCards lists several cards by name, for an entry that narrates a group.
func namedCards(n Namer, ids []LocalID) string {
	names := make([]string, len(ids))
	for i, id := range ids {
		names[i] = n.Name(id)
	}
	return strings.Join(names, ", ")
}

// namedCardsAnd lists several cards by name as a prose series — "A", "A and B",
// or "A, B, and C" — for an entry that reads as a sentence.
func namedCardsAnd(n Namer, ids []LocalID) string {
	names := make([]string, len(ids))
	for i, id := range ids {
		names[i] = n.Name(id)
	}
	if len(names) < 2 {
		return strings.Join(names, "")
	}
	if len(names) == 2 {
		return names[0] + " and " + names[1]
	}
	return strings.Join(names[:len(names)-1], ", ") + ", and " + names[len(names)-1]
}

// because appends the event a lasting reaction fired on, so a lasting entry
// reads as its plain outcome plus what triggered it and the two never drift.
func because(text string, on Event) string {
	return text + " (" + on.clause() + ")"
}

// Frame is the attribution a run of entries shares: who acted, which card they
// acted with, under which of that card's ability categories, and which card
// granted that ability when it was not the source's own text. Attribution is not
// passed to each emission site; resolution opens a frame and the entries emitted
// inside it inherit it.
type Frame struct {
	// Actor is the player whose action or ability is resolving.
	Actor int
	// Source is the card being played or used, when there is one.
	Source LocalID
	// HasSource distinguishes "no source card" from LocalID 0, which is a real card.
	HasSource bool
	// Trigger is the ability category resolving (Play:, Reap:, Destroyed:, …), or
	// the unset trigger for an outcome that is not an ability.
	Trigger Trigger
	// Grantor is the card whose text granted the ability, when an upgrade or a
	// constant ability granted it rather than the source carrying it.
	Grantor LocalID
	// HasGrantor distinguishes "the source's own text" from LocalID 0.
	HasGrantor bool
	// Depth is how many frames are open above this one. A client groups the entries
	// under one depth-0 frame into a single bubble.
	Depth int
}

// Record is one entry together with the attribution it was emitted under.
type Record struct {
	Frame Frame
	Entry LogEntry
}

// Text renders the record's outcome under its frame, so a card ability's outcome
// reads with the source card as its subject ("Batdrone deals 2 damage") rather
// than the acting player (ADR 0011).
func (r Record) Text(n Namer) string {
	return r.Entry.Text(framedNamer{Namer: n, frame: r.Frame})
}

// sourced is a Namer that knows the source card of the frame an entry is
// rendering under. An outcome entry asks for its subject through subject, which
// resolves to that card when a card ability is resolving.
type sourced interface {
	frameSource() (LocalID, bool)
}

// framedNamer names an outcome under the frame it was recorded in. It resolves
// subject to the frame's source card, so an ability's line reads "Batdrone deals
// 2 damage" instead of "Player 0 deals 2 damage"; an outcome recorded outside a
// card frame — a player's own play, a turn event, a lasting effect — has no
// source and falls back to the player.
type framedNamer struct {
	Namer
	frame Frame
}

// frameSource reports the source card of the frame, if any.
func (f framedNamer) frameSource() (LocalID, bool) {
	return f.frame.Source, f.frame.HasSource
}

// framedSource names the source card of the frame an outcome renders under, when
// a card ability is behind it. An entry whose affected player is never its own
// actor — a pool a card drains, a hand a card forces discarded — uses this to
// read "Card has Player 1 …" under an ability and "Player 1 …" on its own.
func framedSource(n Namer) (string, bool) {
	if s, ok := n.(sourced); ok {
		if id, has := s.frameSource(); has {
			return n.Name(id), true
		}
	}
	return "", false
}

// subject renders the subject of an outcome: the frame's source card when a card
// ability is resolving, and the acting player otherwise. An outcome entry names
// its subject through this so its line reads "Card does X" inside a bubble whose
// header already named the player who acted.
func subject(n Namer, player int) string {
	if s, ok := framedSource(n); ok {
		return s
	}
	return n.PlayerName(player)
}

// actorPossessive names the party acting on a zone and the possessive determiner
// for zones that party owns: under a card's ability the source card acts and the
// owner is named ("Dew Faerie … from Player 2's discard pile"), while on a
// player's own action the player acts and owns ("Player 2 … from their discard
// pile"). A zone move names the owner outright rather than "their" so a card
// acting on its own or an opponent's zone both read right.
func actorPossessive(n Namer, player int) (who, owner string) {
	if s, ok := framedSource(n); ok {
		return s, n.PlayerName(player) + "'s"
	}
	return n.PlayerName(player), "their"
}

// openFrame pushes an attribution frame; every entry recorded until the returned
// function runs inherits it. Frames nest, so an ability that causes another to
// resolve keeps its own attribution on its own entries.
func (g *Game) openFrame(f Frame) func() {
	f.Depth = len(g.frames)
	g.frames = append(g.frames, f)
	return func() { g.frames = g.frames[:len(g.frames)-1] }
}

// frame returns the attribution now in force, or the bare active-player frame
// when nothing has opened one.
func (g *Game) frame() Frame {
	if len(g.frames) == 0 {
		return Frame{Actor: g.State.ActivePlayer}
	}
	return g.frames[len(g.frames)-1]
}

// record narrates one resolved outcome under the frame in force. Recording is
// switchable off so a simulated state — the MCTS bot cloning a position — costs
// no allocation; log behavior can therefore never be load-bearing for game rules.
func (g *Game) record(e LogEntry) {
	if !g.recording {
		return
	}
	g.Log = append(g.Log, Record{Frame: g.frame(), Entry: e})
	if g.Verbose {
		fmt.Println(e.Text(g))
	}
}

// LogText renders the whole log as prose, for a reader that wants lines rather
// than entries: a terminal dump, a test assertion, a soak failure report.
func (g *Game) LogText() []string {
	out := make([]string, len(g.Log))
	for i, rec := range g.Log {
		out[i] = rec.Text(g)
	}
	return out
}

// SetRecording turns log recording on or off. It is on by default; a bot that
// only wants the rules turns it off.
func (g *Game) SetRecording(on bool) { g.recording = on }

// RestoredEntry is an entry read back from a persisted log, whose text was
// narrated by an earlier run of the game. A typed entry does not survive JSON
// (the same reason a CardDefinition does not), so a frontend that persists the
// log keeps each entry's rendered text and restores it as one of these. It reads
// the same but carries no ids, so a restored line renders without card links.
type RestoredEntry struct{ Line string }

// Text returns the text the entry was narrated with when it was recorded.
func (e RestoredEntry) Text(Namer) string { return e.Line }

// Restore replaces the log with entries read back from persistence.
func (g *Game) Restore(log []Record) { g.Log = log }

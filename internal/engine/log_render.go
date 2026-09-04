package engine

import "strings"

// This file holds the extrinsic pass that turns a rendered entry back into its
// parts, so a client can link the cards an entry names and tint the players it
// names without matching prose against a card index (ADR 0011). It is a
// standalone function over LogEntry, not a second method on every variant
// (ADR 0006).

// LogSegment is one piece of a rendered log entry: plain text, the name of a
// card the entry named, the name of a player it named, or a keyword that stands
// for a game concept with an emblem — each carrying what it was named by, so a
// client never has to recognise it from the words.
type LogSegment struct {
	Text string
	Card LocalID
	// HasCard distinguishes a card named by LocalID 0 from plain text.
	HasCard bool
	Player  int
	// HasPlayer distinguishes player 0 from plain text.
	HasPlayer bool
	// Icon is the concept key of the emblem this segment's keyword stands for
	// ("aember", "house-brobnar", …), or "" when the segment is not a keyword.
	Icon string
}

// RenderEntry renders an entry and splits the result into segments, marking the
// spans that are names or keywords. What each name stands for comes from the
// entry itself: rendering is watched, so every name the entry asked for is known
// together with the id or player it asked by, and no client has to guess which
// words in a sentence are cards, players, or emblems.
func RenderEntry(e LogEntry, n Namer) []LogSegment {
	spy := &namerSpy{Namer: n}
	text := e.Text(spy)
	var out []LogSegment
	plain := 0
	for i := 0; i < len(text); {
		marked, ok := segmentAt(spy, text, i)
		if !ok {
			i++
			continue
		}
		if i > plain {
			out = append(out, LogSegment{Text: text[plain:i]})
		}
		out = append(out, marked)
		i += len(marked.Text)
		plain = i
	}
	if plain < len(text) {
		out = append(out, LogSegment{Text: text[plain:]})
	}
	return out
}

// segmentAt returns the marked segment starting at text[i], if any. A name the
// entry asked for beats an icon keyword, so a card called "Æmber Imp" stays one
// card link instead of splitting into an emblem and prose.
func segmentAt(spy *namerSpy, text string, i int) (LogSegment, bool) {
	if named, ok := spy.nameAt(text, i); ok {
		return named.segment(), true
	}
	return iconAt(text, i)
}

// aemberWord is the keyword the Æmber emblem stands in front of.
const aemberWord = "Æmber"

// iconAt returns the icon segment for a keyword starting at text[i]. The icon
// vocabulary is closed and owned by the engine — Æmber, the seven houses, being
// stunned, chains, and keys — so a client draws emblems from what the engine
// reported rather than scanning prose for words it happens to recognise
// (ADR 0011).
func iconAt(text string, i int) (LogSegment, bool) {
	if wordAt(text, i, aemberWord) {
		return LogSegment{Text: aemberWord, Icon: "aember"}, true
	}
	for h := Brobnar; h <= Untamed; h++ {
		if name := h.String(); wordAt(text, i, name) {
			return LogSegment{Text: name, Icon: houseIconKey(h)}, true
		}
	}
	if wordAt(text, i, "stunned") {
		return LogSegment{Text: "stunned", Icon: "stun"}, true
	}
	if wordAt(text, i, "chains") {
		return LogSegment{Text: "chains", Icon: "chains"}, true
	}
	if wordAt(text, i, "chain") {
		return LogSegment{Text: "chain", Icon: "chains"}, true
	}
	if wordAt(text, i, "keys") {
		return LogSegment{Text: "keys", Icon: "key"}, true
	}
	// "key phase" names the turn phase, not an actual key, so it stays plain text.
	if wordAt(text, i, "key") && !strings.HasPrefix(text[i+len("key"):], " phase") {
		return LogSegment{Text: "key", Icon: "key"}, true
	}
	return LogSegment{}, false
}

// houseIconKey is the concept key of a house's emblem.
func houseIconKey(h House) string {
	return "house-" + strings.ToLower(h.String())
}

// namedThing is one card or player an entry asked to have named while it
// rendered, remembered together with what it was named by.
type namedThing struct {
	name   string
	card   LocalID
	player int
	isCard bool
}

// segment turns a remembered name into the segment that stands for it.
func (t namedThing) segment() LogSegment {
	if t.isCard {
		return LogSegment{Text: t.name, Card: t.card, HasCard: true}
	}
	return LogSegment{Text: t.name, Player: t.player, HasPlayer: true}
}

// namerSpy is a Namer that remembers what it was asked to name, so the rendered
// text can be split back into the ids that produced it.
type namerSpy struct {
	Namer
	named []namedThing
}

func (s *namerSpy) Name(id LocalID) string {
	name := s.Namer.Name(id)
	s.named = append(s.named, namedThing{name: name, card: id, isCard: true})
	return name
}

func (s *namerSpy) PlayerName(player int) string {
	name := s.Namer.PlayerName(player)
	s.named = append(s.named, namedThing{name: name, player: player})
	return name
}

// nameAt returns the longest name that starts at text[i] on a word boundary,
// among the names this entry actually used.
func (s *namerSpy) nameAt(text string, i int) (namedThing, bool) {
	var best namedThing
	for _, c := range s.named {
		if len(c.name) > len(best.name) && wordAt(text, i, c.name) {
			best = c
		}
	}
	return best, best.name != ""
}

// wordAt reports whether word sits at text[i] as a whole word rather than as
// part of a longer one.
func wordAt(text string, i int, word string) bool {
	if word == "" || !strings.HasPrefix(text[i:], word) {
		return false
	}
	if i > 0 && isNameByte(text[i-1]) {
		return false
	}
	end := i + len(word)
	return end >= len(text) || !isNameByte(text[end])
}

// isNameByte reports whether b can sit inside a word, which is what makes a
// name match a whole word rather than part of one.
func isNameByte(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}

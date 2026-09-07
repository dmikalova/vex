package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/dmikalova/vactrol/internal/cards/provenance"
)

// importSet pairs a source set with its Master Vault expansion id. There is no
// live expansions endpoint (the decks feed links only houses), so the ids are
// pinned here from decksofkeyforge.com's expansion table. Order is release order,
// the order `import-provenance all` walks the sets.
type importSet struct {
	Set       provenance.SourceSet
	Expansion int
}

// importSets is every set the importer can build a catalog for, with its Master
// Vault expansion id.
var importSets = []importSet{
	{provenance.CallOfTheArchons, 341},
	{provenance.AgeOfAscension, 435},
	{provenance.WorldsCollide, 452},
	{provenance.MassMutation, 479},
	{provenance.DarkTidings, 496},
	{provenance.WindsOfExchange, 600},
	{provenance.GrimReminders, 700},
	{provenance.AemberSkies, 800},
	{provenance.TokensOfChange, 855},
	{provenance.MoreMutation, 874},
	{provenance.Menagerie, 722},
	{provenance.VaultMasters2025, 939},
	{provenance.PropheticVisions, 886},
	{provenance.CrucibleClash, 918},
	{provenance.DraconianMeasures, 928},
}

// importSetBySlug returns the importSet for a slug, or false when unknown.
func importSetBySlug(slug string) (importSet, bool) {
	for _, is := range importSets {
		if is.Set.Slug == slug {
			return is, true
		}
	}
	return importSet{}, false
}

// importSlugs lists every importable set slug, for error messages.
func importSlugs() string {
	slugs := make([]string, len(importSets))
	for i, is := range importSets {
		slugs[i] = is.Set.Slug
	}
	return strings.Join(slugs, ", ")
}

// importProvenance builds a source catalog (internal/cards/provenance/<slug>.json)
// by paging the Master Vault decks feed for a set and folding each card it links
// into the catalog shape. With no slug it opens an interactive picker offering
// every set plus "All sets"; the slug "all" imports every set in release order.
// Each set is fetched on its own — there is no combined query.
func importProvenance(args []string) error {
	var slug string
	switch len(args) {
	case 0:
	case 1:
		slug = args[0]
	default:
		return fmt.Errorf("usage: cardlookup import-provenance [setSlug|all]")
	}
	if slug == "" {
		if isInteractive() {
			picked, err := pickImportSet()
			if err != nil {
				return err
			}
			if picked == "" {
				return nil // cancelled
			}
			slug = picked
		} else {
			return fmt.Errorf(
				"import-provenance needs a set slug or `all` when stdin is not a terminal",
			)
		}
	}

	targets := importSets
	if slug != "all" {
		is, ok := importSetBySlug(slug)
		if !ok {
			return fmt.Errorf("unknown set slug %q; known slugs: %s, all", slug, importSlugs())
		}
		targets = []importSet{is}
	}

	// A fresh connection per request (no keep-alive) mirrors curl: reused pooled
	// connections go half-dead under the feed's aggressive throttling and then
	// hang the next request until the client timeout.
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			TLSHandshakeTimeout: 15 * time.Second,
		},
	}
	for _, is := range targets {
		if err := importOneSet(client, is); err != nil {
			return fmt.Errorf("importing %s: %w", is.Set.Name, err)
		}
	}
	return nil
}

// importOneSet fetches, transforms, and writes one set's catalog.
func importOneSet(client *http.Client, is importSet) error {
	fmt.Printf("%s (expansion %d): fetching…\n", is.Set.Name, is.Expansion)
	raw, err := fetchSetCards(client, is.Set.Slug, is.Expansion)
	if err != nil {
		return err
	}
	cards := make([]catalogCard, 0, len(raw))
	for _, r := range raw {
		c, err := transformCard(r)
		if err != nil {
			return err
		}
		cards = append(cards, c)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].Number < cards[j].Number })

	path := filepath.Join("internal", "cards", "provenance", is.Set.Slug+".json")
	if err := writeCatalog(path, cards); err != nil {
		return err
	}
	fmt.Printf("%s: wrote %d cards to %s\n", is.Set.Name, len(cards), path)
	return nil
}

// --- Master Vault decks feed -------------------------------------------------

const mvDecksURL = "https://www.keyforgegame.com/api/decks/"

// mvResponse is the slice of the decks feed the importer reads: the total deck
// count and the card objects the page's decks link to (links=cards folds them
// into _linked.cards).
type mvResponse struct {
	Count  int `json:"count"`
	Linked struct {
		Cards []mvCard `json:"cards"`
	} `json:"_linked"`
}

// mvCard is the subset of a Master Vault card object the catalog needs. Numeric
// stats arrive as strings on some cards and as JSON numbers on others, so they
// are decoded through flexString.
type mvCard struct {
	CardTitle  string     `json:"card_title"`
	House      string     `json:"house"`
	CardType   string     `json:"card_type"`
	CardText   string     `json:"card_text"`
	Traits     string     `json:"traits"`
	Amber      flexString `json:"amber"`
	Power      flexString `json:"power"`
	Armor      flexString `json:"armor"`
	Rarity     string     `json:"rarity"`
	CardNumber string     `json:"card_number"`
	Expansion  int        `json:"expansion"`
	IsMaverick bool       `json:"is_maverick"`
	IsAnomaly  bool       `json:"is_anomaly"`
}

// flexString decodes a JSON value that may be either a string or a number into a
// string, since the decks feed spells the same stat both ways across cards.
type flexString string

func (f *flexString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = flexString(s)
		return nil
	}
	*f = flexString(strings.Trim(string(data), `"`))
	return nil
}

// fetchSetCards pages the decks feed for one expansion, collecting every distinct
// card it links, until the set is covered. It stops when the decks run out, or a
// run of pages adds nothing new: a run of 10 empty pages with no interior gap in
// the collected collector numbers (the set looks complete), or a hard cap of 100
// empty pages (a lingering gap that never fills). Maverick cards are skipped; a
// card belonging to another expansion (a cross-set maverick) is skipped too.
func fetchSetCards(client *http.Client, slug string, expansion int) ([]mvCard, error) {
	const pageSize = 25
	seen := map[string]bool{}
	var collected []mvCard

	// The feed throttles heavily, so pace requests: start at a polite delay and
	// slow down whenever a page gets throttled, easing back toward the baseline on
	// clean pages. A 429 is never fatal — fetchDecksPage waits it out.
	const basePace = time.Second
	const maxPace = 20 * time.Second
	pace := basePace

	// longest is the longest throttle wait seen so far, kept across pages so a
	// fresh throttle resumes from that floor instead of restarting at two seconds.
	longest := 2 * time.Second

	emptyStreak := 0
	for page := 1; ; page++ {
		resp, throttled, err := fetchDecksPage(client, &longest, expansion, page, pageSize)
		if err != nil {
			return nil, err
		}
		if throttled {
			pace += 3 * time.Second
			if pace > maxPace {
				pace = maxPace
			}
		} else if pace > basePace {
			pace -= time.Second
		}
		added := 0
		for _, c := range resp.Linked.Cards {
			if c.IsMaverick || (c.Expansion != 0 && c.Expansion != expansion) {
				continue
			}
			key := c.CardNumber + "|" + c.CardTitle
			if seen[key] {
				continue
			}
			seen[key] = true
			collected = append(collected, c)
			added++
		}
		if added > 0 {
			emptyStreak = 0
		} else {
			emptyStreak++
		}

		decksSeen := page * pageSize
		if decksSeen > resp.Count {
			decksSeen = resp.Count
		}
		fmt.Printf("%s: %d\n", slug, len(collected))
		fmt.Printf(
			"  page %d: %d/%d decks scanned, %d cards so far\n",
			page, decksSeen, resp.Count, len(collected),
		)

		// Out of decks: the last page came back short.
		if page*pageSize >= resp.Count {
			break
		}
		if emptyStreak >= 100 {
			break
		}
		if emptyStreak >= 10 && !hasInteriorGap(collected) {
			break
		}
		time.Sleep(pace)
	}
	return collected, nil
}

// fetchDecksPage requests one page of the decks feed. The feed throttles heavily,
// so a 429 is never fatal: it waits and keeps retrying until the page comes back.
// Each wait starts from *longest — the longest wait seen so far across all pages —
// and doubles up to a one-minute cap, so throttling that returns picks up where the
// backoff left off rather than restarting at two seconds. *longest is updated to
// the longest wait used. It reports whether the page was throttled so the caller
// can slow its pace.
func fetchDecksPage(
	client *http.Client,
	longest *time.Duration,
	expansion, page, pageSize int,
) (*mvResponse, bool, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	q.Set("links", "cards")
	q.Set("expansion", strconv.Itoa(expansion))
	q.Set("ordering", "date")
	reqURL := mvDecksURL + "?" + q.Encode()

	throttled := false
	const maxBackoff = 60 * time.Second
	wait := *longest
	sleepBackoff := func() {
		if wait > *longest {
			*longest = wait
		}
		time.Sleep(wait)
		wait *= 2
		if wait > maxBackoff {
			wait = maxBackoff
		}
	}
	for {
		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, throttled, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "vactrol-cardlookup/1.0")
		resp, err := client.Do(req)
		if err != nil {
			// A timeout or reset under load is transient — wait it out too.
			throttled = true
			fmt.Printf(
				"  page %d request failed (%v), waiting %s then retrying…\n",
				page, err, wait.Round(time.Second),
			)
			sleepBackoff()
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			_ = resp.Body.Close()
			throttled = true
			if ra := retryAfter(resp); ra > wait {
				wait = ra
			}
			fmt.Printf(
				"  page %d throttled (429), waiting %s then retrying…\n",
				page, wait.Round(time.Second),
			)
			sleepBackoff()
			continue
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
			_ = resp.Body.Close()
			return nil, throttled, fmt.Errorf(
				"decks feed status %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}
		var out mvResponse
		err = json.NewDecoder(resp.Body).Decode(&out)
		_ = resp.Body.Close()
		if err != nil {
			return nil, throttled, fmt.Errorf("decoding decks feed: %w", err)
		}
		return &out, throttled, nil
	}
}

// retryAfter reads the feed's Retry-After header (seconds), doubling it as the
// keyteki converter does, and falls back to two seconds when it is absent.
func retryAfter(resp *http.Response) time.Duration {
	if v := resp.Header.Get("Retry-After"); v != "" {
		if secs, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return time.Duration(secs*2) * time.Second
		}
	}
	return 2 * time.Second
}

// hasInteriorGap reports whether the collected numeric collector numbers skip a
// value below their maximum — a hole that means a card is still missing, so the
// importer keeps paging. Lettered reference numbers (S01, A21, P07) are ignored;
// they fill as they appear.
func hasInteriorGap(cards []mvCard) bool {
	maxNum := 0
	present := map[int]bool{}
	for _, c := range cards {
		n, err := strconv.Atoi(strings.TrimLeft(c.CardNumber, "0"))
		if err != nil {
			continue
		}
		present[n] = true
		if n > maxNum {
			maxNum = n
		}
	}
	for n := 1; n < maxNum; n++ {
		if !present[n] {
			return true
		}
	}
	return false
}

// --- transform ---------------------------------------------------------------

// catalogCard is one card as written to a source catalog. Field order is the
// order the JSON is emitted; empty optional fields are omitted. Printed carries
// the original card title only when it differs from the ASCII-folded name.
type catalogCard struct {
	Number  string   `json:"number"`
	Name    string   `json:"name"`
	Printed string   `json:"printed,omitempty"`
	House   string   `json:"house"`
	Type    string   `json:"type"`
	Rarity  string   `json:"rarity"`
	Traits  []string `json:"traits,omitempty"`
	Power   int      `json:"power,omitempty"`
	Armor   int      `json:"armor,omitempty"`
	Amber   int      `json:"amber,omitempty"`
	Text    string   `json:"text,omitempty"`
}

// transformCard folds one Master Vault card object into the catalog shape: it
// ASCII-folds the title and text, expands the amber/damage markup, normalizes the
// house and rarity, and classifies an anomaly. It fails loudly on an unknown
// markup token so a new icon is caught rather than silently dropped.
func transformCard(m mvCard) (catalogCard, error) {
	name := asciiFold(m.CardTitle)
	text, err := expandAndFoldText(m.CardText)
	if err != nil {
		return catalogCard{}, fmt.Errorf("card %q: %w", m.CardTitle, err)
	}

	house := strings.ToLower(strings.ReplaceAll(m.House, " ", ""))
	rarity := m.Rarity
	if strings.EqualFold(rarity, "FIXED") {
		rarity = "Special"
	}
	if m.IsAnomaly {
		house = "brobnar"
		rarity = "Special"
	}

	c := catalogCard{
		Number: normalizeCatalogNumber(m.CardNumber),
		Name:   name,
		House:  house,
		Type:   strings.ToLower(m.CardType),
		Rarity: rarity,
		Traits: foldTraits(m.Traits),
		Amber:  atoiSafe(string(m.Amber)),
		Text:   text,
	}
	if name != m.CardTitle {
		c.Printed = m.CardTitle
	}
	if strings.EqualFold(m.CardType, "creature") {
		c.Power = atoiSafe(string(m.Power))
		c.Armor = atoiSafe(string(m.Armor))
	}
	return c, nil
}

// normalizeCatalogNumber zero-pads an all-digit collector number to three digits
// (the printed form, "4" -> "004") and leaves a lettered reference number as is
// ("S01", "A21", "P07").
func normalizeCatalogNumber(s string) string {
	s = strings.TrimSpace(s)
	if n, err := strconv.Atoi(s); err == nil {
		return fmt.Sprintf("%03d", n)
	}
	return s
}

// foldTraits splits a Master Vault traits string on its bullet separator, folds
// each to lowercase ASCII, and drops the empties. A card with no traits yields a
// nil slice, which is omitted from the JSON.
func foldTraits(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, " • ")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToLower(strings.TrimSpace(asciiFold(p)))
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// atoiSafe parses a Master Vault numeric string, treating blank or non-numeric as
// zero (the stat is absent).
func atoiSafe(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// markupWords maps a Master Vault card-text icon token to the word the catalog
// spells it as. Only the tokens the feed is known to use are listed; expandMarkup
// fails on any other so a new icon surfaces instead of vanishing.
var markupWords = map[string]string{
	"A": "Aember",
	"D": "Damage",
}

// expandAndFoldText expands a card's markup, ASCII-folds the result, strips the
// templated keyword reminder text, and tidies stray whitespace.
func expandAndFoldText(s string) (string, error) {
	expanded, err := expandMarkup(s)
	if err != nil {
		return "", err
	}
	return tidyText(stripKeywordReminders(asciiFold(expanded))), nil
}

// keywordReminderTexts is the exact parenthetical reminder text KeyForge prints
// after a keyword. Only these templated blurbs are stripped; any other
// parenthetical (a card-specific clarification) is left untouched.
var keywordReminderTexts = []string{
	"The first time this creature is attacked each turn, no damage is dealt.",
	"When you use this creature to fight, it is dealt no damage in return.",
	"When you use a creature with skirmish to fight, it is dealt no damage in return.",
	"This creature's neighbors cannot be attacked unless they have taunt.",
	"Any damage dealt by this creature's power during a fight destroys the damaged creature.",
}

// keywordReminderPatterns matches the reminder texts that carry a number
// (Assault N, Hazardous N), anchored so only the whole blurb qualifies.
var keywordReminderPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^Before this creature attacks, deal \d+ Damage to the attacked enemy\.$`),
	regexp.MustCompile(
		`^Before this creature is attacked, deal \d+ Damage to the attacking enemy\.$`,
	),
}

// parenReminderRe captures a parenthetical along with any spaces or tabs before
// it, so removing a keyword reminder leaves no dangling separator.
var parenReminderRe = regexp.MustCompile(`[ \t]*\([^)]*\)`)

// stripKeywordReminders removes each parenthetical whose text is a templated
// keyword reminder, leaving card-specific parentheticals in place.
func stripKeywordReminders(s string) string {
	return parenReminderRe.ReplaceAllStringFunc(s, func(m string) string {
		open := strings.IndexByte(m, '(')
		inner := strings.TrimSpace(m[open+1 : len(m)-1])
		if isKeywordReminder(inner) {
			return ""
		}
		return m
	})
}

// isKeywordReminder reports whether inner is one of the templated keyword blurbs.
func isKeywordReminder(inner string) bool {
	for _, t := range keywordReminderTexts {
		if inner == t {
			return true
		}
	}
	for _, re := range keywordReminderPatterns {
		if re.MatchString(inner) {
			return true
		}
	}
	return false
}

// tidyText drops trailing whitespace from each line (including space left before a
// newline) and trims the leading and trailing whitespace of the whole text.
func tidyText(s string) string {
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		lines[i] = strings.TrimRight(ln, " \t")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// expandMarkup turns Master Vault card-text markup into plain words: the vertical
// tab line break becomes a newline, and each icon token (<A>, <D>) becomes its
// word, with a space inserted where the icon abutted a letter or digit so "2<A>"
// reads "2 Aember". An unrecognized or unterminated token is an error.
func expandMarkup(s string) (string, error) {
	s = strings.ReplaceAll(s, "\u000b", "\n")
	var b strings.Builder
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if rs[i] != '<' {
			b.WriteRune(rs[i])
			continue
		}
		j := i + 1
		for j < len(rs) && rs[j] != '>' {
			j++
		}
		if j >= len(rs) {
			return "", fmt.Errorf("unterminated markup token in %q", s)
		}
		tok := string(rs[i+1 : j])
		word, ok := markupWords[tok]
		if !ok {
			return "", fmt.Errorf("unknown markup token <%s> in %q", tok, s)
		}
		if out := b.String(); out != "" && isAlnumASCII(rune(out[len(out)-1])) {
			b.WriteByte(' ')
		}
		b.WriteString(word)
		if j+1 < len(rs) && isAlnumASCII(rs[j+1]) {
			b.WriteByte(' ')
		}
		i = j
	}
	return b.String(), nil
}

// isAlnumASCII reports whether r is a plain ASCII letter or digit.
func isAlnumASCII(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
}

// asciiFold rewrites a string into plain ASCII: curly quotes straighten, Æ becomes
// Ae, the multiplication sign becomes x, en/em dashes and the minus sign become a
// hyphen, and accented letters lose their diacritics (é -> e). English card text
// and names fold cleanly; anything left non-ASCII is passed through unchanged.
func asciiFold(s string) string {
	s = strings.NewReplacer(
		"\u2019", "'", // ’ right single quote
		"\u2018", "'", // ‘ left single quote
		"\u02bc", "'", // ʼ modifier letter apostrophe
		"\u201c", `"`, // “ left double quote
		"\u201d", `"`, // ” right double quote
		"\u00c6", "Ae", // Æ
		"\u00e6", "ae", // æ
		"\u00d7", "x", // × multiplication sign
		"\u2013", "-", // – en dash
		"\u2014", "-", // — em dash
		"\u2212", "-", // − minus sign
		"\u2026", "...", // … ellipsis
	).Replace(s)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	folded, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return folded
}

// --- output ------------------------------------------------------------------

// writeCatalog writes the cards as a source catalog: a two-space-indented JSON
// array with a trailing newline, matching the existing catalog files. HTML escaping
// is off so card text keeps its literal &, <, and > rather than \u escapes.
func writeCatalog(path string, cards []catalogCard) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cards); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

package parseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestEmojiCatalogFiltersAndBody(t *testing.T) {
	ctx := context.Background()
	c, call := newTestClient(t, okJSON(`{"q":"","emojis":[],"total":0,"next_offset":null}`))
	include := false
	offset := 2
	_, err := c.EmojiSearch(ctx, "face", EmojiSearchOptions{Category: "Smileys & Emotion", Version: "17.0", MaxVersion: "18.0", IncludeSkinTones: &include, Offset: &offset})
	if err != nil {
		t.Fatal(err)
	}
	want := url.Values{"q": {"face"}, "category": {"Smileys & Emotion"}, "version": {"17.0"}, "max_version": {"18.0"}, "include_skin_tones": {"false"}, "offset": {"2"}}
	got, _ := url.ParseQuery(call.rawQry)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("params: %v", got)
	}
	zero := 0
	_, err = c.EmojiBrowse(ctx, EmojiBrowseOptions{IncludeSkinTones: &include, Offset: &zero})
	if err != nil || call.rawQry != "include_skin_tones=false&offset=0" {
		t.Fatal(err, call.rawQry)
	}
	_, err = c.EmojiCategories(ctx)
	if err != nil || call.path != "/emoji/categories" || call.rawQry != "" {
		t.Fatal(err, call)
	}
	_, err = c.Emoji(ctx, ":cat:", EmojiOptions{ShortcodeSource: "github"})
	if err != nil || call.rawQry != "shortcode_source=github" {
		t.Fatal(err, call)
	}
	c, _ = newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/emoji/parse" || r.URL.Query().Get("deep") != "true" || r.URL.Query().Get("lang") != "fr" {
			t.Fatalf("bad parse request: %v", r)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body["text"] != "x 👍🏿 👍🏿" {
			t.Fatal(err, body)
		}
		okJSON(`{"emojis":[{"emoji":"👍🏿","name":"thumbs up","shortcodes":[],"category":null,"positions":[{"start":2,"end":6},{"start":7,"end":11}],"deep":{"base":"👍","variants":[{"emoji":"👍🏿","tones":["dark"]}],"shortcode_sources":{"cldr":[],"github":[],"slack":[],"discord":[]}}}]}`)(w, r)
	})
	parsed, err := c.EmojiParse(ctx, "x 👍🏿 👍🏿", EmojiParseOptions{Deep: true, Lang: "fr"})
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Emojis) != 1 || len(parsed.Emojis[0].Positions) != 2 || parsed.Emojis[0].Positions[1].End != 11 || *parsed.Emojis[0].Deep.Base != "👍" || parsed.Emojis[0].Deep.Variants[0].Tones[0] != "dark" {
		t.Fatalf("lost parse result: %#v", parsed)
	}
}

func TestEmojiOldResponseCompatibility(t *testing.T) {
	for _, body := range []string{`{"emoji":"😀","name":"face","shortcodes":[]}`, `{"emoji":"😀","name":"face","shortcodes":[],"deep":{}}`, `{"emoji":"😀","name":"face","shortcodes":[],"deep":{"base":null,"variants":null,"shortcode_sources":null}}`} {
		var out Emoji
		if err := json.Unmarshal([]byte(body), &out); err != nil {
			t.Fatal(err)
		}
		if out.Deep != nil && (out.Deep.Base != nil || out.Deep.Variants != nil) {
			t.Fatal("changed old nullable data")
		}
	}
}

func TestEmojiPrivateConversion(t *testing.T) {
	c, call := newTestClient(t, okJSON(`{"text":"🚀","unresolved":[]}`))
	result, err := c.EmojiConvert(context.Background(), ":rocket:", "to_emoji", "github")
	if err != nil || result.Text != "🚀" || call.path != "/emoji/convert" {
		t.Fatal(err, result, call)
	}
}

func TestEmojiDownloadMetadataAndArtwork(t *testing.T) {
	body, err := os.ReadFile("testdata/emoji-catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	var bundle EmojiCatalog
	if err := json.Unmarshal(body, &bundle); err != nil {
		t.Fatal(err)
	}
	c, call := newTestClient(t, okJSON(string(body)))
	result, err := c.EmojiCatalog(context.Background())
	if err != nil || call.path != "/emoji/catalog" || call.rawQry != "deep=true" {
		t.Fatal(err, call)
	}
	no := false
	_, err = c.EmojiCatalog(context.Background(), EmojiCatalogOptions{Deep: &no, IncludeSkinTones: &no, MaxVersion: "17.0"})
	if err != nil || call.rawQry != "deep=false&include_skin_tones=false&max_version=17.0" {
		t.Fatal(err, call)
	}
	if result.Forms["❤"] != "❤️" || result.Shortcodes["github"]["heart"][0] != "❤️" || len(result.Emojis[0].Deep.Artwork) != 0 {
		t.Fatalf("lost export metadata: %#v", result)
	}
	if result.Scope == nil || result.Total != 1 || !result.Scope.CompleteInventory {
		t.Fatalf("lost export scope: %#v", result.Scope)
	}
	if result.Search == nil || result.Localization.ResolvedLocale == nil || *result.Localization.ResolvedLocale != "fr" || result.Search.Native["❤️"].Keywords[0] != "love" {
		t.Fatalf("lost native search facts: %#v", result.Search)
	}
	metadataJSON, _ := json.Marshal(bundle.Metadata)
	c, call = newTestClient(t, okJSON(string(metadataJSON)))
	metadata, err := c.EmojiMetadata(context.Background())
	if err != nil || call.path != "/emoji/metadata" || call.rawQry != "" || metadata.Total != 1 || metadata.Sources[0].SourceURL != "https://unicode.org/emoji" || metadata.Notices[0].Notice != "Full source notice" {
		t.Fatal(err, metadata)
	}
	if _, err = c.EmojiCatalog(context.Background(), EmojiCatalogOptions{}, EmojiCatalogOptions{}); err == nil {
		t.Fatal("multiple catalog options accepted")
	}
}

func TestEmojiActualArtwork(t *testing.T) {
	body, err := os.ReadFile("testdata/emoji-artwork.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Metadata EmojiMetadata `json:"metadata"`
		Emoji    Emoji         `json:"emoji"`
	}
	if err = json.Unmarshal(body, &fixture); err != nil {
		t.Fatal(err)
	}
	if len(fixture.Metadata.ArtworkSources) != 3 || len(fixture.Emoji.Deep.Artwork) != 3 {
		t.Fatal("lost rendering facts")
	}
	if *fixture.Emoji.Deep.Artwork["openmoji"].SVG != "https://cdn.parseapi.com/emoji/openmoji/17.0.0/1f680.svg" || fixture.Emoji.Deep.Artwork["openmoji"].PNG["72"] != "https://cdn.parseapi.com/emoji/openmoji/17.0.0/72/1f680.png" {
		t.Fatal("lost keyed artwork URLs")
	}
	if len(fixture.Metadata.ArtworkSources[0].Notices) == 0 {
		t.Fatal("lost licence notice links")
	}
}

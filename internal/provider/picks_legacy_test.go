package provider

// A provider saved before PickedFrom was kept: its picks are taken as they
// are. Read as "everything released after the newest pick", a single pick
// of an older model served every newer one, and the next Save wrote them
// all down as the user's picks — a relay's 3 models all picked after one
// was clicked, and "no matter how I save, it saves every model".

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/yetone/magpie/internal/catalog"
)

func TestLegacyPicksServeWhatWasPicked(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, ".cache"))

	const chat = "http://127.0.0.1:1/v1"
	var ms []catalog.Model
	for i := range 10 {
		ms = append(ms, catalog.Model{ID: fmt.Sprintf("m%d", i), Released: fmt.Sprintf("2026-%02d-01", i+1)})
	}
	if err := catalog.SaveLive("relay", chat, ms); err != nil {
		t.Fatal(err)
	}
	// the state of a provider saved before PickedFrom: one old pick, no
	// PickedFrom, the list fetched since
	dir := filepath.Join(home, ".config", "magpie")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(map[string]any{"providers": []map[string]any{{
		"id": "relay", "name": "Relay", "key": "k", "chat": chat,
		"models": []string{"m0"},
	}}})
	if err := os.WriteFile(filepath.Join(dir, "providers.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Find("relay")
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Exposed(); len(got) != 1 || got[0].ID != "m0" {
		t.Fatalf("picked m0 alone, %d models served: %v", len(got), got)
	}
	if got := p.Picks(); len(got) != 1 || got[0] != "m0" {
		t.Fatalf("picked m0 alone, picks = %v", got)
	}
}

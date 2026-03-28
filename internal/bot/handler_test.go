package bot

import (
	"strings"
	"testing"

	"grimoire/internal/domain/player"
)

type fakeRepo struct{}

func (fakeRepo) SavePlayer(p *player.Player) error { return nil }

func (fakeRepo) LoadPlayers(names []string) (map[string]*player.Player, error) {
	m := make(map[string]*player.Player)
	for _, n := range names {
		m[n] = player.New(n)
	}
	return m, nil
}

func TestRenderTableFocusLine(t *testing.T) {
	names := []string{"Ada"}
	stats := map[string]*player.Player{"Ada": player.New("Ada")}
	b := NewGrimoireBot(names, stats, fakeRepo{})

	outEmpty := b.RenderTable("")
	if strings.Contains(outEmpty, "Jogador Selecionado:") {
		t.Fatal("expected no focus line when focus empty")
	}

	outFocus := b.RenderTable("Ada")
	if !strings.Contains(outFocus, "Jogador Selecionado: Ada") {
		t.Fatalf("expected focus line in output: %q", outFocus)
	}
}

func TestParseModalCustomID(t *testing.T) {
	id, stats, ok := parseModalCustomID("modal_data:123456789")
	if !ok || !stats || id != "123456789" {
		t.Fatalf("modal_data: got id=%q stats=%v ok=%v", id, stats, ok)
	}
	id2, stats2, ok2 := parseModalCustomID("modal_custom:999")
	if !ok2 || stats2 || id2 != "999" {
		t.Fatalf("modal_custom: got id=%q stats=%v ok=%v", id2, stats2, ok2)
	}
	_, _, ok3 := parseModalCustomID("other")
	if ok3 {
		t.Fatal("expected false for unknown prefix")
	}
}

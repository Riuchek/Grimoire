package player

import "testing"

func TestPlayerUpdateStats(t *testing.T) {
	p := New("x")
	p.UpdateStats(10, 9, 8, 7)
	if p.DanoTotal() != 10 || p.DanoMax() != 9 || p.CuraTotal() != 8 || p.CuraMax() != 7 {
		t.Fatalf("UpdateStats: got dano=%d/%d cura=%d/%d", p.DanoTotal(), p.DanoMax(), p.CuraTotal(), p.CuraMax())
	}
}

func TestPlayerLoadStats(t *testing.T) {
	p := New("y")
	p.LoadStats(1, 2, 3, 4, 5, 6, 7, 8, "note")
	if p.SucessoCritico() != 1 || p.FalhaCritica() != 2 || p.DanoTotal() != 3 || p.Mortes() != 8 || p.Custom() != "note" {
		t.Fatal("LoadStats did not restore fields")
	}
}

func TestPlayerIncrements(t *testing.T) {
	p := New("z")
	p.AddNat20()
	p.AddNat20()
	p.AddNat1()
	p.AddQueda()
	p.AddMorte()
	if p.SucessoCritico() != 2 || p.FalhaCritica() != 1 || p.Quedas() != 1 || p.Mortes() != 1 {
		t.Fatalf("increments: n20=%d n1=%d q=%d m=%d", p.SucessoCritico(), p.FalhaCritica(), p.Quedas(), p.Mortes())
	}
}

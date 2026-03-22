package status

type Player struct {
	name      string
	nat20     int
	nat1      int
	danoTotal int
	danoMax   int
	curaTotal int
	curaMax   int
	quedas    int
	mortes    int
	custom    string
}

func NewPlayer(name string) *Player {
	return &Player{name: name}
}

func (p *Player) Name() string   { return p.name }
func (p *Player) Nat20() int     { return p.nat20 }
func (p *Player) Nat1() int      { return p.nat1 }
func (p *Player) DanoTotal() int { return p.danoTotal }
func (p *Player) DanoMax() int   { return p.danoMax }
func (p *Player) CuraTotal() int { return p.curaTotal }
func (p *Player) CuraMax() int   { return p.curaMax }
func (p *Player) Quedas() int    { return p.quedas }
func (p *Player) Mortes() int    { return p.mortes }
func (p *Player) Custom() string { return p.custom }

func (p *Player) AddNat20() { p.nat20++ }
func (p *Player) AddNat1()  { p.nat1++ }
func (p *Player) AddQueda() { p.quedas++ }
func (p *Player) AddMorte() { p.mortes++ }

func (p *Player) UpdateStats(dTotal, dMax, cTotal, cMax int) {
	p.danoTotal = dTotal
	p.danoMax = dMax
	p.curaTotal = cTotal
	p.curaMax = cMax
}

func (p *Player) SetCustom(text string) {
	p.custom = text
}

func (p *Player) LoadStats(n20, n1, dt, dm, ct, cm, q, m int, custom string) {
	p.nat20 = n20
	p.nat1 = n1
	p.danoTotal = dt
	p.danoMax = dm
	p.curaTotal = ct
	p.curaMax = cm
	p.quedas = q
	p.mortes = m
	p.custom = custom
}

package player

type Repository interface {
	SavePlayer(p *Player) error
	LoadPlayers(names []string) (map[string]*Player, error)
}

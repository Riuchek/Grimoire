package bot

import (
	"log/slog"
	"time"

	"grimoire/internal/domain/player"
)

var _ player.Repository = (*LoggingPlayerRepository)(nil)

type LoggingPlayerRepository struct {
	Inner player.Repository
}

func (l *LoggingPlayerRepository) SavePlayer(p *player.Player) error {
	start := time.Now()
	err := l.Inner.SavePlayer(p)
	ms := time.Since(start).Milliseconds()
	if err != nil {
		slog.Warn("db SavePlayer", "player", p.Name(), "ms", ms, "err", err)
	} else {
		slog.Info("db SavePlayer", "player", p.Name(), "ms", ms)
	}
	return err
}

func (l *LoggingPlayerRepository) LoadPlayers(names []string) (map[string]*player.Player, error) {
	start := time.Now()
	m, err := l.Inner.LoadPlayers(names)
	ms := time.Since(start).Milliseconds()
	if err != nil {
		slog.Warn("db LoadPlayers", "names", len(names), "ms", ms, "err", err)
	} else {
		slog.Info("db LoadPlayers", "names", len(names), "ms", ms)
	}
	return m, err
}

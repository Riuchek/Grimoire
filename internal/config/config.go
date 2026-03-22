package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Token  string
	DBPath string
	Names  []string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		Token:  strings.TrimSpace(os.Getenv("DISCORD_TOKEN")),
		DBPath: dbPath(),
		Names:  playerNames(),
	}
}

func dbPath() string {
	if p := strings.TrimSpace(os.Getenv("GRIMOIRE_DB_PATH")); p != "" {
		return p
	}
	return "./grimoire.db"
}

func defaultPlayerNames() []string {
	return []string{
		"Gustavo", "Mariana", "Pedro", "Joao", "Janis", "Catti", "Maria", "Eric", "Andre",
	}
}

func playerNames() []string {
	raw := os.Getenv("GRIMOIRE_PLAYERS")
	if strings.TrimSpace(raw) == "" {
		return defaultPlayerNames()
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return defaultPlayerNames()
	}
	return out
}

package app

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/config"
)

func NewDatabase(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := strings.Builder{}
	dsn.WriteString(fmt.Sprintf("%s:%s", cfg.Username, cfg.Password))
	dsn.WriteString(fmt.Sprintf("@tcp(%s:%d)", cfg.Host, cfg.Port))
	dsn.WriteString("/" + cfg.Name)
	if len(cfg.Query) > 0 {
		dsn.WriteString("?" + strings.Join(cfg.Query, "&"))
	}

	return sqlx.Open(cfg.Type, dsn.String())
}

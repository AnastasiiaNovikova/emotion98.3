package db

import (
	"cfg"
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type envConfig struct {
	Adapter  string
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

func (c envConfig) ConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

func getEnvConfig() (*envConfig, error) {
	var configs map[string]envConfig
	if err := cfg.GetYamlConfig("db", &configs); err != nil {
		return nil, fmt.Errorf("can't get db config: %s", err)
	}

	c, ok := configs[cfg.GetEnv()]
	if !ok {
		return nil, fmt.Errorf("no current env %q in db config %+v",
			cfg.GetEnv(), configs)
	}

	return &c, nil
}

func init() {
	if err := initDB(); err != nil {
		log.Fatalf("can't init DB: %s", err)
	}
}

func initDB() error {
	cfg, err := getEnvConfig()
	if err != nil {
		return fmt.Errorf("can't get db env config: %s", err)
	}

	db, err = gorm.Open(cfg.Adapter, cfg.ConnString())
	if err != nil {
		return fmt.Errorf("can't open gorm connection for cfg %+v: %s", cfg, err)
	}
	return nil
}

func Get() *gorm.DB {
	return db
}

func Int64FK(v int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: v,
		Valid: v != 0,
	}
}

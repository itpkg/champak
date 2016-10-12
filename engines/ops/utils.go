package ops

import (
	"fmt"
	"path"

	"github.com/spf13/viper"

	"bitbucket.org/liamstask/goose/lib/goose"
)

func dbMigrationsDir() string {
	return path.Join("db", "migrations")
}

func dbConf() *goose.DBConf {
	args := ""
	for k, v := range viper.GetStringMapString("database.args") {
		args += fmt.Sprintf(" %s=%s ", k, v)
	}
	return &goose.DBConf{
		MigrationsDir: dbMigrationsDir(),
		Env:           viper.GetString("env"),
		Driver: goose.DBDriver{
			Name:    viper.GetString("database.driver"),
			OpenStr: args,
			Import:  "github.com/lib/pq",
			Dialect: &goose.PostgresDialect{},
		},
		PgSchema: "",
	}
}

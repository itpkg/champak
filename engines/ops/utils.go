package ops

import (
	"errors"
	"path"
	"strings"
)

func dbMigrationsDir() string {
	return path.Join("db", "migrations")
}

func errors2(ers []error) error {
	var msg []string
	for _, er := range ers {
		msg = append(msg, er.Error())
	}
	return errors.New(strings.Join(msg, "\n"))
}

package auth

import (
	"github.com/facebookgo/inject"
	"github.com/itpkg/champak/web"
	"github.com/itpkg/champak/web/i18n"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

//IocAction ioc action
func IocAction(fn func(*cli.Context, *inject.Graph) error) cli.ActionFunc {
	return Action(func(ctx *cli.Context) error {
		var inj inject.Graph
		logger := OpenLogger()
		if !IsProduction() {
			inj.Logger = logger
		}

		db, err := OpenDatabase()
		if err != nil {
			return err
		}
		rep := OpenRedis()

		i1n := i18n.I18n{Locales: make(map[string]map[string]string)}
		if err := inj.Provide(
			&inject.Object{Value: logger},
			&inject.Object{Value: db},
			&inject.Object{Value: rep},
			&inject.Object{Value: &i18n.GormStore{}},
			&inject.Object{Value: &i1n},
		); err != nil {
			return err
		}
		web.Loop(func(en web.Engine) error {
			if e := en.Map(&inj); e != nil {
				return e
			}
			return inj.Provide(&inject.Object{Value: en})
		})
		if err := inj.Populate(); err != nil {
			return err
		}
		if err := i1n.Load("locales"); err != nil {
			return err
		}
		return fn(ctx, &inj)
	})
}

//Action load config action
func Action(f cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		return f(c)
	}
}

package auth

import (
	"fmt"
	"net/http"

	"github.com/facebookgo/inject"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/itpkg/champak/web"
	"github.com/itpkg/champak/web/i18n"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

//Shell command line
func (p *Engine) Shell() []cli.Command {
	return []cli.Command{

		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the app server",
			Action: IocAction(func(*cli.Context, *inject.Graph) error {
				if IsProduction() {
					gin.SetMode(gin.ReleaseMode)
				}
				rt := gin.Default()
				// rt.LoadHTMLGlob(fmt.Sprintf("themes/%s/**/*", viper.GetString("server.theme")))
				rt.Use(i18n.LocaleHandler)

				web.Loop(func(en web.Engine) error {
					en.Mount(rt)
					return nil
				})

				adr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
				hnd := cors.New(cors.Options{
					AllowCredentials: true,
					AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
					AllowedHeaders:   []string{"*"},
					Debug:            !IsProduction(),
				}).Handler(rt)

				if IsProduction() {
					return endless.ListenAndServe(adr, hnd)
				}
				return http.ListenAndServe(adr, hnd)
			}),
		},
		{
			Name:    "worker",
			Aliases: []string{"w"},
			Usage:   "start the worker progress",
			Action: IocAction(func(_ *cli.Context, inj *inject.Graph) error {
				web.Loop(func(en web.Engine) error {
					en.Worker()
					return nil
				})

				return p.Jobber.Start()
			}),
		},
		{
			Name:    "redis",
			Aliases: []string{"re"},
			Usage:   "open redis connection",
			Action: Action(func(*cli.Context) error {
				return web.Shell(
					"redis-cli",
					"-h", viper.GetString("redis.host"),
					"-p", viper.GetString("redis.port"),
					"-n", viper.GetString("redis.db"),
				)
			}),
		},
		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "database operations",
			Subcommands: []cli.Command{
				{
					Name:    "example",
					Usage:   "scripts example for create database and user",
					Aliases: []string{"e"},
					Action: Action(func(*cli.Context) error {
						drv := viper.GetString("database.driver")
						args := viper.GetStringMapString("database.args")
						var err error
						switch drv {
						case "postgres":
							fmt.Printf("CREATE USER %s WITH PASSWORD '%s';\n", args["user"], args["password"])
							fmt.Printf("CREATE DATABASE %s WITH ENCODING='UTF8';\n", args["dbname"])
							fmt.Printf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;\n", args["dbname"], args["user"])
						default:
							err = fmt.Errorf("unknown driver %s", drv)
						}
						return err
					}),
				},
				{
					Name:    "migrate",
					Usage:   "migrate the database",
					Aliases: []string{"m"},
					Action: Action(func(*cli.Context) error {
						//TODO
						return nil
					}),
				},
				{
					Name:    "rollback",
					Usage:   "rollback the database",
					Aliases: []string{"m"},
					Action: Action(func(*cli.Context) error {
						//TODO
						return nil
					}),
				},
				{
					Name:    "connect",
					Usage:   "connect database",
					Aliases: []string{"c"},
					Action: Action(func(*cli.Context) error {
						drv := viper.GetString("database.driver")
						args := viper.GetStringMapString("database.args")
						var err error
						switch drv {
						case "postgres":
							err = web.Shell("psql",
								"-h", args["host"],
								"-p", args["port"],
								"-U", args["user"],
								args["dbname"],
							)
						default:
							err = fmt.Errorf("unknown driver %s", drv)
						}
						return err
					}),
				},
				{
					Name:    "create",
					Usage:   "create database",
					Aliases: []string{"n"},
					Action: Action(func(*cli.Context) error {
						drv := viper.GetString("database.driver")
						args := viper.GetStringMapString("database.args")
						var err error
						switch drv {
						case "postgres":
							err = web.Shell("psql",
								"-h", args["host"],
								"-p", args["port"],
								"-U", "postgres",
								"-c", fmt.Sprintf(
									"CREATE DATABASE %s WITH ENCODING='UTF8'",
									args["dbname"],
								),
							)
						default:
							err = fmt.Errorf("unknown driver %s", drv)
						}
						return err
					}),
				},
				{
					Name:    "drop",
					Usage:   "drop database",
					Aliases: []string{"d"},
					Action: Action(func(*cli.Context) error {
						drv := viper.GetString("database.driver")
						args := viper.GetStringMapString("database.args")
						var err error
						switch drv {
						case "postgres":
							err = web.Shell("psql",
								"-h", args["host"],
								"-p", args["port"],
								"-U", "postgres",
								"-c", fmt.Sprintf("DROP DATABASE %s", args["dbname"]),
							)
						default:
							err = fmt.Errorf("unknown driver %s", drv)
						}
						return err
					}),
				},
			},
		},
		{
			Name:    "cache",
			Aliases: []string{"c"},
			Usage:   "cache operations",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Usage:   "list all cache keys",
					Aliases: []string{"l"},
					Action: IocAction(func(*cli.Context, *inject.Graph) error {
						keys, err := p.Cache.Keys()
						if err != nil {
							return err
						}
						for _, k := range keys {
							fmt.Println(k)
						}
						return nil
					}),
				},
				{
					Name:    "clear",
					Usage:   "clear cache items",
					Aliases: []string{"c"},
					Action: IocAction(func(*cli.Context, *inject.Graph) error {
						return p.Cache.Flush()
					}),
				},
			},
		},
	}
}

func init() {
	viper.SetEnvPrefix("champak")
	viper.BindEnv("env")
	viper.SetDefault("env", "development")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetDefault("redis", map[string]interface{}{
		"host": "localhost",
		"port": 6379,
		"db":   8,
	})

	viper.SetDefault("home", map[string]interface{}{
		"backend": "http://localhost:8080",
		"front":   "http://localhost:4200",
	})

	viper.SetDefault("database", map[string]interface{}{
		"driver": "postgres",
		"args": map[string]interface{}{
			"host":    "localhost",
			"port":    5432,
			"user":    "postgres",
			"dbname":  "champak_dev",
			"sslmode": "disable",
		},
		"pool": map[string]int{
			"max_open": 180,
			"max_idle": 6,
		},
	})

	viper.SetDefault("server", map[string]interface{}{
		"port":  8080,
		"name":  "www.change-me.com",
		"theme": "bootstrap4",
	})
	viper.SetDefault("secrets", map[string]interface{}{
		"jwt":  web.RandomStr(32),
		"aes":  web.RandomStr(32),
		"hmac": web.RandomStr(32),
	})

	viper.SetDefault("workers", map[string]interface{}{
		"timeout": 30,
	})

	viper.SetDefault("elasticsearch", []string{"http://localhost:9200"})
}

package auth

import (
	"crypto/x509/pkix"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"

	"github.com/BurntSushi/toml"
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
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init config file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "environment, e",
					Value: "development",
					Usage: "environment, like: development, production, stage, test...",
				},
			},
			Action: func(c *cli.Context) error {
				const fn = "config.toml"
				if _, err := os.Stat(fn); err == nil {
					return fmt.Errorf("file %s already exists", fn)
				}
				fmt.Printf("generate file %s\n", fn)

				viper.Set("env", c.String("environment"))
				args := viper.AllSettings()
				fd, err := os.Create(fn)
				if err != nil {
					return err
				}
				defer fd.Close()
				end := toml.NewEncoder(fd)
				err = end.Encode(args)

				return err

			},
		},
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
			Name:    "openssl",
			Aliases: []string{"ssl"},
			Usage:   "generate ssl certificates",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "name",
				},
				cli.StringFlag{
					Name:  "country, c",
					Value: "Earth",
					Usage: "country",
				},
				cli.StringFlag{
					Name:  "organization, o",
					Value: "Mother Nature",
					Usage: "organization",
				},
				cli.IntFlag{
					Name:  "years, y",
					Value: 1,
					Usage: "years",
				},
			},
			Action: Action(func(c *cli.Context) error {
				name := c.String("name")
				if len(name) == 0 {
					cli.ShowCommandHelp(c, "openssl")
					return nil
				}
				root := path.Join("etc", "ssl", name)

				key, crt, err := CreateCertificate(
					true,
					pkix.Name{
						Country:      []string{c.String("country")},
						Organization: []string{c.String("organization")},
					},
					c.Int("years"),
				)
				if err != nil {
					return err
				}

				fnk := path.Join(root, "key.pem")
				fnc := path.Join(root, "crt.pem")

				fmt.Printf("generate pem file %s\n", fnk)
				err = WritePemFile(fnk, "RSA PRIVATE KEY", key)
				fmt.Printf("test: openssl rsa -noout -text -in %s\n", fnk)

				if err == nil {
					fmt.Printf("generate pem file %s\n", fnc)
					err = WritePemFile(fnc, "CERTIFICATE", crt)
					fmt.Printf("test: openssl x509 -noout -text -in %s\n", fnc)
				}
				if err == nil {
					fmt.Printf(
						"verify: diff <(openssl rsa -noout -modulus -in %s) <(openssl x509 -noout -modulus -in %s)",
						fnk,
						fnc,
					)
				}
				fmt.Println()
				return err
			}),
		},

		{
			Name:    "nginx",
			Aliases: []string{"ng"},
			Usage:   "init nginx config file",
			Action: Action(func(*cli.Context) error {
				const tpl = `
server {
  listen 80;
  server_name {{.Name}};
  rewrite ^(.*) https://$host$1 permanent;
}

upstream {{.Name}}_prod {
  server localhost:{{.Port}} fail_timeout=0;
}

server {
  listen 443;

  ssl  on;
  ssl_certificate  /etc/ssl/certs/{{.Name}}.crt;
  ssl_certificate_key  /etc/ssl/private/{{.Name}}.key;
  ssl_session_timeout  5m;
  ssl_protocols  SSLv2 SSLv3 TLSv1;
  ssl_ciphers  RC4:HIGH:!aNULL:!MD5;
  ssl_prefer_server_ciphers  on;

  client_max_body_size 4G;
  keepalive_timeout 10;
  proxy_buffers 16 64k;
  proxy_buffer_size 128k;

  server_name {{.Name}};
  root {{.Root}}/public;
  index index.html;
  access_log /var/log/nginx/{{.Name}}.access.log;
  error_log /var/log/nginx/{{.Name}}.error.log;
  location / {
    try_files $uri $uri/ /index.html?/$request_uri;
  }
#  location ^~ /assets/ {
#    gzip_static on;
#    expires max;
#    access_log off;
#    add_header Cache-Control "public";
#  }
  location ~* \.(?:css|js)$ {
    gzip_static on;
    expires max;
    access_log off;
    add_header Cache-Control "public";
  }
  location ~* \.(?:jpg|jpeg|gif|png|ico|cur|gz|svg|svgz|mp4|ogg|ogv|webm|htc)$ {
    expires 1M;
    access_log off;
    add_header Cache-Control "public";
  }
  location ~* \.(?:rss|atom)$ {
    expires 12h;
    access_log off;
    add_header Cache-Control "public";
  }
  location ~ ^/api/{{.Version}}(/?)(.*) {
    proxy_set_header X-Forwarded-Proto https;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_redirect off;
    proxy_pass http://{{.Name}}_prod/$2$is_args$args;
    # limit_req zone=one;
  }
}
`
				t, err := template.New("").Parse(tpl)
				if err != nil {
					return err
				}
				pwd, err := os.Getwd()
				if err != nil {
					return err
				}

				name := viper.GetString("server.name")
				fn := path.Join("etc", "nginx", "sites-enabled", name+".conf")
				if err = os.MkdirAll(path.Dir(fn), 0700); err != nil {
					return err
				}
				fmt.Printf("generate file %s\n", fn)
				fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
				if err != nil {
					return err
				}
				defer fd.Close()

				return t.Execute(fd, struct {
					Name    string
					Port    int
					Root    string
					Version string
				}{
					Name:    name,
					Port:    viper.GetInt("http.port"),
					Root:    pwd,
					Version: "v1",
				})
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
								"-U", args["user"],
								"-c", fmt.Sprintf("create database %s WITH ENCODING='UTF8'", args["dbname"]),
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
								"-U", args["user"],
								"-c", fmt.Sprintf("drop database %s", args["dbname"]),
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
		"jwt": web.RandomStr(32),
		"aes": web.RandomStr(32),
	})

	viper.SetDefault("workers", map[string]interface{}{
		"timeout": 30,
	})

	viper.SetDefault("elasticsearch", []string{"http://localhost:9200"})
}

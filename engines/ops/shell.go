package ops

import (
	"crypto/x509/pkix"
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"path"
	"time"

	"golang.org/x/text/language"

	"bitbucket.org/liamstask/goose/lib/goose"

	"github.com/BurntSushi/toml"
	"github.com/itpkg/champak/engines/auth"
	"github.com/itpkg/champak/web"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

//Shell command options
func (p *Engine) Shell() []cli.Command {
	return []cli.Command{

		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "database operations",
			Subcommands: []cli.Command{
				{
					Name:    "example",
					Usage:   "scripts example for create database and user",
					Aliases: []string{"e"},
					Action: auth.Action(func(*cli.Context) error {
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
					Action: auth.Action(func(*cli.Context) error {
						ver, err := goose.GetMostRecentDBVersion(dbMigrationsDir())
						if err != nil {
							return err
						}
						return goose.RunMigrations(dbConf(), dbMigrationsDir(), ver)
					}),
				},
				{
					Name:    "rollback",
					Usage:   "rollback the database",
					Aliases: []string{"r"},
					Action: auth.Action(func(*cli.Context) error {
						cnf := dbConf()
						crt, err := goose.GetDBVersion(cnf)
						if err != nil {
							return err
						}
						ver, err := goose.GetPreviousDBVersion(dbMigrationsDir(), crt)
						if err != nil {
							return err
						}
						return goose.RunMigrations(cnf, dbMigrationsDir(), ver)
					}),
				},
				{
					Name:    "version",
					Usage:   "show database scheme version",
					Aliases: []string{"v"},
					Action: auth.Action(func(*cli.Context) error {
						cnf := dbConf()
						migs, err := goose.CollectMigrations(dbMigrationsDir(), 0, ((1 << 63) - 1))
						if err != nil {
							return err
						}
						db, err := goose.OpenDBFromDBConf(cnf)
						if err != nil {
							return err
						}
						defer db.Close()
						if _, err = goose.EnsureDBVersion(cnf, db); err != nil {
							return err
						}
						fmt.Println("    Applied At                  Migration")
						for _, mig := range migs {
							var row goose.MigrationRecord
							q := fmt.Sprintf("SELECT tstamp, is_applied FROM goose_db_version WHERE version_id=%d ORDER BY tstamp DESC LIMIT 1", mig.Version)
							e := db.QueryRow(q).Scan(&row.TStamp, &row.IsApplied)

							if e != nil && e != sql.ErrNoRows {
								return e
							}

							var appliedAt string

							if row.IsApplied {
								appliedAt = row.TStamp.Format(time.ANSIC)
							} else {
								appliedAt = "Pending"
							}

							fmt.Printf("    %-24s -- %v\n", appliedAt, mig.Source)
						}
						return nil
					}),
				},
				{
					Name:    "connect",
					Usage:   "connect database",
					Aliases: []string{"c"},
					Action: auth.Action(func(*cli.Context) error {
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
					Action: auth.Action(func(*cli.Context) error {
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
					Action: auth.Action(func(*cli.Context) error {
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
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate file template",
			Subcommands: []cli.Command{
				{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "generate config file",
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
					Name:    "nginx",
					Aliases: []string{"ng"},
					Usage:   "generate nginx.conf",
					Action: auth.Action(func(*cli.Context) error {
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
					Action: auth.Action(func(c *cli.Context) error {
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
						err = WritePemFile(fnk, "RSA PRIVATE KEY", key, 0600)
						fmt.Printf("test: openssl rsa -noout -text -in %s\n", fnk)

						if err == nil {
							fmt.Printf("generate pem file %s\n", fnc)
							err = WritePemFile(fnc, "CERTIFICATE", crt, 0444)
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
					Name:    "migration",
					Usage:   "generate migration file",
					Aliases: []string{"m"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name",
						},
					},
					Action: auth.Action(func(c *cli.Context) error {
						name := c.String("name")
						if len(name) == 0 {
							cli.ShowCommandHelp(c, "migration")
							return nil
						}
						root := dbMigrationsDir()
						if err := os.MkdirAll(root, 0700); err != nil {
							return err
						}
						pth, err := goose.CreateMigration(
							c.String("name"),
							"sql",
							root,
							time.Now(),
						)
						fmt.Printf("generate file %s\n", pth)
						return err
					}),
				},

				{
					Name:    "locale",
					Usage:   "generate locale file",
					Aliases: []string{"l"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "locale name",
						},
					},
					Action: auth.Action(func(c *cli.Context) error {
						name := c.String("name")
						if len(name) == 0 {
							cli.ShowCommandHelp(c, "locale")
							return nil
						}
						lng, err := language.Parse(name)
						if err != nil {
							return err
						}
						const root = "locales"
						if err = os.MkdirAll(root, 0700); err != nil {
							return err
						}
						file := path.Join(root, fmt.Sprintf("%s.ini", lng.String()))
						fmt.Printf("generate file %s\n", file)
						fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
						if err != nil {
							return err
						}
						defer fd.Close()
						return err
					}),
				},
			},
		},
	}
}

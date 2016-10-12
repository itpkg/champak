package auth

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

//IsProduction is production mode?
func IsProduction() bool {
	return viper.GetString("env") == "production"
}

//OpenDatabase open database
func OpenDatabase() (*gorm.DB, error) {
	//postgresql: "user=%s password=%s host=%s port=%d dbname=%s sslmode=%s"
	args := ""
	for k, v := range viper.GetStringMapString("database.args") {
		args += fmt.Sprintf(" %s=%s ", k, v)
	}
	db, err := gorm.Open(viper.GetString("database.driver"), args)
	if err != nil {
		return nil, err
	}
	if !IsProduction() {
		db.LogMode(true)
	}

	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(viper.GetInt("database.pool.max_idle"))
	db.DB().SetMaxOpenConns(viper.GetInt("database.pool.max_open"))
	return db, nil

}

//OpenRedis get redis connection pool
func OpenRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial(
				"tcp",
				fmt.Sprintf(
					"%s:%d",
					viper.GetString("redis.host"),
					viper.GetInt("redis.port"),
				),
			)
			if e != nil {
				return nil, e
			}
			if _, e = c.Do("SELECT", viper.GetInt("redis.db")); e != nil {
				c.Close()
				return nil, e
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

//OpenLogger open logger
func OpenLogger() *logging.Logger {
	var bkd logging.Backend
	if IsProduction() {
		var err error
		bkd, err = logging.NewSyslogBackend("champak")
		if err != nil {
			bkd = logging.NewLogBackend(os.Stdout, "", 0)
		}
	} else {
		bkd = logging.NewLogBackend(os.Stdout, "", 0)
	}

	//`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`
	if IsProduction() {
		logging.SetFormatter(logging.MustStringFormatter(`%{color}%{level:.4s} %{id:03x} %{color:reset} [%{shortfunc}] %{message}`))
		logging.SetLevel(logging.INFO, "")
	} else {
		logging.SetFormatter(logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{level:.4s} %{id:03x} %{color:reset} [%{longfunc}] %{message}`))
	}
	logging.SetBackend(bkd)
	return logging.MustGetLogger("backend")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

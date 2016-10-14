package auth

import (
	"crypto/aes"

	"github.com/SermoDigital/jose/crypto"
	"github.com/facebookgo/inject"
	"github.com/itpkg/champak/web"
	"github.com/itpkg/champak/web/cache"
	"github.com/itpkg/champak/web/i18n"
	"github.com/itpkg/champak/web/jobber"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

//Engine engine model
type Engine struct {
	Cache     cache.Store     `inject:""`
	Dao       *Dao            `inject:""`
	Db        *gorm.DB        `inject:""`
	Jobber    jobber.Jobber   `inject:""`
	Logger    *logging.Logger `inject:""`
	Encryptor *web.Encryptor  `inject:""`
	I18n      *i18n.I18n      `inject:""`
}

//Map map objects
func (p *Engine) Map(inj *inject.Graph) error {
	cip, err := aes.NewCipher([]byte(viper.GetString("secrets.aes")))
	if err != nil {
		return err
	}

	return inj.Provide(
		&inject.Object{Value: cip},
		&inject.Object{Value: cip, Name: "aes.cip"},
		&inject.Object{Value: []byte(viper.GetString("secrets.hmac")), Name: "hmac.key"},
		&inject.Object{Value: []byte(viper.GetString("secrets.jwt")), Name: "jwt.key"},
		&inject.Object{Value: crypto.SigningMethodHS512, Name: "jwt.method"},
		&inject.Object{Value: &cache.RedisStore{}},
		&inject.Object{Value: &jobber.RedisJobber{
			Timeout:  viper.GetInt("workers.timeout"),
			Handlers: make(map[string]jobber.Handler),
		}},
	)

}

//Worker do background job
func (p *Engine) Worker() {

}

func init() {
	web.Register(&Engine{})
}

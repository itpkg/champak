package auth

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/itpkg/champak/web"
	"github.com/itpkg/champak/web/cache"
	"github.com/itpkg/champak/web/i18n"
	"github.com/itpkg/champak/web/jobber"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

//Engine engine model
type Engine struct {
	Cache     cache.Store     `inject:""`
	Db        *gorm.DB        `inject:""`
	Jobber    jobber.Jobber   `inject:""`
	Logger    *logging.Logger `inject:""`
	Encryptor *web.Encryptor  `inject:""`
	I18n      *i18n.I18n      `inject:""`
}

//Map map objects
func (p *Engine) Map(*inject.Graph) error {
	return nil
}

//Mount mount web point
func (p *Engine) Mount(*gin.Engine) {}

//Worker do background job
func (p *Engine) Worker() {

}

func init() {
	web.Register(&Engine{})
}

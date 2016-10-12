package ops

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/itpkg/champak/web"
	logging "github.com/op/go-logging"
)

//Engine engine model
type Engine struct {
	Logger *logging.Logger `inject:""`
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

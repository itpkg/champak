package reading

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/itpkg/champak/web"
	"github.com/urfave/cli"
)

//Engine engine model
type Engine struct {
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

//Shell command options
func (p *Engine) Shell() []cli.Command {
	return []cli.Command{}
}
func init() {
	web.Register(&Engine{})
}

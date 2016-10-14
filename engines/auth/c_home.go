package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itpkg/champak/web"
)

func (p *Engine) getSiteInfo(c *gin.Context) {
	lang := c.MustGet("locale").(string)
	info := make(map[string]interface{})
	info["lang"] = lang
	for _, k := range []string{"title", "subTitle", "keywords", "description", "copyright"} {
		info[k] = p.I18n.T(lang, fmt.Sprintf("site.%s", k))
	}
	author := make(map[string]string)
	for _, k := range []string{"name", "email"} {
		var v string
		if err := p.Dao.Get(fmt.Sprintf("author.%s", k), &v); err != nil {
			p.Logger.Error(err)
		}
		author[k] = v
	}
	info["author"] = author
	for _, k := range []string{"top", "bottom"} {
		var links []web.Link
		if err := p.Dao.Get(fmt.Sprintf("links.%s", k), &links); err != nil {
			p.Logger.Error(err)
			for i := 1; i <= 3; i++ {
				links = append(
					links,
					web.Link{Href: "index", Label: fmt.Sprintf("Link %d", i)},
				)
			}
		}
		info[fmt.Sprintf("%sLinks", k)] = links
	}
	c.JSON(http.StatusOK, info)
}

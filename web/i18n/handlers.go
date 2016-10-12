package i18n

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	logging "github.com/op/go-logging"
	"golang.org/x/text/language"
)

//LocaleHandler detect locale from http header
func LocaleHandler(lg *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const key = "locale"
		// 1. Check URL arguments.
		lng := c.Request.URL.Query().Get(key)

		// 2. Get language information from cookies.
		if len(lng) == 0 {
			if ck, er := c.Request.Cookie(key); er == nil {
				lng = ck.Value
			}
		}

		// 3. Get language information from 'Accept-Language'.
		if len(lng) == 0 {
			al := c.Request.Header.Get("Accept-Language")
			if len(al) > 4 {
				lng = al[:5]
			}
		}
		tag, err := language.Parse(lng)
		if err != nil {
			lg.Error(err)
			tag = language.AmericanEnglish
		}

		//Write cookie
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    key,
			Value:   tag.String(),
			Expires: time.Now().Add(7 * 24 * time.Hour),
			Path:    "/",
		})

		c.Set(key, tag.String())

	}
}

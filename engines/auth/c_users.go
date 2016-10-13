package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Engine) getUsersSignIn(c *gin.Context) {
	lang := c.MustGet("locale").(string)
	c.HTML(http.StatusOK, "users/non-sign-in", gin.H{
		"locale": lang,
		"form": gin.H{
			"title": p.I18n.T(lang, "auth.users.sign-in"),
			"fields": []gin.H{
				gin.H{"type": "email", "id": "email"},
				gin.H{"type": "password", "id": "password"},
				gin.H{"type": "password", "id": "passwordConfirm"},
			},
		},
	})
}

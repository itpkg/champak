package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Engine) getUsersSignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "users/sign-in", gin.H{})
}

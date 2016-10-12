package auth

import "github.com/gin-gonic/gin"

//Mount mount web point
func (p *Engine) Mount(rt *gin.Engine) {
	ug := rt.Group("/users")
	ug.GET("sign-in", p.getUsersSignIn)
}

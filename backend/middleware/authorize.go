package middleware

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthorizeRequest is used to authorize a request for a certain end-point group.

func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-token-gen-on-server-side")
		if v == nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "Please login."})
			c.Abort()
		}
		c.Next()
	}
}

func AuthorizeAdminRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminSession := sessions.Default(c)
		v := adminSession.Get("agent-token-gen-on-server-side")
		if v == nil {
			c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{"message": "Please login."})
			c.Abort()
		}
		c.Next()
	}
}

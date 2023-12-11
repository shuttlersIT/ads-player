package middleware

import (
	"gorm.io/gorm"

	//"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
		c.Next()
	}
}

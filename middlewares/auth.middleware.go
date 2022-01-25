package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"messagewith-server/sessions"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := sessions.GetSessionFromCookie(c)
		if err != nil {
			c.Next()
			return
		}

		user, err := sessions.GetUserFromSession(c, session)
		if err != nil {
			c.Next()
			return
		}

		ctx := context.WithValue(c.Request.Context(), "LoggedUser", user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
